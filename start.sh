#!/bin/bash

go run /master/master.go &

#slave of first level
for (( i=0; i<4; i++ )); do
   go run /slave/slave.go &
done

#slave of second level
#for (( i=0; i<4; i++ )); do
#   go run /slave/secondLevelSlave.go &
#done