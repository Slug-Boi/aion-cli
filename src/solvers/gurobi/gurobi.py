#!/usr/bin/env python3.11

# Copyright 2024, Gurobi Optimization, LLC

# This code is very heavily based on the gurobi workforce1 python example:
# https://docs.gurobi.com/projects/examples/en/current/examples/workforce.html#subsectionworkforce1examples


# Objective: Assign groups to timeslots;
# Each group must be assigned to exactly one timeslot
# Each timeslot must at most have one group assigned to it or none
# The objective is to minimize the cost of assigning groups to timeslots
# based on the wishes of the groups


# note from gurobi example: if the problem cannot be solved, use IIS to
# find a set of conflicting constraints. Note that there may
# be additional conflicts besides what is reported via IIS.

# debug test command:
# $ python gurobi.py "gr1,gr2" "Tue2;gr1;4,Tue2;gr2;1,Wed1;gr2;1,Wed2;gr1;1"

#TODO: Do this again from the bottom up without this example and use the GRB binary variable instead

import gurobipy as gp
from gurobipy import GRB
import sys

args = sys.argv

if len(args) != 3:
    print("Usage: gurobi.py \"<group1>,<group2>..."" \"<timeslot1:group1:cost1>,<timeslot2:group2:cost2>...\"")
    sys.exit(1)

groups = args[1].split(",")
timeslots = args[2].split(",")

groupdict = {}
for group in groups:
    groupdict[group] = 1

dayGroupTupleList = []
timeslotCostDict = {}

for timeslot in timeslots:
    day, group, cost = timeslot.split(";")
    dayGroupTupleList.append((day, group))
    # TODO: might need to change to a float but wait
    timeslotCostDict[(day, group)] = float(cost)

# TODO: Change variable names to fit the problem better
# Number of groups needed for each timeslot (This will always be 1 as we need
# a 1-1 relation between groups
# and timeslots)
groups, _ = gp.multidict(groupdict)
# append groups to shifts and shiftRequirements


# The cost of assigning a group to a timeslot this is calculated from
# the wishes given by the groups
# in the google form
# Green = 0, Yellow = 10/sum, red = 100/sum
# All values have a heuristic tie breaker value added to them
cost = gp.tupledict(timeslotCostDict)

# Mapping from timeslot to group this is what we need to minimize the cost of
availability = gp.tuplelist(dayGroupTupleList)

# Model
m = gp.Model("assignment")

# Assignment variables: x[w,s] == 1 if worker w is assigned to shift s.
# Since an assignment model always produces integer solutions, we use
# continuous variables and solve as an LP.
# vtype=GRB.BINARY can be used to specify binary variables.
x = m.addVars(availability, ub=1, name="x")

# The objective is to minimize the total pay costs
m.setObjective(gp.quicksum(cost[w, s] * x[w, s]
               for w, s in availability), GRB.MINIMIZE)


# Constraints: assign exactly shiftRequirements[s] workers to each shift s
reqCts = m.addConstrs(
    (x.sum("*", s) == 1 for s in groups))

# Constraints: each group is assigned at most one timeslot
m.addConstrs(x.sum(s, "*") <= 1 for s, _ in availability)


# Save model
m.write("aion.lp")


# TODO: Figure out how to get the relevant data out of the solver afterwards
# variables we want:
# Path from group to timeslot (e.g gr1 -> Tue2)
# Cost of the assignment (e.g. gr1 -> Tue2 = 4)
# The total cost of the Assignment (e.g. 4 + 1 + 1 + 1)
# Solve time (e.g. 0.01s)

# Optimize
m.optimize()
status = m.Status
if status == GRB.UNBOUNDED:
    print("The model cannot be solved because it is unbounded")
    sys.exit(0)
if status == GRB.OPTIMAL:
    print(f"The optimal objective is {m.ObjVal:g}")
    if x is not None:
        for v in x:
            if x[v].X != 0 and x[v].X != 1:
                print("Warning: solution is not binary")
                print(f"{v} = {x[v].X}")
                exit(1)
            if x[v].X > 0.5:
                print(f"{v[0]}->{v[1]}")
    sys.exit(0)
if status != GRB.INF_OR_UNBD and status != GRB.INFEASIBLE:
    print(f"Optimization was stopped with status {status}")
    sys.exit(0)

# do IIS
print("The model is infeasible; computing IIS")
m.computeIIS()
if m.IISMinimal:
    print("IIS is minimal\n")
else:
    print("IIS is not minimal\n")
print("\nThe following constraint(s) cannot be satisfied:")
for c in m.getConstrs():
    if c.IISConstr:
        print(c.ConstrName)
