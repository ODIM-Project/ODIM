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

package com.odim.simulator.tree.templates

import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.Resource.Companion.embeddedObject
import com.odim.simulator.tree.structure.Resource.Companion.resourceObject
import com.odim.simulator.tree.structure.ResourceType.ROLE

fun managerAccount() = resourceObject(
        "Id"                    to 0,
        "Name"                  to "Manager Account",
        "Description"           to "Manager Account description",
        "Password"              to "",
        "UserName"              to "",
        "RoleId"                to "",
        "Locked"                to false,
        "Enabled"               to true,
        "Links"                 to embeddedObject(
                "Role"                  to LinkableResource(ROLE),
                "Oem"                   to embeddedObject()
        ),
        "Oem"                   to embeddedObject()
)
