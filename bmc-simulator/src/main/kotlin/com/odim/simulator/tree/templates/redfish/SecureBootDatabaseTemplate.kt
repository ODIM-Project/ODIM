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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET_KEYS
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CERTIFICATES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.SIGNATURES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.SECURE_BOOT_DATABASE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(SECURE_BOOT_DATABASE)
open class SecureBootDatabaseTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Secure Boot Database Description",
                "Name" to "Secure Boot Database",
                "DatabaseId" to "",
                "Certificates" to ResourceCollection(CERTIFICATES_COLLECTION),
                "Signatures" to LinkableResource(SIGNATURES_COLLECTION),
                "Actions" to Actions(
                        Action(RESET_KEYS, "ResetKeysType", mutableListOf())
                )
        ))
    }
}
