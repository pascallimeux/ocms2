#!/bin/bash
echo "Build OCMS"
go build -ldflags "-s" $GOPATH/src/github.com/pascallimeux/ocms2/ocms.go
if [ ! -d "$GOPATH/src/github.com/pascallimeux/ocms2/dist" ]; then
  mkdir $GOPATH/src/github.com/pascallimeux/ocms2/dist
fi
mv ocms $GOPATH/src/github.com/pascallimeux/ocms2/dist/ocms
cp $GOPATH/src/github.com/pascallimeux/ocms2/settings.toml $GOPATH/src/github.com/pascallimeux/ocms2/dist/settings.toml
cp $GOPATH/src/github.com/pascallimeux/ocms2/modules/auth/authsettings.toml $GOPATH/src/github.com/pascallimeux/ocms2/dist/authsettings.toml
cp *.sh $GOPATH/src/github.com/pascallimeux/ocms2/dist
chmod u+x $GOPATH/src/github.com/pascallimeux/ocms2/dist/*.sh
