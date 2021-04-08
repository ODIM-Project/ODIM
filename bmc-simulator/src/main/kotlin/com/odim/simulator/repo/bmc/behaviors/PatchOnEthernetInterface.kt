/*
 * Copyright (c) Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.odim.simulator.repo.bmc.behaviors

import com.fasterxml.jackson.databind.node.ArrayNode
import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.EmbeddedArray
import com.odim.simulator.tree.structure.EmbeddedObjectType
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V4_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_ADDRESS
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_ADDRESS_POLICY_ENTRY
import com.odim.simulator.tree.structure.EmbeddedObjectType.IP_V6_GATEWAY_STATIC_ADDRESS
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.utils.getArrayOrNull
import com.odim.utils.getBooleanOrNull
import com.odim.utils.getNumberOrNull
import com.odim.utils.getStringOrNull

@Suppress("StringLiteralDuplication")
class PatchOnEthernetInterface : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val ethernetInterface = item as Resource
        updateAutoNeg(request, ethernetInterface)
        updateFQDN(request, ethernetInterface)
        updateFullDuplex(request, ethernetInterface)
        updateHostName(request, ethernetInterface)
        updateIPv4Addresses(request, tree, ethernetInterface)
        updateIPv4StaticAddresses(request, tree, ethernetInterface)
        updateIPv6Addresses(request, tree, ethernetInterface)
        updateIPv6DefaultGateway(request, ethernetInterface)
        updateIPv6StaticAddresses(request, tree, ethernetInterface)
        updateIPv6StaticDefaultGateways(request, tree, ethernetInterface)
        updateIPv6AddressPolicyTable(request, tree, ethernetInterface)
        updateInterfaceEnabled(request, ethernetInterface)
        updateMTUSize(request, ethernetInterface)
        return terminal(success(ethernetInterface))
    }

    private fun updateIPv6AddressPolicyTable(request: Request, tree: ResourceTree, ethernetInterface: Resource) {
        request.json?.getArrayOrNull("IPv6AddressPolicyTable")?.let { ipsJson ->
            updateArray(ipsJson, tree, ethernetInterface, "IPv6AddressPolicyTable", IP_V6_ADDRESS_POLICY_ENTRY) {
                mapOf(
                        "Prefix" to it.getStringOrNull("Prefix"),
                        "Precedence" to it.getNumberOrNull("Precedence"),
                        "Label" to it.getNumberOrNull("Label")
                )
            }
        }
    }

    private fun updateMTUSize(request: Request, ethernetInterface: Resource) {
        request.json?.getNumberOrNull("MTUSize")?.let {
            ethernetInterface { "MTUSize" to it }
        }
    }

    private fun updateInterfaceEnabled(request: Request, ethernetInterface: Resource) {
        request.json?.getBooleanOrNull("InterfaceEnabled")?.let {
            ethernetInterface { "InterfaceEnabled" to it }
        }
    }

    private fun updateIPv6StaticDefaultGateways(request: Request, tree: ResourceTree, ethernetInterface: Resource) {
        request.json?.getArrayOrNull("IPv6StaticDefaultGateways")?.let { ipsJson ->
            updateArray(ipsJson, tree, ethernetInterface, "IPv6StaticDefaultGateways", IP_V6_GATEWAY_STATIC_ADDRESS) {
                mapOf(
                        "Address" to it.getStringOrNull("Address"),
                        "PrefixLength" to it.getNumberOrNull("PrefixLength")
                )
            }
        }
    }

    private fun updateIPv6StaticAddresses(request: Request, tree: ResourceTree, ethernetInterface: Resource) {
        request.json?.getArrayOrNull("IPv6StaticAddresses")?.let { ipsJson ->
            updateArray(ipsJson, tree, ethernetInterface, "IPv6StaticAddresses", IP_V6_ADDRESS) {
                mapOf(
                        "Address" to it.getStringOrNull("Address"),
                        "PrefixLength" to it.getNumberOrNull("PrefixLength"),
                        "AddressOrigin" to it.getStringOrNull("AddressOrigin"),
                        "AddressState" to it.getStringOrNull("AddressState")
                )
            }
        }
    }

    private fun updateIPv6DefaultGateway(request: Request, ethernetInterface: Resource) {
        request.json?.getStringOrNull("IPv6DefaultGateway")?.let {
            ethernetInterface { "IPv6DefaultGateway" to it }
        }
    }

    private fun updateIPv6Addresses(request: Request, tree: ResourceTree, ethernetInterface: Resource) {
        request.json?.getArrayOrNull("IPv6Addresses")?.let { ipsJson ->
            updateArray(ipsJson, tree, ethernetInterface, "IPv6Addresses", IP_V6_ADDRESS) {
                mapOf(
                        "Address" to it.getStringOrNull("Address"),
                        "PrefixLength" to it.getNumberOrNull("PrefixLength"),
                        "AddressOrigin" to it.getStringOrNull("AddressOrigin"),
                        "AddressState" to it.getStringOrNull("AddressState")
                )
            }
        }
    }

    private fun updateIPv4StaticAddresses(request: Request, tree: ResourceTree, ethernetInterface: Resource) {
        request.json?.getArrayOrNull("IPv4StaticAddresses")?.let { ipsJson ->
            updateArray(ipsJson, tree, ethernetInterface, "IPv4StaticAddresses", IP_V4_ADDRESS) {
                mapOf(
                        "Address" to it.getStringOrNull("Address"),
                        "SubnetMask" to it.getStringOrNull("SubnetMask"),
                        "AddressOrigin" to it.getStringOrNull("AddressOrigin"),
                        "Gateway" to it.getStringOrNull("Gateway")
                )
            }
        }
    }

    private fun updateIPv4Addresses(request: Request, tree: ResourceTree, ethernetInterface: Resource) {
        request.json?.getArrayOrNull("IPv4Addresses")?.let { ipsJson ->
            updateArray(ipsJson, tree, ethernetInterface, "IPv4Addresses", IP_V4_ADDRESS) {
                mapOf(
                        "Address" to it.getStringOrNull("Address"),
                        "SubnetMask" to it.getStringOrNull("SubnetMask"),
                        "AddressOrigin" to it.getStringOrNull("AddressOrigin"),
                        "Gateway" to it.getStringOrNull("Gateway")
                )
            }
        }
    }

    private fun updateHostName(request: Request, ethernetInterface: Resource) {
        request.json?.getStringOrNull("HostName")?.let {
            ethernetInterface { "HostName" to it }
        }
    }

    private fun updateFullDuplex(request: Request, ethernetInterface: Resource) {
        request.json?.getBooleanOrNull("FullDuplex")?.let {
            ethernetInterface { "FullDuplex" to it }
        }
    }

    private fun updateFQDN(request: Request, ethernetInterface: Resource) {
        request.json?.getStringOrNull("FQDN")?.let {
            ethernetInterface { "FQDN" to it }
        }
    }

    private fun updateAutoNeg(request: Request, ethernetInterface: Resource) {
        request.json?.getBooleanOrNull("AutoNeg")?.let {
            ethernetInterface { "AutoNeg" to it }
        }
    }

    private fun updateArray(json: ArrayNode, tree: ResourceTree,
                            ethernetInterface: Resource,
                            path: String,
                            type: EmbeddedObjectType,
                            objConstructFun: (json: ObjectNode) -> Map<String, Any?>) {
        val array = ethernetInterface.traverse<EmbeddedArray<ResourceObject>>(path)
        array.clear()
        json.forEach { ipJson ->
            array.add(tree.createEmbeddedObject(type) {
                objConstructFun(ipJson as ObjectNode).forEach { (key, value) ->
                    key to value
                }
            })
        }
    }
}
