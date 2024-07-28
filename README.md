# aion-cli
This is a CLI tool that takes in requests/timeslot wishes from a strawpoll form and schedules them in a way that minimizes the number of conflicts.  

The tool uses one of two different solvers, the [min cost flow algorithm (path augmentation based on SPFA)](https://cp-algorithms.com/graph/min_cost_flow.html) or Gurobi depending the subcommand used. The tool will output the schedule in a html file that is hosted locally on your machine, where it can also be downloaded as a CSV file that can be imported into Excell or Google Sheets.  

The intention is to make the html file display in a transparent manner the way the schedule was generated, so that the user can understand the reasoning behind the schedule.

## Features
The tool currently features a couple of commands listed below are the most important if you want a full list you can use the -h flag or check out the docs folder on the repo

### Generate 
The generate command will run the solver on a given form and will display the output using an HTML file that will be hosted locally on your machine. The page will also allow the user to save the data as a CSV file that can be imported in any CSV compatible application.

### Solve
The solve command is used mostly as a debugging command. It functions almost the same as generate but instead of showing the result as an html file it outputs the result in the terminal. This result breakdown is much more simplistic and as such is not recommended for anything other than testing or debugging purposes

### Config
The tool has a config configuration command set that will allow the user to edit the config that the program uses this is where the default solver is stored as well as a default form_ID to use when fetching data and more. 
The program automatically creates a config file the first time it runs. To see an example of a config file you can see the examples folder 
### Form 
This command is extremely simple. It will curl GET the google sheets ID and display the form data in the terminal. This is mostly for testing that your configuration is setup correctly.


## Usage 
### Solve
The solve command lets the user pull data from a google form by supplying the ID to a linked google sheets as an argument or saving that ID to the config file. The solve command will use the solver selected via a subcommand (or the default solver if no subcommand is used) to find an optimal timeslot allocation based on wishes and output the solution in the terminal. For details about the different output types see ...  

Example using the min cost flow algorithm:  
```
aion solve min_cost [form_ID]
```  

Example using the Gurobi solver:  
```
aion solve gurobi [form_ID]
```  

If the form_ID is saved to the config file using the optional --save flag then you can omit the form_ID argument for all future calls to the same form



## Installation
### requiremnts:
latest Aion CLI binary 
Python 3 or later (used for gurobi solver)
curl (used to fetch google form data)  

### manual install
to install the program you just need to run the binary file. For ease of use you can create an alias in your shell rc file or powershell profile


Coming soon
### Nix support
Coming soon  
The tool can easily be installed using the nix package manager.  

## Future Work
The hope for this program for future development (if we pick it back up) is to add the ability of creating groups with wishes as well and combining these two features so they can be used in tandem. This is a much harder problem to solve and thus is out of scope for the initial release of this project.



