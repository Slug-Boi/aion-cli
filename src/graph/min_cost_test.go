package graph_test

import (
	"slices"
	"testing"

	"github.com/Slug-Boi/aion-cli/graph"
)

// This is a debugging graph

func debugGraphBuilder() []graph.Edge {
	// Edge values
	// From, To, Capacity, Cost
	g := []graph.Edge{}

	// 0 source
	// 8 sink
	// Groups -> 1, 2, 3
	// Timeslots -> 4, 5, 6, 7
	// Add edges to the graph
	g = append(g, graph.Edge{From: 0, To: 1, Capacity: 100, Cost: 1})
	g = append(g, graph.Edge{From: 0, To: 2, Capacity: 100, Cost: 1})
	g = append(g, graph.Edge{From: 0, To: 3, Capacity: 100, Cost: 1})

	g = append(g, graph.Edge{From: 1, To: 4, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 2, To: 4, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 2, To: 5, Capacity: 2, Cost: 2})
	g = append(g, graph.Edge{From: 3, To: 6, Capacity: 3, Cost: 3})
	g = append(g, graph.Edge{From: 3, To: 7, Capacity: 1, Cost: 1})

	g = append(g, graph.Edge{From: 4, To: 8, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 5, To: 8, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 6, To: 8, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 7, To: 8, Capacity: 1, Cost: 1})

	return g
}

func TestMinCost(t *testing.T) {
	g := debugGraphBuilder()

	groups := 3

	// Values are:
	// 9 nodes, 3 is minimum flow required, 0 is source, 8 is sink, g is the graph
	cost, paths := graph.MinCostPath(9, 3, 0, 8, g)

	// Check number of paths
	if len(paths) != 3 {
		t.Error("Expected 3 paths, got", len(paths))
	}

	// Check found paths
	p1 := []int{8,4,1,0}
	p2 := []int{8,7,3,0}
	p3 := []int{8,5,2,0}

	for _, path := range paths {
			i := 8
			act_paths := []int{8}
			for i != 0 {
				act_paths = append(act_paths, path[i])
				i = path[i]
			}
			for i := 0 ; i < 3 ; i++{
				if !slices.Equal(act_paths, p1) && !slices.Equal(act_paths, p2) && !slices.Equal(act_paths, p3) {
				t.Error("Path", p1, "is not valid ", act_paths)
				}
			}	
		}
	
	// Check cost
	cost = ((cost - groups) - len(paths))
	if cost != 4 {
		t.Error("Expected 4 cost, got", cost)
	}

}
