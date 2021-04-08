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
import com.odim.simulator.dsl.DSL.Companion.makeTypeSafeMap
import com.odim.simulator.tree.structure.EmbeddedArray
import com.odim.simulator.tree.structure.EmbeddedObject
import com.odim.simulator.tree.structure.ExtensibleEmbeddedObjectType
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.Resource.Companion.searchForActions
import com.odim.simulator.tree.structure.Resource.Companion.searchForTreeElementsAsList
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceIdSequence
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.ResourceTypeBase
import com.odim.simulator.tree.structure.SingletonResource
import org.springframework.data.util.AnnotatedTypeScanner
import kotlin.reflect.KClass
import kotlin.reflect.full.createInstance

class ResourceFactory(private val packageNames: List<String>, private val defaultVersion: ResourceVersion) {
    private val templatePackagesNamespace = "com.odim.simulator.tree.templates"
    private val templateMap = mutableMapOf<ResourceTypeBase, KClass<out Any>>()
    private val embeddedObjectsMap = mutableMapOf<ExtensibleEmbeddedObjectType, KClass<out Any>>()
    private val resourceIdSequence = ResourceIdSequence()

    init {
        scanTemplatesPackages()
        scanEmbeddedObjectsPackages()
    }

    fun setTemplate(type: ResourceTypeBase, templateClass: KClass<out ResourceTemplate>) {
        templateMap[type] = templateClass
    }

    fun setTemplateEmbedded(type: ExtensibleEmbeddedObjectType, templateClass: KClass<out ResourceTemplate>) {
        embeddedObjectsMap[type] = templateClass
    }

    fun create(type: ResourceTypeBase, version: ResourceVersion? = null,
               override: ((DSL).() -> Unit)? = null, oemResource: Boolean = false): Resource {
        val resourceVersion = version ?: defaultVersion
        val resource = getResource(type, resourceVersion)
        searchForTreeElementsAsList(resource.data).forEach {
            it.meta.parent = resource
            it.meta.resourceFactory = this
        }
        searchForActions(resource.data).forEach { it.parent = resource }
        resource.meta.resourceFactory = this
        override?.let { resource(override) }
        resource.isOemResource = oemResource
        return resource
    }

    fun createEmbeddedObject(type: ExtensibleEmbeddedObjectType, version: ResourceVersion? = null, override: ((DSL).() -> Unit)? = null): ResourceObject {
        val objectVersion = version ?: defaultVersion
        val embeddedObject = getEmbeddedObject(type, objectVersion)
        fillMapTemplate(embeddedObject, objectVersion)
        override?.run {
            val dsl = DSL()
            dsl.this()
            dsl.applyTo(embeddedObject)
        }
        return embeddedObject
    }

    fun generateId(resource: Resource) = resourceIdSequence.next(resource, resource.meta.parent!! as ResourceCollection)

    private fun scanTemplatesPackages() {
        packageNames.flatMap { packageName ->
            TypeScanner.findTypes(Template::class.java, "$templatePackagesNamespace.$packageName").map { it.kotlin }
        }.forEach { kclazz ->
            val type = (kclazz.annotations.first { it is Template } as Template).type
            templateMap[type] = kclazz
        }
    }

    private fun scanEmbeddedObjectsPackages() {
        packageNames.flatMap { packageName ->
            TypeScanner.findTypes(EmbeddedObjectTemplate::class.java, "$templatePackagesNamespace.$packageName.embedded").map { it.kotlin }
        }.forEach { kclazz ->
            val type = (kclazz.annotations.first { it is EmbeddedObjectTemplate } as EmbeddedObjectTemplate).type
            embeddedObjectsMap[type] = kclazz
        }
    }

    private fun getResource(type: ResourceTypeBase, version: ResourceVersion): Resource {
        val templateClass = templateMap[type] ?: throw TreeBuildingException("No template registered with type: $type. Maybe forgot @Template?")
        val template = templateClass.createInstance() as ResourceTemplate
        val resource = Resource(type, template.getUsedRedfishVersion(version))

        resource.data = template.get(version)
        fillMapTemplate(resource.data, version)

        resource.setOdataCreateStrategy(template.getOdataCreateStrategy(version))

        return resource
    }

    private fun getEmbeddedObject(type: ExtensibleEmbeddedObjectType, version: ResourceVersion): ResourceObject {
        val templateClass = embeddedObjectsMap[type] ?: throw TreeBuildingException("No template registered with type: $type. " +
                "Maybe forgot @EmbeddedObjectTemplate?")
        val template = templateClass.createInstance() as ResourceTemplate
        return template.get(version)
    }

    private fun fillMapTemplate(data: ResourceObject, version: ResourceVersion): ResourceObject {
        for ((key, value) in data) {
            data[key] = fillTemplates(value, version)
        }
        return data
    }

    private fun fillTemplates(value: Any?, version: ResourceVersion): Any? {
        return when (value) {
            is SingletonResource -> create(value.type, version, override = null, oemResource = value.isOemResource)
            is EmbeddedObject -> createEmbeddedObject(value.type, version)
            is Map<*, *> -> fillMapTemplate(makeTypeSafeMap(value).toMutableMap(), version)
            is EmbeddedArray<*> -> value
            is Iterable<*> -> value.map { fillTemplates(it, version) }.toMutableList()
            else -> value
        }
    }

    object TypeScanner {
        private val cache = HashMap<Pair<Class<*>, String>, Set<Class<*>>>()

        fun findTypes(type: Class<out Annotation>, basePackage: String) =
                cache.getOrPut(Pair(type, basePackage)) {
                    AnnotatedTypeScanner(type).findTypes(basePackage)
                }
    }
}
