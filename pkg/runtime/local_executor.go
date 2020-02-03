package runtime

import (
	"github.com/SimonBaeumer/cmd"
	"log"
	"strings"
)

type LocalExecutor struct {
}

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
		Stdout:   strings.TrimSpace(strings.Replace(cut.Stdout(), "\r\n", "\n", -1)),
		Stderr:   strings.TrimSpace(strings.Replace(cut.Stderr(), "\r\n", "\n", -1)),
	}

	log.Println("title: '"+test.Title+"'", " ExitCode: ", test.Result.ExitCode)
	log.Println("title: '"+test.Title+"'", " Stdout: ", test.Result.Stdout)
	log.Println("title: '"+test.Title+"'", " Stderr: ", test.Result.Stderr)

	return Validate(test)
}
