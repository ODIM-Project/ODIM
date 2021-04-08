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

package component.tasks

import com.odim.simulator.RedfishSimulator
import com.odim.simulator.execution.ExecutionMode.ASYNC
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.HttpStatusCode.ACCEPTED
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.utils.getString
import com.odim.utils.retry
import com.odim.utils.url
import component.ComponentSimulatorTest
import org.testng.Assert.assertEquals
import org.testng.annotations.BeforeClass
import org.testng.annotations.Test
import java.time.Duration.ofMinutes


class AsyncTasksTest : ComponentSimulatorTest() {
    private lateinit var simulator: RedfishSimulator
    private lateinit var system: Resource

    @BeforeClass
    fun beforeClazz() {

        simulator = RedfishSimulator()
        system = simulator.create(COMPUTER_SYSTEM)

        with(simulator) {
            root(system)
            provider.serve(this, serverPort)
        }
    }

    @Test
    fun `Set async execution mode for System GET create and complete task`() {
        simulator.setExecutionModeFor(system, GET) {
            mode = ASYNC
        }

        val response = httpClient.get("/redfish/v1/Systems/1")
        val taskMonitorUrl = response.headers["Location"].orEmpty().first()

        assertEquals(response.code, ACCEPTED)

        var task = response.json
        assertEquals(task.getString("TaskState"), "New")

        val responseFromTask = retry(ofMinutes(1)) {
            val response = httpClient.get(taskMonitorUrl)

            if (response.code == OK) success(response)
            else notYet("Cannot get task monitor with url $taskMonitorUrl.")
        }

        task = httpClient.get(task.url()).json
        assertEquals(task.getString("TaskState"), "Completed")

        assertEquals(responseFromTask.code, OK)
        assertEquals(responseFromTask.json.url(), "/redfish/v1/Systems/1")
    }

    @Test
    fun `Return Accepted code after set async execution mode with specified url`() {
        simulator.setExecutionModeFor("/redfish/v1/Systems/.+", GET) {
            mode = ASYNC
        }
        val response = httpClient.get("/redfish/v1/Systems/1")
        assertEquals(response.code, ACCEPTED)
    }

    @Test
    fun `After unset async mode GET request should be synchronous`() {
        val executionBinding = simulator.setExecutionModeFor(system, GET) {
            mode = ASYNC
        }

        // remove binding and check SYNC mode
        simulator.unsetExecutionMode(executionBinding)

        val syncResponse = httpClient.get("/redfish/v1/Systems/1")
        assertEquals(syncResponse.code, OK)
    }
}
