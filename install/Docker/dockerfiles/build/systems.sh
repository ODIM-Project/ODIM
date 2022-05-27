#!/bin/bash -x
# (C) Copyright [2022] Hewlett Packard Enterprise Development LP
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
protos=("account" "aggregator" "auth" "chassis" "events" "fabrics" "managers" "role" "session" "systems" "task" "telemetry" "update" "compositionservice" "licenses")
for str in ${protos[@]}; do
  proto_path="$(pwd)/lib-utilities/proto/$str"
  proto_file_name="$str.proto"
  if [ $str == 'auth' ]
  then
    proto_file_name="odim_auth.proto"
  fi
  if [ $str == 'compositionservice' ]
  then
    proto_file_name="composition_service.proto"
  fi
  protoc --go_opt=M$proto_file_name=./ --go_out=plugins=grpc:$proto_path --proto_path=$proto_path $proto_file_name
done

LIST=("svc-systems")
echo $LIST
for i in $LIST; do
    cd $i
    go build -i .
    if [ $? -eq 0 ]; then
        echo Successfully build $i service
    else
        echo Failed to build $i service
    fi
    cd ../
done
