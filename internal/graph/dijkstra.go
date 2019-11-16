package graph

import "math"

/*
https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm#Pseudocode
simple

 1  function Dijkstra(Graph, source):
 2
 3      create vertex set Q
 4
 5      for each vertex v in Graph:
 6          dist[v] ← INFINITY
 7          prev[v] ← UNDEFINED
 8          add v to Q
10      dist[source] ← 0
11
12      while Q is not empty:
13          u ← vertex in Q with min dist[u]
14
15          remove u from Q
16
17          for each neighbor v of u:           // only v that are still in Q
18              alt ← dist[u] + length(u, v)
19              if alt < dist[v]:
20                  dist[v] ← alt
21                  prev[v] ← u
22
23      return dist[], prev[]

priority queue-based

1  function Dijkstra(Graph, source):
2      dist[source] ← 0                           // Initialization
3
4      create vertex priority queue Q
5
6      for each vertex v in Graph:
7          if v ≠ source
8              dist[v] ← INFINITY                 // Unknown distance from source to v
9          prev[v] ← UNDEFINED                    // Predecessor of v
10
11         Q.add_with_priority(v, dist[v])
12
13
14     while Q is not empty:                      // The main loop
15         u ← Q.extract_min()                    // Remove and return best vertex
16         for each neighbor v of u:              // only v that are still in Q
17             alt ← dist[u] + length(u, v)
18             if alt < dist[v]
19                 dist[v] ← alt
20                 prev[v] ← u
21                 Q.decrease_priority(v, alt)
22
23     return dist, prev
*/

type DijkstraSearcher struct {
}

func (d *DijkstraSearcher) ShortestPath(graph *DirectedGraph, startKey, targetKey string) []*Node {
	if startKey == targetKey {
		return []*Node{}
	}

	pendingNodes := make(map[string]struct{})
	distanceTo := make(map[string]int)
	prevHop := make(map[string]*Node)

	for vertexKey := range graph.adjList {
		distanceTo[vertexKey] = math.MaxInt32
		prevHop[vertexKey] = nil
		pendingNodes[vertexKey] = struct{}{}
	}

	distanceTo[startKey] = 0

	for len(pendingNodes) > 0 {
		// u ← vertex in Q with min dist[u]
		minDistance := math.MaxInt32
		var minDistanceKey string

		for vertexKey := range pendingNodes {
			if distance := distanceTo[vertexKey]; distance <= minDistance {
				minDistance = distance
				minDistanceKey = vertexKey
			}
		}

		// remove u from Q
		delete(pendingNodes, minDistanceKey)

		if minDistanceKey == targetKey {
			path := []*Node{}
			tmp, err := graph.GetNode(targetKey)
			if err != nil {
				panic(err)
			}

			if prevHop[targetKey] != nil || tmp.Key == startKey {
				for tmp != nil {
					path = append([]*Node{tmp}, path...)
					tmp = prevHop[tmp.Key]
				}
			}

			return path
		}

		edges := graph.adjList[minDistanceKey].edges
		for _, e := range edges {
			candidateDistance := minDistance + e.Cost

			k := e.Dest.Key
			if candidateDistance < distanceTo[k] {
				distanceTo[k] = candidateDistance
				node, err := graph.GetNode(minDistanceKey)
				if err != nil {
					// this will never happen..
					panic(err)
				}

				prevHop[k] = node
			}
		}
	}

	return nil
}
