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
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.utils.getObject
import com.odim.utils.getString

@Suppress("ReturnCount")
class PatchOnComputerSystem : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val targetAllowableValues = listOf("None", "Pxe", "Hdd", "Cd", "BiosSetup", "UefiShell", "Usb")
        val enabledAllowableValues = listOf("Disabled", "Once", "Continuous")
        val modeAllowableValues = listOf("Legacy", "UEFI")

        val bootNode = request.json?.getObject("Boot") ?: return terminal(Response.badRequest("Boot not found"))

        if (bootNode.has("BootSourceOverrideTarget")) {
            val bootSourceTargetValue = bootNode.getString("BootSourceOverrideTarget")
            if (!targetAllowableValues.contains(bootSourceTargetValue)) {
                return terminal(Response.badRequest("Value for BootSourceOverrideTarget is not valid"))
            }
        }
        if (bootNode.has("BootSourceOverrideEnabled")) {
            val bootSourceEnabledValue = bootNode.getString("BootSourceOverrideEnabled")
            if (!enabledAllowableValues.contains(bootSourceEnabledValue)) {
                return terminal(Response.badRequest("Value for BootSourceOverrideEnabled is not valid"))
            }
        }
        if (bootNode.has("BootSourceOverrideMode")) {
            val bootSourceModeValue = bootNode.getString("BootSourceOverrideMode")
            if (!modeAllowableValues.contains(bootSourceModeValue)) {
                return terminal(Response.badRequest("Value for BootSourceOverrideMode is not valid"))
            }
        }
        return nonTerminal(success())
    }
}
