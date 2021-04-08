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

package com.odim.simulator.dsl

import java.util.ArrayList

sealed class TreeNode {
    companion object {
        internal fun build(value: Any?): TreeNode = when (value) {
            is Map<*, *> -> MapNode.build(value)
            is Iterable<*> -> ArrayNode.build(value)
            else -> PrimitiveNode(value)
        }
    }

    abstract override fun toString(): String
}

internal class PrimitiveNode(val value: Any?) : TreeNode() {
    override fun toString() = if (value is String) """"$value"""" else value.toString()
}

internal class ArrayNode(values: List<TreeNode> = emptyList()) : TreeNode(), MutableList<TreeNode> by values.toMutableList() {
    companion object {
        internal fun build(source: Iterable<*>) = ArrayNode(source.map(TreeNode.Companion::build))
    }

    override fun toString() = toList().toString()
}

internal class MapNode(values: Map<String, TreeNode> = emptyMap()) : TreeNode(), MutableMap<String, TreeNode> by LinkedHashMap(values) {
    companion object {
        internal fun build(source: Map<*, *>) = MapNode(source.map { (key, value) -> Pair(key.toString(), TreeNode.build(value)) }.toMap())
    }

    override fun toString() = toMap().toString()
}

fun flatten(src: TreeNode?): Any? = when (src) {
    is MapNode -> flattenMap(src)
    is ArrayNode -> flattenArray(src)
    is PrimitiveNode -> src.value
    null -> null
}

internal fun flattenMap(src: MapNode): MutableMap<String, Any?> = LinkedHashMap(src.mapValues {
    when (it.value) {
        is MapNode, is ArrayNode -> flatten(it.value)
        else -> (it.value as PrimitiveNode).value
    }
})

internal fun flattenArray(src: ArrayNode): MutableList<Any?> = ArrayList(src.map {
    when (it) {
        is MapNode, is ArrayNode -> flatten(it)
        is PrimitiveNode -> it.value
    }
})
