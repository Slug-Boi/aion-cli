package html

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"embed"

	"github.com/Slug-Boi/aion-cli/solvers/graph"
)

//go:embed root.html
var htmlFiles embed.FS

//go:embed css/*
var cssFiles embed.FS

type WebData struct {
	GroupNumber string
	Timeslot    string
	Day         string
	Date        string
	WishLevel   string
	//Path        []int
}

// HTML template code inspired by https://gowebexamples.com/templates/
func GenerateHTML(args []string) {
	iterations := 0

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webData []WebData
		iterations++
		// Run the solver
		sink, users, _, paths, nodeToTimeslot := graph.SolveMin_Cost(args)

		finalPaths := map[int]int{}

		// Create a map that maps the group number to the timeslot
		// This is required since some paths will have residual graphs in them which will overwrite
		// Previous paths chosen
		for _, path := range paths {
			i := sink
			timeslotNode := -1
			for i != 0 {
				if _, ok := users[i]; !ok {
					timeslotNode = i
				} else {
					if timeslotNode != -1 {
						finalPaths[i] = timeslotNode
					}
					timeslotNode = -1
				}
				i = path[i]
			}
		}

		// TODO: Convert path slice to an actual linear path to reduce nested calls as it looks disgusting
		// Might be as simple as using the uncommented code below try later
		for user, timeslot := range finalPaths {
			timeslotStr := nodeToTimeslot[timeslot]
			Date_Day_Timeslot := strings.Split(timeslotStr, " ")
			trimmedTimeslot := strings.Trim(Date_Day_Timeslot[2], "[]")
			day := Date_Day_Timeslot[1]
			date := Date_Day_Timeslot[0]
			wishLevel := users[user].Votes[timeslotStr]

			webData = append(webData, WebData{GroupNumber: users[user].GroupNumber, Timeslot: trimmedTimeslot, Day: day, Date: date, WishLevel: wishLevel})
		}

		// sort by group number
		sort.Slice(webData, func(i, j int) bool {
			return webData[i].GroupNumber < webData[j].GroupNumber
		})
		
		t, err := template.ParseFS(htmlFiles, "root.html")
		if err != nil {
			//TODO: change to Zap logger later
			log.Fatal(err)
		}

		data := struct {
			WebData    []WebData
			Iterations int
		}{
			WebData:    webData,
			Iterations: iterations,
		}

		err = t.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	// Change the working directory to the html folder
	os.Chdir("html")

	// Serve the css file for the html
	var css http.FileSystem = http.FS(cssFiles)
	http.Handle("/css/", http.FileServer(css))
	//fs := http.FileServer(http.Dir("./css/"))
	//http.Handle("/css/", http.StripPrefix("/css", fs))

	// Launch the server on port 80
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		// TODO: Use zaplogger to log the error.
		return
	}

}
