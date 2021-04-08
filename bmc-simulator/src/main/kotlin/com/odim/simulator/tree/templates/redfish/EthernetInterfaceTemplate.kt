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
import com.odim.simulator.tree.RedfishVersion.V1_0_10
import com.odim.simulator.tree.RedfishVersion.V1_0_11
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_0_8
import com.odim.simulator.tree.RedfishVersion.V1_0_9
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_1_7
import com.odim.simulator.tree.RedfishVersion.V1_1_8
import com.odim.simulator.tree.RedfishVersion.V1_1_9
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_2_6
import com.odim.simulator.tree.RedfishVersion.V1_2_7
import com.odim.simulator.tree.RedfishVersion.V1_2_8
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_3_4
import com.odim.simulator.tree.RedfishVersion.V1_3_5
import com.odim.simulator.tree.RedfishVersion.V1_3_6
import com.odim.simulator.tree.RedfishVersion.V1_3_7
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_4_3
import com.odim.simulator.tree.RedfishVersion.V1_4_4
import com.odim.simulator.tree.RedfishVersion.V1_4_5
import com.odim.simulator.tree.RedfishVersion.V1_4_6
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_5_1
import com.odim.simulator.tree.RedfishVersion.V1_5_2
import com.odim.simulator.tree.RedfishVersion.V1_5_3
import com.odim.simulator.tree.RedfishVersion.V1_5_4
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.DHCPV4_CONFIGURATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.DHCPV6_CONFIGURATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V4_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_ADDRESS_POLICY_ENTRY
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_GATEWAY_STATIC_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_STATIC_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATELESS_ADDRESS_AUTO_CONFIGURATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedObjectType.VLAN
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.VLAN_NETWORK_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.HOST_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.NETWORK_DEVICE_FUNCTION

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(ETHERNET_INTERFACE)
open class EthernetInterfaceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Ethernet Interface Description",
                "Name" to "Ethernet Interface",
                "UefiDevicePath" to null,
                "Status" to embeddedObject(STATUS),
                "InterfaceEnabled" to null,
                "PermanentMACAddress" to null,
                "MACAddress" to null,
                "SpeedMbps" to null,
                "AutoNeg" to null,
                "FullDuplex" to null,
                "MTUSize" to null,
                "HostName" to null,
                "FQDN" to null,
                "MaxIPv6StaticAddresses" to null,
                "VLAN" to embeddedObject(VLAN),
                "IPv4Addresses" to embeddedArray(IP_V4_ADDRESS),
                "IPv6AddressPolicyTable" to embeddedArray(IP_V6_ADDRESS_POLICY_ENTRY),
                "IPv6Addresses" to embeddedArray(IP_V6_ADDRESS),
                "IPv6StaticAddresses" to embeddedArray(IP_V6_STATIC_ADDRESS),
                "IPv6DefaultGateway" to null,
                "NameServers" to embeddedArray(),
                "VLANs" to ResourceCollection(VLAN_NETWORK_INTERFACES_COLLECTION)
        ))
        version(V1_0_2, V1_0_0)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_0_9, V1_0_8)
        version(V1_0_10, V1_0_9)
        version(V1_0_11, V1_0_10)
        version(V1_1_0, V1_0_2, resourceObject(
                "LinkStatus" to null,
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Endpoints" to LinkableResourceArray(ENDPOINT)
                )
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_1_7, V1_1_6)
        version(V1_1_8, V1_1_7)
        version(V1_1_9, V1_1_8)
        version(V1_2_0, V1_1_1, embeddedObject(
                "Links" to embeddedObject(
                        "HostInterface" to LinkableResource(HOST_INTERFACE)
                )
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_2_6, V1_2_5)
        version(V1_2_7, V1_2_6)
        version(V1_2_8, V1_2_7)
        version(V1_3_0, V1_2_1, resourceObject(
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Chassis" to LinkableResource(CHASSIS)
                )
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_3_4, V1_3_3)
        version(V1_3_5, V1_3_4)
        version(V1_3_6, V1_3_5)
        version(V1_3_7, V1_3_6)
        version(V1_4_0, V1_3_1, resourceObject(
                "DHCPv4" to embeddedObject(DHCPV4_CONFIGURATION),
                "DHCPv6" to embeddedObject(DHCPV6_CONFIGURATION),
                "StatelessAddressAutoConfig" to embeddedObject(STATELESS_ADDRESS_AUTO_CONFIGURATION),
                "IPv6StaticDefaultGateways" to embeddedArray(IP_V6_GATEWAY_STATIC_ADDRESS),
                "StaticNameServers" to embeddedArray(),
                "IPv4StaticAddresses" to embeddedArray(IP_V4_ADDRESS)
        ))
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_4_3, V1_4_2)
        version(V1_4_4, V1_4_3)
        version(V1_4_5, V1_4_4)
        version(V1_4_6, V1_4_5)
        version(V1_5_0, V1_4_2)
        version(V1_5_1, V1_5_0)
        version(V1_5_2, V1_5_1)
        version(V1_5_3, V1_5_2)
        version(V1_5_4, V1_5_3)
        version(V1_6_0, V1_5_2, resourceObject(
                "EthernetInterfaceType" to null,
                "Links" to embeddedObject(
                        "NetworkDeviceFunction" to LinkableResource(NETWORK_DEVICE_FUNCTION)
                )
        ))
        version(V1_6_1, V1_6_0)
        version(V1_6_2, V1_6_1)
    }
}
