#!/bin/bash
echo "Build OCMS"
go build -ldflags "-s" $GOPATH/src/github.com/pascallimeux/ocms/main.go
if [ ! -d "$GOPATH/src/github.com/pascallimeux/ocms/dist" ]; then
  mkdir $GOPATH/src/github.com/pascallimeux/ocms/dist
fi
mv main $GOPATH/src/github.com/pascallimeux/ocms/dist/ocms
cp $GOPATH/src/github.com/pascallimeux/ocms/config/config.json $GOPATH/src/github.com/pascallimeux/ocms/dist/config.json
cp *.sh $GOPATH/src/github.com/pascallimeux/ocms/dist
chmod u+x $GOPATH/src/github.com/pascallimeux/ocms/dist/*.sh
