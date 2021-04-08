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

package com.odim.simulator.repo.bmc.utils

import kotlin.random.Random

class PropertyValueGenerator private constructor() {
    companion object Factory {
        fun generateSerial(): String {
            val charSet: List<Char> = ('A'..'F') + ('0'..'9')
            return (1..12).map { charSet.random() }.joinToString("")
        }

        fun generateMacAddress(): String {
            return (1..6).joinToString("-") {
                Integer.toHexString(Random.nextInt(255)).padStart(2, '0')
            }
        }

        fun generateIpAddressPart(): String {
            return (1..3).joinToString(".") {
                Random.nextInt(255).toString()
            }
        }

        fun generateEntryLogMessage(): String {
            val id = ('a'..'z').map { it }.shuffled().subList(0, 7).joinToString("")
            return "Priority:6 login failed: $id"
        }

        fun generateDateTime() =
                "2019-2020-" + Random.nextInt(1, 12) + "-" + Random.nextInt(1, 30) + "-26T" +
                        Random.nextInt(1, 24) + ":37:" + Random.nextInt(0, 59) + "59+00:00"

        fun generateVersion(): String {
            val charSet: List<Char> = ('a'..'z') + ('0'..'9')
            return Random.nextInt(10).toString().plus(".").plus(Random.nextInt(100).toString()).plus(".")
                    .plus((0..8).map { charSet.random() }.joinToString(""))
        }
    }
}
