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

package com.odim.simulator.events

import com.odim.odimclient.RedfishHttpClient
import com.odim.simulator.events.JsonSerializer.serialize
import com.odim.simulator.events.Subscription.Companion.forEventType
import com.odim.simulator.events.Subscription.Companion.forOriginEvents
import com.odim.simulator.http.Response

data class Subscription(
        val name: String,
        val destination: String,
        val eventTypes: List<EventType>,
        val context: String,
        val protocol: String = "Redfish",
        val originResources: List<ODataId> = emptyList()
) {
    companion object {
        fun forEventType(destinationUrl: String, vararg events: EventType) = Subscription(
                name = "Subscription",
                destination = destinationUrl,
                eventTypes = events.toList(),
                context = "Context")


        fun forOriginEvents(destinationUrl: String, originOfContent: List<ODataId>, vararg events: EventType) = Subscription(
                name = "Subscription",
                destination = destinationUrl,
                eventTypes = events.toList(),
                context = "Context",
                originResources = originOfContent)
    }
}

class ClientEventService(val redfishHttpClient: RedfishHttpClient) {

    private val subscriptionUrls: MutableList<String> = mutableListOf()

    fun subscribe(destinationUrl: String, vararg eventTypes: EventType): Response {
        val subscription = forEventType(destinationUrl, *eventTypes)
        return sendSubscription(subscription)
    }

    fun subscribeWithEventsOfOrigin(destinationUrl: String, eventTypes: Array<EventType>, originOfContent: List<ODataId>): Response {
        val subscription = forOriginEvents(destinationUrl, originOfContent, *eventTypes)
        return sendSubscription(subscription)
    }

    private fun sendSubscription(subscription: Subscription): Response {
        val resp = redfishHttpClient.post("/EventService/Subscriptions", serialize(subscription))
        if (resp.headers.containsKey("Location")) {
            subscriptionUrls.addAll(resp.headers["Location"].orEmpty())
        } else {
            throw IllegalArgumentException("There should be a location of subscription")
        }
        return resp
    }

    fun unsubscribeAll() = subscriptionUrls.forEach { url: String -> redfishHttpClient.delete(url) }
}
