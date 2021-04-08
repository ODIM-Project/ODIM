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

package com.odim.odimclient

import com.odim.simulator.http.HttpClient

class RedfishHttpClient(val address: String, overrides: HttpClient.() -> Unit = {}) : HttpClient(overrides) {
    private val redfishPrefix = "/redfish/v1"
    private val shortRedfishPrefix = "redfish/v1"
    private val taskPrefix = "/taskmon"
    private val shortTaskPrefix = "taskmon"

    fun serviceRoot() = get(redfishPrefix).json

    override fun resolvePath(path: String) = address + when {
        path.isEmpty() -> redfishPrefix
        path.startsWith(redfishPrefix) -> path
        path.startsWith(taskPrefix) -> path
        path.startsWith(shortRedfishPrefix) -> "/$path"
        path.startsWith(shortTaskPrefix) -> "/$path"
        path.startsWith("/") -> "$redfishPrefix$path"
        else -> "$redfishPrefix/$path"
    }
}
