// Copyright (c) Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package com.odim.simulator.tree.templates.redfish

import com.odim.simulator.tree.RedfishVersion.V1_0_0
import com.odim.simulator.tree.RedfishVersion.V1_0_1
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESET_METRICS
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.SENSOR_THRESHOLDS
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.SENSOR

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(SENSOR)
open class SensorTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Sensor Description",
                "Name" to "Sensor",
                "ReadingType" to null,
                "DataSourceUri" to null,
                "Status" to embeddedObject(STATUS),
                "Reading" to null,
                "ReadingUnits" to null,
                "PhysicalContext" to null,
                "PhysicalSubContext" to null,
                "PeakReading" to null,
                "MaxAllowableOperatingValue" to null,
                "MinAllowableOperatingValue" to null,
                "AdjustedMaxAllowableOperatingValue" to null,
                "AdjustedMinAllowableOperatingValue" to null,
                "ApparentVA" to null,
                "ReactiveVAR" to null,
                "PowerFactor" to null,
                "LoadPercent" to null,
                "Location" to embeddedObject(RESOURCE_LOCATION),
                "ElectricalContext" to null,
                "VoltageType" to null,
                "Thresholds" to embeddedObject(SENSOR_THRESHOLDS),
                "ReadingRangeMax" to null,
                "ReadingRangeMin" to null,
                "Precision" to null,
                "Accuracy" to null,
                "SensingFrequency" to null,
                "PeakReadingTime" to null,
                "SensorResetTime" to null,
                "Actions" to Actions(
                        Action(RESET_METRICS)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_1_0, V1_0_4, resourceObject(
                "CrestFactor" to null,
                "THDPercent" to null,
                "LifetimeReading" to null,
                "SensingInterval" to null,
                "ReadingTime" to null,
                "Implementation" to null
        ))
        version(V1_1_1, V1_1_0)
    }
}
