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
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.CHASSIS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.COMPUTER_SYSTEMS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MANAGERS_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MESSAGE_REGISTRY_FILES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.SESSIONS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ACCOUNT_SERVICE
import com.odim.simulator.tree.structure.ResourceType.CERTIFICATE_SERVICE
import com.odim.simulator.tree.structure.ResourceType.EVENT_SERVICE
import com.odim.simulator.tree.structure.ResourceType.SERVICE_ROOT
import com.odim.simulator.tree.structure.ResourceType.SESSION_SERVICE
import com.odim.simulator.tree.structure.ResourceType.TASK_SERVICE
import com.odim.simulator.tree.structure.ResourceType.UPDATE_SERVICE
import com.odim.simulator.tree.structure.SingletonResource
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0
import java.util.UUID.randomUUID

@Template(SERVICE_ROOT)
open class ServiceRootBmcTemplate : ResourceTemplate() {
    init {
        version(BMC_1_0, V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Service Root Description",
                "Name" to "Service Root",
                "RedfishVersion" to "1.9.0",
                "UUID" to randomUUID().toString(),
                "Systems" to ResourceCollection(COMPUTER_SYSTEMS_COLLECTION),
                "Chassis" to ResourceCollection(CHASSIS_COLLECTION),
                "Managers" to ResourceCollection(MANAGERS_COLLECTION),
                "SessionService" to SingletonResource(SESSION_SERVICE),
                "AccountService" to SingletonResource(ACCOUNT_SERVICE),
                "UpdateService" to SingletonResource(UPDATE_SERVICE),
                "EventService" to SingletonResource(EVENT_SERVICE),
                "CertificateService" to SingletonResource(CERTIFICATE_SERVICE),
                "TaskService" to SingletonResource(TASK_SERVICE),
                "Registries" to ResourceCollection(MESSAGE_REGISTRY_FILES_COLLECTION),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Sessions" to ResourceCollection(SESSIONS_COLLECTION)
                )
        ))
    }
}
