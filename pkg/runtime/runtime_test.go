package runtime

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestRuntime_Start(t *testing.T) {
    s := getExampleTestSuite()
    got := Start(s)

    assert.IsType(t, make(<-chan TestResult), got)
}

func getExampleTestSuite() []TestCase {
    tests := []TestCase{
        {
            Command: CommandUnderTest{
                Cmd: "echo hello",
            },
            Expected: Expected{
                Stdout: ExpectedOut{
                    Exactly: "hello",
                },
                ExitCode: 0,
            },
            Title: "Output hello",
        },
    }
    return tests
}

