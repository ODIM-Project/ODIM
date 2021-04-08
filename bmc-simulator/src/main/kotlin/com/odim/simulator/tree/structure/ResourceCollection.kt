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

package com.odim.simulator.tree.structure

import com.odim.simulator.tree.TreeBuildingException
import com.odim.simulator.tree.structure.Resource.Companion.collectNotIndexedResources
import java.util.Collections.synchronizedList

open class ResourceCollection(open val type: ResourceCollectionTypeBase, open val jsonName: String? = null, open val onlyLinks: Boolean = false,
                         open var autoExpand: Boolean = false) : TreeElement(type), Searchable {

    override fun findElement(searchPattern: Any?, other: String) = findElementByPart(searchPattern, other)

    val linkPart get() = jsonName ?: this.meta.type.jsonName()
    val members: MutableList<Resource> = synchronizedList(mutableListOf<Resource>())

    open fun findElementByPart(searchPattern: Any?, path: String)
            = members.firstOrNull { it.data["Id"].toString() == searchPattern }

    fun add(child: Resource): ResourceCollection {
        if (type.of() == child.meta.type) {
            members.add(child)
            if (!onlyLinks) {
                child.meta.parent = this
                if (child.isReachableFromRoot()) {
                    indexTreePart(child)
                }
            }
            return this
        }
        throw TreeBuildingException("Cannot append $child to $this!")
    }

    fun remove(child: Resource): ResourceCollection {
        members.remove(child)
        return this
    }

    override fun getOdataValue() = "#${this.meta.type.oDataType()}.${this.meta.type.oDataType()}"
    override fun getOdataContext() = "/redfish/v1/${'$'}metadata${this.getOdataValue()}"

    fun size() = members.size

    override fun toLink(): String = "${this.meta.parent!!.toLink()}/$linkPart"

    override fun print(start: Int): String {
        val indent = "  ".repeat(start)
        var print = "${this.meta.type} {"
        this.members.forEachIndexed { index, resource -> print += "\n  $indent$index:${resource.print(start + 1)}" }
        print += "\n$indent}"
        return print
    }

    override fun toString(): String = "${this.meta.type} of ${this.meta.parent}"

    private fun indexTreePart(resource: Resource) {
        mutableListOf(resource).also {
            it.addAll(collectNotIndexedResources(resource.data))
        }.forEach {
            if (it.data.containsKey("Id") && it.data["Id"] == 0) {
                it.data["Id"] = this.meta.resourceFactory!!.generateId(it)
            }
        }
    }
}
