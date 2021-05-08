#!/bin/bash

PYCO=$(awk '{printf "%s\\n", $0}' ./sample.py | sed 's/"/\\"/g')

docker run --name python -d -i baracode/python:latest 

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

echo $JSON | docker attach  python
docker logs python

docker stop python
docker rm -f python