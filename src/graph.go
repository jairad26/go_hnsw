package main

type Node struct {
	vector     []float64
	edges      []*Node
	uuid       int
	next_layer int
}

func (n *Node) add_value(vec []float64) {
	n.vector = vec
}

func (n *Node) add_edge(node *Node) {
	n.edges = append(n.edges, node)
}

type Graph struct {
	vertices map[int]*Node
}
