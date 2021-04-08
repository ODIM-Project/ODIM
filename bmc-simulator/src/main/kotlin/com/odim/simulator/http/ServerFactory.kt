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

package com.odim.simulator.http

import com.odim.simulator.CoreConfig.EXTERNAL_KEYSTORE_LOCATION
import com.odim.simulator.CoreConfig.EXTERNAL_KEYSTORE_PASSWORD
import com.odim.simulator.CoreConfig.EXTERNAL_TRUSTSTORE_LOCATION
import com.odim.simulator.CoreConfig.EXTERNAL_TRUSTSTORE_PASSWORD
import com.odim.simulator.CoreConfig.JETTY_SERVER_MAX_THREADS
import com.odim.simulator.CoreConfig.USE_SERVER_MTLS
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.http.HttpStatusCode.METHOD_NOT_ALLOWED
import com.odim.simulator.http.HttpStatusCode.NOT_FOUND
import com.odim.simulator.http.HttpStatusCode.UNAUTHORIZED
import com.odim.simulator.http.Response.Companion.errorJson
import io.javalin.Javalin
import io.javalin.websocket.WsHandler
import org.eclipse.jetty.server.Connector
import org.eclipse.jetty.server.Server
import org.eclipse.jetty.server.ServerConnector
import org.eclipse.jetty.util.ssl.SslContextFactory
import org.eclipse.jetty.util.thread.QueuedThreadPool
import org.slf4j.LoggerFactory.getLogger
import java.util.function.Consumer

object ServerFactory {
    private val logger = getLogger(this.javaClass)

    fun create(ip: String, port: Int, useTls: Boolean = false, wsConfig: Map<String, Consumer<WsHandler>> = emptyMap()) = Javalin.create { config ->
        config.defaultContentType = "application/json"
        config.showJavalinBanner = false
        config.server {
            val server = Server(QueuedThreadPool(getConfigProperty(JETTY_SERVER_MAX_THREADS.path, 20)))
            val connector = if (useTls) {
                ServerConnector(server, getSslContextFactory())
            } else {
                ServerConnector(server)
            }
            connector.host = ip
            connector.port = port
            server.connectors = arrayOf<Connector>(connector)
            server
        }
    }.also {
        it.server()!!.serverPort = port
        wsConfig.forEach { (endpoint, wsHandler) ->
            it.ws(endpoint, wsHandler)
        }
    }
            .error(UNAUTHORIZED.value) { it.result(errorJson("Bad credentials.")) }
            .error(NOT_FOUND.value) { it.result(errorJson("Resource not found.")) }
            .error(METHOD_NOT_ALLOWED.value) { it.result(errorJson("Method not allowed.")) }!!

    private fun getSslContextFactory(): SslContextFactory {
        val sslContextFactory = SslContextFactory.Server()

        sslContextFactory.setIncludeCipherSuites(
                "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384")

        val keyStoreLocation: String = getConfigProperty(EXTERNAL_KEYSTORE_LOCATION)
        logger.debug("Using external keystore: $keyStoreLocation")

        sslContextFactory.keyStorePath = keyStoreLocation
        sslContextFactory.setKeyStorePassword(getConfigProperty(EXTERNAL_KEYSTORE_PASSWORD))

        if (getConfigProperty(USE_SERVER_MTLS)) {
            val trustStoreLocation: String = getConfigProperty(EXTERNAL_TRUSTSTORE_LOCATION)
            logger.debug("Using external keystore: $trustStoreLocation")

            sslContextFactory.trustStorePath = trustStoreLocation
            sslContextFactory.setTrustStorePassword(getConfigProperty(EXTERNAL_TRUSTSTORE_PASSWORD))

            sslContextFactory.needClientAuth = true
        }
        return sslContextFactory
    }
}

