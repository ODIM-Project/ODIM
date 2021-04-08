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

import com.odim.simulator.dsl.DSL
import com.odim.simulator.tree.RedfishVersion
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.SERVICE_ROOT

@Suppress("TooManyFunctions")
open class Resource(resourceType: ResourceTypeBase, version: RedfishVersion? = null) : TreeElement(resourceType, version), Searchable {
    override fun findElement(searchPattern: Any?, other: String): Item? {
        val treeElements = searchForTreeElementsAsMap(this.data)
        return if (treeElements[searchPattern] is TreeElement) {
            treeElements[searchPattern] as TreeElement
        } else {
            // This condition is used to properly handle BMC issue, where link to contained Chassis is set by main Chassis
            // main Chassis -> Id RackMount, Link ../Chassis/RackMount
            // subordinate Chassis -> Id Baseboard, Link ../Chassis/RackMount/Baseboard
            if (this.meta.type == CHASSIS) {
                return if (chassisLinkPredicate(this.toLink(), other)) this else null
            }
            searchForTreeElementsAsMap(this.data).filter {
                it.value.meta.type.jsonName() == searchPattern
            }.values.firstOrNull()
        }
    }

    companion object {
        fun searchForActions(data: Map<*, *>): List<ActionElement> {
            val actions = mutableListOf<ActionElement>()
            for ((_, value) in data) {
                if (value is Actions) {
                    actions.addAll(value.getActions())
                } else if (value is Map<*, *>) {
                    actions.addAll(searchForActions(value))
                }
            }
            return actions
        }

        fun searchForTreeElementsAsMap(data: Map<*, *>): Map<String, TreeElement> {
            val elements = mutableMapOf<String, TreeElement>()
            for ((prop, value) in data) {
                if (value is TreeElement) {
                    elements[prop as String] = value
                } else if (value is Map<*, *>) {
                    elements.putAll(searchForTreeElementsAsMap(value))
                }
            }
            return elements
        }

        fun searchForTreeElementsAsList(data: Map<*, *>): List<TreeElement> {
            val elements = mutableListOf<TreeElement>()
            for ((prop, value) in data) {
                if (value is TreeElement) {
                    elements.add(value)
                } else if (value is Map<*, *>) {
                    elements.addAll(searchForTreeElementsAsList(value))
                }
            }
            return elements
        }

        @Suppress("NestedBlockDepth")
        fun collectNotIndexedResources(data: Map<*, *>): List<Resource> {
            val elements = mutableListOf<Resource>()
            for ((_, value) in data) {
                if (value is ResourceCollection && !value.onlyLinks) {
                    value.members.forEach {
                        if (it.data["Id"] == 0) {
                            elements.add(it)
                        }
                        elements.addAll(collectNotIndexedResources(it.data))
                    }
                } else if (value is Map<*, *>) {
                    elements.addAll(collectNotIndexedResources(value))
                }
            }
            return elements
        }

        fun searchForExapndableElements(data: Map<*, *>): Map<String, Any> {
            val elements = mutableMapOf<String, Any>()
            for ((prop, value) in data) {
                if (expandable(value)) {
                    elements[prop as String] = value!!
                } else if (value is Map<*, *>) {
                    elements.putAll(searchForExapndableElements(value))
                }
            }
            return elements
        }

        fun resourceObject(): ResourceObject = mutableMapOf()

        fun resourceObject(vararg properties: Pair<String, Any?>): ResourceObject = mutableMapOf(*properties)

        fun embeddedObject(type: ExtensibleEmbeddedObjectType): EmbeddedObject = EmbeddedObject(type)

        fun embeddedObject(vararg properties: Pair<String, Any?>): ResourceObject = mutableMapOf(*properties)

        fun embeddedArray(type: ExtensibleEmbeddedObjectType): EmbeddedArray<ResourceObject> = ObjectArray(type)

        fun embeddedArray(vararg elements: Any): EmbeddedArray<Any?> = ValueArray(elements.toList())

        private fun expandable(value: Any?) =
                value is Resource || value is ResourceCollection || value is LinkableResource || value is LinkableResourceArray
    }

    var data = resourceObject()

    private var odataCreateStrategy: (RedfishVersion) -> String? = { _ -> null}

    object NoProperty

    fun append(resource: Resource, property: String? = null) = append(this, resource, property)

    fun appendTo(resource: Resource, property: String? = null) = append(resource, this, property)

    operator fun invoke(vararg attachments: Resource): Resource {
        attachments.forEach { append(it) }
        return this
    }

    operator fun invoke(vararg resources: Collection<Resource>): Resource {
        resources.flatMap { it }.forEach { append(it) }
        return this
    }

    operator fun invoke(builder: DSL.() -> Unit): Resource {
        val dsl = DSL()
        dsl.builder()
        dsl.applyTo(this.data)
        return this
    }

    override fun getOdataValue() = odataCreateStrategy(this.meta.createdVersion!!)
            ?: "#${this.meta.type.oDataType()}.${this.meta.createdVersion}.${this.meta.type.oDataType()}"

    override fun getOdataContext() = odataCreateStrategy(this.meta.createdVersion!!)
            ?: "/redfish/v1/${'$'}metadata#${this.meta.type.oDataType()}.${this.meta.type.oDataType()}"

    fun setOdataCreateStrategy(strategy: ((RedfishVersion) -> String?)?) {
        strategy?.let { this.odataCreateStrategy = strategy }
    }

    override fun toLink(): String {
        if (this.meta.parent == null) {
            var baseUri = "/redfish/v1"
            if (this.isOemResource) {
                baseUri += this.oemLinkPrefix
            }
            return baseUri
        }

        val jsonName = this.meta.type.jsonName()
        val link = if (!jsonName.isEmpty()) {
            jsonName
        } else {
            this.data["Id"]
        }
        return "${this.meta.parent!!.toLink()}${if (this.meta.type.isEmbedded()) "#" else ""}/$link"
    }

    fun getCollection(type: ResourceCollectionType) = searchForTreeElementsAsMap(this.data)
            .filter { it.value.meta.type == type && it.value is ResourceCollection }
            .map { it.value as ResourceCollection }
            .firstOrNull() ?: throw NoSuchElementException("Resource $this does not contain collection of type $type")

    fun getSingleton(name: String) = traverse<Resource>(name)

    inline fun <reified T : Any?> traverse(path: String): T {
        val pathChunked = path.split("/")
        var result: Any? = this.data
        for (part in pathChunked.filter { it.isNotEmpty() }) {
            result = this.lookForProperty(result, part)
            if (result == NoProperty) {
                throw TreeTraversalException("Traversing $this for $path failed: part \"$part\" not found!")
            }
        }
        return result as T
    }

    @Suppress("SwallowedException")
    inline fun <reified T : Any?> traverseOrNull(path: String): T? = try {
        traverse<T>(path)
    } catch (e: TreeTraversalException) {
        null
    }

    fun lookForProperty(obj: Any?, prop: String) = when (obj) {
        is Map<*, *> -> obj.getOrDefault(prop, NoProperty)
        is List<*> ->
            obj.getOrElse(prop.toIntOrNull() ?: throw IllegalArgumentException("Wrong path specified: non integer array index!")) { NoProperty }
        else -> throw TreeTraversalException("Encountered unexpected type to traverse: ${obj?.javaClass}!")
    }

    /**
     * Method to append child Resource to parent CollectionResource that accepts it. It will append to the first found
     * valid collection unless you provide property parameter. If property is provided child Resource will be appended
     * to the collection under specified key (in template)
     */
    private fun append(parent: Resource, child: Resource, property: String?) =
            searchForTreeElementsAsMap(parent.data).filter { it.value is ResourceCollection }
                    .filter { property == null || it.key == property }
                    .map { it.value as ResourceCollection }
                    .firstOrNull { it.type.of() == child.meta.type }
                    ?.let {
                        if (!it.onlyLinks && child.meta.parent != null) {
                            throw IllegalArgumentException("Resource of type $child is already appended to resource of type ${child.meta.parent}!")
                        }
                        it.add(child)
                        parent
                    } ?: throw NoSuchElementException("Cannot append $child to $parent!")

    fun isReachableFromRoot(): Boolean {
        var resource: TreeElement? = this
        while (resource != null) {
            if (resource.meta.type == SERVICE_ROOT) {
                return true
            }
            resource = resource.meta.parent
        }
        return false
    }

    override fun print(start: Int): String {
        val indent = "  ".repeat(start)
        var print = "${this.meta.type} {"
        this.data.forEach { prop, value ->
            print += if (value is TreeElement) {
                "\n  $indent$prop:${value.print(start + 1)}"
            } else {
                "\n  $indent$prop: $value"
            }
        }
        print += "\n$indent}"
        return print
    }

    override fun toString(): String {
        if (this.data.containsKey("Id")) {
            return "${this.meta.type}[ID=${this.data["Id"]}]"
        }
        return this.meta.type.toString()
    }

    fun getId() = "${this.data["Id"]}"
}
