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
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_0_8
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_1_7
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_5_1
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_7_1
import com.odim.simulator.tree.RedfishVersion.V1_8_0
import com.odim.simulator.tree.RedfishVersion.V1_8_1
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.RESUME_SUBSCRIPTION
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.EVENT_DESTINATION_SYSLOG_FILTER
import com.odim.simulator.tree.structure.EmbeddedObjectType.HTTP_HEADER_PROPERTY
import com.odim.simulator.tree.structure.EmbeddedObjectType.SNMP_SETTINGS
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CERTIFICATES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ANY
import com.odim.simulator.tree.structure.ResourceType.EVENT_DESTINATION
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT_DEFINITION

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(EVENT_DESTINATION)
open class EventDestinationTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Event Destination Description",
                "Name" to "Event Destination",
                "Destination" to "",
                "EventTypes" to embeddedArray(),
                "Context" to null,
                "Protocol" to "Redfish",
                "HttpHeaders" to embeddedArray(HTTP_HEADER_PROPERTY),
                "Actions" to Actions(
                        Action(RESUME_SUBSCRIPTION)
                )
        ))
        version(V1_0_2, V1_0_0)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_1_0, V1_0_2, resourceObject(
                "OriginResources" to LinkableResourceArray(ANY),
                "MessageIds" to embeddedArray()
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_1_7, V1_1_6)
        version(V1_2_0, V1_1_2, resourceObject(
                "Actions" to Actions()
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_3_0, V1_2_2, resourceObject(
                "SubscriptionType" to null
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_4_0, V1_3_0, resourceObject(
                "RegistryPrefixes" to embeddedArray(),
                "ResourceTypes" to embeddedArray(),
                "SubordinateResources" to null,
                "EventFormatType" to null
        ))
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_5_0, V1_4_1)
        version(V1_5_1, V1_5_0)
        version(V1_6_0, V1_5_0, resourceObject(
                "DeliveryRetryPolicy" to null,
                "Status" to embeddedObject(STATUS),
                "MetricReportDefinitions" to LinkableResourceArray(METRIC_REPORT_DEFINITION)
        ))
        version(V1_6_1, V1_6_0)
        version(V1_7_0, V1_6_0, resourceObject(
                "SNMP" to embeddedObject(SNMP_SETTINGS)
        ))
        version(V1_7_1, V1_7_0)
        version(V1_8_0, V1_7_0, resourceObject(
                "IncludeOriginOfCondition" to null
        ))
        version(V1_8_1, V1_8_0)
        version(V1_9_0, V1_8_1, resourceObject(
                "Certificates" to ResourceCollection(CERTIFICATES_COLLECTION),
                "VerifyCertificate" to null,
                "SyslogFilters" to embeddedArray(EVENT_DESTINATION_SYSLOG_FILTER),
                "OEMProtocol" to "",
                "OEMSubscriptionType" to ""
        ))
    }
}
