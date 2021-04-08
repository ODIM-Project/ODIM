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
import com.odim.simulator.tree.RedfishVersion.V1_11_3
import com.odim.simulator.tree.RedfishVersion.V1_12_0
import com.odim.simulator.tree.RedfishVersion.V1_12_1
import com.odim.simulator.tree.RedfishVersion.V1_12_2
import com.odim.simulator.tree.RedfishVersion.V1_13_0
import com.odim.simulator.tree.RedfishVersion.V1_13_1
import com.odim.simulator.tree.RedfishVersion.V1_14_0
import com.odim.simulator.tree.RedfishVersion.V1_1_0
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
import com.odim.simulator.tree.RedfishVersion.V1_5_9
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.RedfishVersion.V1_6_3
import com.odim.simulator.tree.RedfishVersion.V1_6_4
import com.odim.simulator.tree.RedfishVersion.V1_6_5
import com.odim.simulator.tree.RedfishVersion.V1_6_6
import com.odim.simulator.tree.RedfishVersion.V1_6_7
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_7_1
import com.odim.simulator.tree.RedfishVersion.V1_7_2
import com.odim.simulator.tree.RedfishVersion.V1_7_3
import com.odim.simulator.tree.RedfishVersion.V1_7_4
import com.odim.simulator.tree.RedfishVersion.V1_7_5
import com.odim.simulator.tree.RedfishVersion.V1_7_6
import com.odim.simulator.tree.RedfishVersion.V1_7_7
import com.odim.simulator.tree.RedfishVersion.V1_8_0
import com.odim.simulator.tree.RedfishVersion.V1_8_1
import com.odim.simulator.tree.RedfishVersion.V1_8_2
import com.odim.simulator.tree.RedfishVersion.V1_8_3
import com.odim.simulator.tree.RedfishVersion.V1_8_4
import com.odim.simulator.tree.RedfishVersion.V1_8_5
import com.odim.simulator.tree.RedfishVersion.V1_8_6
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.RedfishVersion.V1_9_1
import com.odim.simulator.tree.RedfishVersion.V1_9_2
import com.odim.simulator.tree.RedfishVersion.V1_9_3
import com.odim.simulator.tree.RedfishVersion.V1_9_4
import com.odim.simulator.tree.RedfishVersion.V1_9_5
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.PHYSICAL_SECURITY
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.DRIVES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.LOG_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MEDIA_CONTROLLERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MEMORY_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MEMORY_DOMAINS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.NETWORK_ADAPTERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.PCIE_DEVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.SENSORS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.ASSEMBLY
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.FACILITY
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.PCIE_SLOTS
import com.odim.simulator.tree.structure.ResourceType.POWER
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.RESOURCE_BLOCK
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.SWITCH
import com.odim.simulator.tree.structure.ResourceType.THERMAL
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(CHASSIS)
open class ChassisTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Chassis Description",
                "Name" to "Chassis",
                "ChassisType" to "Rack",
                "Manufacturer" to null,
                "Model" to null,
                "SKU" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "AssetTag" to null,
                "IndicatorLED" to null,
                "Status" to embeddedObject(STATUS),
                "LogServices" to ResourceCollection(LOG_SERVICES_COLLECTION),
                "Thermal" to SingletonResource(THERMAL),
                "Power" to SingletonResource(POWER),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ComputerSystems" to LinkableResourceArray(COMPUTER_SYSTEM),
                        "ManagedBy" to LinkableResourceArray(MANAGER),
                        "ContainedBy" to LinkableResource(CHASSIS),
                        "Contains" to LinkableResourceArray(CHASSIS),
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
        version(V1_0_1, V1_0_0, resourceObject(
                "PowerState" to null
        ))
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
        version(V1_1_0, V1_0_2, resourceObject(
                "PhysicalSecurity" to embeddedObject(PHYSICAL_SECURITY)
        ))
        version(V1_1_2, V1_1_0)
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
        version(V1_2_0, V1_1_2, resourceObject(
                "Location" to embeddedObject(RESOURCE_LOCATION),
                "Links" to embeddedObject(
                        "ManagersInChassis" to LinkableResourceArray(MANAGER),
                        "Drives" to LinkableResourceArray(DRIVE),
                        "Storage" to LinkableResourceArray(STORAGE)
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
        version(V1_3_0, V1_2_0)
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
        version(V1_4_0, V1_3_1, resourceObject(
                "HeightMm" to null,
                "WidthMm" to null,
                "DepthMm" to null,
                "WeightKg" to null,
                "NetworkAdapters" to ResourceCollection(NETWORK_ADAPTERS_COLLECTION),
                "Links" to embeddedObject(
                        "PCIeDevices" to LinkableResourceArray(PCIE_DEVICE)
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
        version(V1_5_0, V1_4_1, embeddedObject(
                "Links" to embeddedObject(
                        "ResourceBlocks" to LinkableResourceArray(RESOURCE_BLOCK)
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
        version(V1_5_9, V1_5_8)
        version(V1_6_0, V1_5_2, resourceObject(
                "Assembly" to SingletonResource(ASSEMBLY)
        ))
        version(V1_6_1, V1_6_0)
        version(V1_6_2, V1_6_1)
        version(V1_6_3, V1_6_2)
        version(V1_6_4, V1_6_3)
        version(V1_6_5, V1_6_4)
        version(V1_6_6, V1_6_5)
        version(V1_6_7, V1_6_6)
        version(V1_7_0, V1_6_0, resourceObject(
                "UUID" to null,
                "Links" to embeddedObject(
                        "Switches" to LinkableResourceArray(SWITCH)
                )
        ))
        version(V1_7_1, V1_7_0)
        version(V1_7_2, V1_7_1)
        version(V1_7_3, V1_7_2)
        version(V1_7_4, V1_7_3)
        version(V1_7_5, V1_7_4)
        version(V1_7_6, V1_7_5)
        version(V1_7_7, V1_7_6)
        version(V1_8_0, V1_7_1, resourceObject(
                "PCIeSlots" to LinkableResource(PCIE_SLOTS)
        ))
        version(V1_8_1, V1_8_0)
        version(V1_8_2, V1_8_1)
        version(V1_8_3, V1_8_2)
        version(V1_8_4, V1_8_3)
        version(V1_8_5, V1_8_4)
        version(V1_8_6, V1_8_5)
        version(V1_9_0, V1_8_1, resourceObject(
                "EnvironmentalClass" to null,
                "Sensors" to ResourceCollection(SENSORS_COLLECTION),
                "Links" to embeddedObject(
                        "Processors" to LinkableResourceArray(PROCESSOR)
                )
        ))
        version(V1_9_1, V1_9_0)
        version(V1_9_2, V1_9_1)
        version(V1_9_3, V1_9_2)
        version(V1_9_4, V1_9_3)
        version(V1_9_5, V1_9_4)
        version(V1_10_0, V1_9_2, resourceObject(
                "PCIeDevices" to LinkableResource(PCIE_DEVICES_COLLECTION)
        ))
        version(V1_10_1, V1_10_0)
        version(V1_10_2, V1_10_1)
        version(V1_10_3, V1_10_2)
        version(V1_11_0, V1_10_0, resourceObject(
                "MediaControllers" to LinkableResource(MEDIA_CONTROLLERS_COLLECTION),
                "Memory" to ResourceCollection(MEMORY_COLLECTION),
                "MemoryDomains" to ResourceCollection(MEMORY_DOMAINS_COLLECTION),
                "Links" to embeddedObject(
                        "Facility" to LinkableResource(FACILITY)
                )
        ))
        version(V1_11_1, V1_11_0)
        version(V1_11_2, V1_11_1)
        version(V1_11_3, V1_11_2)
        version(V1_12_0, V1_11_1, resourceObject(
                "MaxPowerWatts" to null,
                "MinPowerWatts" to null
        ))
        version(V1_12_1, V1_12_0)
        version(V1_12_2, V1_12_1)
        version(V1_13_0, V1_12_1)
        version(V1_13_1, V1_13_0)
        version(V1_14_0, V1_13_1, resourceObject(
                "LocationIndicatorActive" to null,
                "Drives" to ResourceCollection(DRIVES_COLLECTION)
        ))
    }
}
