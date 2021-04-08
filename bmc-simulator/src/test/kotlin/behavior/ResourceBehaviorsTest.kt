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

package behavior

import com.odim.simulator.RedfishSimulator
import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.Behavior.Companion.behavior
import com.odim.simulator.behaviors.BehaviorDataStore.SharedInformationType.BIOS_SETTINGS
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.behaviors.ResourceBehaviors
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpMethod.DELETE
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.HttpMethod.PATCH
import com.odim.simulator.http.HttpMethod.POST
import com.odim.simulator.http.HttpStatusCode.METHOD_NOT_ALLOWED
import com.odim.simulator.http.HttpStatusCode.NOT_FOUND
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response.Companion.notAllowed
import com.odim.simulator.http.Response.Companion.notFound
import com.odim.simulator.tree.RedfishVersion.Companion.REDFISH_VERSION_LATEST
import com.odim.simulator.tree.ResourceFactory
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollectionType.COMPUTER_SYSTEMS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.ETHERNET_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.FAN
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.POWER_CONTROL
import com.odim.simulator.tree.structure.ResourceType.POWER_SUPPLY
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE
import com.odim.simulator.tree.structure.ResourceType.TEMPERATURE
import com.odim.simulator.tree.structure.ResourceType.VLAN_NETWORK_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.VOLTAGE
import com.odim.utils.JsonMapper.jsonMapper
import com.nhaarman.mockitokotlin2.any
import com.nhaarman.mockitokotlin2.mock
import com.nhaarman.mockitokotlin2.never
import com.nhaarman.mockitokotlin2.times
import com.nhaarman.mockitokotlin2.verify
import com.nhaarman.mockitokotlin2.whenever
import org.testng.Assert.assertEquals
import org.testng.Assert.fail
import org.testng.annotations.BeforeMethod
import org.testng.annotations.Test
import java.io.IOException

val NOT_FOUND_BEHAVIOR = behavior { terminal(notFound()) }
val NOT_ALLOWED_BEHAVIOR = behavior { terminal(notAllowed()) }

class ResourceBehaviorsTest {

    private lateinit var behaviors: ResourceBehaviors
    private lateinit var tree: ResourceTree

    @BeforeMethod
    fun before() {
        tree = ResourceTree(ResourceFactory(listOf("redfish"), REDFISH_VERSION_LATEST))
        behaviors = ResourceBehaviors()
    }

    @Test
    fun `REST DELETE  removes system with id=1`() {
        val systemToDelete = tree.create(COMPUTER_SYSTEM)
        tree.root.append(systemToDelete)

        val req = Request(DELETE, "")
        // before delete
        val systemCollection = tree.root.getCollection(COMPUTER_SYSTEMS_COLLECTION)
        assertEquals(systemCollection.members.count(), 1)
        // delete
        behaviors.createResponse(tree, systemToDelete, req)
        // after delete
        assertEquals(systemCollection.members.count(), 0)
    }

    @Test
    fun `REST PATCH behavior updates system Name and Oem-Custom_Company-UserModeEnabled`() {
        val systemToPatch = tree.create(COMPUTER_SYSTEM)
        tree.root.append(systemToPatch)
        val systemPatchJson = makeJson {
            "Name" to "System_updated_by_PATCH"
            "Oem" to {
                "Custom_Company" to {
                    "UserModeEnabled" to true
                }
            }
        }.toString()

        val req = Request(PATCH, systemPatchJson)
        behaviors.createResponse(tree, systemToPatch, req)

        assertEquals(systemToPatch.traverse("Name"), "System_updated_by_PATCH")
        assertEquals(systemToPatch.traverse("Oem/Custom_Company/UserModeEnabled"), true)
    }

    @Test
    fun `REST PATCH behavior does not alter other fields`() {

        val systemToPatch = tree.create(COMPUTER_SYSTEM) {
            "Name" to "PatchSystem1"
            "Description" to "PATCH verification system"
            "Oem" to {
                "Custom_Company" to {
                    "TrustedExecutionTechnologyEnabled" to false
                    "UserModeEnabled" to true
                }
            }
        }
        tree.root.append(systemToPatch)
        assertEquals(systemToPatch.traverse("Name"), "PatchSystem1")
        assertEquals(systemToPatch.traverse("Description"), "PATCH verification system")
        assertEquals(systemToPatch.traverse("Oem/Custom_Company/UserModeEnabled"), true)
        assertEquals(systemToPatch.traverse("Oem/Custom_Company/TrustedExecutionTechnologyEnabled"), false)
        val systemPatchJson = makeJson {
            "Description" to "PATCHED"
            "Oem" to {
                "Custom_Company" to {
                    "TrustedExecutionTechnologyEnabled" to true
                }
            }
        }.toString()
        val req = Request(PATCH, systemPatchJson)
        behaviors.createResponse(tree, systemToPatch, req)
        assertEquals(systemToPatch.traverse("Description"), "PATCHED")
        assertEquals(systemToPatch.traverse("Oem/Custom_Company/TrustedExecutionTechnologyEnabled"), true)
        assertEquals(systemToPatch.traverse("Name"), "PatchSystem1")
        assertEquals(systemToPatch.traverse("Oem/Custom_Company/UserModeEnabled"), true)

    }

    @Test
    fun `REST POST behavior adds ethernet interface`() {
        val system = tree.create(COMPUTER_SYSTEM)
        tree.root.append(system)
        val ethInfCollection = system.getCollection(ETHERNET_INTERFACES_COLLECTION)
        val ethernetInterfaceString = makeJson {
            "Name" to "EI_1"
            "Version" to 10
            "MacAddress" to "11:22:33:44:55:66"
        }.toString()

        val req = Request(POST, ethernetInterfaceString)
        behaviors.createResponse(tree, ethInfCollection, req)
        assertEquals(ethInfCollection.members.count(), 1)
    }

    @Test
    fun `REST POST behavior adds system to non-empty collection`() {
        tree.root.append(tree.create(COMPUTER_SYSTEM))
        val systemPostJson = makeJson {
            "Name" to "System_from_POST"
            "Oem" to {
                "Custom_Company" to {
                    "UserModeEnabled" to true
                }
            }
        }.toString()

        val req = Request(POST, systemPostJson)
        val systemCollection = tree.root.getCollection(COMPUTER_SYSTEMS_COLLECTION)
        behaviors.createResponse(tree, systemCollection, req)
        assertEquals(systemCollection.members.count(), 2)
        val recentlyAddedSystem = systemCollection.members.first { it.traverse<String>("Name") == "System_from_POST" }
        assertEquals(recentlyAddedSystem.traverse("Id"), 2)
        assertEquals(recentlyAddedSystem.traverse("Oem/Custom_Company/UserModeEnabled"), true)
    }

    @Test
    fun `REST GET behavior is added by default`() {

        val storageService = tree.create(STORAGE_SERVICE)

        val response = behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(response.code, OK)
    }

    @Test
    fun `bind NotFound behavior to resource override defaults REST GET`() {
        val storageService = tree.create(STORAGE_SERVICE)
        behaviors.appendBehavior(storageService, GET, NOT_FOUND_BEHAVIOR)

        val response = behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(response.code, NOT_FOUND)
    }

    @Test
    fun `bind NotFound behavior to resource type override default REST GET`() {
        val storageService = tree.create(STORAGE_SERVICE)
        behaviors.appendBehavior(STORAGE_SERVICE, GET, NOT_FOUND_BEHAVIOR)

        val response = behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(response.code, NOT_FOUND)
    }

    @Test
    fun `final 405 behavior is not overridden by 404 behavior`() {
        val storageService = tree.create(STORAGE_SERVICE)
        behaviors.replaceBehavior(STORAGE_SERVICE, GET, NOT_ALLOWED_BEHAVIOR)
        behaviors.appendBehavior(STORAGE_SERVICE, GET, NOT_FOUND_BEHAVIOR)

        val response = behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(response.code, METHOD_NOT_ALLOWED)
    }

    @Test
    fun `last final 404 behavior is taken`() {
        val storageService = tree.create(STORAGE_SERVICE)
        behaviors.replaceBehavior(STORAGE_SERVICE, GET, NOT_ALLOWED_BEHAVIOR)
        behaviors.replaceBehavior(STORAGE_SERVICE, GET, NOT_FOUND_BEHAVIOR)

        val response = behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(response.code, NOT_FOUND)
    }

    @Test
    fun `prepend action is invoke before default get behavior`() {
        val storageService = tree.create(STORAGE_SERVICE)
        val checkList = mutableListOf<String>()

        val prependBehavior: Behavior = behavior {
            checkList.add("BEFORE")
            nonTerminal(response)
        }

        val appendBehavior: Behavior = behavior {
            checkList.add("AFTER")
            nonTerminal(response)
        }

        behaviors.prependBehavior(storageService, GET, prependBehavior)
        behaviors.appendBehavior(storageService, GET, appendBehavior)

        behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(checkList, listOf("BEFORE", "AFTER"))
    }

    @Test
    fun `behavior shared data should be accessible for all behaviors`() {
        val storageService = tree.create(STORAGE_SERVICE)
        val quietBoot: String = "QuietBoot"
        var accessible: Boolean = false

        val postBehavior: Behavior = behavior {
            dataStore.insert(BIOS_SETTINGS, quietBoot)
            nonTerminal(response)
        }

        val getBehavior: Behavior = behavior {
            if (dataStore.readAndRemove(BIOS_SETTINGS) == quietBoot) {
                accessible = true
            }
            assertEquals(dataStore.storeSize(), 0)
            nonTerminal(response)
        }

        behaviors.appendBehavior(storageService, POST, postBehavior)
        behaviors.appendBehavior(storageService, GET, getBehavior)

        behaviors.createResponse(tree, storageService, Request(POST))
        behaviors.createResponse(tree, storageService, Request(GET))
        assertEquals(accessible, true)
    }

    @Test
    fun `final behavior is the latest applied`() {
        val storageService = tree.create(STORAGE_SERVICE)

        val beforeFinalBehavior: Behavior = mock { whenever(it.run(any(), any(), any(), any(), any())).thenReturn(nonTerminal(mock())) }
        val finalBehavior: Behavior = mock { whenever(it.run(any(), any(), any(), any(), any())).thenReturn(terminal(mock())) }
        val afterFinalBehavior: Behavior = mock { whenever(it.run(any(), any(), any(), any(), any())).thenReturn(nonTerminal(mock())) }

        behaviors.appendBehavior(storageService, GET, beforeFinalBehavior)
        behaviors.appendBehavior(storageService, GET, finalBehavior)
        behaviors.appendBehavior(storageService, GET, afterFinalBehavior)

        behaviors.createResponse(tree, storageService, Request(GET))

        verify(beforeFinalBehavior, times(1)).run(any(), any(), any(), any(), any())
        verify(finalBehavior, times(1)).run(any(), any(), any(), any(), any())
        verify(afterFinalBehavior, never()).run(any(), any(), any(), any(), any())
    }

    @Test
    fun `default REST GET behavior for several resources`() {
        with(RedfishSimulator()) {
            val (system, processor, ethInf2,
                    storageService, chassis, manager1,
                    ethInf1) = memorize {
                root(
                        +create(COMPUTER_SYSTEM)(+create(PROCESSOR),
                                +create(ETHERNET_INTERFACE)
                        ),
                        +create(STORAGE_SERVICE),
                        +create(CHASSIS),
                        +create(MANAGER)(+create(ETHERNET_INTERFACE)),
                        +create(MANAGER)
                )
            }

            link(system, manager1)
            link(system, chassis)
            val power = chassis.traverse<Resource>("Power")
            val powerControl = power.traverse<EmbeddedResourceArray>("PowerControl")
            val pc = create(POWER_CONTROL)
            link(pc, chassis, propertySecond = "PoweredBy")
            powerControl.add(pc)
            powerControl.add(create(POWER_CONTROL))
            val voltages = power.traverse<EmbeddedResourceArray>("Voltages")
            voltages.add(create(VOLTAGE))
            voltages.add(create(VOLTAGE))
            val pSupps = power.traverse<EmbeddedResourceArray>("PowerSupplies")
            pSupps.add(create(POWER_SUPPLY))
            val thermal = chassis.traverse<Resource>("Thermal")
            thermal.traverse<EmbeddedResourceArray>("Temperatures").add(create(TEMPERATURE))
            thermal.traverse<EmbeddedResourceArray>("Fans").add(create(FAN))
            listOf(root, system, manager1, storageService, processor, ethInf1, ethInf2)
                    .forEach { isResourceValidJson(it) }
        }
    }

    private fun isResourceValidJson(resource: Resource) {
        val response = behaviors.createResponse(tree, resource, Request(GET))
        assertEquals(response.code, OK)

        try {
            jsonMapper.readTree(response.body)
        } catch (e: IOException) {
            fail("Parsing JSON failed.", e)
        }
    }
}
