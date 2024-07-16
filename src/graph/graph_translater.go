package graph

import (
	"math/rand"
	"strings"
	"time"

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

		// Translate timeslots to participant linked nodes
		for _, timeslot := range participant.Votes {
			// Add edge from participant to timeslot
			var cap float64
			if timeslot == 0 {
				cap = 5.0
			} else if timeslot == 2 {
				cap = 3.0
			} else {
				cap = 1.0
			}
			if cap == cap+floats[i] {
				println("yes")
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

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
