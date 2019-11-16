package graph

type Searcher interface {
	ShortestPath(graph *DirectedGraph, startKey, targetKey string) []*Node
}
