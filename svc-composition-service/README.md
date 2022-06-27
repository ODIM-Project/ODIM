#  Composition Service

Resource Aggregator for ODIM exposes Redfish APIs to view and manage ResourceBlocks and ResourceZones that are used for composability. The Redfish `CompositionService` APIs allow you to create and remove ResourceBlocks and ResourceZones.



**Supported endpoints**


|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/CompositionService|GET|`Login` |
|/redfish/v1/CompositionService/ResourceBlocks|GET, POST|`Login` , `ConfigureComponents` |
|/redfish/v1/CompositionService/ResourceBlocks/\{ResourceBlockId\}|GET, DELETE|`Login` , `ConfigureComponents` |
|/redfish/v1/CompositionService/ResourceZones|GET, POST|`Login` , `ConfigureComponents` |
|/redfish/v1/CompositionService/ResourceZones/\{ResourceZoneId\}|GET, DELETE|`Login` , `ConfigureComponents`|
|/redfish/v1/CompositionService/ActivePool|GET|`Login` |
|/redfish/v1/CompositionService/FreePool|GET|`Login` |
|/redfish/v1/CompositionService/Actions/CompositionService.Compose|POST|`ConfigureManager` |



##  Modifying Configurations of composition Service
  
Config file of CompositionService is located at: **odimra/svc-composition-service/app/config/config.py**  

  
**Specific configurations for Composition Service are:**
  
##  Log Location of the Composition Service
  
/var/log/odimra/composition_service.log
  
  





##  Composition Service

|||
|---------------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/CompostionService` |
|**Description** |CompostionService.|
|**Returns** |Composition service resources and theri details.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService'

```

>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#CompositionService",
    "@odata.type": "#CompositionService.v1_2_0.CompositionService",
    "@odata.id": "/redfish/v1/CompositionService",
    "Id": "CompositionService",
    "Name": "Composition Service",
    "Status": {
        "State": "Enabled",
        "Health": "OK"
    },
    "ServiceEnabled": true,
    "AllowOverprovisioning": true,
    "AllowZoneAffinity": true,
    "ResourceBlocks": {
        "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks"
    },
    "ResourceZones": {
        "@odata.id": "/redfish/v1/CompositionService/ResourceZones"
    },
    "ActivePool": {
        "@odata.id": "/redfish/v1/CompositionService/ActivePool"
    },
    "CompositionReservations": {
        "@odata.id": "/redfish/v1/CompositionService/CompositionReservations"
    },
    "FreePool": {
        "@odata.id": "/redfish/v1/CompositionService/FreePool"
    },
    "Actions": {
        "#CompositionService.Compose": {
            "target": "/redfish/v1/CompostionService/Actions/CompositionService.Compose"
        }
    },
    "ReservationDuration": null,
    "Oem": {}
}

```









## Collection of ResourceBlocks
Resource Blocks are the lowest level building blocks for composition requests. Resource Block instance contain the list of components found within the Resource Block instance. For example, if a Resource Block contains 1 Processor and 4 DIMMs, then all of those components will be part of the same composition request, even if only one of them is needed.

|||
|------------------|----------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/CompositionService/ResourceBlocks` |
|**Description** |A collection of resource blocks that participate in composability.|
|**Returns** |Links to ResourceBlock instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceBlocks'

```

>**Sample response body**

```
{
    "@odata.type": "#ResourceBlockCollection.ResourceBlockCollection",
    "Name": "Resource Block Collection",
    "Members@odata.count": 1,
    "Members": [
        {
            "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/0ebefd98-8a88-11ec-9a3b-cabb961e4309"
        }
    ],
    "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks"
}

```









## Single ResourceBlock

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/CompositionService/ResourceBlocks/{ResourceBlockID}` |
|**Description** |JSON schema representing a particular resource block.|
|**Returns** |Details of this resource blocks and links to zones and other resources.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceBlocks/{ResourceBlockID}'

```

>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#ResourceBlock.ResourceBlock.",
    "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/76e3eff0-5d8d-11ec-ad64-56546eecab63",
    "@odata.type": "#ResourceBlock.v1_4_0.ResourceBlock",
    "Client": null,
    "CompositionStatus": {
        "CompositionState": "Composed",
        "MaxCompositions": 1,
        "NumberOfCompositions": 1,
        "Reserved": false,
        "SharingCapable": false
    },
    "ComputerSystems": [
        {
            "@odata.id": "/redfish/v1/Systems/a85ab2c9-2934-43fe-96f4-e0c56c3216ad:1"
        }
    ],
    "Description": "ComputerSystem Block",
    "Id": "76e3eff0-5d8d-11ec-ad64-56546eecab63",
    "Name": "ComputerSystem Block",
    "Pool": "Active",
    "ResourceBlockType": [
        "ComputerSystem"
    ],
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    }
}
```






## Collection of ResourceZones
Resource Zones describe to the client the different composition restrictions of the Resource Blocks reported by the service; Resource Blocks that are reported in the same Resource Zone are allowed to be composed together. This enables the clients to not perform try and fail logic to figure out the different restrictions that are in place for a given implementation.

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/CompositionService/ResourceZones/{ResourceZoneID}`` |
|**Description** |A collection of ResourceZones.|
|**Returns** |Links to ResourceZone instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceZones'

```

>**Sample response body**

```
{
    "@odata.type": "#ZoneCollection.ZoneCollection",
    "Name": "Resource Zone Collection",
    "Members@odata.count": 2,
    "Members": [
        {
            "@odata.id": "/redfish/v1/CompositionService/ResourceZones/385316c6-8e9d-11ec-942f-caea6c362014"
        },
        {
            "@odata.id": "/redfish/v1/CompositionService/ResourceZones/247c8884-8e9e-11ec-942f-caea6c362014"
        }
    ],
    "@odata.id": "/redfish/v1/CompositionService/ResourceZones"
}
	
```








## Single ResourceZone

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/CompositionService/ResourceZones/{ResourceZoneID}` |
|**Description** |JSON schema representing a specific resource zone.|
|**Returns** |Properties of this zone.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceZones/{ResourceZoneID}'

```

>**Sample response body**

```
{
    "@odata.id": "/redfish/v1/CompositionService/ResourceZones/385316c6-8e9d-11ec-942f-caea6c362014",
    "@odata.type": "#Zone.v1_6_0.Zone",
    "Id": "385316c6-8e9d-11ec-942f-caea6c362014",
    "Name": "ResourceZone1",
    "Status": {
        "State": "Enabled",
        "Health": "OK"
    },
    "ZoneType": "ZoneOfResourceBlocks",
    "Links": {
        "ResourceBlocks": [
            {
                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/11260ff0-8b54-11ec-83c3-e6e57d72922b"
            }
        ]
    }
}

```









## Collection of FreePool

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/CompositionService/FreePool`` |
|**Description** |A collection of resource blocks in free pool.|
|**Returns** |Links to the resource blocks that belong to free pool.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/FreePool'

```


>**Sample response body**

```
{
  "@odata.id": "/redfish/v1/CompositionService/FreePool",
  "@odata.type": "#FreePoolCollection.FreePoolCollection",
  "Members": [
    {
      "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/da2f9f40-5997-11ec-beb4-b26e7be3bd29"
    },
    {
      "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/0d1a98dc-5d8b-11ec-a584-0630e493c686"
    }
  ],
  "Members@odata.count": 2,
  "Name": "Free Pool Collection"
}
	
```








## Collection of ActivePool


|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/CompositionService/ActivePool`` |
|**Description** |A collection of ResourceBlocks in Active Pool.|
|**Returns** |Links to the resource blocks that belong to Active pool.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ActivePool'

```



>**Sample response body**

```
{
  "@odata.id": "/redfish/v1/CompositionService/ActivePool",
  "@odata.type": "#ActivePoolCollection.ActivePoolCollection",
  "Members": [
    {
      "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/da2f9f40-5997-11ec-beb4-b26e7be3bd29"
    },
    {
      "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/0d1a98dc-5d8b-11ec-a584-0630e493c686"
    }
  ],
  "Members@odata.count": 2,
  "Name": "Active Pool Collection"
}
	
```







## Creating a ResourceBlock

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/CompositionService/ResourceBlocks` |
|**Description** |This operation creates a resource block.|
|**Returns** |<ul><li>Link to the created ResourceBlock in the `Location` header.</li><li>JSON schema representing the created ResourceBlock.</li></ul>|
|**Response code** |`201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"ComputerSystem ResourceBlock",
   "Description":"ComputerSystem ResourceBlock",
   "ResourceBlockType": [ "ComputerSystem" ], 
   "ComputerSystems":[ 
      { "@odata.id": "/redfish/v1/Systems/4eac063d-60d7-4fe8-a3d8-bf70afb6d228:1" } 
   ]
}'
'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceBlocks'

```





**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(optional\) |Name for the Resource Block. |
|Description|String \(optional\) |Description for the Resource Block. |
|ResourceBlockType| \(required\) |This property shall contain an array of enumerated values that describe the types of resource block. "Compute", "Processor", "Memory", "Network", "Storage", "ComputerSystem", "Expansion", "IndependentResource" are ResourceBlockTypes defined in Redfish. |
|ComputerSystems\[\{| Array \(required\) | Required if ResourceBlockType is ComputerSystm. Represents an array of computer systems that are used for this ResourceBlock. |
|@odata.id\}\]|String \(required\) | Link to a computer system. |




>**Sample response header**

```
HTTP/1.1 201 Created
Allow: "GET", "POST"
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Location: "/redfish/v1/CompositionService/ResourceBlocks/0ebefd98-8a88-11ec-9a3b-cabb961e4309"
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Thu, 10 Feb 2022 15:42:31 GMT
Content-Length: 707


```

>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#ResourceBlock.ResourceBlock.",
    "@odata.type": "#ResourceBlock.v1_4_0.ResourceBlock",
    "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/0ebefd98-8a88-11ec-9a3b-cabb961e4309",
    "Id": "0ebefd98-8a88-11ec-9a3b-cabb961e4309",
    "Name": "ComputerSystem Block",
    "Description": "ComputerSystem Block",
    "ResourceBlockType": [
        "ComputerSystem"
    ],
    "Status": {
        "State": "Enabled",
        "Health": "OK"
    },
    "CompositionStatus": {
        "Reserved": false,
        "CompositionState": "Composed",
        "SharingCapable": false,
        "MaxCompositions": 1,
        "NumberOfCompositions": 1
    },
    "Client": null,
    "Pool": "Active",
    "ComputerSystems": [
        {
            "@odata.id": "/redfish/v1/Systems/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8"
        }
    ]
}


```






## Creating a Resourcezone

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/CompositionService/ResourceZones` |
|**Description** |This operation creates a ResourceZone.|
|**Returns** |<ul><li>Link to the created ResourceZone in the `Location` header.</li><li>JSON schema representing the created ResourceZone.</li></ul>|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
    "Name": "ResourceZone1", 
    "Description": "Resource Zone 1", 
    "ZoneType": "ZoneOfResourceBlocks", 
    "Links": {
        "ResourceBlocks": [
            {"@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/11260ff0-8b54-11ec-83c3-e6e57d72922b"}
        ]
    }
}'
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceZones'

```



**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(optional\) |Name for the ResourceZone.|
|Description|String \(optional\) |Description for the resource zone.|
|ZoneType|String \(required\) |ZoneOfResourceBlocks: "This value shall indicate a zone that contains resources of type ResourceBlock.  This value shall only be used for zones subordinate to the composition service."|
|Links\{|Object \(required\) |Contains references to other resources that are related to the zone. |
|ResourceBlocks \[\{|Array \(required\) | Represents an array of ResourceBlocks that are used for this ResourceZone. |
|@odata.id\]\}|String \(required\) |Link to ResourceBlocks|

**NOTE:** For ODIM implementation we only have ResourceBlock of type “ComputerSystem”. So, each of the ResourceBlock is used to create a new ResourceZone. In other words, ODIM environment has only one ResourceBlock for each ResourceZone.

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow: "GET", "POST"
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Location: /redfish/v1/CompositionService/ResourceZones/385316c6-8e9d-11ec-942f-caea6c362014
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Thu, 16 Dec 2021 00:31:40 GMT
Content-Length: 421


```

>**Sample response body**

```
{
    "@odata.type": "#Zone.v1_6_0.Zone",
    "Id": "385316c6-8e9d-11ec-942f-caea6c362014",
    "Name": "ResourceZone1",
    "Status": {
        "State": "Enabled",
        "Health": "OK"
    },
    "ZoneType": "ZoneOfResourceBlocks",
    "Links": {
        "ResourceBlocks": [
            {
                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/11260ff0-8b54-11ec-83c3-e6e57d72922b"
            }
        ]
    },
    "@odata.id": "/redfish/v1/CompositionService/ResourceZones/385316c6-8e9d-11ec-942f-caea6c362014"
}

```






## Deleting a ResourceBlock

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/CompositionService/ResourceBlocks/{ResourceBlockId}` |
|**Description** |This operation deletes a specific resource block.|
|**Response code** |`204 NO Content` |
|**Authentication** |Yes|


>**curl command**

```
curl -i DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odim_hosts}:{port}/redfish/v1/CompositionService/ResourceBlocks/{ResourceBlockId}'

```







## Deleting a ResourceZone

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/CompositionService/ResourceZones/{ResourcezoneId}` |
|**Description** |This operation deletes a resource zone instance.|
|**Response code** | `204 No Content` |
|**Authentication** |Yes|


>**curl command**


```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/ResourceZones/{ResourceZoneId}'

```








## Compose Action
Compose Action is used to compose and decompose of a system. “StanzaType” property explains if it is composing or decomposition of a system. 

### Compose System

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/CompositionService/Actions/CompositionService.Compose` |
|**Description** |This operation creates a composed system.|
|**Returns** |JSON schema defined for the redfish compose action.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
    "RequestFormat": "Manifest", 
    "RequestType": "Apply", 
    "Manifest": { 
        "Description": "Specific composition example", 
        "Timestamp": "2021-1216T10:35:16+06:00", 
        "Expand": "None", 
        "Stanzas": [ 
            { 
                "StanzaType": "ComposeSystem", 
                "StanzaId": "The identifier of the stanza.  This is a unique identifier specified by the client and is not used by the service.", 
                "Request": { 
                    "Links": { 
                        "ResourceBlocks": [ 
                            { "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/ComputerSystem-1 " }, 
                            { "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/SASSled-1 " }, 
                            { "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/NVMe-Target-1 " } 
                        ] 
                    } 
                } 
            } 
        ] 
    }
}'
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/Actions/CompositionService.Compose'

```

**NOTE:** Currently we are supporting “Specific” Composition Type only.

>**Sample response header** 

```
HTTP/1.1 200 OK
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Thu, 16 Dec 2021 12:14:32 GMT
Transfer-Encoding: chunked
```

>**Sample response body**

```
{
    "RequestFormat": "Manifest",
    "RequestType": "Apply",
    "Manifest": {
        "Timestamp": "2021-1216T10:35:16+06:00",
        "Description": "Specific composition example",
        "Expand": "None",
        "Stanzas": [
            {
                "StanzaId": "The identifier of the stanza.  This is a unique identifier specified by the client and is not used by the service.",
                "StanzaType": "ComposeSystem",
                "Request": {
                    "Links": {
                        "ResourceBlocks": [
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/ComputerSystem-1 "
                            },
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/SASSled-1 "
                            },
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/NVMe-Target-1 "
                            }
                        ]
                    }
                },
                "Response": {
                    "@odata.id": "/redfish/v1/Systems/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8",
                    "@odata.type": "#ComputerSystem.v1_18_0.ComputerSystem",
                    "Id": "6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8",
                    "Name": "Computer System",
                    "SystemType": "Physical",
                    "Links": {
                        "Chassis": [
                            {
                                "@odata.id": "/redfish/v1/Chassis/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8"
                            }
                        ],
                        "ManagedBy": [
                            {
                                "@odata.id": "/redfish/v1/Managers/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8"
                            }
                        ],
                        "ResourceBlocks": [
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/ComputerSystem-1 "
                            },
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/SASSled-1 "
                            },
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/NVMe-Target-1 "
                            }
                        ]
                    }
                }
            }
        ]
    }
}


```






### Deompose System

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/CompositionService/Actions/CompositionService.Compose` |
|**Description** |This operation Decomposes a system.|
|**Returns** |JSON schema defined for the redfish compose action.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{ 
    "RequestFormat": "Manifest", 
    "RequestType": "Apply", 
    "Manifest": { 
        "Description": "Specific composition example", 
        "Timestamp": "2021-12-16T10:45:10+06:00", 
        "Expand": "None", 
        "Stanzas": [ 
            { 
                "StanzaType": "DecomposeSystem", 
                "StanzaId": "The identifier of the stanza.  This is a unique identifier specified by the client and is not used by the service.", 
                "Request": { 
                    "Links": { 
                        "ComputerSystems": [ 
                            { "@odata.id": "/redfish/v1/Systems/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8" } 
                        ] 
                    } 
                } 
            } 
        ]
    }
}'
 'https://{odimra_host}:{port}/redfish/v1/CompositionService/Actions/CompositionService.Compose'

```



>**Sample response header** 

```
HTTP/1.1 200 OK
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Thu, 16 Dec 2021 12:14:32 GMT
Transfer-Encoding: chunked
```

>**Sample response body**

```
{
    "RequestFormat": "Manifest",
    "RequestType": "Apply",
    "Manifest": {
        "Description": "Specific composition example",
        "Timestamp": "2021-12-16T10:45:10+06:00",
        "Expand": "None",
        "Stanzas": [
            {
                "StanzaId": "The identifier of the stanza.  This is a unique identifier specified by the client and is not used by the service.",
                "StanzaType": "DecomposeSystem",
                "Request": {
                    "Links": {
                        "ComputerSystems": [
                            {
                                "@odata.id": "/redfish/v1/Systems/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8"
                            }
                        ]
                    }
                },
                "Response": {
                   "@odata.id": "/redfish/v1/Systems/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8",
                    "@odata.type": "#ComputerSystem.v1_18_0.ComputerSystem",
                    "Id": "6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8",
                    "Name": "Computer System",
                    "SystemType": "Physical",
                    "Links": {
                        "Chassis": [
                            {
                                "@odata.id": "/redfish/v1/Chassis/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8"
                            }
                        ],
                        "ManagedBy": [
                            {
                                "@odata.id": "/redfish/v1/Managers/6eac063d-60d7-4fe8-a3d8-bf70afb6d235.8"
                            }
                        ],
                        "ResourceBlocks": [
                            {
                                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/ComputerSystem-1 "
                            }
                        ]
                    }
                }
            }
        ]
    }
}
```

**NOTE:** Please refer to following redfish document for action parameters https://www.dmtf.org/sites/default/files/standards/documents/DSP0268_2021.2.pdf

