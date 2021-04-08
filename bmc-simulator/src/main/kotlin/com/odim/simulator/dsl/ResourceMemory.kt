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

import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.TreeElement

class ResourceMemory {
    var open = false
        private set
    private val memory = mutableListOf<Resource>()

    fun open() {
        memory.clear()
        open = true
    }

    fun close() {
        open = false
    }

    fun record(resource: Resource) {
        memory.add(resource)
    }

    fun fetchOrderedResources() = MemorySorter(memory).sort()
}

private class MemorySorter(val unsorted: List<Resource>) {
    private val memoryRoot = MemoNode(null)
    private val sorted = mutableListOf<MemoNode>()
    private val alreadyAdded = mutableSetOf<TreeElement>()

    fun sort(): List<Resource> {
        unsorted.forEach(::addToTree)
        treeWalk(memoryRoot).forEach { sorted.add(it) }
        return sorted.mapNotNull { unsorted.elementAtOrNull(unsorted.indexOf(it.item)) }.toList()
    }

    private fun buildMemoryPath(item: Resource): List<MemoNode> {
        val path = mutableListOf<MemoNode>()
        var resource: TreeElement? = item
        while (resource != null) {
            path.add(0, MemoNode(resource, unsorted.indexOf(resource), path.take(1).toMutableList()))
            resource = resource.meta.parent
        }
        return path
    }

    private fun addToTree(item: Resource) {
        if (item in alreadyAdded) {
            return
        }
        val path = buildMemoryPath(item)
        var node = memoryRoot
        for (segment in path) {
            node = node.findChild(segment) ?: node.addAndGet(segment)
            segment.item?.let(alreadyAdded::add)
        }
    }

    private fun treeWalk(startFrom: MemoNode): Iterator<MemoNode> = object : AbstractIterator<MemoNode>() {
        val stack = mutableListOf(startFrom)
        override fun computeNext() = when {
            stack.isEmpty() -> done()
            else -> {
                val node = stack.removeAt(0)
                stack.addAll(0, node.children)
                setNext(node)
            }
        }
    }
}

private class MemoNode(val item: TreeElement?, var order: Int = -1, val children: MutableList<MemoNode> = mutableListOf()) {
    fun findChild(reference: MemoNode) = children.firstOrNull { it.item == reference.item }
    fun addAndGet(child: MemoNode) = child.also { children.add(it) }
}
