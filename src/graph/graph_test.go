package graph_test

import (
	"encoding/json"
	"slices"
	"sort"
	"testing"

	"github.com/Slug-Boi/aion-cli/forms"
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

	var form []forms.Form

	// Unmarshal json data
	err := json.Unmarshal(data, &form)
	if err != nil {
		t.Error("Error unmarshalling json data")
	}

	// Create a graph from form
	g, sink, _, _ := graph.Translate(form)

	// Check number of edges
	if len(g) != 18 {
		t.Error("Expected 18 edges, got", len(g))
	}

	// Check sink value
	if sink != 11 {
		t.Error("Expected sink value of 11, got", sink)
	}

}

func TestHashHeuristic(t *testing.T) {
	// Check hash heuristic
	heuristic := graph.HashHeuristic("Group 4", "Group 4Group 5")

	if heuristic > 0.00005 {
		t.Error("Expected something greater than 0 heuristic, got", heuristic)
	}

	heuristic2 := graph.HashHeuristic("Group 5", "Group 4Group 5")
	if heuristic2 > 0.00005 {
		t.Error("Expected something greater than 0 heuristic2, got", heuristic)
	}
	if heuristic == heuristic2 {
		t.Error("heuristic and heuristic2 are equal to the same\nHeuristic:", heuristic, "\nHeuristic2:", heuristic2)
	}

}

func TestGraphTieBreaking(t *testing.T) {
	// Create json data for form
	data := []byte(`{"participant_count":2,"poll_options":[{"id":"NPgxbaN4oy2","start_time":1720436400,"end_time":1720440000},{"id":"wAg39ORa8y8","start_time":1720440000,"end_time":1720443600}],"poll_participants":[{"name":"4","id":"Jnv72xxzmgv","poll_votes":[1]},{"name":"5","id":"jn1jJllGLgQ","poll_votes":[1]}]}`)

	var form forms.Form

	// Unmarshal json data
	err := json.Unmarshal(data, &form)
	if err != nil {
		t.Error("Error unmarshalling json data")
	}

	// Create a graph from form
	g, sink, users := graph.Translate(form)

	// Check heuristic values of the two users:

	// Convert map to slice
	usersSlice := make([]graph.User, 0)
	for _, user := range users {
		usersSlice = append(usersSlice, user)
	}

	// Sort users by id to ensure consistent ordering when generating the concatenated string
	sort.Slice(usersSlice, func(i, j int) bool {
		return users[i].Id < users[j].Id
	})

	// Generate two heuristics from the two users
	allStrings := users[0].Id + users[1].Id
	heur1 := graph.HashHeuristic(users[0].Id, allStrings)
	heur2 := graph.HashHeuristic(users[1].Id, allStrings)

	// Check that heuristic 1 is lesser than heuristic 2
	if heur1 < heur2 {
		t.Error("Expected heuristic 1 to be greater than heuristic 2, got", heur1, heur2)
	}

	// Check number of edges
	if len(g) != 6 {
		t.Error("Expected 6 edges, got", len(g))
	}

	// Check sink value
	if sink != 5 {
		t.Error("Expected sink value of 5, got", sink)
	}

	// Values are:
	// 12 nodes, 2 is minimum flow required, 0 is source, 11 is sink, g is the graph
	_, paths := graph.MinCostPath(len(g), 2, 0, 5, g)

	// Check number of paths
	if len(paths) != 2 {
		t.Error("Expected 2 paths, got", len(paths))
	}

	// Check found paths
	// This indirectly confirms that node 1 always gets the preferred timeslot that being node 3
	p1 := []int{5, 4, 2, 0}
	p2 := []int{5, 3, 1, 0}

	for _, path := range paths {
		i := 5
		act_paths := []int{5}
		for i != 0 {
			act_paths = append(act_paths, path[i])
			i = path[i]
		}
		for i := 0; i < 2; i++ {
			if !slices.Equal(act_paths, p1) && !slices.Equal(act_paths, p2) {
				t.Error("Path", p1, "is not valid ", act_paths)
			}
		}
	}



}

//TODO: Needs integration test of the two parts working together might want to do in seperate test folder
