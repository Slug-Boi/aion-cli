package gurobi

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Slug-Boi/aion-cli/forms"
	libfuncs "github.com/Slug-Boi/aion-cli/lib_funcs"
	"github.com/Slug-Boi/aion-cli/logger"
	"github.com/Slug-Boi/aion-cli/solvers/graph"
)

var Sugar = logger.SetupLogger()

//go:embed gurobi.py
var gurobiFiles embed.FS

// Translates the form data to Gurobi syntax (this is proprietary to how the gurobi python program works)
// a new translator will more than likely need to be created for any other optimization program
func TranslateGurobi(data []forms.Form) (string, string, map[string]forms.Form) {
	// Create number to form map
	users := map[string]forms.Form{}

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
		// add user to map
		users[participant.GroupNumber] = participant

		// Add group to string builder
		sbGroups.WriteString(participant.GroupNumber + ",")

		// Create heuristic for the participant
		heuristic := graph.HashHeuristic(participant.HashString, allStrings)

		// Translate timeslots to participant linked nodes:
		// Caps are all the individual costs for each timeslot
		// SumCap is the sum of all the costs for each timeslot
		caps := map[string]float64{}
		sumCap := 0.0

		// Calculation is done this way to make sure that Want are all weighted equally
		// Can do and Cannot are weighted different between groups with the idea being that groups
		// with many wishes will get lower values overall on their wishes meaning they are more likely to get
		// their wishes granted as they are more flexible. Can do will always be weighted lower than cannot
		// due to the sum division. These calculations give a more fair distribution of timeslots between groups
		// and will incentivize groups to be more flexible with their timeslots and answer truthfully
		for timeslot, vote := range participant.Votes {
			caps, sumCap = graph.CostSummer(timeslot, vote, caps, sumCap)
		}

		if sumCap == 0 {
			sumCap = 1
		}
		// add the timeslot costs to the string builder
		for timeslot := range participant.Votes {
			sbTimeslots.WriteString(timeslot + ";" + participant.GroupNumber + ";" + fmt.Sprintf("%v", (caps[timeslot]/sumCap)+heuristic) + ",")
		}
	}

	return sbGroups.String()[:sbGroups.Len()-1], sbTimeslots.String()[:sbTimeslots.Len()-1], users
}

// Currently this runs the Gurobi optimization through python. There is a gurobi library for Go
// but it is different syntax wise so this is a temporary solution until the Go library is implemented
// (if time permits)
func RunGurobi(data []forms.Form) (string, map[string]forms.Form, error) {

	// translate
	groups, timeslots, users := TranslateGurobi(data)

	gurobiCode, err := gurobiFiles.ReadFile("gurobi.py")
	if err != nil {
		return "", users, err
	}

	// run the gurobi python script
	cmd := exec.Command("python", "-c", string(gurobiCode), groups, timeslots)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), users, err
	}

	// remove most of the output as it is not needed
	// TODO: figure out if extended output should be enabled through a flag

	splitOut := strings.Split(string(out), "Optimal objective ")
	if len(splitOut) < 2 {
		return "",users, fmt.Errorf("No optimal solution found")
	}

	out = []byte(splitOut[1])

	return string(out), users, nil
}

func cleanupLP() {
	os.Remove("aion.lp")
}

func SolveGurobi(args []string) (string,map[string]string, map[string]string) {

	// Get the config file
	conf := libfuncs.SetupConfig(args)

	Sugar.Debugln("Form is being processed with the following Form ID:", conf.FormID)

	// Get the form data
	data := forms.GetForm(conf)

	// Run the gurobi python program
	out, users, err := RunGurobi(data)
	if err != nil {
		Sugar.Panicf("Error running gurobi: %v\n%s", err,out)
	}

	// Parse the output
	groupTimeslot := map[string]string{}
	wishLevels := map[string]string{}

	costTimeslot := strings.Split(out[1:], "\n")
	for _, timeslot := range costTimeslot[2:len(costTimeslot)-1] {
		split := strings.Split(timeslot, "->")
		groupTimeslot[split[1]] = split[0]
		wishLevels[split[1]] = users[split[1]].Votes[split[0]]
	}

	defer cleanupLP()

	//TODO: add return values for the gurobi solver
	return costTimeslot[0], groupTimeslot, wishLevels
}
