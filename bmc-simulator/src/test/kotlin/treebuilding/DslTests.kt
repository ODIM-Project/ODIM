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
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE
import org.testng.Assert.assertEquals
import org.testng.Assert.assertNotNull
import org.testng.Assert.assertTrue
import org.testng.annotations.Test

class DslTests {
    @Test
    fun `Create resource and override its properties`() {
        with(RedfishSimulator()) {
            val system = create(COMPUTER_SYSTEM) {
                "Status" to {
                    "Health" to "Sick"
                    "HealthRollup" to "Unrolled"
                    "State" to "Indeterminate"
                }
            }
            assertEquals(system.traverse("Status/Health"), "Sick")
            assertEquals(system.traverse("Status/HealthRollup"), "Unrolled")
            assertEquals(system.traverse("Status/State"), "Indeterminate")
            system {
                "Status" to {
                    "State" to "Disabled"
                }
            }
            assertEquals(system.traverse("Status/State"), "Disabled")
            assertEquals(system.traverse("Status/Health"), "Sick", "Overriding properties should not erase unreferenced properties (it's additive)")
        }
    }

    @Test
    fun `Create and combine resources at a later time`() {
        with(RedfishSimulator()) {
            val system = create(COMPUTER_SYSTEM) { "Name" to "System1" }
            val processor = create(PROCESSOR) { "Name" to "Cpu1" }
            root(
                    system(processor)
            )
            assertEquals(system.traverse("Name"), "System1")
            val processors = system.traverse<ResourceCollection>("Processors")
            assertNotNull(processors.members.singleOrNull { it.traverse<String>("Name") == "Cpu1" })
        }
    }

    @Test
    fun `Create and combine resources at once`() {
        with(RedfishSimulator()) {
            val system = create(COMPUTER_SYSTEM) { "Name" to "System1" }
            root(
                    system(
                            create(ETHERNET_INTERFACE) { "Name" to "Eth1" }
                    )
            )
            assertEquals(system.traverse("Name"), "System1")
            val interfaces = system.traverse<ResourceCollection>("EthernetInterfaces")
            val eth22 = interfaces.members.singleOrNull { it.traverse<String>("Name") == "Eth1" }
            assertNotNull(eth22)
        }
    }

    @Test
    fun `Override nested array property`() {
        with(RedfishSimulator()) {
            val system = create(COMPUTER_SYSTEM) {
                "Oem" to {
                    "Custom_Company" to {
                        "PCIeConnectionId" to array[
                                "Id1",
                                "Id2",
                                "Id3"
                        ]
                    }
                }
            }

            fun getConnectionIds() = system.traverse<List<*>>("Oem/Custom_Company/PCIeConnectionId")
            assertEquals(getConnectionIds(), listOf("Id1", "Id2", "Id3"))
            system {
                "Oem" to {
                    "Custom_Company" to {
                        "PCIeConnectionId" to array[empty]
                    }
                }
            }
            assertEquals(getConnectionIds().size, 0)
        }
    }

    @Test
    fun `Create multiple resources in one go`() {
        with(RedfishSimulator()) {
            val (sys1, sys2, sys3) = createMany(COMPUTER_SYSTEM, COMPUTER_SYSTEM, COMPUTER_SYSTEM)
            root(sys1, sys2, sys3)
            assertEquals(root.traverse<ResourceCollection>("Systems").members.size, 3)
        }
    }

    @Test
    fun `Create and memorize multiple resources`() {
        val (system, chassis, manager) = RedfishSimulator().memorize {
            root(
                    create(COMPUTER_SYSTEM) { "Name" to "System1" },
                    +create(COMPUTER_SYSTEM) { "Name" to "System2" },
                    create(STORAGE_SERVICE),
                    +create(CHASSIS) { "Name" to "Chassis1" },
                    create(CHASSIS) { "Name" to "Chassis2" },
                    create(MANAGER) { "Name" to "Manager1" },
                    +create(MANAGER) { "Name" to "Manager2" }
            )
        }
        assertEquals(system.traverse("Name"), "System2")
        assertEquals(chassis.traverse("Name"), "Chassis1")
        assertEquals(manager.traverse("Name"), "Manager2")
    }

    @Test
    fun `Create and memorize resources on multiple depths, retaining top-down memorization order`() {
        val rsrc = RedfishSimulator().memorize {
            +root { "Name" to "root" }(
                    create(COMPUTER_SYSTEM)(
                            create(ETHERNET_INTERFACE)
                    ),
                    +create(COMPUTER_SYSTEM) { "Name" to "System1" }(
                            +create(ETHERNET_INTERFACE) { "Name" to "Eth1" },
                            create(ETHERNET_INTERFACE)
                    ),
                    +create(CHASSIS) { "Name" to "Chassis1" },
                    create(MANAGER)(
                            create(ETHERNET_INTERFACE)
                    ),
                    +create(MANAGER) { "Name" to "Manager1" }(
                            +create(ETHERNET_INTERFACE) { "Name" to "Eth2" }
                    )
            )
        }
        assertEquals(rsrc.map { it.traverse<String>("Name") }, listOf("root", "System1", "Eth1", "Chassis1", "Manager1", "Eth2"))
    }

    @Test
    fun `Create complex tree with resources memorized on different levels and branches`() {
        val rsrc = RedfishSimulator().memorize {
            +root { "Name" to "root" }(
                    create(COMPUTER_SYSTEM),
                    +create(COMPUTER_SYSTEM) { "Name" to "System1" }(
                            +create(ETHERNET_INTERFACE) { "Name" to "Eth1" },
                            create(ETHERNET_INTERFACE)
                    ),
                    +create(CHASSIS) { "Name" to "Chassis1" },
                    create(MANAGER),
                    +create(MANAGER) { "Name" to "Manager1" }(
                            +create(ETHERNET_INTERFACE) { "Name" to "Eth2" }
                    )
            )
            +create(MANAGER) { "Name" to "Manager2" }(
                    create(ETHERNET_INTERFACE) { "Name" to "Eth3" },
                    +create(ETHERNET_INTERFACE) { "Name" to "Eth4" }
            )
        }
        assertEquals(rsrc.map { it.traverse<String>("Name") }, listOf("root", "System1", "Eth1", "Chassis1", "Manager1", "Eth2", "Manager2", "Eth4"))
    }

    @Test
    fun `Create unconnected trees and confirm memorized resources are in top-down order`() {
        val rsrc = RedfishSimulator().memorize {
            create(COMPUTER_SYSTEM)(
                    create(ETHERNET_INTERFACE),
                    create(ETHERNET_INTERFACE),
                    +create(ETHERNET_INTERFACE) { "Name" to "Eth1" },
                    create(ETHERNET_INTERFACE)
            )
            +create(CHASSIS) { "Name" to "Chassis1" }
            create(MANAGER) { "Name" to "m19" }(
                    +create(ETHERNET_INTERFACE) { "Name" to "Eth2" }
            )
        }
        assertEquals(rsrc.map { it.traverse<String>("Name") }, listOf("Eth1", "Chassis1", "Eth2"))
    }

    @Test
    fun `Create complex DSL with nested objects within array`() {
        with(RedfishSimulator()) {
            val processor = create(PROCESSOR) {
                "Name" to "SpeedSelect_2"
                "Oem" to {
                    "Custom_Company" to {
                        "SpeedSelect" to {
                            "CurrentConfiguration" to "1"
                            "Configurations" to array[
                                    {
                                        "ConfigurationId" to 1
                                        "HighPriorityCoreCount" to 4
                                        "HighPriorityBaseFrequency" to 25
                                        "LowPriorityCoreCount" to 0
                                        "TDP" to 120
                                        "LowPriorityBaseFrequency" to "null"
                                        "MaxJunctionTempCelsius" to 90L
                                    }
                            ]
                        }
                    }
                }
            }
            assertEquals(processor.traverse("Name"), "SpeedSelect_2")
            assertEquals(processor.traverse("Oem/Custom_Company/SpeedSelect/Configurations/0/HighPriorityCoreCount"), 4)
            assertEquals(processor.traverse("Oem/Custom_Company/SpeedSelect/Configurations/0/MaxJunctionTempCelsius"), 90L)
        }
    }

    @Test
    fun `DSL can override and define properties to null`() {
        with (RedfishSimulator()) {
            val (processor) = memorize {
                create(COMPUTER_SYSTEM) (
                        +create(PROCESSOR) {
                            "Oem" to {
                                "Custom_Company" to {
                                    "AlternativeFirmware" to null
                                "Brand" to "X505"
                            }
                        }
                    }
                )
            }
            assertEquals(processor.traverse<Any?>("Oem/Custom_Company/AlternativeFirmware"), null)
            assertEquals(processor.traverse<Any?>("Oem/Custom_Company/Brand"), "X505")
            processor {
                "Oem" to {
                    "Custom_Company" to {
                        "Brand" to null
                    }
                }
            }
            assertEquals(processor.traverse<Any?>("Oem/Custom_Company/Brand"), null)
        }
    }

    @Test
    fun `DSL can override property to an empty object`() {
        with (RedfishSimulator()) {
            val (processor) = memorize {
                create(COMPUTER_SYSTEM)(
                        +create(PROCESSOR) {
                            "Oem" to {
                                "Custom_Company" to {
                                    "AlternativeFirmware" to null
                                }
                            }
                        }
                )
            }
            assertEquals(processor.traverse<Any?>("Oem/Custom_Company/AlternativeFirmware"), null)
            processor {
                "Oem" to {
                    "Custom_Company" to {
                        "AlternativeFirmware" to {}
                    }
                }
            }
            assertEquals(processor.traverse("Oem/Custom_Company/AlternativeFirmware"), mapOf<Any?, Any?>())
        }
    }

    @Test
    fun `Override keeps resource property types declared in templates`() {
        with (RedfishSimulator()) {
            val processor = create(PROCESSOR)
            create(COMPUTER_SYSTEM) (processor)
            assertTrue(processor.traverse<Any>("HighSpeedCoreIDs") is EmbeddedArray<*>)
            processor {
                "HighSpeedCoreIDs" to array[empty]
            }
            assertTrue(processor.traverse<Any>("HighSpeedCoreIDs") is EmbeddedArray<*>)
        }
    }
}
