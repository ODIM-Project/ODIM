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

import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0
import com.odim.simulator.tree.templates.redfish.EthernetInterfaceTemplate

@Template(ETHERNET_INTERFACE)
open class EthernetInterfaceBmcTemplate : EthernetInterfaceTemplate() {
    init {
        version(BMC_1_0, V1_6_2)
    }
}
