#  Host to fabric networking

Resource Aggregator for ODIM exposes Redfish APIs to view and manage simple fabrics. A fabric is a network topology consisting of entities such as interconnecting switches, zones, endpoints, and address pools. The Redfish `Fabrics` APIs allow you to create and remove these entities in a fabric.

When creating fabric entities, ensure to create them in the following order:

1.  Zone-specific address pools

2.  Address pools for zone of zones

3.  Zone of zones

4.  Endpoints

5.  Zone of endpoints


When deleting fabric entities, ensure to delete them in the following order:

1.  Zone of endpoints

2.  Endpoints

3.  Zone of zones

4.  Address pools for zone of zones

5.  Zone-specific address pools

<blockquote>
IMPORTANT:

Before using the `Fabrics` APIs, ensure that the fabric manager is installed, its plugin is deployed, and added into the Resource Aggregator for ODIM framework.

</blockquote>


**Supported endpoints**



|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Fabrics|GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}|GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}/Switches|GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}|GET|`Login` |
| /redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}/Ports<br> |GET|`Login` |
| /redfish/v1/Fabrics/\{fabricId\} /Switches/\{switchId\}/Ports/\{portid\}<br> |GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}/Zones|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/Zones/\{zoneId\}|GET, PATCH, DELETE|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/AddressPools|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/AddressPools/\{addresspoolid\}|GET, DELETE|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/Endpoints|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/Endpoints/\{endpointId\}|GET, DELETE|`Login`, `ConfigureComponents` |



##  Modifying Configurations of fabric Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md file to change the configurations of an odimra service.
  
**Specific configurations for Fabric Service are:**
  
##  Log Location of the Fabric Service
  
/var/log/ODIMRA/fabric.log
  
  





##  Collection of fabrics

|||
|---------------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Fabrics` |
|**Description** |A collection of simple fabrics.|
|**Returns** |Links to the fabric instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics",
   "Id":"FabricCollection",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6"
      }
   ],
   "Members@odata.count":1,
   "Name":"Fabric Collection",
   "RedfishVersion":"1.14.0",
   "@odata.type":"#FabricCollection.FabricCollection"
}
```












## Single fabric

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}` |
|**Description** |Schema representing a specific fabric.|
|**Returns** |Links to various components contained in this fabric instance - address pools, endpoints, switches, and zones.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}'


```



>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5",
   "@odata.type":"#Fabric.v1_3_0.Fabric",
   "AddressPools":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/AddressPools"
   },
   "Description":"test",
   "Endpoints":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Endpoints"
   },
   "FabricType":"Ethernet",
   "Id":"f4d1578a-d16f-43f2-bb81-cd6db8866db5",
   "Name":"cfm-test",
   "Status":{ 
      "Health":"OK",
      "State":"Enabled"
   },
   "Switches":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches"
   },
   "Zones":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Zones"
   }
}
```







## Collection of switches

|||
|------------------|----------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Switches` |
|**Description** |A collection of switches located in this fabric.|
|**Returns** |Links to the switch instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches",
   "Id":"SwitchCollection",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/fb7dc9fd-d0f1-474e-b849-77262f5d73b7"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/e97a3f0b-cc89-40d8-af3f-9b9bdd793d73"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/a4ca3161-db95-487d-a930-1b13dc697ed0"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/bc95a9aa-8447-4b89-a99d-25235f7bae92"
      }
   ],
   "Members@odata.count":4,
   "Name":"Switch Collection",
   "@odata.type":"#SwitchCollection.SwitchCollection"
}
```











## Single switch

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}` |
|**Description** |JSON schema representing a particular fabric switch.|
|**Returns** |Details of this switch and links to its ports.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/fb7dc9fd-d0f1-474e-b849-77262f5d73b7",
   "@odata.type":"#Switch.v1_8_0.Switch",
   "Id":"fb7dc9fd-d0f1-474e-b849-77262f5d73b7",
   "Manufacturer":"Aruba",
   "Model":"Aruba 8325",
   "Name":"Switch_172.10.20.1",
   "PartNumber":"JL636A",
   "Ports":{ 
      "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/fb7dc9fd-d0f1-474e-b849-77262f5d73b7/Ports"
   },
   "SerialNumber":"TW8BKM302H",
   "Status":{ 
      "Health":"Ok",
      "State":"Enabled"
   },
   "SwitchType":"Ethernet"
}
```




## Collection of ports

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports`` |
|**Description** |A collection of ports of this switch.|
|**Returns** |Links to the port instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports",
   "Id":"PortCollection",
   "Members":[ 
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/0cb2ff96-b7a7-4627-a7b4-274d915f2524",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/54096ea1-cfb8-4a6c-b7a3-d6263db729a6",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/699b8f82-a6bf-47fa-a594-d73c95a8f81e",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/7dd6b8c6-de72-4499-98dc-568a16e28a88",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/9d097004-e034-4772-98c5-fa695688cc4d",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/b82e663b-6d1e-43c9-9d49-63b68b2a5b06"
   ],
   "Members@odata.count":6,
   "Name":"PortCollection",
   "@odata.type":"#PortCollection.PortCollection"
}
	
```








## Single port

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports/{portid}` |
|**Description** |JSON schema representing a specific switch port.|
|**Returns** |Properties of this port.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports/{portid}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/a4ca3161-db95-487d-a930-1b13dc697ed0/Ports/80b5f999-25e9-4b37-992c-de2f065ee0e3",
   "@odata.type":"#Port.v1_5_0.Port",
   "CurrentSpeedGbps":0,
   "Description":"single port",
   "Id":"80b5f999-25e9-4b37-992c-de2f065ee0e3",
   "Links":{ 
      "ConnectedPorts":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/768f9da7-56fc-4f13-b6f8-a1cd241e2313:1/EthernetInterfaces/3"
         }
      ],
      "ConnectedSwitches":[ 

      ]
   },
   "MaxSpeedGbps":25,
   "Name":"1/1/3",
   "PortId":"1/1/3",
   "PortProtocol":"Ethernet",
   "PortType":"UpstreamPort",
   "Status":{ 
      "Health":"Ok",
      "State":"Enabled"
   },
   "Width":1
}
```









## Collection of address pools

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/AddressPools`` |
|**Description** |A collection of address pools.|
|**Returns** |Links to the address pool instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```


>**Sample response body**

```
{
	"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools",
	"Id": "AddressPool Collection",
	"Members": [{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/54a6d41b-6ed2-460b-90c7-cc5fdd74e6ad"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/f936ba02-fa82-456b-a7d7-3d006228f63c"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/95dec77a-c393-4391-8943-29e3ce03c6ca"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/21d4c00e-0c7c-4af3-af76-fd66df5d5831"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/35698e79-a765-4052-86ed-e290d7b6fd01"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/612aeda9-5cca-4f51-b755-b90008467bad"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/470515c8-6089-4d97-ba4f-e7dabc9d7e6a"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/062f6464-4f0e-4a6b-bb6b-c1857bba1533"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/7b740372-2f88-46d8-af84-7b66fee87695"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/7a98eb5d-99e9-4924-b647-057e3ad772bf"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/9f33f532-0796-42fc-819c-7938a4d6de7c"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/3a821e9d-901a-4469-913b-917351d6eef7"
		}
	],
	"Members@odata.count": 12,
	"Name": "AddressPool Collection",
	"RedfishVersion": "1.14.0",
	"@odata.type": "#AddressPoolCollection.AddressPoolCollection"
}
	
```














## Single address pool

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}`` |
|**Description** |JSON schema representing a specific address pool.|
|**Returns** |Properties of this address pool.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}'

```


>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/AddressPools/44c44b52-a784-48e5-9f26-b833d42cf455",
   "@odata.type":"#AddressPool.vxx.AddressPool",
   "BgpEvpn":{ 
      "BgpEvpnEviNumberLowerAddress":200,
      "BgpEvpnEviNumberUpperAddress":220
   },
   "Description":"",
   "ExternalBgp":{ 
      "EbgpAsNumberLowerAddress":65000,
      "EbgpAsNumberUpperAddress":65100
   },
   "IPv4":{ 
      "IPv4FabricLinkLowerAddress":"172.10.20.1",
      "IPv4FabricLinkUpperAddress":"172.10.20.10",
      "IPv4GatewayAddress":"",
      "IPv4HostLowerAddress":"",
      "IPv4HostUpperAddress":"",
      "IPv4LoopbackLowerAddress":"172.10.20.1",
      "IPv4LoopbackUpperAddress":"172.10.20.10",
      "NativeVlan":0,
      "VlanIdentifierLowerAddress":0,
      "VlanIdentifierUpperAddress":0
   },
   "Id":"44c44b52-a784-48e5-9f26-b833d42cf455",
   "Links":{ 
      "Zones":[ 

      ]
   },
   "MultiProtocolIbgp":{ 
      "MPIbgpAsNumberLowerAddress":1,
      "MPIbgpAsNumberUpperAddress":1
   },
   "Name":""
}
	
```






## Collection of endpoints


|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/Endpoints`` |
|**Description** |A collection of fabric endpoints.|
|**Returns** |Links to the endpoint instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints'

```



>**Sample response body**

```
{
	"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints",
 "@odata.type": "#EndpointCollection.EndpointCollection",
	"Members": [{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints/f59d59f3-d2ec-4cc1-9255-f35b5b09a31a"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints/952f0049-d639-4a00-820a-353a95564d37"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints/8a9c27b0-d4d7-4eef-a731-f92aedc49c69"
		}
	],
	"Members@odata.count": 3,
	"Name": "Endpoint Collection"
	
}
	
```





##  Single endpoint

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}` |
|**Description** |JSON schema representing a specific fabric endpoint.|
|**Returns** |Details of this endpoint.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/8f6d8828-a21a-464f-abf9-ed062fa08cd9/Endpoints/b21f3e57-e46d-4a8e-92c8-8658edd107cb",
   "@odata.type":"#Endpoint.v1_11_0.Endpoint",
   "Description":"NK Endpoint Collection Description",
   "EndpointProtocol":"Ethernet",
   "Id":"b21f3e57-e46d-4a8e-92c8-8658edd107cb",
   "Links":{ 
      "ConnectedPorts":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/8f6d8828-a21a-464f-abf9-ed062fa08cd9/Switches/3f4ac957-90ec-4676-91b9-90f9d78ef56c/Ports/7f708d4f-795d-401d-8bc1-c797fb3ce20b"
         }
      ],
      "Zones":[ 

      ]
   },
   "Name":"NK Endpoint Collection1",
   "Redundancy":[ 
      { 
         "MaxNumSupported":2,
         "MemberId":"Bond0",
         "MinNumNeeded":1,
         "Mode":"",
         "RedundancySet":[ 
            [ 

            ]
         ],
         "Status":{ 
            "Health":"",
            "State":""
         }
      }
   ]
}
```






## Collection of zones

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/Zones`` |
|**Description** |A collection of fabric zones.|
|**Returns** |Links to the zone instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```


>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/afa007b7-7ea6-4ab3-b5f1-ad37c8aebed7"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/85462c4f-028d-45d6-99d8-73c7889ea263"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/e906a6ab-18ef-4617-a151-420265e7d0f9"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/f310bf40-5163-4cbf-be5b-ac574fe87863"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/5c0b60a0-55f7-43f0-9b23-bfbba9130743"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/ac424042-7524-4c04-acbd-2d1af0a4832f"
      }
   ],
   "Members@odata.count":6,
   "Name":"Zone Collection",
   "@odata.type":"#ZoneCollection.ZoneCollection"
	
```






## Single zone


|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}` |
|**Description** |JSON schema representing a specific fabric zone.|
|**Returns** |Details of this zone.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/f310bf40-5163-4cbf-be5b-ac574fe87863",
   "@odata.type":"#Zone.v1_6_1.Zone",
   "DefaultRoutingEnabled":false,
   "Description":"",
   "Id":"f310bf40-5163-4cbf-be5b-ac574fe87863",
   "Links":{ 
      "AddressPools":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/AddressPools/e8edcc87-81f9-43a9-b1ce-20a895a60014"
         }
      ],
      "ContainedByZones":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/ac424042-7524-4c04-acbd-2d1af0a4832f"
         }
      ],
      "ContainsZones":[ 

      ],
      "Endpoints":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Endpoints/a9a41a01-ac4c-460d-923e-98ad9cc7abef"
         }
      ],
      "InvolvedSwitches":[ 

      ]
   },
   "Name":"NK Zone 1",
   "ZoneType":"ZoneOfEndpoints"
}
```









## Creating a zone-specific address pool

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/AddressPools` |
|**Description** |This operation creates an address pool that can be used by a zone of endpoints.|
|**Returns** |<ul><li>Link to the created address pool in the `Location` header.</li><li>JSON schema representing the created address pool.</li></ul>|
|**Response code** |`201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"FC 18 vlan_102 - AddressPools",
   "Description":"vlan_102",
   "IPv4":{
      "VlanIdentifierAddressRange":{
         "Lower":102,
         "Upper":102
      }
   },
   "BgpEvpn":{
      "GatewayIPAddressList":[
         "10.18.102.2/24",
         "10.18.102.3/24"
      ],
      "AnycastGatewayIPAddress":"10.18.102.1"
   }
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```







>**Sample request body**

```
{
   "Name":"FC 18 vlan_102 - AddressPools",
   "Description":"vlan_102",
   "IPv4":{
      "VlanIdentifierAddressRange":{
         "Lower":102,
         "Upper":102
      }
   },
   "BgpEvpn":{
      "GatewayIPAddressList":[
         "10.18.102.2/24",
         "10.18.102.3/24"
      ],
      "AnycastGatewayIPAddress":"10.18.102.1"
   }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(optional\)<br> |Name for the address pool.|
|Description|String \(optional\)<br> |Description for the address pool.|
|IPv4\{| \(required\)<br> | |
|VlanIdentifierAddressRange\{| \(optional\)<br> | A single VLAN to assign on the ports or lags.<br> |
|Lower|Integer \(required\)<br> |VLAN lower address.|
|Upper\}\}|Integer \(required\)<br> |VLAN upper address.<br>**NOTE:**<br> `Lower` and `Upper` must have the same value. Ensure that IP range is accurate and it does not overlap with other pools.|
|BgpEvpn\{| \(required\)<br> | |
|GatewayIPAddressList|Array \(required\)<br> | IP pool to assign IPv4 address to the IP interface for VLAN per switch.<br> |
|AnycastGatewayIPAddress\}|String \(required\)<br> | A single active gateway IP address for the IP interface.<br> |
| | | |



>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:18:54 GMT
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
  "@odata.id": "/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf",
  "@odata.type": "#AddressPool.vxx.AddressPool",
  "BgpEvpn": {
    "AnycastGatewayIPAddress": "10.18.102.1",
    "AnycastGatewayMACAddress": "",
    "GatewayIPAddressList": [
      "10.18.102.2/24",
      "10.18.102.3/24"
    ],
    "RouteDistinguisherList": "",
    "RouteTargetList": [
      
    ]
  },
  "Description": "vlan_102",
  "Ebgp": {
    
  },
  "IPv4": {
    "EbgpAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "FabricLinkAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "IbgpAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "LoopbackAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "NativeVlan": 0,
    "VlanIdentifierAddressRange": {
      "Lower": 102,
      "Upper": 102
    }
  },
  "Id": "e2ec196d-4b55-44b3-b928-8273de9fb8bf",
  "Links": {
    "Zones": [
      
    ]
  },
  "Name": "FC 18 vlan_102 - AddressPools"
}

```






## Creating an address pool for zone of zones

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/AddressPools` |
|**Description** |This operation creates an address pool for a zone of zones in a specific fabric.|
|**Returns** |<ul><li>Link to the created address pool in the `Location` header.</li><li>JSON schema representing the created address pool.</li></ul>|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
  "Name": "AddressPool for ZoneOfZones - Vlan3002",
  "IPv4": {  
    "VlanIdentifierAddressRange": {
        "Lower": 3002,
        "Upper": 3002
    },
    "IbgpAddressRange": {
              "Lower": "192.12.1.10",
              "Upper": "192.12.1.15"
    },
    "EbgpAddressRange": {
              "Lower": "172.12.1.10",
              "Upper": "172.12.1.15"
    }
  },
  "Ebgp": {
    "AsNumberRange": {
              "Lower": 65120,
              "Upper": 65125
    }
  },
  "BgpEvpn": {
    "RouteDistinguisherList": ["65002:102"],  
    "RouteTargetList": ["65002:102", "65002:102"],
    "GatewayIPAddressList": ["192.168.18.122/31", "192.168.18.123/31"]
  }
}'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```







>**Sample request body**

```
{
  "Name": "AddressPool for ZoneOfZones - Vlan3002",
  "IPv4": {  
    "VlanIdentifierAddressRange": {
        "Lower": 3002,
        "Upper": 3002
    },
    "IbgpAddressRange": {
              "Lower": "192.12.1.10",
              "Upper": "192.12.1.15"
    },
    "EbgpAddressRange": {
              "Lower": "172.12.1.10",
              "Upper": "172.12.1.15"
    }
  },
  "Ebgp": {
    "AsNumberRange": {
              "Lower": 65120,
              "Upper": 65125
    }
  },
  "BgpEvpn": {
    "RouteDistinguisherList": ["65002:102"],  
    "RouteTargetList": ["65002:102", "65002:102"],
    "GatewayIPAddressList": ["192.168.18.122/31", "192.168.18.123/31"]
  }
}

```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String|Name for the address pool.|
|Description|String \(optional\)<br> |Description for the address pool.|
|IPv4\{| \(required\)<br> | |
|VlanIdentifierAddressRange\{| \(required\)<br> | A single VLAN \(virtual LAN\) used for creating the IP interface for the user Virtual Routing and Forwarding \(VRF\).<br> |
|Lower|Integer \(required\)<br> |VLAN lower address|
|Upper\}|Integer \(required\)<br> |VLAN upper address|
|IbgpAddressRange\{| \(required\)<br> | IPv4 address used as the Router Id for the VRF per switch.<br> |
|Lower|String \(required\)<br> |IPv4 lower address|
|Upper\}|String \(required\)<br> |IPv4 upper address|
|EbgpAddressRange\{| \(optional\)<br> |External neighbor IPv4 addresses.|
|Lower|String \(required\)<br> |IPv4 lower address|
|Upper\} \}|String \(required\)<br> |IPv4 upper address|
|Ebgp\{| \(optional\)<br> | |
|AsNumberRange\{| \(optional\)<br> |External neighbor ASN.<br>**NOTE:**<br> `EbgpAddressRange` and `AsNumberRange` values should be a matching sequence and should be of same length.|
|Lower|Integer \(optional\)<br> | |
|Upper\} \}|Integer \(optional\)<br> | |
|BgpEvpn\{| \(required\)<br> | |
|RouteDistinguisherList|Array \(required\)<br> | Single route distinguisher value for the VRF.<br> |
|RouteTargetList|Array \(optional\)<br> | Route targets. By default, the route targets will be configured as both import and export.<br> |
|GatewayIPAddressList\}|Array \(required\)<br> | IP pool to assign IPv4 address to the IP interface used by the VRF per switch.<br> |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:18:58 GMT
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "@odata.type":"#AddressPool.vxx.AddressPool",
   "BgpEvpn":{
      "AnycastGatewayIPAddress":"",
      "AnycastGatewayMACAddress":"",
      "GatewayIPAddressList":[
         "192.168.18.122/31",
         "192.168.18.123/31"
      ],
      "RouteDistinguisherList":[
         "65002:102"
      ],
      "RouteTargetList":[
         "65002:102",
         "65002:102"
      ]
   },
   "Description":"",
   "Ebgp":{
      "AsNumberRange":{
         "Lower":65120,
         "Upper":65125
      }
   },
   "IPv4":{
      "EbgpAddressRange":{
         "Lower":"172.12.1.10",
         "Upper":"172.12.1.15"
      },
      "FabricLinkAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "IbgpAddressRange":{
         "Lower":"192.12.1.10",
         "Upper":"192.12.1.15"
      },
      "LoopbackAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "NativeVlan":0,
      "VlanIdentifierAddressRange":{
         "Lower":3002,
         "Upper":3002
      }
   },
   "Id":"84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "Links":{
      "Zones":[

      ]
   },
   "Name":"AddressPool for ZoneOfZones - Vlan3002"
}
```


##  Adding a zone of zones



|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones` |
|**Description** |This operation creates an empty container zone for all the other zones in a specific fabric. To assign address pools, endpoints, other zones, or switches to this zone, perform HTTP `PATCH` on this zone. See [Updating a Zone](#updating-a-zone).|
|**Returns** |JSON schema representing the created zone.|
|**Response code** |`201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"Fabric Zone of Zones:",
   "ZoneType":"ZoneOfZones",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
         }
      ]
   }
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```

>**Sample request body**

```
{
   "Name":"Fabric Zone of Zones:",
   "ZoneType":"ZoneOfZones",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
         }
      ]
   }
}
```

**Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String \(optional\)<br> |Name for the zone.|
|Description|String \(optional\)<br> |Description for the zone.|
|ZoneType|String|The type of the zone to be created. Options include: `ZoneofZones` and `ZoneofEndpoints`<br> The type of the zone for a zone of zones is `ZoneofZones`.<br> |
|Links\{| \(optional\)<br> | |
|AddressPools|Array \(optional\)<br> | `AddressPool` links supported for the Zone of Zones \(`AddressPool` links created for `ZoneofZones`\).<br> |


>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:19:00 GMT
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98",
   "@odata.type":"#Zone.v1_6_1.Zone",
   "Description":"",
   "Id":"a2dc8760-ea05-4cab-8f95-866c1c380f98",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
         }
      ],
      "ContainedByZones":[

      ],
      "ContainsZones":[

      ],
      "Endpoints":[

      ],
      "InvolvedSwitches":[

      ]
   },
   "Name":"Fabric Zone of Zones:",
   "ZoneType":"ZoneOfZones"
}
```


## Adding an endpoint


|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Endpoints` |
|**Description** |This operation creates an endpoint in a specific fabric.|
|**Returns** | <ul><li>Link to the created endpoint in the `Location` header.</li><li>JSON schema representing the created endpoint.</li></ul>|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"Host 2 Endpoint 1 Collection",
   "Description":"Host 2 Endpoint 1 Collection Description",
   "Redundancy":[
      {
         "Mode":"Sharing",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/81f5ed9a-a4a1-4383-a450-a7f98b792ca2/Ports/29f077b0-e7a5-495f-a3d2-643937f600de"
            },
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/b7ba4ece-1716-4d8d-af2c-7aea1682bf91/Ports/62a32f83-c7b1-4cb7-9b47-2f444588d29b"
            }
         ]
      }
   ]
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints'

```

>**Sample request body** for a single endpoint

```
{
   "Name":"NK Endpoint Collection",
   "Description":"NK Endpoint Collection Description",
   "Links":{
      "ConnectedPorts":[
         {
            "@odata.id":"/redfish/v1/Fabrics/113a30e3-f312-4221-8f7f-49943c5ff07d/Switches/f4a37f55-be1e-400b-93be-7d7c0afd4cbd/Ports/0d22b201-30d5-43e8-90ab-277c87624c05"
         }
      ]
   }
}
```

>**Sample request body** for a redundant endpoint

```
{
   "Name":"Host 2 Endpoint 1 Collection",
   "Description":"Host 2 Endpoint 1 Collection Description",
   "Redundancy":[
      {
         "Mode":"Sharing",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/81f5ed9a-a4a1-4383-a450-a7f98b792ca2/Ports/29f077b0-e7a5-495f-a3d2-643937f600de"
            },
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/b7ba4ece-1716-4d8d-af2c-7aea1682bf91/Ports/62a32f83-c7b1-4cb7-9b47-2f444588d29b"
            }
         ]
      }
   ]
}
```

**Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String \(optional\)<br> |Name for the endpoint.|
|Description|String \(optional\)<br> |Description for the endpoint.|
|Links\{| \(required\)<br> | |
|ConnectedPorts|Array \(required\)<br> | Switch port connected to the switch.<br>  <br> |
|Zones\}|Array \(optional\)<br> | Endpoint is part of `ZoneofEndpoints` and only one zone is permitted in the zones list.<br> |
|Redundancy\[|Array| |
|Mode|String|Redundancy mode.|
|RedundancySet\]|Array| Set of redundancy ports connected to the switches.<br> |

>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:19:02 GMT
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97",
   "@odata.type":"#Endpoint.v1_11_0.Endpoint",
   "Description":"Host 2 Endpoint 1 Collection Description",
   "EndpointProtocol":"Ethernet",
   "Id":"fe34aff2-e81f-4167-a0c3-9bf5a67e2a97",
   "Links":{
      "ConnectedPorts":[

      ],
      "Zones":[

      ]
   },
   "Name":"Host 2 Endpoint 1 Collection",
   "Redundancy":[
      {
         "MaxNumSupported":2,
         "MemberId":"Bond0",
         "MinNumNeeded":1,
         "Mode":"Sharing",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/81f5ed9a-a4a1-4383-a450-a7f98b792ca2/Ports/29f077b0-e7a5-495f-a3d2-643937f600de"
            },
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/b7ba4ece-1716-4d8d-af2c-7aea1682bf91/Ports/62a32f83-c7b1-4cb7-9b47-2f444588d29b"
            }
         ],
         "Status":{
            "Health":"",
            "State":""
         }
      }
   ]
}
```






## Creating a zone of endpoints

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones` |
|**Description** |This operation creates a zone of endpoints in a specific fabric.<br>**NOTE:**<br> Ensure that the endpoints are created first before assigning them to the zone of endpoints.|
|**Returns** |JSON schema representing the created zone.|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"Fabric ZoneofEndpoint",
   "ZoneType":"ZoneOfEndpoints",
   "Links":{
      "ContainedByZones":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98"
         }
      ],
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf"
         }
      ],
      "Endpoints":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97"
         }
      ]
   }
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```






>**Sample request body**

```
{
   "Name":"Fabric ZoneofEndpoint",
   "ZoneType":"ZoneOfEndpoints",
   "Links":{
      "ContainedByZones":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98"
         }
      ],
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf"
         }
      ],
      "Endpoints":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97"
         }
      ]
   }
}
```


**Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String \(optional\)<br> |The name for the zone.|
|Description|String \(optional\)<br> |The description for the zone.|
|DefaultRoutingEnabled|Boolean \(required\)<br> |Set to `false`.|
|ZoneType|String \(required\)<br> |The type of the zone to be created. Options include: `ZoneofZones`and `ZoneofEndpoints`<br> The type of the zone for a zone of endpoints is `ZoneofEndpoints`.<br> |
|Links\{|Object \(optional\)<br> |Contains references to other resources that are related to the zone.|
|ContainedByZones \[\{|Array \(optional\)<br> |Represents an array of `ZoneofZones` for the zone being created \(applicable when creating a zone of endpoints\).|
|@odata.id \}\]|String|Link to a Zone of zones.|
|AddressPools \[\{|Array \(optional\)<br> |Represents an array of address pools linked with the zone \(zone-specific address pools\).|
|@odata.id \}\]|String|Link to an address pool.|
|Endpoints \[\{|Array \(optional\)<br> |Represents an array of endpoints to be included in the zone.|
|@odata.id \}\]|String|Link to an endpoint.|



>**Sample response header**

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control: no-cache
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Location: /redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/06d344bb-cce1-4b0c-8414-6f6df1ea373f
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Thu, 14 May 2020 16:19:37 GMT
Transfer-Encoding: chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/06d344bb-cce1-4b0c-8414-6f6df1ea373f",
   "@odata.type":"#Zone.v1_6_1.Zone",
   "Description":"",
   "Id":"06d344bb-cce1-4b0c-8414-6f6df1ea373f",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf"
         }
      ],
      "ContainedByZones":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98"
         }
      ],
      "ContainsZones":[

      ],
      "Endpoints":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97"
         }
      ],
      "InvolvedSwitches":[

      ]
   },
   "Name":"Fabric ZoneofEndpoint",
   "ZoneType":"ZoneOfEndpoints"
}
```



##  Updating a zone

|||
|---------------|---------------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}` |
|**Description** |This operation assigns or unassigns a collection of endpoints, address pools, zone of zones, or switches to a zone of endpoints or a collection of address pools to a zone of zones in a specific fabric.|
|**Returns** |JSON schema representing an updated zone.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
	"Links": {
		"Endpoints": [{
			"@odata.id": "/redfish/v1/Fabrics/d76f4c66-aa60-4693-bea1-feac44fb9f81/Endpoints/a9d7f926-fc3c-465f-9724-928ba9becdb2"
		}]
	}
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}'

```








>**Sample request body** \(assigning links for a zone of endpoints\)

```
{
	"Links": {
		"Endpoints": [{
			"@odata.id": "/redfish/v1/Fabrics/d76f4c66-aa60-4693-bea1-feac44fb9f81/Endpoints/a9d7f926-fc3c-465f-9724-928ba9becdb2"
		}]
	}
}
```

>**Sample request body** \(unassigning links for a zone of endpoints\)

```
{
	"Links": {
		"Endpoints": []
	}
}
```

>**Sample request body** \(assigning links for a zone of zone\)

```
{
   "Links":{
      "AddressPools":[
        "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
      ]
   }
}
```

>**Sample request body** \(unassigning links for a zone of zone\)

```
{
   "Links":{
      "AddressPools":[

      ]
   }
}
```

>**Sample response body** \(assigned links in a zone of endpoints\)

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/143476dc-0ac1-4352-96f3-e0782aeed84a/Zones/57c325f0-eda4-4754-b8da-826d5e266c04",
   "@odata.type":"#Zone.v1_6_1.Zone",
   "Description":"NK Zone Collection Description",
   "Id":"57c325f0-eda4-4754-b8da-826d5e266c04",
   "Links":{ 
      "AddressPools":[ 

      ],
      "ContainedByZones":[ 

      ],
      "ContainsZones":[ 

      ],
      "Endpoints":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Endpoints/e8edcc87-81f9-43a9-b1ce-20a895a60014"
         }
      ],
      "InvolvedSwitches":[ 

      ]
   },
   "Name":"SS Zone Collection default",
   "ZoneType":"ZoneOfEndpoints"
}
```


>**Sample response body** \(unassigned links in a zone of endpoints\)

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/143476dc-0ac1-4352-96f3-e0782aeed84a/Zones/57c325f0-eda4-4754-b8da-826d5e266c04",
   "@odata.type":"#Zone.v1_6_1.Zone",
   "Description":"NK Zone Collection Description",
   "Id":"57c325f0-eda4-4754-b8da-826d5e266c04",
   "Links":{ 
      "AddressPools":[ 

      ],
      "ContainedByZones":[ 

      ],
      "ContainsZones":[ 

      ],
      "Endpoints":[ 

      ],
      "InvolvedSwitches":[ 

      ]
   },
   "Name":"SS Zone Collection default",
   "ZoneType":"ZoneOfEndpoints"
}
```


## Deleting a zone

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}` |
|**Description** |This operation deletes a zone in a specific fabric.<br>**NOTE:**<br> If you delete a non-empty zone \(a zone which contains links to address pools, other zones, endpoints, or switches\), you will receive an HTTP `400` error. Before attempting to delete, unassign all links in the zone. See [updating a zone](#updating-a-zone).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}'

```







## Deleting an endpoint

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}` |
|**Description** |This operation deletes an endpoint in a specific fabric.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odim_hosts}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}'

```






## Deleting an address pool

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}` |
|**Description** |This operation deletes an address pool in a specific fabric.<br>**NOTE:**<br> If you delete an address pool that is being used in any zone, you will receive an HTTP `400` error. Before attempting to delete, ensure that the address pool you want to delete is not present in any zone. To get the list of address pools in a zone, see links to `addresspools` in the sample response body for a [single zone](#single-zone).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}'

```