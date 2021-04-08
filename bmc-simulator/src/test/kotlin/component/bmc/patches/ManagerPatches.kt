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

class ManagerPatches : ComponentSimulatorTest() {
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
    fun `Patching Manager should update patched properties`() {
        val response = httpClient.patch("Managers/1", makeJson {
            "CommandShell" to {
                "ServiceEnabled" to true
            }
            "GraphicalConsole" to {
                "ServiceEnabled" to false
            }
            "SerialConsole" to {
                "ServiceEnabled" to false
            }
        })
        assertEquals(response.code, OK)
        checkResourcePropertiesSet("Managers/1") {
            "CommandShell" to {
                "ServiceEnabled" to true
            }
            "GraphicalConsole" to {
                "ServiceEnabled" to false
            }
            "SerialConsole" to {
                "ServiceEnabled" to false
            }
        }
    }

    @Test
    fun `Patching NetworkProtocol should update patched properties`() {
        val response = httpClient.patch("Managers/1/NetworkProtocol", makeJson {
            "FQDN" to "new fqdn"
            "HTTPS" to {
                "ProtocolEnabled" to false
                "Port" to 222
            }
            "NTP" to {
                "ProtocolEnabled" to true
                "Port" to 999
                "NTPServers" to array["abc", "def"]
            }
            "SSDP" to {
                "ProtocolEnabled" to true
                "Port" to 2000
                "NotifyIPv6Scope" to "Link"
                "NotifyMulticastIntervalSeconds" to 10
                "NotifyTTL" to 10
            }
            "SSH" to {
                "ProtocolEnabled" to true
                "Port" to 111
            }
        })
        assertEquals(response.code, OK)
        checkResourcePropertiesSet("Managers/1/NetworkProtocol") {
            "FQDN" to "new fqdn"
            "HTTPS" to {
                "ProtocolEnabled" to false
                "Port" to 222
            }
            "NTP" to {
                "ProtocolEnabled" to true
                "Port" to 999
                "NTPServers" to array["abc", "def"]
            }
            "SSDP" to {
                "ProtocolEnabled" to true
                "Port" to 2000
                "NotifyIPv6Scope" to "Link"
                "NotifyMulticastIntervalSeconds" to 10
                "NotifyTTL" to 10
            }
            "SSH" to {
                "ProtocolEnabled" to true
                "Port" to 111
            }
        }
    }

    @Test
    fun `Patching EthernetInterface should update patched properties`() {
        val response = httpClient.patch("Managers/1/EthernetInterfaces/eth0", makeJson {
            "FQDN" to "new fqdn"
            "AutoNeg" to false
            "FullDuplex" to false
            "HostName" to "houstnejm"
            "IPv6DefaultGateway" to "::0"
            "InterfaceEnabled" to false
            "MTUSize" to 10
            "IPv4Addresses" to array[
                    {
                        "Address" to "251.243.1.111"
                        "SubnetMask" to "255.255.254.0"
                        "Gateway" to "251.243.1.2"
                        "AddressOrigin" to "DHCP"
                    },
                    {
                        "Address" to "251.243.1.112"
                        "SubnetMask" to "255.255.254.0"
                        "Gateway" to "251.243.1.3"
                        "AddressOrigin" to "DHCP"
                    }
            ]
            "IPv4StaticAddresses" to array[
                    {
                        "Address" to "124.177.85.11"
                        "SubnetMask" to "255.255.254.0"
                        "Gateway" to "124.177.85.1"
                    }
            ]
            "IPv6Addresses" to array[
                    {
                        "Address" to "::BB"
                        "PrefixLength" to 64
                        "AddressOrigin" to "Static"
                    },
                    {
                        "Address" to "AA::"
                        "PrefixLength" to 64
                        "AddressOrigin" to "Static"
                    }
            ]
            "IPv6StaticAddresses" to array[
                    {
                        "Address" to "::XX"
                        "PrefixLength" to 64
                    }
            ]
            "IPv6StaticDefaultGateways" to array[
                    {
                        "Address" to "TT::"
                        "PrefixLength" to 62
                    },
                    {
                        "Address" to "::LL"
                        "PrefixLength" to 64
                    }
            ]
            "IPv6AddressPolicyTable" to array[
                    {
                        "Prefix" to "pre"
                        "Precedence" to 1
                        "Label" to 1
                    }
            ]
        })
        assertEquals(response.code, OK)
        checkResourcePropertiesSet("Managers/1/EthernetInterfaces/eth0") {
            "FQDN" to "new fqdn"
            "AutoNeg" to false
            "FullDuplex" to false
            "HostName" to "houstnejm"
            "IPv6DefaultGateway" to "::0"
            "InterfaceEnabled" to false
            "MTUSize" to 10
            "IPv4Addresses" to array[
                    {
                        "Address" to "251.243.1.111"
                        "SubnetMask" to "255.255.254.0"
                        "Gateway" to "251.243.1.2"
                        "AddressOrigin" to "DHCP"
                    },
                    {
                        "Address" to "251.243.1.112"
                        "SubnetMask" to "255.255.254.0"
                        "Gateway" to "251.243.1.3"
                        "AddressOrigin" to "DHCP"
                    }
            ]
            "IPv4StaticAddresses" to array[
                    {
                        "Address" to "124.177.85.11"
                        "SubnetMask" to "255.255.254.0"
                        "Gateway" to "124.177.85.1"
                    }
            ]
            "IPv6Addresses" to array[
                    {
                        "Address" to "::BB"
                        "PrefixLength" to 64
                        "AddressOrigin" to "Static"
                    },
                    {
                        "Address" to "AA::"
                        "PrefixLength" to 64
                        "AddressOrigin" to "Static"
                    }
            ]
            "IPv6StaticAddresses" to array[
                    {
                        "Address" to "::XX"
                        "PrefixLength" to 64
                    }
            ]
            "IPv6StaticDefaultGateways" to array[
                    {
                        "Address" to "TT::"
                        "PrefixLength" to 62
                    },
                    {
                        "Address" to "::LL"
                        "PrefixLength" to 64
                    }
            ]
            "IPv6AddressPolicyTable" to array[
                    {
                        "Prefix" to "pre"
                        "Precedence" to 1
                        "Label" to 1
                    }
            ]
        }
    }
}
