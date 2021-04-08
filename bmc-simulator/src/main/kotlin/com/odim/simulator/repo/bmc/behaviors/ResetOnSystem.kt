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
import com.odim.simulator.behaviors.BehaviorDataStore.SharedInformationType.BIOS_SETTINGS
import com.odim.simulator.behaviors.BehaviorDataStore.SharedInformationType.FIRMWARE_UPDATE_MESSAGES
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.noContent
import com.odim.simulator.repo.bmc.createSelEntry
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.utils.getString
import com.odim.utils.getStringOrNull

class ResetOnSystem(private val logService: Resource, private val bios: Resource) : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        item as Action
        if(!checkAllowableValues(item, request)) {
            return nonTerminal(badRequest())
        }
        return when (item.parent!!.meta.type) {
            COMPUTER_SYSTEM -> resetOnSystem(item, request, dataStore, tree)
            else -> terminal(noContent())
        }
    }

    private fun resetOnSystem(item: Action, request: Request, dataStore: BehaviorDataStore, tree: ResourceTree): BehaviorResponse {
        var entryMessage = "OEM System Boot Event"
        when (resetState(request)) {
            "On", "ForceOn", "ForceRestart" -> {
                setPowerState(item, "On")
                updateBiosSettings(dataStore, bios, tree)
            }
            "ForceOff", "GracefulShutdown" -> {
                setPowerState(item, "Off")
                entryMessage = "Power Off / Power Down"
            }
            "PushPowerButton" -> when (getSystemPowerState(item)) {
                "On" -> setPowerState(item, "Off")
                "Off" -> setPowerState(item, "On")
            }
        }
        createSelEntry(logService, tree, entryMessage)
        appendFirmwareUpdateMessages(dataStore, tree)
        return nonTerminal(noContent())
    }

    private fun checkAllowableValues(action: Action, request: Request) =
        action.allowableValues["ResetType"]?.contains(request.json?.getStringOrNull("ResetType")) ?: false


    private fun appendFirmwareUpdateMessages(dataStore: BehaviorDataStore, tree: ResourceTree) {
        val messages = dataStore.readAndRemove(FIRMWARE_UPDATE_MESSAGES) as? List<*>

        messages?.forEach {
            createSelEntry(logService, tree, it.toString())
        }
    }

    private fun updateBiosSettings(dataStore: BehaviorDataStore, bios: Resource, tree: ResourceTree) {
        dataStore.readAndRemove(BIOS_SETTINGS)?.let {
            (it as MutableMap<String, Int>).iterator().forEach {
                bios {
                    "Attributes" to {
                        it.key to it.value
                    }
                }
            }
            createSelEntry(logService, tree, "Bios settings updated")
        }
    }

    private fun setPowerState(item: Item, powerState: String) {
        ((item as Action).parent as Resource) {
            "PowerState" to powerState
        }
    }

    private fun getSystemPowerState(item: Item) = ((item as Action).parent as Resource).traverse("PowerState") as String
    private fun resetState(request: Request) = request.json!!.getString("ResetType")
}
