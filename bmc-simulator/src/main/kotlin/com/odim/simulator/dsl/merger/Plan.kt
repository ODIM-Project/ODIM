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
import com.odim.simulator.dsl.ArrayNode
import com.odim.simulator.dsl.DSL
import com.odim.simulator.dsl.MapNode
import com.odim.simulator.dsl.PrimitiveNode
import com.odim.simulator.dsl.TreeNode
import com.odim.simulator.dsl.flatten
import com.odim.simulator.dsl.flattenArray
import com.odim.simulator.dsl.merger.Plan.ModifyMode.APPLY
import com.odim.simulator.dsl.merger.Plan.ModifyMode.REMOVE
import com.odim.simulator.dsl.merger.Plan.ModifyMode.SKIP
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.ExtensibleEmbeddedObjectType
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.ObjectArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.ResourceType
import com.odim.simulator.tree.structure.SingletonResource
import com.odim.simulator.tree.structure.TreeElement
import com.odim.simulator.tree.structure.ValueArray
import com.odim.utils.toMap

class PlannerBuildingException(message: String) : MergeException(message)
class ResourceOverrideException(msg: String, cause: Throwable? = null) : MergeException(msg, cause)

@Suppress("TooManyFunctions")
class Plan private constructor(private val base: Map<String, Any?>, private val override: TreeNode, private val config: (MergeDslConfiguration.() -> Unit)) {

    private enum class ModifyMode {
        APPLY,
        REMOVE,
        SKIP
    }

    private val baseContext: Value<Any?>
    private val dslConfig = MergeDslConfiguration()
    private val plan: MergeSteps = MergeSteps()

    init {
        dslConfig.config()
        baseContext = Value(base)
        plan()
    }

    @Suppress("UnnecessaryApply")
    data class Builder(
            var config: (MergeDslConfiguration.() -> Unit)? = null,
            var base: Map<String, Any?>? = null,
            var override: TreeNode? = null) {

        fun base(base: Map<String, Any?>) = apply { this.base = base }
        fun override(override: TreeNode) = apply { this.override = override }
        fun override(override: ObjectNode) = apply { this.override = TreeNode.build(override.toMap()) }
        fun withConfig(config: (MergeDslConfiguration.() -> Unit)?) = apply { this.config = config }
        fun build(): MergeSteps {
            base ?: throw PlannerBuildingException("Base resource for merger plan cannot be null")
            override ?: throw PlannerBuildingException("Override resource for merger plan cannot be null")
            return Plan(base!!, override!!, config ?: {}).plan
        }
    }

    private fun plan() {
        (override as MapNode).forEach { key, value ->
            when (modifyMode(key, value)) {
                REMOVE -> planForPropertyRemoval(key.substringBefore("@"))
                APPLY -> planForModifications(key, value)
                SKIP -> return@forEach
            }
        }
        removeAbsentPropertiesOnBaseStructure()
    }

    private fun modifyMode(key: String, value: TreeNode): ModifyMode {
        try {
            return when {
                key.endsWith("@removed") -> checkValueForRemoved(key, value)
                base[key] is SingletonResource -> checkFailForNonMergeable(base[key])
                base[key] is ResourceCollection -> checkFailForNonMergeable(base[key])
                base[key] is Resource -> checkFailForNonMergeable(base[key])
                base[key] is Actions -> checkFailForNonMergeable(base[key])
                else -> APPLY
            }
        } catch (e: MergeException) {
            throw MergeException("Can't override property '$key' of resource '${extractId(base)}'", e)
        }
    }

    private fun checkValueForRemoved(key: Any?, value: TreeNode): ModifyMode {
        if (value is PrimitiveNode && value.value == true) return REMOVE
        throw ResourceOverrideException("Key indicating property removal ($key) can have only 'true' value")
    }

    private fun checkFailForNonMergeable(target: Any?): ModifyMode {
        if (dslConfig.failForNonMergeable) throw ResourceOverrideException("Property type ${target?.javaClass?.simpleName} cannot be overridden")
        return SKIP
    }

    private fun extractId(resourceData: Map<String, Any?>) = DSL.makeTypeSafeMap(resourceData).let {
        it["@odata.id"]?.toString() ?: """${it["Name"]}(${it["Id"]})/${it["UUID"]}"""
    }

    private fun planForModifications(key: String, value: TreeNode) {
        when (base[key]) {
            is LinkableResource -> planForLinkableResource(MapProperty(key, baseContext), flatten(value))
            is LinkableResourceArray -> planForLinkableResourceCollection(MapProperty(key, baseContext), flattenArray(ensureArrayNode(value)))
            is ObjectArray -> planForObjectArray(MapProperty(key, baseContext), ensureArrayNode(value))
            is EmbeddedResourceArray -> planEmbeddedResourceArray(MapProperty(key, baseContext), ensureArrayNode(value))
            is ValueArray -> plan.add(planForValueArray(MapProperty(key, baseContext), ensureArrayNode(value)))
            is MutableMap<*, *> -> plan.add(Plan.Builder(config)
                    .base(base[key] as? Map<String, Any?> ?: mutableMapOf())
                    .override(ensureMapNode(value))
                    .build())
            else -> plan.add(planForValue(MapProperty(key, baseContext), flatten(value)))
        }
    }

    private fun removeAbsentPropertiesOnBaseStructure() {
        if (dslConfig.removeAbsentProperties) {
            base.filter {
                !(override as MapNode).containsKey(it.key)
            }.filter {
                it.key != "Id"
            }.forEach {
                planDiscardLinks(makeMapProperty(it.key, baseContext))
                planForPropertyRemoval(it.key)
            }
        }
    }

    private fun planForValueArray(dest: Property<ValueArray>, node: ArrayNode) = UpdateValueArray(dest, Value(flattenArray(node)))

    private fun planForObjectArray(dest: Property<ObjectArray>, node: ArrayNode) {
        val destArray = dest.get() ?: return
        val tailObjects = mutableListOf<ResourceObject>()
        destArray.forEachIndexed { index, _ -> planDiscardLinks(ListProperty<Any?>(index, Value(dest))) }
        plan.add(TruncateObjectArray(dest, Value(0)))
        node.forEach { tailObjects.add(mutableMapOf<String, Any?>()) }
        plan.add(ExtendObjectArray(dest, Value(tailObjects)))
        node.forEachIndexed { index, childrenNode ->
            val arrayItem = ListProperty<Map<String, Any?>>(index, Value(tailObjects))
            plan.add(Plan.Builder().withConfig(config).base(arrayItem.get() ?: mutableMapOf())
                    .override(childrenNode as MapNode)
                    .build())
        }
        plan.add(planAlignEmbeddedObject(Value(tailObjects), Value(destArray.templateType)))
    }

    private fun planEmbeddedResourceArray(dest: Property<EmbeddedResourceArray>, node: ArrayNode) {
        val destArray = dest.get() ?: return
        val tailObjects = mutableListOf<ResourceObject>()
        destArray.data.forEachIndexed { index, _ -> planDiscardLinks(ListProperty<Any?>(index, Value(dest))) }
        plan.add(TruncateResourceArray(dest, Value(0)))
        node.forEach { tailObjects.add(mutableMapOf<String, Any?>()) }
        plan.add(ExtendResourceArray(dest, Value(tailObjects)))
        node.forEachIndexed { index, childrenNode ->
            val arrayItem = ListProperty<Map<String, Any?>>(index, Value(tailObjects))
            plan.add(Plan.Builder().withConfig(config).base(arrayItem.get() ?: mutableMapOf())
                    .override(childrenNode as MapNode)
                    .build())
        }
        plan.add(planAlignEmbeddedResourceArray(Value(tailObjects), Value(destArray.acceptType)))
    }

    private fun planAlignEmbeddedResourceArray(dest: Value<MutableList<ResourceObject>>, templateType: Value<ResourceType>): Operation =
            AlignEmbeddedResourceArray(dest, templateType)

    private fun planAlignEmbeddedObject(dest: Value<MutableList<ResourceObject>>, templateType: Value<ExtensibleEmbeddedObjectType>): Operation =
            AlignEmbeddedObjects(dest, templateType)

    private fun makeMapProperty(key: String, context: Getter<Any?>): MapProperty<Any?> = MapProperty(key, context)

    private fun makeListProperty(index: Int, context: Getter<Any?>): ListProperty<Any?> = ListProperty(index, context)

    private fun planDiscardLinks(dest: Property<*>) {
        when (val destSlot = dest.get()) {
            is Map<*, *> -> destSlot.keys.map { key -> makeMapProperty(key.toString(), Value(destSlot)) }.forEach { property ->
                planDiscardLinks(property)
            }
            is List<*> -> destSlot.indices.map { index -> makeListProperty(index, Value(destSlot)) }.forEach { property ->
                planDiscardLinks(property)
            }
            is LinkableResource -> planForLinkableResource(dest as Property<LinkableResource>, null)
            is LinkableResourceArray -> planForLinkableResourceCollection(dest as Property<LinkableResourceArray>, emptyList())
        }
    }

    private fun ensureMapNode(node: TreeNode) = when (node) {
        is MapNode -> node
        else -> throw MergeException("Expected object value")
    }

    private fun ensureArrayNode(node: TreeNode) = when (node) {
        is ArrayNode -> node
        else -> throw MergeException("Expected array value")
    }

    private fun planForValue(dest: Property<Any?>, srcValue: Any?) = SetConstant(dest, Value(srcValue))

    private fun planForPropertyRemoval(key: String) = plan.add(RemoveProperty(baseContext, key))

    private fun extractLinkUri(srcValue: Any?) = when (srcValue) {
        is Map<*, *> -> DSL.makeTypeSafeMap(srcValue)
        else -> emptyMap()
    }.mapValues { (_, v) -> v?.toString() }.getOrDefault("@odata.id", null)

    private fun planForLinkableResource(dest: Property<LinkableResource>, srcValue: Any?) {
        val oldLink = dest.get() ?: return
        val oldLinkUri = oldLink.getElement()?.toLink()
        val newLinkUri = extractLinkUri(srcValue)
        oldLinkUri?.filter { newLinkUri != oldLinkUri }?.let { link -> plan.add(LinkNeedsRemoval(Value(link))) }
        plan.add(SetResourceLink(dest, BuildResourceLink(Value(oldLink.type), Value(newLinkUri))))
        newLinkUri?.filter { newLinkUri != oldLinkUri }?.let { link -> plan.add(LinkNeedsAdding(Value(link))) }
    }

    private fun planForLinkableResourceCollection(dest: Property<LinkableResourceArray>, srcValue: List<Any?>) {
        val linksCollection = dest.get() ?: return
        val linksToRemove = linksCollection.getElements().map(TreeElement::toLink)
        val linksToAdd = srcValue.map(::extractLinkUri).mapNotNull { it }
        linksToRemove.forEach { plan.add(LinkNeedsRemoval(Value(it))) }
        plan.add(SetListOfResourceLinks(dest, BuildListOfResourceLinks(Value(linksCollection.type), Value(linksToAdd))))
        linksToAdd.forEach { plan.add(LinkNeedsAdding(Value(it))) }
    }
}
