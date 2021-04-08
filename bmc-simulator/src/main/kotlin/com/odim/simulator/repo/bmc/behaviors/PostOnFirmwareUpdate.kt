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

package com.odim.simulator.repo.bmc.behaviors

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorDataStore.SharedInformationType.FIRMWARE_UPDATE_MESSAGES
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.behaviors.TreeJsonRenderer
import com.odim.simulator.dsl.merger.MergeException
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.accepted
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.internalServerError
import com.odim.simulator.repo.bmc.FirmwareUpdateMessages.UPDATE_COMPLETED_SUFFIX
import com.odim.simulator.repo.bmc.FirmwareUpdateMessages.UPDATE_STARTED_SUFFIX
import com.odim.simulator.repo.bmc.createSelEntry
import com.odim.simulator.repo.bmc.utils.COULDNT_PARSE
import com.odim.simulator.repo.bmc.utils.extractSDRLine
import com.odim.simulator.repo.bmc.utils.filenameWithoutExtension
import com.odim.simulator.repo.bmc.utils.getFileMetadata
import com.odim.simulator.repo.bmc.utils.parseBIOSVersionAndBuildId
import com.odim.simulator.repo.bmc.utils.parseBMCVersionAndBuildId
import com.odim.simulator.repo.bmc.utils.parseSDRVersion
import com.odim.simulator.repo.bmc.utils.uploadFile
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.ActionOem
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import java.io.ByteArrayOutputStream
import java.io.InputStream

class PostOnFirmwareUpdate(private val logService: Resource, val manager: Resource, private val deferUpdateAfterReboot: Boolean = false) : Behavior {
    @Suppress("ReturnCount")
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        request.request ?: return terminal(internalServerError())

        val resource = (item as ActionOem).parent as Resource
        val resourceId = resource.getId()

        val task = createTask(tree, resourceId)
        try {
            Merger.merge(tree, task, request)
        } catch (e: MergeException) {
            return terminal(badRequest("Bad POST request: ${e.message}"))
        }

        createSelEntry(logService, tree, "Firmware Update - $resourceId image file  - is uploaded")

        val fileInfo = getFileMetadata(request.request)
        val uploadedFileBytes = getUploadedFileBytes(fileInfo.content)

        uploadFile(fileInfo.filename, uploadedFileBytes)

        val versionAndBuildIdSource = when (resourceId) {
            "SDR" -> extractSDRLine(uploadedFileBytes)
            else -> filenameWithoutExtension(fileInfo)
        }

        val (version, buildId) = parseVersionAndBuildId(resourceId, versionAndBuildIdSource)

        val messagePrefix = createMessagePrefix(resourceId, version, buildId)

        updateResourceVersion(resource, version)
        if (resourceId == "BMC") {
            manager {
                "FirmwareVersion" to version
            }
        }

        updateSel(messagePrefix, dataStore, tree)

        return terminal(accepted("/redfish/v1/TaskService/Tasks/${task.getId()}/TaskMonitor",
                TreeJsonRenderer().toJson(task)))
    }

    private fun getUploadedFileBytes(inputStream: InputStream): ByteArray? {
        val outputStream = ByteArrayOutputStream()
        inputStream.use { input ->
            outputStream.use { output ->
                input.copyTo(output)
            }
        }
        return outputStream.toByteArray()
    }

    private fun updateSel(messagePrefix: String, dataStore: BehaviorDataStore, tree: ResourceTree) {
        val updateStartedMessage = "$messagePrefix$UPDATE_STARTED_SUFFIX"
        val updateCompletedMessage = "$messagePrefix$UPDATE_COMPLETED_SUFFIX"

        when (deferUpdateAfterReboot) {
            true -> dataStore.insert(FIRMWARE_UPDATE_MESSAGES, listOf(updateStartedMessage, updateCompletedMessage))
            false -> {
                createSelEntry(logService, tree, updateStartedMessage)
                createSelEntry(logService, tree, updateCompletedMessage)
            }
        }
    }

    private fun createMessagePrefix(resourceId: String, version: String, buildId: String): String =
            when (resourceId) {
                "BMC" -> "Target: BMC Version:$version,Build ID:$buildId"
                "BIOS" -> "Target: BIOS Version:$version,Release Number:$buildId"
                "ME" -> "Target: ME Version:$version,Build Number:$buildId"
                "SDR" -> "Target: SDR Version:$version"
                else -> "Unknown update type"
            }

    private fun parseVersionAndBuildId(resourceId: String, versionSource: String): Pair<String, String> =
            when (resourceId) {
                "BMC" -> parseBMCVersionAndBuildId(versionSource)
                "BIOS" -> parseBIOSVersionAndBuildId(versionSource)
                "ME" -> parseBIOSVersionAndBuildId(versionSource)
                "SDR" -> parseSDRVersion(versionSource)
                else -> Pair(COULDNT_PARSE, COULDNT_PARSE)
            }

    private fun updateResourceVersion(resource: Resource, version: String) {
        resource {
            "Version" to version
        }
    }

    private fun createTask(tree: ResourceTree, resourceId: String): Resource {
        val tasks = tree.root.traverse<Resource>("TaskService").traverse<ResourceCollection>("Tasks")
        val task = tree.create(tasks.type.of())
        tasks.add(task)

        return task {
            "Name" to "$resourceId update"
        }
    }
}
