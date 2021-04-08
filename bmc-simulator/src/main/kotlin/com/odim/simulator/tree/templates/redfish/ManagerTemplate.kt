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
import com.odim.simulator.tree.RedfishVersion.V1_0_10
import com.odim.simulator.tree.RedfishVersion.V1_0_11
import com.odim.simulator.tree.RedfishVersion.V1_0_12
import com.odim.simulator.tree.RedfishVersion.V1_0_13
import com.odim.simulator.tree.RedfishVersion.V1_0_2
import com.odim.simulator.tree.RedfishVersion.V1_0_3
import com.odim.simulator.tree.RedfishVersion.V1_0_4
import com.odim.simulator.tree.RedfishVersion.V1_0_5
import com.odim.simulator.tree.RedfishVersion.V1_0_6
import com.odim.simulator.tree.RedfishVersion.V1_0_7
import com.odim.simulator.tree.RedfishVersion.V1_0_8
import com.odim.simulator.tree.RedfishVersion.V1_0_9
import com.odim.simulator.tree.RedfishVersion.V1_10_0
import com.odim.simulator.tree.RedfishVersion.V1_1_0
import com.odim.simulator.tree.RedfishVersion.V1_1_1
import com.odim.simulator.tree.RedfishVersion.V1_1_10
import com.odim.simulator.tree.RedfishVersion.V1_1_11
import com.odim.simulator.tree.RedfishVersion.V1_1_2
import com.odim.simulator.tree.RedfishVersion.V1_1_3
import com.odim.simulator.tree.RedfishVersion.V1_1_4
import com.odim.simulator.tree.RedfishVersion.V1_1_5
import com.odim.simulator.tree.RedfishVersion.V1_1_6
import com.odim.simulator.tree.RedfishVersion.V1_1_7
import com.odim.simulator.tree.RedfishVersion.V1_1_8
import com.odim.simulator.tree.RedfishVersion.V1_1_9
import com.odim.simulator.tree.RedfishVersion.V1_2_0
import com.odim.simulator.tree.RedfishVersion.V1_2_1
import com.odim.simulator.tree.RedfishVersion.V1_2_10
import com.odim.simulator.tree.RedfishVersion.V1_2_11
import com.odim.simulator.tree.RedfishVersion.V1_2_2
import com.odim.simulator.tree.RedfishVersion.V1_2_3
import com.odim.simulator.tree.RedfishVersion.V1_2_4
import com.odim.simulator.tree.RedfishVersion.V1_2_5
import com.odim.simulator.tree.RedfishVersion.V1_2_6
import com.odim.simulator.tree.RedfishVersion.V1_2_7
import com.odim.simulator.tree.RedfishVersion.V1_2_8
import com.odim.simulator.tree.RedfishVersion.V1_2_9
import com.odim.simulator.tree.RedfishVersion.V1_3_0
import com.odim.simulator.tree.RedfishVersion.V1_3_1
import com.odim.simulator.tree.RedfishVersion.V1_3_10
import com.odim.simulator.tree.RedfishVersion.V1_3_2
import com.odim.simulator.tree.RedfishVersion.V1_3_3
import com.odim.simulator.tree.RedfishVersion.V1_3_4
import com.odim.simulator.tree.RedfishVersion.V1_3_5
import com.odim.simulator.tree.RedfishVersion.V1_3_6
import com.odim.simulator.tree.RedfishVersion.V1_3_7
import com.odim.simulator.tree.RedfishVersion.V1_3_8
import com.odim.simulator.tree.RedfishVersion.V1_3_9
import com.odim.simulator.tree.RedfishVersion.V1_4_0
import com.odim.simulator.tree.RedfishVersion.V1_4_1
import com.odim.simulator.tree.RedfishVersion.V1_4_2
import com.odim.simulator.tree.RedfishVersion.V1_4_3
import com.odim.simulator.tree.RedfishVersion.V1_4_4
import com.odim.simulator.tree.RedfishVersion.V1_4_5
import com.odim.simulator.tree.RedfishVersion.V1_4_6
import com.odim.simulator.tree.RedfishVersion.V1_4_7
import com.odim.simulator.tree.RedfishVersion.V1_5_0
import com.odim.simulator.tree.RedfishVersion.V1_5_1
import com.odim.simulator.tree.RedfishVersion.V1_5_2
import com.odim.simulator.tree.RedfishVersion.V1_5_3
import com.odim.simulator.tree.RedfishVersion.V1_5_4
import com.odim.simulator.tree.RedfishVersion.V1_5_5
import com.odim.simulator.tree.RedfishVersion.V1_5_6
import com.odim.simulator.tree.RedfishVersion.V1_6_0
import com.odim.simulator.tree.RedfishVersion.V1_6_1
import com.odim.simulator.tree.RedfishVersion.V1_6_2
import com.odim.simulator.tree.RedfishVersion.V1_6_3
import com.odim.simulator.tree.RedfishVersion.V1_7_0
import com.odim.simulator.tree.RedfishVersion.V1_7_1
import com.odim.simulator.tree.RedfishVersion.V1_7_2
import com.odim.simulator.tree.RedfishVersion.V1_7_3
import com.odim.simulator.tree.RedfishVersion.V1_8_0
import com.odim.simulator.tree.RedfishVersion.V1_8_1
import com.odim.simulator.tree.RedfishVersion.V1_8_2
import com.odim.simulator.tree.RedfishVersion.V1_9_0
import com.odim.simulator.tree.RedfishVersion.V1_9_1
import com.odim.simulator.tree.ResourceTemplate
import com.odim.simulator.tree.Template
import com.odim.simulator.tree.structure.Action
import com.odim.simulator.tree.structure.ActionType.FORCE_FAILOVER
import com.odim.simulator.tree.structure.ActionType.MODIFY_REDUNDANCY_SET
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.ActionType.RESET_TO_DEFAULTS
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.EmbeddedObjectType.COMMAND_SHELL
import com.odim.simulator.tree.structure.EmbeddedObjectType.GRAPHICAL_CONSOLE
import com.odim.simulator.tree.structure.EmbeddedObjectType.SERIAL_CONSOLE
import com.odim.simulator.tree.structure.EmbeddedObjectType.STATUS
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.ETHERNET_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.HOST_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.LOG_SERVICES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.SERIAL_INTERFACES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.VIRTUAL_MEDIAS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.ACCOUNT_SERVICE
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.MANAGER_NETWORK_PROTOCOL
import com.odim.simulator.tree.structure.ResourceType.REDUNDANCY
import com.odim.simulator.tree.structure.ResourceType.SOFTWARE_INVENTORY
import com.odim.simulator.tree.structure.ResourceType.SWITCH
import com.odim.simulator.tree.structure.SingletonResource

/**
 * This is generated class. Please don't edit it's contents.
 */
@Template(MANAGER)
open class ManagerTemplate : ResourceTemplate() {
    init {
        version(V1_0_0, resourceObject(
                "Oem" to embeddedObject(),
                "Id" to 0,
                "Description" to "Manager Description",
                "Name" to "Manager",
                "ManagerType" to "ManagementController",
                "EthernetInterfaces" to ResourceCollection(ETHERNET_INTERFACES_COLLECTION),
                "SerialInterfaces" to ResourceCollection(SERIAL_INTERFACES_COLLECTION),
                "NetworkProtocol" to SingletonResource(MANAGER_NETWORK_PROTOCOL),
                "LogServices" to ResourceCollection(LOG_SERVICES_COLLECTION),
                "VirtualMedia" to ResourceCollection(VIRTUAL_MEDIAS_COLLECTION),
                "ServiceEntryPointUUID" to null,
                "UUID" to null,
                "Model" to null,
                "DateTime" to null,
                "DateTimeLocalOffset" to null,
                "FirmwareVersion" to null,
                "SerialConsole" to embeddedObject(SERIAL_CONSOLE),
                "CommandShell" to embeddedObject(COMMAND_SHELL),
                "GraphicalConsole" to embeddedObject(GRAPHICAL_CONSOLE),
                "Status" to embeddedObject(STATUS),
                "Redundancy" to EmbeddedResourceArray(REDUNDANCY),
                "Links" to embeddedObject(
                        "Oem" to embeddedObject(),
                        "ManagerForServers" to LinkableResourceArray(COMPUTER_SYSTEM),
                        "ManagerForChassis" to LinkableResourceArray(CHASSIS)
                ),
                "Actions" to Actions(
                        Action(RESET, "ResetType", mutableListOf(
                                "On",
                                "ForceOff",
                                "GracefulShutdown",
                                "GracefulRestart",
                                "ForceRestart",
                                "Nmi",
                                "ForceOn",
                                "PushPowerButton",
                                "PowerCycle"
                        )),
                        Action(FORCE_FAILOVER, "NewManager", mutableListOf()),
                        Action(MODIFY_REDUNDANCY_SET)
                )
        ))
        version(V1_0_1, V1_0_0)
        version(V1_0_2, V1_0_1)
        version(V1_0_3, V1_0_2)
        version(V1_0_4, V1_0_3)
        version(V1_0_5, V1_0_4)
        version(V1_0_6, V1_0_5)
        version(V1_0_7, V1_0_6)
        version(V1_0_8, V1_0_7)
        version(V1_0_9, V1_0_8)
        version(V1_0_10, V1_0_9)
        version(V1_0_11, V1_0_10)
        version(V1_0_12, V1_0_11)
        version(V1_0_13, V1_0_12)
        version(V1_1_0, V1_0_2, embeddedObject(
                "Links" to embeddedObject(
                        "ManagerInChassis" to LinkableResource(CHASSIS)
                )
        ))
        version(V1_1_1, V1_1_0)
        version(V1_1_2, V1_1_1)
        version(V1_1_3, V1_1_2)
        version(V1_1_4, V1_1_3)
        version(V1_1_5, V1_1_4)
        version(V1_1_6, V1_1_5)
        version(V1_1_7, V1_1_6)
        version(V1_1_8, V1_1_7)
        version(V1_1_9, V1_1_8)
        version(V1_1_10, V1_1_9)
        version(V1_1_11, V1_1_10)
        version(V1_2_0, V1_1_0, resourceObject(
                "PowerState" to null
        ))
        version(V1_2_1, V1_2_0)
        version(V1_2_2, V1_2_1)
        version(V1_2_3, V1_2_2)
        version(V1_2_4, V1_2_3)
        version(V1_2_5, V1_2_4)
        version(V1_2_6, V1_2_5)
        version(V1_2_7, V1_2_6)
        version(V1_2_8, V1_2_7)
        version(V1_2_9, V1_2_8)
        version(V1_2_10, V1_2_9)
        version(V1_2_11, V1_2_10)
        version(V1_3_0, V1_2_1, resourceObject(
                "HostInterfaces" to ResourceCollection(HOST_INTERFACES_COLLECTION)
        ))
        version(V1_3_1, V1_3_0)
        version(V1_3_2, V1_3_1)
        version(V1_3_3, V1_3_2)
        version(V1_3_4, V1_3_3)
        version(V1_3_5, V1_3_4)
        version(V1_3_6, V1_3_5)
        version(V1_3_7, V1_3_6)
        version(V1_3_8, V1_3_7)
        version(V1_3_9, V1_3_8)
        version(V1_3_10, V1_3_9)
        version(V1_4_0, V1_3_3, resourceObject(
                "AutoDSTEnabled" to false,
                "Links" to embeddedObject(
                        "ManagerForSwitches" to LinkableResourceArray(SWITCH)
                )
        ))
        version(V1_4_1, V1_4_0)
        version(V1_4_2, V1_4_1)
        version(V1_4_3, V1_4_2)
        version(V1_4_4, V1_4_3)
        version(V1_4_5, V1_4_4)
        version(V1_4_6, V1_4_5)
        version(V1_4_7, V1_4_6)
        version(V1_5_0, V1_4_1, resourceObject(
                "RemoteRedfishServiceUri" to null,
                "RemoteAccountService" to SingletonResource(ACCOUNT_SERVICE)
        ))
        version(V1_5_1, V1_5_0)
        version(V1_5_2, V1_5_1)
        version(V1_5_3, V1_5_2)
        version(V1_5_4, V1_5_3)
        version(V1_5_5, V1_5_4)
        version(V1_5_6, V1_5_5)
        version(V1_6_0, V1_5_3, embeddedObject(
                "Links" to embeddedObject(
                        "ActiveSoftwareImage" to LinkableResource(SOFTWARE_INVENTORY),
                        "SoftwareImages" to LinkableResourceArray(SOFTWARE_INVENTORY)
                )
        ))
        version(V1_6_1, V1_6_0)
        version(V1_6_2, V1_6_1)
        version(V1_6_3, V1_6_2)
        version(V1_7_0, V1_6_0, resourceObject(
                "Manufacturer" to null,
                "SerialNumber" to null,
                "PartNumber" to null
        ))
        version(V1_7_1, V1_7_0)
        version(V1_7_2, V1_7_1)
        version(V1_7_3, V1_7_2)
        version(V1_8_0, V1_7_1, embeddedObject(
                "Actions" to Actions(
                        Action(RESET_TO_DEFAULTS, "ResetType", mutableListOf())
                )
        ))
        version(V1_8_1, V1_8_0)
        version(V1_8_2, V1_8_1)
        version(V1_9_0, V1_8_1, resourceObject(
                "LastResetTime" to "2017-04-14T06:35:05Z",
                "Links" to embeddedObject(
                        "ManagedBy" to LinkableResourceArray(MANAGER),
                        "ManagerForManagers" to LinkableResourceArray(MANAGER)
                )
        ))
        version(V1_9_1, V1_9_0)
        version(V1_10_0, V1_9_1, resourceObject(
                "TimeZoneName" to ""
        ))
    }
}
