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

import com.fasterxml.jackson.annotation.JsonProperty

data class EventArray(var events: List<Event>)

class ODataId(oDataId: String) {
    @JsonProperty("@odata.id")
    private val oDataId: String = oDataId

    fun oDataId() = this.oDataId
}

data class Event(
        val eventType: EventType,
        val eventId: String = "Event ID",
        val eventTimestamp: String = "Event timestamp",
        val severity: String = "Severe",
        val context: String = "Context",
        val originOfCondition: ODataId = ODataId("0")
)

// Upper camel case to avoid JSON annotation and adding jackson-databind-annotation dependency to project
enum class EventType {
    StatusChange,
    ResourceUpdated,
    ResourceAdded,
    ResourceRemoved,
    Alert
}
