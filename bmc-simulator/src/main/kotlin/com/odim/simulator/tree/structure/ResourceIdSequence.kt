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

class ResourceIdSequence {

    private class Sequence {
        fun findFirstAvailable(reservedIds: List<Int>): Int {
            return generateSequence(0) { it + 1 }
                    .takeWhile { it <= (reservedIds.max()?:0)+1 }
                    .filter { !reservedIds.contains(it) }
                    .first()
        }
    }

    private val idGenerators: MutableMap<String, Sequence> = HashMap()

    fun next(resource: Resource, resourceCollection: ResourceCollection): Int {
        val reservedIds: List<Int> = resourceCollection.members.mapNotNull {
            val element = it.data["Id"]
            if (element is String) {
                element.toIntOrNull()
            } else {
                element as Int
            }
        }

        val sequencer = "${resourceCollection.toLink()}/${resource.meta.type}"
        if (!this.idGenerators.containsKey(sequencer)) {
            this.idGenerators[sequencer] = Sequence()
        }
        return this.idGenerators[sequencer]!!.findFirstAvailable(reservedIds)
    }
}
