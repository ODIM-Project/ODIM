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

package com.odim.simulator.repo.bmc.configurators

import com.odim.simulator.dsl.DSL
import com.odim.simulator.tree.structure.ActionOem
import com.odim.simulator.tree.structure.ActionType.UPDATE_BIOS
import com.odim.simulator.tree.structure.ActionType.UPDATE_BIOS_BACKUP
import com.odim.simulator.tree.structure.ActionType.UPDATE_BMC
import com.odim.simulator.tree.structure.ActionType.UPDATE_ME
import com.odim.simulator.tree.structure.ActionType.UPDATE_SDR
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.Resource

class UpdateServiceConfigurator private constructor() {
    companion object Factory {
        fun configureUpdateService(updateService: Resource, baseSoftwareInventories: List<Resource>) {

            val softwareInventoryBMC = configureServiceInventory(baseSoftwareInventories[0], "BMC") {
                "Oem" to {
                    "@odata.type" to "#CustomCompany.Oem.SoftwareInventory"
                    "UpdateState" to "Idle"
                    "Progress" to 0
                }
            }

            val softwareInventoryBIOS = configureServiceInventory(baseSoftwareInventories[1], "BIOS") {
                "Oem" to {
                    "@odata.type" to "#CustomCompany.Oem.SoftwareInventory"
                }
            }

            val softwareInventoryME = configureServiceInventory(baseSoftwareInventories[2], "ME") {
                "Oem" to {
                    "@odata.type" to "#CustomCompany.Oem.SoftwareInventory"
                }
            }

            val softwareInventorySDR = configureServiceInventory(baseSoftwareInventories[3], "SDR") {
                "Oem" to {
                    "@odata.type" to "#CustomCompany.Oem.SoftwareInventory"
                    "UpdateState" to "SDR is uploaded"
                }

            }

            updateService(softwareInventoryBMC, softwareInventoryBIOS, softwareInventoryME, softwareInventorySDR)

            val updateBiosBackup = object : ActionOem(UPDATE_BIOS_BACKUP) {
                override fun toLink(): String {
                    return "${this.parent!!.toLink()}/Actions/Oem/Intel.Oem.${UPDATE_BIOS.actionName}/?imageflag=Backup"
                }
            }
            softwareInventoryBMC.traverse<Actions>("Actions").addAction(ActionOem(UPDATE_BMC), softwareInventoryBMC)
            softwareInventoryBIOS.traverse<Actions>("Actions").addAction(ActionOem(UPDATE_BIOS), softwareInventoryBIOS)
            softwareInventoryBIOS.traverse<Actions>("Actions").addAction(updateBiosBackup, softwareInventoryBIOS)
            softwareInventoryME.traverse<Actions>("Actions").addAction(ActionOem(UPDATE_ME), softwareInventoryME)
            softwareInventorySDR.traverse<Actions>("Actions").addAction(ActionOem(UPDATE_SDR), softwareInventorySDR)
        }

        private fun configureServiceInventory(softwareInventory: Resource, id: String, additionalParameters: DSL.() -> Unit): Resource {
            return softwareInventory {
                "Id" to id
                "Version" to "1.74.ee39402a"
            }.invoke(additionalParameters)
        }
    }
}
