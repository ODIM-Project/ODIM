#!/bin/bash
docker image ls | grep plugin_ur_builddep > /dev/null 2>&1
if [ ${?} -eq 0 ]; then
        echo "builddep already exists"
        exit 0
else
        cd build && docker build -t plugin_ur_builddep:tst -f Dockerfile.builddep .
        exit 0
fi

