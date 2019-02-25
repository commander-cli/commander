package output

import (
	"github.com/SimonBaeumer/commander/pkg"
)

type TestCase commander.TestCase

type Output interface {
	BuildHeader()
	BuildTestResult(test TestCase)
	BuildSuiteResult()
	GetBuffer() []string
}
