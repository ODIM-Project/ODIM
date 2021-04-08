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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.RESOURCE_BLOCKS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.ZONES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.COMPOSITION_SERVICE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(COMPOSITION_SERVICE)
open class CompositionServiceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to "CompositionService",
                "Description" to "Composition Service Description",
                "Name" to "Composition Service",
                "Status" to embeddedObject(STATUS),
                "ServiceEnabled" to null,
                "Actions" to Actions(),
                "ResourceBlocks" to ResourceCollection(RESOURCE_BLOCKS_COLLECTION),
                "ResourceZones" to ResourceCollection(ZONES_COLLECTION)
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_1_0, V1_0_2, resourceObject(
                "AllowOverprovisioning" to null,
                "AllowZoneAffinity" to null
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
    }
}
