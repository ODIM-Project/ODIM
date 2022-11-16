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

RETURN=`/checkdb.sh | grep '0' > /dev/null`
if [ $? -eq 0 ]; then
	  echo "Updating the db with default entries"
    redis_password=$(openssl pkeyutl -decrypt -in cipher -inkey ${ODIMRA_RSA_PRIVATE_FILE} -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512)
    redis-cli -a ${redis_password} -h ${master} -p ${REDIS_HA_REDIS_SERVICE_PORT} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} <<HERE
Set  "registry:assignedprivileges"  '{"List":["Login", "ConfigureManager", "ConfigureUsers", "ConfigureSelf", "ConfigureComponents"]}'
Set "roles:redfishdefined"  '{"List":["Administrator", "Operator", "ReadOnly"]}'
Set "User:admin"  '{"UserName":"admin","Password":"O01bKrP7Tzs7YoO3YvQt4pRa2J_R6HI34ZfP4MxbqNIYAVQVt2ewGXmhjvBfzMifM7bHFccXKGmdHvj3hY44Hw==","RoleId":"Administrator", "AccountTypes":["Redfish"]}'
Set "role:Administrator"  '{"@odata.type":"","RoleId":"Administrator","Name":"","Description":"","IsPredefined":true,"AssignedPrivileges":["ConfigureSelf","Login","ConfigureUsers","ConfigureComponents","ConfigureManager"],"OemPrivileges":null,"@odata.context":"","@odata.id":""}'
Set "role:Operator"  '{"@odata.type":"","RoleId":"Operator","Name":"","Description":"","IsPredefined":true,"AssignedPrivileges":["ConfigureSelf","Login","ConfigureComponents"],"OemPrivileges":null,"@odata.context":"","@odata.id":""}'
Set "role:ReadOnly"  '{"@odata.type":"","RoleId":"ReadOnly","Name":"","Description":"","IsPredefined":true,"AssignedPrivileges":["ConfigureSelf","Login"],"OemPrivileges":null,"@odata.context":"","@odata.id":""}'
keys *
SAVE
HERE
fi
