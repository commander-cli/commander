package runtime

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestRuntime_Start(t *testing.T) {
    s := getExampleTestSuite()
    got := Start(s)

    assert.IsType(t, make(<-chan TestResult), got)

    count := 0
    for r := range got {
        assert.Equal(t, "Output hello", r.TestCase.Title)
        assert.True(t, r.ValidationResult.Success)
        count++
    }
    assert.Equal(t, 1, count)
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


