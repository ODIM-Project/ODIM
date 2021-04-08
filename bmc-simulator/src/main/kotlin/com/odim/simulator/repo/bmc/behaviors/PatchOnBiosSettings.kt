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
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.utils.getString

class PatchOnBiosSettings : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore) =
        try {
            val updatedSettings: MutableMap<String, Int> = mutableMapOf()
            request.json!!.fieldNames().forEach { updatedSettings[it] = request.json!!.getString(it).toInt() }
            dataStore.insert(BIOS_SETTINGS, updatedSettings)?.let {
                (it as MutableMap<String, Int>).putAll(updatedSettings)
            }
            nonTerminal(success())
        } catch (e: NumberFormatException) {
            nonTerminal(badRequest())
        }
}
