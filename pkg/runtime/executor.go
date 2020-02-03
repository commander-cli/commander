package runtime

type Executor interface {
	Execute(test TestCase) TestResult
}
