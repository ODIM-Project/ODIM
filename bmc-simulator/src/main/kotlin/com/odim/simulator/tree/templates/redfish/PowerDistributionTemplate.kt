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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.TRANSFER_CONTROL
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CIRCUITS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.OUTLETS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.OUTLET_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.SENSORS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.FACILITY
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.POWER_DISTRIBUTION
import com.odim.simulator.tree.structure.ResourceType.POWER_DISTRIBUTION_METRICS
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(POWER_DISTRIBUTION)
open class PowerDistributionTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Power Distribution Description",
                "Name" to "Power Distribution",
                "EquipmentType" to "RackPDU",
                "Model" to null,
                "Manufacturer" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "Version" to null,
                "FirmwareVersion" to "",
                "ProductionDate" to null,
                "AssetTag" to null,
                "UUID" to null,
                "Location" to embeddedObject(RESOURCE_LOCATION),
                "TransferConfiguration" to null,
                "TransferCriteria" to null,
                "Sensors" to ResourceCollection(SENSORS_COLLECTION),
                "Status" to embeddedObject(STATUS),
                "Mains" to LinkableResource(CIRCUITS_COLLECTION),
                "Branches" to LinkableResource(CIRCUITS_COLLECTION),
                "Feeders" to LinkableResource(CIRCUITS_COLLECTION),
                "Subfeeds" to LinkableResource(CIRCUITS_COLLECTION),
                "Outlets" to LinkableResource(OUTLETS_COLLECTION),
                "OutletGroups" to LinkableResource(OUTLET_GROUPS_COLLECTION),
                "Metrics" to SingletonResource(POWER_DISTRIBUTION_METRICS),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Chassis" to LinkableResourceArray(CHASSIS),
                        "Facility" to LinkableResource(FACILITY),
                        "ManagedBy" to LinkableResourceArray(MANAGER)
                ),
                "Actions" to Actions(
                        Action(TRANSFER_CONTROL)
                )
        ))
        version(V1_0_1, V1_0_0)
    }
}
