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

package treebuilding

import com.odim.simulator.RedfishSimulator
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.Request
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import org.testng.Assert.assertFalse
import org.testng.Assert.assertTrue
import org.testng.annotations.BeforeMethod
import org.testng.annotations.Test

class LinkerTest {
    private lateinit var simulator: LinkerSimulator

    @BeforeMethod
    fun prepareTest() {
        simulator = LinkerSimulator()
    }

    @Test
    fun `one way link works`() {
        with (simulator) {
            assertFalse(getResourceJson(system2).toString().contains(chassis2.toLink()))
            assertFalse(getResourceJson(chassis2).toString().contains(system2.toLink()))
            oneWayLink(system2, chassis2)
            assertTrue(getResourceJson(system2).toString().contains(chassis2.toLink()))
            assertFalse(getResourceJson(chassis2).toString().contains(system2.toLink()))
        }
    }

    @Test
    fun `one way link works both ways`() {
        with (simulator) {
            oneWayLink(system2, chassis2)
            assertTrue(getResourceJson(system2).toString().contains(chassis2.toLink()))
            oneWayLink(chassis2, system2)
            assertTrue(getResourceJson(chassis2).toString().contains(system2.toLink()))
        }
    }

    @Test
    fun `find link between resources`() {
        with (simulator) {
            assertTrue(areLinked(chassis, system))
        }
    }

    @Test
    fun `unlink resources works`() {
        with (simulator) {
            removeLink(chassis, system)
            assertFalse(getResourceJson(chassis).toString().contains(system.toLink()))
            assertFalse(getResourceJson(system).toString().contains(chassis.toLink()))
        }
    }

    @Test
    fun `link two resources from either side and verify the link`() {
        with (simulator) {
            assertFalse(areLinkedOnAnySide(system2, chassis2))

            oneWayLink(system2, chassis2)
            assertTrue(oneWayLinkExist(system2, chassis2))
            assertFalse(oneWayLinkExist(chassis2, system2))
            assertTrue(areLinkedOnAnySide(system2, chassis2))

            oneWayLink(chassis2, system2)
            assertTrue(oneWayLinkExist(chassis2, system2))
            assertTrue(areLinked(system2, chassis2))
        }
    }

    private fun RedfishSimulator.getResourceJson(item: Resource) = createResponse(item, Request(GET, url = item.toLink())).json
}

class LinkerSimulator: RedfishSimulator() {
    val system = create(COMPUTER_SYSTEM)
    val system2 = create(COMPUTER_SYSTEM)
    val chassis = create(CHASSIS)
    val chassis2 = create(CHASSIS)

    init {
        root(
                system,
                system2,
                chassis,
                chassis2
        )
        link(chassis, system)
    }
}
