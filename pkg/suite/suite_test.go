package suite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetNodes(t *testing.T) {
	s := Suite{
		Nodes: []Node{{}, {}},
	}

	assert.Len(t, s.GetNodes(), 2)
}

func Test_GetNodesByName(t *testing.T) {
	s := Suite{
		Nodes: []Node{{}, {Name: "node1"}},
	}

	node, e := s.GetNodeByName("node1")
	assert.Equal(t, node.Name, "node1")
	assert.Nil(t, e)

	node, e = s.GetNodeByName("doesnt-exist")
	assert.EqualError(t, e, "could not find node with name doesnt-exist")
}

func Test_AddTest(t *testing.T) {
	s := Suite{TestCases: []TestCase{{Title: "exists"}}}
	s.AddTest(TestCase{Title: "test"})

	assert.Len(t, s.GetTests(), 1)
}

func Test_GetTestByTitle(t *testing.T) {
	s := Suite{TestCases: []TestCase{{Title: "exists"}}}
	test, err := s.GetTestByTitle("exists")

	assert.Nil(t, err)
	assert.Equal(t, "exists", test.Title)
}

func Test_GetGlobalConfig(t *testing.T) {
	s := Suite{Config: GlobalTestConfig{Dir: "/tmp"}}
	assert.Equal(t, "/tmp", s.GetGlobalConfig().Dir)
}

func Test_FindTests(t *testing.T) {
	s := Suite{TestCases: []TestCase{
		{Title: "exists"},
		{Title: "another"},
		{Title: "another one"},
	}}

	test, _ := s.FindTests("exists")
	assert.Len(t, test, 1)
}

func Test_FindMultipleTests(t *testing.T) {
	s := Suite{TestCases: []TestCase{
		{Title: "exists"},
		{Title: "another"},
		{Title: "another one"},
	}}

	test, _ := s.FindTests("another")
	assert.Len(t, test, 2)

	test, _ = s.FindTests("another$")
	assert.Len(t, test, 1)
}
