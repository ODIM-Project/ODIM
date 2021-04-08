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
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.PCIE_INTERFACE
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollectionType.PCIE_FUNCTIONS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ASSEMBLY
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(PCIE_DEVICE)
open class PCIeDeviceTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "P C Ie Device Description",
                "Name" to "P C Ie Device",
                "Manufacturer" to null,
                "Model" to null,
                "SKU" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "AssetTag" to null,
                "DeviceType" to "SingleFunction",
                "FirmwareVersion" to null,
                "Status" to embeddedObject(STATUS),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Chassis" to LinkableResourceArray(CHASSIS),
                        "PCIeFunctions" to LinkableResourceArray(PCIE_FUNCTION)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_1_0, V1_0_2, resourceObject(
                "Actions" to Actions()
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_2_0, V1_1_0, resourceObject(
                "Assembly" to SingletonResource(ASSEMBLY)
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_3_0, V1_2_1, resourceObject(
                "PCIeInterface" to embeddedObject(PCIE_INTERFACE)
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_4_0, V1_3_2, resourceObject(
                "PCIeFunctions" to LinkableResource(PCIE_FUNCTIONS_COLLECTION)
        ))
        version(V1_5_0, V1_4_0, resourceObject(
                "UUID" to null
        ))
    }
}
