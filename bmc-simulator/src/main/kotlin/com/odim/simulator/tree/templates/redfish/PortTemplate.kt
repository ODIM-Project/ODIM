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
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.LINK_CONFIGURATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.PORT_GEN_Z
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.PORT
import com.odim.simulator.tree.structure.ResourceType.PORT_METRICS
import com.odim.simulator.tree.structure.ResourceType.SWITCH
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(PORT)
open class PortTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Port Description",
                "Name" to "Port",
                "Status" to embeddedObject(STATUS),
                "PortId" to null,
                "PortProtocol" to null,
                "PortType" to null,
                "CurrentSpeedGbps" to null,
                "MaxSpeedGbps" to null,
                "Width" to null,
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "AssociatedEndpoints" to LinkableResourceArray(ENDPOINT),
                        "ConnectedSwitches" to LinkableResourceArray(SWITCH),
                        "ConnectedSwitchPorts" to LinkableResourceArray(PORT)
                ),
                "Actions" to Actions(
                        Action(RESET, "ResetType", mutableListOf(
                                "On",
                                "ForceOff",
                                "GracefulShutdown",
                                "GracefulRestart",
                                "ForceRestart",
                                "Nmi",
                                "ForceOn",
                                "PushPowerButton",
                                "PowerCycle"
                        ))
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_1_0, V1_0_3, resourceObject(
                "Location" to embeddedObject(RESOURCE_LOCATION)
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_2_0, V1_1_3, resourceObject(
                "PortMedium" to null,
                "LinkNetworkTechnology" to null,
                "InterfaceEnabled" to null,
                "SignalDetected" to null,
                "LinkTransitionIndicator" to 0,
                "ActiveWidth" to 0,
                "LinkState" to "Enabled",
                "LinkStatus" to "LinkUp",
                "GenZ" to embeddedObject(PORT_GEN_Z),
                "Metrics" to SingletonResource(PORT_METRICS),
                "Links" to embeddedObject(
                        "ConnectedPorts" to LinkableResourceArray(PORT)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_3_0, V1_2_2, resourceObject(
                "LocationIndicatorActive" to null,
                "MaxFrameSize" to null,
                "LinkConfiguration" to embeddedArray(LINK_CONFIGURATION),
                "FibreChannel" to null,
                "Ethernet" to null
        ))
    }
}
