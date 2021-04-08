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

package com.odim.simulator.repo.bmc.behaviors

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.dsl.merger.MergeException
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.errorJson
import com.odim.simulator.http.Response.Companion.internalServerError
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.utils.getString

@Suppress("ReturnCount")
class PostOnAccounts : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val accountsCollection = item as ResourceCollection
        val violations: MutableList<String> = ArrayList()
        validateRequest(request, violations)
        if (violations.isNotEmpty()) return nonTerminal(badRequest(errorJson(violations.first())))

        validateUserNameDuplication(request, accountsCollection, violations)
        if (violations.isNotEmpty()) return nonTerminal(internalServerError())

        val account = tree.create(accountsCollection.type.of())
        accountsCollection.add(account)

        try {
            Merger.merge(tree, account, request)
        } catch (e: MergeException) {
            return nonTerminal(badRequest("Bad POST request: ${e.message}"))
        }

        return nonTerminal(success())
    }

    private fun validateUserNameDuplication(request: Request, accountsCollection: ResourceCollection, violations: MutableList<String>) {
        val username = request.json!!.getString("UserName")
        accountsCollection.members.forEach {
            if (it.data["UserName"] == username) {
                violations.add("The property UserName already exists.")
            }
        }
    }

    private fun validateRequest(request: Request, violations: MutableList<String>) {
        if (request.json?.getString("UserName") == null) {
            violations.add("The property UserName is a required property and must be included in the request.") }
        if (request.json?.getString("RoleId") == null) {
            violations.add("The property RoleId is a required property and must be included in the request.") }
        if (request.json?.getString("Password") == null) {
            violations.add("The property Password is a required property and must be included in the request.") }
    }
}
