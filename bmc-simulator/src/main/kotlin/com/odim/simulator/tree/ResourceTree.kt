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

package com.odim.simulator.tree

import com.odim.simulator.dsl.DSL
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.ExtensibleEmbeddedObjectType
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.SERVICE_ROOT
import com.odim.simulator.tree.structure.ResourceTypeBase
import com.odim.simulator.tree.structure.Searchable
import com.odim.simulator.tree.structure.TreeElement
import kotlin.reflect.KClass

class ResourceTree(private val resourceFactory: ResourceFactory) {
    private val serviceRoot = create(SERVICE_ROOT)

    val root get() = serviceRoot

    fun link(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            ResourceLinker.link(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun oneWayLink(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null) =
            ResourceLinker.oneWayLink(resourceFirst, resourceSecond, propertyFirst)

    fun removeLink(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            ResourceLinker.removeLink(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun areLinked(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            ResourceLinker.areLinked(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun areLinkedOnAnySide(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            ResourceLinker.areLinkedOnAnySide(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun oneWayLinkExist(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null) =
            ResourceLinker.checkIfOneSideLinkExist(resourceFirst, resourceSecond, propertyFirst)

    fun create(type: ResourceTypeBase, version: RedfishVersion? = null, oemResource: Boolean = false, override: (DSL.() -> Unit)? = null) =
            resourceFactory.create(type, version, override, oemResource)

    fun createEmbeddedObject(type: ExtensibleEmbeddedObjectType, version: ResourceVersion? = null, override: ((DSL).() -> Unit)? = null) =
            resourceFactory.createEmbeddedObject(type, version, override)

    fun setTemplate(type: ResourceTypeBase, templateClass: KClass<out ResourceTemplate>) {
        resourceFactory.setTemplate(type, templateClass)
        if (type == SERVICE_ROOT) {
            recreateRoot()
        }
    }

    fun setTemplateEmbedded(type: ExtensibleEmbeddedObjectType, templateClass: KClass<out ResourceTemplate>) {
        resourceFactory.setTemplateEmbedded(type, templateClass)
    }

    fun search(path: String): Item? {
        val ignoredLinkParts = setOf("", "Oem", "Custom_Company")
        val pathChunked = path.replace("/redfish/v1", "")
                .split("/")
                .filter { it !in ignoredLinkParts }

        var element: Item? = root
        for (part in pathChunked) {
            if (element is Resource && element.data[part] is Actions) {
                element = element.data[part] as Actions
            }
            element?.let {
                element = (it as Searchable).findElement(part, path)
            }
        }
        return when (element) {
            null -> null
            else -> element
        }
    }

    private fun recreateRoot() {
        serviceRoot.data.clear()
        serviceRoot.data.putAll(create(SERVICE_ROOT).data)
    }
}
