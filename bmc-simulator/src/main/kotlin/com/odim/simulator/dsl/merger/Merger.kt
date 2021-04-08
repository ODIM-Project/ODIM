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

package com.odim.simulator.dsl.merger

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.http.Request
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Resource
import com.odim.utils.JsonMapper.emptyJson

class Merger(private val config: (MergeDslConfiguration.() -> Unit)? = null) {
    companion object {
        fun merge(tree: ResourceTree, resource: Resource, request: Request) = Merger().merge(tree, resource, request)
        fun merge(tree: ResourceTree, resource: Resource, override: ObjectNode) = Merger().merge(tree, resource, override)
    }

    fun merge(tree: ResourceTree, resource: Resource, request: Request) {
        if (request.body != null) {
            merge(tree, resource, request.json ?: emptyJson)
        }
    }

    fun merge(tree: ResourceTree, resource: Resource, override: ObjectNode) {
        val plan = Plan.Builder().withConfig(config).base(resource.data).override(override).build()
        MergeExecutor.buildLinking(tree, resource.toLink()).execute(plan)
    }
}

open class MergeException(override val message: String, cause: Throwable? = null) : Exception(message, cause)
