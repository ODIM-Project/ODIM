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

package component

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.odimclient.RedfishHttpClient
import com.odim.simulator.SimulatorConfig
import com.odim.simulator.SimulatorConfig.Config.loadTestConfigFile
import com.odim.simulator.dsl.DSL
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpStatusCode
import com.odim.simulator.http.Response
import com.odim.simulator.http.ServerProvider
import com.odim.utils.getString
import com.odim.utils.retry
import org.assertj.core.api.Assertions
import org.springframework.web.client.ResourceAccessException
import org.testng.annotations.AfterClass
import org.testng.annotations.BeforeClass
import org.testng.annotations.Test
import java.time.Duration.ofMinutes

@Test
abstract class ComponentSimulatorTest {
    protected val serverIp = "127.0.0.1"
    protected val serverPort = 12345
    protected lateinit var httpClient : RedfishHttpClient
    protected lateinit var provider : ServerProvider

    protected val login = "admin"
    protected val password = "admin"

    @BeforeClass
    fun setConfigWithoutTLS() {
        loadTestConfigFile("/simulator-config.json")
        httpClient = RedfishHttpClient("http://$serverIp:$serverPort")
        provider = ServerProvider()
    }

    @AfterClass
    fun afterClazz() {
        provider.stopAll()
    }

    protected fun getSessionToken(login: String, password: String): String? {
        val response = retry(ofMinutes(1)) {
            notYetIfThrows<Response, ResourceAccessException> {
                val response: Response = httpClient.post("SessionService/Sessions", makeJson {
                    "UserName" to login
                    "Password" to password
                })
                if (response.code.is2XX) success(response)
                else notYet("The session was not created")
            }
        }

        return response.headers["X-Auth-Token"]?.first()
    }

    protected fun checkResourcePropertiesSet(resourceLocation: String, dsl: DSL.() -> Unit) {
        val resource = httpClient.get(resourceLocation)
        Assertions.assertThat(resource.code).isEqualTo(HttpStatusCode.OK)
        val propertiesMap = DSL().run { dsl(); toMap() }
        checkPropertiesSet("", resource.json, propertiesMap)
    }

    private fun checkPropertiesSet(currentPointer: String, resourceJson: ObjectNode, propertiesMap: Map<*, *>) {
        propertiesMap.forEach { (key, value) ->
            val pointer = "$currentPointer/$key"
            when (value) {
                is Map<*, *> -> {
                    checkPropertiesSet(pointer, resourceJson, value)
                }
                is List<*> -> {
                    value.forEachIndexed { index, element ->
                        when (element) {
                            is Map<*, *> -> checkPropertiesSet("$pointer/$index", resourceJson, element)
                            else -> checkField(resourceJson, "$pointer/$index", element)
                        }
                    }
                }
                else -> {
                    checkField(resourceJson, pointer, value)
                }
            }
        }
    }

    private fun checkField(resourceJson: ObjectNode, pointer: String, value: Any?) {
        Assertions.assertThat(resourceJson.getString(pointer))
                .`as`("Incorrect `$pointer` property value")
                .isEqualTo(value.toString())
    }
}
