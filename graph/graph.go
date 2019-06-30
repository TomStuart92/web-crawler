// Package graph provides an implementation of a directed graph, based upon simple string labels
package graph

import (
	"fmt"
	"io"
)

// Node is a single object in a directed graph
type Node struct {
	label string
}

// Graph implements a directed graph with BFS functionality
type Graph struct {
	nodes map[string]*Node
	edges map[Node][]*Node
}

// NewGraph returns a new instance of Graph
func NewGraph() *Graph {
	return &Graph{make(map[string]*Node), make(map[Node][]*Node)}
}

// AddNode creates a new node based on the given label, or does nothing if the node already exists
func (g *Graph) AddNode(label string) {
	if g.HasNode(label) {
		return
	}
	node := Node{label}
	g.nodes[label] = &node
}

// HasNode returns a boolean based on whether the given string label has been added as a node
func (g *Graph) HasNode(label string) bool {
	_, ok := g.nodes[label]
	return ok
}

// AddEdge creates an edge from the from node to the to node.
// returns an error if either node has not be created before.
func (g *Graph) AddEdge(from string, to string) error {
	fromNode := g.nodes[from]
	toNode := g.nodes[to]
	if fromNode == nil {
		return fmt.Errorf("%s Node does not exist", from)
	}
	if toNode == nil {
		return fmt.Errorf("%s Node does not exist", to)
	}
	g.edges[*fromNode] = append(g.edges[*fromNode], toNode)
	return nil
}

// BFS performs a breadth-first search of the given graph,
// starting at the given labelled node, and
// calling the provided function on each edge in the graph.
func (g *Graph) BFS(start string, f func(string, string)) error {
	queue := make([]*Node, 1)
	queue[0] = g.nodes[start]

	if queue[0] == nil {
		return fmt.Errorf("%s Node does not exist", start)
	}

	visited := make(map[*Node]bool)

	for {
		if len(queue) == 0 {
			return nil
		}
		node := queue[0]
		queue = queue[1:len(queue)]
		visited[node] = true

		edges := g.edges[*node]

		for i := 0; i < len(edges); i++ {
			neighbor := edges[i]
			f(node.label, neighbor.label)
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}
}

// Print is a helper wrapper around BFS, which logs a record for each edge discovered in the BFS
func (g *Graph) Print(output io.Writer, start string) error {
	f := func(from string, to string) {
		fmt.Fprintf(output, "\n%s ---> %s\n", from, to)
	}
	return g.BFS(start, f)
}
