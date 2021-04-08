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

package com.odim.simulator.dsl.merger

import com.odim.simulator.RedfishSimulator
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import org.testng.Assert.assertEquals
import org.testng.Assert.assertFalse
import org.testng.Assert.assertNotNull
import org.testng.annotations.Test

class MergerTest {

    @Test
    fun `Merger should remove properties annotated with @removed annotation`() {
        val simulator = RedfishSimulator()
        val system = simulator.create(COMPUTER_SYSTEM)
        simulator.root.append(system)
        assertNotNull(system.data["Boot"])
        Merger.merge(simulator.tree, simulator.tree.search("Systems/1") as Resource, makeJson {
            "UUID" to "1234"
            "Boot@removed" to true
        })
        assertEquals(system.data["UUID"], "1234")
        assertFalse(system.data.containsKey("Boot"))
    }

    @Test(expectedExceptions = [MergeException::class])
    fun `Merger should throw exception when @removed property has value other than true`() {
        val simulator = RedfishSimulator()
        Merger.merge(simulator.tree, simulator.tree.root, makeJson {
            "Systems@removed" to "any value other than true"
        })
    }

    @Test
    fun `Merger tests - check override simple value`() {
        val simulator = RedfishSimulator()
        val system = simulator.create(COMPUTER_SYSTEM)
        Merger.merge(simulator.tree, system, makeJson {
            "UUID" to "1234"
        })
        assertEquals(system.data["UUID"], "1234")
    }

    @Test
    fun `Merger should override MutableMap values`() {
        val simulator = RedfishSimulator()
        Merger.merge(simulator.tree, simulator.root, makeJson {
            "Oem" to {
                "Custom_Company" to {
                    "ApiVersion" to "MutableMapTest"
                }
            }
        })
        simulator.root.traverse<Map<*, *>>("Oem/Custom_Company")
                .filter { (key, _) -> key == "ApiVersion" }
                .forEach { (_, value) ->
                    assertEquals(value, "MutableMapTest")
                }
    }

    @Test
    fun `Merger should read config from failForNonMergeable and skip non mergeable properties`() {
        val simulator = RedfishSimulator()
        Merger {
            failForNonMergeable = false
        }.merge(simulator.tree, simulator.root, makeJson {
            "TelemetryService" to "NonModifiable"
        })
        assertNotNull(simulator.root.data["TelemetryService"])
    }
}