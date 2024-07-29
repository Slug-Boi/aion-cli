package graph_test

import (
	"os"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/solvers/graph"
)

// Create json data for form
var data = []byte(`Timestamp,Group Number,Lottery String,12-04-24 Monday [8:00-10:00],12-04-24 Monday [10:00-12:00],15-04-24 Thursday [10:00-12:00],15-04-24 Thursday [14:30-16:30]
26/07/2024 10:50:58,Group 1,JavaBois,Want,Can do,,
26/07/2024 10:51:21,Group 2,PartyInTheSewers,Want,Want,,`)

func cleanup() {
	os.Remove("form.csv")
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

	os.WriteFile("form.csv", data, 0644)

	defer cleanup()

	var conf forms.Config

	form := forms.GetForm(conf, true)

	// Create a graph from form
	g, sink, _, _ := graph.Translate(form)

	// Check number of edges
	if len(g) != 14 {
		t.Error("Expected 14 edges, got", len(g))
	}

	// Check sink value
	if sink != 7 {
		t.Error("Expected sink value of 7, got", sink)
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

	var dataTieBreaker = []byte(`Timestamp,Group Number,Lottery String,12-04-24 Monday [8:00-10:00],12-04-24 Monday [10:00-12:00]
					26/07/2024 10:50:58,Group 1,JavaBois,Want,Can do
					26/07/2024 10:51:21,Group 2,PartyInTheSewers,Want,Can do`)

	os.WriteFile("form.csv", dataTieBreaker, 0644)

	defer cleanup()

	var conf forms.Config

	form := forms.GetForm(conf, true)

	// Create a graph from form
	g, sink, users, nodeToTimeslot := graph.Translate(form)

	// Check heuristic values of the two users:

	// Convert map to slice
	usersSlice := []forms.Form{}
	for _, user := range users {
		usersSlice = append(usersSlice, user)
	}

	// Sort users by id to ensure consistent ordering when generating the concatenated string
	sort.Slice(usersSlice, func(i, j int) bool {
		return usersSlice[i].HashString < usersSlice[j].HashString
	})

	// Generate two heuristics from the two users
	// Heur 1: 0.00000000000003216899887849856
	// Heur 2: 0.000000000000030574509877317846

	sb := strings.Builder{}
	allStrings := graph.BaseHashString(form, sb)
	heur1 := graph.HashHeuristic(usersSlice[0].HashString, allStrings)
	heur2 := graph.HashHeuristic(users[1].HashString, allStrings)

	// Check that heuristic 2 is lesser than heuristic 1
	if heur1 < heur2 {
		t.Error("Expected heuristic 2 to be greater than heuristic 1, got", heur1, heur2)
	}

	// Check number of edges
	if len(g) != 8 {
		t.Error("Expected 14 edges, got", len(g))
	}

	// Check sink value
	if sink != 5 {
		t.Error("Expected sink value of 7, got", sink)
	}

	// Values are:
	// 12 nodes, 2 is minimum flow required, 0 is source, 7 is sink, g is the graph
	_, paths := graph.MinCostPath(len(g), 2, 0, 5, g)


	// Check number of paths
	if len(paths) != 2 {
		t.Error("Expected 2 paths, got", len(paths))
	}

	// Check found paths
	// This indirectly confirms that node 1 always gets the preferred timeslot that being node 3
	for _, path := range paths {
		i := 5
		for i != 0 {
			if nodeToTimeslot[i] == "12-04-24 Monday [10:00-12:00]" {
				if users[path[i]].GroupNumber != "Group 1" {
					t.Error("Expected group 1 to be get Monday 10:00-12:00, instead group: ", users[path[i]].GroupNumber, "got it")
				} 
				if nodeToTimeslot[i] == "12-04-24 Monday [8:00-10:00]" {
					if users[path[i]].GroupNumber != "Group 2" {
						t.Error("Expected group 2 to be get Monday 8:00-10:00, instead group: ", users[path[i]].GroupNumber, "got it")
					}
				}

			}
			i = path[i]
		}
	}

}

//TODO: Needs integration test of the two parts working together might want to do in seperate test folder
