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

import com.fasterxml.jackson.databind.node.ObjectNode
import com.odim.odimclient.security.ClientSecurityConfigurator
import com.odim.simulator.CoreConfig.HTTP_CLIENT_BASIC_CREDENTIALS
import com.odim.simulator.CoreConfig.HTTP_CLIENT_WITH_BASIC
import com.odim.simulator.CoreConfig.SERVER_CERTIFICATE_VERIFICATION_ENABLED
import com.odim.simulator.CoreConfig.USE_CLIENT_TLS
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.http.HttpStatusCode.Companion.of
import java.util.Base64.getEncoder
import org.springframework.http.HttpEntity
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpMethod
import org.springframework.http.HttpMethod.DELETE
import org.springframework.http.HttpMethod.GET
import org.springframework.http.HttpMethod.PATCH
import org.springframework.http.HttpMethod.POST
import org.springframework.http.MediaType
import org.springframework.http.client.ClientHttpResponse
import org.springframework.http.client.HttpComponentsClientHttpRequestFactory
import org.springframework.http.converter.FormHttpMessageConverter
import org.springframework.util.MultiValueMap
import org.springframework.web.client.ResponseErrorHandler
import org.springframework.web.client.RestTemplate

class ErrorHandlerWithoutExceptions : ResponseErrorHandler {
    override fun handleError(response: ClientHttpResponse?) {
        // dont rise errors of status codes different than 2XX
    }
    override fun hasError(response: ClientHttpResponse?) = false
}

open class HttpClient(overrides: HttpClient.() -> Unit = {}) {
    private val logger = NonRepeatingHttpLogger(this.javaClass)
    private val rest: RestTemplate
    private val defaultHeaders: HttpHeaders

    var useTls = getConfigProperty<Boolean>(USE_CLIENT_TLS)
    var basicCredentialEnabled = getConfigProperty<Boolean>(HTTP_CLIENT_WITH_BASIC)
    var basicCredentials = getConfigProperty<String>(HTTP_CLIENT_BASIC_CREDENTIALS)
    var serverCertVerificationEnabled = getConfigProperty<Boolean>(SERVER_CERTIFICATE_VERIFICATION_ENABLED)

    init {
        overrides()
        defaultHeaders = HttpHeaders().apply {
            if (basicCredentialEnabled) {
                add("Authorization", "Basic ${getEncoder()
                        .encodeToString(basicCredentials.toByteArray())}")
            }
            add("Content-Type", "application/json")
            add("Accept", "application/json")
        }
        val requestFactory = HttpComponentsClientHttpRequestFactory().also {
            if (useTls) {
                it.httpClient = ClientSecurityConfigurator(serverCertVerificationEnabled).httpClient
            }
        }
        rest = RestTemplate(requestFactory).also {
            it.errorHandler = ErrorHandlerWithoutExceptions()
            it.messageConverters.add(FormHttpMessageConverter())
        }
    }

    protected open fun resolvePath(path: String) = path

    private fun extendDefaultHeaders(headers: HttpHeaders): HttpHeaders = HttpHeaders().apply {
        putAll(headers)
        putAll(defaultHeaders)
    }

    private fun HttpHeaders.setWwwFormUrlencodedToHeaders(): HttpHeaders {
        contentType = MediaType.APPLICATION_FORM_URLENCODED
        return this
    }

    fun get(path: String): Response = sendRequest(HttpEntity("", defaultHeaders), resolvePath(path), GET)

    fun get(path: String, headers: HttpHeaders): Response {
        val requestHeaders = extendDefaultHeaders(headers)
        return sendRequest(HttpEntity("", requestHeaders), resolvePath(path), GET)
    }

    fun post(path: String, body: ObjectNode): Response = sendRequest(HttpEntity(body.toString(), defaultHeaders), resolvePath(path), POST)

    fun post(path: String, body: String): Response = sendRequest(HttpEntity(body, defaultHeaders), resolvePath(path), POST)

    fun post(path: String, body: ObjectNode, headers: HttpHeaders): Response {
        val requestHeaders = extendDefaultHeaders(headers)
        return sendRequest(HttpEntity(body.toString(), requestHeaders), resolvePath(path), POST)
    }

    fun post(path: String, body: Map<String, String>, headers: HttpHeaders): Response = sendRequest(
            HttpEntity(mapToFormString(body),
            headers.setWwwFormUrlencodedToHeaders()),
            resolvePath(path),
            POST)

    private fun mapToFormString(map: Map<String, String>): String {
        var formString = ""
        for (element in map) {
            formString = formString.plus("${element.key}=${element.value}&")
        }
        formString.removeSuffix("&")
        return formString
    }

    fun patch(path: String, body: ObjectNode): Response = sendRequest(HttpEntity(body.toString(), defaultHeaders), resolvePath(path), PATCH)

    fun patch(path: String, body: ObjectNode, headers: HttpHeaders): Response {
        val requestHeaders = extendDefaultHeaders(headers)
        return sendRequest(HttpEntity(body.toString(), requestHeaders), resolvePath(path), PATCH)
    }

    fun delete(path: String): Response = sendRequest(HttpEntity("", defaultHeaders), resolvePath(path), DELETE)

    fun delete(path: String, headers: HttpHeaders): Response {
        val requestHeaders = extendDefaultHeaders(headers)
        return sendRequest(HttpEntity("", requestHeaders), resolvePath(path), DELETE)
    }

    private fun sendRequest(requestEntity: HttpEntity<*>, path: String, method: HttpMethod): Response {
        logger.logRequestOrCount(method, path, requestEntity)
        val responseEntity = rest.exchange(path, method, requestEntity, String::class.java)
        logger.logResponseOrCount(method, path, responseEntity)
        return Response(of(responseEntity.statusCodeValue), responseEntity.body ?: "{}", responseEntity.headers)
    }

    fun postMultipart(path: String, httpEntity: HttpEntity<MultiValueMap<String, Any>>): Response {
        val responseEntity = rest.exchange(path, HttpMethod.POST, httpEntity, String::class.java)
        return Response(HttpStatusCode.of(responseEntity.statusCodeValue), responseEntity.body ?: "{}",
                responseEntity.headers)
    }
}
