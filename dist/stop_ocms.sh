#!/bin/bash
PID=`pidof ocms`
if [ -n "$PID" ]
then
   kill -9 $PID
   echo "OCMS process stopped"
else
   echo "Could not send SIGTERM to kill OCMS, probably it does not work:" >&2
fi