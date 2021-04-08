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
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_0_8
import com.odim.simulator.tree.RedfishVersion.V1_0_9
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_1_7
import com.odim.simulator.tree.RedfishVersion.V1_1_8
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_2_6
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_3_4
import com.odim.simulator.tree.RedfishVersion.V1_3_5
import com.odim.simulator.tree.RedfishVersion.V1_3_6
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_4_3
import com.odim.simulator.tree.RedfishVersion.V1_4_4
import com.odim.simulator.tree.RedfishVersion.V1_4_5
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_5_1
import com.odim.simulator.tree.RedfishVersion.V1_5_2
import com.odim.simulator.tree.RedfishVersion.V1_5_3
import com.odim.simulator.tree.RedfishVersion.V1_5_4
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.RedfishVersion.V1_6_3
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_7_1
import com.odim.simulator.tree.RedfishVersion.V1_7_2
import com.odim.simulator.tree.RedfishVersion.V1_8_0
import com.odim.simulator.tree.RedfishVersion.V1_8_1
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.SET_ENCRYPTION_KEY
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CONSISTENCY_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.ENDPOINT_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.FILE_SYSTEMS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_CONTROLLERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_POOLS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.VOLUMES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.REDUNDANCY
import com.odim.simulator.tree.structure.ResourceType.SIMPLE_STORAGE
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.STORAGE_CONTROLLER
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(STORAGE)
open class StorageTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Storage Description",
                "Name" to "Storage",
                "StorageControllers" to EmbeddedResourceArray(STORAGE_CONTROLLER),
                "Drives" to LinkableResourceArray(DRIVE),
                "Volumes" to ResourceCollection(VOLUMES_COLLECTION),
                "Status" to embeddedObject(STATUS),
                "Redundancy" to EmbeddedResourceArray(REDUNDANCY),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Enclosures" to LinkableResourceArray(CHASSIS)
                ),
                "Actions" to Actions(
                        Action(SET_ENCRYPTION_KEY, "EncryptionKey", mutableListOf())
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_0_9, V1_0_8)
        version(V1_1_0, V1_0_1)
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_1_7, V1_1_6)
        version(V1_1_8, V1_1_7)
        version(V1_2_0, V1_1_2)
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_2_6, V1_2_5)
        version(V1_3_0, V1_2_0)
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_3_4, V1_3_3)
        version(V1_3_5, V1_3_4)
        version(V1_3_6, V1_3_5)
        version(V1_4_0, V1_3_1)
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_4_3, V1_4_2)
        version(V1_4_4, V1_4_3)
        version(V1_4_5, V1_4_4)
        version(V1_5_0, V1_4_1)
        version(V1_5_1, V1_5_0)
        version(V1_5_2, V1_5_1)
        version(V1_5_3, V1_5_2)
        version(V1_5_4, V1_5_3)
        version(V1_6_0, V1_5_1)
        version(V1_6_1, V1_6_0)
        version(V1_6_2, V1_6_1)
        version(V1_6_3, V1_6_2)
        version(V1_7_0, V1_6_1)
        version(V1_7_1, V1_7_0)
        version(V1_7_2, V1_7_1)
        version(V1_8_0, V1_7_1, resourceObject(
                "FileSystems" to ResourceCollection(FILE_SYSTEMS_COLLECTION),
                "StoragePools" to ResourceCollection(STORAGE_POOLS_COLLECTION),
                "StorageGroups" to ResourceCollection(STORAGE_GROUPS_COLLECTION),
                "EndpointGroups" to ResourceCollection(ENDPOINT_GROUPS_COLLECTION),
                "ConsistencyGroups" to LinkableResource(CONSISTENCY_GROUPS_COLLECTION)
        ))
        version(V1_8_1, V1_8_0)
        version(V1_9_0, V1_8_1, resourceObject(
                "Controllers" to ResourceCollection(STORAGE_CONTROLLERS_COLLECTION),
                "Identifiers" to embeddedArray(RESOURCE_IDENTIFIER),
                "Links" to embeddedObject(
                        "SimpleStorage" to LinkableResource(SIMPLE_STORAGE),
                        "StorageServices" to LinkableResourceArray(STORAGE_SERVICE)
                )
        ))
    }
}
