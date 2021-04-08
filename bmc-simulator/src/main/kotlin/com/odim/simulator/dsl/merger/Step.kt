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

import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.ExtensibleEmbeddedObjectType
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.ObjectArray
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.ResourceType
import com.odim.simulator.tree.structure.Type
import com.odim.simulator.tree.structure.ValueArray

sealed class Step

sealed class Operation : Step()

sealed class Expression : Step()

class MergeSteps(val steps: MutableList<Operation> = mutableListOf()) : Operation() {
    fun add(step: Operation) {
        when (step) {
            is MergeSteps -> if (step.steps.isNotEmpty()) {
                steps.add(step)
            }
            else -> steps.add(step)
        }
    }
}

class SetConstant(val left: Property<Any?>, val right: Getter<Any?>) : Operation()

class SetResourceLink(val left: Property<LinkableResource>, val right: Expression) : Operation()

class TruncateObjectArray(val left: Property<ObjectArray>, val right: Getter<Int>) : Operation()

class TruncateResourceArray(val left: Property<EmbeddedResourceArray>, val right: Getter<Int>) : Operation()

class ExtendObjectArray(val left: Property<ObjectArray>, val right: Getter<List<ResourceObject>>) : Operation()

class ExtendResourceArray(val left: Property<EmbeddedResourceArray>, val right: Getter<List<ResourceObject>>) : Operation()

class UpdateValueArray(val left: Property<ValueArray>, val right: Getter<List<Any?>>) : Operation()

class RemoveProperty(val context: Value<Any?>, val key: String) : Operation()

class SetListOfResourceLinks(val left: Property<LinkableResourceArray>, val right: Expression) : Operation()

class LinkNeedsRemoval(val link: Getter<String>) : Operation()

class LinkNeedsAdding(val link: Getter<String>) : Operation()

class BuildResourceLink(val linkType: Getter<Type>, val linkUri: Getter<String?>) : Expression()

class BuildListOfResourceLinks(val linkType: Getter<Type>, val linkUris: Getter<List<String>>) : Expression()

class AlignEmbeddedObjects(val left: Value<MutableList<ResourceObject>>, val right: Getter<ExtensibleEmbeddedObjectType>) : Operation()

class AlignEmbeddedResourceArray(val left: Value<MutableList<ResourceObject>>, val right: Getter<ResourceType>) : Operation()
