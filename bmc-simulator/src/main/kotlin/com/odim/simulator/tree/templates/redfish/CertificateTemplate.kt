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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.REKEY
import com.odim.simulator.tree.structure.ActionType.RENEW
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CERTIFICATE_IDENTIFIER
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.CERTIFICATE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(CERTIFICATE)
open class CertificateTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Certificate Description",
                "Name" to "Certificate",
                "CertificateString" to null,
                "CertificateType" to null,
                "Issuer" to embeddedObject(CERTIFICATE_IDENTIFIER),
                "Subject" to embeddedObject(CERTIFICATE_IDENTIFIER),
                "ValidNotBefore" to "2017-04-14T06:35:05Z",
                "ValidNotAfter" to "2017-04-14T06:35:05Z",
                "KeyUsage" to embeddedArray(),
                "Actions" to Actions()
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_1_0, V1_0_1, embeddedObject(
                "Actions" to Actions(
                        Action(REKEY),
                        Action(RENEW, "ChallengePassword", mutableListOf())
                )
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_2_0, V1_1_1, resourceObject(
                "UefiSignatureOwner" to null
        ))
        version(V1_2_1, V1_2_0)
    }
}
