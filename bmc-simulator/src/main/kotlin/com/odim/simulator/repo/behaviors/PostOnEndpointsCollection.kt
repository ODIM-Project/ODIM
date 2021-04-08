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
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.dsl.merger.MergeException
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.created
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.utils.getStringOrNull

class PostOnEndpointsCollection : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val endpointsCollection = item as ResourceCollection
        val endpoint = tree.create(endpointsCollection.type.of())
        endpointsCollection.add(endpoint)

        return terminal(try {
            Merger.merge(tree, endpoint, request)

            request.json!!.getStringOrNull("ConnectedEntities/0/EntityLink/@odata.id")?.let {
                val drive = tree.search(it) as Resource
                tree.link(endpoint, drive, "EntityLink")
            }

            created(endpoint)
        } catch (e: MergeException) {
            badRequest("Bad POST request: ${e.message}")
        })
    }
}
