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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.EXPOSE_VOLUMES
import com.odim.simulator.tree.structure.ActionType.HIDE_VOLUMES
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CHAP_INFORMATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.DHCHAP_INFORMATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.MAPPED_VOLUME
import com.odim.simulator.tree.structure.EmbeddedObjectType.REPLICA_INFO
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT_GROUP
import com.odim.simulator.tree.structure.ResourceType.STORAGE_GROUP
import com.odim.simulator.tree.structure.ResourceType.VOLUME

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(STORAGE_GROUP)
open class StorageGroupTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Storage Group Description",
                "Name" to "Storage Group",
                "Identifier" to embeddedObject(RESOURCE_IDENTIFIER),
                "AccessState" to null,
                "MembersAreConsistent" to false,
                "VolumesAreExposed" to false,
                "Status" to embeddedObject(STATUS),
                "ReplicaInfo" to embeddedObject(REPLICA_INFO),
                "ClientEndpointGroups" to LinkableResourceArray(ENDPOINT_GROUP),
                "ServerEndpointGroups" to LinkableResourceArray(ENDPOINT_GROUP),
                "Volumes" to LinkableResourceArray(VOLUME),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ParentStorageGroups" to LinkableResourceArray(STORAGE_GROUP),
                        "ChildStorageGroups" to LinkableResourceArray(STORAGE_GROUP),
                        "ClassOfService" to LinkableResource(CLASS_OF_SERVICE)
                ),
                "Actions" to Actions(
                        Action(EXPOSE_VOLUMES),
                        Action(HIDE_VOLUMES)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_1_0, V1_0_2, resourceObject(
                "MappedVolumes" to embeddedArray(MAPPED_VOLUME)
        ))
        version(V1_1_1, V1_1_0, resourceObject(
                "ReplicaTargets" to LinkableResourceArray(ANY)
        ))
        version(V1_1_2, V1_1_1)
        version(V1_2_0, V1_1_1, resourceObject(
                "AuthenticationMethod" to null,
                "ChapInfo" to embeddedArray(CHAP_INFORMATION)
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_3_0, V1_2_2, resourceObject(
                "DHChapInfo" to embeddedArray(DHCHAP_INFORMATION)
        ))
        version(V1_4_0, V1_3_0)
    }
}
