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

import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.forEach
import java.util.TreeSet

private data class Delta(val parent: ResourceVersion? = null,
                         val otherDeltas: List<RedfishVersion> = listOf(),
                         val additive: ResourceObject = mutableMapOf(),
                         val changes: (ResourceObject) -> Unit = {},
                         var odataTypeStrategy: ((ResourceVersion) -> String?)? = null)

@Suppress("UnnecessaryAbstractClass")
abstract class ResourceTemplate {
    private val deltas = mutableMapOf<ResourceVersion, Delta>()

    fun version(version: ResourceVersion, parent: ResourceVersion, otherDeltas: List<RedfishVersion>, additive: ResourceObject = mutableMapOf(),
                changes: (ResourceObject) -> Unit = {}) {
        deltas[version] = Delta(parent, otherDeltas, additive, changes)
    }

    fun version(version: ResourceVersion, parent: ResourceVersion, additive: ResourceObject = mutableMapOf(),
                changes: (ResourceObject) -> Unit = {}) {
        deltas[version] = Delta(parent, additive = additive, changes = changes)
    }

    fun version(version: RedfishVersion, additive: ResourceObject = mutableMapOf(),
                changes: (ResourceObject) -> Unit = {}) {
        deltas[version] = Delta(additive = additive, changes = changes)
    }

    fun oDataTypeStrategy(version: ResourceVersion, odataTypeStrategy: (ResourceVersion) -> String?) {
        deltas[version]?.let {
            it.odataTypeStrategy = odataTypeStrategy
        } ?: throw TreeBuildingException("Create delta for version $version before setting oDataTypeStrategy.")
    }

    fun updateVersion(version: ResourceVersion, additive: ResourceObject) {
        if (!deltas.containsKey(version)) {
            throw TreeBuildingException("You cannot update not existing version ($version)")
        }
        mergeMaps(deltas[version]!!.additive, additive)
    }

    fun get(version: ResourceVersion): ResourceObject {
        val result: ResourceObject = mutableMapOf()
        val deltaList = getDeltasForVersion(version)

        deltaList.forEach { del ->
            del.otherDeltas.forEach {
                if (deltas.containsKey(it)) {
                    val d = deltas[it]!!
                    mergeMaps(result, d.additive)
                    d.changes(result)
                }
            }
            mergeMaps(result, del.additive)
            del.changes(result)
        }

        sortMapKeys(result)
        return result
    }

    fun getOdataCreateStrategy(version: ResourceVersion): ((RedfishVersion) -> String?)? {
        var strategy: ((RedfishVersion) -> String?)? = null
        val deltaList = getDeltasForVersion(version)

        deltaList.forEach { del ->
            del.odataTypeStrategy?.let { strategy = it }
        }

        return strategy
    }

    fun getUsedRedfishVersion(version: ResourceVersion): RedfishVersion {
        var ver: RedfishVersion? = null

        val resourceVersion = getResourceVersion(version)

        if (resourceVersion is RedfishVersion) {
            return resourceVersion
        }

        var delta: Delta? = deltas[resourceVersion]

        while (ver == null && delta != null) {
            if (delta.parent is RedfishVersion) {
                ver = delta.parent as RedfishVersion
            }
            if (delta.otherDeltas.isNotEmpty()) {
                ver = delta.otherDeltas.max()
            }
            delta = deltas[delta.parent]
        }

        if (ver == null) {
            throw TreeBuildingException("Versions inheritance tree of template class ${this::class.simpleName} does not include any Redfish version")
        }
        return ver
    }

    private fun getDeltasForVersion(version: ResourceVersion): List<Delta> {
        val deltaList = mutableListOf<Delta>()

        var delta: Delta? = deltas[getResourceVersion(version)]

        while (delta != null) {
            deltaList.add(delta)
            delta = deltas[delta.parent]
        }
        return deltaList.reversed()
    }

    @Suppress("UNCHECKED_CAST")
    private fun mergeMaps(source: ResourceObject, changes: ResourceObject) {
        changes.forEach {
            when {
                it.value is Actions -> {
                    if (source["Actions"] == null) {
                        source["Actions"] = Actions()
                    }
                    (it.value as Actions).forEach { action ->
                        val actions = source["Actions"] as Actions
                        actions.getActionOrNull(action.actionType)?.let { actionElement ->
                            actions.removeAction(actionElement.actionType)
                        }
                        actions.addAction(action, null)
                    }
                }
                it.value is Map<*, *> && source.containsKey(it.key) -> mergeMaps(source[it.key] as ResourceObject, it.value as ResourceObject)
                else -> source[it.key] = it.value
            }
        }
    }

    private fun sortMapKeys(result: ResourceObject) {
        if (result.containsKey("Links")) {
            result["Links"] = result.remove("Links")
        }
        if (result.containsKey("Actions")) {
            result["Actions"] = result.remove("Actions")
        }
        if (result.containsKey("Oem")) {
            result["Oem"] = result.remove("Oem")
        }
    }

    private fun getResourceVersion(version: ResourceVersion) = TreeSet<ResourceVersion>(deltas.filter {
        it.key.branch == version.branch
    }.map { it.key }).filter {
        it.value <= version.value
    }.maxBy {
        it.value
    } ?: if (deltas.any { it.key.branch == version.branch }) {
        throw TreeBuildingException("No compatible version to ${version::class.simpleName}($version) could not be found. " +
                "Check if template class ${this::class.simpleName} defines version $version or lower.")
    } else {
        deltas.map { it.key }.last()
    }
}
