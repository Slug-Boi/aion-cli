package graph

import (
	"hash/fnv"
	"math"
	"math/rand"

	"github.com/golang-collections/go-datastructures/queue"
)

// MinCostPath returns the minimum cost path from the start node to the end node
// It is based on this implmentation https://cp-algorithms.com/graph/min_cost_flow.html using SPFA

type Edge struct {
	From, To, Capacity int
	Cost               float64
}

var adjacency, capacity [][]int
var cost [][]float64

// https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go
// var inf = int(^uint(0) >> 1)
var inf = math.MaxFloat64

// Shortest_paths returns the shortest path from the start node to all other nodes????
// Variable explanation
// n - number of nodes
// v0 - start node
// d - distance from start node to all other nodes ??? (unsure about this)
// p - the path walked to get to the sink. This is used to backtrack the path
func shortest_paths(n, v0 int, d *[]float64, p *[]int) {
	// Assign a slice of size n to d
	*d = make([]float64, n)

	// Assign a slice of size n to p
	*p = make([]int, n)

	// Assign all d values to infinity
	for i := 0; i < n; i++ {
		(*d)[i] = inf
	}

	// The distance from the start node to itself is 0
	(*d)[v0] = 0

	// Golang defaults to false for []bool
	// This should be the same as visited
	inq := make([]bool, n)

	// Create a queue of size 0
	q := queue.New(int64(n))

	q.Put(v0)

	// Assign all p values to -1
	for i := 0; i < n; i++ {
		(*p)[i] = -1
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
		for _, v := range adjacency[u] {
			if capacity[u][v] > 0 && (*d)[v] > ((*d)[u]+cost[u][v]) {
				(*d)[v] = (*d)[u] + cost[u][v]
				(*p)[v] = u
				if !inq[v] {
					inq[v] = true
					q.Put(v)
				}
			}
		}
	}
}

// MinCostPath returns the minimum cost path from the start node to the end node
// Variable explanation
// N - number of nodes
// K - the minimum required flow. Setting this to group count seems to work well
// s - source node
// t - sink node
// edges - slice of edges
func MinCostPath(N, K, s, t int, edges []Edge) (float64, [][]int) {
	// Assign a path variable to backtrack later
	paths := [][]int{}

	// Assign empty slices to adjacency of size N
	adjacency = make([][]int, N)

	// For each empty slice in cost and capacity assign a slice of size N zeroed out
	cost = make([][]float64, N)
	capacity = make([][]int, N)

	for i := 0; i < N; i++ {
		cost[i] = make([]float64, N)
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
	flow, cost := 0, 0.0

	var p []int
	var d []float64
	// While the flow is less than the capacity
	for flow < K {
		shortest_paths(N, s, &d, &p)
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
		cost += float64(f) * d[t]
		cur = t
		for cur != s {
			//println("cur: ",cur,"\nbefore cap:",capacity[p[cur]][cur])
			capacity[p[cur]][cur] -= f
			//println("after cap:",capacity[p[cur]][cur])
			capacity[cur][p[cur]] += f
			cur = p[cur]
		}
		paths = append(paths, p)
	}

	if flow < K {
		return -1, [][]int{}
	} else {
		return cost, paths
	}

}

// TODO: Figure out if this is doable with a rolling hash function
func HashHeuristic(groupHash, FullHash string) float64 {
	// Combine the two hash strings from input
	combined_str := groupHash + FullHash

	// convert to byte array
	combined := []byte(combined_str)

	// Create hash value of 32 bits
	hasher := fnv.New32a()
	hasher.Write(combined)
	hash := hasher.Sum32()

	// The hash is used to seed the random number generator resulting in the same number every time
	random := rand.New(rand.NewSource(int64(hash)))

	// bound the random number between 0 and 0.5
	random_float := (random.Float64() * 0.5) + 0

	// convert the hash to a binary string of 10 bits by shifting
	// Then we convert the binary string to a float64 for the heuristic
	// The binary number has 54 0s in front of it to make it a decimal number of minimal size
	// This is to make the heuristic as small as possible to avoid messing with the flow algorithm

	//NOTE: This code is extemely cursed and broken I will start with random library and maybe circle back
	// To this later if I have time
	// parsed, _ := strconv.ParseUint(binaryConvert(hash), 2, 64)
	// println(parsed)

	//float := math.Float64frombits(parsed)

	return random_float
}

// func binaryConvert(n uint32) string {
// 	// Shift the number to the right until it is less than 1024 to ensure it is 10 bits or less
// 	for {
// 		if n > 1024 {
// 			n = n >> 1
// 			// or shift the number to the left until it is 10 bits
// 		} else if len(strconv.FormatInt(int64(n), 2)) < 10 {
// 			n = n << 1
// 		} else {
// 			break
// 		}
// 	}
// 	// Convert the binary number to a very small decimal value with 54 0s in front
// 	//lead := "00000"

// 	return strconv.FormatInt(int64(n), 2)
// }
