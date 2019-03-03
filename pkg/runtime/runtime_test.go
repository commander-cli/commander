package runtime

import (
    "github.com/SimonBaeumer/commander/pkg/suite"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestRuntime_Start(t *testing.T) {
    s := getExampleTestSuite()
    got := Start(s)

    assert.True(t, got.Success)
    assert.Len(t, got.TestResults, 1)
}

func TestRuntime_StartFail(t *testing.T) {
    s := getExampleTestSuite()
    s.Tests[0].Command.Cmd = "echo not expected"
    got := Start(s)

    assert.False(t, got.Success)
    assert.Len(t, got.TestResults, 1)
}

func getExampleTestSuite() suite.Suite {
    tests := []suite.TestCase{
        {
            Command: suite.CommandUnderTest{
                Cmd: "echo hello",
            },
            Expected: suite.Expected{
                Stdout: suite.ExpectedOut{
                    Exactly: "hello",
                },
                ExitCode: 0,
            },
            Title: "Output hello",
        },
    }
    return *suite.NewSuite(tests)
}

