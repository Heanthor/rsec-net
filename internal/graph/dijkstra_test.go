package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDijkstraSearcher_ShortestPath_SameNode(t *testing.T) {
	var err error
	g := NewDirectedGraph()
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}

	err = g.AddNode(n1)
	assert.NoError(t, err)

	err = g.AddNode(n2)
	assert.NoError(t, err)

	err = g.AddEdge("n1", "n2", 5)
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n1")
	assert.Equal(t, []*Node{}, result)
}

func TestDijkstraSearcher_ShortestPath_SimpleGraph1(t *testing.T) {
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}
	g, err := NewDirectedGraphChain().
		AddNode(n1).
		AddNode(n2).
		AddEdge("n1", "n2", 5).
		DirectedGraph()
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n2")
	assert.Equal(t, []*Node{n1, n2}, result)
}

func TestDijkstraSearcher_ShortestPath_SimpleGraph2(t *testing.T) {
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}
	n3 := &Node{"n3", nil}

	g, err := NewDirectedGraphChain().
		AddNode(n1).
		AddNode(n2).
		AddNode(n3).
		AddEdge("n1", "n2", 5).
		AddEdge("n2", "n3", 7).
		DirectedGraph()
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n3")
	assert.Equal(t, []*Node{n1, n2, n3}, result)
}

func TestDijkstraSearcher_ShortestPath_TwoOptions1(t *testing.T) {
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}
	n3 := &Node{"n3", nil}
	n4 := &Node{"n4", nil}

	g, err := NewDirectedGraphChain().
		AddNode(n1).
		AddNode(n2).
		AddNode(n3).
		AddNode(n4).
		AddEdge("n1", "n2", 5).
		AddEdge("n2", "n3", 7).
		AddEdge("n3", "n4", 12).
		AddEdge("n1", "n4", 5).
		DirectedGraph()
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n4")
	assert.Equal(t, []*Node{n1, n4}, result)
}

func TestDijkstraSearcher_ShortestPath_TwoOptions2(t *testing.T) {
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}
	n3 := &Node{"n3", nil}
	n4 := &Node{"n4", nil}

	g, err := NewDirectedGraphChain().
		AddNode(n1).
		AddNode(n2).
		AddNode(n3).
		AddNode(n4).
		AddEdge("n1", "n2", 1).
		AddEdge("n2", "n3", 2).
		AddEdge("n3", "n4", 3).
		AddEdge("n1", "n4", 7).
		DirectedGraph()
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n4")
	assert.Equal(t, []*Node{n1, n2, n3, n4}, result)
}

func TestDijkstraSearcher_ShortestPath_Challenge1(t *testing.T) {
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}
	n3 := &Node{"n3", nil}
	n4 := &Node{"n4", nil}
	n5 := &Node{"n5", nil}
	n6 := &Node{"n6", nil}
	n7 := &Node{"n7", nil}
	n8 := &Node{"n8", nil}
	n9 := &Node{"n9", nil}
	n10 := &Node{"n10", nil}
	n11 := &Node{"n11", nil}
	n12 := &Node{"n12", nil}

	g, err := NewDirectedGraphChain().
		AddNode(n1).
		AddNode(n2).
		AddNode(n3).
		AddNode(n4).
		AddNode(n5).
		AddNode(n6).
		AddNode(n7).
		AddNode(n8).
		AddNode(n9).
		AddNode(n10).
		AddNode(n11).
		AddNode(n12).
		AddEdge("n1", "n2", 2).
		AddEdge("n1", "n3", 5).
		AddEdge("n1", "n4", 3).
		AddEdge("n2", "n5", 2).
		AddEdge("n3", "n5", 1).
		AddEdge("n3", "n6", 6).
		AddEdge("n4", "n6", 10).
		AddEdge("n5", "n6", 4).
		AddEdge("n5", "n7", 2).
		AddEdge("n5", "n8", 10).
		AddEdge("n6", "n9", 1).
		AddEdge("n7", "n10", 50).
		AddEdge("n8", "n10", 7).
		AddEdge("n8", "n11", 4).
		AddEdge("n9", "n11", 1).
		AddEdge("n10", "n12", 1).
		AddEdge("n11", "n10", 3).
		DirectedGraph()
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n12")
	assert.Equal(t, []*Node{n1, n2, n5, n6, n9, n11, n10, n12}, result)
}

func TestDijkstraSearcher_ShortestPath_Loop1(t *testing.T) {
	n1 := &Node{"n1", nil}
	n2 := &Node{"n2", nil}
	n3 := &Node{"n3", nil}
	n4 := &Node{"n4", nil}
	n5 := &Node{"n5", nil}

	g, err := NewDirectedGraphChain().
		AddNode(n1).
		AddNode(n2).
		AddNode(n3).
		AddNode(n4).
		AddNode(n5).
		AddEdge("n1", "n2", 2).
		AddEdge("n2", "n3", 3).
		AddEdge("n3", "n4", 2).
		AddEdge("n4", "n5", 1).
		AddEdge("n4", "n5", 10).
		DirectedGraph()
	assert.NoError(t, err)

	searcher := DijkstraSearcher{}
	result := searcher.ShortestPath(g, "n1", "n5")
	assert.Equal(t, []*Node{n1, n2, n3, n4, n5}, result)
}
