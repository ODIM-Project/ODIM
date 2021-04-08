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

package com.odim.simulator.repo.behaviors

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.TreeJsonRenderer
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import java.util.regex.Pattern.compile

class GetOnChassis : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        updateChassisId(item as Resource)
        return nonTerminal(Response(response.code, TreeJsonRenderer(request.expandType, request.expandLevels).toJson(item), response.headers))
    }

    private fun updateChassisId(chassis: Resource) {
        val chassisId = chassis.getId()
        val chassisPattern = compile("/redfish/v1/Chassis(/(?<mainChassisId>[A-Za-z0-9-.%:+_]+)(/(?<chassisId>[A-Za-z0-9-.%:+_]+))?)?")
        val groupMatcher = chassisPattern.matcher(chassisId)
        if (groupMatcher.matches() && groupMatcher.group("chassisId") != null) {
            chassis {
                "Id" to groupMatcher.group("chassisId")
            }
        }
    }
}
