package graph

import (
	// https://pkg.go.dev/github.com/golang-collections/go-datastructures/queue Queue library
	"github.com/golang-collections/go-datastructures/queue"
)

// MinCostPath returns the minimum cost path from the start node to the end node
// It is based on this implmentation https://cp-algorithms.com/graph/min_cost_flow.html using SPFA

type Edge struct {
	From, To, Capacity, Cost int
}

var adjacency, cost, capacity [][]int

// https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
var inf = int(^uint(0) >> 1)

// Shortest_paths returns the shortest path from the start node to all other nodes????
// Variable explanation
// n - number of nodes
// v0 - start node
// d - distance from start node to all other nodes!!!
// p - parent node of each node!!!
func shortest_paths(n, v0 int, d, p []int) {

	// Assign all d values to infinity
	for i := range d {
		d[i] = inf
	}

	// The distance from the start node to itself is 0
	d[v0] = 0

	// Golang defaults to false for []bool
	// This should be the same as visited
	inq := make([]bool, n)

	// Create a queue of size 0
	q := queue.New(0)

	q.Put(v0)

	// Assign all p values to -1
	for i := range p {
		p[i] = -1
	}

	// While the queue is not empty
	for q.Len() > 0 {
		// Get the first element in the queue
		uu, err := q.Get(1)
		if err != nil {
			panic(err)
		}

		// Cast the first element to an int and assign it to u
		u := uu[0].(int)

		inq[u] = false

		// For each vertix in the adjacency list
		for v := range adjacency[u] {
			if capacity[u][v] > 0 && d[v] > d[u]+cost[u][v] {
				d[v] = d[u] + cost[u][v]
				p[v] = u
				if !inq[v] {
					q.Put(v)
					inq[v] = true
				}
			}
		}
	}
}

// MinCostPath returns the minimum cost path from the start node to the end node
// Variable explanation
// N - number of nodes
// K - number of paths???
// s - source node
// t - sink node
// edges - slice of edges
func MinCostPath(N, K, s, t int, edges []Edge) int {
	// Assign empty slices to adjacency of size N
	adjacency = make([][]int, N)

	// For each empty slice in cost and capacity assign a slice of size N zeroed out
	cost = make([][]int, N)
	for i := range cost {
		cost[i] = make([]int, N)
	}

	capacity = make([][]int, N)
	for i := range capacity {
		capacity[i] = make([]int, N)
	}

	// For each edge in the edges slice
	for _, e := range edges {
		// Assign the edge to the adjacency list
		adjacency[e.From] = append(adjacency[e.From], e.To)
		adjacency[e.To] = append(adjacency[e.To], e.From)

		// Assign the cost of the edge to the cost list
		cost[e.From][e.To] = e.Cost
		cost[e.To][e.From] = -e.Cost
		// ^This looks to be the residual cost of edges

		// Assign the capacity of the edge to the capacity list
		capacity[e.From][e.To] = e.Capacity
	}

	// Assign the minimum cost path to the shortest path
	flow, cost := 0, 0
	d := make([]int, N)
	p := make([]int, N)

	// While the flow is less than the capacity
	for flow < K {
		shortest_paths(N, s, d, p)
		if d[t] == inf {
			break
		}

		// find max flow on that path
		f := K - flow
		cur := t
		for cur != s {
			f = min(f, capacity[p[cur]][cur])
			cur = p[cur]
		}

		// apply flow
		flow += f
		cost += f * d[t]
		cur = t
		for cur != s {
			capacity[p[cur]][cur] -= f
			capacity[cur][p[cur]] += f
			cur = p[cur]
		}
	}

	if flow < K {
		return -1
	} else {
		return cost
	}

}