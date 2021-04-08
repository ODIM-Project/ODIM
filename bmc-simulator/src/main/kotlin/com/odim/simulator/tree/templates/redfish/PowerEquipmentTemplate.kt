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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollectionType.POWER_DISTRIBUTIONS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.POWER_EQUIPMENT

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(POWER_EQUIPMENT)
open class PowerEquipmentTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Power Equipment Description",
                "Name" to "Power Equipment",
                "Status" to embeddedObject(STATUS),
                "FloorPDUs" to LinkableResource(POWER_DISTRIBUTIONS_COLLECTION),
                "RackPDUs" to LinkableResource(POWER_DISTRIBUTIONS_COLLECTION),
                "Switchgear" to LinkableResource(POWER_DISTRIBUTIONS_COLLECTION),
                "TransferSwitches" to LinkableResource(POWER_DISTRIBUTIONS_COLLECTION),
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ManagedBy" to LinkableResourceArray(MANAGER)
                )
        ))
    }
}
