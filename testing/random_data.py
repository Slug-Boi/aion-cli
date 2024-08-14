import datetime
import sys
import random
from time import gmtime, strftime

# Add worst case where all groups wish for everything and randomize the input data

# Group class to store group data
class Group:
    def __init__(self, group_number, hashstring, timestamp, votes):
        self.group_number = group_number
        self.hashstring = hashstring
        self.timestamp = timestamp
        self.votes = votes

def shuffle(file):
    with open(file, "r") as f:
        lines = f.readlines()
    
    # Remove first line
    header = lines.pop(0)
    with open(file, "w") as f:
        f.write(header)
        random.shuffle(lines)
        f.writelines(lines)
    
    exit(0)


# Function to return a vote based on random number
def vote(num):
    match num:
        case 0:
            return "Want"
        case 1:
            return "Can do"
        case 2:
            return ""

args = sys.argv[1:]

totalGroups = 0
totalTimeslots = 0

# Set the total number of groups and timeslots based on args
if len(args) == 1 and args[0] == "shuffle":
    shuffle("form.csv")
else:
    if len(args) == 1 and 1 < int(args[0]) <= 20:
        totalGroups = int(args[0])
    if len(args) == 2 and 1 < int(args[0]) <= 20 and int(args[0]) <= int(args[1]) <= 40:
        totalGroups = int(args[0])
        totalTimeslots = int(args[1])

# if no args provided generate random number for each
# Groups bounded to 20 and timeslots bounded to 40
if len(args) == 0:
    totalGroups = random.randint(2, 20)
    totalTimeslots = random.randint(totalGroups, 40)

timeslots = []
for i in range(totalTimeslots):
    timeslots.append("timeslot "+str(i+1))

timestamp = datetime.datetime.now()

groups = []
# Generate random data for groups and timeslots
for i in range(totalGroups):
    group_number = "group "+str(i + 1)
    le = random.randint(1, 32)
    hashstring = random.randbytes(le).hex()
    votes = []
    for j in range(totalTimeslots):
        ran = random.randint(0, 2)
        v = vote(ran)
        votes.append(v)
    shiftedTime = timestamp + datetime.timedelta(seconds=20*i)
    strTime = shiftedTime.strftime("%d/%m/%Y %H:%M:%S")
    group = Group(group_number, hashstring, strTime, votes)
    groups.append(group)
    
# Write the data to a csv file named random_data.csv
lines = ["Timestamp,Group Number,Lottery String,"]
for i in range(totalTimeslots):
    if i != totalTimeslots-1:
        lines.append(timeslots[i]+",")
    else:    
        lines.append(timeslots[i]+"\n")

for i in range(totalGroups):
    lines.append(groups[i].timestamp+","+groups[i].group_number+","+groups[i].hashstring+",")
    for j in range(totalTimeslots):
        if j != totalTimeslots-1:
            lines.append(groups[i].votes[j]+",")
        else:
            lines.append(groups[i].votes[j])
    lines.append("\n")

with open("form.csv", "w") as f:
    f.writelines(lines)