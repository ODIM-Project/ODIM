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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.ACCELERATION_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(ACCELERATION_FUNCTION)
open class AccelerationFunctionTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Acceleration Function Description",
                "Name" to "Acceleration Function",
                "Status" to embeddedObject(STATUS),
                "UUID" to null,
                "FpgaReconfigurationSlots" to embeddedArray(),
                "AccelerationFunctionType" to null,
                "Manufacturer" to "",
                "Version" to "",
                "PowerWatts" to 0,
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Endpoints" to LinkableResourceArray(ENDPOINT),
                        "PCIeFunctions" to LinkableResourceArray(PCIE_FUNCTION)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
    }
}
