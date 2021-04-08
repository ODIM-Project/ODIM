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

import com.odim.simulator.behaviors.ExecutionMode.CHAINED
import com.odim.simulator.behaviors.ExecutionMode.EXCLUSIVE
import com.odim.simulator.behaviors.ExecutionMode.PREPEND
import com.odim.simulator.behaviors.RestBehaviors.onDelete
import com.odim.simulator.behaviors.RestBehaviors.onGet
import com.odim.simulator.behaviors.RestBehaviors.onOptions
import com.odim.simulator.behaviors.RestBehaviors.onPatch
import com.odim.simulator.behaviors.RestBehaviors.onPost
import com.odim.simulator.http.HttpMethod
import com.odim.simulator.http.HttpMethod.DELETE
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.HttpMethod.OPTIONS
import com.odim.simulator.http.HttpMethod.PATCH
import com.odim.simulator.http.HttpMethod.POST
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.notFound
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.ActionType
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.TreeElement
import com.odim.simulator.tree.structure.Type

enum class ExecutionMode(var order: Int) {
    EXCLUSIVE(1),
    PREPEND(2),
    CHAINED(3)
}
typealias ElementPredicate = (Item, HttpMethod) -> Boolean

class BehaviorResponse private constructor(val response: Response, val terminal: Boolean) {
    companion object {
        fun terminal(response: Response) = BehaviorResponse(response, terminal = true)
        fun nonTerminal(response: Response) = BehaviorResponse(response, terminal = false)
    }
}

data class BehaviorHandler(val predicate: ElementPredicate, val behavior: Behavior, val executionMode: ExecutionMode = CHAINED)

@Suppress("TooManyFunctions")
class ResourceBehaviors {

    private val behaviorsMap: MutableMap<ExecutionMode, MutableList<BehaviorHandler>> = mutableMapOf()
    private val dataStore: BehaviorDataStore = BehaviorDataStore()

    init {
        ExecutionMode.values().forEach {
            behaviorsMap[it] = mutableListOf()
        }

        // default Rest behaviors
        behaviorsMap[CHAINED]?.addAll(mutableListOf(
                BehaviorHandler({ _, httpMethod -> httpMethod == GET }, onGet()),
                BehaviorHandler({ _, httpMethod -> httpMethod == POST }, onPost()),
                BehaviorHandler({ _, httpMethod -> httpMethod == PATCH }, onPatch()),
                BehaviorHandler({ _, httpMethod -> httpMethod == DELETE }, onDelete()),
                BehaviorHandler({ _, httpMethod -> httpMethod == OPTIONS }, onOptions())))
    }

    fun appendBehavior(item: Item, httpMethod: HttpMethod, behavior: Behavior) =
            bind(item, httpMethod, behavior, executionMode = CHAINED)

    fun appendBehavior(type: Type, httpMethod: HttpMethod, behavior: Behavior) =
            bind(type, httpMethod, behavior, executionMode = CHAINED)

    fun appendBehavior(httpMethod: HttpMethod, behavior: Behavior) =
            bind(httpMethod, behavior, executionMode = CHAINED)

    fun appendActionBehavior(actionType: ActionType, httpMethod: HttpMethod, behavior: Behavior) =
            bind(actionType, httpMethod, behavior, executionMode = CHAINED)

    fun appendActionBehavior(behavior: Behavior) = bindActions(behavior)

    fun prependBehavior(item: Item, httpMethod: HttpMethod, behavior: Behavior) =
            bind(item, httpMethod, behavior, executionMode = PREPEND)

    fun prependBehavior(type: Type, httpMethod: HttpMethod, behavior: Behavior) =
            bind(type, httpMethod, behavior, executionMode = PREPEND)

    fun prependBehavior(httpMethod: HttpMethod, behavior: Behavior) =
            bind(httpMethod, behavior, executionMode = PREPEND)

    fun replaceBehavior(item: Item, httpMethod: HttpMethod, behavior: Behavior) =
            bind(item, httpMethod, behavior, executionMode = EXCLUSIVE)

    fun replaceBehavior(type: Type, httpMethod: HttpMethod, behavior: Behavior) =
            bind(type, httpMethod, behavior, executionMode = EXCLUSIVE)

    fun replaceBehavior(httpMethod: HttpMethod, behavior: Behavior) =
            bind(httpMethod, behavior, executionMode = EXCLUSIVE)

    fun replaceActionBehavior(actionType: ActionType, behavior: Behavior) =
            bind(actionType, POST, behavior, executionMode = EXCLUSIVE)

    fun replaceActionBehavior(behavior: Behavior) = bindActions(behavior, executionMode = EXCLUSIVE)

    fun removeBehavior(handler: BehaviorHandler) = behaviorsMap[handler.executionMode]?.remove(handler)

    private fun bind(predicate: ElementPredicate, behavior: Behavior, executionMode: ExecutionMode = CHAINED) =
            BehaviorHandler(predicate, behavior, executionMode).also {
                behaviorsMap[executionMode]?.add(it)
            }

    private fun bind(item: Item, on: HttpMethod, behavior: Behavior, executionMode: ExecutionMode = CHAINED): BehaviorHandler {
        val predicate: ElementPredicate = { res, onMethod -> res == item && onMethod == on }
        return bind(predicate, behavior, executionMode)
    }

    private fun bind(type: Type, on: HttpMethod, behavior: Behavior, executionMode: ExecutionMode = CHAINED): BehaviorHandler {
        val predicate: ElementPredicate = { item, onMethod -> item is TreeElement && item.meta.type == type && onMethod == on }
        return bind(predicate, behavior, executionMode)
    }

    private fun bind(on: HttpMethod, behavior: Behavior, executionMode: ExecutionMode = CHAINED): BehaviorHandler {
        val predicate: ElementPredicate = { item, onMethod -> item is TreeElement && onMethod == on }
        return bind(predicate, behavior, executionMode)
    }

    private fun bind(actionType: ActionType, on: HttpMethod, behavior: Behavior, executionMode: ExecutionMode = CHAINED): BehaviorHandler {
        val predicate: ElementPredicate = { item, onMethod -> item is ActionElement && item.actionType == actionType && onMethod == on }
        return bind(predicate, behavior, executionMode)
    }

    private fun bindActions(behavior: Behavior, executionMode: ExecutionMode = CHAINED): BehaviorHandler {
        val predicate: ElementPredicate = { item, onMethod -> item is ActionElement && onMethod == POST }
        return bind(predicate, behavior, executionMode)
    }

    fun createResponse(tree: ResourceTree, element: Item, req: Request): Response {
        val response = notFound()
        val handlers = getBehaviorsForItem(element, req)
        return applyBehaviors(tree, handlers, element, req, response)
    }

    private fun applyBehaviors(tree: ResourceTree, handlers: List<BehaviorHandler>, element: Item, req: Request, response: Response): Response {
        var result = response
        handlers.forEach {
            val behaviorResponse = it.behavior.run(tree, element, req, result, dataStore)
            result = behaviorResponse.response
            if (behaviorResponse.terminal) return result
        }
        return result
    }

    private fun getBehaviorsForItem(element: Item, req: Request): List<BehaviorHandler> =
            behaviorsMap.mapValues { (_, v) ->
                v.filter { it.predicate(element, req.method) }.toMutableList()
            }.toSortedMap(compareBy { it.order })
                    .flatMap {
                        if (EXCLUSIVE == it.key && it.value.isNotEmpty()) return listOf(it.value.last())
                        it.value
                    }
}
