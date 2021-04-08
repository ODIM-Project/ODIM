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

package com.odim.simulator.repo.bmc.configurators

import com.odim.simulator.repo.bmc.utils.PropertyValueGenerator.Factory.generateDateTime
import com.odim.simulator.repo.bmc.utils.PropertyValueGenerator.Factory.generateEntryLogMessage
import com.odim.simulator.repo.bmc.utils.PropertyValueGenerator.Factory.generateIpAddressPart
import com.odim.simulator.repo.bmc.utils.PropertyValueGenerator.Factory.generateMacAddress
import com.odim.simulator.repo.bmc.utils.PropertyValueGenerator.Factory.generateSerial
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import java.util.UUID.randomUUID
import kotlin.random.Random

@SuppressWarnings("LongMethod", "LargeClass", "StringLiteralDuplication")
class ManagerConfigurator private constructor() {
    companion object Factory {
        fun configureManager(manager: Resource): Resource {
            val uuid = randomUUID().toString()
            manager {
                "UUID" to uuid
                "ServiceEntryPointUUID" to uuid
                "PowerState" to "On"
                "DateTime" to "2006-01-02T15:04:05-07:00"
                "Model" to "SimulatorBmc"
                "FirmwareVersion" to "wht-0.59-0-gbb4e18-83ecb75"
                "Status" to {
                    "State" to "Enabled"
                    "Health" to "OK"
                    "HealthRollup" to "OK"
                }
                "GraphicalConsole" to {
                    "ServiceEnabled" to true
                    "MaxConcurrentSessions" to Random.nextInt(0, 10)
                    "ConnectTypesSupported" to array["KVMIP"]
                }
                "SerialConsole" to {
                    "ServiceEnabled" to true
                    "MaxConcurrentSessions" to Random.nextInt(0, 10)
                    "ConnectTypesSupported" to array["SSH"]
                }
            }

            val networkProtocol = manager.traverse<Resource>("NetworkProtocol")
            networkProtocol {
                "Status" to {
                    "State" to "Enabled"
                    "Health" to "OK"
                    "HealthRollup" to "OK"
                }
                "FQDN" to "bmcsimulator.customcompany.com"
                "HostName" to "DCMIA4BF014B4" + Random.nextInt(100, 300)
                "HTTPS" to {
                    "ProtocolEnabled" to true
                    "Port" to 443
                }
                "SSH" to {
                    "ProtocolEnabled" to false
                    "Port" to 66
                }
                "KVMIP" to {
                    "ProtocolEnabled" to false
                    "Port" to 5902
                }
                "VirtualMedia" to {
                    "ProtocolEnabled" to false
                    "Port" to 627
                }
                "SSDP" to {
                    "ProtocolEnabled" to false
                    "Port" to 1900
                }
                "Telnet" to {
                    "ProtocolEnabled" to false
                }
                "DHCP" to {
                    "ProtocolEnabled" to false
                }
                "NTP" to {
                    "ProtocolEnabled" to false
                }
                "DHCPv6" to {
                    "ProtocolEnabled" to false
                }
                "RDP" to {
                    "ProtocolEnabled" to false
                }
                "RFB" to {
                    "ProtocolEnabled" to false
                }
            }

            return manager
        }

        fun configureLogServiceForManager(logService: Resource, baseLogEntries: MutableList<Resource>): Resource {
            logService {
                "OverWritePolicy" to "WrapsWhenFull"
                "MaxNumberOfRecords" to Random.nextInt(100, 500)
                "DateTime" to generateDateTime()
            }

            baseLogEntries.map {
                it {
                    "Message" to generateEntryLogMessage()
                    "Created" to generateDateTime()
                }
                logService.traverse<ResourceCollection>("Entries").add(it)
            }

            return logService
        }

        fun configureEthernetInterfaceForManager(ethIfc: Resource): Resource {
            val ipBeginPart = generateIpAddressPart()
            val ipLastPart = Random.nextInt(2, 255)
            val macAddress = generateMacAddress()

            return ethIfc {
                "Id" to "eth0"
                "AutoNeg" to true
                "FQDN" to "bmcsimulator.customcompany.com"
                "Status" to {
                    "State" to "Enabled"
                    "Health" to "OK"
                    "HealthRollup" to "OK"
                }
                "UefiDevicePath" to "Acpi(HWP0002,PNP0A03,1)"
                "FullDuplex" to true
                "MACAddress" to macAddress
                "PermanentMACAddress" to macAddress
                "IPv4Addresses" to array[{
                    "Address" to generateIpAddressPart() + ipLastPart
                }]
                "HostName" to generateSerial()
                "IPv4Addresses" to array[{
                    "Address" to "$ipBeginPart.$ipLastPart"
                    "SubnetMask" to "255.255.254.0"
                    "Gateway" to "$ipBeginPart.1"
                    "AddressOrigin" to "DHCP"
                }]
                "IPv4StaticAddresses" to array[{
                    "Address" to "$ipBeginPart.$ipLastPart"
                    "SubnetMask" to "255.255.254.0"
                    "Gateway" to "$ipBeginPart.1"
                }]
                "DHCPv4" to {
                    "DHCPEnabled" to true
                }
                "MTUSize" to 1500
                "MaxIPv6StaticAddresses" to 1
                "IPv6Addresses" to array[{
                    "Address" to "::"
                    "PrefixLength" to 64L
                    "AddressOrigin" to "Static"
                }]
                "IPv6StaticAddresses" to array[{
                    "Address" to "::"
                    "PrefixLength" to 64L
                }]
                "IPv6DefaultGateway" to "::"
                "IPv6StaticDefaultGateways" to array[{
                    "Address" to "::"
                    "PrefixLength" to 64L
                }]
                "DHCPv6" to {
                    "OperatingMode" to "Disabled"
                }
                "LinkStatus" to "LinkUp"
                "InterfaceEnabled" to true
                "SpeedMbps" to Random.nextInt(50, 100)
                "VLAN" to {
                    "VLANEnable" to false
                    "VLANId" to 0
                }
            }
        }
    }
}
