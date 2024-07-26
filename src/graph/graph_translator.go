package graph

import (
	"hash/fnv"
	"math/rand"
	"sort"
	"strings"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/thanhpk/randstr"
)

// Translates data from the forms package to the graph package
func Translate(data []forms.Form) ([]Edge, int, map[int]forms.Form, map[int]string) {
	// Create a map to store the user node to user data
	nodeToUser := map[int]forms.Form{}

	// Timeslot node to Unix Time map (First value is starting time and second value is ending time)
	nodeToTime := map[int]string{}
	// Users are always node 1 to len(data.PollResults)
	userNodeInc := 1
	// Timeslots are always node len(data.PollResults) + 1 to len(data.PollResults) + len(participants.Votes)
	intialTimeslotNodeInc := len(data) + 1
	// Start at the initial timeslot node for the incrementor
	timeslotNodeInc := intialTimeslotNodeInc

	// Cache is used to store the group name going to the ID for consistent hashing later
	cache := map[string]string{}

	// Sort users by id to ensure consistent ordering when generating the concatenated string
	sort.Slice(data, func(i, j int) bool {
		return data[i].HashString < data[j].HashString
	})

	// Create base string for creating heuristic
	sb := strings.Builder{}

	// Get the string of all group inputs
	allStrings, cache := BaseHashString(data, cache, sb)


	graph := []Edge{}

	// Translate participants to source linked nodes
	for i, participant := range data {

		heuristic := HashHeuristic(cache[participant.GroupNumber], allStrings)

		// Add edge from source to participant
		graph = append(graph, Edge{From: 0, To: userNodeInc, Capacity: 1, Cost: 0})

		// Translate timeslots to participant linked nodes:
		// Caps are all the individual costs for each timeslot
		// SumCap is the sum of all the costs for each timeslot
		var caps map[string]float64
		var sumCap float64

		// Calculation is done this way to make sure that Want are all weighted equally
		// Can do and Cannot are weighted different between groups with the idea being that groups
		// with many wishes will get lower values overall on their wishes meaning they are more likely to get
		// their wishes granted as they are more flexible. Can do will always be weighted lower than cannot
		// due to the sum division. These calculations give a more fair distribution of timeslots between groups
		// and will incentivize groups to be more flexible with their timeslots and answer truthfully
		for timeslot, vote := range participant.Votes {
			// Add edge from participant to timeslot
			if _, ok := nodeToTime[timeslotNodeInc]; !ok {
				nodeToTime[timeslotNodeInc] = timeslot
				timeslotNodeInc++
			}
			caps, sumCap = CostSummer(timeslot,vote, caps, sumCap)
		}
		timeslotNodeInc = intialTimeslotNodeInc
		// Add edge from participant to timeslot
		//TODO: Check that this still works now that caps is map and not a float slice
		for _, cap := range caps {
			graph = append(graph, Edge{From: userNodeInc, To: timeslotNodeInc, Capacity: 1, Cost: cap/sumCap + heuristic})
			timeslotNodeInc++
		}

		if i != len(data)-1 {
			timeslotNodeInc = intialTimeslotNodeInc
		}

		// Add user to map and increment user node
		nodeToUser[userNodeInc] = participant
		userNodeInc++
	}

	// Link timeslots to sink
	for i := intialTimeslotNodeInc; i < timeslotNodeInc; i++ {
		graph = append(graph, Edge{From: i, To: timeslotNodeInc, Capacity: 1, Cost: 0})
	}
	return graph, timeslotNodeInc, nodeToUser, nodeToTime
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
	random_float := (random.Float64() * 0.00000005) + 0

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

func BaseHashString(data []forms.Form, cache map[string]string, sb strings.Builder) (string, map[string]string) {

	for i := 0; i < len(data); i++ {
		//TODO: Probably redo this to be more modular
		if data[i].HashString != "" {
			//TODO: Add a regex check for valid characters
			cache[data[i].GroupNumber] = data[i].HashString
			sb.WriteString(data[i].HashString)
		} else {
			// If no string is provided then use their id as it is a pseudo random string
			// that will result in the same value each time
			// Random string package: https://github.com/thanhpk/randstr
			token := randstr.String(16) // generate a random 16 character length string
			cache[data[i].GroupNumber] = token
			sb.WriteString(token)
		}
	}

	return sb.String(), cache
}

func CostSummer(timeslot,vote string, caps map[string]float64, sumCap float64) (map[string]float64, float64) {
			// Translate wishes to float values
			if vote == "Want" {
				caps[timeslot] = 0.0
				// Implicit cost of 0 added to sum
			} else if vote == "Can do" {
				caps[timeslot] = 10.0
				sumCap += 10.0
			} else {
				caps[timeslot] = 100.0
				sumCap += 100.0
			}

			return caps, sumCap
}
