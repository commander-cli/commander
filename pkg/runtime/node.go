package runtime

import (
	"fmt"
	"os"
	"strings"
)

// Node represents a configured node with everything needed  to connect to the given host
// which is defined in the type property
// If the type is not available the test will fail and stop its execution
type Node struct {
	Name           string
	Type           string
	User           string
	Pass           string
	Addr           string
	Image          string
	IdentityFile   string
	Privileged     bool
	DockerExecUser string
}

func (n *Node) ExpandEnv() {
	n.Name = os.ExpandEnv(n.Name)
	n.Type = os.ExpandEnv(n.Type)
	n.User = os.ExpandEnv(n.User)
	n.Addr = os.ExpandEnv(n.Addr)
	n.Image = os.ExpandEnv(n.Image)
	n.IdentityFile = os.ExpandEnv(n.IdentityFile)
	n.DockerExecUser = os.ExpandEnv(n.DockerExecUser)

	if n.Pass != "" && !strings.Contains(n.Pass, "$") {
		fmt.Printf("WARNING: Consider using env variables with $VAR or ${VAR} in node %s instead of directly adding passwords to config files.\n", n.Name)
	}
	n.Pass = os.ExpandEnv(n.Pass)
}
