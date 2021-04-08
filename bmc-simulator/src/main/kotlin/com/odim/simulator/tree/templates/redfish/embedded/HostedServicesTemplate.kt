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
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.HOSTED_SERVICES
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.HOSTED_STORAGE_SERVICES

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(HOSTED_SERVICES)
open class HostedServicesTemplate : ResourceTemplate() {
    init {
        version(V1_2_0, embeddedObject(
                "StorageServices" to ResourceCollection(HOSTED_STORAGE_SERVICES),
                "Oem" to embeddedObject()
        ))
    }
}
