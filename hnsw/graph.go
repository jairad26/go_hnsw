package hnsw

type Node struct {
	Vector     []float64
	edges      []*Node
	Uuid       uint64
	next_layer int
}

func (n *Node) add_value(vec []float64) {
	n.Vector = vec
}

func (n *Node) add_edge(node *Node) {
	n.edges = append(n.edges, node)
}

type Graph struct {
	Vertices map[int]*Node
}
