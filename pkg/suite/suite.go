package suite

import "github.com/SimonBaeumer/commander/pkg/runtime"

type Suite interface {
	GetTests()					 []runtime.TestCase
	GetTestByTitle(title string) (runtime.TestCase, error)
}
