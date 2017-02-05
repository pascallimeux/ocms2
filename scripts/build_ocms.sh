#!/bin/bash
echo "Build OCMS"
. env.sh

go build -ldflags "-s" $SRCPATH/ocms.go

if [ ! -d "$SRCBIN" ]; then
  sudo mkdir -p $SRCBIN
  sudo chown -R $USER $DATAREPO
fi


mv ocms $SRCBIN/ocms
cp $SRCPATH/ocms.toml $SRCBIN/ocms.toml
cp $SRCPATH/modules/auth/auth.toml $SRCBIN/auth.toml
cp *.sh $SRCBIN
chmod u+x $SRCBIN/*.sh
