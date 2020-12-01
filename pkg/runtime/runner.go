package runtime

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/commander-cli/commander/pkg/suite"
)

// Runner holds the config and executes the desired runtime env
type Runner struct {
	Nodes []suite.Node
}

// Run the runner
func (r *Runner) Run(tests []suite.TestCase) <-chan TestResult {
	in := make(chan suite.TestCase)
	out := make(chan TestResult)

	go func(tests []suite.TestCase) {
		defer close(in)
		for _, t := range tests {
			in <- t
		}
	}(tests)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(tests chan suite.TestCase) {
		defer wg.Done()

		for t := range tests {
			// If no node was set use local mode as default
			if len(t.Nodes) == 0 {
				t.Nodes = []string{"local"}
			}

			for _, n := range t.Nodes {
				result := TestResult{}
				for i := 1; i <= t.Command.GetRetries(); i++ {

					if t.Skip {
						result = TestResult{TestCase: t, Skipped: true, Node: n}
						break
					}

					e := r.getExecutor(n)
					result = e.Execute(t)
					result.Node = n
					result.Tries = i

					if result.ValidationResult.Success {
						break
					}

					executeRetryInterval(t)
				}
				out <- result
			}
		}
	}(in)

	go func(results chan TestResult) {
		wg.Wait()
		close(results)
	}(out)

	return out
}

// getExecutor gets the node by the name it matches within the runner config
func (r *Runner) getExecutor(node string) Executor {
	for _, n := range r.Nodes {
		if n.Name == node {
			switch n.Type {
			case "ssh":
				return NewSSHExecutor(n.Addr, n.User, WithIdentityFile(n.IdentityFile), WithPassword(n.Pass))
			case "local", "":
				return NewLocalExecutor()
			case "docker":
				log.Println("Use docker executor")
				return DockerExecutor{
					Image:        n.Image,
					Privileged:   n.Privileged,
					ExecUser:     n.DockerExecUser,
					RegistryPass: n.Pass,
					RegistryUser: n.User,
				}
			default:
				log.Fatal(fmt.Sprintf("Node type %s not found for node %s", n.Type, n.Name))
			}
		}
	}

	log.Fatal(fmt.Sprintf("Node %s not found", node))
	return NewLocalExecutor()
}

func executeRetryInterval(t suite.TestCase) {
	if t.Command.GetRetries() > 1 && t.Command.Interval != "" {
		interval, err := time.ParseDuration(t.Command.Interval)
		if err != nil {
			panic(fmt.Sprintf("'%s' interval error: %s", t.Command.Cmd, err))
		}
		time.Sleep(interval)
	}
}
