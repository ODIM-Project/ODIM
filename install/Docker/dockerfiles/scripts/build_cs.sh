#(C) Copyright [2022] American Megatrends International LLC
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

protos=("auth" "compositionservice")

for str in ${protos[@]}; do
  proto_path="$(pwd)/lib-utilities/proto/$str"
  proto_out_path="$(pwd)/svc-composition-service/app/proto/$str" 
  proto_file_name="$str.proto"
  if [ $str == 'auth' ]
  then
    proto_file_name="odim_auth.proto"
  fi
  if [ $str == 'compositionservice' ]
  then
    proto_file_name="composition_service.proto"
  fi
  python3 -m grpc_tools.protoc -I$proto_path --python_out=$proto_out_path --grpc_python_out=$proto_out_path $proto_file_name

  proto_grpc_file_name="${proto_file_name/.proto/"_pb2_grpc.py"}"
  protoc_pb2_name="${proto_file_name/.proto/"_pb2"}"
  proto_grpc_file=$proto_out_path/$proto_grpc_file_name
  if [[ -f "$proto_grpc_file" ]]; then
    sed -i 's/import '$protoc_pb2_name'/import proto.'$str'.'$protoc_pb2_name'/gI' $proto_grpc_file
  fi
done

LIST=`ls | grep -E '^svc-'`
echo $LIST
for i in $LIST; do
    cd $i
    if [[ "$i" == "svc-composition-service" ]]; then
        /bin/bash build.sh
    fi
done