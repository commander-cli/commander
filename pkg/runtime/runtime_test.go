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

func TestRuntime_WithEnvVariables(t *testing.T) {
    s := getExampleTestSuite()
    s[0].Command.Env = []string{"KEY=value"}
    s[0].Command.Cmd = "echo $KEY"
    s[0].Expected.Stdout.Exactly = "value"
    s[0].Title = "Output env variable"

    got := Start(s)

    assert.IsType(t, make(<-chan TestResult), got)

    count := 0
    for r := range got {
        assert.Equal(t, "value", r.TestCase.Result.Stdout)
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


