package runtime

type DockerExecutor struct {
	Image string
}

func (e DockerExecutor) Execute(test TestCase) TestResult {
	return TestResult{}
}
