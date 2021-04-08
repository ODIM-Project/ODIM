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
import com.odim.simulator.tree.structure.ActionType.ADD_ELEMENTS
import com.odim.simulator.tree.structure.ActionType.REMOVE_ELEMENTS
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.ActionType.SET_DEFAULT_BOOT_ORDER
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.AGGREGATE
import com.odim.simulator.tree.structure.ResourceType.ANY

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(AGGREGATE)
open class AggregateTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Aggregate Description",
                "Name" to "Aggregate",
                "Elements" to LinkableResourceArray(ANY),
                "ElementsCount" to null,
                "Actions" to Actions(
                        Action(RESET),
                        Action(SET_DEFAULT_BOOT_ORDER),
                        Action(ADD_ELEMENTS, "Elements", mutableListOf()),
                        Action(REMOVE_ELEMENTS, "Elements", mutableListOf())
                )
        ))
    }
}
