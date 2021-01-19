#!/bin/bash
if [ -a build/docker-compose.yml ]; then
        cd build
        docker-compose down
        LIST=`docker image ls | grep -v REPOSITORY | awk '{print $3}'`
        docker rmi $LIST
        echo "Cleanup done"
        cd ../
#       exit 0
else
        echo "docker-compose.yml doesn't exist, are you in the ur-plugin directory?"
        exit 1
fi
sudo rm -rf /var/log/UR_PLUGIN
