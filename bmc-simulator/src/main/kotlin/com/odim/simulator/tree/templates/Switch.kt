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
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.PORTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.CHASSIS

fun switch() = resourceObject(
        "Id"                    to 0,
        "Name"                  to "Switch",
        "Description"           to "Switch description",
        "SwitchType"            to "PCIe",
        "Status"                to embeddedObject(STATUS),
        "Manufacturer"          to "Microsemi",
        "Model"                 to "Switchtec PSX",
        "SKU"                   to "sku-as-string",
        "SerialNumber"          to "XAJ70002712345",
        "PartNumber"            to "10",
        "AssetTag"              to "SwitchTAG",
        "DomainID"              to 1,
        "IsManaged"             to false,
        "TotalSwitchWidth"      to 80,
        "IndicatorLED"          to null,
        "PowerState"            to "On",
        "Ports"                 to ResourceCollection(PORTS_COLLECTION),
        "Redundancy"            to embeddedArray(),
        "Links"                 to embeddedObject(
                "Chassis"               to LinkableResourceArray(CHASSIS),
                "Oem"                   to embeddedObject()
        ),
        "Actions"               to Actions(
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
        "Oem"                   to embeddedObject()
)
