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
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.LOG_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.PORTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.REDUNDANCY
import com.odim.simulator.tree.structure.ResourceType.SWITCH

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(SWITCH)
open class SwitchTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Switch Description",
                "Name" to "Switch",
                "SwitchType" to null,
                "Status" to embeddedObject(STATUS),
                "Manufacturer" to null,
                "Model" to null,
                "SKU" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "AssetTag" to null,
                "DomainID" to null,
                "IsManaged" to null,
                "TotalSwitchWidth" to null,
                "IndicatorLED" to null,
                "PowerState" to null,
                "Ports" to ResourceCollection(PORTS_COLLECTION),
                "Redundancy" to EmbeddedResourceArray(REDUNDANCY),
                "LogServices" to ResourceCollection(LOG_SERVICES_COLLECTION),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Chassis" to LinkableResource(CHASSIS),
                        "ManagedBy" to LinkableResourceArray(MANAGER)
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
        version(V1_1_0, V1_0_3, resourceObject(
                "Location" to embeddedObject(RESOURCE_LOCATION)
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_2_0, V1_1_3, resourceObject(
                "FirmwareVersion" to null
        ))
        version(V1_2_1, V1_2_0)
        version(V1_3_0, V1_2_0, resourceObject(
                "SupportedProtocols" to embeddedArray(),
                "UUID" to null,
                "Links" to embeddedObject(
                        "Endpoints" to LinkableResourceArray(ENDPOINT)
                )
        ))
        version(V1_3_1, V1_3_0)
        version(V1_4_0, V1_3_1, resourceObject(
                "LocationIndicatorActive" to null,
                "CurrentBandwidthGbps" to null,
                "MaxBandwidthGbps" to null,
                "Links" to embeddedObject(
                        "PCIeDevice" to LinkableResource(PCIE_DEVICE)
                )
        ))
    }
}
