package suite

import (
	"github.com/commander-cli/commander/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetNodes(t *testing.T) {
	s := Suite{
		Nodes: []runtime.Node{runtime.Node{}, runtime.Node{}},
	}

	assert.Len(t, s.GetNodes(), 2)
}

func Test_GetNodesByName(t *testing.T) {
	s := Suite{
		Nodes: []runtime.Node{runtime.Node{}, runtime.Node{Name: "node1"}},
	}

	node, e := s.GetNodeByName("node1")
	assert.Equal(t, node.Name, "node1")
	assert.Nil(t, e)

	node, e = s.GetNodeByName("doesnt-exist")
	assert.EqualError(t, e, "could not find node with name doesnt-exist")
}

func Test_AddTest(t *testing.T) {
	s := Suite{TestCases: []runtime.TestCase{{Title: "exists"}}}
	s.AddTest(runtime.TestCase{Title: "test"})

	assert.Len(t, s.GetTests(), 1)
}

func Test_GetTestByTitle(t *testing.T) {
	s := Suite{TestCases: []runtime.TestCase{{Title: "exists"}}}
	test, err := s.GetTestByTitle("exists")

	assert.Nil(t, err)
	assert.Equal(t, "exists", test.Title)
}

func Test_GetGlobalConfig(t *testing.T) {
	s := Suite{Config: runtime.GlobalTestConfig{Dir: "/tmp"}}
	assert.Equal(t, "/tmp", s.GetGlobalConfig().Dir)
}
