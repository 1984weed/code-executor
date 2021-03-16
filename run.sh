#!/bin/bash

# PYCO="$(cat ./sample.py | sed 's/"//g')"
# PYCO=$(cat ./sample.py | sed -E ':a;N;$!ba;s/\r{0,1}\n/\\n/g')
PYCO=$(awk '{printf "%s\\n", $0}' ./sample.py | sed 's/"/\\"/g')
# P=echo "$PYCO"

JSON=$(cat <<-END
{
    "language": "python",
    "content": "$PYCO",
    "inputCount": 3,
    "argumentCount": 2,
    "inputs": [["[2,7,11,15]", "9"], ["[3,2,4]", "6"], ["[3,3]", "6"]],
    "expectOutputs": ["[0,1]", "[1,2]", "[0,1]"]
}
END)
echo $JSON | go run .
