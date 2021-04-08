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
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.CLEAR_LOG
import com.odim.simulator.tree.structure.ActionType.COLLECT_DIAGNOSTIC_DATA
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.LOG_SERVICE_SYSLOG_FILTER
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.LOG_ENTRIES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.LOG_SERVICE

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(LOG_SERVICE)
open class LogServiceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Log Service Description",
                "Name" to "Log Service",
                "ServiceEnabled" to null,
                "MaxNumberOfRecords" to 0,
                "OverWritePolicy" to "Unknown",
                "DateTime" to null,
                "DateTimeLocalOffset" to null,
                "Entries" to ResourceCollection(LOG_ENTRIES_COLLECTION),
                "Status" to embeddedObject(STATUS),
                "Actions" to Actions(
                        Action(CLEAR_LOG)
                )
        ))
        version(V1_0_2, V1_0_0)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_1_0, V1_0_5, resourceObject(
                "LogEntryType" to null
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_2_0, V1_1_3, resourceObject(
                "SyslogFilters" to embeddedArray(LOG_SERVICE_SYSLOG_FILTER),
                "Actions" to Actions(
                        Action(COLLECT_DIAGNOSTIC_DATA)
                )
        ))
    }
}
