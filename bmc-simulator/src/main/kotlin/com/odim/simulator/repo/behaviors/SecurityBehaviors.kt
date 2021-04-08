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

package com.odim.simulator.repo.behaviors

import com.odim.simulator.CoreConfig.SERVER_BASIC_CREDENTIALS
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.HttpMethod.POST
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.noContent
import com.odim.simulator.http.Response.Companion.unauthorized
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.TreeElement

class SecuredBehavior(private val sessionBehavior: PostOnSessionsCollection) : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val link = generateLink(item)

        if (actionOnRootOrSessions(link, request) || isUserLogged(request, sessionBehavior)) {
            return nonTerminal(response)
        }
        return terminal(unauthorized())
    }

    private fun actionOnRootOrSessions(link: String, request: Request) = (link.endsWith("/redfish/v1")
            || request.method == POST && link.endsWith("/redfish/v1/SessionService/Sessions"))
}

class SystemActionBehavior(private val sessionBehavior: PostOnSessionsCollection) : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val link = generateLink(item)
        if (!link.contains("/redfish/v1/Systems") || isUserLogged(request, sessionBehavior)) {
            return nonTerminal(noContent())
        }
        return terminal(unauthorized())
    }
}

private fun isUserLogged(request: Request, sessionBehavior: PostOnSessionsCollection): Boolean {
    return loggedByBasicAuth(request)
            || sessionBehavior.isLogged(request.headers["X-Auth-Token"])
}

private fun generateLink(item: Item): String {
    return when (item) {
        is ActionElement -> item.toLink()
        else -> (item as TreeElement).toLink()
    }
}

private fun loggedByBasicAuth(request: Request): Boolean {
    request.basicAuthCredentials?.let {
        val user = request.basicAuthCredentials.username
        val password = request.basicAuthCredentials.password
        return getConfigProperty(SERVER_BASIC_CREDENTIALS.path, "") == "$user:$password"
    }
    return false
}
