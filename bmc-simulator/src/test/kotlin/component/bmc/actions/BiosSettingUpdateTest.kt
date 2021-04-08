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

package component.bmc.actions

import com.fasterxml.jackson.databind.node.ObjectNode
import com.nhaarman.mockitokotlin2.mock
import com.odim.simulator.behaviors.ResourceBehaviors
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpMethod.PATCH
import com.odim.simulator.http.HttpMethod.POST
import com.odim.simulator.http.HttpStatusCode.NO_CONTENT
import com.odim.simulator.http.HttpStatusCode.OK
import com.odim.simulator.http.Request
import com.odim.simulator.repo.bmc.behaviors.PatchOnBiosSettings
import com.odim.simulator.repo.bmc.behaviors.ResetOnSystem
import com.odim.simulator.tree.ResourceFactory
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.Actions
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceObject
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.templates.bmc.BmcVersion.Companion.BMC_VERSION_LATEST
import org.testng.Assert
import org.testng.Assert.assertEquals
import org.testng.annotations.BeforeMethod
import org.testng.annotations.Test
import kotlin.random.Random.Default.nextInt

class BiosSettingUpdateTest {
    private lateinit var behaviors: ResourceBehaviors
    private lateinit var tree: ResourceTree
    private lateinit var system: Resource
    private val bios
        get() = system.traverse<Resource>("Bios")
    private val biosSettings
        get() = bios.traverse<Resource>("@Redfish.Settings/SettingsObject")
    private val resetSystem
        get() = system.traverse<Actions>("Actions").getAction(RESET)

    @BeforeMethod
    fun setup() {
        tree = ResourceTree(ResourceFactory(listOf("redfish", "bmc"), BMC_VERSION_LATEST))
        system = tree.create(COMPUTER_SYSTEM)
        tree.root.append(system)
        behaviors = ResourceBehaviors()
        behaviors.replaceBehavior(biosSettings, PATCH, PatchOnBiosSettings())
        behaviors.appendActionBehavior(RESET, POST, ResetOnSystem(mock(), bios))
    }

    @Test
    fun `PATCH on Bios Settings should return OK status`() {
        val response = behaviors.createResponse(tree, biosSettings, Request(PATCH, createBiosAttributesJson().toString()))
        assertEquals(response.code, OK)
    }

    @Test
    fun `PATCH on Bios Settings should not modify BiosSettings or Bios resource`() {
        val json = createBiosAttributesJson()
        val response = behaviors.createResponse(tree, biosSettings, Request(PATCH, json.toString()))
        assertEquals(response.code, OK)

        verifyBiosSettings(json)
        verifyBiosAttributes(json, Assert::assertNotEquals)
    }

    @Test
    fun `PATCH on Bios Settings and then ForceRestart on System should apply new settings on Bios`() {
        val json = createBiosAttributesJson()
        val patchResponse = behaviors.createResponse(tree, biosSettings, Request(PATCH, json.toString()))
        assertEquals(patchResponse.code, OK)

        val resetResponse = behaviors.createResponse(tree, resetSystem, Request(POST, createResetJson("ForceRestart")))
        assertEquals(resetResponse.code, NO_CONTENT)
        verifyBiosAttributes(json, Assert::assertEquals)
    }

    @Test
    fun `New Bios settings should be applied only when System's PowerState is On after power cycle`() {
        val json = createBiosAttributesJson()
        val patchResponse = behaviors.createResponse(tree, biosSettings, Request(PATCH, json.toString()))
        assertEquals(patchResponse.code, OK)

        val forceOffResponse = behaviors.createResponse(tree, resetSystem, Request(POST, createResetJson("ForceOff")))
        assertEquals(forceOffResponse.code, NO_CONTENT)
        assertEquals(system.data["PowerState"], "Off")
        verifyBiosAttributes(json, Assert::assertNotEquals)

        val onResetResponse = behaviors.createResponse(tree, resetSystem, Request(POST, createResetJson("On")))
        assertEquals(onResetResponse.code, NO_CONTENT)
        assertEquals(system.data["PowerState"], "On")
        verifyBiosAttributes(json, Assert::assertEquals)
    }

    private fun verifyBiosSettings(json: ObjectNode) {
        json.fieldNames().iterator().forEach { key ->
            assertEquals(biosSettings.data[key], null)
        }
    }

    private fun verifyBiosAttributes(json: ObjectNode, assert: (Any?, Any?) -> Unit) {
        val biosAttributes = bios.data["Attributes"] as ResourceObject
        json.fieldNames().iterator().forEach { key ->
            assertEquals(biosAttributes.containsKey(key), true)
            biosAttributes[key]?.let {
                assert(it, json[key].toString().toInt())
            }
        }
    }

    private fun createBiosAttributesJson(): ObjectNode {
        return makeJson {
            "QuietBoot" to nextInt(10, 100)
            "IB1_PXE" to nextInt(10, 100)
            "MmiohBase" to nextInt(10, 100)
        }
    }

    private fun createResetJson(type: String): String {
        return makeJson {
            "ResetType" to type
        }.toString()
    }
}
