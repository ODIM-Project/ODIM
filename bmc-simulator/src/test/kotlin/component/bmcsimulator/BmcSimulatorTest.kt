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

package component.bmcsimulator

import com.odim.odimclient.RedfishHttpClient
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpStatusCode.NO_CONTENT
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.repo.bmc.BMCSimulator
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import component.ComponentSimulatorTest
import org.testng.Assert.assertEquals
import org.testng.Assert.assertNotNull
import org.testng.annotations.AfterClass
import org.testng.annotations.BeforeClass
import org.testng.annotations.Test

class BmcSimulatorTest : ComponentSimulatorTest() {
    private lateinit var simulator: BMCSimulator
    private lateinit var ip: String
    private var port = 0
    private lateinit var client: RedfishHttpClient
    private var systemUri = ""

    @BeforeClass
    fun config() {
        simulator = BMCSimulator()
        val (ip, port) = provider.serve(simulator, 7896)
        this.ip = ip
        this.port = port
        this.client = RedfishHttpClient(ip) {
            basicCredentialEnabled = true
            basicCredentials = "${simulator.basicAuthUsername}:${simulator.basicAuthPassword}"
        }

        this.systemUri = tryFindSystemUri()
    }

    private fun tryFindSystemUri(): String {
        val systems = simulator.tree.search("/redfish/v1/Systems") as ResourceCollection
        assertNotNull(systems)
        return systems.members.firstOrNull()?.toLink()!!
    }

    @AfterClass
    fun stop() {
        provider.stopAll()
    }

    @Test(priority = 0)
    fun `Simulator should be exposed`() {
        assertEquals(ip.removeSurrounding("http://", ":$port/redfish/v1"), "127.0.0.1")
        assertEquals(port, 7896)
    }

    @Test
    fun `Get on System should return System resource`() {
        checkComponentExists(systemUri)
    }

    @Test
    fun `Get on Memory should return Memory resource`() {
        checkCollectionExists("$systemUri/Memory")
    }

    @Test
    fun `Get on Processor Metrics should return Processor Metrics resource`() {
        checkComponentExists("$systemUri/Processors/cpu_1/Metrics")
    }

    @Test
    fun `Get on Storage should return Storage resource`() {
        checkComponentExists("$systemUri/Storage/BMCStorage")
    }

    @Test
    fun `Get on Chassis collection should return Chassis collection`() {
        checkCollectionExists("/redfish/v1/Chassis")
    }

    @Test
    fun `Get on Chassis should return Chassis resource`() {
        checkComponentExists("/redfish/v1/Chassis/RackMount")
    }

    @Test
    fun `Get on Chassis under chassis should return Chassis resource`() {
        checkComponentExists("/redfish/v1/Chassis/Baseboard")
    }

    @Test
    fun `Get on Thermal on Chassis under chassis should return Thermal resource`() {
        checkComponentExists("/redfish/v1/Chassis/Baseboard/Thermal")
        checkComponentExists("/redfish/v1/Chassis/Baseboard/Power")
    }

    @Test
    fun `System action should be found in tree`() {
        val element = simulator.tree.search("$systemUri/Actions/ComputerSystem.Reset") as Action
        assertNotNull(element)
    }

    @Test
    fun `Post on System should return status 204`() {
        val requestBody = makeJson {
            "ResetType" to "On"
        }

        val response = this.client.post("$systemUri/Actions/ComputerSystem.Reset", requestBody)
        assertEquals(response.code, NO_CONTENT, response.body)
    }

    @Test
    fun `Patch on Chassis should return status 200`() {
        val requestBody = makeJson {
            "Attributes" to {
                "QuietBoot" to 2
                "OnboardNICEnable_0" to 3
                "OnboardNICPortEnable_0" to 3
                "OnboardNICPortEnable_1" to 2
                "OnboardNICPortEnable_2" to 2
                "OnboardNICPortEnable_3" to 2
            }
        }

        val response = this.client.patch("$systemUri/Bios", requestBody)
        assertEquals(response.code, OK, response.body)
    }

    private fun checkComponentExists(componentUri: String) {
        val element = simulator.tree.search(componentUri) as Resource
        assertNotNull(element)
    }

    private fun checkCollectionExists(componentUri: String) {
        val element = simulator.tree.search(componentUri) as ResourceCollection
        assertNotNull(element)
    }
}
