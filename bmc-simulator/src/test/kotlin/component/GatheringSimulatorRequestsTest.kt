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

package component

import com.odim.simulator.RedfishSimulator
import com.odim.simulator.tree.structure.ResourceType
import org.testng.Assert.assertEquals
import org.testng.annotations.Test

class GatheringSimulatorRequestsTest: ComponentSimulatorTest() {

    @Test
    fun `All received requests are stored in simulator`() {
        with(RedfishSimulator()) {
            val system = create(ResourceType.COMPUTER_SYSTEM)
            this.root(system)
            provider.serve(this, serverPort)

            httpClient.get(root.toLink())
            httpClient.get(system.toLink())

            assertEquals(requests.size, 2)
            assertEquals(requests.first().url, root.toLink())
            assertEquals(requests.last().url, system.toLink())
        }
    }
}