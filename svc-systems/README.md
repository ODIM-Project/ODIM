#  Resource inventory

Resource Aggregator for ODIM allows you to view the inventory of compute and local storage resources through Redfish `Systems` and `Chassis` endpoints. It also offers the capability to search inventory information based on one or more configuration parameters and exposes APIs to manage the added resources.

To discover crucial configuration information about a resource, including chassis, perform `GET` on these endpoints.


<aside class="notice">
To access Redfish `Systems` and `Chassis` endpoints, ensure that you have a minimum privilege of `Login`. If you do not have the necessary privileges, you will receive an HTTP `403 Forbidden` error.
</aside>
  
##  Modifying Configurations of Systems Service
  
Config File of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer the section **Modifying Configurations** in the README.md to change the configurations of a odimra service
  
**Specific configurations for Systems Service are:**
  
##  Log Location of the Systems Service
  
/var/log/ODIMRA/system.log
    
  

## Supported endpoints

|||
|-------|--------------------|
|/redfish/v1/Systems|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}|`GET`, `PATCH`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory/\{memoryId\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/MemoryDomains|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/NetworkInterfaces|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces/\{Id\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Bios|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/SecureBoot|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes|`GET` , `POST`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes/\{volumeId\}|`GET`, `DELETE`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors/\{Id\}|`GET`|
|/redfish/v1/Systems?$filter=\{searchKeys\}%20\{conditionKeys\}%20\{value\}|`GET`|
| /redfish/v1/Systems/\{ComputerSystemId\}/Bios/Settings<br> |`GET`, `PATCH`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.Reset|`POST`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.SetDefaultBootOrder|`POST`|

|||
|-------|--------------------|
|/redfish/v1/Chassis|`GET`|
|/redfish/v1/Chassis/\{chassisId\}|`GET`|
|/redfish/v1/Chassis/\{chassisId\}/Thermal|`GET`|
|/redfish/v1/Chassis/\{chassisId\}/NetworkAdapters|`GET`|







##  Collection of computer systems

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems'


```

> Sample response body 

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


|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems` |
|**Description** |This operation lists all systems available with Resource Aggregator for ODIM.|
|**Returns** |A collection of links to computer system instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|




## Single computer system

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}'

```

> Sample response body 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
   "@odata.etag":"W/\"8C36EBD2\"",
   "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1",
   "@odata.type":"#ComputerSystem.v1_4_0.ComputerSystem",
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



|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}` |
|**Description** |This endpoint fetches information about a specific system.|
|**Returns** |JSON schema representing this computer system instance.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

 



##  Memory collection


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory'


```

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Memory` |
|**Description** |This operation lists all memory devices of a specific server.|
|**Returns** |List of memory resource endpoints.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|



## Single memory

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}'


```

> Sample response body 

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



|||
|---------|-------|
|**Method** |GET|
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}` |
|**Description** |This endpoint retrieves configuration information of specific memory.|
|**Returns** |JSON schema representing this memory resource.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

 








##  Memory domains

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains'


```
|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains` |
|**Description** | This endpoint lists memory domains of a specific system.<br> Memory Domains indicate to the client which Memory \(DIMMs\) can be grouped in Memory Chunks to form interleave sets or otherwise grouped.<br> |
|**Returns** |List of memory domain endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|



##  BIOS

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Bios'


```


|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Bios` |
|**Description** | Use this endpoint to discover system-specific information about a BIOS resource and actions for changing to BIOS settings.<br>**NOTE:**<br> Changes to the BIOS typically require a system reset before they take effect.|
|**Returns** |<ul><li>Actions for changing password and resetting BIOS.</li><li>BIOS attributes.</li></ul> |
|**Response code** |`200 OK` |
|**Authentication** |Yes|


## Network interfaces

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces'


```


|||
|--------|---------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces` |
|**Description** | This endpoint lists network interfaces of a specific system.<br> A network interface contains links to network adapter, network port, and network device function resources.<br> |
|**Returns** |List of network interface endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|



##  Ethernet interfaces


```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces'


```

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces` |
|**Description** |This endpoint lists Ethernet interfaces or network interface controller \(NIC\) of a specific system.|
|**Returns** |List of Ethernet interface endpoints.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|


## Single Ethernet interface


```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{ethernetInterfaceId}'


```


> Sample response body 

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



|||
|-----------|----------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces/{ethernetInterfaceId}` |
|**Description** |This endpoint retrieves information on a single, logical Ethernet interface or network interface controller \(NIC\).|
|**Returns** |JSON schema representing this Ethernet interface.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|





##  Storage

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage'


```


|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage` |
|**Description** | This operation lists storage subsystems.<br> A storage subsystem is a set of storage controllers \(physical or virtual\) and the resources such as volumes that can be accessed from that subsystem.<br> |
|**Returns** |Links to storage subsystems.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


##  Storage subsystem


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}'


```

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}` |
|**Description** | This operation lists resources such as drives and storage controllers in a storage subsystem.<br> |
|**Returns** |Links to the drives and storage controllers of a storage subsystem.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|




##  Storage drive

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}'


```

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}` |
|**Description** | This operation retrieves information about a specific storage drive.<br> |
|**Returns** |JSON schema representing this drive.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


## Volumes

### A collection of volumes

| | | 
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>  |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes` |
|<strong>Description</strong>  |This endpoint retrieves a collection of volumes in a specific storage subsystem.|
|<strong>Returns</strong> |A list of links to volumes.|
|<strong>Response Code</strong> |On success, `200 OK` |
|<strong>Authentication</strong> |Yes|

 

 

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'


```

> Sample response body 

```
{
   ​   "@odata.context":"/redfish/v1/$metadata#VolumeCollection.VolumeCollection",
   ​   "@odata.etag":"W/\"AA6D42B0\"",
   ​   "@odata.id":"/redfish/v1/Systems/eb452cf4-306c-4b21-96fb-698a067da407:1/Storage/ArrayControllers-0/Volumes",
   ​   "@odata.type":"#VolumeCollection.VolumeCollection",
   ​   "Description":"Volume Collection view",
   ​   "Members":​[
      ​      {
         ​         "@odata.id":"/redfish/v1/Systems/eb452cf4-306c-4b21-96fb-698a067da407:1/Storage/ArrayControllers-0/Volumes/1"         ​
      }      ​
   ],
   ​   "Members@odata.count":1,
   ​   "Name":"Volume Collection"   ​
}
```



### Single volume


| | | 
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>   |This endpoint retrieves information about a specific volume in a storage subsystem.|
|<strong>Returns</strong>  |JSON schema representing this volume.|
|<strong>Response Code</strong>  |On success, `200 OK` |
|<strong>Authentication</strong>  |Yes|

 

 

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/{volumeId}'


```

> Sample request body 

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
|<strong>Response Code</strong>   |On success, `200 Ok` |
|<strong>Authentication</strong>|Yes|



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

> Sample request body 

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

### Request parameters 

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(required\)<br> |Name of the new volume.|
|RAIDType|String \(required\)<br> |The RAID type of the volume you want to create.|
|Drives\[\{|Array \(required\)<br> |An array of links to drive resources that the new volume contains.|
|@odata.id \}\]<br> |String|A link to a drive resource.|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the operation is carried out.<br> Supported value is: `OnReset` and `Immediate`. `OnReset` indicates that the operation will be carried out only after you reset the system.|

> Sample response body 

```
 {
   ​   "error":{
      ​      "@Message.ExtendedInfo":[
         ​         {
            ​            "MessageId":"iLO.2.13.SystemResetRequired"            ​
         }         ​
      ],
      ​      "code":"iLO.0.10.ExtendedInfo",
      ​      "message":"See @Message.ExtendedInfo for more information."      ​
   }   ​
}
```






### Deleting a volume


| | | 
|----------|-----------|
|<strong>Method</strong>  | `DELETE` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>  | This operation removes a volume in a specific storage subsystem.<br> |
|<strong>Response Code</strong>|On success, `204 No Content` |
|<strong>Authentication</strong>  |Yes|



```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}'


```

> Sample request body 

```
{
  
   "@Redfish.OperationApplyTime":"OnReset"
}
```

### Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the operation is carried out.<br> Supported value is: `OnReset`. Supported values are: `OnReset` and `Immediate`. `OnReset` indicates that the volume will be deleted only after you reset the system.<br> |





##  SecureBoot


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/SecureBoot'


```


|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/SecureBoot` |
|**Description** |Use this endpoint to discover information on `UEFI Secure Boot` and managing the `UEFI Secure Boot` functionality of a specific system.|
|**Returns** | <ul><li>Action for resetting keys.</li><li> `UEFI Secure Boot` properties.<br>**NOTE:**<br> Use URI in the Actions group to discover information about resetting keys.</li></ul>|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


##  Processors

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors'


```


|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors` |
|**Description** |This endpoint lists processors of a specific system.|
|**Returns** |List of processor resource endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|



##  Single processor



```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}'


```

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}` |
|**Description** |This endpoint fetches information about the properties of a processor attached to a specific server.|
|**Returns** |JSON schema representing this processor.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|



##  Collection of chassis

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis'

```

> Sample response body 

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

|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis` |
|**Description** | This operation lists chassis instances available with Resource Aggregator for ODIM.<br> Chassis represents the physical components of a system - sheet-metal confined spaces, logical zones such as racks, enclosures, chassis and all other containers, and subsystems \(like sensors\).<br> |
|**Returns** |A collection of links to chassis instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

 




> Sample response body 

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





## Single chassis



```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}'


```

> Sample Response body 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.etag":"W/\"50540B90\"",
   "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1",
   "@odata.type":"#Chassis.v1_6_0.Chassis",
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
               "ProductName":"HPE Smart Storage Battery ",
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



|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}` |
|**Description** |This operation fetches information on a specific chassis.|
|**Returns** |JSON schema representing this chassis instance.|
|**Response code** |On success, `200 OK` |
|**Authentication** |Yes|



 
##  Thermal metrics

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Thermal'


```

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Thermal` |
|**Description** |This operation discovers information on the temperature and cooling of a specific chassis.|
|**Returns** |<ul><li>List of links to Fans</li><li>List of links to Temperatures</li></ul>|
| **Response code** | `200 OK` |
|**Authentication** |Yes|


##  Network adapters

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/NetworkAdapters'


```


|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters` |
|**Description** | Use this endpoint to discover information on network adapters. Some examples of network adapters include Ethernet, fibre channel, and converged network adapters.<br> A `NetworkAdapter` represents the physical network adapter capable of connecting to a computer network.<br> |
|**Returns** |Links to network adapter instances available in this chassis.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|



##  Power

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Power'


```


|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Power` |
|**Description** |This operation retrieves power metrics specific to a server.|
|**Returns** |Information on power consumption and power limiting.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


##  Searching the inventory

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value/regular_expression}%20{logicalOperand}%20{searchKeys}%20{conditionKeys}%20{value}'

```

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value}` |
|**Description** | Use this endpoint to search servers based on a filter - combination of a keyword, a condition, and a value.<br> Two ore more filters can be combined in a single request with the help of logical operands.<br>**NOTE:**<br> Only a user with `Login` privilege can perform this operation.|
|**Returns** |Server endpoints based on the specified filter.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

 



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




 
 
 

> Sample response body

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


|||
|---------|-------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset` |
|**Description** |This action shuts down, powers up, and restarts a specific system.<br>**NOTE:**<br> To reset an aggregate of systems, use this URI:<br>`/redfish/v1/AggregationService/Actions/AggregationService.Reset` <br> See [Resetting servers](#resetting-servers).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. Example registry file name: Base.1.4\). See [Message Registries](#message-registries).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

 




> Sample request body

```
{
  "ResetType":"ForceRestart"
}
```

### Request parameters

Refer to [Resetting Servers](#resetting-servers) to know about `ResetType.` 



> Sample response body

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


```
 curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder'

```

> Sample response body

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

|||
|--------|------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder` |
|**Description** |This action changes the boot order of a specific system to default settings.<br>**NOTE:**<br> To change the boot order of an aggregate of systems, use this URI:<br> `/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder` <br> See [Changing the Boot Order of Servers to Default Settings](#changing-the-boot-order-of-servers-to-default-settings).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. Example registry file name: Base.1.4\). See [Message Registries](#message-registries).|
|**Response code** |`200 OK` |
|**Authentication** |Yes|




##  Changing BIOS settings


```
 curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{"Attributes": {"BootMode": "LegacyBios"}}' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{system_id}/Bios/Settings'

```


|||
|-------|-------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Bios/Settings` |
|**Description** |This action changes BIOS configuration.<br>**NOTE:**<br> Any change in BIOS configuration will be reflected only after the system resets. To see the change, [reset the computer system](#resetting-a-computer-system).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id\). See [Message registries](#message-registries). For example,`MessageId` in the sample response body is `iLO.2.8.SystemResetRequired`. The registry to look up is `iLO.2.8`.<br> |
|**Response code** | `200 OK` |
|**Authentication** |Yes|

 




> Sample request body

```
{
	"Attributes": {
		"BootMode": "LegacyBios"
	}
}
```

### Request parameters

`Attributes` are the list of BIOS attributes specific to the manufacturer or provider. To get a full list of attributes, perform `GET` on:


`https://{odimra_host}:{port}/redfish/v1/Systems/1/Bios/Settings`. 


Some of the attributes include:

-   `BootMode` 

-   `NicBoot1` 

-   `PowerProfile` 

-   `AdminPhone` 

-   `ProcCoreDisable` 

-   `UsbControl` 


> Sample response body

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

```
 curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{ 
   "Boot":{ 
      "BootSourceOverrideTarget":"Usb"
   }
}' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}}'

```

|||
|---------|-------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}` |
|**Description** |This action changes the boot order settings of a specific system.|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id\). See [Message Registries](#message-registries). For example,`MessageId` in the sample response body is `Base.1.0.Success`. The registry to look up is `Base.1.0`.<br> |
|**Response code** |`200 OK` |
|**Authentication** |Yes|

 





> Sample request body

```
{ 
   "Boot":{ 
      "BootSourceOverrideTarget":"Usb"
   }
}
```

### Request parameters

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



> Sample response body

```
{ 
   "error":{ 
      "@Message.ExtendedInfo":[ 
         { 
            "MessageId":"Base.1.0.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
```
