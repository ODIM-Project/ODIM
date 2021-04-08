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

import com.fasterxml.jackson.annotation.JsonIgnore
import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.behaviors.ExpandType
import com.odim.simulator.behaviors.ExpandType.Companion.of
import com.odim.simulator.behaviors.ExpandType.NONE
import com.odim.utils.toJson
import io.javalin.core.security.BasicAuthCredentials
import javax.servlet.http.HttpServletRequest

data class Request(val method: HttpMethod,
                   val body: String = "",
                   val url: String = "",
                   val headers: Map<String, String> = mutableMapOf(),
                   val queryParams: Map<String, List<String>> = mutableMapOf(),
                   val queryString: String? = null,
                   val basicAuthCredentials: BasicAuthCredentials? = null,
                   @JsonIgnore val request: HttpServletRequest? = null) {
    val json: ObjectNode? by lazy { body.toJson() }
    val expandType: ExpandType by lazy {
        queryString?.let { qs ->
            "\\\$expand=(.)".toRegex().find(qs)?.groupValues?.get(1)?.let {
                of(it)
            } ?: NONE
        } ?: NONE
    }
    val expandLevels: Int by lazy {
        queryString?.let {
            "\\\$expand=.\\(\\\$levels=([0-9]+)\\)".toRegex().find(it)?.groupValues?.get(1)?.toInt() ?: 1
        } ?: 1
    }
}
