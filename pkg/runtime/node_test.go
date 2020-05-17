package runtime

import (
	"github.com/SimonBaeumer/cmd"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNodeExpandEnv(t *testing.T) {
	os.Setenv("NAME", "test")
	os.Setenv("TYPE", "docker")
	os.Setenv("USER", "user")
	os.Setenv("PASS", "pass")
	os.Setenv("ADDR", "addr")
	os.Setenv("IMAGE", "image")
	os.Setenv("IDENTITY_FILE", "identity-file")
	os.Setenv("DOCKER_EXEC_USER", "docker-exec-user")

	n := Node{
		Name:           "$NAME",
		Type:           "$TYPE",
		User:           "$USER",
		Pass:           "$PASS",
		Addr:           "$ADDR",
		Image:          "$IMAGE",
		IdentityFile:   "$IDENTITY_FILE",
		DockerExecUser: "$DOCKER_EXEC_USER",
	}

	n.ExpandEnv()

	assert.Equal(t, "test", n.Name)
	assert.Equal(t, "docker", n.Type)
	assert.Equal(t, "user", n.User)
	assert.Equal(t, "pass", n.Pass)
	assert.Equal(t, "addr", n.Addr)
	assert.Equal(t, "image", n.Image)
	assert.Equal(t, "identity-file", n.IdentityFile)
	assert.Equal(t, "docker-exec-user", n.DockerExecUser)
}

func TestExpandEnv_PrintWarning(t *testing.T) {
	n := Node{
		Name: "Node01",
		Pass: "clean password",
	}

	out, _ := cmd.CaptureStandardOutput(func() interface{} {
		n.ExpandEnv()
		return nil
	})

	assert.Equal(t, "WARNING: Consider using env variables with $VAR or ${VAR} in node Node01 instead of directly adding passwords to config files.\n", out)
}
