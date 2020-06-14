package runtime

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Runner holds the config and executes the desired runtime env
type Runner struct {
	Nodes []Node
}

// Execute the runner
func (r *Runner) Execute(tests []TestCase) <-chan TestResult {
	in := make(chan TestCase)
	out := make(chan TestResult)

	go func(tests []TestCase) {
		defer close(in)
		for _, t := range tests {
			in <- t
		}
	}(tests)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(tests chan TestCase) {
		defer wg.Done()

		for t := range tests {
			// If no node was set use local mode as default
			if len(t.Nodes) == 0 {
				t.Nodes = []string{"local"}
			}

			for _, n := range t.Nodes {
				result := TestResult{}
				for i := 1; i <= t.Command.GetRetries(); i++ {

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
	if len(r.Nodes) == 0 {
		return NewLocalExecutor()
	}

	for _, n := range r.Nodes {
		if n.Name == node {
			switch n.Type {
			case "ssh":
				return NewSSHExecutor(n.Addr, n.User, WithIdentityFile(n.IdentityFile), WithPassword(n.Pass))
			case "local":
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
			case "":
				return NewLocalExecutor()
			default:
				log.Fatal(fmt.Sprintf("Node type %s not found for node %s", n.Type, n.Name))
			}
		}
	}

	log.Fatal(fmt.Sprintf("Node %s not found", node))
	return NewLocalExecutor()
}

func executeRetryInterval(t TestCase) {
	if t.Command.GetRetries() > 1 && t.Command.Interval != "" {
		interval, err := time.ParseDuration(t.Command.Interval)
		if err != nil {
			panic(fmt.Sprintf("'%s' interval error: %s", t.Command.Cmd, err))
		}
		time.Sleep(interval)
	}
}
