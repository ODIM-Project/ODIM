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

package com.odim.simulator.dsl.merger

interface Getter<T> {
    val context: Getter<Any?>?
    val root: Getter<*>?
        get() {
            var that: Getter<*> = this
            while (that.context != null) {
                that = that.context ?: break
            }
            return that
        }

    fun get(): T?
    fun exists(): Boolean
}

interface Property<T> : Getter<T> {
    override val context: Getter<Any?>
    fun set(value: T)
}

class Value<T>(private val value: T, override val context: Getter<Any?>? = null) : Getter<T> {
    override fun get(): T? = value

    override fun exists() = true

    override fun toString() = "Value(value=$value, context=$context)"
}

@Suppress("UNCHECKED_CAST")
class MapProperty<T>(private val property: String, override val context: Getter<Any?>) : Property<T> {
    override fun get(): T? = (context.get() as? Map<String, T>)?.getOrDefault(property, null)

    override fun exists() = (context.get() as? Map<String, *>)?.containsKey(property) ?: false

    override fun set(value: T) {
        (context.get() as? MutableMap<String, T> ?: mutableMapOf())[property] = value
    }

    override fun toString() = "MapProperty(property='$property', context=$context)"
}

@Suppress("UNCHECKED_CAST")
class ListProperty<T>(private val index: Int, override val context: Getter<Any?>) : Property<T> {
    override fun get(): T? = (context.get() as? List<T>)?.getOrNull(index)

    override fun exists() = index < (context.get() as? List<*>)?.size ?: 0

    override fun set(value: T) {
        (context.get() as? MutableList<T> ?: mutableListOf())[index] = value
    }

    override fun toString() = "ListProperty(index=$index, context=$context)"
}
