package main

import (
	"container/heap"

	"github.com/barweiss/go-tuple"
)

type MinTupleHeap []tuple.T2[float64, int]

func (h MinTupleHeap) Len() int {
	return len(h)
}

func (h MinTupleHeap) Less(i, j int) bool {
	return h[i].V1 < h[j].V1
}

func (h MinTupleHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinTupleHeap) Push(x interface{}) {
	*h = append(*h, x.(tuple.T2[float64, int]))
}

func (h *MinTupleHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Time: O(n)
func buildHeapByInit(array []tuple.T2[float64, int]) *MinTupleHeap {
	// initialize the MinTupleHeap that has implement the heap.Interface
	MinTupleHeap := &MinTupleHeap{}
	*MinTupleHeap = array
	heap.Init(MinTupleHeap)
	return MinTupleHeap
}

// Time: O(nlogn)
func buildHeapByPush(array []int) *MinTupleHeap {
	// initialize the MinTupleHeap that has implement the heap.Interface
	MinTupleHeap := &MinTupleHeap{}
	for _, num := range array {
		heap.Push(MinTupleHeap, num)
	}
	return MinTupleHeap
}
