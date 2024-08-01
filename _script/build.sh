#!/bin/bash
echo "- Compiling websocket backend, please wait"
CGO_ENABLED=0 GOOS=linux go build -o bin/sportnex-websocket ./cmd
if [ $? -ne 0 ]; then
    echo "- Building failed"
    exit -1
fi
echo "- Websocket backend compiled, enjoy! ;)"
