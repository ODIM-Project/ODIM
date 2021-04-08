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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CAPACITY
import com.odim.simulator.tree.structure.EmbeddedObjectType.IO_STATISTICS
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CLASS_OF_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_POOLS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.VOLUMES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CAPACITY_SOURCE
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.SPARE_RESOURCE_SET
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.STORAGE_POOL

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(STORAGE_POOL)
open class StoragePoolTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Storage Pool Description",
                "Name" to "Storage Pool",
                "Identifier" to embeddedObject(RESOURCE_IDENTIFIER),
                "BlockSizeBytes" to null,
                "Capacity" to embeddedObject(CAPACITY),
                "CapacitySources" to LinkableResourceArray(CAPACITY_SOURCE),
                "LowSpaceWarningThresholdPercents" to embeddedArray(),
                "AllocatedVolumes" to ResourceCollection(VOLUMES_COLLECTION),
                "AllocatedPools" to ResourceCollection(STORAGE_POOLS_COLLECTION),
                "ClassesOfService" to ResourceCollection(CLASS_OF_SERVICES_COLLECTION),
                "Status" to embeddedObject(STATUS),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "DefaultClassOfService" to LinkableResource(CLASS_OF_SERVICE)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_1_0, V1_0_1, resourceObject(
                "RemainingCapacityPercent" to null
        ))
        version(V1_1_1, V1_1_0, resourceObject(
                "MaxBlockSizeBytes" to null
        ))
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_2_0, V1_1_1, resourceObject(
                "IOStatistics" to embeddedObject(IO_STATISTICS),
                "RecoverableCapacitySourceCount" to null,
                "DefaultClassOfService" to LinkableResource(CLASS_OF_SERVICE),
                "Links" to embeddedObject(
                        "DedicatedSpareDrives" to LinkableResourceArray(DRIVE),
                        "SpareResourceSets" to LinkableResourceArray(SPARE_RESOURCE_SET)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_3_0, V1_2_1, resourceObject(
                "SupportedRAIDTypes" to embeddedArray(),
                "SupportedProvisioningPolicies" to embeddedArray(),
                "Deduplicated" to null,
                "Compressed" to null,
                "Encrypted" to null,
                "Actions" to Actions()
        ))
        version(V1_3_1, V1_3_0)
        version(V1_4_0, V1_3_1, resourceObject(
                "NVMeSetProperties" to null,
                "NVMeEnduranceGroupProperties" to null,
                "Links" to embeddedObject(
                        "OwningStorageResource" to LinkableResource(STORAGE)
                )
        ))
        version(V1_5_0, V1_4_0)
    }
}
