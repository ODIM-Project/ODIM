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

package com.odim.simulator.tree.templates.bmc

import com.odim.simulator.tree.RedfishVersion.V1_0_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.TEMPERATURE
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0

@Template(TEMPERATURE)
open class TemperatureBmcTemplate : ResourceTemplate() {
    init {
        version(BMC_1_0, V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Name" to null,
                "SensorNumber" to null,
                "Status" to embeddedObject(STATUS),
                "ReadingCelsius" to null,
                "UpperThresholdNonCritical" to null,
                "UpperThresholdCritical" to null,
                "LowerThresholdNonCritical" to null,
                "LowerThresholdCritical" to null
        ))
    }
}
