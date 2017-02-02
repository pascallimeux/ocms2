#!/bin/bash
LOGDIR="./logs"
DBDIR="./db" 
USER="pascal"
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
     rm auth.log
 	 echo "start auth init"
     ./ocms init &
 else
 	 echo "start auth"
     ./ocms &
fi

read -rst 0.5
tail -f ./logs/ocms.log
