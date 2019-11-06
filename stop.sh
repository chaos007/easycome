#!/bin/bash

ps -ef|grep _test|grep -v grep|awk '{print "kill -15 "$2}'|sh

ps aux|grep _test|grep -v grep

echo "shell finished."



