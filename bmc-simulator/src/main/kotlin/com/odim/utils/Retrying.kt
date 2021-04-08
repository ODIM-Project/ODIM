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

package com.odim.utils

import java.time.Duration
import java.time.LocalDateTime.now
import java.util.concurrent.TimeoutException

fun <R, E : Throwable> retry(timeout: Duration, code: TryResults.() -> TryResult<R, E>) =
        retry(timeout, code, { throw TimeoutException("Timeout = $timeout exceeded. Reason: ${it.reason}") })


fun <R, E : Throwable> retryOrDefault(default: R, timeout: Duration, code: TryResults.() -> TryResult<R, E>) =
        retry(timeout, code, { default })

private fun <R, E : Throwable>
        retry(timeout: Duration, code: TryResults.() -> TryResult<R, E>, onTimeout: (last: NotYet) -> R): R {
    val timeoutTime = now() + timeout

    while (true) {
        val tryResult = requireToNotThrow("Retryable code") { TryResults.code() }
        when (tryResult) {
            is Successful -> return tryResult.result
            is Failure -> throw tryResult.throwable
            is NotYet -> if (now().isAfter(timeoutTime)) return onTimeout(tryResult)
        }
    }
}

sealed class TryResult<out R, out E : Throwable>
class Successful<R>(val result: R) : TryResult<R, Nothing>()
class Failure<E : Throwable>(val throwable: E) : TryResult<Nothing, E>()
class NotYet(val reason: String) : TryResult<Nothing, Nothing>()

object TryResults {
    fun <R> success(result: R) = Successful(result)
    fun successIf(success: Boolean) = if (success) Successful(true) else notYetWithLastValue(false)
    fun success() = Successful(Unit)

    fun <R> notYetWithLastValue(last: R) = NotYet("Last returned value was: $last")

    fun notYet(reason: String) = NotYet(reason)

    @Suppress("TooGenericExceptionCaught", "InstanceOfCheckForException")
    inline fun <R, reified E : Throwable> notYetIfThrows(code: TryResults.() -> TryResult<R, E>) =
            try {
                code()
            } catch (e: Throwable) {
                if (e is E) {
                    notYet("Exception has been thrown: $e")
                } else {
                    throw e
                }
            }

    fun <E : Throwable> fail(throwable: E) = Failure(throwable)
}

@Suppress("TooGenericExceptionCaught")
private fun <T> requireToNotThrow(name: String, code: () -> T): T =
        try {
            code()
        } catch (e: Exception) {
            throw IllegalArgumentException("$name should never throw", e)
        }
