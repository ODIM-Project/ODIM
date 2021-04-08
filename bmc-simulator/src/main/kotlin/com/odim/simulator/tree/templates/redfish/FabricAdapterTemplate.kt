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
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.FABRIC_ADAPTER_GEN_Z
import com.odim.simulator.tree.structure.EmbeddedObjectType.PCIE_INTERFACE
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.PORTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.FABRIC_ADAPTER

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(FABRIC_ADAPTER)
open class FabricAdapterTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Fabric Adapter Description",
                "Name" to "Fabric Adapter",
                "Status" to embeddedObject(STATUS),
                "Ports" to ResourceCollection(PORTS_COLLECTION),
                "Manufacturer" to null,
                "Model" to null,
                "SKU" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "SparePartNumber" to null,
                "ASICRevisionIdentifier" to null,
                "ASICPartNumber" to null,
                "ASICManufacturer" to null,
                "FirmwareVersion" to null,
                "UUID" to null,
                "PCIeInterface" to embeddedObject(PCIE_INTERFACE),
                "GenZ" to embeddedObject(FABRIC_ADAPTER_GEN_Z),
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Endpoints" to LinkableResourceArray(ENDPOINT)
                )
        ))
    }
}
