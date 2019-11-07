#!/bin/bash

sh ./stop.sh

./agent_test --id=agent1 &>log/agent1.log &
./agent_test --id=agent2 &>log/agent2.log &
./game_test --id=game1 &>log/game1.log &
./game_test --id=game2 &>log/game2.log &
./unique_test --id=unique &>log/unique.log &
./center_test --id=center &>log/center.log &


echo "process list:"
ps aux|grep _test|grep -v grep


#   ll /tmp/lol_*.log|awk '{print $9}'|xargs rm -rf
