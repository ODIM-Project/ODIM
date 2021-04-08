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

package com.odim.odimclient.searcher

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.tree.structure.Resource
import com.odim.utils.query

open class FieldPredicate {
    private var excludedFieldNames = setOf("Id")

    open fun test(resource: ObjectNode, simulatorResource: Resource): Boolean =
            toMap(simulatorResource).all { resource.query<Any?>(it.first) == it.second }

    private fun toMap(resource: Resource) = resource.data
            .filter { !isExcludedField(it) }
            .map { Pair(it.key, it.value) }
            .toTypedArray()

    private fun isExcludedField(field: Map.Entry<String, Any?>) = isExcludedFieldName(field.key) || isExcludedFieldValueType(field.value)

    private fun isExcludedFieldName(key: String) = key in excludedFieldNames

    private fun isExcludedFieldValueType(value: Any?): Boolean = !(value is Number || value is String)
}
