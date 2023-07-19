package hnsw

import (
	// "fmt"
	"math"
	"math/rand"

	"github.com/barweiss/go-tuple"
)

type HNSW struct {
	Index          []*Graph
	max_levels     int
	mult_factor    float64
	efConstruction int
	max_neighbors  int
}

func HNSW_init(max_levels int, mult_factor float64, efConstruction int, max_neighbors int) *HNSW {
	hnsw := HNSW{
		Index:          []*Graph{},
		max_levels:     max_levels,
		mult_factor:    mult_factor,
		efConstruction: efConstruction,
		max_neighbors:  max_neighbors,
	}
	for i := 0; i < hnsw.max_levels; i++ {
		graph := Graph{
			Vertices: map[int]*Node{},
		}
		hnsw.Index = append(hnsw.Index, &graph)
	}
	return &hnsw
}

func (hnsw *HNSW) Create(vectors [][]float64, uuids []uint64) {
	for i, _ := range vectors {
		hnsw.Insert(vectors[i], uuids[i])
	}
}

func (hnsw *HNSW) get_insert_layer() int {
	// mult_factor is a multiplicative factor used to normalize the distribution
	var level int
	for i := range hnsw.Index {
		// calculate level based on section 3.1 here https://www.pdl.cmu.edu/PDL-FTP/BigLearning/mod0246-liA.pdf
		if rand.Float64() < math.Pow(1.0/float64(hnsw.max_neighbors), float64(hnsw.max_levels-1-i)) {
			level = i
			break
		}
	}
	return level
	// fmt.Println(math.Log(rand.Float64()) * hnsw.mult_factor)
	// level := -int(math.Log(rand.Float64()) * hnsw.mult_factor)
	// return min(level, hnsw.max_levels)
}

func (hnsw *HNSW) Insert(vec []float64, uuid uint64) []*Graph {
	if len(hnsw.Index[0].Vertices) == 0 {
		next_layer := -1
		for i := len(hnsw.Index) - 1; i >= 0; i-- {
			node := Node{
				Vector:     vec,
				next_layer: next_layer,
			}
			hnsw.Index[i].Vertices[0] = &node
			next_layer = 0
		}
		return hnsw.Index
	}
	level := hnsw.get_insert_layer()
	start_v := 0

	for i, graph := range hnsw.Index {
		// perform insertion for layers [level, max_level) only
		if i < level {
			// fmt.Println("THIS IS AN INSERT SEARCH_LAYER")
			start_v = search_layer(graph, start_v, vec, 1)[0].V2
		} else {
			var node Node
			node.Vector = vec
			if i < hnsw.max_levels-1 {
				node.next_layer = len(hnsw.Index[i+1].Vertices)
			} else {
				node.next_layer = -1
			}
			nns := search_layer(graph, start_v, vec, hnsw.max_neighbors) //check up to efConstruction neighbors, only use closest ones up to max_neighbors
			for _, nn := range nns {
				node.edges = append(node.edges, graph.Vertices[nn.V2])                   // outbound edge
				graph.Vertices[nn.V2].edges = append(graph.Vertices[nn.V2].edges, &node) // inbound edge
			}
			graph.Vertices[len(graph.Vertices)] = &node
		}
		start_v = graph.Vertices[start_v].next_layer
	}

	return hnsw.Index
}

/*
Implement priority queue using heap to order nearest neighbor vectors in graph. Using euclidian distance to identify nearest neighbors.
nns: output list of nearest neighbors
candid: heap of candidate nodes
evaluate all nearest neighbors against the best (closest) vector in candid, updating candid & nns as you go.
stop when there are no more candidate points to evaluate, or when you know you can't do any better in this layer
*/
func search_layer(graph *Graph, entry int, query []float64, expected_neighbors int) []tuple.T2[float64, int] {
	//create a new tuple (vector_dist, graph_index)
	best := tuple.New2(euclidian_distance(graph.Vertices[entry].Vector, query), entry)
	nns := []tuple.T2[float64, int]{best}

	//create set using map to append to on future visited nodes
	visited := map[tuple.T2[float64, int]]bool{best: true}

	candidateHeap := *buildHeapByInit([]tuple.T2[float64, int]{best})
	// fmt.Println(candidateHeap)

	for candidateHeap.Len() != 0 {
		curr_candidate := candidateHeap.Pop().(tuple.T2[float64, int])
		// fmt.Println(curr_candidate)
		if nns[len(nns)-1].V1 < curr_candidate.V1 {
			break
		}
		for _, node := range graph.Vertices[curr_candidate.V2].edges {
			curr_dist := euclidian_distance(node.Vector, query)
			curr_key, key_exists := mapkey(graph.Vertices, node)
			if !key_exists {
				panic("value does not exist in map")
			}
			curr_tuple := tuple.New2(curr_dist, curr_key)
			_, node_exists := visited[curr_tuple]
			if !node_exists {
				visited[curr_tuple] = true

				// push only better vectors into candidate heap and add to nearest neighbors
				if curr_dist < nns[len(nns)-1].V1 || len(nns) < expected_neighbors {
					candidateHeap.Push(curr_tuple)
					nns = insort_tuple(nns, curr_tuple)
					if len(nns) > expected_neighbors {
						nns = nns[:len(nns)-1]
					}
				}

			}
		}
	}
	// fmt.Println(nns)
	return nns
}

func Search(index []*Graph, query []float64, expected_neighbors int) []tuple.T2[float64, int] {
	if len(index[0].Vertices) > 0 {
		best_v := 0
		for _, graph := range index {
			// fmt.Println("--------------------NEW GRAPH-----------------")
			// for _, node := range graph.vertices {
			// 	fmt.Println(node.next_layer)
			// }
			curr_best := search_layer(graph, best_v, query, expected_neighbors)[0]
			// fmt.Println(curr_best.V2)
			if graph.Vertices[curr_best.V2].next_layer != -1 {
				best_v = graph.Vertices[curr_best.V2].next_layer
			}
			// fmt.Println(best_v)
		}
		return search_layer(index[len(index)-1], best_v, query, expected_neighbors)
	} else {
		return []tuple.T2[float64, int]{}
	}
}
