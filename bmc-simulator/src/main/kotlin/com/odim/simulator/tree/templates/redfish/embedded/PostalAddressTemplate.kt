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
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.POSTAL_ADDRESS
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(POSTAL_ADDRESS)
open class PostalAddressTemplate : ResourceTemplate() {
    init {
        version(V1_3_0, embeddedObject(
                "Country" to null,
                "Territory" to null,
                "District" to null,
                "City" to null,
                "Division" to null,
                "Neighborhood" to null,
                "LeadingStreetDirection" to null,
                "Street" to null,
                "TrailingStreetSuffix" to null,
                "StreetSuffix" to null,
                "HouseNumber" to null,
                "HouseNumberSuffix" to null,
                "Landmark" to null,
                "Location" to null,
                "Floor" to null,
                "Name" to null,
                "PostalCode" to null,
                "Building" to null,
                "Unit" to null,
                "Room" to null,
                "Seat" to null,
                "PlaceType" to null,
                "Community" to null,
                "POBox" to null,
                "AdditionalCode" to null,
                "Road" to null,
                "RoadSection" to null,
                "RoadBranch" to null,
                "RoadSubBranch" to null,
                "RoadPreModifier" to null,
                "RoadPostModifier" to null,
                "GPSCoords" to null
        ))
        version(V1_7_0, V1_3_0, embeddedObject(
                "AdditionalInfo" to null
        ))
    }
}
