it=5
iit=5
formID=""

# false=gurobi, true=minCost
solver=true

while getopts "i:ii:f:s" opt; do
    case $opt in
        i) it=$OPTARG;;
        ii) iit=$OPTARG;;
        f) formID=$OPTARG;;
        s) solver=false
          break;;
        \?) echo "Invalid option -$OPTARG" >&2;;
    esac
done
echo "Running $it iterations"

python random_data.py

rm min_cost.log gurobi.log min_err.log gur_err.log

# Generate the logs
if [ "$solver" == true ]; then
    echo "Using minCost"
    ./aion-cli solve minCost ${formID} > min_cost.log 2> min_err.log
    MINCOST=$(grep -F -- "User:" min_cost.log | cut -d ' ' -f 2-10) 
    COSTMINCOST=$(grep -F -- "Min Cost:" min_cost.log | cut -d ' ' -f 4) 
    for i in $(seq 1 $iit); do
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
fi
if [ "$solver" == false ]; then
    echo "Using gurobi"
    ./aion-cli solve gurobi ${formID} > gurobi.log 2> gurobi.log

    GUROBI=$(grep -F -- "group" gurobi.log | cut -d ' ' -f 1-10)

    COSTGUROBI=$(grep -F -- "Min Cost:" gurobi.log | cut -d ' ' -f 3) 

    for i in $(seq 1 $iit); do
      ./aion-cli solve gurobi ${formID} > gurobi.log 2> gur_err.log

      GUROBI2=$(grep -F -- "group" gurobi.log | cut -d ' ' -f 1-10)

      COSTGUROBI2=$(grep -F -- "Min Cost:" gurobi.log | cut -d ' ' -f 3)

      if (( $(echo "$COSTGUROBI > $COSTGUROBI2" |bc -l) )) || (( $(echo "$COSTGUROBI < $COSTGUROBI2" |bc -l) )); then
          echo "Mismatch in iteration $i"
          echo "Gurobi: $COSTGUROBI"
          echo "Gurobi2: $COSTGUROBI2"
      fi

      if [ "$GUROBI" != "$GUROBI2" ]; then
          echo "Mismatch in iteration $i"
          echo "Gurobi: $GUROBI"
          echo "Gurobi: $GUROBI2"
          exit 1
      fi
    done
fi

rm min_cost.log gurobi.log min_err.log gur_err.log form.csv

echo "All iterations passed"
