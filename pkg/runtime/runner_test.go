package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getExecutor(t *testing.T) {
	r := Runner{
		Nodes: getExampleNodes(),
	}

	exec := r.getExecutor("ssh")
	assert.IsType(t, SSHExecutor{}, exec)
	exec = r.getExecutor("local")
	assert.IsType(t, LocalExecutor{}, exec)
	exec = r.getExecutor("docker")
	assert.IsType(t, DockerExecutor{}, exec)
}

func getExampleNodes() []Node {
	n1 := Node{
		Name: "ssh",
		Type: "ssh",
	}
	n2 := Node{
		Name: "local",
		Type: "local",
	}
	n3 := Node{
		Name: "docker",
		Type: "docker",
	}

	nodes := []Node{
		n1, n2, n3,
	}
	return nodes
}

func getExampleTestCases() []TestCase {
	tests := []TestCase{
		{
			Command: CommandUnderTest{
				Cmd:     "echo hello",
				Timeout: "5s",
			},
			Expected: Expected{
				Stdout: ExpectedOut{
					Exactly: "hello",
				},
				ExitCode: 0,
			},
			Title: "Output hello",
		},
	}
	return tests
}
