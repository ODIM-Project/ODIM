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

declare ACTION
declare DIR_PATH
declare -A images_list=(\
		["account-session"]="2.0" \
		["aggregation"]="2.0" \
		["api"]="2.0" \
		["etcd"]="1.16" \
		["events"]="2.0" \
		["fabrics"]="2.0" \
		["kafka"]="1.0" \
		["managers"]="2.0" \
		["redis"]="2.0" \
		["systems"]="2.0" \
		["task"]="2.0" \
		["telemetry"]="1.0"\
		["update"]="2.0" \
		["zookeeper"]="1.0" \
		)

eval_cmd_exec()
{
        if [[ $# -lt 3 ]]; then
                echo "[$(date)] -- ERROR -- eval_cmd_exec syntax error $2"
                exit 1
        fi
        if [[ $1 -eq 0 ]]; then
		echo "[$(date)] -- INFO  -- $4"
	else
		echo "[$(date)] -- ERROR -- $2"
		echo "$3"
        fi
}

load_images()
{
        for image in "${!images_list[@]}"; do
		output=$(docker load -i ${DIR_PATH}/${image}.tar 2>&1)
		eval_cmd_exec $? "failed to load ${DIR_PATH}/${image}.tar" \
		"$output" "Successfully loaded ${DIR_PATH}/${image}.tar"
        done

	echo "[$(date)] -- INFO  -- Cleaning up any dangling images created while loading images"
	docker rmi $(/usr/bin/docker images -f "dangling=true" -q) > /dev/null 2>&1
}
 
save_images()
{
        for image in "${!images_list[@]}"; do
		output=$(docker save -o ${DIR_PATH}/${image}.tar ${image}:${images_list[${image}]} 2>&1)
		eval_cmd_exec $? "failed to save ${image}:${images_list[${image}]} as ${DIR_PATH}/${image}.tar" \
		"$output" "Successfully saved ${image}:${images_list[${image}]} as ${DIR_PATH}/${image}.tar"
        done
}
 
usage()
{
        echo -e "$(basename $BASH_SOURCE) load|save <dir_path_to_save_or_load_from>"
        exit 1
}
 
##############################################
###############  MAIN  #######################
##############################################

ACTION=$1
DIR_PATH=$2

if [[ $# -gt 2 ]]; then
        usage
fi

if [[ -z ${DIR_PATH} ]]; then
	DIR_PATH="."
fi

case $1 in
        load)
                load_images
                ;;
        save)
                save_images
                ;;
        *)
                usage
                ;;
esac
 
exit 0
