package output

import (
	"fmt"
	"io"
	"log"
	"os"
	run "runtime"
	"sort"
	"time"

	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/logrusorgru/aurora"
)

var au aurora.Aurora

// OutputWriter represents the output
type OutputWriter struct {
	out   io.Writer
	color bool
	isDir bool
	order bool
}

// NewCliOutput creates a new OutputWriter with a stdout writer
func NewCliOutput(color bool, order bool) OutputWriter {
	return OutputWriter{
		out:   os.Stdout,
		color: color,
		order: order,
	}
}

// Start starts the writing sequence
func (w *OutputWriter) Start(results <-chan runtime.TestResult) bool {
	au = aurora.NewAurora(w.color)
	if run.GOOS == "windows" {
		au = aurora.NewAurora(false)
	}

	fileErrors := []string{}
	testResults := []runtime.TestResult{}

	start := time.Now()
	for r := range results {
		if r.FileError != nil {
			str := fmt.Sprintf("[%s] %s!", r.FileName, r.FileError.Error())
			fileErrors = append(fileErrors, str)
			continue
		}
		testResults = append(testResults, r)
	}
	duration := time.Since(start)

	if w.order { //maintain file order
		sort.SliceStable(testResults, func(i, j int) bool {
			return testResults[i].FileName < testResults[j].FileName
		})
	}

	failed := 0
	//Actually print the results
	for _, r := range testResults {
		if r.ValidationResult.Success {
			str := fmt.Sprintf("✓ [%s] [%s] %s", r.FileName, r.Node, r.TestCase.Title)
			s := w.addTries(str, r)
			w.fprintf(s)
			continue
		}

		if !r.ValidationResult.Success {
			failed++
			str := fmt.Sprintf("✗ [%s] [%s] %s", r.FileName, r.Node, r.TestCase.Title)
			s := w.addTries(str, r)
			w.fprintf(au.Red(s))
			continue
		}
	}

	if failed > 0 {
		w.printFailures(testResults)
	}

	w.fprintf("")
	w.fprintf(fmt.Sprintf("Duration: %.3fs", duration.Seconds()))
	summary := fmt.Sprintf("Count: %d, Failed: %d", len(testResults), failed)
	if failed > 0 {
		w.fprintf(au.Red(summary))
	} else {
		w.fprintf(au.Green(summary))
	}

	w.fprintf("")
	for _, e := range fileErrors {
		w.fprintf("File Errors:")
		w.fprintf(e)
	}
	return failed == 0
}

func (w *OutputWriter) addTries(s string, r runtime.TestResult) string {
	if r.Tries > 1 {
		s = fmt.Sprintf("%s, retries %d", s, r.Tries)
	}
	return s
}

func (w *OutputWriter) printFailures(results []runtime.TestResult) {
	w.fprintf("")
	w.fprintf(au.Bold("Results"))
	w.fprintf(au.Bold(""))

	for _, r := range results {
		if r.TestCase.Result.Error != nil {
			str := fmt.Sprintf("✗ [%s][%s] '%s' could not be executed with error message:", r.FileName, r.Node, r.TestCase.Title)
			w.fprintf(au.Bold(au.Red(str)))
			w.fprintf(r.TestCase.Result.Error.Error())
			continue
		}

		if !r.ValidationResult.Success {
			str := fmt.Sprintf("✗ [%s][%s] '%s', on property '%s'", r.FileName, r.Node, r.TestCase.Title, r.FailedProperty)
			w.fprintf(au.Bold(au.Red(str)))
			w.fprintf(r.ValidationResult.Diff)
		}
	}
}

func (w *OutputWriter) fprintf(a ...interface{}) {
	if _, err := fmt.Fprintln(w.out, a...); err != nil {
		log.Fatal(err)
	}
}
