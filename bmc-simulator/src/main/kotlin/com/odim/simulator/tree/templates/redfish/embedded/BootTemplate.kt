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
import com.odim.simulator.tree.RedfishVersion.V1_0_0
import com.odim.simulator.tree.RedfishVersion.V1_11_0
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.BOOT
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.BOOT_OPTIONS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.CERTIFICATES_COLLECTION

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(BOOT)
open class BootTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, embeddedObject(
                "BootSourceOverrideTarget" to null,
                "BootSourceOverrideEnabled" to null,
                "UefiTargetBootSourceOverride" to null
        ))
        version(V1_1_0, V1_0_0, embeddedObject(
                "BootSourceOverrideMode" to null
        ))
        version(V1_5_0, V1_1_0, embeddedObject(
                "BootOptions" to ResourceCollection(BOOT_OPTIONS_COLLECTION),
                "BootNext" to null,
                "BootOrder" to embeddedArray()
        ))
        version(V1_6_0, V1_5_0, embeddedObject(
                "AliasBootOrder" to embeddedArray(),
                "BootOrderPropertySelection" to null
        ))
        version(V1_7_0, V1_6_0, embeddedObject(
                "Certificates" to ResourceCollection(CERTIFICATES_COLLECTION)
        ))
        version(V1_9_0, V1_7_0, embeddedObject(
                "HttpBootUri" to null
        ))
        version(V1_11_0, V1_9_0, embeddedObject(
                "AutomaticRetryConfig" to null,
                "AutomaticRetryAttempts" to null,
                "RemainingAutomaticRetryAttempts" to null
        ))
    }
}
