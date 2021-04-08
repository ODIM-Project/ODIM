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
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CAPACITY
import com.odim.simulator.tree.structure.EmbeddedObjectType.IMPORTED_SHARE
import com.odim.simulator.tree.structure.EmbeddedObjectType.IO_STATISTICS
import com.odim.simulator.tree.structure.EmbeddedObjectType.REPLICA_INFO
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.FILE_SHARES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.CAPACITY_SOURCE
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.FILE_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.SPARE_RESOURCE_SET

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(FILE_SYSTEM)
open class FileSystemTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "File System Description",
                "Name" to "File System",
                "BlockSizeBytes" to null,
                "Capacity" to embeddedObject(CAPACITY),
                "RemainingCapacity" to embeddedObject(CAPACITY),
                "CapacitySources" to LinkableResourceArray(CAPACITY_SOURCE),
                "LowSpaceWarningThresholdPercents" to embeddedArray(),
                "AccessCapabilities" to embeddedArray(),
                "CaseSensitive" to null,
                "CasePreserved" to null,
                "CharacterCodeSet" to embeddedArray(),
                "MaxFileNameLengthBytes" to null,
                "ClusterSizeBytes" to null,
                "ReplicaInfo" to embeddedObject(REPLICA_INFO),
                "ExportedShares" to ResourceCollection(FILE_SHARES_COLLECTION),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ReplicaCollection" to LinkableResourceArray(FILE_SYSTEM),
                        "ClassOfService" to LinkableResource(CLASS_OF_SERVICE)
                )
        ))
        version(V1_0_1, V1_0_0, resourceObject(
                "ImportedShares" to embeddedArray(IMPORTED_SHARE)
        ))
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_1_0, V1_0_2, resourceObject(
                "RemainingCapacityPercent" to null,
                "Actions" to Actions()
        ))
        version(V1_1_1, V1_1_0, resourceObject(
                "Identifiers" to embeddedArray(RESOURCE_IDENTIFIER)
        ))
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_2_0, V1_1_3, resourceObject(
                "IOStatistics" to embeddedObject(IO_STATISTICS),
                "RecoverableCapacitySourceCount" to null,
                "Links" to embeddedObject(
                        "SpareResourceSets" to LinkableResourceArray(SPARE_RESOURCE_SET)
                )
        ))
        version(V1_2_1, V1_2_0, resourceObject(
                "ReplicaTargets" to LinkableResourceArray(ANY)
        ))
        version(V1_2_2, V1_2_1)
    }
}
