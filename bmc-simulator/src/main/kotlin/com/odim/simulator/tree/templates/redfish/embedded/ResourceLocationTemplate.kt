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
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.CONTACT_INFO
import com.odim.simulator.tree.structure.EmbeddedObjectType.PART_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.PLACEMENT
import com.odim.simulator.tree.structure.EmbeddedObjectType.POSTAL_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(RESOURCE_LOCATION)
open class ResourceLocationTemplate : ResourceTemplate() {
    init {
        version(V1_1_0, embeddedObject(
                "Info" to null,
                "InfoFormat" to null,
                "Oem" to embeddedObject()
        ))
        version(V1_3_0, V1_1_0, embeddedObject(
                "PostalAddress" to embeddedObject(POSTAL_ADDRESS),
                "Placement" to embeddedObject(PLACEMENT)
        ))
        version(V1_5_0, V1_3_0, embeddedObject(
                "PartLocation" to embeddedObject(PART_LOCATION)
        ))
        version(V1_6_0, V1_5_0, embeddedObject(
                "Longitude" to 0.00,
                "Latitude" to 0.00,
                "AltitudeMeters" to 0.00
        ))
        version(V1_7_0, V1_6_0, embeddedObject(
                "Contacts" to embeddedArray(CONTACT_INFO)
        ))
    }
}
