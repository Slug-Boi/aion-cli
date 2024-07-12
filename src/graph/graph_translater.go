package graph

import (
	"github.com/Slug-Boi/aion-cli/forms"
)

type User struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Votes []int  `json:"poll_votes"`
}
//TODO: Add nodeToUser to the return when we need to use it

// Translates data from the forms package to the graph package
func Translate(data forms.Form) ([]Edge, int) {
	nodeToUser := map[int]User{}
	userNodeInc := 1
	intialTimeslotNodeInc := len(data.PollResults) + 1
	timeslotNodeInc := intialTimeslotNodeInc

	graph := []Edge{}

	// Translate participants to source linked nodes
	for _, participant := range data.PollResults {
		// Add edge from source to participant
		graph = append(graph, Edge{From: 0, To: userNodeInc, Capacity: 10, Cost: 10})

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

			graph = append(graph, Edge{From: userNodeInc, To: timeslotNodeInc, Capacity: cap, Cost: cap})
			timeslotNodeInc++
		}

		nodeToUser[userNodeInc] = participant
		userNodeInc++
	}

	// Link timeslots to sink
	for i := intialTimeslotNodeInc; i < timeslotNodeInc; i++ {
		graph = append(graph, Edge{From: i, To: timeslotNodeInc, Capacity: 1, Cost: 1})
	}
	return graph, timeslotNodeInc
}
