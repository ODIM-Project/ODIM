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
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.SUBMIT_TEST_METRIC_REPORT
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.METRIC_DEFINITIONS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.METRIC_REPORTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.METRIC_REPORT_DEFINITIONS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.TRIGGERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.LOG_SERVICE
import com.odim.simulator.tree.structure.ResourceType.TELEMETRY_SERVICE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(TELEMETRY_SERVICE)
open class TelemetryServiceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to "TelemetryService",
                "Description" to "Telemetry Service Description",
                "Name" to "Telemetry Service",
                "Status" to embeddedObject(STATUS),
                "MaxReports" to null,
                "MinCollectionInterval" to null,
                "SupportedCollectionFunctions" to embeddedArray(),
                "MetricDefinitions" to ResourceCollection(METRIC_DEFINITIONS_COLLECTION),
                "MetricReportDefinitions" to ResourceCollection(METRIC_REPORT_DEFINITIONS_COLLECTION),
                "MetricReports" to ResourceCollection(METRIC_REPORTS_COLLECTION),
                "Triggers" to ResourceCollection(TRIGGERS_COLLECTION),
                "LogService" to LinkableResource(LOG_SERVICE),
                "Actions" to Actions(
                        Action(SUBMIT_TEST_METRIC_REPORT)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_1_0, V1_0_0)
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_2_0, V1_1_2, resourceObject(
                "ServiceEnabled" to null
        ))
        version(V1_2_1, V1_2_0)
    }
}
