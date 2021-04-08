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

package com.odim.simulator.http

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.behaviors.ExpandType
import com.odim.simulator.behaviors.ExpandType.ALL
import com.odim.simulator.behaviors.TreeJsonRenderer
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpStatusCode.ACCEPTED
import com.odim.simulator.http.HttpStatusCode.BAD_REQUEST
import com.odim.simulator.http.HttpStatusCode.CREATED
import com.odim.simulator.http.HttpStatusCode.INTERNAL_SERVER_ERROR
import com.odim.simulator.http.HttpStatusCode.METHOD_NOT_ALLOWED
import com.odim.simulator.http.HttpStatusCode.NOT_FOUND
import com.odim.simulator.http.HttpStatusCode.NOT_MODIFIED
import com.odim.simulator.http.HttpStatusCode.NO_CONTENT
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.http.HttpStatusCode.PRECONDITION_FAILED
import com.odim.simulator.http.HttpStatusCode.UNAUTHORIZED
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.utils.JsonMapper.emptyJson
import com.odim.utils.toJson
import java.util.Collections.singletonList

const val LOCATION_HEADER_NAME = "Location"

data class Response constructor(
    val code: HttpStatusCode,
    val body: String,
    val headers: MutableMap<String, List<String>> = mutableMapOf()
) {
    val json: ObjectNode by lazy {
        if (body.isBlank()) emptyJson else body.toJson()
    }
    val location = headers[LOCATION_HEADER_NAME].orEmpty().firstOrNull()

    companion object {
        fun created(resource: Resource, extraHeaders: Map<String, List<String>> = mapOf()): Response {
            val headers = mutableMapOf(LOCATION_HEADER_NAME to singletonList(resource.toLink()))
            headers.putAll(extraHeaders)
            return Response(CREATED, TreeJsonRenderer().toJson(resource), headers)
        }
        fun success(
            item: Item,
            headers: MutableMap<String, List<String>> = mutableMapOf(),
            expandType: ExpandType = ALL,
            expandLevels: Int = 0
        ) = Response(OK, TreeJsonRenderer(expandType, expandLevels).toJson(item), headers)
        fun success(headers: MutableMap<String, List<String>> = mutableMapOf()) = Response(OK, "", headers)
        fun accepted(location: String, body: String) = Response(ACCEPTED, body, mutableMapOf(LOCATION_HEADER_NAME to listOf(location)))
        fun unauthorized() = predefinedResponse(UNAUTHORIZED)
        fun notAllowed() = predefinedResponse(METHOD_NOT_ALLOWED)
        fun notFound() = predefinedResponse(NOT_FOUND)
        fun noContent() = Response(NO_CONTENT, "")
        fun internalServerError() = Response(INTERNAL_SERVER_ERROR,
                errorJson("The request failed due to an internal service error.  The service is still operational."))
        fun badRequest(message: String = "Bad request") = Response(BAD_REQUEST, errorJson(message, "InvalidPayload"))
        fun notModified() = predefinedResponse(NOT_MODIFIED)
        fun preconditionFailed() = predefinedResponse(PRECONDITION_FAILED)
        fun errorJson(message: String, errorType: String = "GeneralError") = makeJson {
            "error" to makeJson {
                "code" to "Base.1.0.$errorType"
                "message" to message
            }
        }.toString()

        private fun predefinedResponse(httpStatusCode: HttpStatusCode) = Response(httpStatusCode, "")
    }
}
