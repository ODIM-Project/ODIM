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
import com.odim.simulator.tree.structure.EmbeddedArray
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollectionType.PROCESSORS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.FABRIC
import com.odim.simulator.tree.structure.ResourceType.PORT
import com.odim.simulator.tree.structure.ResourceType.POWER_CONTROL
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.SWITCH
import com.odim.simulator.tree.structure.StringArray
import org.testng.Assert.assertEquals
import org.testng.Assert.assertNotNull
import org.testng.Assert.assertTrue
import org.testng.annotations.Test

class TreeBuildingTests {

    private val simulator = RedfishSimulator()

    @Test
    fun `Create simple resource`() {
        val system = simulator.create(COMPUTER_SYSTEM)
        assertEquals(system.meta.type, COMPUTER_SYSTEM)
    }

    @Test
    fun `Append other resource to resource`() {
        val system = simulator.create(COMPUTER_SYSTEM)
        simulator.root.append(system)
        assertNotNull(system.meta.parent)

        val processor = simulator.create(PROCESSOR)
        system.append(processor)
        assertNotNull(processor.meta.parent)

        assertEquals(processor.meta.parent!!.meta.type, PROCESSORS_COLLECTION)
    }

    @Test
    fun `Append resource to other resource`() {
        val system = simulator.create(COMPUTER_SYSTEM)
        val processor = simulator.create(PROCESSOR)

        processor.appendTo(system).appendTo(simulator.root)

        assertNotNull(system.meta.parent)
        assertNotNull(processor.meta.parent)
        assertEquals(processor.meta.parent!!.meta.type, PROCESSORS_COLLECTION)
    }

    @Test
    fun `Building links`() {
        val chassis = simulator.create(CHASSIS)
        val powerControl = simulator.create(POWER_CONTROL)
        val power = chassis.traverse<Resource>("Power")

        power.traverse<EmbeddedResourceArray>("PowerControl").add(powerControl)
        chassis.appendTo(simulator.root)

        assertEquals(simulator.root.toLink(), "/redfish/v1")
        assertEquals(chassis.toLink(), "/redfish/v1/Chassis/1")
        assertEquals(power.toLink(), "/redfish/v1/Chassis/1/Power")
        assertEquals(powerControl.toLink(), "/redfish/v1/Chassis/1/Power#/PowerControl")
    }

    @Test(enabled = false)
    fun `Traversing resources`() {
        val system = simulator.create(COMPUTER_SYSTEM)

        assertEquals(system.traverse<String>("Description").toLowerCase(), "Computer System description".toLowerCase())
        assertEquals(system.traverse("Status/Health"), "OK")
        assertEquals(system.traverse("Oem/Custom_Company/@odata.type"), "#CustomCompany.Oem.ComputerSystem")
        assertTrue(system.traverse<StringArray>("Oem/Custom_Company/PCIeConnectionId").isEmpty())
    }

    @Test
    fun `Linking resources`() {
        with(simulator) {
            val chassis = create(CHASSIS)
            val system1 = create(COMPUTER_SYSTEM)
            val system2 = create(COMPUTER_SYSTEM)

            simulator.root.append(chassis).append(system1).append(system2)

            simulator.link(chassis, system1)
            simulator.link(chassis, system2)

            assertTrue(chassis.traverse<LinkableResourceArray>("Links/ComputerSystems").getElements().containsAll(listOf(system1, system2)))
            assertTrue(system1.traverse<LinkableResourceArray>("Links/Chassis").getElements().contains(chassis))
            assertTrue(system2.traverse<LinkableResourceArray>("Links/Chassis").getElements().contains(chassis))
        }
    }

    @Test(enabled = false)
    fun `Resource properties retain types declared in templates`() {
        with(RedfishSimulator()) {
            val processor = create(PROCESSOR)
            create(COMPUTER_SYSTEM)(processor)
            assertTrue(processor.traverse<Any>("Oem/Custom_Company/Capabilities") is EmbeddedArray<*>)
            assertTrue(processor.traverse<Any>("Oem/Custom_Company/OnPackageMemory") is EmbeddedArray<*>)
        }
    }

    @Test
    fun `Resource indexing`() {
        with(RedfishSimulator()) {
            val fabric = create(FABRIC)
            val switch1 = create(SWITCH)
            val switch2 = create(SWITCH)
            val port11 = create(PORT)
            val port12 = create(PORT)
            val port21 = create(PORT)

            assertNotIndexed(listOf(fabric, switch1, switch2, port11, port12, port21))

            fabric(
                    switch1(port11, port12),
                    switch2(port21)
            )

            assertNotIndexed(listOf(fabric, switch1, switch2, port11, port12, port21))

            root(fabric)

            assertEquals(fabric.traverse("Id"), 1)
            assertEquals(switch1.traverse("Id"), 1)
            assertEquals(switch2.traverse("Id"), 2)
            assertEquals(port11.traverse("Id"), 1)
            assertEquals(port12.traverse("Id"), 2)
            assertEquals(port21.traverse("Id"), 1)
        }
    }

    private fun assertNotIndexed(resources: List<Resource>) {
        resources.forEach {
            assertEquals(it.traverse("Id"), 0)
        }
    }
}
