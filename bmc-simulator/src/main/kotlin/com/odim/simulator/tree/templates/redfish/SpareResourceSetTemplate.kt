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
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.SPARE_RESOURCE_SET

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(SPARE_RESOURCE_SET)
open class SpareResourceSetTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Spare Resource Set Description",
                "Name" to "Spare Resource Set",
                "ResourceType" to null,
                "TimeToProvision" to null,
                "TimeToReplenish" to null,
                "OnHandLocation" to embeddedObject(RESOURCE_LOCATION),
                "OnLine" to null,
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "OnHandSpares" to LinkableResourceArray(ANY),
                        "ReplacementSpareSets" to LinkableResourceArray(SPARE_RESOURCE_SET)
                )
        ))
        version(V1_0_1, V1_0_0, resourceObject(
                "Actions" to Actions()
        ))
    }
}
