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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.EXTERNAL_ACCOUNT_PROVIDER_AUTHENTICATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.EXTERNAL_ACCOUNT_PROVIDER_LDAP_SERVICE
import com.odim.simulator.tree.structure.EmbeddedObjectType.EXTERNAL_ACCOUNT_PROVIDER_ROLE_MAPPING
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CERTIFICATES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.EXTERNAL_ACCOUNT_PROVIDER

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(EXTERNAL_ACCOUNT_PROVIDER)
open class ExternalAccountProviderTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "External Account Provider Description",
                "Name" to "External Account Provider",
                "AccountProviderType" to null,
                "ServiceEnabled" to null,
                "ServiceAddresses" to embeddedArray(),
                "Authentication" to embeddedObject(EXTERNAL_ACCOUNT_PROVIDER_AUTHENTICATION),
                "LDAPService" to embeddedObject(EXTERNAL_ACCOUNT_PROVIDER_LDAP_SERVICE),
                "RemoteRoleMapping" to embeddedArray(EXTERNAL_ACCOUNT_PROVIDER_ROLE_MAPPING),
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject()
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_1_0, V1_0_1, resourceObject(
                "Certificates" to ResourceCollection(CERTIFICATES_COLLECTION)
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
    }
}
