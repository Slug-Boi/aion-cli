# aion-cli
<p align="center">
<img src="src/html/css/aion_logo.png" alt="aion_logo" width="200"/>
</p>
<a href="https://github.com/Slug-Boi/aion-cli/releases/latest"><img src="https://img.shields.io/badge/dynamic/yaml?url=https%3A%2F%2Fraw.githubusercontent.com%2FSlug-Boi%2Faion-cli%2Fmaster%2F.github%2Fbadges%2Frelease_badge.yml&query=%24.version&logo=github&label=Release
" alt="github latest release"></a>
<a href="https://pkg.go.dev/github.com/Slug-Boi/aion-cli"><img src="https://img.shields.io/badge/_-reference-blue?logo=go&label=%E2%80%8E%20
" alt="golang package reference"></a>

This is a CLI tool that takes in requests/timeslot wishes from a google form and schedules them in a way that minimizes the number of conflicts.  

The tool uses one of two different solvers, the [Min Cost flow algorithm (path augmentation based on SPFA)](https://cp-algorithms.com/graph/min_cost_flow.html) or [Gurobi](https://www.gurobi.com/resources/mixed-integer-programming-mip-a-primer-on-the-basics/) depending the subcommand used. The tool will output the schedule in a html file that is hosted locally on your machine, where it can also be downloaded as a CSV file that can be imported into Excell or Google Sheets. The program can also create an ICal file (.ics) which can be imported into most calendars.

The intention is to make the html file display in a transparent manner the way the schedule was generated, so that the user can understand the reasoning behind the schedule.

## Features
The tool currently features a couple of commands listed below are the most important if you want a full list you can use the -h flag

### Generate 
The generate command will run the solver on a given form and will display the output using an HTML file that will be hosted locally on your machine. The page will also allow the user to save the data as a CSV or the output generated as an ICal ics file

### Config
The tool has a config command that will allow the user to edit the configuration file that the program uses this is where the default solver is stored as well as a default form_ID to use when fetching data and more. 
The program will ask to create a config file the first time it runs. To see an example of a config file you can see the [examples folder](https://github.com/Slug-Boi/aion-cli/blob/master/examples/example_config.json)

### Solve
The solve command is used mostly as a debugging command. It functions almost the same as generate but instead of showing the result as an html file it outputs the result in the terminal. This is recommended for testing or debugging purposes

### Form 
This command will curl GET the google sheets ID and display the form data in the terminal. This is mostly for testing that your configuration is setup correctly.


## Usage 
[...]: <u>Optional arguments</u>  
<...>: <u>Required Arguments</u>

### Generate
The generate command runs the default solver on a google form by supplying a google sheets ID which is linked to the form and will output the results in a locally hosted HTML file. The ID can either be entered as an argument or saved to the config file to be used as the default  
Example:
```
$ aion-cli generate [form ID]
```
The generate command has 3 flags:  
- --cal - will save the output as an ICal ics file
- --minCost - overrides the default solver to use min_cost flow solver
- --gurobi - overrides the default solver to use gurobi solver

### Solve
The solve command lets the user pull data from a google form by supplying the ID to a linked google sheets as an argument or saving that ID to the config file. The solve command will use the solver selected via a subcommand (or the default solver if no subcommand is used) to find an optimal timeslot allocation based on wishes and output the solution in the terminal. For details about the different solvers please see the [docs folder](https://github.com/Slug-Boi/aion-cli/blob/master/docs/useful_links.md)

Example using the min cost flow algorithm:  
```
$ aion-cli solve min_cost [form_ID]
```  

Example using the Gurobi solver:  
```
$ aion-cli solve gurobi [form_ID]
```  

If the form_ID is saved to the config file using the optional --save flag then you can omit the form_ID argument for all future calls to the same form

### Config
The config command allows the user to edit or remove the config file the program uses. These are the subcommands:  
- solver - select the default solver the program uses
- formID - edit the default formID used when no argument is provided to certain commands
- ical - edit whether the program auto saves the solution as an ICal ics file when using the generate command 

Examples:
```
$ aion-cli config ical true
$ aion-cli config solver min_cost
$ aion-cli config formID abcdefg
```


## Installation
### requirements:
Latest aion-cli CLI binary  
Python 3 or later (used for gurobi solver)  
  

### Manual installation
To install the program you just need to run the binary file. For ease of use you can create an alias in your shell rc file or powershell profile

### Nix support
Coming soon  
The tool can easily be installed using the nix package manager.  

## Future Work
The hope for this program for future development (if we pick it back up) is to add the ability of creating groups with wishes as well and combining these two features so they can be used in tandem. This is a much harder problem to solve and thus is out of scope for the initial release of this project.



