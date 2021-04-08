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

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.dsl.merger.MergeException
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.created
import com.odim.simulator.http.Response.Companion.unauthorized
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.utils.doLaunch
import com.odim.utils.getString
import kotlinx.coroutines.Job
import kotlinx.coroutines.isActive
import java.lang.Thread.sleep
import java.time.LocalTime.now
import java.util.UUID.randomUUID

private const val CREDENTIALS_USERNAME = "UserName"
private const val CREDENTIALS_PASSWORD = "Password"

data class SessionData(val session: Resource, var job: Job)

class PostOnSessionsCollection(private val sessionService: Resource) : Behavior {
    private val sessions = mutableMapOf<String, SessionData>()

    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val reqJson = request.json!!

        val username = reqJson.getString(CREDENTIALS_USERNAME)
        val password = reqJson.getString(CREDENTIALS_PASSWORD)

        val accounts = tree.root.traverse<Resource>("AccountService").traverse<ResourceCollection>("Accounts")
        if (!accounts.members.any {
                    it.traverse<String>("UserName") == username && it.traverse<String>("Password") == password
                }) {
            val returnedResponse = unauthorized()
            return terminal(returnedResponse)
        }

        val sessionsCollection = item as ResourceCollection

        val authToken = randomUUID().toString()
        val sessionId = randomUUID().toString()
        val sessionResource = tree.create(sessionsCollection.type.of()) {
            "Id"        to sessionId
            "UserName"  to username
        }
        sessionsCollection.add(sessionResource)
        sessions[authToken] = SessionData(sessionResource, sessionTimeoutJob(sessionResource))

        return terminal(try {
            Merger.merge(tree, sessionResource, request)

            created(sessionResource, mapOf(
                    "X-Auth-Token" to listOf(authToken)
            ))
        } catch (e: MergeException) {
            badRequest("Bad POST request: ${e.message}")
        })
    }

    fun isLogged(token: String?) = sessions.containsKey(token)

    fun regenerateSession(token: String?) {
        sessions[token]?.let {sessionData ->
            sessionData.job.cancel()
            sessionData.job = sessionTimeoutJob(sessionData.session)
        }
    }

    private fun sessionTimeoutJob(sessionResource: Resource): Job {
        return doLaunch { scope ->
            val timeout = now().plusSeconds(sessionService.traverse<Int>("SessionTimeout").toLong())
            while (true) {
                sleep(100)
                if (!scope.isActive) {
                    return@doLaunch
                } else if (now().isAfter(timeout)) {
                    break
                }
            }
            val sessionsCollection = sessionResource.meta.parent!! as ResourceCollection
            sessionsCollection.remove(sessionResource)

            sessions.filter {
                it.value.session.traverse<String>("Id") == sessionResource.traverse("Id")
            }.map { it.key }.singleOrNull()?.let {
                sessions.remove(it)
            }
        }
    }
}
