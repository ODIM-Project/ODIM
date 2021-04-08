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
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.ActionType.SET_DEFAULT_BOOT_ORDER
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollectionType.AGGREGATES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.AGGREGATION_SOURCES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.CONNECTION_METHODS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.AGGREGATION_SERVICE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(AGGREGATION_SERVICE)
open class AggregationServiceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Aggregation Service Description",
                "Name" to "Aggregation Service",
                "ServiceEnabled" to null,
                "Status" to embeddedObject(STATUS),
                "Aggregates" to LinkableResource(AGGREGATES_COLLECTION),
                "AggregationSources" to LinkableResource(AGGREGATION_SOURCES_COLLECTION),
                "ConnectionMethods" to LinkableResource(CONNECTION_METHODS_COLLECTION),
                "Actions" to Actions(
                        Action(RESET),
                        Action(SET_DEFAULT_BOOT_ORDER, "Systems", mutableListOf())
                )
        ))
    }
}
