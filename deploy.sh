#!/bin/bash
host=$1
port=$2
remote="circleci@${host}"
echo "deploying..."
ssh -p $port $remote sudo systemctl stop logmania;
RETVAL=$?
[ $RETVAL -ne 0 ] && exit 1
scp -q -P $port ~/.go_workspace/bin/logmania $remote:~/bin/logmania;
RETVAL=$?
ssh -p $port $remote sudo systemctl start logmania;
[ $RETVAL -eq 0 ] && RETVAL=$?
[ $RETVAL -ne 0 ] && exit 1
echo "deployed"
