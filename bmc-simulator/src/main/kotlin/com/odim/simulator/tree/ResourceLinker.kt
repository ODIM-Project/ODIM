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

import com.odim.simulator.tree.structure.LinkableElement
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.TreeElement

object ResourceLinker {
    /**
     * Method to link two simulator resources together. It tries to append link to both sides and success only if particular
     * resource is able to hold the link. Otherwise just ignore.
     */
    fun link(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) {
        singleSimulatorCheck(resourceFirst, resourceSecond)

        if (resourceFirst is Resource) {
            searchForLinkableElement(resourceFirst.data, resourceSecond, propertyFirst)?.addLink(resourceSecond, resourceFirst)
        }

        if (resourceSecond is Resource) {
            searchForLinkableElement(resourceSecond.data, resourceFirst, propertySecond)?.addLink(resourceFirst, resourceSecond)
        }
    }

    /**
     * Method to link two simulator resources together. It tries to append link to one side and success only if particular
     * resource is able to hold the link. Otherwise just ignore. Useful when linking resources that are already linked on the other side.
     */
    fun oneWayLink(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null) {
        singleSimulatorCheck(resourceFirst, resourceSecond)

        if (resourceFirst is Resource) searchForLinkableElement(resourceFirst.data, resourceSecond, propertyFirst)?.addLink(resourceSecond, resourceFirst)
    }

    /**
     * Method to remove link between two resources. It tries to replace link with null to both sides.
     */
    fun removeLink(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) {
        singleSimulatorCheck(resourceFirst, resourceSecond)

        if (resourceFirst is Resource) {
            searchForLinkableElement(resourceFirst.data, resourceSecond, propertyFirst)?.removeLink(resourceSecond, resourceFirst)
        }

        if (resourceSecond is Resource) {
            searchForLinkableElement(resourceSecond.data, resourceFirst, propertySecond)?.removeLink(resourceFirst, resourceSecond)
        }
    }

    /**
     * Check if resources are linked on both sides.
     * @return true if link exist on both sides, false if doesn't exist or it is one way link
     */
    fun areLinked(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null): Boolean =
            checkIfOneSideLinkExist(resourceFirst, resourceSecond, propertyFirst) && checkIfOneSideLinkExist(resourceSecond, resourceFirst, propertySecond)

    /**
     * Check if resources are linked on any side
     * @return true if link exist on at least one side, false if doesn't exist at all
     */
    fun areLinkedOnAnySide(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            checkIfOneSideLinkExist(resourceFirst, resourceSecond, propertyFirst) || checkIfOneSideLinkExist(resourceSecond, resourceFirst, propertySecond)

    /**
     * Check if there is link to resource
     * @return true if link exist
     */

    fun checkIfOneSideLinkExist(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String?): Boolean {
        singleSimulatorCheck(resourceFirst, resourceSecond)
        return searchForLinkableElement((resourceFirst as Resource).data, resourceSecond, propertyFirst)?.isThereLink(resourceSecond, resourceFirst) ?: false
    }

    private fun singleSimulatorCheck(resourceFirst: TreeElement, resourceSecond: TreeElement) {
        if (resourceFirst.meta.resourceFactory != resourceSecond.meta.resourceFactory) {
            throw LinkingException("You attempt to link resources from different simulators!")
        }
    }

    private fun searchForLinkableElement(data: Map<*, *>, accepts: TreeElement, property: String? = null): LinkableElement? {
        val linkableResources = searchForLinkableElements(data, accepts, property)
        val linkableSpecified = linkableResources.filter { !it.value.acceptsAny() }
        val linkableAny = linkableResources.filter { it.value.acceptsAny() }
        val message = "More than one LinkableResource %s accepts %s! Specify particular property."
        when {
            linkableSpecified.size > 1 -> throw TreeBuildingException(message.format(linkableSpecified.keys, accepts.meta.type))
            containsFewAndOnlyLinkableAny(linkableSpecified, linkableAny) ->
                throw TreeBuildingException(message.format(linkableAny.keys, accepts.meta.type))
            else -> return linkableSpecified.values.firstOrNull() ?: linkableAny.values.firstOrNull()
        }
    }

    private fun containsFewAndOnlyLinkableAny(linkableSpecified: Map<String, LinkableElement>, linkableAny: Map<String, LinkableElement>) =
            linkableSpecified.isEmpty() && linkableAny.size > 1

    private fun searchForLinkableElements(data: Map<*, *>, accepts: TreeElement, property: String? = null): Map<String, LinkableElement> {
        val linkableResources = mutableMapOf<String, LinkableElement>()
        for ((prop, value) in data) {
            if (value is LinkableElement && value.accepts(accepts)) {
                if (property != null && prop != property) {
                    continue
                }
                linkableResources[prop as String] = value
            } else if (value is Map<*, *>) {
                linkableResources.putAll(searchForLinkableElements(value, accepts, property))
            } else if (value is List<*>) {
                value.forEach {
                    if (it is Map<*, *>) {
                        linkableResources.putAll(searchForLinkableElements(it, accepts, property))
                    }
                }
            }
        }
        return linkableResources
    }

    private class LinkingException(message: String): IllegalArgumentException(message)
}
