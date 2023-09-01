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
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.odim -t odim:6.0 ." "odim"

# base image for building composition service Image
#eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.cs -t cs:1.0 ." "cs"

# third party docker images
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.etcd -t etcd:1.16 ." "etcd"

eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.redis -t redis:5.0 ." "redis"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.kafka -t kafka:3.0 ." "kafka"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.zookeeper -t zookeeper:3.0 ." "zookeeper"


# ODIMRA services image
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.accountSession -t account-session:5.0 ." "account session"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.aggregation -t aggregation:6.0 ". "aggregation"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.api -t api:7.0 ." "api"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.events -t events:7.0 ." "events"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.fabrics -t fabrics:5.0 ." "fabrics"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.licenses -t licenses:4.0 ." "licenses"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.managers -t managers:7.0 ." "managers"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.systems -t systems:7.0 ." "systems"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.task -t task:6.0 ." "task"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.telemetry -t telemetry:4.0 ." "telemetry"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.update -t update:6.0 ." "update"
#eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.compositionService -t composition-service:1.0 ." "composition service"

# ODIMRA plugins image
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.urplugin -t urplugin:3.4 ." "urplugin"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.grfplugin -t grfplugin:4.1 ." "grfplugin"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.dellplugin -t dellplugin:3.1 ." "dellplugin"
eval_cmd_exec "/usr/bin/docker build -f install/Docker/dockerfiles/Dockerfile.lenovoplugin -t lenovoplugin:1.3 ." "lenovoplugin"
