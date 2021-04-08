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

enum class HttpStatusCode(val value: Int) {
    OK(200),
    CREATED(201),
    ACCEPTED(202),
    NO_CONTENT(204),
    PARTIAL_CONTENT(206),
    MOVED_PERMANENTLY(301),
    FOUND(302),
    NOT_MODIFIED(304),
    BAD_REQUEST(400),
    UNAUTHORIZED(401),
    FORBIDDEN(403),
    NOT_FOUND(404),
    METHOD_NOT_ALLOWED(405),
    NOT_ACCEPTABLE(406),
    CONFLICT(409),
    PRECONDITION_FAILED(412),
    INTERNAL_SERVER_ERROR(500),
    NOT_IMPLEMENTED(501),
    BAD_GATEWAY(502),
    SERVICE_UNAVAILABLE(503),
    GATEWAY_TIMEOUT(504);

    val is2XX get() = this.value in 200..299

    companion object {
        fun of(code: Int): HttpStatusCode {
            val httpCode = values().firstOrNull { it.value == code }
            if (httpCode != null) {
                return httpCode
            }
            throw IllegalStateException("Cannot get HttpCode entity for $code code.")
        }
    }
}
