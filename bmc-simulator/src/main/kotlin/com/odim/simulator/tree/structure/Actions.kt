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

class Actions(vararg actions: ActionElement) : Item, Searchable {
    override fun findElement(searchPattern: Any?, other: String) = this.getActions().find {
        it.getActionQualifiedName() == other.split("/").last()
    }

    private val actions = mutableListOf<ActionElement>()

    init {
        this.actions.addAll(actions)
    }

    fun getActions() = this.actions

    fun getAction(type: ActionType): ActionElement = actions.single { actionElement -> actionElement.actionType == type }

    fun getActionOrNull(type: ActionType): ActionElement? = actions.singleOrNull { actionElement -> actionElement.actionType == type }

    fun addAction(action: ActionElement, addTo: Resource?) {
        action.parent = addTo
        this.actions.add(action)
    }

    fun removeAction(type: ActionType) {
        this.actions.removeIf { it.actionType == type }
    }

    operator fun iterator() = this.actions.iterator()
}

inline fun Actions.forEach(action: (ActionElement) -> Unit) {
    for (element in this) action(element)
}
