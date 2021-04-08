package com.odim.simulator.tree.templates.bmc

import com.odim.simulator.tree.RedfishVersion.V1_8_2
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.FIRMWARE_INVENTORIES_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.UPDATE_SERVICE
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0
import com.odim.simulator.tree.templates.redfish.UpdateServiceTemplate

@Template(UPDATE_SERVICE)
open class UpdateServiceBmcTemplate: UpdateServiceTemplate() {
    init {
        version(
            BMC_1_0, V1_8_2, resourceObject(
                "FirmwareInventory" to ResourceCollection(FIRMWARE_INVENTORIES_COLLECTION)
                )
        )
    }
}
