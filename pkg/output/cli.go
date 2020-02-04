package output

import (
	"fmt"
	"io"
	"log"
	"os"
	run "runtime"
	"time"

	"github.com/SimonBaeumer/commander/pkg/runtime"
	"github.com/logrusorgru/aurora"
)

var au aurora.Aurora

// OutputWriter represents the output
type OutputWriter struct {
	out   io.Writer
	color bool
}

// NewCliOutput creates a new OutputWriter with a stdout writer
func NewCliOutput(color bool) OutputWriter {
	return OutputWriter{
		out:   os.Stdout,
		color: color,
	}
}

// Start starts the writing sequence
func (w *OutputWriter) Start(results <-chan runtime.TestResult) bool {
	au = aurora.NewAurora(w.color)
	if run.GOOS == "windows" {
		au = aurora.NewAurora(false)
	}

	failed := 0
	testResults := []runtime.TestResult{}
	start := time.Now()

	for r := range results {
		testResults = append(testResults, r)
		if r.ValidationResult.Success {
			str := fmt.Sprintf("✓ %s", r.TestCase.Title)
			s := w.addTries(str, r)
			w.fprintf(s)
		}

		if !r.ValidationResult.Success {
			failed++
			str := fmt.Sprintf("✗ %s", r.TestCase.Title)
			s := w.addTries(str, r)
			w.fprintf(au.Red(s))
		}
	}

	duration := time.Since(start)

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
			str := fmt.Sprintf("✗ '%s' could not be executed with error message:", r.TestCase.Title)
			w.fprintf(au.Bold(au.Red(str)))
			w.fprintf(r.TestCase.Result.Error.Error())
			continue
		}

		if !r.ValidationResult.Success {
			str := fmt.Sprintf("✗ '%s', on property '%s'", r.TestCase.Title, r.FailedProperty)
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
