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

package com.odim.simulator.repo.bmc.behaviors

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.success
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.utils.getNumberOrNull
import com.odim.utils.getObjectOrNull
import com.odim.utils.getString
import com.odim.utils.getStringOrNull

class PatchOnChassis : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val chassis = item as Resource
        updateSystemIndicatorLed(request, item)
        updateAssetTag(request, chassis)
        updatePhysicalSecurity(request, chassis)
        updateLocation(request, chassis)
        return terminal(success(chassis))
    }

    private fun updateLocation(request: Request, chassis: Resource) {
        request.json?.getObjectOrNull("Location")?.let { json ->
            patchFirstLevelLocation(json.getStringOrNull("Info"), "Info", chassis)
            patchFirstLevelLocation(json.getStringOrNull("InfoFormat"), "InfoFormat", chassis)
            patchFirstLevelLocation(json.getNumberOrNull("Longitude"), "Longitude", chassis)
            patchFirstLevelLocation(json.getNumberOrNull("Latitude"), "Latitude", chassis)
            patchFirstLevelLocation(json.getNumberOrNull("AltitudeMeters"), "AltitudeMeters", chassis)
            updateLocationPostalAddress(json, chassis)
            updateLocationPlacement(json, chassis)
            updateLocationPartLocation(json, chassis)
        }
    }

    private fun updateLocationPartLocation(json: ObjectNode, chassis: Resource) {
        json.getObjectOrNull("PartLocation")?.let { partLocationJson ->
            listOf(
                    "ServiceLabel",
                    "LocationType",
                    "Reference",
                    "Orientation"
            ).forEach {
                patchSecondLevelLocation(partLocationJson.getStringOrNull(it), "PartLocation", it, chassis)
            }
            listOf(
                    "LocationOrdinalValue"
            ).forEach {
                patchSecondLevelLocation(partLocationJson.getNumberOrNull(it), "PartLocation", it, chassis)
            }
        }
    }

    private fun updateLocationPlacement(json: ObjectNode, chassis: Resource) {
        json.getObjectOrNull("Placement")?.let { placementJson ->
            listOf(
                    "Row",
                    "Rack",
                    "RackOffsetUnits",
                    "AdditionalInfo"
            ).forEach {
                patchSecondLevelLocation(placementJson.getStringOrNull(it), "Placement", it, chassis)
            }
            listOf(
                    "RackOffset"
            ).forEach {
                patchSecondLevelLocation(placementJson.getNumberOrNull(it), "Placement", it, chassis)
            }
        }
    }

    @Suppress("LongMethod")
    private fun updateLocationPostalAddress(json: ObjectNode, chassis: Resource) {
        json.getObjectOrNull("PostalAddress")?.let { postalAddJson ->
            listOf(
                    "Country",
                    "Territory",
                    "District",
                    "City",
                    "Division",
                    "Neighborhood",
                    "LeadingStreetDirection",
                    "Street",
                    "TrailingStreetSuffix",
                    "StreetSuffix",
                    "HouseNumberSuffix",
                    "Landmark",
                    "Floor",
                    "Name",
                    "PostalCode",
                    "Building",
                    "Unit",
                    "Room",
                    "Seat",
                    "PlaceType",
                    "Community",
                    "POBox",
                    "AdditionalCode",
                    "Road",
                    "RoadSection",
                    "RoadBranch",
                    "RoadSubBranch",
                    "RoadPreModifier",
                    "RoadPostModifier",
                    "GPSCoords",
                    "AdditionalInfo"
            ).forEach {
                patchSecondLevelLocation(postalAddJson.getStringOrNull(it), "PostalAddress", it, chassis)
            }
            listOf(
                    "HouseNumber"
            ).forEach {
                patchSecondLevelLocation(postalAddJson.getNumberOrNull(it), "PostalAddress", it, chassis)
            }
        }
    }

    private fun updatePhysicalSecurity(request: Request, chassis: Resource) {
        request.json?.getObjectOrNull("PhysicalSecurity")?.let {
            it.getStringOrNull("IntrusionSensor")?.let {
                chassis {
                    "PhysicalSecurity" to {
                        "IntrusionSensor" to it
                    }
                }
            }
        }
    }

    private fun updateAssetTag(request: Request, chassis: Resource) {
        request.json?.getStringOrNull("AssetTag")?.let {
            chassis { "AssetTag" to it }
        }
    }

    private fun patchFirstLevelLocation(value: Any?, key: String, chassis: Resource) {
        value?.let {
            chassis {
                "Location" to {
                    key to it
                }
            }
        }
    }

    private fun patchSecondLevelLocation(value: Any?, group: String, key: String, chassis: Resource) {
        value?.let {
            chassis {
                "Location" to {
                    group to {
                        key to it
                    }
                }
            }
        }
    }

    private fun updateSystemIndicatorLed(request: Request, item: Item) {
        val indicatorFieldName = "IndicatorLED"
        val newIndicatorState = request.json?.getString(indicatorFieldName)

        newIndicatorState.let {
            val chassis = item as Resource
            val system = chassis.traverseOrNull<LinkableResourceArray>("Links/ComputerSystems")?.first as Resource

            system {
                indicatorFieldName to newIndicatorState
            }
        }
    }
}
