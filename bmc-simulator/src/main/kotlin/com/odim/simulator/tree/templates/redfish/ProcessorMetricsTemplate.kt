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
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.CACHE_METRICS
import com.odim.simulator.tree.structure.EmbeddedObjectType.CORE_METRICS
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR_METRICS

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(PROCESSOR_METRICS)
open class ProcessorMetricsTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to "ProcessorMetrics",
                "Description" to "Processor Metrics Description",
                "Name" to "Processor Metrics",
                "BandwidthPercent" to null,
                "AverageFrequencyMHz" to null,
                "ThrottlingCelsius" to null,
                "TemperatureCelsius" to null,
                "ConsumedPowerWatt" to null,
                "FrequencyRatio" to null,
                "Cache" to embeddedArray(CACHE_METRICS),
                "LocalMemoryBandwidthBytes" to null,
                "RemoteMemoryBandwidthBytes" to null,
                "KernelPercent" to null,
                "UserPercent" to null,
                "CoreMetrics" to embeddedArray(CORE_METRICS),
                "Actions" to Actions()
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_1_0, V1_0_2, resourceObject(
                "OperatingSpeedMHz" to null
        ))
        version(V1_1_1, V1_1_0)
    }
}
