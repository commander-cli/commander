package output

import (
	"fmt"
	"io"
	"log"
	"os"
	run "runtime"

	"github.com/commander-cli/commander/pkg/runtime"
	"github.com/logrusorgru/aurora"
)

// OutputWriter represents the output
type OutputWriter struct {
	out      io.Writer
	au       aurora.Aurora
	template cliTemplate
}

// NewCliOutput creates a new OutputWriter with a stdout writer
func NewCliOutput(color bool) OutputWriter {
	au := aurora.NewAurora(color)
	if run.GOOS == "windows" {
		au = aurora.NewAurora(false)
	}

	t := newCliTemplate()

	return OutputWriter{
		out:      os.Stdout,
		au:       au,
		template: t,
	}
}

// TestResult for output
type TestResult struct {
	FileName       string
	Title          string
	Node           string
	Tries          int
	Success        bool
	FailedProperty string
	Diff           string
	Error          error
}

// GetEventHandler create a new runtime.EventHandler
func (w *OutputWriter) GetEventHandler() *runtime.EventHandler {
	handler := runtime.EventHandler{}
	handler.TestFinished = func(testResult runtime.TestResult) {
		tr := convertTestResult(testResult)
		w.printResult(tr)
	}
	return &handler
}

// PrintSummary prints summary
func (w *OutputWriter) PrintSummary(result runtime.Result) bool {
	if result.Failed > 0 {
		w.printFailures(result.TestResults)
	}

	w.fprintf("")
	w.fprintf(w.template.duration(result))
	summary := w.template.summary(result)
	if result.Failed > 0 {
		w.fprintf(w.au.Red(summary))
	} else {
		w.fprintf(w.au.Green(summary))
	}

	return result.Failed == 0
}

// PrintResult prints the simple output form of a TestReault
func (w *OutputWriter) printResult(r TestResult) {
	if !r.Success {
		w.fprintf(w.au.Red(w.template.testResult(r)))
		return
	}
	w.fprintf(w.template.testResult(r))
}

func (w *OutputWriter) printFailures(results []runtime.TestResult) {
	w.fprintf("")
	w.fprintf(w.au.Bold("Results"))
	w.fprintf(w.au.Bold(""))

	for _, tr := range results {
		r := convertTestResult(tr)
		if r.Error != nil {
			w.fprintf(w.au.Bold(w.au.Red(w.template.errors(r))))
			w.fprintf(r.Error.Error())
			continue
		}
		if !r.Success {
			w.fprintf(w.au.Bold(w.au.Red(w.template.failures(r))))
			w.fprintf(r.Diff)
		}
	}
}

func (w *OutputWriter) fprintf(a ...interface{}) {
	if _, err := fmt.Fprintln(w.out, a...); err != nil {
		log.Fatal(err)
	}
}

// convert runtime.TestResult to output.TestResult
func convertTestResult(tr runtime.TestResult) TestResult {
	testResult := TestResult{
		FileName:       tr.TestCase.FileName,
		Title:          tr.TestCase.Title,
		Node:           tr.Node,
		Tries:          tr.Tries,
		Success:        tr.ValidationResult.Success,
		FailedProperty: tr.FailedProperty,
		Diff:           tr.ValidationResult.Diff,
		Error:          tr.TestCase.Result.Error,
	}

	return testResult
}
