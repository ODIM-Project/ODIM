#!/bin/bash

#(C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
#under the License.

declare SOURCE_DIR
declare DESTINATION_DIR

create_packages()
{
	packages_list=(odimra-platformconfig odimra-config configure-hosts odimra-pv-pvc redis-ha redis zookeeper-ha zookeeper kafka-secret zookeeper-secret odimra-secret reloader update task systems managers fabrics events api aggregation telemetry account-session kafka-ha kafka odimra-k8s-access-config etcd etcd-ha)

	for chart in "${packages_list[@]}"; do
		helm package ${SOURCE_DIR}/$chart -d ${DESTINATION_DIR}
		if [[ $? -ne 0 ]]; then
			echo "[$(date)] -- ERROR -- Failed to create $chart helm package, exiting"
			exit 1
		fi
	done
}

create_package()
{
	helm package ${SOURCE_DIR} -d ${DESTINATION_DIR}
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- Failed to create $(basename ${SOURCE_DIR}) helm package, exiting"
		exit 1
	fi
}

usage()
{
        echo -e "$(basename $BASH_SOURCE) (<helm_chart_source_dir_path>|<helm_charts_source_dir_path>) <dest_dir_path_to_store_pkgs>"
        exit 1
}

##############################################
###############  MAIN  #######################
##############################################

if [[ $# -ne 2 ]]; then
        usage
fi

SOURCE_DIR=$1
DESTINATION_DIR=$2

if [[ -z ${SOURCE_DIR} ]] || [[ ! -d ${SOURCE_DIR} ]]; then
	echo "[$(date)] -- ERROR -- invalid directory path [${SOURCE_DIR}] passed, exiting!!!"
	usage
fi

if [[ -z ${DESTINATION_DIR} ]] || [[ ! -d ${DESTINATION_DIR} ]]; then
        echo "[$(date)] -- ERROR -- invalid directory path [${DESTINATION_DIR}] passed, exiting!!!"
        usage
fi

if [[ -d ${SOURCE_DIR}/templates ]]; then
	create_package
else
	create_packages
fi

exit 0
