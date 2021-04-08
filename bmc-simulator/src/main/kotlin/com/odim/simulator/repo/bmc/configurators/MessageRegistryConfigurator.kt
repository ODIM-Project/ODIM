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
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection

@SuppressWarnings("LongMethod")
class MessageRegistryConfigurator private constructor() {
    companion object Factory {
        private val location = makeJson {
            "Location" to array[
                    {
                        "Language" to "en"
                        "PublicationUri" to "http://www.dmtf.org/sites/default/files/standards/documents/DSP8011_1.0.0a.json"
                        "ArchiveUri" to ""
                        "Uri" to ""
                        "ArchiveFile" to ""
                    }
            ]
        }

        fun configureMessageRegistry(messageRegistryService: ResourceCollection, messageRegistries: List<Resource>) {
            val baseMessages = createRegistry(
                    messageRegistries[0],
                    "Base Message",
                    "BaseMessages",
                    "Base.1.5.0") {
                location
            }

            val eventingMessages = createRegistry(
                    messageRegistries[1],
                    "Eventing",
                    "EventingMessages",
                    "Alert.1.0.0") {
                location
            }

            val commonMessages = createRegistry(
                    messageRegistries[2],
                    "Common Message",
                    "CommonMessages",
                    "Common.1.0.0") {
                location
            }

            val statusChangeMessages = createRegistry(
                    messageRegistries[3],
                    "StatusChange Message",
                    "StatusChangeMessages",
                    "StatusChange.1.0.0") {
                location
            }

            val biosAttributeRegistryMessages = createRegistry(
                    messageRegistries[4],
                    "Bios",
                    "BiosAttributeRegistry",
                    "BiosAttributeRegistry.1.0.0") {
                location
            }

            messageRegistryService
                    .add(baseMessages)
                    .add(eventingMessages)
                    .add(commonMessages)
                    .add(statusChangeMessages)
                    .add(biosAttributeRegistryMessages)
        }

        private fun createRegistry(registryResource: Resource, name: String, id: String, registry: String, override: DSL.() -> Unit): Resource {
            return registryResource {
                "Name" to "$name Registry File"
                "Description" to "$name Registry File locations"
                "Id" to id
                "Registry" to registry
                "Languages" to array["en"]
            }.invoke(override)
        }
    }
}
