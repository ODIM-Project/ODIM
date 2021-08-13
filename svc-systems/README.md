#  Resource inventory

Resource Aggregator for ODIM allows you to view the inventory of compute and local storage resources through Redfish `Systems`, `Chassis`, and `Managers` endpoints. 
It also offers the capability to:	
- Search inventory information based on one or more configuration parameters.
	
- Manage the resources added in the inventory. 

To discover crucial configuration information about a resource, including chassis, perform `GET` on these endpoints.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Systems|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}|GET, PATCH|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory/\{memoryId\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/MemoryDomains|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/NetworkInterfaces|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces/\{Id\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Bios|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/SecureBoot|GET|`Login` |
|/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}|`GET`|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Drives/\{driveId\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes/\{volumeId\}|GET, DELETE|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors/\{Id\}|GET|`Login` |
|/redfish/v1/Systems?$filter=\{searchKeys\}%20\{conditionKeys\}%20\{value\}|GET|`Login` |
| /redfish/v1/Systems/\{ComputerSystemId\}/Bios/Settings<br> |GET, PATCH|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.Reset|POST|`ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.SetDefaultBootOrder|POST|`ConfigureComponents` |

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Chassis|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Chassis/\{chassisId\}|GET, PATCH, DELETE|`Login`, `ConfigureComponents`|
|/redfish/v1/Chassis/\{chassisId\}/Thermal|GET|`Login`|
|/redfish/v1/Chassis/\{chassisId\}/NetworkAdapters|GET|`Login` |
|/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{networkadapterId}|GET|`Login`|




>**NOTE:**
To view system, chassis, and manager resources, ensure that you have a minimum privilege of `Login`. If you do not have the necessary privileges, you will receive an HTTP `403 Forbidden` error.
  
##  Modifying Configurations of Systems Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md file to change the configurations of an odimra service.
  
**Specific configurations for Systems Service are:**
  
##  Log Location of the Systems Service
  
/var/log/ODIMRA/system.log
    
  









##  Collection of computer systems

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems` |
|**Description** |This operation lists all systems available with Resource Aggregator for ODIM.|
|**Returns** |A collection of links to computer system instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems'


```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
   "@odata.id":"/redfish/v1/Systems/",
   "@odata.type":"#ComputerSystemCollection.ComputerSystemCollection",
   "Description":"Computer Systems view",
   "Name":"Computer Systems",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Systems/ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1"
      },
      { 
         "@odata.id":"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":2
}
```







## Single computer system

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}` |
|**Description** |This endpoint fetches information about a specific system.|
|**Returns** |JSON schema representing this computer system instance.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}'

```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
   "@odata.etag":"W/\"8C36EBD2\"",
   "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1",
   "@odata.type":"#ComputerSystem.v1_14_1.ComputerSystem",
   "Id":"e24fb205-6669-4080-b53c-67d4923aa73e:1",
   "Actions":{ 
      "#ComputerSystem.Reset":{ 
         "ResetType@Redfish.AllowableValues":[ 
            "On",
            "ForceOff",
            "GracefulShutdown",
            "ForceRestart",
            "Nmi",
            "PushPowerButton"
         ],
         "target":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Actions/ComputerSystem.Reset"
      }
   },
   "AssetTag":"",
   "Bios":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Bios"
   },
   "BiosVersion":"U32 v2.00 (02/02/2019)",
   "Boot":{ 
      "BootOptions":{ 
         "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/BootOptions"
      },
      "BootOrder":[ 
         "Boot000A",
         "Boot000B",
         "Boot000C",
         "Boot0012",
         "Boot000D",
         "Boot000F",
         "Boot000E",
         "Boot0010",
         "Boot0011",
         "Boot0013",
         "Boot0015",
         "Boot0014",
         "Boot0016"
      ],
      "BootSourceOverrideEnabled":"Disabled",
      "BootSourceOverrideMode":"UEFI",
      "BootSourceOverrideTarget":"None",
      "BootSourceOverrideTarget@Redfish.AllowableValues":[ 
         "None",
         "Cd",
         "Hdd",
         "Usb",
         "SDCard",
         "Utilities",
         "Diags",
         "BiosSetup",
         "Pxe",
         "UefiShell",
         "UefiHttp",
         "UefiTarget"
      ],
      "UefiTargetBootSourceOverride":"None",
      "UefiTargetBootSourceOverride@Redfish.AllowableValues":[ 
         "UsbClass(0xFFFF,0xFFFF,0xFF,0xFF,0xFF)",
         "PciRoot(0x0)/Pci(0x14,0x0)/USB(0x13,0x0)",
         "PciRoot(0x3)/Pci(0x0,0x0)/Pci(0x0,0x0)/Scsi(0x0,0x4000)",
         "PciRoot(0x3)/Pci(0x0,0x0)/Pci(0x0,0x0)/Scsi(0x1,0x4000)",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv4(0.0.0.0)/Uri()",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv4(0.0.0.0)",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
         "HD(2,GPT,E0698C18-D9A0-4F58-93CA-A6AEA6BFC93B,0x96800,0x32000)/\\EFI\\Microsoft\\Boot\\bootmgfw.efi",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv4(0.0.0.0)/Uri()",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv4(0.0.0.0)",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)"
      ]
   },
   "EthernetInterfaces":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/EthernetInterfaces"
   },
   "HostName":"",
   "IndicatorLED":"Off",
   "Links":{ 
      "ManagedBy":[ 
         { 
            "@odata.id":"/redfish/v1/Managers/e24fb205-6669-4080-b53c-67d4923aa73e:1"
         }
      ],
      "Chassis":[ 
         { 
            "@odata.id":"/redfish/v1/Chassis/e24fb205-6669-4080-b53c-67d4923aa73e:1"
         }
      ]
   },
   "LogServices":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/LogServices"
   },
   "Manufacturer":"HPE",
   "Memory":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Memory"
   },
   "MemoryDomains":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/MemoryDomains"
   },
   "MemorySummary":{ 
      "Status":{ 
         "HealthRollup":"OK"
      },
      "TotalSystemMemoryGiB":384,
      "TotalSystemPersistentMemoryGiB":0
   },
   "Model":"ProLiant DL360 Gen10",
   "Name":"Computer System",
   "NetworkInterfaces":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/NetworkInterfaces"
   },
   "Oem":{ 
      
         },
         "AggregateHealthStatus":{ 
            "AgentlessManagementService":"Unavailable",
            "BiosOrHardwareHealth":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "FanRedundancy":"Redundant",
            "Fans":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Memory":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Network":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "PowerSupplies":{ 
               "PowerSuppliesMismatch":false,
               "Status":{ 
                  "Health":"OK"
               }
            },
            "PowerSupplyRedundancy":"Redundant",
            "Processors":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "SmartStorageBattery":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Storage":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Temperatures":{ 
               "Status":{ 
                  "Health":"OK"
               }
            }
         },
         "Bios":{ 
            "Backup":{ 
               "Date":"10/02/2018",
               "Family":"U32",
               "VersionString":"U32 v1.46 (10/02/2018)"
            },
            "Current":{ 
               "Date":"02/02/2019",
               "Family":"U32",
               "VersionString":"U32 v2.00 (02/02/2019)"
            },
            "UefiClass":2
         },
         "CurrentPowerOnTimeSeconds":38039,
         "DeviceDiscoveryComplete":{ 
            "AMSDeviceDiscovery":"NoAMS",
            "DeviceDiscovery":"vMainDeviceDiscoveryComplete",
            "SmartArrayDiscovery":"Complete"
         },
         "ElapsedEraseTimeInMinutes":0,
         "EndOfPostDelaySeconds":null,
         "EstimatedEraseTimeInMinutes":0,
         "IntelligentProvisioningAlwaysOn":true,
         "IntelligentProvisioningIndex":9,
         "IntelligentProvisioningLocation":"System Board",
         "IntelligentProvisioningVersion":"3.20.154",
         "IsColdBooting":false,
         "Links":{ 
            "PCIDevices":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIDevices"
            },
            "PCISlots":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCISlots"
            },
            "NetworkAdapters":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/BaseNetworkAdapters"
            },
            "SmartStorage":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/SmartStorage"
            },
            "USBPorts":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/USBPorts"
            },
            "USBDevices":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/USBDevices"
            },
            "EthernetInterfaces":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/EthernetInterfaces"
            },
            "WorkloadPerformanceAdvisor":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/WorkloadPerformanceAdvisor"
            }
         },
         "PCAPartNumber":"847479-001",
         "PCASerialNumber":"PVZEK0ARHBV392",
         "PostDiscoveryCompleteTimeStamp":"2020-02-23T23:09:45Z",
         "PostDiscoveryMode":null,
         "PostMode":null,
         "PostState":"InPostDiscoveryComplete",
         "PowerAllocationLimit":1000,
         "PowerAutoOn":"Restore",
         "PowerOnDelay":"Minimum",
         "PowerOnMinutes":463,
         "PowerRegulatorMode":"Dynamic",
         "PowerRegulatorModesSupported":[ 
            "OSControl",
            "Dynamic",
            "Max",
            "Min"
         ],
         "SMBIOS":{ 
            "extref":"/smbios"
         },
         "ServerFQDN":"",
         "SmartStorageConfig":[ 
            { 
               "@odata.id":"/redfish/v1/systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/smartstorageconfig"
            }
         ],
         "SystemROMAndiLOEraseComponentStatus":{ 
            "BIOSSettingsEraseStatus":"Idle",
            "iLOSettingsEraseStatus":"Idle"
         },
         "SystemROMAndiLOEraseStatus":"Idle",
         "SystemUsage":{ 
            "AvgCPU0Freq":126,
            "AvgCPU1Freq":0,
            "CPU0Power":62,
            "CPU1Power":54,
            "CPUICUtil":0,
            "CPUUtil":2,
            "IOBusUtil":0,
            "JitterCount":0,
            "MemoryBusUtil":0
         },
         "UserDataEraseComponentStatus":{ 

         },
         "UserDataEraseStatus":"Idle",
         "VirtualProfile":"Inactive"
      }
   },
   "PCIeDevices":[
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/1"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/2"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/3"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/4"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/5"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/6"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/7"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/8"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/9"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/10"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/11"
    }
   ],
   "PCIeDevices@odata.count": 11,
   "PowerState":"On",
   "ProcessorSummary":{ 
      "Count":2,
      "Model":"Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz",
      "Status":{ 
         "HealthRollup":"OK"
      }
   },
   "Processors":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Processors"
   },
   "SKU":"867959-B21",
   "SecureBoot":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/SecureBoot"
   },
   "SerialNumber":"MXQ91100T9",
   "Status":{ 
      "Health":"OK",
      "HealthRollup":"OK",
      "State":"Starting"
   },
   "Storage":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Storage"
   },
   "SystemType":"Physical",
   "TrustedModules":[ 
      { 
         "Oem":{ 
            
         },
         "Status":{ 
            "State":"Absent"
         }
      }
   ],
   "UUID":"39373638-3935-584D-5139-313130305439"
```





 



##  Memory collection

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Memory` |
|**Description** |This operation lists all memory devices of a specific server.|
|**Returns** |List of memory resource endpoints.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory'


```





## Single memory


|||
|---------|-------|
|**Method** |GET|
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}` |
|**Description** |This endpoint retrieves configuration information of specific memory.|
|**Returns** |JSON schema representing this memory resource.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}'


```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Memory.Memory",
   "@odata.etag":"W/\"E6EC3A2C\"",
   "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Memory/proc1dimm1",
   "@odata.type":"#Memory.v1_7_1.Memory",
   "Id":"proc1dimm1",
   "BaseModuleType":"RDIMM",
   "BusWidthBits":72,
   "CacheSizeMiB":0,
   "CapacityMiB":32768,
   "DataWidthBits":64,
   "DeviceLocator":"PROC 1 DIMM 1",
   "ErrorCorrection":"MultiBitECC",
   "LogicalSizeMiB":0,
   "Manufacturer":"HPE",
   "MemoryDeviceType":"DDR4",
   "MemoryLocation":{ 
      "Channel":6,
      "MemoryController":2,
      "Slot":1,
      "Socket":1
   },
   "MemoryMedia":[ 
      "DRAM"
   ],
   "MemoryType":"DRAM",
   "Name":"proc1dimm1",
   "NonVolatileSizeMiB":0,
   "Oem":{ 
      
   },
   "OperatingMemoryModes":[ 
      "Volatile"
   ],
   "OperatingSpeedMhz":2666,
   "PartNumber":"840758-091",
   "PersistentRegionSizeLimitMiB":0,
   "RankCount":2,
   "SecurityCapabilities":{ 

   },
   "Status":{ 
      "Health":"OK",
      "State":"Enabled"
   },
   "VendorID":"52736",
   "VolatileRegionSizeLimitMiB":32768,
   "VolatileSizeMiB":32768
}
```





 




##  Memory domains

|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains` |
|**Description** | This endpoint lists memory domains of a specific system.<br> Memory Domains indicate to the client which Memory \(DIMMs\) can be grouped in Memory Chunks to form interleave sets or otherwise grouped.<br> |
|**Returns** |List of memory domain endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains'


```




##  BIOS

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Bios` |
|**Description** | Use this endpoint to discover system-specific information about a BIOS resource and actions for changing to BIOS settings.<br>**NOTE:**<br> Changes to the BIOS typically require a system reset before they take effect.|
|**Returns** |<ul><li>Actions for changing password and resetting BIOS.</li><li>BIOS attributes.</li></ul> |
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Bios'


```





## Network interfaces

|||
|--------|---------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces` |
|**Description** | This endpoint lists network interfaces of a specific system.<br> A network interface contains links to network adapter, network port, and network device function resources.<br> |
|**Returns** |List of network interface endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces'


```






##  Ethernet interfaces

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces` |
|**Description** |This endpoint lists Ethernet interfaces or network interface controller \(NIC\) of a specific system.|
|**Returns** |List of Ethernet interface endpoints.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces'


```




## Single Ethernet interface

|||
|-----------|----------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces/{ethernetInterfaceId}` |
|**Description** |This endpoint retrieves information on a single, logical Ethernet interface or network interface controller \(NIC\).|
|**Returns** |JSON schema representing this Ethernet interface.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{ethernetInterfaceId}'


```


>**Sample response body** 

```
{
	"@odata.context": "/redfish/v1/$metadata#EthernetInterface.EthernetInterface",
	"@odata.etag": "W/\"5DEAF04A\"",
	"@odata.id": "/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1/EthernetInterfaces/1",
	"@odata.type": "#EthernetInterface.v1_4_1.EthernetInterface",
	"Id": "1",
	"FullDuplex": true,
	"IPv4Addresses": [],
	"IPv4StaticAddresses": [],
	"IPv6AddressPolicyTable": [],
	"IPv6Addresses": [],
	"IPv6StaticAddresses": [],
	"IPv6StaticDefaultGateways": [],
	"LinkStatus": null,
	"MACAddress": "80:30:e0:32:0a:58",
	"Name": "",
	"NameServers": [],
	"SpeedMbps": null,
	"StaticNameServers": [],
	"Status": {
		"Health": null,
		"State": null
	},
	"UefiDevicePath": "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)"
}
```



##  PCIeDevice

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}` |
|**Description** | This operation fetches information about a specific PCIe device.<br> |
|**Returns** |Properties of a PCIe device attached to a computer system such as type, the version of the PCIe specification in use by this device, and more.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}'


```
> Sample response body

```
{
    "@odata.context": "/redfish/v1/$metadata#PCIeDevice.PCIeDevice",
    "@odata.etag": "W/\"33150E20\"",
    "@odata.id": "/redfish/v1/Systems/1b77fcdd-b6a2-44b4-83f9-cfb4926fcd79:1/PCIeDevices/1",
    "@odata.type": "#PCIeDevice.v1_5_0.PCIeDevice",
    "Id": "1",
    "Name": "HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeServerPciDevice.HpeServerPciDevice",
            "@odata.etag": "W/\"33150E20\"",
            "@odata.id": "/redfish/v1/Systems/1b77fcdd-b6a2-44b4-83f9-cfb4926fcd79:1/PCIDevices/1",
            "@odata.type": "#HpeServerPciDevice.v2_0_0.HpeServerPciDevice",
            "BusNumber": 2,
            "ClassCode": 2,
            "DeviceID": 5719,
            "DeviceInstance": 1,
            "DeviceLocation": "Embedded",
            "DeviceNumber": 0,
            "DeviceSubInstance": 1,
            "DeviceType": "Embedded LOM",
            "FunctionNumber": 0,
            "Id": "1",
            "LocationString": "Embedded LOM 1",
            "Name": "HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
            "SegmentNumber": 0,
            "StructuredName": "NIC.LOM.1.1",
            "SubclassCode": 0,
            "SubsystemDeviceID": 8894,
            "SubsystemVendorID": 4156,
            "UEFIDevicePath": "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)",
            "VendorID": 5348
        }
    }
}

```





##  Storage

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage` |
|**Description** | This operation lists storage subsystems.<br> A storage subsystem is a set of storage controllers \(physical or virtual\) and the resources such as volumes that can be accessed from that subsystem.<br> |
|**Returns** |Links to storage subsystems.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage'


```





##  Storage subsystem

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}` |
|**Description** | This operation lists resources such as drives and storage controllers in a storage subsystem.<br> |
|**Returns** |Links to the drives and storage controllers of a storage subsystem.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}'


```

> Sample response body 

```
{
    "@odata.context": "/redfish/v1/$metadata#Storage.Storage",
    "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0",
    "@odata.type": "#Storage.v1_10_0.Storage",
    "Description": "HPE Smart Storage Array Controller View",
    "Drives": [
        {
            "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0/Drives/0"
        },
        {
            "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0/Drives/1"
        }
    ],
    "Id": "ArrayController-0",
    "Name": "HpeSmartStorageArrayController",
    "StorageControllers": [
        {
            "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0#/StorageControllers/0",
            "FirmwareVersion": "1.98",
            "Manufacturer": "HPE",
            "MemberId": "0",
            "Model": "HPE Smart Array P408i-a SR Gen10",
            "Name": "HpeSmartStorageArrayController",
            "PartNumber": "836260-001",
            "PhysicalLocation": {
                "PartLocation": {
                    "LocationOrdinalValue": 0,
                    "LocationType": "Slot",
                    "ServiceLabel": "Slot=0"
                }
            },
            "SerialNumber": "PEYHC0DRHBV3CZ ",
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            }
        }
    ],
    "Volumes": {
        "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0/Volumes"
    }
}

```




##  Storage drive

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}` |
|**Description** | This operation retrieves information about a specific storage drive.<br> |
|**Returns** |JSON schema representing this drive.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}'


```




## Volumes



### Collection of volumes

| | |
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>  |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes` |
|<strong>Description</strong>  |This endpoint retrieves a collection of volumes in a specific storage subsystem.|
|<strong>Returns</strong> |A list of links to volumes.|
|<strong>Response code</strong> |On success, `200 OK` |
|<strong>Authentication</strong> |Yes|

 

>**curl command** 

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'


```

>**Sample response body** 

```
{
      "@odata.context":"/redfish/v1/$metadata#VolumeCollection.VolumeCollection",
      "@odata.etag":"W/\"AA6D42B0\"",
      "@odata.id":"/redfish/v1/Systems/eb452cf4-306c-4b21-96fb-698a067da407:1/Storage/ArrayControllers-0/Volumes",
      "@odata.type":"#VolumeCollection.VolumeCollection",
      "Description":"Volume Collection view",
      "Members":[
            {
                  "@odata.id":"/redfish/v1/Systems/eb452cf4-306c-4b21-96fb-698a067da407:1/Storage/ArrayControllers-0/Volumes/1"         
      }      
   ],
      "Members@odata.count":1,
      "Name":"Volume Collection"   
}
```



### Single volume


| | |
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>   |This endpoint retrieves information about a specific volume in a storage subsystem.|
|<strong>Returns</strong>  |JSON schema representing this volume.|
|<strong>Response code</strong>  |On success, `200 OK` |
|<strong>Authentication</strong>  |Yes|

 
>**curl command**
 

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/{volumeId}'


```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#Volume.Volume",
   "@odata.etag":"W/\"46916D5D\"",
   "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Volumes/1",
   "@odata.type":"#Volume.v1_4_1.Volume",
   "CapacityBytes":1200209526784,
   "Encrypted":false,
   "Id":"1",
   "Identifiers":[
      {
         "DurableName":"600508B1001C2AFE083D7F9026B2E994",
         "DurableNameFormat":"NAA"
      }
   ],
   "Links":{
      "Drives":[
         {
            "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Drives/0"
         }
      ]
   },
   "Name":"Drive_Volume_Link",
   "RAIDType":"RAID0",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```


### Creating a volume

| | |
|----------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong>  |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes` |
|<strong>Description</strong>  | This operation creates a volume in a specific storage subsystem.|
|<strong>Response code</strong>   |On success, `200 Ok` |
|<strong>Authentication</strong>|Yes|

>**curl command**

```
curl -i -X POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
   "Name":"Volume_Demo",
   "RAIDType":"RAID1",
   "Drives":[
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/0"
      },
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/1"
      }
   ],
   "@Redfish.OperationApplyTime":"OnReset"
}}' \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'


```

>**Sample request body** 

```
{
   "Name":"Volume_Demo",
   "RAIDType":"RAID1",
   "Drives":[
      {
         "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Drives/0"
      },
      {
         "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Drives/1"
      }
   ],
   "@Redfish.OperationApplyTime":"OnReset"
}}
```

**Request parameters** 

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(required\)<br> |Name of the new volume.|
|RAIDType|String \(required\)<br> |The RAID type of the volume you want to create.|
|Drives\[\{|Array \(required\)<br> |An array of links to drive resources that the new volume contains.|
|@odata.id \}\]<br> |String|A link to a drive resource.|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the operation is carried out.<br> Supported values: `OnReset` and `Immediate`.<br> `OnReset` indicates that the new volume will be available only after you have successfully reset the system. To know how to reset a system, see [Resetting a computer system](#resetting-a-computer-system).<br>`Immediate` indicates that the created volume will be available in the system immediately after the operation is successfully completed.|

>**Sample response body** 

```
 {
      "error":{
            "@Message.ExtendedInfo":[
                  {
                        "MessageId":"iLO.2.13.SystemResetRequired"            
         }         
      ],
            "code":"iLO.0.10.ExtendedInfo",
            "message":"See @Message.ExtendedInfo for more information."      
   }   
}
```





### Deleting a volume


| | |
|----------|-----------|
|<strong>Method</strong>  | `DELETE` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>  | This operation removes a volume in a specific storage subsystem.|
|<strong>Response code</strong>|On success, `204 No Content` |
|<strong>Authentication</strong>  |Yes|

>**curl command**

```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}'


```

>**Sample request body** 

```
{
  
   "@Redfish.OperationApplyTime":"OnReset"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the operation is carried out.<br> Supported values are: `OnReset` and `Immediate`.<br> `OnReset` indicates that the volume will be deleted only after you have successfully reset the system.<br> `Immediate` indicates that the volume will be deleted immediatley after the operation is completed successfully. |




##  SecureBoot

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/SecureBoot` |
|**Description** |Use this endpoint to discover information on `UEFI Secure Boot` and managing the `UEFI Secure Boot` functionality of a specific system.|
|**Returns** | <ul><li>Action for resetting keys.</li><li> `UEFI Secure Boot` properties.<br>**NOTE:**<br> Use URI in the Actions group to discover information about resetting keys.</li></ul>|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/SecureBoot'


```





##  Processors

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors` |
|**Description** |This endpoint lists processors of a specific system.|
|**Returns** |List of processor resource endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors'


```






##  Single processor

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}` |
|**Description** |This endpoint fetches information about the properties of a processor attached to a specific server.|
|**Returns** |JSON schema representing this processor.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}'


```


## Chassis

Chassis represents the physical components of a system—sheet-metal confined spaces, logical zones such as racks, enclosures, chassis and all other containers, and subsystems \(like sensors\).

To view, create, and manage racks or rack groups, ensure that the URP \(Unmanaged Rack Plugin\) is running and is added into the Resource Aggregator for ODIM framework. To know how to add a plugin, see [Adding a plugin as an aggregation source](#adding-a-plugin-as-an-aggregation-source).

>**NOTE:**
URP is installed automatically during the deployment of the resource aggregator.


###  Collection of chassis

|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis` |
|**Description** | This operation lists chassis instances available with Resource Aggregator for ODIM.|
|**Returns** |A collection of links to chassis instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis'

```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
   "@odata.id":"/redfish/v1/Chassis/",
   "@odata.type":"#ChassisCollection.ChassisCollection",
   "Description":"Computer System Chassis view",
   "Name":"Computer System Chassis",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Chassis/ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1"
      },
      { 
         "@odata.id":"/redfish/v1/Chassis/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":2
}
```



 




>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
   "@odata.id":"/redfish/v1/Chassis/",
   "@odata.type":"#ChassisCollection.ChassisCollection",
   "Description":"Computer System Chassis view",
   "Name":"Computer System Chassis",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Chassis/ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1"
      },
      { 
         "@odata.id":"/redfish/v1/Chassis/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":2
}
```





### Single chassis

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}` |
|**Description** |This operation fetches information on a specific computer system chassis, rack group, or a rack.|
|**Returns** |JSON schema representing this chassis instance.|
|**Response code** |On success, `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}'


```

>**Sample response body** 

1. **Computer system chassis**

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.etag":"W/\"50540B90\"",
   "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1",
   "@odata.type":"#Chassis.v1_16_0.Chassis",
   "Id":"192083d2-c60a-4318-967b-cb5890c6dfe4:1",
   "ChassisType":"RackMount",
   "Links":{ 
      "ManagedBy":[ 
         { 
            "@odata.id":"/redfish/v1/Managers/141cbba9-1e99-4272-b855-1781730bfe1c:1"
         }
      ],
      "ComputerSystems":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/192083d2-c60a-4318-967b-cb5890c6dfe4:1"
         }
      ]
   },
   "Manufacturer":"HPE",
   "Model":"ProLiant DL380 Gen10",
   "Name":"Computer System Chassis",
   "NetworkAdapters":{ 
      "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/NetworkAdapters"
   },
   "Oem":{ 
      
         },
         "Firmware":{ 
            "PlatformDefinitionTable":{ 
               "Current":{ 
                  "VersionString":"8.9.0 Build 38"
               }
            },
            "PowerManagementController":{ 
               "Current":{ 
                  "VersionString":"1.0.4"
               }
            },
            "PowerManagementControllerBootloader":{ 
               "Current":{ 
                  "Family":"25",
                  "VersionString":"1.1"
               }
            },
            "SPSFirmwareVersionData":{ 
               "Current":{ 
                  "VersionString":"4.1.4.251"
               }
            },
            "SystemProgrammableLogicDevice":{ 
               "Current":{ 
                  "VersionString":"0x2A"
               }
            }
         },
         "Links":{ 
            "Devices":{ 
               "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/Devices"
            }
         },
         "MCTPEnabledOnServer":true,
         "SmartStorageBattery":[ 
            { 
               "ChargeLevelPercent":100,
               "FirmwareVersion":"0.70",
               "Index":1,
               "MaximumCapWatts":96,
               "Model":"875241-B21",
               "ProductName":"Smart Storage Battery ",
               "RemainingChargeTimeSeconds":7,
               "SerialNumber":"6WQXL0CB2BV63K",
               "SparePartNumber":"878643-001",
               "Status":{ 
                  "Health":"OK",
                  "State":"Enabled"
               }
            }
         ],
         "SystemMaintenanceSwitches":{ 
            "Sw1":"Off",
            "Sw10":"Off",
            "Sw11":"Off",
            "Sw12":"Off",
            "Sw2":"Off",
            "Sw3":"Off",
            "Sw4":"Off",
            "Sw5":"Off",
            "Sw6":"Off",
            "Sw7":"Off",
            "Sw8":"Off",
            "Sw9":"Off"
         }
      }
   },
   "Power":{ 
      "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/Power"
   },
   "SKU":"868704-B21",
   "SerialNumber":"2M291101JX",
   "Status":{ 
      "Health":"OK",
      "State":"Disabled"
   },
   "Thermal":{ 
      "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/Thermal"
   }
}
```

2. **Rack group chassis**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/f4e24c1c-dd2f-5a17-91b7-71620eb070df",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"f4e24c1c-dd2f-5a17-91b7-71620eb070df",
   "Description":"My RackGroup",
   "Name":"RG8",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```

3. **Rack chassis**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```





###  Thermal metrics

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Thermal` |
|**Description** |This operation discovers information on the temperature and cooling of a specific chassis.|
|**Returns** |<ul><li>List of links to Fans</li><li>List of links to Temperatures</li></ul>|
| **Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Thermal'


```




### Collection of network adapters

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters|
|**Description** | This endpoint lists network adapters contained in a chassis. A `NetworkAdapter` represents the physical network adapter capable of connecting to a computer network.<br> Some examples include Ethernet, fibre channel, and converged network adapters.|
|**Returns** |Links to network adapter instances available in this chassis.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/NetworkAdapters'


```


### Single network adapter

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{NetworkAdapterId}` |
|**Description** | This endpoint retrieves information on a specific network adapter.|
|**Returns** |JSON schema representing this network adapter.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|



>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{NetworkAdapterId}'


```


>**Sample response body**

```

{
   "@odata.context":"/redfish/v1/$metadata#NetworkAdapter.NetworkAdapter",
   "@odata.etag":"W/\"F303ECE9\"",
   "@odata.id":"/redfish/v1/Chassis/a022faa5-107c-496d-874e-89c9f3e2df1c:1/NetworkAdapters/{rid}",
   "@odata.type":"#NetworkAdapter.v1_5_0.NetworkAdapter",
   "Description":"The network adapter resource instances available in this chassis.",
   "Name":"Network Adapter View",
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapter.HpeBaseNetworkAdapter",
         "@odata.etag":"W/\"7A9A9CE7\"",
         "@odata.id":"/redfish/v1/Systems/1/BaseNetworkAdapters/1/",
         "@odata.type":"#HpeBaseNetworkAdapter.v2_0_0.HpeBaseNetworkAdapter",
         "Id":"1",
         "FcPorts":[
            
         ],
         "Firmware":{
            "Current":{
               "VersionString":"20.14.54"
            }
         },
         "Name":"HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
         "PhysicalPorts":[
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a4",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            },
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a5",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            },
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a6",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            },
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a7",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            }
         ],
         "Status":{
            "State":null
         },
         "StructuredName":"NIC.LOM.1.1",
         "UEFIDevicePath":"PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)"
      }
   }
}


```






###  Power

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Power` |
|**Description** |This operation retrieves power metrics specific to a server.|
|**Returns** |Information on power consumption and power limiting.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Power'


```

### Creating a rack group

|||
|---------|-------|
|Method | `POST` |
|URI |`/redfish/v1/Chassis`|
|Description |This operation creates a rack group.|
|Returns |<ul><li>`Location` header that contains a link to the created rack group \(highlighted in bold in the sample response header\).</li><li>JSON schema representing the created rack group.<br></li></ul>|
|Response code |On success, `201 Created`|
|Authentication |Yes|

 

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "ChassisType": "RackGroup",
  "Description": "My RackGroup",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/{managerId}"
      }
    ]
  },
  "Name": "RG5"
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis'


```

>**Sample request body**

```
{
  "ChassisType": "RackGroup",
  "Description": "My RackGroup",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
      }
    ]
  },
  "Name": "RG5"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ChassisType|String \(required\)<br> |The type of chassis. The type to be used to create a rack group is RackGroup.<br> |
|Description|String \(optional\)<br> |Description of this rack group.|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack group.|
|ManagedBy \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the manager resources that manage this chassis. The manager resource for racks and rack groups is the URP \(Unmanaged Rack Plugin\) manager. Provide the link to the URP manager.<br> |
|Name|String \(required\)<br> |Name for this rack group.|


>**Sample response header**

```
Connection:keep-alive
Content-Type:application/json; charset=UTF-8
Date:Wed,06 Jan 2021 09:37:43 GMT+15m 26s
**Location:/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d**
Odata-Version:4.0
X-Frame-Options:sameorigin
Content-Length:462 bytes
```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"c2459269-011c-58d3-a217-ef914c4c295d",
   "Description":"My RackGroup",
   "Name":"RG5",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```


### Creating a rack

|||
|---------|-------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/Chassis`|
|**Description**|This operation creates a rack.|
|**Returns** |<ul><li>`Location` header that contains a link to the created rack \(highlighted in bold in the sample response header\).</li><li>JSON schema representing the created rack.<br></li></ul>|
|**Response code** |On success, `201 Created`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "ChassisType": "Rack",
  "Description": "rack number one",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/{managerId}"
      }
    ],
    "ContainedBy": [
      {
	    "@odata.id":"/redfish/v1/Chassis/{chassisId}"
	  }
    ]
  },
  "Name": "RACK#1"
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis'


```

>**Sample request body**

```
{
  "ChassisType": "Rack",
  "Description": "rack number one",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/675560ae-e903-41d9-bfb2-561951999999"
      }
    ],
    "ContainedBy": [
      {
	    "@odata.id":"/redfish/v1/Chassis/1be678f0-86dd-58ac-ac38-16bf0f6dafee"
	  }
    ]
  },
  "Name": "RACK#1"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ChassisType|String \(required\)<br> |The type of chassis. The type to be used to create a rack is Rack.<br> |
|Description|String \(optional\)<br> |Description of this rack.|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack.|
|ManagedBy \[\{<br> @odata.id<br> \}\]<br> |Array \(required\)<br> |An array of links to the manager resources that manage this chassis. The manager resource for racks and rack groups is the URP \(Unmanaged Rack Plugin\) manager. Provide the link to the URP manager.<br> |
|ContainedBy \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the rack groups for containing this rack.|
|Name|String \(required\)<br> |Name for this rack group.|


>**Sample response header**

```
Connection:keep-alive
Content-Type:application/json; charset=UTF-8
Date:Wed,06 Jan 2021 09:37:43 GMT+15m 26s
**Location:/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2**
Odata-Version:4.0
X-Frame-Options:sameorigin
Content-Length:462 bytes
```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```



### Attaching chassis to a rack

|||
|---------|-------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/Chassis/{rackId}`|
|**Description** |This operation attaches chassis to a specific rack.|
|**Returns** |JSON schema for the modified rack having links to the attached chassis.|
|**Response code** |On success, `200 Ok`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "Links": {
    "Contains": [
      {
        "@odata.id": "/redfish/v1/Chassis/{chassisId}"
      }
    ]
  }
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'


```

>**Sample request body**

```
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

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack.|
|Contains \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the computer system chassis resources to be attached to this rack.|




>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "Contains":[
         {
            "@odata.id":"/redfish/v1/Chassis/4159c951-d0d0-4263-858b-0294f5be6377:1"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```


### Detaching chassis from a rack

|||
|---------|-------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/Chassis/{rackId}`|
|**Description** |This operation detaches chassis from a specific rack.|
|**Returns** |JSON schema representing the modified rack.|
|**Response code** |On success, `200 Ok`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "Links": {
    "Contains": []
  }
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'


```

>**Sample request body**

```
{
  "Links": {
    "Contains": []
  }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack.|
|Contains \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the computer system chassis resources to be attached to this rack. To detach chassis from this rack, provide an empty array as value.|




>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```

### Deleting a rack

|||
|---------|-------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/Chassis/{rackId}`|
|**Description** |This operation deletes a specific rack.<br>**IMPORTANT:**<br> If you try to delete a nonempty rack, you will receive an HTTP `409 Conflict` error. Ensure to detach the chassis attached to a rack before deleting it.<br>|
|**Response code** |On success, `204 No Content`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'
```
   

>**Sample request body**

None.

### Deleting a rack group

|||
|---------|-------|
|**Method**| `DELETE` |
|**URI**|`/redfish/v1/Chassis/{rackGroupId}``|
|**Description**|This operation deletes a specific rack group.<br>**IMPORTANT:**<br>If you try to delete a nonempty rack group, you will receive an HTTP `409 Conflict` error. Ensure to remove all the racks contained in a rack group before deleting it.<br>|
|**Response code**|On success, `204 No Content`|
|**Authentication**|Yes|

 

>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   'https://{odim_host}:{port}/redfish/v1/Chassis/{rackGroupId}`'
```
   

>**Sample request body**

None.





##  Searching the inventory

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value}` |
|**Description** | Use this endpoint to search servers based on a filter - combination of a keyword, a condition, and a value.<br> Two ore more filters can be combined in a single request with the help of logical operands.<br>**NOTE:**<br> Only a user with `Login` privilege can perform this operation.|
|**Returns** |Server endpoints based on the specified filter.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value/regular_expression}%20{logicalOperand}%20{searchKeys}%20{conditionKeys}%20{value}'

```



> Sample usage 

```
curl -i GET \
      -H "X-Auth-Token:{X-Auth-Token}" \
    'http://{odimra_host}:{port}/redfish/v1/Systems?$filter=MemorySummary/TotalSystemMemoryGiB%20eq%20384'

```

### Request URI parameters

-  `{searchkeys}` refers to `ComputerSystem` parameters. Following are the allowed search keys:

       -    `ProcessorSummary/Count` 
        
       -   `ProcessorSummary/Model` 
        
       -   `ProcessorSummary/sockets` 
        
       -    `SystemType` 
        
       -   `MemorySummary/TotalSystemMemoryGiB` 
        
       -   `FirmwareVersion` 
        
       -   `Storage/Drives/Quantity` 
        
       -   `Storage/Drives/Capacity` 
        
       -   `Storage/Drives/Type` 
	
-  `{conditionKeys}` refers to Redfish-specified conditions. Following are the allowed condition keys:

    |Condition Key|Meaning|Supported data type|
    |-------------|-------|-------------------|
    |"eq"|Equal to|All data types|
    |"ne"|Not equal to|All data types|
    |"gt"|Greater than|All numeric data types|
    |"ge"|Greater than or equal to|All numeric data types|
    |"le"|Lesser than or equal to|All numeric data types|
    |"lt"|Lesser than|All numeric data types|

     

- `{value}` refers to the actual value of the search parameter or a regular expression. Allowed regular expressions are as follows:

    `*, ?, ., $,%,^,&, /,!` 

      Examples:

         1. `$filter=TotalSystemMemoryGiB%20eq%20**384**` 
        
         2. `$filter=ProcessorSummary/Model%20eq%20**int\***` 
        
         3. `$filter=Storage/Drives/Type%20eq%20HDD` 

-  `{logicalOperands}` refers to the logical operands that are used to combine two or more filters in a request. Allowed logical operands are `and`, `or`, and `not`.



**Sample filters**

**Simple filter examples:**


1. `$filter=TotalSystemMemoryGiB%20eq%20384`
   
	This filter searches a server having total physical memory of 384 GB.



2. `$filter=ProcessorSummary/Model%20eq%20int*`

    This filter searches a server whose processor model name starts with `int` and ends with any combination of letters/numbers/special characters.

  


**Compound filter example:**


`$filter=(ProcessorSummary/Count%20eq%202%20and%20ProcessorSummary/Model%20eq%20intel)%20and%20(TotalSystemMemoryGiB%20eq%20384)`

This filter searches a server having total physical memory of 384 GB and two Intel processors.





 


>**Sample response body**

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
   "@odata.id":"/redfish/v1/Systems/",
   "@odata.type":"#ComputerSystemCollection.ComputerSystemCollection",
   "Description":"Computer Systems view",
   "Name":"Computer Systems",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":1
}
```








# Actions on a computer system



##  Resetting a computer system

|||
|---------|-------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset` |
|**Description** |This action shuts down, powers up, and restarts a specific system.<br>**NOTE:**<br> To reset an aggregate of systems, use this URI:<br>`/redfish/v1/AggregationService/Actions/AggregationService.Reset` <br> See [Resetting servers](#resetting-servers).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. Example registry file name: Base.1.4\). See [Message Registries](#message-registries).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
 curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
  "ResetType":"ForceRestart"
}
' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset'

```



>**Sample request body**

```
{
  "ResetType":"ForceRestart"
}
```

**Request parameters**

Refer to [Resetting Servers](#resetting-servers) to know about `ResetType.` 



>**Sample response body**

```
{
	"error": {
		"@Message.ExtendedInfo": [{
			"MessageId": "Base.1.4.Success"
		}],
		"code": "iLO.0.10.ExtendedInfo",
		"message": "See @Message.ExtendedInfo for more information."
	}
}
```



##  Changing the boot order of a computer system to default settings

|||
|--------|------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder` |
|**Description** |This action changes the boot order of a specific system to default settings.<br>**NOTE:**<br> To change the boot order of an aggregate of systems, use this URI:<br> `/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder` <br> See [Changing the Boot Order of Servers to Default Settings](#changing-the-boot-order-of-servers-to-default-settings).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. Example registry file name: Base.1.4\). See [Message Registries](#message-registries).|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
 curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder'

```

>**Sample response body**

```
{
	"error": {
		"@Message.ExtendedInfo": [{
			"MessageId": "Base.1.4.Success"
		}],
		"code": "iLO.0.10.ExtendedInfo",
		"message": "See @Message.ExtendedInfo for more information."
	}
}
```






##  Changing BIOS settings

|||
|-------|-------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Bios/Settings` |
|**Description** |This action changes BIOS configuration.<br>**NOTE:**<br> Any change in BIOS configuration will be reflected only after the system resets. To see the change, [reset the computer system](#resetting-a-computer-system).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id\). See [Message registries](#message-registries). For example,`MessageId` in the sample response body is `iLO.2.8.SystemResetRequired`. The registry to look up is `iLO.2.8`.<br> |
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
 curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{"Attributes": {"BootMode": "LegacyBios"}}' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{system_id}/Bios/Settings'

```




>**Sample request body**

```
{
	"Attributes": {
		"BootMode": "LegacyBios"
	}
}
```

**Request parameters**

`Attributes` are the list of BIOS attributes specific to the manufacturer or provider. To get a full list of attributes, perform `GET` on:


`https://{odimra_host}:{port}/redfish/v1/Systems/1/Bios/Settings`. 




>**Sample response body**

```
{ 
   "error":{ 
      "@Message.ExtendedInfo":[ 
         { 
            "MessageId":"iLO.2.8.SystemResetRequired"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
}
```








## Changing the boot order settings

|||
|---------|-------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}` |
|**Description** |This action changes the boot order settings of a specific system.|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id\). See [Message Registries](#message-registries). For example,`MessageId` in the sample response body is `Base.1.10.0.Success`. The registry to look up is `Base.1.10.0`.<br> |
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
 curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{ 
   "Boot":{ 
      "BootSourceOverrideTarget":"UefiHttp"
   }
}' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}}'

```




>**Sample request body**

```
{ 
   "Boot":{ 
      "BootSourceOverrideTarget":"UefiHttp"
   }
}
```

**Request parameters**

To get a full list of boot attributes that you can update, perform `GET` on:


`https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}`.


Check attributes under `Boot` in the JSON response. 

Some of the attributes include:

-   `BootSourceOverrideTarget` 

-   `UefiTargetBootSourceOverride` 


For possible values, see values listed under `{attribute}.AllowableValues`. 

> Example:

```
BootSourceOverrideTarget@Redfish.AllowableValues":[
"None",
"Cd",
"Hdd",
"Usb",
"SDCard",
"Utilities",
"Diags",
"BiosSetup",
"Pxe",
"UefiShell",
"UefiHttp",
"UefiTarget"
],
"UefiTargetBootSourceOverride@Redfish.AllowableValues":[
"HD(1,GPT,A6AC4D57-9D6D-46C8-8533-E7787450280D,0x800,0x4E2000)/\\EFI\\redhat\\shimx64.efi",
"PciRoot(0x3)/Pci(0x0,0x0)/Pci(0x0,0x0)/Scsi(0x0,0x0)",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv4(0.0.0.0)",
"UsbClass(0xFFFF,0xFFFF,0xFF,0xFF,0xFF)",
"PciRoot(0x0)/Pci(0x14,0x0)/USB(0x13,0x0)",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv4(0.0.0.0)/Uri()",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv4(0.0.0.0)",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv4(0.0.0.0)/Uri()",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
"HD(1,GPT,F998DA94-45F7-4877-B907-10EDB6E65B07,0x800,0x64000)/\\EFI\\ubuntu\\shimx64.efi",
"HD(2,GPT,0AF5A707-BBEB-4D8A-8016-D840C2516753,0x40800,0x7700000)/\\EFI\\Microsoft\\Boot\\bootmgfw.efi"
]
}
```


**NOTE:**

If you attempt to update `BootSourceOverrideTarget` to `UefiTarget`, when `UefiTargetBootSourceOverride` is set to `None`, you will receive an HTTP `400 Bad Request` error. Update `UefiTargetBootSourceOverride` before setting `BootSourceOverrideTarget` to `UefiTarget`.



>**Sample response body**

```
{ 
   "error":{ 
      "@Message.ExtendedInfo":[ 
         { 
            "MessageId":"Base.1.10.0.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
```
