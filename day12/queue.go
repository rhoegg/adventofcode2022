package main

import "container/heap"

type WeightedLocation struct {
	Location
	Weight int
	index int
}

type LocationPQ []*WeightedLocation

func (pq LocationPQ) Len() int {
	return len(pq)
}

func (pq LocationPQ) Less(i, j int) bool {
	return pq[i].Weight < pq[j].Weight
}

func (pq LocationPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *LocationPQ) Push(item interface{}) {
	n := len(*pq)
	location := item.(*WeightedLocation)
	location.index = n
	*pq = append(*pq, location)
}

func (pq *LocationPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the weight of a WeightedLocation in the queue.
func (pq *LocationPQ) update(item *WeightedLocation, weight int) {
	item.Weight = weight
	heap.Fix(pq, item.index)
}
