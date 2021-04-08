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

package com.odim.utils

import com.fasterxml.jackson.core.JsonParseException
import com.fasterxml.jackson.core.JsonPointer.compile
import com.fasterxml.jackson.databind.JsonNode
import com.fasterxml.jackson.databind.ObjectMapper
import com.fasterxml.jackson.databind.node.ObjectNode
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import com.odim.simulator.http.Request
import com.odim.utils.JsonMapper.jsonMapper
import com.odim.utils.JsonMapper.log
import org.slf4j.LoggerFactory.getLogger
import java.lang.System.getenv

object JsonMapper {
    val jsonMapper = jacksonObjectMapper().registerModule(JavaTimeModule())
    val emptyJson = "{}".toJson()
    val log = getLogger(this.javaClass)
}

fun String.toJson() = jsonMapper.readTree(if (isNotEmpty()) this else "{}") as ObjectNode

@SuppressWarnings("TooGenericExceptionCaught")
fun Request.isValidBodyJson() = try {
    this.json
    ObjectMapper().readTree(this.body)
    true
} catch (e: JsonParseException) {
    log.warn(e.message)
    false
} catch (e: ClassCastException) {
    log.warn(e.message)
    false
} catch (e: Error) {
    log.error(e.message)
    false
}


/**
 * Prints well-formatted JSON. Creates new ObjectMapper locally to apply PrettyPrinter only for this usage.
 */
fun Any.prettyPrint() = jacksonObjectMapper().writer().withDefaultPrettyPrinter().writeValueAsString(this)!!

fun JsonNode.isNull(key: String): Boolean = this.query<Any?>(key) == null

fun JsonNode.getString(jsonPointer: String) = this.query<String>(jsonPointer)
fun JsonNode.getStringOrNull(jsonPointer: String) = this.query<String?>(jsonPointer)

inline fun <reified T> JsonNode.query(jsonPointer: String): T {
    val pointerStart = if (jsonPointer.startsWith("/")) "" else "/"
    return jsonMapper.convertValue(this.at(compile("$pointerStart$jsonPointer")), T::class.java)
}

internal fun getValuePrefEnv(value: String?): String? = when {
    value != null -> getenv(value) ?: value
    else -> value
}

fun JsonNode.merge(updateNode: JsonNode): ObjectNode {
    val fieldNames = updateNode.fieldNames()
    while (fieldNames.hasNext()) {
        val fieldName = fieldNames.next()
        val jsonNode = this.get(fieldName)
        when {
            jsonNode != null && jsonNode.isObject -> jsonNode.merge(updateNode.get(fieldName))
            this is ObjectNode -> this.set(fieldName, updateNode.get(fieldName))
        }
    }
    return this as ObjectNode
}
