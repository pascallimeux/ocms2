#!/bin/bash
echo "Build OCMS"
SRCPATH="$GOPATH/src/github.com/pascallimeux/ocms2"
SRCBIN="/data/ocms2/dist"
go build -ldflags "-s" $SRCPATH/ocms.go
if [ ! -d "$SRCBIN" ]; then
  mkdir $SRCBIN
fi
mv ocms $SRCBIN/ocms
cp $SRCPATH/ocms.toml $SRCBIN/ocms.toml
cp $SRCPATH/modules/auth/auth.toml $SRCBIN/auth.toml
cp *.sh $SRCBIN
chmod u+x $SRCBIN/*.sh
