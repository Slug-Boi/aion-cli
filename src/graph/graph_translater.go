package graph

import (
	"hash/fnv"
	"math/rand"
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
	// Create a map to store the user node to user data
	nodeToUser := map[int]User{}
	// Users are always node 1 to len(data.PollResults)
	userNodeInc := 1
	// Timeslots are always node len(data.PollResults) + 1 to len(data.PollResults) + len(participants.Votes)
	intialTimeslotNodeInc := len(data.PollResults) + 1
	// Start at the initial timeslot node for the incrementor
	timeslotNodeInc := intialTimeslotNodeInc

	// Cache is used to store the group name going to the ID for consistent hashing later
	cache := map[string]string{}

	// Sort users by id to ensure consistent ordering when generating the concatenated string
	sort.Slice(data.PollResults, func(i, j int) bool {
		return data.PollResults[i].Id < data.PollResults[j].Id
	})

	// Create base string for creating heuristic
	sb := strings.Builder{}
	for i := 0; i < data.Participant_count; i++ {
		//TODO: Probably redo this to be more modular
		split := strings.Fields(data.PollResults[i].Name)
		if len(split) > 1 {
			//TODO: Add a regex check for valid characters
			cache[data.PollResults[i].Name] = split[1]
			sb.WriteString(split[1])
		} else {
			// If no string is provided then use their id as it is a pseudo random string
			// that will result in the same value each time
			cache[data.PollResults[i].Name] = data.PollResults[i].Id
			sb.WriteString(data.PollResults[i].Id)
		}
	}
	// Get the string of all group inputs
	allStrings := sb.String()

	graph := []Edge{}

	// Translate participants to source linked nodes
	for _, participant := range data.PollResults {

		heuristic := HashHeuristic(cache[participant.Name], allStrings)

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

			graph = append(graph, Edge{From: userNodeInc, To: timeslotNodeInc, Capacity: 1, Cost: cap + heuristic})
			timeslotNodeInc++
		}

		// Add user to map and increment user node
		nodeToUser[userNodeInc] = participant
		userNodeInc++
	}

	// Link timeslots to sink
	for i := intialTimeslotNodeInc; i < timeslotNodeInc; i++ {
		graph = append(graph, Edge{From: i, To: timeslotNodeInc, Capacity: 1, Cost: 0})
	}
	return graph, timeslotNodeInc, nodeToUser
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
	random_float := (random.Float64() * 0.00005) + 0

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
