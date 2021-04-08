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

package com.odim.simulator.behaviors

import com.fasterxml.jackson.databind.ObjectMapper
import com.odim.simulator.behaviors.ExpandType.ALL
import com.odim.simulator.behaviors.ExpandType.DEPENDENT
import com.odim.simulator.behaviors.ExpandType.NONE
import com.odim.simulator.behaviors.ExpandType.SUBORDINATE
import com.odim.simulator.tree.TreeBuildingException
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.ActionOem
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.GeneratedValue
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.Resource.Companion.searchForExapndableElements
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.TreeElement
import java.util.LinkedList
import java.util.Stack

@Suppress("LargeClass", "LongMethod")
class TreeJsonRenderer(private val expandType: ExpandType = NONE, private val expandLevels: Int = 1) {
    private val edges = mutableSetOf<Pair<Resource, Resource>>()
    private val context = Stack<String>()
    private val resourceContext = Stack<Resource>()
    val expandedResources = mutableListOf<Resource>()

    fun toJson(element: Item): String {
        if (expandType != NONE) {
            buildEdges(element)
        }
        return when (element) {
            is ResourceCollection -> getCollectionJson(element, expandLevels - 1)
            is Resource -> getResourceJson(element, expandLevels)
            else -> throw IllegalArgumentException("Can't extract Json from element of type ${element::class}")
        }
    }

    fun toJson(resourceObject: ResourceObject) = jsonConverter(toJsonString(resourceObject, this.expandLevels))

    @Suppress("ComplexMethod", "NestedBlockDepth")
    private fun buildEdges(element: Item) {
        if (edges.isNotEmpty()) {
            throw TreeBuildingException("TreeJsonRenderer.toJson() can be used only once!")
        }

        val resourcesToVisit = LinkedList<Resource>()
        val visitedResources = mutableListOf<Resource>()

        when (element) {
            is Resource -> resourcesToVisit.add(element)
            is ResourceCollection -> element.members.forEach {
                resourcesToVisit.add(it)
            }
        }

        while (resourcesToVisit.isNotEmpty()) {
            val resource = resourcesToVisit.removeFirst()
            searchForExapndableElements(resource.data).forEach { (_, item) ->
                when (item) {
                    is Resource -> saveEdge(resource, item, visitedResources, resourcesToVisit, element)
                    is ResourceCollection -> item.members.forEach { res ->
                        saveEdge(resource, res, visitedResources, resourcesToVisit, element)
                    }
                    is LinkableResource -> item.getElement()?.let { res ->
                        if (res is Resource) {
                            saveEdge(resource, res, visitedResources, resourcesToVisit, element)
                        }
                    }
                    is LinkableResourceArray -> item.getElements().filter { it is Resource }
                            .map { it as Resource }
                            .forEach { res ->
                                saveEdge(resource, res, visitedResources, resourcesToVisit, element)
                            }
                }
            }
        }
    }

    private fun saveEdge(resFrom: Resource, resTo: Resource,
                         dontAddAgain: MutableList<Resource>,
                         resourcesToVisit: LinkedList<Resource>,
                         startingElement: Item) {
        if (!dontAddAgain.contains(resTo)) {
            resourcesToVisit.add(resTo)
            dontAddAgain.add(resTo)
        }
        if (resTo != startingElement && !edges.any { it.second == resTo }) {
            edges.add(resFrom to resTo)
        }
    }

    private fun getResourceJson(resource: Resource, expandLevel: Int): String {
        var response = """
                    "@odata.context": "${resource.getOdataContext()}",
                    "@odata.id": "${resource.toLink()}",
                    "@odata.type": "${resource.getOdataValue()}",
                """
        resourceContext.push(resource)
        response += toJsonString(resource.data, expandLevel)
        resourceContext.pop()
        return jsonConverter(response)
    }

    private fun getEmbeddedResourceJson(resource: Resource, index: Int): String {
        val response =
                """
                    "@odata.id": "${resource.toLink()}/$index",
                    "MemberId": "$index",
                """.run { plus(toJsonString(resource.data, 0)) }

        return jsonConverter(response)
    }

    private fun jsonConverter(response: String): String {
        val potentiallyWrongJson = "{ $response }"

        val objectMapper = ObjectMapper()
        val betterJson = objectMapper.readTree(potentiallyWrongJson)
        return objectMapper.writeValueAsString(betterJson)
    }

    private fun getCollectionJson(collection: ResourceCollection, expandLevel: Int): String {
        val jsonName = collection.linkPart
        val arr = if (listOf(ALL, SUBORDINATE).contains(expandType)) getExpandedMemberList(collection, expandLevel) else getMemberList(collection)

        return """{
                    "@odata.context": "${collection.getOdataContext()}",
                    "@odata.id": "${collection.toLink()}",
                    "@odata.type": "${collection.getOdataValue()}",
                    "Name": "$jsonName Collection",
                    "Description": "$jsonName Collection",
                    "Members@odata.count": ${arr.size},
                    "Members": [${arr.joinToString(",")}],
                    "Oem": { }
        }"""
    }

    private fun getMemberList(collection: ResourceCollection): MutableList<String> {
        val arr = mutableListOf<String>()
        for (resource in collection.members) {
            arr.add("""{"@odata.id":"${resource.toLink()}"}""")
        }

        return arr
    }

    private fun getExpandedMemberList(collection: ResourceCollection, expandLevel: Int): MutableList<String> {
        val arr = mutableListOf<String>()
        for (resource in collection.members) {
            expandedResources.add(resource)
            arr.add(getResourceJson(resource, expandLevel))
        }

        return arr
    }

    private fun getActionsJson(actions: List<ActionElement>): String {
        val arr = mutableListOf<String>()
        val arrOem = mutableListOf<String>()
        actions.filter { it is Action }
                .forEach { arr.add(getActionJson(it)) }
        actions.filter { it is ActionOem }
                .forEach { arrOem.add(getActionJson(it)) }
        return """${arr.joinToString(",")}${if (arr.isEmpty()) "" else ","}"Oem":{${arrOem.joinToString(",")}}"""
    }

    // danger: there is assumption that parent of action will always be a resource
    private fun getActionJson(action: ActionElement): String {
        val arr = mutableListOf<String>().apply {
            add(""""target": "${action.toLink()}"""")
            for ((paramName, allowableValues) in action.allowableValues) {
                add(""""$paramName@Redfish.AllowableValues":${allowableValues.map { s -> """"$s"""" }}""")
            }
        }
        return """
            "#${action.getActionQualifiedName()}": {
                ${arr.joinToString(",")}
            }
        """
    }

    private fun shouldExpand(element: TreeElement?, expandLevel: Int) = expandLevel > 0
            && (expandType == ALL || expandType == SUBORDINATE && !context.contains("Links") || expandType == DEPENDENT && context.contains("Links"))
            && (element is ResourceCollection || edges.any { it.first == resourceContext.peek() && it.second == element })

    @Suppress("ComplexMethod", "NestedBlockDepth")
    private fun toJsonString(data: Any, expandLevel: Int): String = when (data) {
        is Map<*, *> -> {
            val arr = mutableListOf<String>()
            for ((prop, value) in data) {
                context.push(prop as String)
                when (value) {
                    is Actions -> arr.add(""""$prop":{${getActionsJson(value.getActions())}}""")
                    is EmbeddedResourceArray -> {
                        arr.add(""""$prop":[${toJsonString(value.data, expandLevel)}]""")
                    }
                    is LinkableResource -> {
                        val element = value.getElement()
                        if (shouldExpand(element, expandLevel)) {
                            arr.add(""""$prop":${if (element != null) {
                                element as Resource
                                expandedResources.add(element)
                                getResourceJson(element, expandLevel - 1)
                            } else "null"}""")
                        } else {
                            val link = value.getElement()?.toLink()
                            arr.add(""""$prop":${if (link != null) """{"@odata.id": "$link"}""" else "null"}""")
                        }
                    }
                    is LinkableResourceArray -> {
                        arr.add(""""$prop":[${toJsonString(value.getElements(), expandLevel)}]""")
                    }
                    is TreeElement -> {
                        if (shouldExpand(value, expandLevel)) {
                            if (value is Resource) {
                                expandedResources.add(value)
                                arr.add(""""$prop":${getResourceJson(value, expandLevel - 1)}""")
                            } else if (value is ResourceCollection) {
                                arr.add(""""$prop":${getCollectionJson(value, expandLevel - 1)}""")
                            }
                        } else {
                            val link = """{"@odata.id": "${value.toLink()}"}"""
                            arr.add(""""$prop":$link""")
                        }
                    }
                    is GeneratedValue<*> -> {
                        val generatedValue = value.getValue()
                        arr.add(""""$prop":${if (generatedValue is String) """"$generatedValue"""" else generatedValue}""")
                    }
                    is Map<*, *> -> arr.add(""""$prop":{${toJsonString(value, expandLevel)}}""")
                    is List<*> -> arr.add(""""$prop":[${toJsonString(value, expandLevel)}]""")
                    else -> {
                        val part = if (value is String || value is Enum<*>  || prop == "Id") """"$value"""" else value
                        arr.add(""""$prop":$part""")
                    }
                }
                context.pop()
            }
            arr.joinToString(",")
        }
        is List<*> -> {
            val arr = mutableListOf<String>()
            for ((index, value) in data.withIndex()) {
                when (value) {
                    is TreeElement -> when {
                        value.meta.type.isEmbedded() -> arr.add(getEmbeddedResourceJson(value as Resource, index))
                        else -> {
                            if (shouldExpand(value, expandLevel)) {
                                if (value is Resource) {
                                    expandedResources.add(value)
                                    arr.add(getResourceJson(value, expandLevel - 1))
                                } else if (value is ResourceCollection) {
                                    arr.add(getCollectionJson(value, expandLevel - 1))
                                }
                            } else {
                                arr.add("""{"@odata.id": "${value.toLink()}"}""")
                            }
                        }
                    }
                    is Map<*, *> -> arr.add("""{${toJsonString(value, expandLevel)}}""")
                    is List<*> -> arr.add("""[${toJsonString(value, expandLevel)}]""")
                    else -> arr.add(when (value) {
                        is String -> """"$value""""
                        else -> "$value"
                    })
                }
            }
            arr.joinToString(",")
        }
        else -> ""
    }
}
