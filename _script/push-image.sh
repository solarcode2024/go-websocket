#!/bin/bash
version=`cat version`
if [ "$1" = "dev" ]; then 
    version=$version"_dev"
fi

echo "- Pushing lpiexecutive/sportnex_websocket_backend:$version, please wait."

docker push lpiexecutive/sportnex_websocket_backend:$version
if [ $? -ne 0 ]; then
    echo "* Push Failed."
    exit -1
fi

echo "- Image pushed to hub!"
