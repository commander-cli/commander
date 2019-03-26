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
			w.fprintf("✓ " + r.TestCase.Title)
		}

		if !r.ValidationResult.Success {
			failed++
			w.fprintf(au.Red("✗ " + r.TestCase.Title))
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

func (w *OutputWriter) printFailures(results []runtime.TestResult) {
	w.fprintf("")
	w.fprintf(au.Bold("Results"))
	w.fprintf(au.Bold(""))

	for _, r := range results {
		if r.TestCase.Result.Error != nil {
			w.fprintf(au.Bold(au.Red("✗ '" + r.TestCase.Title + "' could not be executed with error message:")))
			w.fprintf(r.TestCase.Result.Error.Error())
			continue
		}

		if !r.ValidationResult.Success {
			w.fprintf(au.Bold(au.Red("✗ '" + r.TestCase.Title + "', on property '" + r.FailedProperty + "'")))
			w.fprintf(r.ValidationResult.Diff)
		}
	}
}

func (w *OutputWriter) fprintf(a ...interface{}) {
	if _, err := fmt.Fprintln(w.out, a...); err != nil {
		log.Fatal(err)
	}
}
