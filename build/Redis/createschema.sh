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

RETURN=`checkdb.sh | grep '0' > /dev/null`
if [ $? -eq 0 ]; then
redis-cli -p 6380 <<HERE
Set  "registry:assignedprivileges"  '{"List":["Login", "ConfigureManager", "ConfigureUsers", "ConfigureSelf", "ConfigureComponents"]}'
Set "roles:redfishdefined"  '{"List":["Administrator", "Operator", "ReadOnly"]}'
Set "User:admin"  '{"UserName":"admin","Password":"O01bKrP7Tzs7YoO3YvQt4pRa2J_R6HI34ZfP4MxbqNIYAVQVt2ewGXmhjvBfzMifM7bHFccXKGmdHvj3hY44Hw==","RoleId":"Administrator", "AccountTypes":["Redfish"]}'
Set "role:Administrator"  '{"@odata.type":"","RoleId":"Administrator","Name":"","Description":"","IsPredefined":true,"AssignedPrivileges":["ConfigureSelf","Login","ConfigureUsers","ConfigureComponents","ConfigureManager"],"OemPrivileges":null,"@odata.context":"","@odata.id":""}'
Set "role:Operator"  '{"@odata.type":"","RoleId":"Operator","Name":"","Description":"","IsPredefined":true,"AssignedPrivileges":["ConfigureSelf","Login","ConfigureComponents"],"OemPrivileges":null,"@odata.context":"","@odata.id":""}'
Set "role:ReadOnly"  '{"@odata.type":"","RoleId":"ReadOnly","Name":"","Description":"","IsPredefined":true,"AssignedPrivileges":["ConfigureSelf","Login"],"OemPrivileges":null,"@odata.context":"","@odata.id":""}'
ZAdd "Subscription" 0 '{"UserName":"","SubscriptionID":"0","Destination":"","Name":"default","Context":"","EventTypes":["Alert"],"MessageIds":null,"Protocol":"Redfish","SubscriptionType":"RedfishEvent","EventFormatType":"","SubordinateResources":true,"ResourceTypes":null,"OriginResources":[],"Hosts":[],"DeliveryRetryPolicy":"RetryForever"}'
keys *
SAVE
HERE
fi
