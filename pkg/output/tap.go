package output

import (
	"fmt"
	"github.com/commander-cli/commander/pkg/runtime"
	"io"
	"log"
	"os"
)

var _ Output = (*CLIOutputWriter)(nil)

// TAPOutputWriter writes TAP results
type TAPOutputWriter struct {
	out io.Writer
}

// NewTAPOutputWriter represents the output, defaults to stdout
func NewTAPOutputWriter() Output {
	return TAPOutputWriter{
		out: os.Stdout,
	}
}

func (w TAPOutputWriter) GetEventHandler() *runtime.EventHandler {
	return runtime.NewEmptyEventHandler()
}

func (w TAPOutputWriter) PrintSummary(result runtime.Result) {
	counter := 0
	for _, r := range result.TestResults {
		if r.Skipped {
			// skipped tests are not specified in the TAP specification
			continue
		}

		counter++
		if r.FailedProperty != "" {
			w.fprintf("%d ok - %s", counter+1, r.TestCase.Title)
		} else {
			w.fprintf("%d not ok - %s", counter+1, r.TestCase.Title)
		}
	}

	w.fprintf("1..%d", counter)
}

func (w TAPOutputWriter) fprintf(msg string , a ...interface{}) {
	if _, err := fmt.Fprintln(w.out, fmt.Sprintf(msg, a...)); err != nil {
		log.Fatal(err)
	}
}
