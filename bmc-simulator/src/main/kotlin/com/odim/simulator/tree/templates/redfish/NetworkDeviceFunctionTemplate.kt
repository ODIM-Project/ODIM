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
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_2_6
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_3_4
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.FIBRE_CHANNEL
import com.odim.simulator.tree.structure.EmbeddedObjectType.INFINI_BAND
import com.odim.simulator.tree.structure.EmbeddedObjectType.ISCSI_BOOT
import com.odim.simulator.tree.structure.EmbeddedObjectType.NETWORK_DEVICE_FUNCTION_ETHERNET
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.NETWORK_DEVICE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.NETWORK_PORT
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.PORT

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(NETWORK_DEVICE_FUNCTION)
open class NetworkDeviceFunctionTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Network Device Function Description",
                "Name" to "Network Device Function",
                "Status" to embeddedObject(STATUS),
                "NetDevFuncType" to null,
                "DeviceEnabled" to null,
                "NetDevFuncCapabilities" to embeddedArray(),
                "Ethernet" to embeddedObject(NETWORK_DEVICE_FUNCTION_ETHERNET),
                "iSCSIBoot" to embeddedObject(ISCSI_BOOT),
                "FibreChannel" to embeddedObject(FIBRE_CHANNEL),
                "AssignablePhysicalPorts" to LinkableResourceArray(NETWORK_PORT),
                "PhysicalPortAssignment" to LinkableResource(NETWORK_PORT),
                "BootMode" to null,
                "VirtualFunctionsEnabled" to null,
                "MaxVirtualFunctions" to null,
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "PCIeFunction" to LinkableResource(PCIE_FUNCTION)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_1_0, V1_0_1, resourceObject(
                "Actions" to Actions()
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_2_0, V1_1_1, embeddedObject(
                "Links" to embeddedObject(
                        "Endpoints" to LinkableResourceArray(ENDPOINT)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_2_6, V1_2_5)
        version(V1_3_0, V1_2_2, embeddedObject(
                "Links" to embeddedObject(
                        "PhysicalPortAssignment" to LinkableResource(NETWORK_PORT)
                )
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_3_4, V1_3_3)
        version(V1_4_0, V1_3_3, embeddedObject(
                "Links" to embeddedObject(
                        "EthernetInterface" to LinkableResource(ETHERNET_INTERFACE)
                )
        ))
        version(V1_4_1, V1_4_0)
        version(V1_5_0, V1_4_1, resourceObject(
                "AssignablePhysicalNetworkPorts" to LinkableResourceArray(PORT),
                "PhysicalNetworkPortAssignment" to LinkableResource(PORT),
                "InfiniBand" to embeddedObject(INFINI_BAND),
                "Links" to embeddedObject(
                        "PhysicalNetworkPortAssignment" to LinkableResource(PORT)
                )
        ))
    }
}
