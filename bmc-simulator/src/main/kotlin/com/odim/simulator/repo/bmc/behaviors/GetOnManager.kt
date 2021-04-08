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
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.utils.TIME_ZONE_PATTERN
import java.time.ZonedDateTime.now
import java.time.format.DateTimeFormatter.ISO_OFFSET_DATE_TIME
import java.time.format.DateTimeFormatter.ofPattern

class GetOnManager : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        updateTimeProperties(item as Resource)
        return nonTerminal(success(item))
    }

    private fun updateTimeProperties(resource: Resource) {
        resource.data.replace("DateTime", now().format(ISO_OFFSET_DATE_TIME))
        resource.data.replace("DateTimeLocalOffset", now().format(ofPattern(TIME_ZONE_PATTERN)))
    }
}
