package runtime

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var _ Executor = (*DockerExecutor)(nil)

// DockerExecutor executes the test inside a docker container
type DockerExecutor struct {
	Image        string // Image which is started to execute the test
	Privileged   bool   // Enable privileged mode for the container
	ExecUser     string // ExecUser defines which user executes the docker container
	RegistryUser string
	RegistryPass string
}

// Execute executes the script inside a docker container
func (e DockerExecutor) Execute(test TestCase) TestResult {
	log.Printf("DOCKER_HOST: %s \n", os.Getenv("DOCKER_HOST"))
	log.Printf("DOCKER_CERT_PATH: %s \n", os.Getenv("DOCKER_CERT_PATH"))
	log.Printf("DOCKER_API_VERSION: %s \n", os.Getenv("DOCKER_API_VERSION"))

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithAPIVersionNegotiation(), client.FromEnv)
	if err != nil {
		test.Result.Error = err
		return TestResult{
			TestCase: test,
		}
	}

	authConfig := types.AuthConfig{
		Username: e.RegistryUser,
		Password: e.RegistryPass,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		panic(err)
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	log.Printf("Pulling image %s\n", e.Image)
	reader, err := cli.ImagePull(ctx, e.Image, types.ImagePullOptions{
		RegistryAuth: authStr,
	})
	if err != nil {
		test.Result.Error = fmt.Errorf("could not pull image '%s' with error: '%s'", e.Image, err)
		return TestResult{
			TestCase: test,
		}
	}
	buf := bytes.Buffer{}
	buf.ReadFrom(reader)
	log.Printf("Pull log image'%s':\n %s\n", e.Image, buf.String())

	var env []string
	for k, v := range test.Command.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	log.Printf("Create container %s\n", e.Image)
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:      e.Image,
			WorkingDir: test.Command.Dir,
			Env:        env,
			User:       e.ExecUser,
			Cmd:        []string{"/bin/sh", "-c", test.Command.Cmd},
			Tty:        false,
		}, nil, nil, nil, "")
	if err != nil {
		test.Result.Error = fmt.Errorf("could not pull image '%s' with error: '%s'", e.Image, err)
		return TestResult{
			TestCase: test,
		}
	}

	log.Printf("Started container %s %s\n", e.Image, resp.ID)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		test.Result.Error = fmt.Errorf("could not pull image '%s' with error: '%s'", e.Image, err)
		return TestResult{
			TestCase: test,
		}
	}

	duration := time.Duration(1 * time.Second)
	defer cli.ContainerStop(ctx, resp.ID, &duration)

	status := container.ContainerWaitOKBody{}
	statusCh, errC := cli.ContainerWait(ctx, resp.ID, "")
	select {
	case err := <-errC:
		if err != nil {
			panic(err)
		}
	case s := <-statusCh:
		status = s
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	_, err = stdcopy.StdCopy(stdout, stderr, out)
	if err != nil {
		panic(err)
	}

	log.Println("title: '"+test.Title+"'", " Command: ", test.Command.Cmd)
	log.Println("title: '"+test.Title+"'", " Directory: ", test.Command.Dir)
	log.Println("title: '"+test.Title+"'", " Env: ", test.Command.Env)

	// status := <-waitBody
	// Write test result
	test.Result = CommandResult{
		ExitCode: int(status.StatusCode),
		Stdout:   strings.TrimSpace(strings.Replace(stdout.String(), "\r\n", "\n", -1)),
		Stderr:   strings.TrimSpace(strings.Replace(stderr.String(), "\r\n", "\n", -1)),
	}

	log.Println("title: '"+test.Title+"'", " ExitCode: ", test.Result.ExitCode)
	log.Println("title: '"+test.Title+"'", " Stdout: ", test.Result.Stdout)
	log.Println("title: '"+test.Title+"'", " Stderr: ", test.Result.Stderr)

	return Validate(test)
}
