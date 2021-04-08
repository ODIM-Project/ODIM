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

package com.odim.simulator.behaviors

import com.odim.simulator.behaviors.Behavior.Companion.behavior
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.dsl.merger.MergeException
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.created
import com.odim.simulator.http.Response.Companion.noContent
import com.odim.simulator.http.Response.Companion.notAllowed
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection

object RestBehaviors {
    fun onGet() = behavior { nonTerminal(success(item, expandType = request.expandType, expandLevels = request.expandLevels)) }
    fun onPost() = behavior { nonTerminal(post(tree, item, request)) }
    fun onPatch() = behavior { nonTerminal(patch(tree, item, request)) }
    fun onDelete() = behavior { nonTerminal(delete(item)) }
    fun onOptions() = behavior { terminal(options(item)) }

    private fun post(sim: ResourceTree, res: Item, request: Request) =
            when (res) {
                is ResourceCollection -> postOnCollection(sim, res, request)
                is ActionElement -> noContent()
                else -> notAllowed()
            }

    private fun patch(sim: ResourceTree, res: Item, req: Request) =
            if (res is Resource) patchResource(sim, res, req) else notAllowed()

    private fun delete(res: Item) =
            if (res is Resource) deleteOnResource(res) else notAllowed()

    private fun options(res: Item): Response {
        return success(mutableMapOf(
                "Allow" to listOf("OPTIONS of $res")
        ))
    }

    private fun deleteOnResource(resource: Resource): Response {
        val parent = resource.meta.parent
        if (parent is ResourceCollection) {
            parent.members.remove(resource)
            return noContent()
        }
        throw IllegalStateException("Parent is not a collection")
    }

    private fun patchResource(resourceTree: ResourceTree, resource: Resource, request: Request) = try {
        Merger.merge(resourceTree, resource, request)
        success(resource)
    } catch (e: MergeException) {
        badRequest("Bad PATCH request: ${e.message}")
    }

    private fun postOnCollection(resourceTree: ResourceTree, collection: ResourceCollection, request: Request): Response {
        val newResource = resourceTree.create(collection.type.of())
        return try {
            collection.add(newResource)
            Merger.merge(resourceTree, newResource, request)
            created(newResource)
        } catch (e: MergeException) {
            badRequest("Bad POST request: ${e.message}")
        }
    }
}
