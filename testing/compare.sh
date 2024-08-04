it=1
formID=""

setData=true

while getopts "i:f:" opt; do
    case $opt in
        i) it=$OPTARG;;
        f) formID=$OPTARG;;
        rand) setData=false;;
        \?) echo "Invalid option -$OPTARG" >&2;;
    esac
done
echo "Running $it iterations"

if [ "$setData" = false ]; then
    python random_data.py
fi

rm min_cost.log gurobi.log min_err.log gur_err.log
# Generate the logs
for i in $(seq 1 $it); do
    if [ "$setData" = true ]; then
        python random_data.py
    fi
    ./aion-cli solve minCost ${formID} > min_cost.log 2> min_err.log
    ./aion-cli solve gurobi ${formID} > gurobi.log 2> gur_err.log

    MINCOST=$(grep -F -- "User:" min_cost.log | cut -d ' ' -f 2-15) 
    GUROBI=$(grep -F -- "group" gurobi.log | cut -d ' ' -f 1-15)

    COSTMINCOST=$(grep -F -- "Min Cost:" min_cost.log | cut -d ' ' -f 4)
    COSTGUROBI=$(grep -F -- "Min Cost:" gurobi.log | cut -d ' ' -f 3)

    if (( $(echo "$COSTMINCOST > $COSTGUROBI" |bc -l) )) || (( $(echo "$COSTMINCOST < $COSTGUROBI" |bc -l) )); then
        echo "Mismatch in iteration $i"
        echo "MinCost: $COSTMINCOST"
        echo "Gurobi: $COSTGUROBI"
    fi

    if [ "$MINCOST" != "$GUROBI" ]; then
        echo "Mismatch in iteration $i"
        echo "MinCost: $MINCOST"
        echo "Gurobi: $GUROBI"
        exit 1
    fi
done

rm min_cost.log gurobi.log min_err.log gur_err.log form.csv

echo "All iterations passed"
