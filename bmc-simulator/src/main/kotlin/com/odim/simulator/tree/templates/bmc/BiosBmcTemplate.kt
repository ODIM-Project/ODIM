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

import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.BIOS
import com.odim.simulator.tree.structure.ResourceType.BIOS_SETTINGS
import com.odim.simulator.tree.structure.SingletonResource
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0
import com.odim.simulator.tree.templates.redfish.BiosTemplate

@Template(BIOS)
open class BiosBmcTemplate : BiosTemplate() {
    init {
        version(BMC_1_0, V1_1_1, resourceObject(
                "@Redfish.Settings" to resourceObject(
                        "@odata.type" to "#Settings.v1_1_0.Settings",
                        "SettingsObject" to SingletonResource(BIOS_SETTINGS)
                )
        ))
    }
}
