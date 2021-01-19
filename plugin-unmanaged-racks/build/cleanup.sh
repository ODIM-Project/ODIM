#!/bin/bash
if [ -a build/docker-compose.yml ]; then
	cd build
	docker-compose down
	LIST=`docker image ls | grep -E 'ur-plugin|plugin_ur_builddep' | awk '{print $3}'`
	docker rmi $LIST
        echo "Cleanup Done"
        cd ../
        exit 0
else
	echo "docker-compose.yml doesn't exist, are you in the ur-plugin directory?"
	exit 1
fi

