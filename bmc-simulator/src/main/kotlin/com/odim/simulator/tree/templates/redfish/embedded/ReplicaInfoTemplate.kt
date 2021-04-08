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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.structure.EmbeddedObjectType.REPLICA_INFO
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.DATA_PROTECTION_LINE_OF_SERVICE

/**
 * This is generated class. Please don't edit it's contents.
 */
@EmbeddedObjectTemplate(REPLICA_INFO)
open class ReplicaInfoTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, embeddedObject(
                "ReplicaPriority" to null,
                "ReplicaReadOnlyAccess" to null,
                "UndiscoveredElement" to null,
                "WhenSynced" to null,
                "SyncMaintained" to false,
                "ReplicaRecoveryMode" to null,
                "ReplicaUpdateMode" to null,
                "PercentSynced" to null,
                "FailedCopyStopsHostIO" to false,
                "WhenActivated" to null,
                "WhenDeactivated" to null,
                "WhenEstablished" to null,
                "WhenSuspended" to null,
                "WhenSynchronized" to null,
                "ReplicaSkewBytes" to null,
                "ReplicaType" to null,
                "ReplicaProgressStatus" to null,
                "ReplicaState" to null,
                "RequestedReplicaState" to null,
                "ConsistencyEnabled" to false,
                "ConsistencyType" to null,
                "ConsistencyState" to null,
                "ConsistencyStatus" to null,
                "ReplicaRole" to null,
                "Replica" to LinkableResource(ANY)
        ))
        version(V1_1_0, V1_0_0, embeddedObject(
                "DataProtectionLineOfService" to LinkableResource(DATA_PROTECTION_LINE_OF_SERVICE)
        ))
        version(V1_2_0, V1_1_0, embeddedObject(
                "SourceReplica" to LinkableResource(ANY)
        ))
        version(V1_3_0, V1_2_0, embeddedObject(
                "ReplicaFaultDomain" to null
        ))
    }
}
