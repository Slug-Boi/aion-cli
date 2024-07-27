package gurobi

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/graph"
)

// Translates the form data to Gurobi syntax (this is proprietary to how the gurobi python program works)
// a new translator will more than likely need to be created for any other optimization program
func TranslateGurobi(data []forms.Form) (string, string) {
	// Sort users by id to ensure consistent ordering when generating the concatenated string
	sort.Slice(data, func(i, j int) bool {
		return data[i].HashString < data[j].HashString
	})

	// Create base string for creating heuristic
	sb := strings.Builder{}

	// Get the string of all group inputs and the cache
	allStrings := graph.BaseHashString(data, sb)

	// Create string builders for two return strings
	sbGroups := strings.Builder{}
	sbTimeslots := strings.Builder{}

	for _, participant := range data {
		// Add group to string builder
		sbGroups.WriteString(participant.GroupNumber + ",")

		// Create heuristic for the participant
		heuristic := graph.HashHeuristic(participant.HashString, allStrings)

		// Translate timeslots to participant linked nodes:
		// Caps are all the individual costs for each timeslot
		// SumCap is the sum of all the costs for each timeslot
		caps :=  map[string]float64{}
		var sumCap float64

		// Calculation is done this way to make sure that Want are all weighted equally
		// Can do and Cannot are weighted different between groups with the idea being that groups
		// with many wishes will get lower values overall on their wishes meaning they are more likely to get
		// their wishes granted as they are more flexible. Can do will always be weighted lower than cannot
		// due to the sum division. These calculations give a more fair distribution of timeslots between groups
		// and will incentivize groups to be more flexible with their timeslots and answer truthfully
		for timeslot, vote := range participant.Votes {
			caps, sumCap = graph.CostSummer(timeslot,vote, caps, sumCap)
		}

		// add the timeslot costs to the string builder
		for timeslot, cap := range caps {
			sbTimeslots.WriteString(timeslot+";"+participant.GroupNumber+";"+fmt.Sprintf("%f",(cap/sumCap)+heuristic)+",")
		}
	}

	return sbGroups.String()[:sbGroups.Len()-1], sbTimeslots.String()[:sbTimeslots.Len()-1]
}

// Currently this runs the Gurobi optimization through python. There is a gurobi library for Go
// but it is different syntax wise so this is a temporary solution until the Go library is implemented
// (if time permits)
func RunGurobi(data []forms.Form) (string, error) {
	
	// translate 
	groups, timeslots := TranslateGurobi(data)

	cmd := exec.Command("python", "./gurobi/gurobi.py", groups, timeslots)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
