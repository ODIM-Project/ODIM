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

package component.bmc.patches

import com.odim.odimclient.RedfishHttpClient
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.repo.bmc.BMCSimulator
import component.ComponentSimulatorTest
import org.testng.Assert.assertEquals
import org.testng.annotations.BeforeClass
import org.testng.annotations.Test

class ChassisPatches : ComponentSimulatorTest() {
    private lateinit var simulator: BMCSimulator
    private lateinit var ip: String

    @BeforeClass
    fun config() {
        simulator = BMCSimulator()
        val (ip, _) = provider.serve(simulator, 12345)
        this.ip = ip
        this.httpClient = RedfishHttpClient(ip) {
            basicCredentialEnabled = true
            basicCredentials = "${simulator.basicAuthUsername}:${simulator.basicAuthPassword}"
        }
    }

    @Test
    fun `Patching Chassis should update patched properties`() {
        val response = httpClient.patch("Chassis/Baseboard", makeJson {
            "AssetTag" to "test"
            "PhysicalSecurity" to {
                "IntrusionSensor" to "xxx"
            }
            "Location" to {
                "Info" to "aaa"
                "InfoFormat" to "bbb"
                "PostalAddress" to {
                    "Country" to "PL"
                }
                "Placement" to {
                    "RackOffset" to 4
                    "Row" to "3"
                    "Rack" to "6"
                }
                "PartLocation" to {
                    "ServiceLabel" to "label"
                    "LocationOrdinalValue" to 2
                }
                "Latitude" to 12.0005
            }
        })
        assertEquals(response.code, OK)
        checkResourcePropertiesSet("Chassis/Baseboard") {
            "AssetTag" to "test"
            "PhysicalSecurity" to {
                "IntrusionSensor" to "xxx"
            }
            "Location" to {
                "Info" to "aaa"
                "InfoFormat" to "bbb"
                "PostalAddress" to {
                    "Country" to "PL"
                }
                "Placement" to {
                    "RackOffset" to 4
                    "Row" to "3"
                    "Rack" to "6"
                }
                "PartLocation" to {
                    "ServiceLabel" to "label"
                    "LocationOrdinalValue" to 2
                }
                "Latitude" to 12.0005
            }
        }
    }

    @Test
    fun `Patching Power should update patched properties`() {
        val response = httpClient.patch("Chassis/Baseboard/Power", makeJson {
            "PowerControl" to array[
                    {
                        "PowerLimit" to {
                            "LimitException" to 10
                            "CorrectionInMs" to 110
                        }
                    }
            ]
        })
        assertEquals(response.code, OK)
        checkResourcePropertiesSet("Chassis/Baseboard/Power") {
            "PowerControl" to array[
                    {
                        "PowerLimit" to {
                            "LimitException" to 10
                            "CorrectionInMs" to 110
                        }
                    }
            ]
        }
    }
}
