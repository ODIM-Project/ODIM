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
import com.odim.simulator.tree.structure.ResourceCollection

class GetOnComputerSystem : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        updateComputerSystem(item as Resource)
        return nonTerminal(Response(response.code, TreeJsonRenderer(request.expandType, request.expandLevels).toJson(item), response.headers))
    }

    private fun updateComputerSystem(system: Resource) {
        system {
            "ProcessorSummary" to {
                "Count" to system.traverse<ResourceCollection>("Processors").size()
                "Model" to system.traverse<ResourceCollection>("Processors").members.firstOrNull()?.traverse<String?>("Model")
            }
            "MemorySummary" to {
                "TotalSystemMemoryGiB" to system.traverse<ResourceCollection>("Memory").members.sumBy {
                    it.traverse("CapacityMiB") ?: 0
                } / 1024
            }
        }
    }
}
