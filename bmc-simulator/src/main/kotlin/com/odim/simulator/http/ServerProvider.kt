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

import com.odim.simulator.CoreConfig.USE_SERVER_TLS
import com.odim.simulator.RequestProcessor
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.http.HttpMethod.DELETE
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.HttpMethod.OPTIONS
import com.odim.simulator.http.HttpMethod.PATCH
import com.odim.simulator.http.HttpMethod.POST
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.ServerFactory.create
import com.odim.utils.IpV4Address
import com.odim.utils.ipAddressOf
import com.odim.utils.isValidBodyJson
import io.javalin.Javalin
import io.javalin.core.security.BasicAuthCredentials
import io.javalin.http.Context
import io.javalin.websocket.WsHandler
import org.slf4j.LoggerFactory.getLogger
import java.io.EOFException
import java.util.function.Consumer

private const val REDFISH_ROOT = "/redfish/v1"
private const val REDFISH_ROOT_WILDCARD = "/redfish/v1/*"

class ServerProvider(val port: Int = 0, ipAddressRange: Sequence<IpV4Address> = generateSequence { ipAddressOf("127.0.0.1") }) {
    private val deployments = mutableSetOf<Deployment>()
    private val nonRepeatingLogger = NonRepeatingHttpLogger(this.javaClass)
    private val ipAddresses = IpAddressPool(ipAddressRange)
    private val log = getLogger(this.javaClass)


    /**
     * Deploys and serves Simulator on given port.
     * @param requestProcessor - Simulator to be served
     * @param port - port to be served on. 0 means random port
     * @param useTls - if simulator should use SSL protocol
     * @return Address of simulator and port serving on. Use deconstruction to retrieve both.
     */
    fun serve(requestProcessor: RequestProcessor,
              port: Int = this.port,
              useTls: Boolean = false,
              wsConfig: Map<String, Consumer<WsHandler>> = emptyMap()): Pair<String, Int> {
        val usesTls = useTls || getConfigProperty(USE_SERVER_TLS)

        val ipAddress = ipAddresses.poll()

        val server = with(create(ipAddress.toString(), port, usesTls, wsConfig)) {
            get(REDFISH_ROOT) { getRequestHandler(requestProcessor, it, GET) }
            get(REDFISH_ROOT_WILDCARD) { getRequestHandler(requestProcessor, it, GET) }
            post(REDFISH_ROOT_WILDCARD) { getRequestHandler(requestProcessor, it, POST) }
            patch(REDFISH_ROOT_WILDCARD) { getRequestHandler(requestProcessor, it, PATCH) }
            delete(REDFISH_ROOT_WILDCARD) { getRequestHandler(requestProcessor, it, DELETE) }
            options(REDFISH_ROOT_WILDCARD) { getRequestHandler(requestProcessor, it, OPTIONS) }
        }
        val deployment = Deployment(requestProcessor, server, ipAddress.toString(), usesTls)
        deployments.add(deployment)
        deployment.start()
        return deployment.address to deployment.port
    }

    /**
     * Stops and removes deployment of given Simulator.
     * @param bmcAddress - Address IP of BMC simulator to be stopped
     */
    fun stop(bmcAddress: IpV4Address) {
        deployments.singleOrNull {
            it.ipV4Address == bmcAddress
        }?.let {
            it.stop()
            deployments.remove(it)
        }
    }


    /**
     * Stops and removes deployment of given Simulator.
     * @param requestProcessor - Simulator to be stopped
     */
    fun stop(requestProcessor: RequestProcessor) {
        deployments.singleOrNull {
            it.requestProcessor == requestProcessor
        }?.let {
            it.stop()
            deployments.remove(it)
        }
    }

    /**
     * Stops and removes all Simulator deployments.
     */
    fun stopAll() {
        deployments.forEach { it.stop() }
        deployments.clear()
    }

    private fun getRequestHandler(requestProcessor: RequestProcessor, ctx: Context, httpMethod: HttpMethod) {
        val basicAuthCredentials = if (ctx.header("Authorization").isNullOrEmpty()) null else ctx.basicAuthCredentials()
        val request = createRequest(ctx, httpMethod, basicAuthCredentials)
        val response = createResponse(request, ctx, requestProcessor)
        ctx.status(response.code.value).result(response.body)
        applyHeaders(response, ctx)
    }

    private fun createResponse(
        request: Request?,
        ctx: Context,
        requestProcessor: RequestProcessor
    ): Response {
        val response: Response
        if (isRequestCorrect(request, ctx)) {
            response = requestProcessor.createResponse(request!!)
            nonRepeatingLogger.logSimulatorRequest(requestProcessor.className, request, response)
        } else response = badRequest("Invalid Json")
        return response
    }

    private fun isRequestCorrect(request: Request?, ctx: Context) =
        if (request == null) false else request.isValidBodyJson() || ctx.isMultipart()

    private fun createRequest(
        ctx: Context,
        httpMethod: HttpMethod,
        basicAuthCredentials: BasicAuthCredentials?
    ) = try {
        if (ctx.isMultipartFormData()) {
            Request(
                httpMethod, "", ctx.path(), ctx.headerMap(), ctx.queryParamMap(), ctx.queryString(),
                basicAuthCredentials, ctx.req
            )
        } else {
            Request(
                httpMethod, ctx.body(), ctx.path(), ctx.headerMap(), ctx.queryParamMap(), ctx.queryString(),
                basicAuthCredentials
            )
        }
    } catch (e: EOFException) {
        log.warn("{} for {} request at {}", e.message, ctx.method(), ctx.url())
        null
    }

    private fun applyHeaders(response: Response, context: Context) {
        response.headers.forEach { headerEntry ->
            val key = headerEntry.key
            headerEntry.value.forEach {
                context.header(key, it)
            }
        }
    }
}

private class Deployment(val requestProcessor: RequestProcessor,
                         private val server: Javalin,
                         private val ip: String,
                         private val usesTls: Boolean) {
    private val logger = getLogger(this.javaClass)

    val ipV4Address: IpV4Address = ipAddressOf(ip)
    val address: String by lazy { "${if (usesTls) "https" else "http"}://$ip:${server.port()}$REDFISH_ROOT" }
    val port: Int by lazy { server.port() }

    init {
        server.events { listener ->
            listener.serverStarted {
                logger.info("${requestProcessor.className} started at: $address")
                requestProcessor.welcomeMessage?.let { logger.info(it) }
            }
            listener.serverStopped {
                logger.info("${requestProcessor.className} at $address stopped")
            }
        }
    }

    fun start() {
        requestProcessor.onStart()
        server.start()
    }

    fun stop() {
        requestProcessor.onStop()
        server.stop()
    }
}

private class IpAddressPool(addresses: Sequence<IpV4Address>) {
    private val ipAddressIterator = addresses.iterator()

    fun poll() = when {
        ipAddressIterator.hasNext() -> ipAddressIterator.next()
        else -> throw error("All available IP Addresses have been polled")
    }
}
