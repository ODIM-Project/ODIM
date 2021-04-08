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
import com.odim.simulator.tree.structure.ActionType.POWER_CONTROL
import com.odim.simulator.tree.structure.ActionType.RESET_METRICS
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.OUTLET
import com.odim.simulator.tree.structure.ResourceType.OUTLET_GROUP
import com.odim.simulator.tree.structure.ResourceType.SENSOR

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(OUTLET_GROUP)
open class OutletGroupTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Outlet Group Description",
                "Name" to "Outlet Group",
                "Status" to embeddedObject(STATUS),
                "CreatedBy" to null,
                "PowerOnDelaySeconds" to null,
                "PowerOffDelaySeconds" to null,
                "PowerCycleDelaySeconds" to null,
                "PowerRestoreDelaySeconds" to null,
                "PowerRestorePolicy" to "AlwaysOn",
                "PowerState" to null,
                "PowerEnabled" to null,
                "PowerWatts" to LinkableResource(SENSOR),
                "EnergykWh" to LinkableResource(SENSOR),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Outlets" to LinkableResourceArray(OUTLET)
                ),
                "Actions" to Actions(
                        Action(POWER_CONTROL, "PowerState", mutableListOf(
                                "On",
                                "Off"
                        )),
                        Action(RESET_METRICS)
                )
        ))
        version(V1_0_1, V1_0_0)
    }
}
