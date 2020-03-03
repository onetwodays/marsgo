#!/usr/bin/env bash
APPNAME="spider.exe"
sudo killall -9 ${APPNAME}
nohup ./${APPNAME}  > /dev/null  2>&1 &
ps -ef|grep  ${APPNAME}
