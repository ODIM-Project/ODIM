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

import com.odim.simulator.events.JsonSerializer.deserialize
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.http.ServerFactory.create
import io.javalin.Javalin
import org.slf4j.LoggerFactory

class EventsCatcher(topicName: String, private val ip: String = "127.0.0.1") {
    private val logger = LoggerFactory.getLogger(this.javaClass)

    private val eventsPath = "/events/$topicName"
    private val eventArrays = mutableListOf<EventArray>()
    private var server: Javalin? = null

    val events get() = eventArrays.toList()
    val destinationUrl get() = "http://$ip:${server?.port()}$eventsPath"

    fun startListening(port: Int = 0): EventsCatcher {
        server = create(ip, port, false)
                .post(eventsPath) {
                    it.status(OK.value).result("Received")
                    logger.info("Events received: ${it.body()}")
                    eventArrays.add(deserialize(it.body()))
                }.start()
        return this
    }

    fun stop() {
        server?.stop()
    }
}
