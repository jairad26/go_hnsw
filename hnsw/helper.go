package hnsw

import (
	"math"
	"reflect"
	"sort"

	"github.com/barweiss/go-tuple"
)

func norm(v []float64) float64 {
	var sum float64
	for _, value := range v {
		sum += math.Pow(value, 2)
	}
	return math.Sqrt(sum)
}

func max(a int, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func add_slices(a []float64, b []float64) []float64 {
	var output []float64
	for i := 0; i < max(len(a), len(b)); i++ {
		output = append(output, a[i]+b[i])
	}
	return output
}

func sub_slices(a []float64, b []float64) []float64 {
	var output []float64
	for i := 0; i < max(len(a), len(b)); i++ {
		output = append(output, a[i]-b[i])
	}
	return output
}

func mapkey(m map[int]*Node, value *Node) (key int, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func insort_tuple(slice []tuple.T2[float64, int], val tuple.T2[float64, int]) []tuple.T2[float64, int] {
	i := sort.Search(len(slice), func(i int) bool { return slice[i].V1 > val.V1 })
	slice = append(slice, tuple.New2(0.0, 0))
	copy(slice[i+1:], slice[i:])
	slice[i] = val
	return slice
}

func position(vec_arr [][]float64, vec []float64) int {
	for index, element := range vec_arr {
		if reflect.DeepEqual(element, vec) {
			return index
		}
	}
	return -1
}
