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

package com.odim.simulator.tree.templates

import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.NETWORK_ADAPTERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.PCIE_DEVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_SWITCH
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.POWER
import com.odim.simulator.tree.structure.ResourceType.RESOURCE_BLOCK
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.SWITCH
import com.odim.simulator.tree.structure.ResourceType.THERMAL
import com.odim.simulator.tree.structure.SingletonResource
import java.util.UUID.randomUUID

fun chassis() = resourceObject(
        "Id" to 0,
        "ChassisType" to "Module",
        "Name" to "Chassis",
        "Description" to "Chassis description",
        "Manufacturer" to null,
        "Model" to null,
        "SKU" to "sku-as-string",
        "SerialNumber" to null,
        "PartNumber" to null,
        "UUID" to randomUUID().toString(),
        "AssetTag" to null,
        "IndicatorLED" to null,
        "PowerState" to "On",
        "Status" to embeddedObject(STATUS),
        "PCIeDevices" to ResourceCollection(PCIE_DEVICES_COLLECTION),
        "PhysicalSecurity" to embeddedObject(
                "IntrusionSensorNumber" to 64,
                "IntrusionSensor" to "Normal",
                "IntrusionSensorReArm" to "Manual"

        ),
        "HeightMm"              to 2000,
        "WidthMm"               to 1000,
        "DepthMm"               to 1000,
        "WeightKg"              to 500,
        "NetworkAdapters"       to ResourceCollection(NETWORK_ADAPTERS_COLLECTION),
        "Links"                 to embeddedObject(
                "@odata.type"           to "#Chassis.v1_2_0.Links",
                "ComputerSystems"       to LinkableResourceArray(COMPUTER_SYSTEM),
                "ManagedBy"             to LinkableResourceArray(MANAGER),
                "ManagersInChassis"     to LinkableResourceArray(MANAGER),
                "ContainedBy"           to LinkableResource(CHASSIS),
                "Contains"              to LinkableResourceArray(CHASSIS),
                "PoweredBy"             to SingletonResource(POWER),
                "CooledBy"              to SingletonResource(THERMAL),
                "Drives"                to LinkableResourceArray(DRIVE),
                "Storage"               to LinkableResourceArray(STORAGE),
                "PCIeDevices"           to LinkableResourceArray(PCIE_DEVICE),
                "ResourceBlocks"        to LinkableResourceArray(RESOURCE_BLOCK),
                "EthernetSwitches"      to LinkableResourceArray(ETHERNET_SWITCH),
                "Switches"              to LinkableResourceArray(SWITCH),
                "Oem"                   to embeddedObject()
        ),
        "Actions" to Actions(
                Action(RESET, "ResetType", mutableListOf(
                        "On",
                        "ForceOff",
                        "GracefulShutdown",
                        "GracefulRestart",
                        "ForceRestart",
                        "Nmi",
                        "ForceOn",
                        "PushPowerButton",
                        "PowerCycle"
                ))
        ),
        "Oem" to embeddedObject(
                "Custom_Company" to embeddedObject(
                        "@odata.type" to "#CustomCompany.Oem.Chassis",
                        "Location" to embeddedObject(
                                "Id" to "Drawer1",
                                "ParentId" to null
                        )
                )
        )
)
