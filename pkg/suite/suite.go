package suite

import "github.com/SimonBaeumer/commander/pkg/runtime"

// Suite represents the current tests and configs
type Suite interface {
	GetTests() []runtime.TestCase
	GetTestByTitle(title string) (runtime.TestCase, error)
	GetGlobalConfig() runtime.TestConfig
}
