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
import com.odim.simulator.tree.structure.EmbeddedObjectType.NVME_CONTROLLER_PROPERTIES
import com.odim.simulator.tree.structure.EmbeddedObjectType.PCIE_INTERFACE
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_IDENTIFIER
import com.odim.simulator.tree.structure.EmbeddedObjectType.RESOURCE_LOCATION
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedObjectType.STORAGE_CONTROLLER_CACHE_SUMMARY
import com.odim.simulator.tree.structure.EmbeddedObjectType.STORAGE_CONTROLLER_RATES
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.PORTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ASSEMBLY
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.STORAGE_CONTROLLER
import com.odim.simulator.tree.structure.ResourceType.VOLUME
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(STORAGE_CONTROLLER)
open class StorageControllerTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Storage Controller Description",
                "Name" to "Storage Controller",
                "Status" to embeddedObject(STATUS),
                "SpeedGbps" to null,
                "FirmwareVersion" to null,
                "Manufacturer" to null,
                "Model" to null,
                "SKU" to null,
                "SerialNumber" to null,
                "PartNumber" to null,
                "AssetTag" to null,
                "SupportedControllerProtocols" to embeddedArray(),
                "SupportedDeviceProtocols" to embeddedArray(),
                "Identifiers" to embeddedArray(RESOURCE_IDENTIFIER),
                "Location" to embeddedObject(RESOURCE_LOCATION),
                "Assembly" to SingletonResource(ASSEMBLY),
                "CacheSummary" to embeddedObject(STORAGE_CONTROLLER_CACHE_SUMMARY),
                "PCIeInterface" to embeddedObject(PCIE_INTERFACE),
                "SupportedRAIDTypes" to embeddedArray(),
                "Ports" to ResourceCollection(PORTS_COLLECTION),
                "ControllerRates" to embeddedObject(STORAGE_CONTROLLER_RATES),
                "NVMeControllerProperties" to embeddedObject(NVME_CONTROLLER_PROPERTIES),
                "Actions" to Actions(),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "Endpoints" to LinkableResourceArray(ENDPOINT),
                        "PCIeFunctions" to LinkableResourceArray(PCIE_FUNCTION),
                        "AttachedVolumes" to LinkableResourceArray(VOLUME)
                )
        ))
    }
}
