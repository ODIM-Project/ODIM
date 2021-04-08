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

class EmbeddedResourceArray(val acceptType: ResourceType) : TreeElement(acceptType) {
    val data = mutableListOf<Resource>()

    fun add(resource: Resource): Boolean {
        if (!accepts(resource)) {
            throw IllegalArgumentException("EmbeddedResourceArray does not accept resource of type ${resource.meta.type}.! Desired type: $acceptType")
        }
        val added = data.add(resource)
        if (added) {
            resource.meta.parent = this.meta.parent
        }
        return added
    }

    fun remove(resource: Resource): Boolean {
        val removed = data.remove(resource)
        if (removed) {
            resource.meta.parent = null
        }
        return removed
    }

    private fun accepts(resource: Resource) = resource.meta.type == acceptType

    override fun print(start: Int) = this.data.toString()

    override fun toLink(): String {
        throw UnsupportedOperationException("EmbeddedResourceArray is only container. Use toLink of its contents!")
    }
}
