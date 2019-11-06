#!/bin/bash

sh ./stop.sh

go build -o agent_test ./agent
go build -o game_test ./game
go build -o unique_test ./unique
go build -o center_test ./center

echo "process list:"
ps aux|grep _test|grep -v grep


#   ll /tmp/lol_*.log|awk '{print $9}'|xargs rm -rf
