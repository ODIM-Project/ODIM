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
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpMethod.PATCH
import com.odim.simulator.http.HttpStatusCode.BAD_REQUEST
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.http.Request
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.FABRIC
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.ZONE
import org.testng.Assert.assertEquals
import org.testng.Assert.assertFalse
import org.testng.Assert.assertTrue
import org.testng.annotations.BeforeMethod
import org.testng.annotations.Test

class ResourceOverrideTest {
    lateinit var simulator: ResourceOverrideSimulator

    @BeforeMethod
    fun prepareTest() {
        simulator = ResourceOverrideSimulator()
    }

    @Test
    fun `Patching embedded array of objects creates objects from template`() {
        val onLocationPatchJSON = makeJson {
            "Location" to array[
                    { "Info" to null },
                    { "InfoFormat" to null },
                    { "Oem" to null },
                    { "PostalAddress" to null },
                    { "Placement" to null },
                    { "PartLocation" to null },
                    { "Longitude" to 0.00 },
                    { "Latitude" to 0.00 },
                    { "AltitudeMeters" to 0.00 },
                    { "Contacts" to null }
            ]
        }.toString()
        with (simulator) {
            val request = Request(PATCH, onLocationPatchJSON)
            val response = behaviors.createResponse(tree, drive, request)
            assertEquals(response.code, OK)
            val referenceObject = createEmbeddedObject(RESOURCE_LOCATION)
            assertEquals(drive.traverse<List<*>>("Location").size, 10)
            for (n in 0 until 10) {
                val protocolN = drive.traverse<Map<*, *>>("Location/$n")
                assertEquals(protocolN.keys, referenceObject.keys, "When overriding array of EmbeddedObjects with new members, their keys should be aligned with template")
            }
        }
    }

    @Test
    fun `Patching another embedded array of objects creates objects from template`() {
        val onIdentifiersPatchJSON = makeJson {
            "Identifiers" to array[
                    { "DurableName" to null },
                    { "DurableNameFormat" to null }
            ]
        }.toString()
        with (simulator) {
            val request = Request(PATCH, onIdentifiersPatchJSON)
            val response = behaviors.createResponse(tree, drive, request)
            assertEquals(response.code, OK)
            val referenceObject = createEmbeddedObject(RESOURCE_IDENTIFIER)
            assertEquals(drive.traverse<List<*>>("Identifiers").size, 2)
            for (n in 0 until 2) {
                val configurationN = drive.traverse<Map<*, *>>("Identifiers/$n")
                assertEquals(configurationN.keys, referenceObject.keys, "When overriding array of EmbeddedObjects with new members, their keys should be aligned with template")
            }
        }
    }

    @Test
    fun `Discarding link to Resource removes link in manager simulator`() {
        val unlinkChassisPatch = makeJson {
            "Links" to {
                "Chassis" to array[empty]
            }
        }.toString()
        with (simulator) {
            assertTrue(areLinked(chassis, system))
            val request = Request(PATCH, unlinkChassisPatch)
            val response = behaviors.createResponse(tree, system, request)
            assertEquals(response.code, OK)
            assertFalse(areLinkedOnAnySide(chassis, system))
        }
    }

    @Test
    fun `Replacing link to Resource removes old and add new link in simulator`() {
        with (simulator) {
            val relinkChassisPatch = makeJson {
                "Links" to {
                    "Chassis" to array[{
                        "@odata.id" to chassis2.toLink()
                    }]
                }
            }.toString()
            assertTrue(areLinked(system, chassis))
            assertFalse(areLinkedOnAnySide(system, chassis2))
            val request = Request(PATCH, relinkChassisPatch)
            val response = behaviors.createResponse(tree, system, request)
            assertEquals(response.code, OK)
            assertFalse(areLinkedOnAnySide(system, chassis))
            assertTrue(areLinked(system, chassis2))
        }
    }

    @Test
    fun `Overriding list of links with empty array removes links in simulator`() {
        with (simulator) {
            val removeAllLinksJson = makeJson {
                "Links" to {
                    "Endpoints" to array[empty]
                }
            }.toString()
            assertTrue(areLinkedOnAnySide(zone1, endpoint1))
            val request = Request(PATCH, removeAllLinksJson)
            val response = behaviors.createResponse(tree, zone1, request)
            assertEquals(response.code, OK)
            assertFalse(areLinkedOnAnySide(zone1, endpoint1))
        }
    }

    @Test
    fun `Overriding list of links with partially new entries removes links that were defined`() {
        with (simulator) {
            val replaceLinkJson = makeJson {
                "Links" to {
                    "Endpoints" to array[
                        { "@odata.id" to endpoint2.toLink() }
                    ]
                }
            }.toString()
            assertFalse(areLinkedOnAnySide(zone2, endpoint2))
            val replaceLinkRequest = Request(PATCH, replaceLinkJson)
            val replaceLinkResponse = behaviors.createResponse(tree, zone1, replaceLinkRequest)
            assertEquals(replaceLinkResponse.code, OK)
            assertFalse(areLinkedOnAnySide(zone2, endpoint2))
            assertFalse(areLinkedOnAnySide(zone2, endpoint1))
        }
    }

    @Test
    fun `Trying to override array with object fails`() {
        val destroyLinksChassisPatch = makeJson {
            "Links" to {
                "Chassis" to {
                    "@odata.id" to null
                }
            }
        }.toString()
        with (simulator) {
            val request = Request(PATCH, destroyLinksChassisPatch)
            val response = behaviors.createResponse(tree, system, request)
            assertEquals(response.code, BAD_REQUEST)
        }
    }

    @Test
    fun `Trying to override object with array fails`() {
        val destroyStatePatch = makeJson {
            "Status" to array[
                    { "State" to null },
                    { "HealthRollup" to null }
            ]
        }.toString()
        with (simulator) {
            val request = Request(PATCH, destroyStatePatch)
            val response = behaviors.createResponse(tree, zone1, request)
            assertEquals(response.code, BAD_REQUEST)
        }
    }

    @Test
    fun `Trying to override Resource embedded within another fails`() {
        val destroyMetricsResourcePatch = makeJson {
            "Links" to {
                "PoweredBy" to {
                    "ProcessorSummary" to null
                }
            }
        }.toString()
        with (simulator) {
            val request = Request(PATCH, destroyMetricsResourcePatch)
            val response = behaviors.createResponse(tree, chassis, request)
            assertEquals(response.code, BAD_REQUEST)
        }
    }

    @Test
    fun `Trying to override ResourceCollection fails`() {
        val destroyCollectionPatch = makeJson {
            "PCIeDevices" to null
        }.toString()
        with (simulator) {
            val request = Request(PATCH, destroyCollectionPatch)
            val response = behaviors.createResponse(tree, system, request)
            assertEquals(response.code, BAD_REQUEST)
        }
    }

    @Test
    fun `Trying to override Array of Resources embedded in another Resource fails`() {
        val destroyPowerControlPatch = makeJson {
            "PowerControl" to null
        }.toString()
        with (simulator) {
            val power = chassis.traverse<Resource>("Power")
            val request = Request(PATCH, destroyPowerControlPatch)
            val response = behaviors.createResponse(tree, power, request)
            assertEquals(response.code, BAD_REQUEST)
        }
    }
}

class ResourceOverrideSimulator: RedfishSimulator() {
    val system = create(COMPUTER_SYSTEM) {
        "Name" to "System1"
        "Oem" to {
            "Custom_Company" to {
                "PerformanceConfiguration" to {
                    "CurrentConfigurationId" to "0"
                    "Configurations" to array[
                            {
                                "@odata.type" to "#CustomCompany.Oem.SpeedSelectConfiguration"
                                "ConfigurationId" to "0"
                                "Type" to "StaticSpeedSelect"
                                "TDPPerCpu" to 120
                                "MaxCpuJunctionTemp" to 90
                                "ActiveCoresPerCpu" to 10
                                "BaseCoreFrequency" to 1700
                            },
                            {
                                "@odata.type" to "#CustomCompany.Oem.SpeedSelectConfiguration"
                                "ConfigurationId" to "1"
                                "Type" to "StaticSpeedSelect"
                                "TDPPerCpu" to 110
                                "MaxCpuJunctionTemp" to 80
                                "ActiveCoresPerCpu" to 8
                                "BaseCoreFrequency" to 1600L
                            }
                    ]
                }
            }
        }
    }
    val processor = create(PROCESSOR)
    val drive = create(DRIVE)
    val manager = create(MANAGER) { "Name" to "Manager1" }
    val chassis = create(CHASSIS)
    val chassis2 = create(CHASSIS)
    val endpoint1 = create(ENDPOINT) { "Name" to "Endpoint1" }
    val endpoint2 = create(ENDPOINT) { "Name" to "Endpoint2" }
    val fabric = create(FABRIC)
    val zone1 = create(ZONE) { "Name" to "Zone1" }
    val zone2 = create(ZONE) { "Name" to "Zone2" }

    init {
        root(
                system(
                    processor
                ),
                manager,
                chassis,
                chassis2,
                fabric(
                    zone1,
                    zone2,
                    endpoint1,
                    endpoint2
                )
        )
        link(chassis, system)
        link(zone1, endpoint1)
        link(chassis, manager, "ManagedBy", "ManagerInChassis")
    }
}
