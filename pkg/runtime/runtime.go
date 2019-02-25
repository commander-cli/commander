package runtime

import (
    "bytes"
    "fmt"
    "github.com/SimonBaeumer/commander/pkg"
    "github.com/SimonBaeumer/commander/pkg/output"
    "log"
    "os"
    "os/exec"
    "strings"
    "syscall"
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
    o := &output.HumanOutput{}
    counter := 0
    for r := range c {
        s := o.BuildTestResult(output.TestCase(r))
        fmt.Println(s)

        counter++
        if counter >= len(suite.GetTestCases()) {
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
    test.Result.FailureProperties = result.Properties

    // Send to result channel
    results <- test
}

func compile(command string) *Command {
    return &Command{
        cmd: command,
    }
}

// Execute executes a command on the system
func (c *Command) Execute() error {
    cmd := exec.Command("sh", "-c", c.cmd)
    env := os.Environ()
    cmd.Env = env

    var (
        outBuff bytes.Buffer
        errBuff bytes.Buffer
    )
    cmd.Stdout = &outBuff
    cmd.Stderr = &errBuff

    err := cmd.Start()
    log.Println("Started command " + c.cmd)
    if err != nil {
        log.Println("Started command " + c.cmd + " err: " + err.Error())
    }

    if err := cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
                c.exitCode = status.ExitStatus()
                //log.Printf("Exit Status: %d", status.ExitStatus())
            }
        }
    } else {
        c.exitCode = 0
    }
    c.stderr = strings.Trim(errBuff.String(), "\n")
    c.stdout = strings.Trim(outBuff.String(), "\n")

    return nil
}
