package runtime

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// DockerExecutor executes the test inside a docker container
type DockerExecutor struct {
	Image string
}

// Execute executes the script inside a docker container
func (e DockerExecutor) Execute(test TestCase) TestResult {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		test.Result.Error = err
		return TestResult{
			TestCase: test,
		}
	}

	log.Printf("Pulling image %s\n", e.Image)
	reader, err := cli.ImagePull(ctx, e.Image, types.ImagePullOptions{})
	if err != nil {
		test.Result.Error = fmt.Errorf("could not pull image '%s' with error: '%s'", e.Image, err)
		return TestResult{
			TestCase: test,
		}
	}
	io.Copy(os.Stdout, reader)

	var env []string
	for k, v := range test.Command.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	log.Printf("Started container")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      e.Image,
		WorkingDir: test.Command.Dir,
		Env:        env,
		Cmd:        []string{"/bin/sh", "-c", test.Command.Cmd},
		Tty:        false,
	}, nil, nil, "")
	if err != nil {
		test.Result.Error = fmt.Errorf("could not pull image '%s' with error: '%s'", e.Image, err)
		return TestResult{
			TestCase: test,
		}
	}

	log.Printf("Started container %s\n", resp.ID)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		test.Result.Error = fmt.Errorf("could not pull image '%s' with error: '%s'", e.Image, err)
		return TestResult{
			TestCase: test,
		}
	}
	duration := time.Duration(1 * time.Second)
	defer cli.ContainerStop(ctx, resp.ID, &duration)

	status, err := cli.ContainerWait(ctx, resp.ID)
	fmt.Printf("status %d \n", status)
	if err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	written, err := stdcopy.StdCopy(stdout, stderr, out)
	if err != nil {
		fmt.Printf("Written %d\n", written)
		panic(err)
	}

	log.Println("title: '"+test.Title+"'", " Command: ", test.Command.Cmd)
	//log.Println("title: '"+test.Title+"'", " Directory: ", cut.Dir)
	//log.Println("title: '"+test.Title+"'", " Env: ", cut.Env)

	// Write test result
	test.Result = CommandResult{
		ExitCode: int(status),
		Stdout:   strings.TrimSpace(strings.Replace(stdout.String(), "\r\n", "\n", -1)),
		Stderr:   strings.TrimSpace(strings.Replace(stderr.String(), "\r\n", "\n", -1)),
	}

	log.Println("title: '"+test.Title+"'", " ExitCode: ", test.Result.ExitCode)
	log.Println("title: '"+test.Title+"'", " Stdout: ", test.Result.Stdout)
	log.Println("title: '"+test.Title+"'", " Stderr: ", test.Result.Stderr)

	return Validate(test)
}
