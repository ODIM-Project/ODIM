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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CREDENTIAL_BOOTSTRAPPING
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.ETHERNET_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.HOST_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.MANAGER_NETWORK_PROTOCOL
import com.odim.simulator.tree.structure.ResourceType.ROLE
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(HOST_INTERFACE)
open class HostInterfaceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Host Interface Description",
                "Name" to "Host Interface",
                "HostInterfaceType" to null,
                "Status" to embeddedObject(STATUS),
                "InterfaceEnabled" to null,
                "ExternallyAccessible" to null,
                "AuthenticationModes" to embeddedArray(),
                "KernelAuthRoleId" to "",
                "KernelAuthEnabled" to null,
                "FirmwareAuthRoleId" to "",
                "FirmwareAuthEnabled" to null,
                "HostEthernetInterfaces" to ResourceCollection(ETHERNET_INTERFACES_COLLECTION),
                "ManagerEthernetInterface" to LinkableResource(ETHERNET_INTERFACE),
                "NetworkProtocol" to SingletonResource(MANAGER_NETWORK_PROTOCOL),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ComputerSystems" to LinkableResourceArray(COMPUTER_SYSTEM),
                        "KernelAuthRole" to LinkableResource(ROLE),
                        "FirmwareAuthRole" to LinkableResource(ROLE)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_1_0, V1_0_1, resourceObject(
                "Actions" to Actions()
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_2_0, V1_1_3, resourceObject(
                "AuthNoneRoleId" to "",
                "Links" to embeddedObject(
                        "AuthNoneRole" to LinkableResource(ROLE)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_3_0, V1_2_2, resourceObject(
                "CredentialBootstrapping" to embeddedObject(CREDENTIAL_BOOTSTRAPPING),
                "Links" to embeddedObject(
                        "CredentialBootstrappingRole" to LinkableResource(ROLE)
                )
        ))
    }
}
