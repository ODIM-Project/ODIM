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

import com.fasterxml.jackson.databind.node.ArrayNode
import com.fasterxml.jackson.databind.node.ObjectNode

fun ObjectNode.url(): String = this.query("@odata.id")

fun ObjectNode.getBoolean(jsonPointer: String) = this.query<Boolean>(jsonPointer)
fun ObjectNode.getBooleanOrNull(jsonPointer: String) = this.query<Boolean?>(jsonPointer)

fun ObjectNode.getStringOrNullPrefEnv(jsonPointer: String) = getValuePrefEnv(this.query(jsonPointer))

fun ObjectNode.getString(jsonPointer: String) = this.query<String>(jsonPointer)
fun ObjectNode.getStringOrNull(jsonPointer: String) = this.query<String?>(jsonPointer)

fun ObjectNode.getNumber(jsonPointer: String) = this.query<Number>(jsonPointer)
fun ObjectNode.getNumberOrNull(jsonPointer: String) = this.query<Number?>(jsonPointer)

fun ObjectNode.getObject(jsonPointer: String) = this.query<ObjectNode>(jsonPointer)
fun ObjectNode.getObjectOrNull(jsonPointer: String) = this.query<ObjectNode?>(jsonPointer)

fun ObjectNode.getArray(jsonPointer: String) = this.query<ArrayNode>(jsonPointer)
fun ObjectNode.getArrayOrNull(jsonPointer: String) = this.query<ArrayNode?>(jsonPointer)

@Suppress("UNCHECKED_CAST")
fun ObjectNode.toMap(): Map<String, Any?> = JsonMapper.jsonMapper.convertValue(this, Map::class.java) as Map<String, Any?>
