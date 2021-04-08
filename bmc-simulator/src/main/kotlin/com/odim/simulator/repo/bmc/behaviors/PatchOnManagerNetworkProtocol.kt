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

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.EmbeddedArray
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.utils.getArrayOrNull
import com.odim.utils.getBooleanOrNull
import com.odim.utils.getNumberOrNull
import com.odim.utils.getObjectOrNull
import com.odim.utils.getStringOrNull

class PatchOnManagerNetworkProtocol : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val managerNetworkProtocol = item as Resource
        updateFQDN(request, managerNetworkProtocol)
        updateHTTPS(request, managerNetworkProtocol)
        updateNTP(request, managerNetworkProtocol)
        updateSSDP(request, managerNetworkProtocol)
        updateSSH(request, managerNetworkProtocol)

        return terminal(success(managerNetworkProtocol))
    }

    private fun updateSSH(request: Request, managerNetworkProtocol: Resource) {
        request.json?.getObjectOrNull("SSH")?.let {
            updateProtocol("SSH", managerNetworkProtocol, mapOf(
                    "ProtocolEnabled" to it.getBooleanOrNull("ProtocolEnabled"),
                    "Port" to it.getNumberOrNull("Port")
            ))
        }
    }

    private fun updateSSDP(request: Request, managerNetworkProtocol: Resource) {
        request.json?.getObjectOrNull("SSDP")?.let {
            updateProtocol("SSDP", managerNetworkProtocol, mapOf(
                    "ProtocolEnabled" to it.getBooleanOrNull("ProtocolEnabled"),
                    "Port" to it.getNumberOrNull("Port"),
                    "NotifyMulticastIntervalSeconds" to it.getNumberOrNull("NotifyMulticastIntervalSeconds"),
                    "NotifyTTL" to it.getNumberOrNull("NotifyTTL"),
                    "NotifyIPv6Scope" to it.getStringOrNull("NotifyIPv6Scope"),
            ))
        }
    }

    private fun updateNTP(request: Request, managerNetworkProtocol: Resource) {
        request.json?.getObjectOrNull("NTP")?.let {
            updateProtocol("NTP", managerNetworkProtocol, mapOf(
                    "ProtocolEnabled" to it.getBooleanOrNull("ProtocolEnabled"),
                    "Port" to it.getNumberOrNull("Port")
            ))
            it.getArrayOrNull("NTPServers")?.let { serversJson ->
                val traverse = managerNetworkProtocol.traverse<EmbeddedArray<String>>("NTP/NTPServers")
                traverse.clear()
                traverse.addAll(serversJson.toList().map { it.asText() }.toList())
            }
        }
    }

    private fun updateHTTPS(request: Request, managerNetworkProtocol: Resource) {
        request.json?.getObjectOrNull("HTTPS")?.let {
            updateProtocol("HTTPS", managerNetworkProtocol, mapOf(
                    "ProtocolEnabled" to it.getBooleanOrNull("ProtocolEnabled"),
                    "Port" to it.getNumberOrNull("Port"),
            ))
        }
    }

    private fun updateFQDN(request: Request, managerNetworkProtocol: Resource) {
        request.json?.getStringOrNull("FQDN")?.let {
            managerNetworkProtocol { "FQDN" to it }
        }
    }

    private fun updateProtocol(protocol: String, managerNetworkProtocol: Resource, values: Map<String, Any?>) {
        values.forEach { (key, value) ->
            value?.let {
                managerNetworkProtocol {
                    protocol to {
                        key to it
                    }
                }
            }
        }
    }
}
