#!/bin/bash

# mandatory to systemctl service
. /data/ocms2/dist/env.sh

if [ ! -d "$LOGDIR" ]; then
	echo "Create log directory"
    sudo mkdir -p $LOGDIR
fi
if [ ! -d "$DBDIR" ]; then
	echo "Create db directory"
    sudo mkdir -p $DBDIR    
fi

sudo chown -R $USER $LOGDIR
sudo chown -R $USER $DBDIR
echo "OCMS process started."
if [ "$1" ==  "init" ]
 then
 	if [ -f "auth.log" ]
	then
	   rm auth.log
	fi
 	echo "start auth init"
    CMD="$OCMSPATH/ocms init "
 else
 	echo "start auth"
    CMD="$OCMSPATH/ocms "
fi

eval "$CMD"

#read -rst 0.5
#tail -f $LOGDIR/ocms.log
