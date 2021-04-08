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
package com.odim.simulator.tree.templates.redfish.embedded

import com.odim.simulator.tree.EmbeddedObjectTemplate
import com.odim.simulator.tree.RedfishVersion.V1_0_0
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.CONNECTED_ENTITY
import com.odim.simulator.tree.structure.EmbeddedObjectType.PCI_ID
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.ResourceType.ANY

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(CONNECTED_ENTITY)
open class ConnectedEntityTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, embeddedObject(
                "EntityType" to null,
                "EntityRole" to null,
                "EntityPciId" to embeddedObject(PCI_ID),
                "PciFunctionNumber" to null,
                "PciClassCode" to null,
                "Identifiers" to embeddedArray(RESOURCE_IDENTIFIER),
                "EntityLink" to LinkableResource(ANY),
                "Oem" to embeddedObject()
        ))
        version(V1_4_0, V1_0_0, embeddedObject(
                "GenZ" to null
        ))
    }
}
