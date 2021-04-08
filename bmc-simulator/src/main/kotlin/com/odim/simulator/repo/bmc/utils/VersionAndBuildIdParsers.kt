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

import java.io.ByteArrayInputStream

const val COULDNT_PARSE = "couldn't parse"
private const val SDR_VERSION_PREFIX = "_LF_VERSION"
private const val BIOS_BUILD_ID_LENGTH = 4
private const val BMC_VERSION_LENGTH = 4

fun parseBIOSVersionAndBuildId(filename: String): Pair<String, String> {
    val versionAndBuildId = trimBIOSDetails(filename)

    val version = extractBIOSVersion(versionAndBuildId)
    val buildId = extractBIOSBuildId(versionAndBuildId)

    return Pair(version, buildId)
}

fun parseBMCVersionAndBuildId(filename: String): Pair<String, String> {
    val versionAndBuildId = trimBMCPlatformName(filename)

    val version = extractBMCVersion(versionAndBuildId)
    val buildId = extractBMCBuildId(versionAndBuildId)

    return Pair(version, buildId)
}

fun parseSDRVersion(lineWithVersion: String): Pair<String, String> {
    val version = lineWithVersion.substringAfter(SDR_VERSION_PREFIX).trim { it.isWhitespace() || it == '\'' }

    return Pair(version, "")
}

fun extractSDRLine(bytes: ByteArray?): String {
    val inputStream = ByteArrayInputStream(bytes)
    return inputStream.bufferedReader().useLines {
        it.find(::containsVersionPrefix)
    } ?: COULDNT_PARSE
}

private fun containsVersionPrefix(line: String) = line.contains(SDR_VERSION_PREFIX)

private fun trimBMCPlatformName(filename: String) = filename.substringAfter('_')

private fun trimBIOSDetails(filename: String) = filename.substringBefore('_')

private fun extractBMCVersion(versionAndBuildId: String) = versionAndBuildId

private fun extractBMCBuildId(versionAndBuildId: String) = versionAndBuildId.drop(BMC_VERSION_LENGTH)

private fun extractBIOSVersion(versionAndBuildId: String) = versionAndBuildId

private fun extractBIOSBuildId(versionAndBuildId: String) = versionAndBuildId.takeLast(BIOS_BUILD_ID_LENGTH)
