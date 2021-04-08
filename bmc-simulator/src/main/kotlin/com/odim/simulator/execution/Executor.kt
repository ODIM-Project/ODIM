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

package com.odim.simulator.execution

import com.odim.simulator.behaviors.ElementPredicate
import com.odim.simulator.behaviors.ResourceBehaviors
import com.odim.simulator.execution.ExecutionMode.ASYNC
import com.odim.simulator.execution.ExecutionMode.SYNC
import com.odim.simulator.http.HttpMethod
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.accepted
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.TASK
import com.odim.simulator.tree.structure.TreeElement
import com.odim.utils.doLaunch
import java.time.LocalDateTime
import java.util.regex.Pattern

enum class ExecutionMode {
    SYNC,
    ASYNC
}

data class ExecutionBinding(
        val predicate: ElementPredicate,
        var mode: ExecutionMode = SYNC,
        var preDelay: Int = 0,
        var postDelay: Int = 0)

private data class ResponseWrapper(var responseSource: () -> Response)

class Executor {
    private val bindings = mutableListOf<ExecutionBinding>()
    private val tasks = mutableMapOf<String, ResponseWrapper>()

    fun setExecutionModeFor(urlRegex: String, httpMethod: HttpMethod, adjust:  (ExecutionBinding).() -> Unit)
            = createExecutionBinding(createItemUrlPredicate(urlRegex, httpMethod), adjust)

    fun setExecutionModeFor(item: Item, httpMethod: HttpMethod, adjust: (ExecutionBinding).() -> Unit)
            = createExecutionBinding(createItemPredicate(item, httpMethod), adjust)

    fun unsetExecutionMode(binding: ExecutionBinding) {
        bindings.remove(binding)
    }

    fun run(tree: ResourceTree, behaviors: ResourceBehaviors, item: Item, request: Request): Response {
        val binding = bindings.firstOrNull { it.predicate(item, request.method) }
        val mode = binding?.mode ?: SYNC

        val job: () -> Response = { behaviors.createResponse(tree, item, request) }

        return when (mode) {
            SYNC -> doSync(job)
            ASYNC -> doAsync(tree, behaviors, job)
        }
    }

    fun getTask(url: String): Response? {
        if (tasks.containsKey(url)) {
            return tasks[url]?.responseSource?.invoke()
        }
        return null
    }

    private fun doAsync(tree: ResourceTree, behaviors: ResourceBehaviors, job: () -> Response): Response {
        val resource = createTaskResource(tree)
        val monitorUrl = generateTaskMonitorUrl(resource.traverse("Id"))

        val taskSource: () -> Response = {
            accepted(monitorUrl, body = behaviors.createResponse(tree, resource, Request(GET)).body)
        }
        val resourceResponse = taskSource()

        val responseWrapper = ResponseWrapper(taskSource)
        tasks[monitorUrl] = responseWrapper
        runAsyncTask(responseWrapper, resource, job)
        return resourceResponse
    }

    private fun createExecutionBinding(predicate: ElementPredicate, adjust: ExecutionBinding.() -> Unit): ExecutionBinding {
        val executionBinding = ExecutionBinding(predicate)
        executionBinding.adjust()
        bindings.add(executionBinding)
        return executionBinding
    }

    private fun runAsyncTask(responseWrapper: ResponseWrapper, taskResource: Resource, job: () -> Response) {
        doLaunch {
            taskResource {
                "TaskState" to "Running"
                "StartTime" to "${LocalDateTime.now()}"
            }
            val result = job()
            responseWrapper.responseSource = { result }
            taskResource {
                "TaskState" to "Completed"
                "EndTime" to "${LocalDateTime.now()}"
            }
        }
    }

    private fun createTaskResource(tree: ResourceTree): Resource {
        val resource = tree.create(TASK)
        tree.root.traverse<Resource>("Tasks").append(resource)
        return resource
    }

    private fun generateTaskMonitorUrl(id: Int) = "/redfish/v1/TaskService/Tasks/$id/Monitor"

    private fun doSync(job: () -> Response) = job()

    private fun createItemPredicate(item: Item, httpMethod: HttpMethod) =
            { res: Item, onMethod: HttpMethod -> res == item && onMethod == httpMethod }

    private fun createItemUrlPredicate(urlRegex: String, httpMethod: HttpMethod) =
            { res: Item, onMethod: HttpMethod -> httpMethod == onMethod && res is TreeElement && Pattern.matches(urlRegex, res.toLink()) }
}
