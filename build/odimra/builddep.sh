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

LIST=`ls -R | grep -v 'lib-rest-client' | grep -E '^svc-|^plugin-'` 
echo $LIST
for i in $LIST; do
	cd $i
	go mod download
	go mod vendor
	go build -i -race .
	if [ $? -eq 0 ]; then
		echo Build for odimra service/lib dependencies $i are Successful !!!!
		arr1+=$i;
	else
		echo Build for odimra service/lib dependency $i Failed !!!!
		arr2+=$i;
		flag=0
	fi
	cd ../
done
