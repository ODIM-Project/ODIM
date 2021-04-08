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

package com.odim.simulator

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.CoreConfig.RESOURCES_CONFIG
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionTypeBase
import org.slf4j.LoggerFactory.getLogger
import kotlin.system.measureTimeMillis

@Suppress("LongMethod")
class ResourcesConfigurator(private val tree: ResourceTree) {
    private val logger = getLogger(this.javaClass)

    fun configure() {
        val resourcesToOverride: MutableMap<Resource, String> = mutableMapOf()
        val execTime = measureTimeMillis {
            getConfigProperty<ObjectNode>(RESOURCES_CONFIG).fields().also {
                logger.debug("Configuring resources")
            }.forEach { (resourceLink, override) ->
                val item = tree.search(resourceLink)
                when (item) {
                    is ResourceCollection -> {
                        if (!override.canConvertToInt()) {
                            throw ResourcesConfigurationException("Configuration of ${item.toLink()} number value was expected.")
                        }
                        val missing = override.asInt() - item.size()
                        if (missing < 0) {
                            throw ResourcesConfigurationException("Expected number of elements in collection ${item.toLink()} is less than already exists.")
                        }
                        val resourceType = (item.meta.type as ResourceCollectionTypeBase).of()
                        repeat(missing) {
                            item.add(tree.create(resourceType))
                        }
                    }
                    is Resource -> {
                        (override as ObjectNode).remove("Id")?.run { resourcesToOverride[item] = this.asText() }
                        Merger.merge(tree, item, override)
                    }
                    is ActionElement -> override.fields().forEach { (param, values) ->
                        val paramName = param.replace("@Redfish.AllowableValues", "")
                        item.allowableValues[paramName] = values.map { it.textValue() }.toMutableList()
                    }
                    else -> throw ResourcesConfigurationException("Resource $resourceLink not found.")
                }
            }
            resourcesToOverride.forEach { (resource, id) ->
                resource {
                    "Id" to id
                }
            }
        }
        logger.debug("Configuring resources time: $execTime ms")
    }
}

class ResourcesConfigurationException(message: String) : RuntimeException(message)
