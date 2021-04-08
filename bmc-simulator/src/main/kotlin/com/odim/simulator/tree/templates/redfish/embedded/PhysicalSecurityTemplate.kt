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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.PHYSICAL_SECURITY
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(PHYSICAL_SECURITY)
open class PhysicalSecurityTemplate : ResourceTemplate() {
    init {
        version(V1_1_0, embeddedObject(
                "IntrusionSensorNumber" to null,
                "IntrusionSensor" to null,
                "IntrusionSensorReArm" to null
        ))
    }
}
