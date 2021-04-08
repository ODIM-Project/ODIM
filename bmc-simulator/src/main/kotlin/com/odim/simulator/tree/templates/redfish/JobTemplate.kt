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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.JOB_PAYLOAD
import com.odim.simulator.tree.structure.EmbeddedObjectType.MESSAGE
import com.odim.simulator.tree.structure.EmbeddedObjectType.SCHEDULE
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.JOBS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.JOB

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(JOB)
open class JobTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Job Description",
                "Name" to "Job",
                "JobStatus" to "OK",
                "JobState" to "New",
                "StartTime" to "2017-04-14T06:35:05Z",
                "EndTime" to "2017-04-14T06:35:05Z",
                "MaxExecutionTime" to null,
                "PercentComplete" to null,
                "CreatedBy" to "",
                "Schedule" to embeddedObject(SCHEDULE),
                "HidePayload" to false,
                "Payload" to embeddedObject(JOB_PAYLOAD),
                "Steps" to ResourceCollection(JOBS_COLLECTION),
                "StepOrder" to embeddedArray(),
                "Messages" to embeddedArray(MESSAGE),
                "Actions" to Actions()
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
    }
}
