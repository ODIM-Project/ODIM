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

class Action(actionType: ActionType) : ActionElement(actionType), Searchable {
    override fun findElement(searchPattern: Any?, other: String) =
            if (this.getActionQualifiedName() == searchPattern) { this } else { null }

    constructor(actionType: ActionType, parameterName: String, allowableValues: MutableList<String>) : this(actionType) {
        this.allowableValues[parameterName] = allowableValues
    }

    constructor(actionType: ActionType, allowableValues: MutableMap<String, MutableList<String>>) : this(actionType) {
        this.allowableValues = allowableValues
    }

    override fun getActionQualifiedName() = "${getActionNamespace()}.${actionType.actionName}"

    override fun toLink() = "${this.parent!!.toLink()}/Actions/${getActionQualifiedName()}"
}
