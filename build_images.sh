#!/bin/bash
#(C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.

export DOCKER_BUILDKIT=1

if [[ -z ${ODIMRA_GROUP_ID} ]] || [[ -z ${ODIMRA_USER_ID} ]]; then
	echo "[$(date)] -- ERROR -- ODIMRA_GROUP_ID or ODIMRA_USER_ID is not set, exiting"
	exit 1
fi

if [[ -n ${http_proxy} ]] || [[ -n ${https_proxy} ]]; then 
	echo "Adding proxy to build arguments"
        build_args="--build-arg ODIMRA_USER_ID=${ODIMRA_USER_ID} \
	            --build-arg ODIMRA_GROUP_ID=${ODIMRA_GROUP_ID} \
		    --build-arg http_proxy=${http_proxy} \
		    --build-arg https_proxy=${https_proxy}"
else
	build_args="--build-arg ODIMRA_USER_ID=${ODIMRA_USER_ID} \
                    --build-arg ODIMRA_GROUP_ID=${ODIMRA_GROUP_ID}"
fi

/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.odim -t odim:1.0 .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.accountSession -t account-session:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.aggregation -t aggregation:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.api -t api:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.consul -t consul:1.6 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.events -t events:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.fabrics -t fabrics:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.grfplugin -t grf-plugin:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.kafka -t odim_kafka:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.managers -t managers:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.redis -t redis:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.systems -t systems:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.task -t task:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.update -t update:1.0 $build_args .
/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.zookeeper -t odim_zookeeper:1.0 $build_args .
