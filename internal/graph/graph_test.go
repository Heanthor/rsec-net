package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectedGraph_AddNodeEmptyGraph(t *testing.T) {
	singleNode := Node{"hello", nil}

	graph := &DirectedGraph{
		adjList: make(map[string]value),
	}

	err := graph.AddNode(&singleNode)
	assert.NoError(t, err)

	expected := &DirectedGraph{
		adjList: map[string]value{
			singleNode.Key: {&singleNode, []edge{}},
		},
	}
	assert.Equal(t, expected, graph)
}

func TestDirectedGraph_AddNodeExistingNode(t *testing.T) {
	singleNode := Node{"hello", nil}

	graph := &DirectedGraph{
		adjList: map[string]value{
			singleNode.Key: {&singleNode, []edge{}},
		},
	}

	err := graph.AddNode(&singleNode)
	assert.Error(t, err)
}

func TestDirectedGraph_AddEdge(t *testing.T) {
	startNode := Node{"start", nil}
	endNode := Node{"end", nil}

	graph := &DirectedGraph{
		adjList: map[string]value{
			startNode.Key: {&startNode, []edge{}},
			endNode.Key:   {&endNode, []edge{}},
		},
	}

	err := graph.AddEdge(&startNode, &endNode, 5)
	assert.NoError(t, err)

	expected := &DirectedGraph{
		adjList: map[string]value{
			startNode.Key: {&startNode, []edge{
				{&endNode, 5},
			}},
			endNode.Key: {&endNode, []edge{}},
		},
	}

	assert.Equal(t, expected, graph)
}

func TestDirectedGraph_RemoveNode(t *testing.T) {
	n1 := Node{"n1", nil}
	n2 := Node{"n2", nil}
	n3 := Node{"n3", nil}

	graph := &DirectedGraph{
		adjList: map[string]value{
			n1.Key: {&n1, []edge{{&n3, 3}, {&n2, 5}}},
			n2.Key: {&n2, []edge{{&n1, 10}}},
			n3.Key: {&n3, []edge{{&n2, 20}}},
		},
	}

	err := graph.RemoveNode(n2.Key)
	assert.NoError(t, err)

	expected := &DirectedGraph{
		adjList: map[string]value{
			n1.Key: {&n1, []edge{{&n3, 3}}},
			n3.Key: {&n3, []edge{}},
		},
	}
	assert.Equal(t, expected, graph)
}

func TestDirectedGraph_RemoveEdge(t *testing.T) {
	n1 := Node{"n1", nil}
	n2 := Node{"n2", nil}
	n3 := Node{"n3", nil}

	graph := &DirectedGraph{
		adjList: map[string]value{
			n1.Key: {&n1, []edge{{&n3, 3}, {&n2, 5}}},
			n2.Key: {&n2, []edge{{&n1, 10}}},
			n3.Key: {&n3, []edge{{&n2, 20}}},
		},
	}

	err := graph.RemoveEdge(n1.Key, n3.Key)
	assert.NoError(t, err)

	expected := &DirectedGraph{
		adjList: map[string]value{
			n1.Key: {&n1, []edge{{&n2, 5}}},
			n2.Key: {&n2, []edge{{&n1, 10}}},
			n3.Key: {&n3, []edge{{&n2, 20}}},
		},
	}
	assert.Equal(t, expected, graph)
}

func TestDirectedGraph_GetEdgeCost(t *testing.T) {
	n1 := Node{"n1", nil}
	n2 := Node{"n2", nil}
	n3 := Node{"n3", nil}

	graph := &DirectedGraph{
		adjList: map[string]value{
			n1.Key: {&n1, []edge{{&n3, 3}, {&n2, 5}}},
			n2.Key: {&n2, []edge{{&n1, 10}}},
			n3.Key: {&n3, []edge{{&n2, 20}}},
		},
	}

	cost, err := graph.GetEdgeCost(n1.Key, n2.Key)
	assert.NoError(t, err)
	assert.Equal(t, 5, cost)
}
