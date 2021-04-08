// Copyright (c) Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package com.odim.simulator.tree.templates.redfish

import com.odim.simulator.tree.RedfishVersion.V1_0_0
import com.odim.simulator.tree.RedfishVersion.V1_0_1
import com.odim.simulator.tree.RedfishVersion.V1_0_10
import com.odim.simulator.tree.RedfishVersion.V1_0_11
import com.odim.simulator.tree.RedfishVersion.V1_0_12
import com.odim.simulator.tree.RedfishVersion.V1_0_13
import com.odim.simulator.tree.RedfishVersion.V1_0_14
import com.odim.simulator.tree.RedfishVersion.V1_0_15
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_0_8
import com.odim.simulator.tree.RedfishVersion.V1_0_9
import com.odim.simulator.tree.RedfishVersion.V1_10_0
import com.odim.simulator.tree.RedfishVersion.V1_10_1
import com.odim.simulator.tree.RedfishVersion.V1_10_2
import com.odim.simulator.tree.RedfishVersion.V1_10_3
import com.odim.simulator.tree.RedfishVersion.V1_11_0
import com.odim.simulator.tree.RedfishVersion.V1_11_1
import com.odim.simulator.tree.RedfishVersion.V1_11_2
import com.odim.simulator.tree.RedfishVersion.V1_12_0
import com.odim.simulator.tree.RedfishVersion.V1_12_1
import com.odim.simulator.tree.RedfishVersion.V1_13_0
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_10
import com.odim.simulator.tree.RedfishVersion.V1_1_11
import com.odim.simulator.tree.RedfishVersion.V1_1_12
import com.odim.simulator.tree.RedfishVersion.V1_1_13
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_1_7
import com.odim.simulator.tree.RedfishVersion.V1_1_8
import com.odim.simulator.tree.RedfishVersion.V1_1_9
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_10
import com.odim.simulator.tree.RedfishVersion.V1_2_11
import com.odim.simulator.tree.RedfishVersion.V1_2_12
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_2_6
import com.odim.simulator.tree.RedfishVersion.V1_2_7
import com.odim.simulator.tree.RedfishVersion.V1_2_8
import com.odim.simulator.tree.RedfishVersion.V1_2_9
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_10
import com.odim.simulator.tree.RedfishVersion.V1_3_11
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_3_4
import com.odim.simulator.tree.RedfishVersion.V1_3_5
import com.odim.simulator.tree.RedfishVersion.V1_3_6
import com.odim.simulator.tree.RedfishVersion.V1_3_7
import com.odim.simulator.tree.RedfishVersion.V1_3_8
import com.odim.simulator.tree.RedfishVersion.V1_3_9
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_10
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_4_3
import com.odim.simulator.tree.RedfishVersion.V1_4_4
import com.odim.simulator.tree.RedfishVersion.V1_4_5
import com.odim.simulator.tree.RedfishVersion.V1_4_6
import com.odim.simulator.tree.RedfishVersion.V1_4_7
import com.odim.simulator.tree.RedfishVersion.V1_4_8
import com.odim.simulator.tree.RedfishVersion.V1_4_9
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_5_1
import com.odim.simulator.tree.RedfishVersion.V1_5_2
import com.odim.simulator.tree.RedfishVersion.V1_5_3
import com.odim.simulator.tree.RedfishVersion.V1_5_4
import com.odim.simulator.tree.RedfishVersion.V1_5_5
import com.odim.simulator.tree.RedfishVersion.V1_5_6
import com.odim.simulator.tree.RedfishVersion.V1_5_7
import com.odim.simulator.tree.RedfishVersion.V1_5_8
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.RedfishVersion.V1_6_3
import com.odim.simulator.tree.RedfishVersion.V1_6_4
import com.odim.simulator.tree.RedfishVersion.V1_6_5
import com.odim.simulator.tree.RedfishVersion.V1_6_6
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_7_1
import com.odim.simulator.tree.RedfishVersion.V1_7_2
import com.odim.simulator.tree.RedfishVersion.V1_7_3
import com.odim.simulator.tree.RedfishVersion.V1_7_4
import com.odim.simulator.tree.RedfishVersion.V1_7_5
import com.odim.simulator.tree.RedfishVersion.V1_8_0
import com.odim.simulator.tree.RedfishVersion.V1_8_1
import com.odim.simulator.tree.RedfishVersion.V1_8_2
import com.odim.simulator.tree.RedfishVersion.V1_8_3
import com.odim.simulator.tree.RedfishVersion.V1_8_4
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.RedfishVersion.V1_9_1
import com.odim.simulator.tree.RedfishVersion.V1_9_2
import com.odim.simulator.tree.RedfishVersion.V1_9_3
import com.odim.simulator.tree.RedfishVersion.V1_9_4
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.ADD_RESOURCE_BLOCK
import com.odim.simulator.tree.structure.ActionType.REMOVE_RESOURCE_BLOCK
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.ActionType.SET_DEFAULT_BOOT_ORDER
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.BOOT
import com.odim.simulator.tree.structure.EmbeddedObjectType.HOSTED_SERVICES
import com.odim.simulator.tree.structure.EmbeddedObjectType.HOST_GRAPHICAL_CONSOLE
import com.odim.simulator.tree.structure.EmbeddedObjectType.HOST_SERIAL_CONSOLE
import com.odim.simulator.tree.structure.EmbeddedObjectType.MEMORY_SUMMARY
import com.odim.simulator.tree.structure.EmbeddedObjectType.PROCESSOR_SUMMARY
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedObjectType.TRUSTED_MODULES
import com.odim.simulator.tree.structure.EmbeddedObjectType.VIRTUAL_MEDIA_CONFIG
import com.odim.simulator.tree.structure.EmbeddedObjectType.WATCHDOG_TIMER
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.ETHERNET_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.FABRIC_ADAPTERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.LOG_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MEMORY_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MEMORY_DOMAINS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.NETWORK_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.PROCESSORS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.SIMPLE_STORAGES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.VIRTUAL_MEDIAS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.BIOS
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.REDUNDANCY
import com.odim.simulator.tree.structure.ResourceType.RESOURCE_BLOCK
import com.odim.simulator.tree.structure.ResourceType.SECURE_BOOT
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(COMPUTER_SYSTEM)
open class ComputerSystemTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Computer System Description",
                "Name" to "Computer System",
                "SystemType" to "Physical",
                "AssetTag" to null,
                "Manufacturer" to null,
                "Model" to null,
                "SKU" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "UUID" to null,
                "HostName" to null,
                "IndicatorLED" to null,
                "PowerState" to null,
                "Boot" to embeddedObject(BOOT),
                "BiosVersion" to null,
                "ProcessorSummary" to embeddedObject(PROCESSOR_SUMMARY),
                "MemorySummary" to embeddedObject(MEMORY_SUMMARY),
                "Processors" to ResourceCollection(PROCESSORS_COLLECTION),
                "EthernetInterfaces" to ResourceCollection(ETHERNET_INTERFACES_COLLECTION),
                "SimpleStorage" to ResourceCollection(SIMPLE_STORAGES_COLLECTION),
                "LogServices" to ResourceCollection(LOG_SERVICES_COLLECTION),
                "Status" to embeddedObject(STATUS),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Chassis" to LinkableResourceArray(CHASSIS),
                        "ManagedBy" to LinkableResourceArray(MANAGER),
                        "PoweredBy" to LinkableResourceArray(ANY),
                        "CooledBy" to LinkableResourceArray(ANY)
                ),
                "Actions" to Actions(
                        Action(RESET, "ResetType", mutableListOf(
                                "On",
                                "ForceOff",
                                "GracefulShutdown",
                                "GracefulRestart",
                                "ForceRestart",
                                "Nmi",
                                "ForceOn",
                                "PushPowerButton",
                                "PowerCycle"
                        ))
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_0_9, V1_0_8)
        version(V1_0_10, V1_0_9)
        version(V1_0_11, V1_0_10)
        version(V1_0_12, V1_0_11)
        version(V1_0_13, V1_0_12)
        version(V1_0_14, V1_0_13)
        version(V1_0_15, V1_0_14)
        version(V1_1_0, V1_0_2, resourceObject(
                "TrustedModules" to embeddedArray(TRUSTED_MODULES),
                "SecureBoot" to SingletonResource(SECURE_BOOT),
                "Bios" to SingletonResource(BIOS),
                "Memory" to ResourceCollection(MEMORY_COLLECTION),
                "Storage" to ResourceCollection(STORAGES_COLLECTION)
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_1_7, V1_1_6)
        version(V1_1_8, V1_1_7)
        version(V1_1_9, V1_1_8)
        version(V1_1_10, V1_1_9)
        version(V1_1_11, V1_1_10)
        version(V1_1_12, V1_1_11)
        version(V1_1_13, V1_1_12)
        version(V1_2_0, V1_1_1, resourceObject(
                "HostingRoles" to embeddedArray(),
                "PCIeDevices" to LinkableResourceArray(PCIE_DEVICE),
                "PCIeFunctions" to LinkableResourceArray(PCIE_FUNCTION),
                "HostedServices" to embeddedObject(HOSTED_SERVICES),
                "MemoryDomains" to ResourceCollection(MEMORY_DOMAINS_COLLECTION),
                "Links" to embeddedObject(
                        "Endpoints" to LinkableResourceArray(ENDPOINT)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_2_6, V1_2_5)
        version(V1_2_7, V1_2_6)
        version(V1_2_8, V1_2_7)
        version(V1_2_9, V1_2_8)
        version(V1_2_10, V1_2_9)
        version(V1_2_11, V1_2_10)
        version(V1_2_12, V1_2_11)
        version(V1_3_0, V1_2_1, resourceObject(
                "NetworkInterfaces" to ResourceCollection(NETWORK_INTERFACES_COLLECTION)
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_3_4, V1_3_3)
        version(V1_3_5, V1_3_4)
        version(V1_3_6, V1_3_5)
        version(V1_3_7, V1_3_6)
        version(V1_3_8, V1_3_7)
        version(V1_3_9, V1_3_8)
        version(V1_3_10, V1_3_9)
        version(V1_3_11, V1_3_10)
        version(V1_4_0, V1_3_1, embeddedObject(
                "Links" to embeddedObject(
                        "ResourceBlocks" to LinkableResourceArray(RESOURCE_BLOCK)
                )
        ))
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_4_3, V1_4_2)
        version(V1_4_4, V1_4_3)
        version(V1_4_5, V1_4_4)
        version(V1_4_6, V1_4_5)
        version(V1_4_7, V1_4_6)
        version(V1_4_8, V1_4_7)
        version(V1_4_9, V1_4_8)
        version(V1_4_10, V1_4_9)
        version(V1_5_0, V1_4_2, resourceObject(
                "Redundancy" to EmbeddedResourceArray(REDUNDANCY),
                "HostWatchdogTimer" to embeddedObject(WATCHDOG_TIMER),
                "SubModel" to null,
                "Links" to embeddedObject(
                        "ConsumingComputerSystems" to LinkableResourceArray(COMPUTER_SYSTEM),
                        "SupplyingComputerSystems" to LinkableResourceArray(COMPUTER_SYSTEM)
                ),
                "Actions" to Actions(
                        Action(SET_DEFAULT_BOOT_ORDER)
                )
        ))
        version(V1_5_1, V1_5_0)
        version(V1_5_2, V1_5_1)
        version(V1_5_3, V1_5_2)
        version(V1_5_4, V1_5_3)
        version(V1_5_5, V1_5_4)
        version(V1_5_6, V1_5_5)
        version(V1_5_7, V1_5_6)
        version(V1_5_8, V1_5_7)
        version(V1_6_0, V1_5_2, resourceObject(
                "PowerRestorePolicy" to "AlwaysOn",
                "Actions" to Actions(
                        Action(ADD_RESOURCE_BLOCK),
                        Action(REMOVE_RESOURCE_BLOCK)
                )
        ))
        version(V1_6_1, V1_6_0)
        version(V1_6_2, V1_6_1)
        version(V1_6_3, V1_6_2)
        version(V1_6_4, V1_6_3)
        version(V1_6_5, V1_6_4)
        version(V1_6_6, V1_6_5)
        version(V1_7_0, V1_6_1)
        version(V1_7_1, V1_7_0)
        version(V1_7_2, V1_7_1)
        version(V1_7_3, V1_7_2)
        version(V1_7_4, V1_7_3)
        version(V1_7_5, V1_7_4)
        version(V1_8_0, V1_7_0)
        version(V1_8_1, V1_8_0)
        version(V1_8_2, V1_8_1)
        version(V1_8_3, V1_8_2)
        version(V1_8_4, V1_8_3)
        version(V1_9_0, V1_8_0)
        version(V1_9_1, V1_9_0)
        version(V1_9_2, V1_9_1)
        version(V1_9_3, V1_9_2)
        version(V1_9_4, V1_9_3)
        version(V1_10_0, V1_9_1, resourceObject(
                "FabricAdapters" to LinkableResource(FABRIC_ADAPTERS_COLLECTION)
        ))
        version(V1_10_1, V1_10_0)
        version(V1_10_2, V1_10_1)
        version(V1_10_3, V1_10_2)
        version(V1_11_0, V1_10_1)
        version(V1_11_1, V1_11_0)
        version(V1_11_2, V1_11_1)
        version(V1_12_0, V1_11_1, resourceObject(
                "LastResetTime" to "2017-04-14T06:35:05Z"
        ))
        version(V1_12_1, V1_12_0)
        version(V1_13_0, V1_12_1, resourceObject(
                "LocationIndicatorActive" to null,
                "BootProgress" to null,
                "PowerOnDelaySeconds" to null,
                "PowerOffDelaySeconds" to null,
                "PowerCycleDelaySeconds" to null,
                "SerialConsole" to embeddedObject(HOST_SERIAL_CONSOLE),
                "GraphicalConsole" to embeddedObject(HOST_GRAPHICAL_CONSOLE),
                "VirtualMediaConfig" to embeddedObject(VIRTUAL_MEDIA_CONFIG),
                "VirtualMedia" to ResourceCollection(VIRTUAL_MEDIAS_COLLECTION)
        ))
    }
}
