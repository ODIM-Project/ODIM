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
sleep 3
echo "Checking if default entries already present"
redis_password=$(openssl pkeyutl -decrypt -in cipher -inkey ${ODIMRA_RSA_PRIVATE_FILE} -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512)
redis-cli -a ${redis_password} -h ${master} -p ${REDIS_HA_REDIS_SERVICE_PORT} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} <<HERE
exists "role:Administrator"
HERE
