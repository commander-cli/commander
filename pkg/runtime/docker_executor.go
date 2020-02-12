package runtime

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"log"
	"strings"
	"time"
)

type DockerExecutor struct {
	Image string
	Name  string
}

func (e DockerExecutor) Execute(test TestCase) TestResult {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		test.Result.Error = err
		return TestResult{
			TestCase: test,
		}
	}

	//	reader, err := cli.ImagePull(ctx, e.Image, types.ImagePullOptions{})/
	//	if err != nil {
	//		panic(err)
	//	}
	//	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: e.Image,
		Cmd:   []string{"/bin/sh", "-c", test.Command.Cmd},
		Tty:   false,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	duration := time.Duration(1 * time.Second)
	defer cli.ContainerStop(ctx, resp.ID, &duration)
	fmt.Printf("Started container %s\n", resp.ID)

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
