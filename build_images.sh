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

eval_cmd_exec()
{
	echo "Building $2 image"
	cmd="$1 $build_args"
        output=$($cmd 1> /dev/null)
        status=$?
        if [[ $status -ne 0 ]]; then
		echo $output
                echo "$2 image creation failed"
                exit 1
        fi
	echo "$2 image build was successful"
        return $status
}

# base image for building ODIMRA services image
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.odim -t odim:2.0 ." "odim"

# third party docker images
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.etcd -t etcd:1.16 ." "etcd"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.redis -t redis:2.0 ." "redis"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.kafka -t kafka:1.0 ." "kafka"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.zookeeper -t zookeeper:1.0 ." "zookeeper"

# ODIMRA services image
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.accountSession -t account-session:2.0 ." "account session"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.aggregation -t aggregation:2.0 ". "aggregation"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.api -t api:2.0 ." "api"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.events -t events:2.0 ." "events"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.fabrics -t fabrics:2.0 ." "fabrics"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.telemetry -t telemetry:1.0 ." "telemetry"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.managers -t managers:2.0 ." "managers"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.systems -t systems:2.0 ." "systems"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.task -t task:2.0 ." "task"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.update -t update:2.0 ." "update"

# ODIMRA plugins image
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.urplugin -t urplugin:2.0 ." "urplugin"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.grfplugin -t grfplugin:2.0 ." "grfplugin"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.dellplugin -t dellplugin:1.0 ." "dellplugin"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.lenovoplugin -t lenovoplugin:1.0 ." "lenovoplugin"
