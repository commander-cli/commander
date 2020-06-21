package output

import (
	"fmt"
	"io"
	"log"
	"os"
	run "runtime"
	"time"

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

// Result respresents the aggregation of all TestResults/summary of a runtime
type Result struct {
	TestResults []TestResult
	Duration    time.Duration
	Failed      int
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

// PrintSummary prints summary
func (w *OutputWriter) PrintSummary(result Result) bool {
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
func (w *OutputWriter) PrintResult(r TestResult) {
	if !r.Success {
		w.fprintf(w.au.Red(w.template.testResult(r)))
		return
	}
	w.fprintf(w.template.testResult(r))
}

func (w *OutputWriter) printFailures(results []TestResult) {
	w.fprintf("")
	w.fprintf(w.au.Bold("Results"))
	w.fprintf(w.au.Bold(""))

	for _, r := range results {
		if r.Error != nil {
			str := fmt.Sprintf("✗ [%s] '%s' could not be executed with error message:", r.Node, r.Title)
			str = w.addFile(str, r)
			w.fprintf(w.au.Bold(w.au.Red(str)))
			w.fprintf(r.Error.Error())
			continue
		}

		if !r.Success {
			str := fmt.Sprintf("✗ [%s] '%s', on property '%s'", r.Node, r.Title, r.FailedProperty)
			str = w.addFile(str, r)
			w.fprintf(w.au.Bold(w.au.Red(str)))
			w.fprintf(r.Diff)
		}
	}
}

func (w *OutputWriter) addFile(s string, r TestResult) string {
	if r.FileName == "" {
		return s
	}
	s = s[:3] + " [" + r.FileName + "]" + s[3:]
	return s
}

func (w *OutputWriter) addTries(s string, r TestResult) string {
	if r.Tries > 1 {
		s = fmt.Sprintf("%s, retries %d", s, r.Tries)
	}
	return s
}

func (w *OutputWriter) fprintf(a ...interface{}) {
	if _, err := fmt.Fprintln(w.out, a...); err != nil {
		log.Fatal(err)
	}
}
