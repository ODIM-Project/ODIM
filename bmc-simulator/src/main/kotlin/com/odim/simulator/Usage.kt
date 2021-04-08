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

package com.odim.simulator

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.Request
import com.odim.simulator.http.ServerProvider
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollectionType.STORAGE_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.FAN
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.METRIC_DEFINITION
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT_DEFINITION
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.POWER_CONTROL
import com.odim.simulator.tree.structure.ResourceType.POWER_SUPPLY
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE
import com.odim.simulator.tree.structure.ResourceType.TEMPERATURE
import com.odim.simulator.tree.structure.ResourceType.VOLTAGE
import com.odim.simulator.tree.structure.ResourceType.VOLUME
import com.odim.simulator.tree.structure.StringArray
import com.odim.utils.query

fun showDSLOverrideWithNull() {
    with(RedfishSimulator()) {
        val system = create(COMPUTER_SYSTEM) {
            "Oem" to {
                "Custom_Company" to {
                    "PCIeConnectionId" to array[
                            "Tip",
                            "Top"
                    ]
                    "AbsoluteAddressableMemory" to null
                }
            }
        }
        println(this.createResponse(system, Request(GET)).json.query<ObjectNode>("Oem").toString())
        system {
            "Oem" to {
                "Custom_Company" to {
                    "DiscoveryState" to null
                    "AbsoluteAddressableMemory" to 262144
                }
            }
        }
        println(this.createResponse(system, Request(GET)).json.query<ObjectNode>("Oem").toString())
    }
}

fun main(args: Array<String>) {
    showDSLOverrideWithNull()
    with(RedfishSimulator()) {
        val system = create(COMPUTER_SYSTEM) {
            "Status" to {
                "Health" to "Sick"
                "HealthRollup" to "Unrolled"
            }
            "Oem" to {
                "Custom_Company" to {
                    "PCIeConnectionId" to array[
                            "Tip",
                            "Top"
                    ]
                }
            }
        }
        val processor = create(PROCESSOR)
        val ethInf1 = create(ETHERNET_INTERFACE)
        val ethInf2 = create(ETHERNET_INTERFACE)
        val chassis = create(CHASSIS)
        val chassis2 = create(CHASSIS)
        val pcieDevice = create(PCIE_DEVICE)
        pcieDevice.appendTo(chassis)
        link(system, pcieDevice)

        system.append(processor)

        val traverse = system.traverse<Int>("Id")
        val traverse2 = system.traverse<ResourceObject>("Oem")
        system.traverse<StringArray>("Oem/Custom_Company/PCIeConnectionId").add("test")
        val traverse4 = system.traverse<String>("Oem/Custom_Company/PCIeConnectionId/0")
        println("traverse = $traverse, traverse2 = $traverse2, traverse4 = $traverse4")

        val ss = create(STORAGE_SERVICE)
        val volume = create(VOLUME)
        ss.append(volume)
        ss.appendTo(root)
        link(system, root.getCollection(STORAGE_SERVICES_COLLECTION))

        val manager = create(MANAGER)
        val manager2 = create(MANAGER)
        val ethInf3 = create(ETHERNET_INTERFACE)
        ethInf3.appendTo(manager).appendTo(root)
        manager2.appendTo(root)

        link(system, manager)
        link(system, manager2)

        ethInf1.appendTo(system).appendTo(root)
        ethInf2.appendTo(system)

        chassis.appendTo(root)
        chassis2.appendTo(root)

        link(system, chassis)
//        ODIM have problem with recognize ContainedBy property
//        link(chassis, chassis2, "ContainedBy", "Contains")

        println(root.print())

        val power = chassis2.traverse<Resource>("Links/PoweredBy")
        val powerControl = power.traverse<EmbeddedResourceArray>("PowerControl")
        val pc = create(POWER_CONTROL)
        link(pc, chassis)
        powerControl.add(pc)
        val pc2 = create(POWER_CONTROL)
        powerControl.add(pc2)

        val voltages = power.traverse<EmbeddedResourceArray>("Voltages")
        val volt = create(VOLTAGE)
        voltages.add(volt)
        val volt2 = create(VOLTAGE)
        voltages.add(volt2)

        val pSupps = power.traverse<EmbeddedResourceArray>("PowerSupplies")
        val pSupp = create(POWER_SUPPLY)
        pSupps.add(pSupp)

        val thermal = chassis2.traverse<Resource>("Links/CooledBy")
        thermal.traverse<EmbeddedResourceArray>("Temperatures").add(create(TEMPERATURE))
        thermal.traverse<EmbeddedResourceArray>("Fans").add(create(FAN))

        val telemetryService = root.traverse<Resource>("TelemetryService")
        val metricDefinition = create(METRIC_DEFINITION)
        val metricReportDefinition = create(METRIC_REPORT_DEFINITION)
        val metricReport = create(METRIC_REPORT)
        telemetryService(
                metricDefinition,
                metricReportDefinition,
                metricReport
        )

        ServerProvider().serve(this, 3000)
    }
}
