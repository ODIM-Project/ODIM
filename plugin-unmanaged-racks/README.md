<!-- 
 Copyright (c) 2020 Intel Corporation

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

 http://www.apache.org/licenses/LICENSE-2.0
 
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
-->

# Unmanaged Racks Plugin (URP) 

This folder contains implementation of Unamanaged Racks Plugin(URP) for ODIMRA. 
Plugin exposes narrowed obligatory REST API described by Plugin Developer’s Guide (PDG).
In addition URP exposes following REST endpoints:

* `GET /ODIM/v1/Chassis` - return collection of unmanaged Chassis(RackGroups/Racks)
* `GET /ODIM/v1/Chassis/{id}` - return instance of unmanaged Chassis(RackGroups/Racks)
* `POST /ODIM/v1/Chassis` - creates new unmanaged Chassis(RackGroups/Racks) 
* `DELETE /ODIM/v1/Chassis/{id}` - deletes existing unmanaged Chassis(RackGroups/Racks)
* `PATCH /ODIM/v1/Chassis/{id}` - updates existing unmanaged Chassis(RackGroups/Racks)

Full specification of URP is available here: https://wiki.odim.io/display/HOME/Plugin+for+Unmanaged+Racks.

Please be aware this plugin is still under development, and some features might be missing.

## Build 

Build URP using following command:
```
cd plugin-unmanaged-racks
make build
``` 

Run URP using run target:
```
make run
```

## Register URP in ODIMRA

1. Make `https://localhost:45000/redfish/v1/AggregationService/ConnectionMethods` endpoint exposes connection method required by URP plugin. `ConnectionMethodVariant` should be `Compute:BasicAuth:URP_v1.0.0`.

```
{
  "@odata.type": "#ConnectionMethod.v1_0_0.ConnectionMethod",
  "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/bea4ab96-edd1-4cce-b57c-f83e218e97b6",
  "@odata.context": "/redfish/v1/$metadata#ConnectionMethod.v1_0_0.ConnectionMethod",
  "Id": "bea4ab96-edd1-4cce-b57c-f83e218e97b6",
  "Name": "Connection Method",
  "Severity": "OK",
  "ConnectionMethodType": "Redfish",
  "ConnectionMethodVariant": "Compute:BasicAuth:URP_v1.0.0",
  "Links": {
    "AggregationSources": []
  }
}
```

2. Execute plugin registration request:
```
POST https://odimra.local.com:45000/redfish/v1/AggregationService/AggregationSources
Authorization:Basic YWRtaW46T2QhbTEyJDQ=

{
 "HostName": "odimra.local.com:45003",
 "Password":"Od!m12$4",
 "UserName":"admin",
 "Links": {
   "ConnectionMethod":{
       "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/bea4ab96-edd1-4cce-b57c-f83e218e97b6"
   }
 }
}
```

## Create RackGroup
```
POST https://localhost:45000/redfish/v1/Chassis
Authorization:Basic YWRtaW46T2QhbTEyJDQ=

{
  "ChassisType": "RackGroup",
  "Description": "My RackGroup",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/675560ae-e903-41d9-bfb2-561951999999"
      }
    ]
  },
  "Name": "RG2"
}
```
## Create Rack
```
POST https://odimra.local.com:45000/redfish/v1/Chassis
Authorization:Basic YWRtaW46T2QhbTEyJDQ=

{
  "ChassisType": "Rack",
  "Description": "rack no 1",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/675560ae-e903-41d9-bfb2-561951999999"
      }
    ],
    "ContainedBy": [
       {"@odata.id":"/redfish/v1/Chassis/1be678f0-86dd-58ac-ac38-16bf0f6dafee"}
    ]
  },

  "Name": "RACK#1"
}

```
## Attach selected Chassis under Rack
```
PATCH https://odimra.local.com:45000/redfish/v1/Chassis/3061416c-5144-5d96-9ec8-69d670a89a8b
Authorization:Basic YWRtaW46T2QhbTEyJDQ=
Content-Type: application/json

{
  "Links": {
    "Contains": [
      {
        "@odata.id": "/redfish/v1/Chassis/46db63a9-2dcb-43b3-bdf2-54ce9c42e9d9:1"
      }
    ]
  }
}
```

## Detach Chassis from Rack
```
PATCH https://odimra.local.com:45000/redfish/v1/Chassis/3061416c-5144-5d96-9ec8-69d670a89a8b
Authorization:Basic YWRtaW46T2QhbTEyJDQ=
Content-Type: application/json

{
  "Links": {
    "Contains": []
  }
}
```

## Delete Rack
```
DELETE https://odimra.local.com:45000/redfish/v1/Chassis/3061416c-5144-5d96-9ec8-69d670a89a8b
Authorization:Basic YWRtaW46T2QhbTEyJDQ=
```

## Delete RackGroup
```
DELETE https://odimra.local.com:45000/redfish/v1/Chassis/1be678f0-86dd-58ac-ac38-16bf0f6dafee
Authorization:Basic YWRtaW46T2QhbTEyJDQ=
```
