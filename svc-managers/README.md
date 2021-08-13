#  Managers

Resource Aggregator for ODIM exposes APIs to retrieve information about managers. Examples of managers include:

-   Resource Aggregator for ODIM itself

-   Baseboard Management Controllers \(BMC\)

-   Enclosure Managers

-   Management Controller

-   Other subsystems like plugins


**Supported endpoints**


|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Managers|`GET`| `Login` |
|/redfish/v1/Managers/\{managerId\}|`GET`| `Login` |
|/redfish/v1/Managers/\{managerId\}/EthernetInterfaces|`GET`| `Login` |
|/redfish/v1/Managers/\{managerId\}/HostInterfaces|`GET`| `Login` |
|/redfish/v1/Managers/\{managerId\}/LogServices|`GET`| `Login` |
|/redfish/v1/Managers/\{managerId\}/NetworkProtocol|`GET`| `Login` |


##  Modifying Configurations of managers Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md to change the configurations of an odimra service.
  
**Specific configurations for managers Service are:**
  
##  Log Location of the managers Service
  
/var/log/ODIMRA/managers.log
  
  





##  Collection of managers

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Managers` |
|**Description** |A collection of managers.|
|**Returns** |Links to the manager instances. This collection includes a manager for Resource Aggregator for ODIM itself and other managers.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Managers'

```


>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#ManagerCollection.ManagerCollection",
   "@odata.id":"/redfish/v1/Managers",
   "@odata.type":"#ManagerCollection.ManagerCollection",
   "Description":"Managers view",
   "Name":"Managers",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Managers/a64fc187-e0e9-4f68-82a8-67a616b84b1d"
      },
      {
         "@odata.id":"/redfish/v1/Managers/141cbba9-1e99-4272-b855-1781730bfe1c:1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6"
      },
      {
         "@odata.id":"/redfish/v1/Managers/0e778112-4684-433d-9998-ca6f399c031f:1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b"
      },
      {
         "@odata.id":"/redfish/v1/Managers/a6ddc4c0-2568-4e16-975d-fa771b0be853"
      }
   ],
   "Members@odata.count":6
}


```









##  Single manager

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Managers/{managerId}` |
|**Description** |A single manager.|
|**Returns** |Information about a specific management control system or a plugin or Resource Aggregator for ODIM itself. In the JSON schema representing a system \(BMC\) manager, there are links to the managers for:<ul><li>EthernetInterfaces:<br>`/redfish/v1/Managers/{managerId}/EthernetInterfaces`</li><li>HostInterfaces:<br>`/redfish/v1/Managers/{managerId}/HostInterfaces` </li><li>LogServices:<br>`/redfish/v1/Managers/{managerId}/LogServices` </li><li>NetworkProtocol:<br>`/redfish/v1/Managers/{managerId}/NetworkProtocol` <br> To know more about each manager, perform HTTP `GET` on these links.</li></ul>|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Managers/{managerId}'

```



>**Sample response body for a system \(BMC\) manager** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.etag":"W/\"2D2866FD\"",
   "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1",
   "@odata.type":"#Manager.v1_12_0.Manager",
   "Actions":{ 
      "#Manager.Reset":{ 
         "target":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/Actions/Manager.Reset"
      }
   },
   "CommandShell":{ 
      "ConnectTypesSupported":[ 
         "SSH",
         "Oem"
      ],
      "MaxConcurrentSessions":9,
      "ServiceEnabled":true
   },
   "EthernetInterfaces":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/EthernetInterfaces"
   },
   "FirmwareVersion":"iLO 5 v1.40",
   "GraphicalConsole":{ 
      "ConnectTypesSupported":[ 
         "KVMIP"
      ],
      "MaxConcurrentSessions":10,
      "ServiceEnabled":true
   },
   "HostInterfaces":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/HostInterfaces"
   },
   "Id":"1",
   "Links":{ 
      "ManagerForChassis":[ 
         { 
            "@odata.id":"/redfish/v1/Chassis/88b36c7c-d708-4a4a-8af5-5d58779d377c:1"
         }
      ],
      "ManagerForServers":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/88b36c7c-d708-4a4a-8af5-5d58779d377c:1"
         }
      ],
      "ManagerInChassis":{ 
         "@odata.id":"/redfish/v1/Chassis/88b36c7c-d708-4a4a-8af5-5d58779d377c:1"
      }
   },
   "LogServices":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/LogServices"
   },
   "ManagerType":"BMC",
   "Name":"Manager",
   "NetworkProtocol":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/NetworkProtocol"
   },
   "Oem":{},
      
   "SerialConsole":{ 
      "ConnectTypesSupported":[ 
         "SSH",
         "IPMI",
         "Oem"
      ],
      "MaxConcurrentSessions":13,
      "ServiceEnabled":true
   },
   "Status":{ 
      "State":"Absent"
   },
   "UUID":"a964d6a9-a45c-57aa-90ec-08b38850b2f3",
   "VirtualMedia":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/VirtualMedia"
   }
}
```

>**Sample response body for Resource Aggregator for ODIM manager** 

```
{
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.id":"/redfish/v1/Managers/a64fc187-e0e9-4f68-82a8-67a616b84b1d",
   "@odata.type":"#Manager.v1_12_0.Manager",
   "Name":"ODIMRA",
   "ManagerType":"Service",
   "Id":"a64fc187-e0e9-4f68-82a8-67a616b84b1d",
   "UUID":"a64fc187-e0e9-4f68-82a8-67a616b84b1d",
   "FirmwareVersion":"1.0",
   "Status":{
      "State":"Enabled"
   }
}
```

>**Sample response body for a plugin manager**

```
{
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.etag":"W/\"AA6D42B0\"",
   "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6",
   "@odata.type":"#Manager.v1_12_0.Manager",
   "FirmwareVersion":"v1.0.0",
   "Id":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b",
   "ManagerType":"Service",
   "Name":"GRF",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   },
   "UUID":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b"
}
```