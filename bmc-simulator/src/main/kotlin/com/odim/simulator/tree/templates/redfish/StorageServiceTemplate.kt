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
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.SET_ENCRYPTION_KEY
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.IO_STATISTICS
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CLASS_OF_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.CONSISTENCY_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.DRIVES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.ENDPOINTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.ENDPOINT_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.FILE_SYSTEMS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.LINES_OF_SERVICE_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_GROUPS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_POOLS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.VOLUMES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.DATA_PROTECTION_LOS_CAPABILITIES
import com.odim.simulator.tree.structure.ResourceType.DATA_SECURITY_LOS_CAPABILITIES
import com.odim.simulator.tree.structure.ResourceType.DATA_STORAGE_LOS_CAPABILITIES
import com.odim.simulator.tree.structure.ResourceType.IO_CONNECTIVITY_LOS_CAPABILITIES
import com.odim.simulator.tree.structure.ResourceType.IO_PERFORMANCE_LOS_CAPABILITIES
import com.odim.simulator.tree.structure.ResourceType.REDUNDANCY
import com.odim.simulator.tree.structure.ResourceType.SPARE_RESOURCE_SET
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(STORAGE_SERVICE)
open class StorageServiceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Storage Service Description",
                "Name" to "Storage Service",
                "Identifier" to embeddedObject(RESOURCE_IDENTIFIER),
                "Status" to embeddedObject(STATUS),
                "StorageGroups" to ResourceCollection(STORAGE_GROUPS_COLLECTION),
                "EndpointGroups" to ResourceCollection(ENDPOINT_GROUPS_COLLECTION),
                "ClientEndpointGroups" to ResourceCollection(ENDPOINT_GROUPS_COLLECTION),
                "ServerEndpointGroups" to ResourceCollection(ENDPOINT_GROUPS_COLLECTION),
                "Volumes" to ResourceCollection(VOLUMES_COLLECTION),
                "FileSystems" to ResourceCollection(FILE_SYSTEMS_COLLECTION),
                "StoragePools" to ResourceCollection(STORAGE_POOLS_COLLECTION),
                "Drives" to ResourceCollection(DRIVES_COLLECTION),
                "Endpoints" to ResourceCollection(ENDPOINTS_COLLECTION),
                "Redundancy" to EmbeddedResourceArray(REDUNDANCY),
                "ClassesOfService" to ResourceCollection(CLASS_OF_SERVICES_COLLECTION),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "HostingSystem" to LinkableResource(ANY),
                        "DefaultClassOfService" to LinkableResource(CLASS_OF_SERVICE),
                        "DataProtectionLoSCapabilities" to SingletonResource(DATA_PROTECTION_LOS_CAPABILITIES),
                        "DataSecurityLoSCapabilities" to SingletonResource(DATA_SECURITY_LOS_CAPABILITIES),
                        "DataStorageLoSCapabilities" to SingletonResource(DATA_STORAGE_LOS_CAPABILITIES),
                        "IOConnectivityLoSCapabilities" to SingletonResource(IO_CONNECTIVITY_LOS_CAPABILITIES),
                        "IOPerformanceLoSCapabilities" to SingletonResource(IO_PERFORMANCE_LOS_CAPABILITIES)
                ),
                "Actions" to Actions(
                        Action(SET_ENCRYPTION_KEY, mutableMapOf(
                                "Storage" to mutableListOf(),
                                "EncryptionKey" to mutableListOf()
                        )
                ))
        ))
        version(V1_0_1, V1_0_0, resourceObject(
                "StorageSubsystems" to ResourceCollection(STORAGES_COLLECTION)
        ))
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_1_0, V1_0_3)
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_2_0, V1_1_0, resourceObject(
                "IOStatistics" to embeddedObject(IO_STATISTICS),
                "SpareResourceSets" to LinkableResourceArray(SPARE_RESOURCE_SET),
                "DataProtectionLoSCapabilities" to SingletonResource(DATA_PROTECTION_LOS_CAPABILITIES),
                "DataSecurityLoSCapabilities" to SingletonResource(DATA_SECURITY_LOS_CAPABILITIES),
                "DataStorageLoSCapabilities" to SingletonResource(DATA_STORAGE_LOS_CAPABILITIES),
                "IOConnectivityLoSCapabilities" to SingletonResource(IO_CONNECTIVITY_LOS_CAPABILITIES),
                "IOPerformanceLoSCapabilities" to SingletonResource(IO_PERFORMANCE_LOS_CAPABILITIES),
                "DefaultClassOfService" to LinkableResource(CLASS_OF_SERVICE)
        ))
        version(V1_2_1, V1_2_0)
        version(V1_3_0, V1_2_1, resourceObject(
                "ConsistencyGroups" to LinkableResource(CONSISTENCY_GROUPS_COLLECTION)
        ))
        version(V1_4_0, V1_2_1, resourceObject(
                "LinesOfService" to LinkableResourceArray(LINES_OF_SERVICE_COLLECTION)
        ))
    }
}
