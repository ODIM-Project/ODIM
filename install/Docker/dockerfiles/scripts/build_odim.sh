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

LIST=`ls | grep -v 'lib-rest-client' | grep -E '^svc-|^plugin-|add-hosts'`
echo $LIST
flag=0
for i in $LIST; do
    cd $i
    go mod download
    go mod vendor
    go build -i .
    if [ $? -eq 0 ]; then
        echo Successfully build $i service
    else
        echo Failed to build $i service
	arr+=$i,;
	flag=1
    fi
    cd ../
done
if [[ "$flag" -eq 1 ]]; then
	echo "Failed to build $arr services"
	exit 1
fi
