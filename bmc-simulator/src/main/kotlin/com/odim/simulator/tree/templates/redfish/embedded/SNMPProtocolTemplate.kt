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
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.SNMP_COMMUNITY
import com.odim.simulator.tree.structure.EmbeddedObjectType.SNMP_PROTOCOL
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(SNMP_PROTOCOL)
open class SNMPProtocolTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, embeddedObject(
                "ProtocolEnabled" to null,
                "Port" to null
        ))
        version(V1_5_0, V1_0_0, embeddedObject(
                "EnableSNMPv1" to null,
                "EnableSNMPv2c" to null,
                "EnableSNMPv3" to null,
                "CommunityStrings" to embeddedArray(SNMP_COMMUNITY),
                "CommunityAccessMode" to null,
                "EngineId" to null,
                "HideCommunityStrings" to null,
                "AuthenticationProtocol" to null,
                "EncryptionProtocol" to null
        ))
    }
}
