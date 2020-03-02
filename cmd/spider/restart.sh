#!/usr/bin/env bash
sudo killall -9 spider.exe
nohup ./spider.exe  > /dev/null  2>&1 &
