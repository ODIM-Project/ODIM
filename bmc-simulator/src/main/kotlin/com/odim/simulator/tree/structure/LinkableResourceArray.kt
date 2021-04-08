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

import com.odim.simulator.tree.TreeBuildingException

class LinkableResourceArray(type: Type) : LinkableElement(type) {
    private val members = mutableListOf<TreeElement>()

    fun getElements(): List<TreeElement> = this.members

    val first get() = this.members.first()

    override fun addLink(element: TreeElement, linkTo: TreeElement) {
        if (!this.members.contains(element)) {
            this.members.add(element)
        } else {
            throw TreeBuildingException("There is already link between $element and $linkTo!")
        }
    }

    override fun removeLink(element: TreeElement, linkTo: TreeElement) {
        if (this.members.contains(element)) {
            this.members.remove(element)
        } else {
            throw TreeBuildingException("There is no link between $element and $linkTo!")
        }
    }

    override fun isThereLink(element: TreeElement, linkTo: TreeElement) = this.members.contains(element)

}
