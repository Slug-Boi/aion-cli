package graph

import (
	"sort"
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
	for i := 0; i < data.Participant_count; i++ {
		//TODO: Probably redo this to be more modular
		split := strings.Fields(data.PollResults[i].Name)
		if len(split) > 1 {
			//TODO: Add a regex check for valid characters
			sb.WriteString(split[1])
		} else {
			// If no string is provided then use their id as it is a pseudo random string
			// that will result in the same value each time
			sb.WriteString(data.PollResults[i].Id)
		}
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
		graph = append(graph, Edge{From: 0, To: userNodeInc, Capacity: 1, Cost: 0})


		// These variables are used to keep track of the total budget of the group
		// And the amount of leftovers nodes once we hit "red" timeslots 
		// The list is sorted so we know all worst timeslots will be at the end
		budget := 50.0
		leftovers := 0.0

		sort.Slice(participant.Votes, func(i, j int) bool {
			return participant.Votes[i] > participant.Votes[j]
		})

		// Translate timeslots to participant linked nodes
		for j, timeslot := range participant.Votes {

			// Add edge from participant to timeslot
			var cap float64

			if timeslot == 0 {
				// Calculate the divider for the budget based on the amount of timeslots left
				if leftovers == 0 {
					leftovers = float64(len(participant.Votes)-j)
				}
				cap = (budget / leftovers)
			} else if timeslot == 2 {
				cap = 5.0
				budget = budget - 5.0
			} else {
				cap = 0.0
			}

			graph = append(graph, Edge{From: userNodeInc, To: timeslotNodeInc, Capacity: 1, Cost: cap + floats[i]})
			timeslotNodeInc++
		}

		nodeToUser[userNodeInc] = participant
		userNodeInc++
	}

	// Link timeslots to sink
	for i := intialTimeslotNodeInc; i < timeslotNodeInc; i++ {
		graph = append(graph, Edge{From: i, To: timeslotNodeInc, Capacity: 1, Cost: 0})
	}
	return graph, timeslotNodeInc, nodeToUser
}
