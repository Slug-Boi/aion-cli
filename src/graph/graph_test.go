package graph_test

import (
	"encoding/json"
	"os"
	"slices"
	"testing"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/graph"
)

func cleanup() {
	os.Remove("data.json")
}

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
	g = append(g, graph.Edge{From: 0, To: 1, Capacity: 1, Cost: 0})
	g = append(g, graph.Edge{From: 0, To: 2, Capacity: 1, Cost: 0})
	g = append(g, graph.Edge{From: 0, To: 3, Capacity: 1, Cost: 0})

	g = append(g, graph.Edge{From: 1, To: 4, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 2, To: 4, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 2, To: 5, Capacity: 1, Cost: 2})
	g = append(g, graph.Edge{From: 3, To: 6, Capacity: 1, Cost: 3})
	g = append(g, graph.Edge{From: 3, To: 7, Capacity: 1, Cost: 1})

	g = append(g, graph.Edge{From: 4, To: 8, Capacity: 1, Cost: 0})
	g = append(g, graph.Edge{From: 5, To: 8, Capacity: 1, Cost: 0})
	g = append(g, graph.Edge{From: 6, To: 8, Capacity: 1, Cost: 0})
	g = append(g, graph.Edge{From: 7, To: 8, Capacity: 1, Cost: 0})

	return g
}

func TestMinCost(t *testing.T) {
	g := debugGraphBuilder()

	// Values are:
	// 9 nodes, 3 is minimum flow required, 0 is source, 8 is sink, g is the graph
	cost, paths := graph.MinCostPath(9, 3, 0, 8, g)

	// Check number of paths
	if len(paths) != 3 {
		t.Error("Expected 3 paths, got", len(paths))
	}

	// Check found paths
	p1 := []int{8, 4, 1, 0}
	p2 := []int{8, 7, 3, 0}
	p3 := []int{8, 5, 2, 0}

	for _, path := range paths {
		i := 8
		act_paths := []int{8}
		for i != 0 {
			act_paths = append(act_paths, path[i])
			i = path[i]
		}
		for i := 0; i < 3; i++ {
			if !slices.Equal(act_paths, p1) && !slices.Equal(act_paths, p2) && !slices.Equal(act_paths, p3) {
				t.Error("Path", p1, "is not valid ", act_paths)
			}
		}
	}

	// Check cost

	if int(cost) != 4 {
		t.Error("Expected 4 cost, got", cost)
	}

}

func TestGraphTranslation(t *testing.T) {
	// Create json data for form
	data := []byte(`{"participant_count":2,"poll_options":[{"id":"NPgxbaN4oy2","start_time":1720436400,"end_time":1720440000},{"id":"wAg39ORa8y8","start_time":1720440000,"end_time":1720443600},{"id":"6QnMoXKEVZe","start_time":1720443600,"end_time":1720447200},{"id":"NoZr4wk7Dn3","start_time":1720447200,"end_time":1720450800}],"poll_participants":[{"name":"Group 4","id":"Jnv72xxzmgv","poll_votes":[1,2,0,2]},{"name":"Group 5","id":"jn1jJllGLgQ","poll_votes":[2,1,1,0]}]}`)

	var form forms.Form

	// Unmarshal json data
	err := json.Unmarshal(data, &form)
	if err != nil {
		t.Error("Error unmarshalling json data")
	}

	// Create a graph from form
	g, sink, _ := graph.Translate(form)

	// Check number of edges
	if len(g) != 18 {
		t.Error("Expected 18 edges, got", len(g))
	}

	// Check sink value
	if sink != 11 {
		t.Error("Expected sink value of 11, got", sink)
	}

}

//TODO: Needs integration test of the two parts working together might want to do in seperate test folder
