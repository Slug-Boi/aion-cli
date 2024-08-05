package html

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Slug-Boi/aion-cli/logger"
	"github.com/Slug-Boi/aion-cli/solvers/graph"
	"github.com/Slug-Boi/aion-cli/solvers/gurobi"
	"github.com/emersion/go-ical"
	"github.com/facette/natsort"
)

var Sugar = logger.SetupLogger()

var calChannel = make(chan []WebData)

//go:embed *.html
var htmlFiles embed.FS

//go:embed css/*
var cssFiles embed.FS

type WebData struct {
	GroupNumber string
	Timeslot    string
	Day         string
	Date        string
	WishLevel   string
	Path        string
	PathCost    float64
}

func graphWebData(args []string, sorterMethod string) ([]WebData, string) {
	var webData []WebData

	// Run the solver
	sink, users, cost, paths, nodeToTimeslot, groupTimeslotCost := graph.SolveMin_Cost(args)

	translatedPaths := map[int]int{}

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
					translatedPaths[i] = timeslotNode
				}
				timeslotNode = -1
			}
			i = path[i]
		}
	}

	for user, timeslot := range translatedPaths {
		timeslotStr := nodeToTimeslot[timeslot]
		Date_Day_Timeslot := strings.Split(timeslotStr, " ")
		trimmedTimeslot := strings.Trim(Date_Day_Timeslot[2], "[]")
		day := Date_Day_Timeslot[1]
		date := Date_Day_Timeslot[0]
		wishLevel := users[user].Votes[timeslotStr]

		webData = append(webData, WebData{GroupNumber: users[user].GroupNumber, Timeslot: trimmedTimeslot, Day: day, Date: date, WishLevel: wishLevel, Path: fmt.Sprintf("[0,%d,%d,%d]", user, timeslot, sink), PathCost: groupTimeslotCost[users[user].GroupNumber+timeslotStr]})
	}

	// sort by group number using natural sorting
	webData = sorter(webData, sorterMethod)

	return webData, strconv.FormatFloat(cost, 'f', -1, 64)
}

func gurobiWebData(args []string, sorterMethod string) ([]WebData, string) {
	var webData []WebData

	// Run the solver
	cost, timeslots, wishLevels, groupTimeslotCost := gurobi.SolveGurobi(args)

	for group, timeslot := range timeslots {
		timeslotStr := strings.Split(timeslot, " ")
		timeslotTrimmed := strings.Trim(timeslotStr[2], "[]")
		webData = append(webData, WebData{GroupNumber: group, Timeslot: timeslotTrimmed, Day: timeslotStr[1], Date: timeslotStr[0], WishLevel: wishLevels[group], Path: "Gurobi does not support paths", PathCost: groupTimeslotCost[group+timeslot]})
	}

	// sort by group number using natural sorting
	webData = sorter(webData, sorterMethod)

	return webData, strconv.FormatFloat(cost, 'f', -1, 64)

}

func sorter(webData []WebData, sortMethod string) []WebData {
	if sortMethod == "timeslot" {
		// time.Parse layout string
		layouts := []string{"02-01-06 Monday 15:04", "02-01-2006 Monday 15:04", "02-01-06 Mon 15:04", "02-01-2006 Mon 15:04"}
		layout := ""
		// figure out which layout to use for the time.Parse
		for i := 0; i < len(layouts); i++ {
			timetest, _ := time.Parse(layouts[i], webData[0].Date+" "+webData[0].Day+" "+strings.Split(webData[0].Timeslot, "-")[0])
			if timetest.String() != "0001-01-01 00:00:00 +0000 UTC" {
				layout = layouts[i]
				break
			}
		}
		// if no layout was found, panic
		if layout == "" {
			Sugar.Panic("Could not find a suitable layout for time.Parse. Data timeslot format is not supported.")
		}
		// sort by timeslot using time.Parse and time.Before
		sort.Slice(webData, func(i, j int) bool {
			t1, _ := time.Parse(layout, webData[i].Date+" "+webData[i].Day+" "+strings.Split(webData[i].Timeslot, "-")[0])
			t2, _ := time.Parse(layout, webData[j].Date+" "+webData[j].Day+" "+strings.Split(webData[j].Timeslot, "-")[0])
			return t1.Before(t2)
		})
	} else if sortMethod == "group_number" {
		// sort by group number using natural sorting
		// e.g. group 1, group 2, group 3 ...., group 10, group 11, ....
		sort.Slice(webData, func(i, j int) bool {
			return natsort.Compare(webData[i].GroupNumber, webData[j].GroupNumber)
		})
	}
	return webData
}

// HTML template code inspired by https://gowebexamples.com/templates/
func GenerateHTML(args []string, solver, sorterMethod string) {
	iterations := 0
	fmt.Println("Starting server")
	go func() {
		for {
			res, _ := http.Get("http://localhost:80")
			if res.StatusCode == 200 {
				fmt.Println("Server is up: http://localhost:80")
				break
			}
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webData []WebData

		iterations++

		if solver == "gurobi" {
			webData, _ = gurobiWebData(args, sorterMethod)
		} else {
			webData, _ = graphWebData(args, sorterMethod)
		}

		// Async call to create the ical file (only works if the flag is set)
		go func() {
			calChannel <- webData
			//defer close(calChannel)
		}()

		t, err := template.ParseFS(htmlFiles, "root.html")
		if err != nil {
			Sugar.Panic(err)
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

	http.HandleFunc("/advanced", func(w http.ResponseWriter, r *http.Request) {
		var webData []WebData
		var cost string

		iterations++

		if solver == "gurobi" {
			webData, cost = gurobiWebData(args, sorterMethod)
		} else {
			webData, cost = graphWebData(args, sorterMethod)
		}

		// Async call to create the ical file (only works if the flag is set)
		go func() {
			calChannel <- webData
			//defer close(calChannel)
		}()

		t, err := template.ParseFS(htmlFiles, "root_advanced.html")
		if err != nil {
			Sugar.Panic(err)
		}

		data := struct {
			WebData    []WebData
			Iterations int
			Cost       string
		}{
			WebData:    webData,
			Iterations: iterations,
			Cost:       cost,
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
		Sugar.Panic(err)
	}
}

// this function will create an ical calendar
func CreateICal() {
	// Create a new calendar
	data := <-calChannel

	cal := ical.NewCalendar()
	cal.Props.SetText(ical.PropVersion, "2.0")
	//TODO: Figure out what this is about
	cal.Props.SetText(ical.PropProductID, "-//Aion CLI Scheduler version 1.0//EN")

	// Add events to the calendar with the data from the solver
	for _, group := range data {
		event := ical.NewEvent()
		event.Props.SetText(ical.PropUID, group.GroupNumber)
		event.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
		//TODO: Add the type of meeting here potentially
		event.Props.SetText(ical.PropSummary, group.GroupNumber)

		// Parse the date and time from the string and add it to the event as the start time
		layout := "02-01-06 15:04"
		startTime := strings.Split(group.Timeslot, "-")[0]
		endTime := strings.Split(group.Timeslot, "-")[1]

		calStartTime, _ := time.Parse(layout, group.Date+" "+startTime)
		calEndTime, _ := time.Parse(layout, group.Date+" "+endTime)

		event.Props.SetDateTime(ical.PropDateTimeStart, calStartTime)
		event.Props.SetDateTime(ical.PropDateTimeEnd, calEndTime)

		cal.Children = append(cal.Children, event.Component)
	}

	// calendar buffer
	var buf bytes.Buffer

	// Encode the calendar to the buffer
	if err := ical.NewEncoder(&buf).Encode(cal); err != nil {
		Sugar.Panicf("Error marshaling the calendar:\n", err.Error())
	}

	//Write the calendar to a file
	err := os.WriteFile("calendar.ics", buf.Bytes(), 0644)
	if err != nil {
		Sugar.Panicf("Error writing the calendar to a file:\n", err.Error())
	}

}
