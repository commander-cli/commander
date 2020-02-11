package runtime

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

type SSHExecutorTestEnv struct {
	User         string
	Pass         string
	Host         string
	IdentityFile string
}

var SSHTestEnv SSHExecutorTestEnv

func TestMain(m *testing.M) {
	v := os.Getenv("COMMANDER_SSH_TEST")
	if v != "1" {
		log.Println("Skip ssh_executor_test, set env COMMANDER_SSH_TEST to 1")
		return
	}

	SSHTestEnv = SSHExecutorTestEnv{
		Host:         os.Getenv("COMMANDER_TEST_SSH_HOST"),
		Pass:         os.Getenv("COMMANDER_TEST_SSH_PASS"),
		User:         os.Getenv("COMMANDER_TEST_SSH_USER"),
		IdentityFile: os.Getenv("COMMANDER_TEST_SSH_IDENTITY_FILE"),
	}

	m.Run()
}

func createExecutor() SSHExecutor {
	s := SSHExecutor{
		Host:         SSHTestEnv.Host,
		User:         SSHTestEnv.User,
		Password:     SSHTestEnv.Pass,
		IdentityFile: SSHTestEnv.IdentityFile,
	}
	return s
}

func Test_SSHExecutor(t *testing.T) {
	s := createExecutor()

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: "echo test",
		},
		Expected: Expected{
			ExitCode: 0,
			Stdout:   ExpectedOut{Exactly: "test"},
		},
	}
	got := s.Execute(test)

	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "test", got.TestCase.Result.Stdout)
}

func Test_SSHExecutor_WithDir(t *testing.T) {
	s := createExecutor()

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: "echo $LC_TEST_KEY1; echo $LC_TEST_KEY2",
			Env: map[string]string{
				"LC_TEST_KEY1": "ENV_VALUE1",
				"LC_TEST_KEY2": "ENV_VALUE2",
			},
		},
	}
	got := s.Execute(test)

	assert.True(t, got.ValidationResult.Success)
	assert.Equal(t, "ENV_VALUE1\nENV_VALUE2", got.TestCase.Result.Stdout)
	assert.Equal(t, 0, got.TestCase.Result.ExitCode)
}

func Test_SSHExecutor_ExitCode(t *testing.T) {
	s := createExecutor()

	test := TestCase{
		Command: CommandUnderTest{
			Cmd: "exit 2;",
		},
	}

	got := s.Execute(test)

	assert.False(t, got.ValidationResult.Success)
	assert.Equal(t, 2, got.TestCase.Result.ExitCode)
}
