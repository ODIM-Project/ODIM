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



## URP deployment instructions

For deploying the Unmanaged Racks plugin and adding the plugin to the Resource Aggregator for ODIM framework, refer to the "Deploying the Unmanaged Rack Plugin" section in the [Resource Aggregator for Open Distributed Infrastructure Management™ Readme](https://github.com/ODIM-Project/ODIM/blob/main/README.md).



## Create RackGroup

```
POST https://{odim_host}:{port}/redfish/v1/Chassis
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
POST https://{odim_host}:{port}/redfish/v1/Chassis
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
PATCH https://{odim_host}:{port}/redfish/v1/Chassis/3061416c-5144-5d96-9ec8-69d670a89a8b
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
PATCH https://{odim_host}:{port}/redfish/v1/Chassis/3061416c-5144-5d96-9ec8-69d670a89a8b
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
DELETE https://{odim_host}:{port}/redfish/v1/Chassis/3061416c-5144-5d96-9ec8-69d670a89a8b
Authorization:Basic YWRtaW46T2QhbTEyJDQ=
```

## Delete RackGroup
```
DELETE https://{odim_host}:{port}/redfish/v1/Chassis/1be678f0-86dd-58ac-ac38-16bf0f6dafee
Authorization:Basic YWRtaW46T2QhbTEyJDQ=
```
