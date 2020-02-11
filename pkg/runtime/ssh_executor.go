package runtime

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

// SSHExecutor
type SSHExecutor struct {
	Host         string
	User         string
	Password     string
	IdentityFile string
}

// Execute executes a command on a remote host viá SSH
func (e SSHExecutor) Execute(test TestCase) TestResult {
	if test.Command.InheritEnv {
		panic("Inherit env is not supported viá SSH")
	}

	// initialize auth methods with pass auth method as the default
	authMethods := []ssh.AuthMethod{
		ssh.Password(e.Password),
	}

	// add public key auth if identity file is given
	if e.IdentityFile != "" {
		signer := e.createSigner()
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	// create ssh config
	sshConf := &ssh.ClientConfig{
		User: e.User,
		Auth: authMethods,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// create ssh connection
	conn, err := ssh.Dial("tcp", e.Host, sshConf)
	if err != nil {
		log.Fatal(err)
	}

	// start session
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
			test.Result = CommandResult{
				Error: fmt.Errorf("Failed setting env variables, maybe ssh server is configured to only accept LC_ prefixed env variables. Error: %s", err),
			}
			return TestResult{
				TestCase: test,
			}
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

func (e SSHExecutor) createSigner() ssh.Signer {
	buffer, err := ioutil.ReadFile(e.IdentityFile)
	if err != nil {
		log.Fatal(err)
	}
	signer, err := ssh.ParsePrivateKey(buffer)
	return signer
}
