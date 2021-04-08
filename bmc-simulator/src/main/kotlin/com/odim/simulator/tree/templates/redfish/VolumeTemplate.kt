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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.ASSIGN_REPLICA_TARGET
import com.odim.simulator.tree.structure.ActionType.CHANGE_RAID_LAYOUT
import com.odim.simulator.tree.structure.ActionType.CHECK_CONSISTENCY
import com.odim.simulator.tree.structure.ActionType.CREATE_REPLICA_TARGET
import com.odim.simulator.tree.structure.ActionType.FORCE_ENABLE
import com.odim.simulator.tree.structure.ActionType.INITIALIZE
import com.odim.simulator.tree.structure.ActionType.REMOVE_REPLICA_RELATIONSHIP
import com.odim.simulator.tree.structure.ActionType.RESUME_REPLICATION
import com.odim.simulator.tree.structure.ActionType.REVERSE_REPLICATION_RELATIONSHIP
import com.odim.simulator.tree.structure.ActionType.SPLIT_REPLICATION
import com.odim.simulator.tree.structure.ActionType.SUSPEND_REPLICATION
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CAPACITY
import com.odim.simulator.tree.structure.EmbeddedObjectType.IO_STATISTICS
import com.odim.simulator.tree.structure.EmbeddedObjectType.OPERATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.REPLICA_INFO
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_POOLS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.CAPACITY_SOURCE
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.CONSISTENCY_GROUP
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.SPARE_RESOURCE_SET
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.STORAGE_GROUP
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE
import com.odim.simulator.tree.structure.ResourceType.VOLUME

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(VOLUME)
open class VolumeTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Volume Description",
                "Name" to "Volume",
                "Status" to embeddedObject(STATUS),
                "CapacityBytes" to null,
                "VolumeType" to null,
                "Encrypted" to null,
                "EncryptionTypes" to embeddedArray(),
                "Identifiers" to embeddedArray(RESOURCE_IDENTIFIER),
                "BlockSizeBytes" to null,
                "Operations" to embeddedArray(OPERATION),
                "OptimumIOSizeBytes" to null,
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Drives" to LinkableResourceArray(DRIVE)
                ),
                "Actions" to Actions(
                        Action(CHECK_CONSISTENCY)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_1_0, V1_0_0, resourceObject(
                "AccessCapabilities" to embeddedArray(),
                "MaxBlockSizeBytes" to null,
                "Capacity" to embeddedObject(CAPACITY),
                "CapacitySources" to LinkableResourceArray(CAPACITY_SOURCE),
                "LowSpaceWarningThresholdPercents" to embeddedArray(),
                "Manufacturer" to null,
                "Model" to null,
                "ReplicaInfo" to embeddedObject(REPLICA_INFO),
                "StorageGroups" to ResourceCollection(STORAGE_GROUPS_COLLECTION),
                "AllocatedPools" to ResourceCollection(STORAGE_POOLS_COLLECTION),
                "Links" to embeddedObject(
                        "ClassOfService" to LinkableResource(CLASS_OF_SERVICE)
                )
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_2_0, V1_1_1, resourceObject(
                "IOStatistics" to embeddedObject(IO_STATISTICS),
                "RemainingCapacityPercent" to null,
                "Links" to embeddedObject(
                        "DedicatedSpareDrives" to LinkableResourceArray(DRIVE)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_3_0, V1_2_1, resourceObject(
                "RecoverableCapacitySourceCount" to null,
                "ReplicaTargets" to LinkableResourceArray(ANY),
                "Links" to embeddedObject(
                        "SpareResourceSets" to LinkableResourceArray(SPARE_RESOURCE_SET)
                )
        ))
        version(V1_3_1, V1_3_0, resourceObject(
                "RAIDType" to null
        ))
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_4_0, V1_3_1, resourceObject(
                "ProvisioningPolicy" to null,
                "StripSizeBytes" to null,
                "ReadCachePolicy" to null,
                "VolumeUsage" to null,
                "WriteCachePolicy" to null,
                "WriteCacheState" to null,
                "LogicalUnitNumber" to null,
                "MediaSpanCount" to null,
                "DisplayName" to null,
                "WriteHoleProtectionPolicy" to "Off",
                "Deduplicated" to null,
                "Compressed" to null,
                "Links" to embeddedObject(
                        "ClientEndpoints" to LinkableResourceArray(ENDPOINT),
                        "ServerEndpoints" to LinkableResourceArray(ENDPOINT),
                        "StorageGroups" to LinkableResourceArray(STORAGE_GROUP),
                        "ConsistencyGroups" to LinkableResourceArray(CONSISTENCY_GROUP),
                        "OwningStorageService" to LinkableResource(STORAGE_SERVICE)
                ),
                "Actions" to Actions(
                        Action(ASSIGN_REPLICA_TARGET),
                        Action(CREATE_REPLICA_TARGET),
                        Action(REMOVE_REPLICA_RELATIONSHIP),
                        Action(RESUME_REPLICATION, "TargetVolume", mutableListOf()),
                        Action(REVERSE_REPLICATION_RELATIONSHIP, "TargetVolume", mutableListOf()),
                        Action(SPLIT_REPLICATION, "TargetVolume", mutableListOf()),
                        Action(SUSPEND_REPLICATION, "TargetVolume", mutableListOf())
                )
        ))
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_5_0, V1_4_2, resourceObject(
                "IOPerfModeEnabled" to null,
                "NVMeNamespaceProperties" to null,
                "Links" to embeddedObject(
                        "JournalingMedia" to LinkableResource(ANY),
                        "OwningStorageResource" to LinkableResource(STORAGE)
                ),
                "Actions" to Actions(
                        Action(INITIALIZE),
                        Action(CHANGE_RAID_LAYOUT),
                        Action(FORCE_ENABLE)
                )
        ))
    }
}
