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

package com.odim.simulator.tree.structure

sealed class EmbeddedArray<T>: MutableList<T>

class ObjectArray(val templateType: ExtensibleEmbeddedObjectType) : EmbeddedArray<ResourceObject>(), MutableList<ResourceObject> by mutableListOf()

class ValueArray(values: Iterable<Any?> = emptyList()) : EmbeddedArray<Any?>(), MutableList<Any?> by values.toMutableList()
