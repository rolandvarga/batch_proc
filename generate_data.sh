#!/bin/bash
echo "{\"objects\": [ " > data.json 
for i in `seq 0 49999`; do echo "{ \"id\": \"object_$i\", \"seq\": $i, \"data\": \"Data for ID object_$i\" },"; done | sort --random-sort | sed '$ s/.$//' >> data.json
echo "]}" >> data.json
