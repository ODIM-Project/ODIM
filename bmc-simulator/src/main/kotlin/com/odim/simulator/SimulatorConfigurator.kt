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

package com.odim.simulator

import com.fasterxml.jackson.databind.node.ArrayNode
import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.simulator.CoreConfig.EXTERNAL_KEYSTORE_LOCATION
import com.odim.simulator.CoreConfig.EXTERNAL_KEYSTORE_PASSWORD
import com.odim.simulator.CoreConfig.EXTERNAL_TRUSTSTORE_LOCATION
import com.odim.simulator.CoreConfig.EXTERNAL_TRUSTSTORE_PASSWORD
import com.odim.simulator.CoreConfig.HTTP_CLIENT_BASIC_CREDENTIALS
import com.odim.simulator.CoreConfig.HTTP_CLIENT_WITH_BASIC
import com.odim.simulator.CoreConfig.JETTY_SERVER_MAX_THREADS
import com.odim.simulator.CoreConfig.REGISTER_ODIM_URL
import com.odim.simulator.CoreConfig.RESOURCES_CONFIG
import com.odim.simulator.CoreConfig.SERVER_BASIC_CREDENTIALS
import com.odim.simulator.CoreConfig.SERVER_CERTIFICATE_VERIFICATION_ENABLED
import com.odim.simulator.CoreConfig.SERVE_IP
import com.odim.simulator.CoreConfig.SERVE_PORT
import com.odim.simulator.CoreConfig.USE_CLIENT_TLS
import com.odim.simulator.CoreConfig.USE_ETAGS
import com.odim.simulator.CoreConfig.USE_SERVER_MTLS
import com.odim.simulator.CoreConfig.USE_SERVER_TLS
import com.odim.utils.JsonMapper.emptyJson
import com.odim.utils.JsonMapper.jsonMapper
import com.odim.utils.getArrayOrNull
import com.odim.utils.getBooleanOrNull
import com.odim.utils.getNumberOrNull
import com.odim.utils.getObjectOrNull
import com.odim.utils.getStringOrNull
import com.odim.utils.getStringOrNullPrefEnv
import com.odim.utils.merge
import com.odim.utils.prettyPrint
import org.slf4j.LoggerFactory.getLogger
import java.io.InputStream
import java.nio.file.Files.exists
import java.nio.file.Path
import java.nio.file.Paths.get

enum class CoreConfig(val path: String) {
    SERVE_IP("binding/ip"),
    SERVE_PORT("binding/port"),
    REGISTER_ODIM_URL("binding/odimUrl"),
    JETTY_SERVER_MAX_THREADS("httpServerMaxThreads"),
    USE_SERVER_TLS("security/server/useTLS"),
    EXTERNAL_KEYSTORE_LOCATION("security/server/externalKeyStoreLocation"),
    EXTERNAL_KEYSTORE_PASSWORD("security/server/externalKeyStorePassword"),
    USE_SERVER_MTLS("security/server/useServerMTLS"),
    EXTERNAL_TRUSTSTORE_LOCATION("security/trustStore/externalTrustStoreLocation"),
    EXTERNAL_TRUSTSTORE_PASSWORD("security/trustStore/externalTrustStorePassword"),
    HTTP_CLIENT_WITH_BASIC("security/httpClient/withBasic"),
    HTTP_CLIENT_BASIC_CREDENTIALS("security/httpClient/basicCredentials"),
    SERVER_BASIC_CREDENTIALS("security/server/basicCredentials"),
    USE_CLIENT_TLS("security/httpClient/useTLS"),
    SERVER_CERTIFICATE_VERIFICATION_ENABLED("security/httpClient/serverCertificateVerificationEnabled"),
    RESOURCES_CONFIG("resourcesConfig"),
    USE_ETAGS("cache/etags")
}

class SimulatorConfig private constructor() {
    companion object Config {
        const val defaultExternalConfigFileName = "simulator-config.json"
        private val logger = getLogger(this::class.java)
        private var defaultConfigFile = this::class.java.getResourceAsStream("/simulator-config.json")
        var config: ObjectNode = loadDefaultConfig().merge(tryGetPropertiesFromExternalFile(defaultExternalConfigFileName))

        init {
            System.setProperty("jdk.tls.namedGroups", "secp521r1")
            System.setProperty("jdk.tls.rejectClientInitiatedRenegotiation", "true")
        }

        fun appendValuesFromExternalConfig(externalConfigFilePathProperty: String?) {
            externalConfigFilePathProperty?.let { config.merge(tryGetPropertiesFromExternalFile(it)) }
        }

        /**
         * Use this method only for component tests purpose
         */
        fun loadTestConfigFile(path: String) {
            defaultConfigFile = this::class.java.getResourceAsStream(path)
            val testConfigFile = get(defaultExternalConfigFileName)
            config = loadDefaultConfig().merge(loadTestConfig(testConfigFile))
        }

        /**
         * Use this method only for component tests purpose
         */
        private fun loadTestConfig(externalConfigFile: Path): ObjectNode = if (exists(externalConfigFile)) {
            load(externalConfigFile.toFile().inputStream())
        } else emptyJson

        /**
         * Method to get properties from core configuration. Values should exist
         */
        @Suppress("IMPLICIT_CAST_TO_ANY", "ComplexMethod")
        inline fun <reified T : Any> getConfigProperty(property: CoreConfig): T {
            return (when (property) {
                SERVE_IP -> config.getStringOrNull(SERVE_IP.path)
                SERVE_PORT -> config.getNumberOrNull(SERVE_PORT.path)?.toInt()
                JETTY_SERVER_MAX_THREADS -> config.getNumberOrNull(JETTY_SERVER_MAX_THREADS.path)?.toInt()
                REGISTER_ODIM_URL -> config.getStringOrNull(REGISTER_ODIM_URL.path)
                USE_SERVER_TLS -> config.getBooleanOrNull(USE_SERVER_TLS.path)
                EXTERNAL_KEYSTORE_LOCATION -> config.getStringOrNullPrefEnv(EXTERNAL_KEYSTORE_LOCATION.path)
                EXTERNAL_KEYSTORE_PASSWORD -> config.getStringOrNullPrefEnv(EXTERNAL_KEYSTORE_PASSWORD.path)
                USE_SERVER_MTLS -> config.getBooleanOrNull(USE_SERVER_MTLS.path)
                EXTERNAL_TRUSTSTORE_LOCATION -> config.getStringOrNullPrefEnv(EXTERNAL_TRUSTSTORE_LOCATION.path)
                EXTERNAL_TRUSTSTORE_PASSWORD -> config.getStringOrNullPrefEnv(EXTERNAL_TRUSTSTORE_PASSWORD.path)
                HTTP_CLIENT_WITH_BASIC -> config.getBooleanOrNull(HTTP_CLIENT_WITH_BASIC.path)
                HTTP_CLIENT_BASIC_CREDENTIALS -> config.getStringOrNull(HTTP_CLIENT_BASIC_CREDENTIALS.path)
                SERVER_BASIC_CREDENTIALS -> config.getStringOrNull(SERVER_BASIC_CREDENTIALS.path)
                USE_CLIENT_TLS -> config.getBooleanOrNull(USE_CLIENT_TLS.path)
                SERVER_CERTIFICATE_VERIFICATION_ENABLED -> config.getBooleanOrNull(SERVER_CERTIFICATE_VERIFICATION_ENABLED.path)
                USE_ETAGS -> config.getBooleanOrNull(USE_ETAGS.path)
                RESOURCES_CONFIG -> config.getObjectOrNull(RESOURCES_CONFIG.path)
            } ?: throw ConfigurationException("Property ${property.path} not found!")) as T
        }

        /**
         * Methods to get properties from external configuration. Provide defaults.
         */
        @Suppress("IMPLICIT_CAST_TO_ANY")
        inline fun <reified T : Any> getConfigProperty(key: String, default: T): T {
            val value = when (default) {
                is Boolean -> config.getBooleanOrNull(key) ?: default
                is String -> config.getStringOrNull(key) ?: default
                is Int -> config.getStringOrNull(key)?.toInt() ?: default
                is Long -> config.getNumberOrNull(key)?.toLong() ?: default
                is ObjectNode -> config.getObjectOrNull(key) ?: default
                is ArrayNode -> config.getArrayOrNull(key) ?: default
                else -> default
            }
            return value as T
        }

        private fun loadDefaultConfig() = load(defaultConfigFile)

        private fun tryGetPropertiesFromExternalFile(externalConfigFileName: String): ObjectNode {
            val externalConfigFile = get(externalConfigFileName)
            return if (exists(externalConfigFile)) {
                val externalProperties = load(externalConfigFile.toFile().inputStream())
                logger.debug("Loaded external config: ${externalProperties.prettyPrint()}")
                externalProperties
            } else emptyJson
        }

        private fun load(inputStream: InputStream) = inputStream.bufferedReader().use {
            jsonMapper.readTree(it) as ObjectNode
        }
    }
}

class ConfigurationException : RuntimeException {
    constructor(message: String) : super(message)
    constructor(message: String, e: Exception) : super(message, e)
}
