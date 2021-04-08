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
import com.odim.simulator.repo.bmc.configurators.AccountServiceConfigurator.Factory.Role.Administrator
import com.odim.simulator.repo.bmc.configurators.AccountServiceConfigurator.Factory.Role.Callback
import com.odim.simulator.repo.bmc.configurators.AccountServiceConfigurator.Factory.Role.NoAccess
import com.odim.simulator.repo.bmc.configurators.AccountServiceConfigurator.Factory.Role.Operator
import com.odim.simulator.repo.bmc.configurators.AccountServiceConfigurator.Factory.Role.User
import com.odim.simulator.tree.ResourceLinker.link
import com.odim.simulator.tree.structure.Resource

class AccountServiceConfigurator private constructor() {
    companion object Factory {
        enum class Role {
            NoAccess,
            Administrator,
            Operator,
            User,
            Callback
        }

        fun configureAccountService(accountService: Resource, roles: MutableList<Resource>, accounts: MutableList<Resource>) {
            accountService {
                "ServiceEnabled" to true
                "MinPasswordLength" to 6
                "MaxPasswordLength" to 20
            }

            val adminRole = configureRole(roles[0], Administrator, "Administrator User Role") {
                "AssignedPrivileges" to array[
                        "Login",
                        "ConfigureManager",
                        "ConfigureUsers",
                        "ConfigureSelf",
                        "ConfigureComponents"
                ]
            }

            val operatorRole = configureRole(roles[1], Operator, "Operator User Role") {
                "AssignedPrivileges" to array[
                        "Login",
                        "ConfigureComponents"
                ]
            }

            val userRole = configureRole(roles[2], User, "User Role") {
                "AssignedPrivileges" to array[
                        "Login"
                ]
            }

            val callbackRole = configureRole(roles[3], Callback, "Callback User Role")
            val noAccessRole = configureRole(roles[4], NoAccess, "NoAccess User Role")

            val noAccessAccount = createAccount(accounts[0], 1, "NoAccess account", "anonymous", NoAccess)
            val adminAccount = createAccount(accounts[1], 2, "Administrator account", "root", Administrator, true)

            accountService(adminRole, operatorRole, userRole, callbackRole, noAccessRole, noAccessAccount, adminAccount)
            link(noAccessAccount, noAccessRole)
            link(adminAccount, adminRole)
        }

        @Suppress("UnnecessaryApply")
        private fun configureRole(role: Resource, id: Role, desc: String, isPredefined: Boolean = true, privileges: (DSL.() -> Unit)? = null): Resource {
            return role {
                "Id" to id.toString()
                "Name" to "User Role"
                "Description" to desc
                "IsPredefined" to isPredefined
            }.apply { privileges?.let { invoke(it) } }
        }

        @Suppress("LongParameterList")
        private fun createAccount(account: Resource, id: Int, name: String, userName: String, roleId: Role, enabled: Boolean = false): Resource {
            return account {
                "Id" to id
                "Name" to name
                "Description" to "User Account"
                "Enabled" to enabled
                "Password" to null
                "UserName" to userName
                "RoleId" to roleId.toString()
            }
        }
    }
}
