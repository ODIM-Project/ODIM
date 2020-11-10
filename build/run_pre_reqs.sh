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

# perform pre-requisites required for creating docker image
pre_reqs()
{
	if [[ -z ${ODIMRA_GROUP_ID} ]] || [[ -z ${ODIMRA_USER_ID} ]]; then
		echo "[$(date)] -- ERROR -- ODIMRA_GROUP_ID or ODIMRA_USER_ID is not set, exiting"
		exit 1
	fi

	if [[ -n "$(getent passwd odimra 2>&1)" ]]; then
		echo "[$(date)] -- INFO  -- user odimra exists"
		sudo userdel odimra
	fi
	if [[ -n "$(getent group odimra 2>&1)" ]]; then
		echo "[$(date)] -- INFO  -- group odimra exists"
		sudo groupdel odimra
	fi
	sudo groupadd -g ${ODIMRA_GROUP_ID} -r odimra
	sudo useradd -u ${ODIMRA_USER_ID} -r -M -g odimra odimra

	if [[ -z $FQDN ]]; then
		echo "[ERROR] Set FQDN to environment of the host machine using the following command: export FQDN=<user_preferred_fqdn_for_host>"
		exit 1
        fi
	
        if [[ -z $HOSTIP ]]; then
		echo "[ERROR] Set the environment variable, HOSTIP to the IP address of your system using following coomand: export HOSTIP=<ip_address_of_your_system>"
        	exit 1
        fi

}

##############################################
###############  MAIN  #######################
##############################################

pre_reqs
