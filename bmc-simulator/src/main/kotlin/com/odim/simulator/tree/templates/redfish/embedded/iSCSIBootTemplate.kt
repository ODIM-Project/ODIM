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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.ISCSI_BOOT
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(ISCSI_BOOT)
open class iSCSIBootTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, embeddedObject(
                "IPAddressType" to null,
                "InitiatorIPAddress" to null,
                "InitiatorName" to null,
                "InitiatorDefaultGateway" to null,
                "InitiatorNetmask" to null,
                "TargetInfoViaDHCP" to null,
                "PrimaryTargetName" to null,
                "PrimaryTargetIPAddress" to null,
                "PrimaryTargetTCPPort" to null,
                "PrimaryLUN" to null,
                "PrimaryVLANEnable" to null,
                "PrimaryVLANId" to null,
                "PrimaryDNS" to null,
                "SecondaryTargetName" to null,
                "SecondaryTargetIPAddress" to null,
                "SecondaryTargetTCPPort" to null,
                "SecondaryLUN" to null,
                "SecondaryVLANEnable" to null,
                "SecondaryVLANId" to null,
                "SecondaryDNS" to null,
                "IPMaskDNSViaDHCP" to null,
                "RouterAdvertisementEnabled" to null,
                "AuthenticationMethod" to null,
                "CHAPUsername" to null,
                "CHAPSecret" to null,
                "MutualCHAPUsername" to null,
                "MutualCHAPSecret" to null
        ))
    }
}
