#!/bin/bash
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

# Script is for generating certificate and private key
# for Client mode connection usage only

if [[ -z $FQDN ]]; then
	echo "[ERROR] Set FQDN to environment of the host machine using the following command: export FQDN=<user_preferred_fqdn_for_host>"
	exit 1
fi
if [[ -z $HOSTIP ]]; then
	echo "[ERROR] Set the environment variable, HOSTIP to the IP address of your system using following coomand: export HOSTIP=<ip_address_of_your_system>"
        exit 1
fi
RootServiceUUID=$(uuidgen)
sed -i "s#\"RootServiceUUID\".*#\"RootServiceUUID\": \"${RootServiceUUID}\",#" build/odimra/odimra_config/odimra_config.json

docker image ls | grep odimra_builddep > /dev/null 2>&1
if [ ${?} -eq 0 ]; then
	echo "builddep already exists"
	exit 0
else
	cd build && docker build -t odimra_builddep:tst -f odimra/Dockerfile.builddep .
	exit 0
fi

