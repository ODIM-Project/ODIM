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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.ETHERNET_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.FILE_SHARE
import com.odim.simulator.tree.structure.ResourceType.FILE_SYSTEM

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(FILE_SHARE)
open class FileShareTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "File Share Description",
                "Name" to "File Share",
                "FileSharePath" to null,
                "FileSharingProtocols" to embeddedArray(),
                "Status" to embeddedObject(STATUS),
                "DefaultAccessCapabilities" to embeddedArray(),
                "ExecuteSupport" to false,
                "RootAccess" to false,
                "WritePolicy" to null,
                "CASupported" to false,
                "FileShareTotalQuotaBytes" to null,
                "FileShareRemainingQuotaBytes" to null,
                "LowSpaceWarningThresholdPercents" to embeddedArray(),
                "FileShareQuotaType" to null,
                "EthernetInterfaces" to ResourceCollection(ETHERNET_INTERFACES_COLLECTION),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ClassOfService" to LinkableResource(CLASS_OF_SERVICE),
                        "FileSystem" to LinkableResource(FILE_SYSTEM)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_1_0, V1_0_0, resourceObject(
                "RemainingCapacityPercent" to null,
                "Actions" to Actions()
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
    }
}
