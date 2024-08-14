it=5
iit=3
formID=""
shuffle=false
worst=false

# currently the worst flag doesn't seem to work as intended please refrain from using it for now

# false=gurobi, true=minCost
solver=true

while getopts "i:l:f:s:w:g" opt; do
    case $opt in
        i) it=$OPTARG;;
        l) iit=$OPTARG;;
        f) formID=$OPTARG;;
        s) shuffle=true;;
        w) worst=true;;
        g) solver=false;;
        \?) echo "Invalid option -$OPTARG" >&2;;
    esac
done
echo "Running $it iterations"

rm min_cost.log gurobi.log min_err.log gur_err.log

# Generate the logs
if [ "$solver" = true ]; then
    echo "Using minCost"
    for i in $(seq 1 $it); do
        python random_data.py 

    ./aion-cli solve minCost ${formID} > min_cost.log 2> min_err.log
    MINCOST=$(grep -F -- "User:" min_cost.log | cut -d ' ' -f 2-10) 
    COSTMINCOST=$(grep -F -- "Min Cost:" min_cost.log | cut -d ' ' -f 4) 
    for i in $(seq 1 $iit); do
    if [ "$shuffle" == true ]; then
        python random_data.py shuffle
    fi
      ./aion-cli solve minCost ${formID} > min_cost.log 2> min_err.log

      MINCOST2=$(grep -F -- "User:" min_cost.log | cut -d ' ' -f 2-10) 

      COSTMINCOST2=$(grep -F -- "Min Cost:" min_cost.log | cut -d ' ' -f 4)

      if (( $(echo "$COSTMINCOST > $COSTMINCOST2" |bc -l) )) || (( $(echo "$COSTMINCOST < $COSTMINCOST2" |bc -l) )); then
          echo "Mismatch in iteration $i"
          echo "MinCost1: $COSTMINCOST"
          echo "MinCost2: $COSTMINCOST2"
      fi

      if [ "$MINCOST" != "$MINCOST2" ]; then
          echo "Mismatch in iteration $i"
          echo "MinCost1: $MINCOST"
          echo "MinCost2: $MINCOST2"
          exit 1
      fi
    done  
    done    
fi
if [ "$solver" = false ]; then
    echo "Using gurobi"
    for i in $(seq 1 $it); do
    python random_data.py

    ./aion-cli solve gurobi ${formID} > gurobi.log 2> gurobi.log

    GUROBI=$(grep -F -- "group" gurobi.log | cut -d ' ' -f 1-10)

    COSTGUROBI=$(grep -F -- "Min Cost:" gurobi.log | cut -d ' ' -f 3) 

    for i in $(seq 1 $iit); do
    if [ "$shuffle" == true ]; then
        python random_data.py shuffle
    fi
      ./aion-cli solve gurobi ${formID} > gurobi.log 2> gur_err.log

      GUROBI2=$(grep -F -- "group" gurobi.log | cut -d ' ' -f 1-10)

      COSTGUROBI2=$(grep -F -- "Min Cost:" gurobi.log | cut -d ' ' -f 3)

      if [ "$GUROBI" != "$GUROBI2" ]; then
          echo "Mismatch in iteration $i"
          echo "Gurobi: $GUROBI"
          echo "Gurobi: $GUROBI2"
          exit 1
      fi
    done
    done
fi

rm min_cost.log gurobi.log min_err.log gur_err.log form.csv

echo "All iterations passed"