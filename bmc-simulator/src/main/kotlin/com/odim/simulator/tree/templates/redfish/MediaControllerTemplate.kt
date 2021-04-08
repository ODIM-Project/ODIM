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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.PORTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.MEDIA_CONTROLLER
import com.odim.simulator.tree.structure.ResourceType.MEMORY_DOMAIN

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(MEDIA_CONTROLLER)
open class MediaControllerTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Media Controller Description",
                "Name" to "Media Controller",
                "Manufacturer" to null,
                "Model" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "Status" to embeddedObject(STATUS),
                "Ports" to ResourceCollection(PORTS_COLLECTION),
                "MediaControllerType" to null,
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Endpoints" to LinkableResourceArray(ENDPOINT),
                        "MemoryDomains" to LinkableResourceArray(MEMORY_DOMAIN)
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
        version(V1_1_0, V1_0_1, resourceObject(
                "UUID" to null
        ))
    }
}
