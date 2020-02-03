package runtime

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
)

// SSHExecutor
type SSHExecutor struct {
	Host     string
	User     string
	Password string
}

// Execute executes a command on a remote host viá SSH
func (e SSHExecutor) Execute(test TestCase) TestResult {
	if test.Command.InheritEnv {
		log.Fatal("Inhereit env is not supported viá SSH")
	}

	sshConf := &ssh.ClientConfig{
		User: e.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(e.Password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	conn, err := ssh.Dial("tcp", e.Host, sshConf)
	if err != nil {
		log.Fatal(err)
	}

	session, err := conn.NewSession()
	defer session.Close()
	if err != nil {
		log.Fatal(err)
	}

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	session.Stdout = &stdoutBuffer
	session.Stderr = &stderrBuffer

	for k, v := range test.Command.Env {
		err := session.Setenv(k, v)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed, maybe ssh server is configured to only accept LC_ prefixed env variables. Error: %s", err))
		}
	}

	dirCmd := ""
	if test.Command.Dir != "" {
		dirCmd = fmt.Sprintf("cd %s; ", test.Command.Dir)
	}

	exitCode := 0
	err = session.Run(fmt.Sprintf("%s %s", dirCmd, test.Command.Cmd))
	switch err.(type) {
	case *ssh.ExitError:
		ee, _ := err.(*ssh.ExitError)
		exitCode = ee.Waitmsg.ExitStatus()
	case nil:
		break
	default:
		log.Println(test.Title, " failed ", err.Error())
		test.Result = CommandResult{
			Error: err,
		}

		return TestResult{
			TestCase: test,
		}
	}

	test.Result = CommandResult{
		ExitCode: exitCode,
		Stdout:   strings.TrimSpace(strings.Replace(stdoutBuffer.String(), "\r\n", "\n", -1)),
		Stderr:   strings.TrimSpace(strings.Replace(stderrBuffer.String(), "\r\n", "\n", -1)),
	}

	log.Println("title: '"+test.Title+"'", " ExitCode: ", test.Result.ExitCode)
	log.Println("title: '"+test.Title+"'", " Stdout: ", test.Result.Stdout)
	log.Println("title: '"+test.Title+"'", " Stderr: ", test.Result.Stderr)

	return Validate(test)
}
