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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.GENERATE_CSR
import com.odim.simulator.tree.structure.ActionType.REPLACE_CERTIFICATE
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.CERTIFICATE_LOCATIONS
import com.odim.simulator.tree.structure.ResourceType.CERTIFICATE_SERVICE
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(CERTIFICATE_SERVICE)
open class CertificateServiceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to "CertificateService",
                "Description" to "Certificate Service Description",
                "Name" to "Certificate Service",
                "CertificateLocations" to SingletonResource(CERTIFICATE_LOCATIONS),
                "Actions" to Actions(
                        Action(GENERATE_CSR),
                        Action(REPLACE_CERTIFICATE)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
    }
}
