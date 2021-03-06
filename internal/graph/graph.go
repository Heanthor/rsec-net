package graph

import "errors"

type Node struct {
	Key  string
	Data interface{}
}

type edge struct {
	Dest *Node
	Cost int
}

type value struct {
	start *Node
	edges []edge
}

// DirectedGraph represents a directed graph as an adjacency list.
type DirectedGraph struct {
	adjList map[string]value
}

// NewDirectedGraph returns an empty directed graph.
func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{make(map[string]value)}
}

// AddNode adds a new node to the graph.
func (d *DirectedGraph) AddNode(n *Node) error {
	if _, ok := d.adjList[n.Key]; ok {
		return errors.New("node with key already in graph")
	}

	d.adjList[n.Key] = value{n, []edge{}}

	return nil
}

// RemoveNode deletes a node (by key) and any associated edges from the graph.
func (d *DirectedGraph) RemoveNode(key string) error {
	if _, ok := d.adjList[key]; !ok {
		return errors.New("node with key not in graph")
	}

	// delete the node
	delete(d.adjList, key)

	// delete any edges that contain the node as an endpoint
	for k, value := range d.adjList {
		for i, e := range value.edges {
			if e.Dest.Key == key {
				// cut element from list
				value.edges = append(value.edges[:i], value.edges[i+1:]...)
				d.adjList[k] = value
			}
		}
	}

	return nil
}

// AddNode adds a new node to the graph.
func (d *DirectedGraph) GetNode(key string) (*Node, error) {
	if n, ok := d.adjList[key]; !ok {
		return nil, errors.New("node with key not in graph")
	} else {
		return n.start, nil
	}
}

// AddEdge adds a new edge to the graph with a cost.
func (d *DirectedGraph) AddEdge(start, end string, cost int) error {
	if _, ok := d.adjList[start]; !ok {
		return errors.New("start node with key not in graph")
	}

	if _, ok := d.adjList[end]; !ok {
		return errors.New("end node with key not in graph")
	}

	value := d.adjList[start]
	endNode, err := d.GetNode(end)
	if err != nil {
		return err
	}

	value.edges = append(value.edges, edge{endNode, cost})

	d.adjList[start] = value

	return nil
}

// RemoveEdge deletes a an edge from the graph.
func (d *DirectedGraph) RemoveEdge(start, end string) error {
	if _, ok := d.adjList[start]; !ok {
		return errors.New("start node with key not in graph")
	}

	if _, ok := d.adjList[end]; !ok {
		return errors.New("end node with key not in graph")
	}

	// iterate through start node's edges, until we find end
	value := d.adjList[start]
	for i, e := range value.edges {
		if e.Dest.Key == end {
			value.edges = append(value.edges[:i], value.edges[i+1:]...)
			d.adjList[start] = value
		}
	}

	return nil
}

// GetEdgeCost returns the cost of the edge between start and end nodes
func (d *DirectedGraph) GetEdgeCost(start, end string) (int, error) {
	if _, ok := d.adjList[start]; !ok {
		return -1, errors.New("start node with key not in graph")
	}

	if _, ok := d.adjList[end]; !ok {
		return -1, errors.New("end node with key not in graph")
	}

	for _, edge := range d.adjList[start].edges {
		if edge.Dest.Key == end {
			return edge.Cost, nil
		}
	}
	return -1, errors.New("no edge exists between nodes")
}

// Chain is a method chaining struct for a directed graph.
type Chain struct {
	g   *DirectedGraph
	err error
}

func NewDirectedGraphChain() *Chain {
	g := NewDirectedGraph()

	return &Chain{g, nil}
}

// Err returns the error from the chain, if any exist
func (c *Chain) Err() error {
	return c.err
}

// DirectedGraph returns the directed graph built by the chain, or the error if any occurred in building.
func (c *Chain) DirectedGraph() (*DirectedGraph, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.g, nil
}

// AddNode calls AddNode on the graph
func (c *Chain) AddNode(n *Node) *Chain {
	if c.err != nil {
		return c
	}

	err := c.g.AddNode(n)
	if err != nil {
		c.err = err
	}

	return c
}

// AddEdge calls addEdge on the graph
func (c *Chain) AddEdge(start, end string, cost int) *Chain {
	if c.err != nil {
		return c
	}

	err := c.g.AddEdge(start, end, cost)
	if err != nil {
		c.err = err
	}

	return c
}
