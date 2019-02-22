package runtime

import (
    "bytes"
    "fmt"
    "github.com/SimonBaeumer/commander/pkg"
    "log"
    "os/exec"
    "strings"
)

type Command struct {
    cmd      string
    args     string
    exitCode int
    stderr   string
    stdout   string

    //env
    //timeout
    //...
}

// Start starts the given test suite
func Start(suite commander.Suite) {
    c := make(chan commander.TestCase)

    for _, t := range suite.GetTestCases() {
        go runTest(t, c)
    }

    printResults(c, suite)
}

func printResults(c chan commander.TestCase, suite commander.Suite) {
    counter := 0
    for r := range c {
        // Validate result
        if !r.Result.Success {
            fmt.Println("✗ ", r.Title)
        } else {
            fmt.Println("✓ ", r.Title)
        }

        counter++
        if (counter >= len(suite.GetTestCases())) {
            close(c)
        }
    }
}

func runTest(test commander.TestCase, results chan<- commander.TestCase) {
    // Create command
    cmd := compile(test.Command)

    // Execute command
    if err := cmd.Execute(); err != nil {
        log.Fatal(err)
    }

    // Write test result
    test.Result = commander.TestResult{
        ExitCode: cmd.exitCode,
        Stdout:   cmd.stdout,
        Stderr:   cmd.stderr,
    }

    result := Validate(test)
    test.Result.Success = result.Success
    test.Result.FailureProperty = result.Property

    // Send to result channel
    results <- test
}

func compile(command string) *Command {
    parts := strings.Split(command, " ")
    executable := parts[0]

    splitArgs := append(parts[:0], parts[1:]...)
    args := strings.Join(splitArgs, " ")

    return &Command{
        cmd: executable,
        args: args,
    }
}

// Execute executes a command on the system
func (c *Command) Execute() error {
    cmd := exec.Command(c.cmd, c.args)

    var (
        outBuff bytes.Buffer
        errBuff bytes.Buffer
    )
    cmd.Stdout = &outBuff
    cmd.Stderr = &errBuff

    err := cmd.Run()
    if err != nil {
        return err
    } else {
        c.exitCode = 0
    }

    c.stderr = errBuff.String()
    c.stdout = outBuff.String()

    return nil
}
