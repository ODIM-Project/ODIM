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

import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType
import kotlin.random.Random
import kotlin.random.Random.Default.nextInt

@Suppress("LargeClass")
class ChassisConfigurator private constructor() {
    companion object Factory {
        private const val MAIN_CHASSIS_ID = "RackMount"
        fun configureMainChassis(chassis: Resource): Resource {
            chassis {
                "Id" to "RackMount"
                "Name" to "Computer System Chassis"
                "Status" to {
                    "State" to "Enabled"
                    "Health" to "OK"
                    "HealthRollup" to "OK"
                }
                "AssetTag" to "RM001"
                "EnvironmentalClass" to "A3"
                "ChassisType" to MAIN_CHASSIS_ID
                "Manufacturer" to "Custom Corporation"
                "Model" to "S2600WFT"
                "SerialNumber" to "BQWF73800" + nextInt(250, 350)
                "PartNumber" to "H48104-" + nextInt(250, 350)
                "PowerState" to "Off"
                "IndicatorLED" to "Off"
            }

            chassis.data.remove("Power")
            chassis.data.remove("Thermal")

            return chassis
        }

        private fun configureTemperature(temperature: Resource) = temperature {
            "Status" to {
                "State" to "Enabled"
                "Health" to "OK"
                "HealthRollup" to "OK"
            }
            "SensorNumber" to nextInt(1, 50)
            "ReadingCelsius" to nextInt(35, 65)
            "LowerThresholdFatal" to nextInt(0, 5)
            "LowerThresholdCritical" to nextInt(5, 20)
            "LowerThresholdNonCritical" to nextInt(20, 40)
            "UpperThresholdNonCritical" to nextInt(40, 60)
            "UpperThresholdCritical" to nextInt(60, 80)
            "UpperThresholdFatal" to nextInt(80, 120)
        }

        private fun configureFan(fan: Resource) = fan {
            "Name" to "MainFan"
            "FanName" to "MainFan"
            "Status" to {
                "State" to "Enabled"
                "Health" to "OK"
                "HealthRollup" to "OK"
            }
            "Manufacturer" to "SilentFans"
            "Model" to "Cool4k"
            "LowerThresholdFatal" to nextInt(0, 100)
            "LowerThresholdCritical" to nextInt(100, 400)
            "LowerThresholdNonCritical" to nextInt(400, 1000)
            "UpperThresholdNonCritical" to nextInt(1000, 1600)
            "UpperThresholdCritical" to nextInt(1600, 1800)
            "UpperThresholdFatal" to nextInt(1800, 2500)
            "MinReadingRange" to nextInt(800, 900)
            "MaxReadingRange" to nextInt(1600, 1700)
            "SensorNumber" to nextInt(1, 10)
            "Reading" to nextInt(900, 1600)
            "ReadingUnits" to "RPM"
            "PhysicalContext" to "Back"
        }

        private fun configurePowerControl(powerControl: Resource) = powerControl {
            "PowerConsumedWatts" to nextInt(1, 50)
            "PowerMetrics" to {
                "IntervalInMin" to nextInt(1, 30000)
                "MinConsumedWatts" to nextInt(1, 50)
                "MaxConsumedWatts" to nextInt(1, 50)
                "AverageConsumedWatts" to nextInt(1, 50)
            }
        }

        private fun configurePowerSupply(powerSupply: Resource) = powerSupply {
            "Name" to "PowerSupply"
            "Status" to {
                "State" to "Enabled"
                "Health" to "OK"
                "HealthRollup" to "OK"
            }
            "Model" to "S040EV1200250"
            "PartNumber" to "PSU300566"
            "FirmwareVersion" to "1.0"
            "Manufacturer" to "GigaWattSolutions"
            "PowerCapacityWatts" to 5000
            "PowerSupplyType" to "AC"
            "LastPowerOutputWatts" to 3500
            "LineInputVoltage" to 277
            "LineInputVoltageType" to "AC277V"
        }

        private fun configureRedundancy(redundancy: Resource): Resource {
            redundancy {
                "Mode" to "N+m"
                "MinNumNeeded" to nextInt(1, 5)
                "MaxNumSupported" to nextInt(1, 5)
            }

            redundancy.data.remove("RedundancyEnabled")
            redundancy.data.remove("Actions")
            redundancy.data.remove("Oem")

            return redundancy
        }

        private fun configureVoltage(voltage: Resource): Resource {
            val minRange = Random.nextFloat()
            val lowerTreshold = Random.nextFloat()
            voltage {
                "Name" to "Voltage"
                "Status" to {
                    "State" to "Enabled"
                    "Health" to "OK"
                    "HealthRollup" to "OK"
                }
                "SensorNumber" to nextInt(1, 10)
                "ReadingVolts" to Random.nextFloat()
                "LowerThresholdNonCritical" to lowerTreshold + Random.nextFloat()
                "LowerThresholdCritical" to lowerTreshold
                "UpperThresholdCritical" to lowerTreshold + 3 * Random.nextFloat()
                "UpperThresholdNonCritical" to lowerTreshold + 2 * Random.nextFloat()
                "MinReadingRange" to minRange
                "MaxReadingRange" to minRange + Random.nextFloat()
            }
            voltage.data.remove("UpperThresholdFatal")
            voltage.data.remove("LowerThresholdFatal")
            voltage.data.remove("Actions")
            voltage.data.remove("Oem")

            return voltage
        }

        fun appendPower(chassis: Resource, baseRedundancy: Resource, basePowerControl: Resource,
                        basePowerSupply: Resource?, system: Resource, chassisRackMount: Resource): Resource {
            val powerControl = configurePowerControl(basePowerControl)
            val power = chassis.traverse<Resource>("Power")
            power.traverse<EmbeddedResourceArray>("PowerControl").add(powerControl)
            basePowerSupply?.let {
                power.traverse<EmbeddedResourceArray>("PowerSupplies").add(configurePowerSupply(basePowerSupply))
            }
            powerControl.traverse<LinkableResourceArray>("RelatedItem").addLink(system, powerControl)
            powerControl.traverse<LinkableResourceArray>("RelatedItem").addLink(chassisRackMount, powerControl)

            val redundancyPower = configureRedundancy(baseRedundancy)
            redundancyPower.data.remove("MinNumNeeded")
            redundancyPower.data.remove("MaxNumSupported")

            power.traverse<EmbeddedResourceArray>("Redundancy").add(redundancyPower)

            basePowerControl {
                "Name" to "PowerControl"
                "PowerAllocatedWatts" to 5000
                "PowerAvailableWatts" to 3500
                "PowerCapacityWatts" to 5000
                "PowerConsumedWatts" to 1500
                "PowerLimit" to {
                    "LimitException" to 20
                    "CorrectionInMs" to 120
                }
                "PowerMetrics" to {
                    "IntervalInMin" to 5
                    "MinConsumedWatts" to 800
                    "MaxConsumedWatts" to 1800
                    "AverageConsumedWatts" to 1200
                }
            }

            return chassis
        }

        fun appendVoltageVolt(chassis: Resource, baseVoltage: Resource, system: Resource,
                              chassisRackMount: Resource): Resource {
            val power = chassis.traverse<Resource>("Power")
            val voltage = configureVoltage(baseVoltage)
            voltage.traverse<LinkableResourceArray>("RelatedItem").addLink(system, voltage)
            voltage.traverse<LinkableResourceArray>("RelatedItem").addLink(chassisRackMount, voltage)
            power.traverse<EmbeddedResourceArray>("Voltages").add(voltage)

            return chassis
        }

        fun configureChassis(id: String, baseChassis: Resource, tag: String, row: String): Resource {
            baseChassis {
                "Id" to id
                "Name" to "Computer System Card"
                "Status" to {
                    "State" to "Enabled"
                    "Health" to "OK"
                    "HealthRollup" to "OK"
                }
                "AssetTag" to tag
                "ChassisType" to "Card"
                "EnvironmentalClass" to "A3"
                "DepthMm" to 800
                "HeightMm" to 140
                "WeightKg" to 26
                "PowerState" to "On"
                "IndicatorLED" to "Lit"
                "Manufacturer" to "Custom Corporation"
                "Model" to "S2600WFT"
                "SerialNumber" to "BQWF73800" + nextInt(250, 350)
                "PowerState" to "Off"
                "PartNumber" to "H48104-" + nextInt(250, 350)
                "Location" to {
                    "Placement" to {
                        "Row" to row
                        "Rack" to "1"
                    }
                }
                "Links" to {
                    "ContainedBy" to LinkableResourceArray(ResourceType.CHASSIS)
                }
            }

            return baseChassis
        }

        fun appendThermal(chassis: Resource, baseTemperature: List<Resource>, baseFan: Resource, baseRedundancy: Resource) {
            val temperature0 = configureTemperature(baseTemperature[0])
            val temperature1 = configureTemperature(baseTemperature[1])
            val temperature2 = configureTemperature(baseTemperature[2])

            val thermal = chassis.traverse<Resource>("Thermal")
            thermal.traverse<EmbeddedResourceArray>("Temperatures").add(temperature0)
            thermal.traverse<EmbeddedResourceArray>("Temperatures").add(temperature1)
            temperature2.data.remove("LowerThresholdNonCritical")
            temperature2.data.remove("LowerThresholdCritical")
            thermal.traverse<EmbeddedResourceArray>("Temperatures").add(temperature2)

            val fan = configureFan(baseFan)
            thermal.traverse<EmbeddedResourceArray>("Fans").add(fan)

            val redundancyThermal = configureRedundancy(baseRedundancy)
            thermal.traverse<EmbeddedResourceArray>("Redundancy").add(redundancyThermal)
        }
    }
}
