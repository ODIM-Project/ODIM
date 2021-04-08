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

package component.repo

import com.odim.simulator.commandline.SimulatorsFinder.getAllNames
import com.odim.simulator.commandline.SimulatorsFinder.getSimulatorByName
import com.odim.utils.getString
import component.ComponentSimulatorTest
import org.testng.Assert.assertEquals
import org.testng.annotations.Test

class RepoSimulatorsTest : ComponentSimulatorTest() {

    @Test
    fun `All simulators from repo start correctly`() {
        getAllNames().forEach {
            val simulator =  getSimulatorByName(it)
            provider.serve(simulator, serverPort)
            val root = httpClient.serviceRoot()
            assertEquals(root.getString("Name"), "Service Root")
            provider.stop(simulator)
        }
    }
}
