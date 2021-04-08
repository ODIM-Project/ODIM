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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.BGP_NEIGHBOR
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(BGP_NEIGHBOR)
open class BGPNeighborTemplate : ResourceTemplate() {
    init {
        version(V1_1_0, embeddedObject(
                "Address" to null,
                "AllowOwnASEnabled" to null,
                "ConnectRetrySeconds" to null,
                "HoldTimeSeconds" to null,
                "KeepaliveIntervalSeconds" to null,
                "MinimumAdvertisementIntervalSeconds" to null,
                "TCPMaxSegmentSizeBytes" to null,
                "PathMTUDiscoveryEnabled" to null,
                "PassiveModeEnabled" to null,
                "TreatAsWithdrawEnabled" to null,
                "ReplacePeerASEnabled" to null,
                "PeerAS" to null,
                "LocalAS" to null,
                "LogStateChangesEnabled" to null,
                "MaxPrefix" to null
        ))
    }
}
