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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CALCULATION_PARAMS_TYPE
import com.odim.simulator.tree.structure.EmbeddedObjectType.METRIC_DEFINITION_WILDCARD
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.METRIC_DEFINITION

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(METRIC_DEFINITION)
open class MetricDefinitionTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Metric Definition Description",
                "Name" to "Metric Definition",
                "MetricType" to null,
                "MetricDataType" to null,
                "Units" to null,
                "Implementation" to null,
                "Calculable" to null,
                "IsLinear" to null,
                "Wildcards" to embeddedArray(METRIC_DEFINITION_WILDCARD),
                "MetricProperties" to embeddedArray(),
                "CalculationParameters" to embeddedArray(CALCULATION_PARAMS_TYPE),
                "PhysicalContext" to null,
                "SensingInterval" to null,
                "DiscreteValues" to embeddedArray(),
                "Precision" to null,
                "Accuracy" to null,
                "Calibration" to null,
                "TimestampAccuracy" to null,
                "MinReadingRange" to null,
                "MaxReadingRange" to null,
                "CalculationAlgorithm" to null,
                "CalculationTimeInterval" to null,
                "Actions" to Actions()
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_1_0, V1_0_5, resourceObject(
                "OEMCalculationAlgorithm" to null
        ))
    }
}
