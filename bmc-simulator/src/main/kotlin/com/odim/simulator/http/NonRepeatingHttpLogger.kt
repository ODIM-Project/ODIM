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

import com.odim.utils.prettyPrint
import com.odim.utils.toJson
import org.slf4j.Logger
import org.slf4j.LoggerFactory.getLogger
import org.springframework.http.HttpEntity
import org.springframework.http.HttpMethod
import org.springframework.http.HttpStatus.I_AM_A_TEAPOT
import org.springframework.http.MediaType
import org.springframework.http.ResponseEntity

class NonRepeatingHttpLogger(classOnBehalfOfWhichIWillLog: Class<*>) {

    private val logger = getLogger(classOnBehalfOfWhichIWillLog)
    private val logRequester = LogRequester(logger)
    private val logResponser = LogResponser(logger)
    private val logSimulatorRequester = LogSimulatorRequester(logger)

    internal fun logRequestOrCount(method: HttpMethod, path: String, requestEntity: HttpEntity<*>) =
            logRequester.logRequestOrCount(method, path, requestEntity)

    internal fun logResponseOrCount(method: HttpMethod, path: String, responseEntity: ResponseEntity<String>) =
            logResponser.logResponseOrCount(method, path, responseEntity)

    internal fun logSimulatorRequest(simulatorName: String?, request: Request, response: Response) =
            logSimulatorRequester.logSimulatorRequest(simulatorName, request, response)
}

private class LogRequester(val logger: Logger) {
    private var lastRequestEntity: HttpEntity<*>? = null
    private var lastPath: String? = null
    private var lastMethod: HttpMethod? = null
    private var requestCounter = 0

    fun logRequestOrCount(method: HttpMethod, path: String, requestEntity: HttpEntity<*>) {
        if (requestEntity == lastRequestEntity && path == lastPath && method == lastMethod) {
            logger.debug("request ({} {}) repeated {} times", method, path, ++requestCounter)
        } else {
            val body = requestEntity.body
            if (body is String) {
                val contentType = requestEntity.headers.contentType
                val message = when {
                    contentType != null && contentType.includes(MediaType.APPLICATION_FORM_URLENCODED) -> body
                    body.isNotEmpty() -> body.toJson().prettyPrint()
                    else -> ""
                    }
                logger.debug("Sending $method $path $message")
            }
        }
            lastRequestEntity = requestEntity
            lastPath = path
            lastMethod = method
            requestCounter = 0
        }
    }

private class LogResponser(val logger: Logger) {
    private var lastResponseEntity = ResponseEntity("", I_AM_A_TEAPOT)
    private var responseCounter = 0

    fun logResponseOrCount(method: HttpMethod, path: String, responseEntity: ResponseEntity<String>) {
        if (lastResponseEntity.body == responseEntity.body
                && responseEntity.statusCode == lastResponseEntity.statusCode) {
            logger.debug("response ({}) repeated {} times", responseEntity.statusCode, ++responseCounter)
        } else {
            logger.debug("Response to $method $path: ${responseEntity.statusCodeValue} - ${responseEntity.statusCodeValue}" +
                    if (!responseEntity.body.isNullOrEmpty()) "\n${responseEntity.body.prettyPrint()}" else "")
            lastResponseEntity = responseEntity
            responseCounter = 0
        }
    }
}

private class LogSimulatorRequester(val logger: Logger) {
    private var lastSimulatorName: String? = null
    private var lastRequest: Request? = null
    private var simulatorRequestCounter = 0

    fun logSimulatorRequest(simulatorName: String?, request: Request, response: Response) {
        if (lastSimulatorName == simulatorName && lastRequest == request) {
            logger.debug("Request to simulator `{}` repeated {} times ({} {})", simulatorName, ++simulatorRequestCounter, request.method, request.url)
        } else {
            logger.debug("${simulatorName} handling request: ${request.prettyPrint()}, response: ${response.prettyPrint()}")
            simulatorRequestCounter = 0
            lastSimulatorName = simulatorName
            lastRequest = request
        }
    }
}
