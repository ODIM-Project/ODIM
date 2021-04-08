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

import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.ExtensibleEmbeddedObjectType
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.ObjectArray
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.ResourceType
import com.odim.simulator.tree.structure.TreeElement
import com.odim.simulator.tree.structure.Type
import com.odim.simulator.tree.structure.ValueArray

class MergeExecutor private constructor(vararg linkProcessors: PlanProcessor) {
    private var processorChain: PlanProcessor = listOf(ValueArrayProcessor(), ObjectArrayProcessor(), *linkProcessors, DefaultsProcessor())
            .reduce(PlanProcessor?::join)

    companion object {
        fun buildNonLinking(): MergeExecutor = MergeExecutor(LinkAgnosticProcessor(), LinkEvaluator(), EmbeddedObjectAgnosticProcessor())

        fun buildLinking(tree: ResourceTree, targetUri: String): MergeExecutor = MergeExecutor(LinkRemovingProcessor(tree, targetUri),
                LinkEvaluator(), EmbeddedObjectProcessor(tree), EmbeddedResourceArrayProcessor(tree))
    }

    inner class Execution {
        fun visit(step: Step) {
            when (step) {
                is MergeSteps -> step.steps.forEach(::visit)
                is Operation -> perform(step)
            }
        }

        private fun perform(operation: Operation) {
            processorChain.perform(operation)
        }

    }

    fun execute(plan: Step) {
        Execution().visit(plan)
    }
}

private fun PlanProcessor?.join(that: PlanProcessor): PlanProcessor = when (this) {
    null -> that
    else -> this.also { it.last.next = that }
}

private interface PlanProcessor {
    var prev: PlanProcessor?
    var next: PlanProcessor?
    val first: PlanProcessor
        get() {
            var that = this
            var then = that.prev
            do {
                that = then ?: break
                then = that.prev
            } while (then != null)
            return that
        }

    val last: PlanProcessor
        get() {
            var that = this
            var then = that.next
            do {
                that = then ?: break
                then = that.next
            } while (then != null)
            return that
        }
    fun perform(operation: Operation)
    fun <T> evaluate(step: Expression): T?
}

private abstract class PlanProcessorBase : PlanProcessor {
    override var prev: PlanProcessor? = null

    override var next: PlanProcessor? = null
        set(value) {
            if (value != field) {
                field = value
                field?.prev = this
            }
        }

    override fun perform(operation: Operation) = next?.perform(operation) ?: Unit

    override fun <T> evaluate(step: Expression): T? = next?.evaluate(step)
}

private class LinkRemovingProcessor(private val tree: ResourceTree, private val targetUri: String) : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is SetResourceLink -> operation.run {
                (first.evaluate(right) as LinkableResource?)?.let { left.set(it) }
            }
            is SetListOfResourceLinks -> operation.run {
                (first.evaluate(right) as LinkableResourceArray?)?.let { left.set(it) }
            }
            is LinkNeedsRemoval -> optRemoveLinkFrom(targetUri, operation.link.get() as String)
            is LinkNeedsAdding -> optAddingLink(targetUri, operation.link.get() as String)
            else -> next?.perform(operation)
        }
    }

    private fun optAddingLink(targetUri: String, sourceUri: String?) {
        if (sourceUri == null) {
            return
        }
        tree.link(requireLinkedResource(targetUri), requireLinkedResource(sourceUri))
    }

    private fun optRemoveLinkFrom(targetUri: String, sourceUri: String?) {
        if (sourceUri == null) {
            return
        }
        tree.removeLink(requireLinkedResource(targetUri), requireLinkedResource(sourceUri))
    }

    private fun requireLinkedResource(linkUri: String): TreeElement {
        return tree.search(linkUri) as? TreeElement
                ?: throw MergeException("Missing linked resource '$linkUri'")
    }

}

private class LinkAgnosticProcessor : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is SetResourceLink -> operation.run {
                (first.evaluate(right) as LinkableResource?)?.let { left.set(it) }
            }
            is SetListOfResourceLinks -> operation.run {
                (first.evaluate(right) as LinkableResourceArray?)?.let { left.set(it) }
            }
            is LinkNeedsRemoval, is LinkNeedsAdding -> Unit
            else -> next?.perform(operation)
        }
    }
}

private class LinkEvaluator : PlanProcessorBase() {
    override fun <T> evaluate(step: Expression): T? {
        return when (step) {
            is BuildResourceLink -> buildLinkableResource(step.linkType.get()) as? T
            is BuildListOfResourceLinks -> buildLinkableResourceCollection(step.linkType.get()) as? T
        }
    }

    private fun buildLinkableResource(linkType: Type?): LinkableResource {
        if (linkType == null) {
            throw MergeException("Unspecified Resource type to create a link")
        }
        return LinkableResource(linkType)
    }

    private fun buildLinkableResourceCollection(linkType: Type?): LinkableResourceArray {
        if (linkType == null) {
            throw MergeException("Unspecified Resource type to create a link array")
        }
        return LinkableResourceArray(linkType)
    }
}

private class EmbeddedObjectProcessor(val tree: ResourceTree) : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is AlignEmbeddedObjects -> alignEmbeddedObjects(operation.left.get() as List<ResourceObject>, operation.right.get() as ExtensibleEmbeddedObjectType)
            else -> next?.perform(operation)
        }
    }

    private fun alignEmbeddedObjects(target: List<ResourceObject>, embeddedObjectType: ExtensibleEmbeddedObjectType) {
        target.forEach { item ->
            val reference = tree.createEmbeddedObject(embeddedObjectType)
            (reference.keys - item.keys).forEach { key ->
                item[key] = reference[key]
            }
        }
    }
}

private class EmbeddedResourceArrayProcessor(val tree: ResourceTree) : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is AlignEmbeddedResourceArray -> alignEmbeddedResourceArray(operation.left.get() as List<ResourceObject>, operation.right.get() as ResourceType)
            is ExtendResourceArray -> extendResourceArray(operation.left.get() as EmbeddedResourceArray, operation.right.get() as List<ResourceObject>)
            is TruncateResourceArray -> truncateResourceArray(operation.left.get() as EmbeddedResourceArray, operation.right.get() as Int)
            else -> next?.perform(operation)
        }
    }

    private fun alignEmbeddedResourceArray(target: List<ResourceObject>, embeddedObjectType: ResourceType) {
        target.forEach { item ->
            val reference = tree.create(embeddedObjectType)
            (reference.data.keys - item.keys).forEach { key ->
                item[key] = reference.data[key]
            }
        }
    }

    private fun extendResourceArray(objectArray: EmbeddedResourceArray, tailObjects: List<ResourceObject>) {
        for (m in tailObjects) {
            val resource = tree.create(objectArray.acceptType)
            resource.data = m
            objectArray.add(resource)
        }
    }

    private fun truncateResourceArray(objectArray: EmbeddedResourceArray, targetSize: Int) {
        objectArray.run {
            data.subList(targetSize, data.size).clear()
        }
    }
}

private class EmbeddedObjectAgnosticProcessor : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is AlignEmbeddedObjects -> Unit
            else -> next?.perform(operation)
        }
    }
}

private class ObjectArrayProcessor : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is TruncateObjectArray -> truncateObjectArray(operation.left.get() as ObjectArray, operation.right.get() as Int)
            is ExtendObjectArray -> extendObjectArray(operation.left.get() as ObjectArray, operation.right.get() as List<ResourceObject>)
            else -> next?.perform(operation)
        }
    }

    private fun extendObjectArray(objectArray: ObjectArray, tailObjects: List<ResourceObject>) {
        objectArray.addAll(tailObjects)
    }

    private fun truncateObjectArray(objectArray: ObjectArray, targetSize: Int) {
        objectArray.run {
            subList(targetSize, size).clear()
        }
    }
}

private class ValueArrayProcessor : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is UpdateValueArray -> updateValueArray(operation.left.get() as ValueArray, operation.right.get() as List<Any?>)
            else -> next?.perform(operation)
        }
    }

    private fun updateValueArray(valueArray: ValueArray, src: List<Any?>) {
        valueArray.clear()
        valueArray.addAll(src)
    }
}

private class DefaultsProcessor : PlanProcessorBase() {
    override fun perform(operation: Operation) {
        when (operation) {
            is SetConstant -> operation.left.set(operation.right.get())
            is RemoveProperty -> (operation.context.get() as HashMap<*, *>).remove(operation.key)
            else -> throw MergeException("Don't know how to handle operation $operation")
        }
    }

    override fun <T> evaluate(step: Expression): T? {
        throw MergeException("Don't know how to evaluate expression $step")
    }
}
