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
import com.odim.simulator.tree.RedfishVersion.V1_0_10
import com.odim.simulator.tree.RedfishVersion.V1_0_11
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_0_8
import com.odim.simulator.tree.RedfishVersion.V1_0_9
import com.odim.simulator.tree.RedfishVersion.V1_10_0
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_1_7
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_2_6
import com.odim.simulator.tree.RedfishVersion.V1_2_7
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_3_4
import com.odim.simulator.tree.RedfishVersion.V1_3_5
import com.odim.simulator.tree.RedfishVersion.V1_3_6
import com.odim.simulator.tree.RedfishVersion.V1_3_7
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_4_3
import com.odim.simulator.tree.RedfishVersion.V1_4_4
import com.odim.simulator.tree.RedfishVersion.V1_4_5
import com.odim.simulator.tree.RedfishVersion.V1_4_6
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_5_1
import com.odim.simulator.tree.RedfishVersion.V1_5_2
import com.odim.simulator.tree.RedfishVersion.V1_5_3
import com.odim.simulator.tree.RedfishVersion.V1_5_4
import com.odim.simulator.tree.RedfishVersion.V1_5_5
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.RedfishVersion.V1_6_3
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_7_1
import com.odim.simulator.tree.RedfishVersion.V1_7_2
import com.odim.simulator.tree.RedfishVersion.V1_8_0
import com.odim.simulator.tree.RedfishVersion.V1_8_1
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.FPGA
import com.odim.simulator.tree.structure.EmbeddedObjectType.PROCESSOR_ID
import com.odim.simulator.tree.structure.EmbeddedObjectType.PROCESSOR_INTERFACE
import com.odim.simulator.tree.structure.EmbeddedObjectType.PROCESSOR_MEMORY
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.ACCELERATION_FUNCTIONS_COLLETION
import com.odim.simulator.tree.structure.ResourceCollectionType.OPERATING_CONFIGS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.PROCESSORS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ASSEMBLY
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.OPERATING_CONFIG
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR_METRICS
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(PROCESSOR)
open class ProcessorTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Processor Description",
                "Name" to "Processor",
                "Socket" to null,
                "ProcessorType" to null,
                "ProcessorArchitecture" to null,
                "InstructionSet" to null,
                "ProcessorId" to embeddedObject(PROCESSOR_ID),
                "Status" to embeddedObject(STATUS),
                "Manufacturer" to null,
                "Model" to null,
                "MaxSpeedMHz" to null,
                "TotalCores" to null,
                "TotalThreads" to null
        ))
        version(V1_0_2, V1_0_0)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_0_9, V1_0_8)
        version(V1_0_10, V1_0_9)
        version(V1_0_11, V1_0_10)
        version(V1_1_0, V1_0_4, resourceObject(
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Chassis" to LinkableResource(CHASSIS)
                )
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_1_7, V1_1_6)
        version(V1_2_0, V1_1_0, resourceObject(
                "Location" to embeddedObject(RESOURCE_LOCATION),
                "Assembly" to SingletonResource(ASSEMBLY)
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_2_6, V1_2_5)
        version(V1_2_7, V1_2_6)
        version(V1_3_0, V1_2_0, resourceObject(
                "SubProcessors" to ResourceCollection(PROCESSORS_COLLECTION)
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_3_4, V1_3_3)
        version(V1_3_5, V1_3_4)
        version(V1_3_6, V1_3_5)
        version(V1_3_7, V1_3_6)
        version(V1_4_0, V1_3_2, resourceObject(
                "TDPWatts" to null,
                "MaxTDPWatts" to null,
                "Metrics" to SingletonResource(PROCESSOR_METRICS),
                "UUID" to null,
                "ProcessorMemory" to embeddedArray(PROCESSOR_MEMORY),
                "FPGA" to embeddedObject(FPGA),
                "AccelerationFunctions" to ResourceCollection(ACCELERATION_FUNCTIONS_COLLETION),
                "Links" to embeddedObject(
                        "Endpoints" to LinkableResourceArray(ENDPOINT),
                        "ConnectedProcessors" to LinkableResourceArray(PROCESSOR),
                        "PCIeDevice" to LinkableResource(PCIE_DEVICE),
                        "PCIeFunctions" to LinkableResourceArray(PCIE_FUNCTION)
                )
        ))
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_4_3, V1_4_2)
        version(V1_4_4, V1_4_3)
        version(V1_4_5, V1_4_4)
        version(V1_4_6, V1_4_5)
        version(V1_5_0, V1_4_1, resourceObject(
                "TotalEnabledCores" to null
        ))
        version(V1_5_1, V1_5_0)
        version(V1_5_2, V1_5_1)
        version(V1_5_3, V1_5_2)
        version(V1_5_4, V1_5_3)
        version(V1_5_5, V1_5_4)
        version(V1_6_0, V1_5_2, embeddedObject(
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
        version(V1_6_1, V1_6_0)
        version(V1_6_2, V1_6_1)
        version(V1_6_3, V1_6_2)
        version(V1_7_0, V1_6_0, resourceObject(
                "SerialNumber" to null,
                "PartNumber" to null,
                "Version" to null,
                "FirmwareVersion" to ""
        ))
        version(V1_7_1, V1_7_0)
        version(V1_7_2, V1_7_1)
        version(V1_8_0, V1_7_1, resourceObject(
                "SystemInterface" to embeddedObject(PROCESSOR_INTERFACE),
                "OperatingSpeedMHz" to null,
                "MinSpeedMHz" to null
        ))
        version(V1_8_1, V1_8_0)
        version(V1_9_0, V1_8_1, resourceObject(
                "TurboState" to null,
                "BaseSpeedPriorityState" to null,
                "HighSpeedCoreIDs" to embeddedArray(),
                "OperatingConfigs" to LinkableResource(OPERATING_CONFIGS_COLLECTION),
                "AppliedOperatingConfig" to LinkableResource(OPERATING_CONFIG)
        ))
        version(V1_10_0, V1_9_0, resourceObject(
                "LocationIndicatorActive" to null,
                "BaseSpeedMHz" to null,
                "SpeedLimitMHz" to null,
                "SpeedLocked" to null
        ))
    }
}
