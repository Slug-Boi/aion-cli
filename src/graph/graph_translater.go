package graph

import (
	"strings"

	"github.com/Slug-Boi/aion-cli/forms"
)

type User struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Votes []int  `json:"poll_votes"`
}

// Translates data from the forms package to the graph package
func Translate(data forms.Form) ([]Edge, int, map[int]User) {
	nodeToUser := map[int]User{}
	userNodeInc := 1
	intialTimeslotNodeInc := len(data.PollResults) + 1
	timeslotNodeInc := intialTimeslotNodeInc

	// Create heuristics for tie breaker
	sb := strings.Builder{}
	for i := 1; i < data.Participant_count; i++ {
		sb.WriteString(strings.Split(data.PollResults[i].Name, "")[1])
	}
	allStrings := sb.String()
	floats := []float64{}

	for i := 0; i < data.Participant_count; i++ {
		floats = append(floats, HashHeuristic(data.PollResults[i].Name, allStrings))
	}

	graph := []Edge{}

	// Translate participants to source linked nodes
	for i, participant := range data.PollResults {
		// Add edge from source to participant
		graph = append(graph, Edge{From: 0, To: userNodeInc, Capacity: 1, Cost: 1})

		// Translate timeslots to participant linked nodes
		for _, timeslot := range participant.Votes {
			// Add edge from participant to timeslot
			var cap int
			if timeslot == 0 {
				cap = 5
			} else if timeslot == 2 {
				cap = 3
			} else {
				cap = 1
			}

			graph = append(graph, Edge{From: userNodeInc, To: timeslotNodeInc, Capacity: 1, Cost: cap})
			timeslotNodeInc++
		}

		nodeToUser[userNodeInc] = participant
		userNodeInc++
	}

	// Link timeslots to sink
	for i := intialTimeslotNodeInc; i < timeslotNodeInc; i++ {
		graph = append(graph, Edge{From: i, To: timeslotNodeInc, Capacity: 1, Cost: 1})
	}
	return graph, timeslotNodeInc, nodeToUser
}
