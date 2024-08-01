#!/bin/bash
version=`cat version`

# read input "dev"
if [ "$1" = "dev" ]; then 
    version=$version"_dev"
fi

echo "- Compiling sportnex-websocket:$version, please wait."
CGO_ENABLED=0 GOOS=linux go build -o bin/sportnex-websocket ./cmd
if [ $? -ne 0 ]; then
    echo "* Build Failed."
    exit -1
fi
echo "- sportnex-websocket:$version compiled!"

echo "- Build docker image"
docker build -t lpiexecutive/sportnex_websocket_backend:$version .
echo "- Build docker image OK! (sportnex_websocket_backend:$version)"
