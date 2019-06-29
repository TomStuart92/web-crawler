package graph_test

import (
	"testing"

	"github.com/TomStuart92/web-crawler/graph"
)

func TestNewGraph(t *testing.T) {
	g := graph.NewGraph()
	if g == nil {
		t.Error("Expected #NewGraph to return new graph")
	}
}

func TestAddNodeThenHasNode(t *testing.T) {
	label := "test"
	g := graph.NewGraph()
	g.AddNode(label)
	hasNode := g.HasNode(label)
	if !hasNode {
		t.Error("Expected #HasNode to return true, after node is added")
	}
}

func TestHasNodeWithNoNode(t *testing.T) {
	label := "test"
	g := graph.NewGraph()
	hasNode := g.HasNode(label)
	if hasNode {
		t.Error("Expected #HasNode to return false, if node has not been added")
	}
}

func TestHasNodeWithWithDuplicateAddNodes(t *testing.T) {
	label := "test"
	g := graph.NewGraph()
	g.AddNode(label)
	g.AddNode(label)
	hasNode := g.HasNode(label)
	if !hasNode {
		t.Error("Expected #HasNode to return true, after node is added")
	}
}

func TestAddEdgeWithValidNodes(t *testing.T) {
	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	g.AddNode(label1)
	g.AddNode(label2)
	err := g.AddEdge(label1, label2)
	if err != nil {
		t.Error(err)
	}
}

func TestAddEdgeWithoutAddingFromNode(t *testing.T) {
	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	g.AddNode(label2)
	err := g.AddEdge(label1, label2)
	if err == nil {
		t.Error("Expected Error after trying to add edge to non-existant node")
	}
}

func TestAddEdgeWithoutAddingToNode(t *testing.T) {
	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	g.AddNode(label1)
	err := g.AddEdge(label1, label2)
	if err == nil {
		t.Error("Expected Error after trying to add edge to non-existant node")
	}
}

func TestAddEdgeWithoutAddingAnyNode(t *testing.T) {
	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	err := g.AddEdge(label1, label2)
	if err == nil {
		t.Error("Expected Error after trying to add edge to non-existant node")
	}
}

func TestBFSWithSimpleSetup(t *testing.T) {
	var edges [][]string

	f := func(from string, to string) {
		testPair := []string{from, to}
		edges = append(edges, testPair)
	}

	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	g.AddNode(label1)
	g.AddNode(label2)
	err := g.AddEdge(label1, label2)

	if err != nil {
		t.Error(err)
		return
	}

	err = g.BFS(label1, f)

	if err != nil {
		t.Error(err)
		return
	}

	if len(edges) != 1 {
		t.Errorf("Expected 1 edge, got %d", len(edges))
		return
	}
	if edges[0][0] != label1 {
		t.Errorf("Expected edge 1 from label to be %s, got %s", label1, edges[0][0])
		return
	}
	if edges[0][1] != label2 {
		t.Errorf("Expected edge 1 to label to be %s, got %s", label2, edges[0][1])
		return
	}
}

func TestBFSWithNoEdges(t *testing.T) {
	var edges [][]string

	f := func(from string, to string) {
		testPair := []string{from, to}
		edges = append(edges, testPair)
	}

	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	g.AddNode(label1)
	g.AddNode(label2)

	err := g.BFS(label1, f)

	if err != nil {
		t.Error(err)
	}

	if len(edges) != 0 {
		t.Errorf("Expected 0 edges, got %d", len(edges))
	}
}

func TestBFSWithNoNodes(t *testing.T) {
	var edges [][]string

	f := func(from string, to string) {
		testPair := []string{from, to}
		edges = append(edges, testPair)
	}

	label1 := "node1"

	g := graph.NewGraph()
	err := g.BFS(label1, f)

	if err == nil {
		t.Error("Expected Node does not exist error, got nil")
	}
}

func TestBFSWithUnknownNode(t *testing.T) {
	var edges [][]string

	f := func(from string, to string) {
		testPair := []string{from, to}
		edges = append(edges, testPair)
	}

	label1 := "node1"
	label2 := "node2"
	g := graph.NewGraph()
	g.AddNode(label1)

	err := g.BFS(label2, f)

	if err == nil {
		t.Error("Expected Node does not exist error, got nil")
	}
}
