package runtime

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/commander-cli/cmd"
)

var _ Executor = (*LocalExecutor)(nil)

// LocalExecutor will be used to execute tests on the local host
type LocalExecutor struct{}

// NewLocalExecutor creates a new local executor
func NewLocalExecutor() Executor {
	return LocalExecutor{}
}

// Execute will execute the given test on the current node
func (e LocalExecutor) Execute(test TestCase) TestResult {
	timeoutOpt, err := createTimeoutOption(test.Command.Timeout)
	if err != nil {
		test.Result = CommandResult{Error: err}
		return TestResult{
			TestCase: test,
		}
	}

	envOpt := createEnvVarsOption(test)

	// cut = command under test
	cut := cmd.NewCommand(
		test.Command.Cmd,
		cmd.WithWorkingDir(test.Command.Dir),
		timeoutOpt,
		envOpt)

	if err := cut.Execute(); err != nil {
		log.Println(test.Title, " failed ", err.Error())
		test.Result = CommandResult{
			Error: err,
		}

		return TestResult{
			TestCase: test,
		}
	}

	log.Println("title: '"+test.Title+"'", " Command: ", test.Command.Cmd)
	log.Println("title: '"+test.Title+"'", " Directory: ", cut.Dir)
	log.Println("title: '"+test.Title+"'", " Env: ", cut.Env)

	// Write test result
	test.Result = CommandResult{
		ExitCode: cut.ExitCode(),
		Stdout:   strings.TrimSpace(strings.ReplaceAll(cut.Stdout(), "\r\n", "\n")),
		Stderr:   strings.TrimSpace(strings.ReplaceAll(cut.Stderr(), "\r\n", "\n")),
	}

	log.Println("title: '"+test.Title+"'", " ExitCode: ", test.Result.ExitCode)
	log.Println("title: '"+test.Title+"'", " Stdout: ", test.Result.Stdout)
	log.Println("title: '"+test.Title+"'", " Stderr: ", test.Result.Stderr)

	return Validate(test)
}

func createEnvVarsOption(test TestCase) func(c *cmd.Command) {
	return func(c *cmd.Command) {
		// Add all env variables from parent process
		if test.Command.InheritEnv {
			for _, v := range os.Environ() {
				split := strings.Split(v, "=")
				c.AddEnv(split[0], split[1])
			}
		}

		// Add custom env variables
		for k, v := range test.Command.Env {
			c.AddEnv(k, v)
		}
	}
}

func createTimeoutOption(timeout string) (func(c *cmd.Command), error) {
	timeoutOpt := cmd.WithoutTimeout
	if timeout != "" {
		d, err := time.ParseDuration(timeout)
		if err != nil {
			return func(c *cmd.Command) {}, err
		}
		timeoutOpt = cmd.WithTimeout(d)
	}

	return timeoutOpt, nil
}
