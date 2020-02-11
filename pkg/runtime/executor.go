package runtime

// Executor interface which will be implemented by all available executors, like ssh or local
type Executor interface {
	Execute(test TestCase) TestResult
}
