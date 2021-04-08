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

import com.odim.simulator.behaviors.ExternalBehavior
import com.odim.simulator.behaviors.ResourceBehaviors
import com.odim.simulator.dsl.DSL
import com.odim.simulator.dsl.ResourceMemory
import com.odim.simulator.events.EventArray
import com.odim.simulator.events.EventSender
import com.odim.simulator.execution.ExecutionBinding
import com.odim.simulator.execution.Executor
import com.odim.simulator.http.HttpClient
import com.odim.simulator.http.HttpMethod
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.notFound
import com.odim.simulator.tree.Branch
import com.odim.simulator.tree.RedfishVersion
import com.odim.simulator.tree.RedfishVersion.Companion.REDFISH_VERSION_LATEST
import com.odim.simulator.tree.ResourceFactory
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.ResourceVersion
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.ExtensibleEmbeddedObjectType
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceType
import com.odim.simulator.tree.structure.ResourceTypeBase
import com.odim.simulator.tree.structure.TreeElement

open class RedfishSimulator(defaultResourceVersion: ResourceVersion = REDFISH_VERSION_LATEST,
                            extendingPackages: List<String> = listOf()) : RequestProcessor {
    override val welcomeMessage: String? = null
    private var packages = mutableListOf("redfish")

    init {
        packages.addAll(extendingPackages)
    }

    private val executor = Executor()
    private val eventSender = EventSender(HttpClient())
    private val memory = ResourceMemory()
    private val resourceTree = ResourceTree(ResourceFactory(getTemplatesPackages(), defaultResourceVersion))
    private val resourceBehaviors = ResourceBehaviors()
    private val receivedRequest = mutableListOf<Request>()

    private val externalEndpoints = mutableMapOf<String, ExternalBehavior>()

    val root get() = resourceTree.root
    val tree get() = resourceTree
    val behaviors get() = resourceBehaviors
    val requests get() = receivedRequest.toList()

    val subscriptions get() = getDestinationUrls(getSubscriptions(root.traverse("EventService")))

    fun link(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            tree.link(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun oneWayLink(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null) =
            tree.oneWayLink(resourceFirst, resourceSecond, propertyFirst)

    fun removeLink(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            tree.removeLink(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun areLinked(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            tree.areLinked(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun areLinkedOnAnySide(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null, propertySecond: String? = null) =
            tree.areLinkedOnAnySide(resourceFirst, resourceSecond, propertyFirst, propertySecond)

    fun oneWayLinkExist(resourceFirst: TreeElement, resourceSecond: TreeElement, propertyFirst: String? = null) =
            tree.oneWayLinkExist(resourceFirst, resourceSecond, propertyFirst)

    fun create(type: ResourceTypeBase, version: RedfishVersion? = null, oemResource: Boolean = false, override: (DSL.() -> Unit)? = null) =
            tree.create(type, version, oemResource, override)

    fun createResourcesList(amount: Int, type: ResourceTypeBase, version: RedfishVersion? = null, oemResource: Boolean = false, override: (DSL.() -> Unit)? = null) =
            (0 until amount).map { tree.create(type, version, oemResource, override) }.toCollection(mutableListOf())

    fun createEmbeddedObject(type: ExtensibleEmbeddedObjectType, version: ResourceVersion? = null, override: ((DSL).() -> Unit)? = null) =
            tree.createEmbeddedObject(type, version, override)

    fun sendEvent(events: EventArray) = subscriptions.map { eventSender.sendEvent(events, it) }

    fun setExecutionModeFor(item: Item, httpMethod: HttpMethod, adjust: (ExecutionBinding).() -> Unit) =
            executor.setExecutionModeFor(item, httpMethod, adjust)

    fun setExecutionModeFor(urlRegex: String, httpMethod: HttpMethod, adjust: (ExecutionBinding).() -> Unit) =
            executor.setExecutionModeFor(urlRegex, httpMethod, adjust)

    fun unsetExecutionMode(executionBinding: ExecutionBinding) = executor.unsetExecutionMode(executionBinding)

    fun createResponse(item: Item, request: Request) = executor.run(tree, behaviors, item, request)

    fun addExternalEndpoint(address: String, behavior: ExternalBehavior) {
        externalEndpoints[address] = behavior
    }

    override fun createResponse(request: Request): Response {
        receivedRequest.add(request)
        return executor.getTask(request.url) ?: findResource(request)
    }

    override fun onStart() {
        ResourcesConfigurator(tree).configure()
    }

    override fun onStop() {}

    private fun getTemplatesPackages() = packages.also { list ->
        (this::class.annotations.firstOrNull { it is Branch } as? Branch)?.names?.toList()?.forEach { list.add(it) }
    }

    private fun findResource(request: Request): Response {
        val element = tree.search(request.url)
        if (element != null && !isGetOnAction(element, request)) {
            return createResponse(element, request)
        }
        if (externalEndpoints.containsKey(request.url)) {
            return externalEndpoints[request.url]!!.run(tree, request)
        }
        return notFound()
    }

    /** Evaluate [builder] statements in context of simulator memorizing marked resources.
     * Returns list of resources marked for memorization.
     * Memorize can't be called twice. Don't memorize on different levels of create. */
    fun memorize(builder: RedfishSimulator.() -> Unit): List<Resource> {
        if (memory.open) {
            throw IllegalStateException("Simulator is already in memorizing mode")
        }
        memory.open()
        this.builder()
        memory.close()
        return memory.fetchOrderedResources()
    }

    /** Create multiple resources corresponding to given [types].
     * Returns list of created resources. */
    fun createMany(vararg types: ResourceType, version: RedfishVersion? = null): List<Resource> {
        val res = ArrayList<Resource>(types.size)
        types.forEach { type -> res.add(tree.create(type, version)) }
        return res
    }

    /** Memorize given resource to be returned from invocation of [memorize]. */
    operator fun Resource.unaryPlus(): Resource {
        if (!memory.open) {
            throw IllegalStateException("Can't add resource to memory as simulator is not memorizing")
        }
        memory.record(this)
        return this
    }

    operator fun List<Resource>.component6() = this[5]
    operator fun List<Resource>.component7() = this[6]
    operator fun List<Resource>.component8() = this[7]
    operator fun List<Resource>.component9() = this[8]
    operator fun List<Resource>.component10() = this[9]

    private fun getDestinationUrls(subscriptions: List<Resource>): MutableList<String> {
        val subscriptionsUrls = mutableListOf<String>()
        subscriptions.forEach { subscription ->
            subscriptionsUrls.add(getEventDestinationUrl(subscription))
        }
        return subscriptionsUrls
    }

    private fun getEventDestinationUrl(subscription: Resource) = subscription.data.getValue("Destination") as String

    private fun getSubscriptions(eventService: Resource) = eventService.traverse<ResourceCollection>("Subscriptions").members

    private fun isGetOnAction(element: Item, request: Request) =
            element is ActionElement && request.method == GET
}
