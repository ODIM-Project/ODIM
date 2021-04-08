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
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.MEMORY
import com.odim.simulator.tree.structure.ResourceType.MEMORY_METRICS
import com.odim.simulator.tree.structure.SingletonResource
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0

@Template(MEMORY)
open class MemoryBmcTemplate : ResourceTemplate() {
    init {
        version(BMC_1_0, V1_0_0, resourceObject(
                "Id" to 0,
                "Description" to "Memory Description",
                "Name" to "Memory",
                "MemoryType" to null,
                "MemoryDeviceType" to null,
                "BaseModuleType" to null,
                "MemoryMedia" to embeddedArray(),
                "CapacityMiB" to null,
                "DataWidthBits" to null,
                "BusWidthBits" to null,
                "Manufacturer" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "AllowedSpeedsMHz" to embeddedArray(),
                "Metrics" to SingletonResource(MEMORY_METRICS),
                "RankCount" to null,
                "DeviceLocator" to null,
                "ErrorCorrection" to null,
                "OperatingSpeedMhz" to null,
                "Status" to embeddedObject(STATUS),
                "Oem" to embeddedObject(
                        "Custom_Company" to embeddedObject(
                                "@odata.type" to "#CustomCompany.Oem.Memory",
                                "VoltageVolt" to null
                        )
                )
        ))
    }
}
