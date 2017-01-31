#!/bin/bash
LOGDIR="/var/log/ocms"
USER="pascal"
if [ ! -d "$LOGDIR" ]; then
	echo "Create log directory"
    sudo mkdir -p $LOGDIR
fi
    sudo chown -R $USER $LOGDIR
echo "OCMS process started."
./ocms &
