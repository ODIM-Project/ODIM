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

import com.odim.simulator.dsl.Marker.ARRAY
import com.odim.simulator.dsl.Marker.EMPTY
import com.odim.simulator.dsl.merger.MergeExecutor
import com.odim.simulator.dsl.merger.Plan
import java.util.LinkedList

enum class Marker {
    EMPTY,
    ARRAY
}

class DSL {
    companion object {
        fun mergeMaps(dest: MutableMap<String, Any?>, src: Map<String, Any?>) {
            val dsl = DSL()
            dsl.treeRoot.putAll(TreeNode.build(src) as MapNode)
            dsl.applyTo(dest)
        }

        fun makeTypeSafeMap(src: Map<*, *>) = src.map { (key, value) -> Pair(key.toString(), value) }.toMap()
    }

    private val valueStack = LinkedList<TreeNode>()
    private val treeRoot = MapNode()
    private val treeStack = LinkedList(listOf(treeRoot))

    val empty = EMPTY
    val array get() = this

    fun toMap() = flattenMap(treeRoot)

    fun isMap() = treeStack.last() is MapNode

    fun applyTo(target: MutableMap<String, Any?>) {
        if (!isMap()) {
            throw UnsupportedOperationException("DSL does not render to a map")
        }
        fun applyMap(dest: MutableMap<String, Any?>, src: MapNode) {
            val plan = Plan.Builder().base(dest).override(src).build()
            MergeExecutor.buildNonLinking().execute(plan)
        }
        applyMap(target, treeRoot)
    }

    infix fun String.to(nested: () -> Any) {
        val newNode = MapNode()
        (treeStack.last() as MapNode)[this] = newNode
        treeStack.add(newNode)
        nested()
        treeStack.removeLast()
    }

    infix fun String.to(item: Marker) = when (item) {
        EMPTY -> throw IllegalArgumentException("can't map $this to <empty>, <empty> can exist only in array")
        ARRAY -> (treeStack.last() as MapNode)[this] = valueStack.removeLast()
    }

    infix fun String.to(item: Any?) {
        (treeStack.last() as MapNode)[this] = PrimitiveNode(item)
    }

    operator fun get(vararg items: Any): Marker {
        val resultArray = ArrayNode()
        val singleArgument = items.singleOrNull()
        when (singleArgument) {
            EMPTY -> Unit
            is Sequence<*> -> wrapSequenceItemsIntoArray(singleArgument, resultArray)
            else -> wrapItemsIntoArray(items, resultArray)
        }
        valueStack.addLast(resultArray)
        return ARRAY
    }

    private fun wrapSequenceItemsIntoArray(items: Sequence<Any?>, target: ArrayNode) = wrapItemsIntoArray(
            items.filterNotNullTo(mutableListOf()).toTypedArray(), target)

    private fun wrapItemsIntoArray(items: Array<out Any>, target: ArrayNode) = items.forEach {
        when {
            it == EMPTY -> throw IllegalArgumentException("can't process array of values, <empty> should be the sole element of array")
            it == ARRAY -> throw UnsupportedOperationException("can't process nested array yet")
            uglyIsLambda(it) -> {
                val newNode = MapNode()
                treeStack.add(newNode)
                uglyInvokeLambda(it)
                treeStack.removeLast()
                target.add(newNode)
            }
            it is Map<*, *> ->
                target.add(MapNode.build(it))
            else -> target.add(PrimitiveNode(it))
        }
    }

    private fun uglyIsLambda(value: Any?) = value != null && (value::class.java.isLocalClass || value::class.java.isAnonymousClass)

    @Suppress("SwallowedException", "TooGenericExceptionCaught")
    private fun uglyInvokeLambda(value: Any) {
        for (m in value::class.java.methods.filter { it.name == "invoke" }) {
            try {
                val isAccessible = m.isAccessible
                m.isAccessible = true
                try {
                    m.invoke(value)
                } finally {
                    m.isAccessible = isAccessible
                }
                break
            } catch (e: Exception) {
                continue
            }
        }
    }
}
