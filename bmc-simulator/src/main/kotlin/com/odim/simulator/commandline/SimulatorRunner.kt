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

package com.odim.simulator.commandline

import com.odim.simulator.RequestProcessor
import com.odim.simulator.commandline.SimulatorsFinder.getSimulatorByName
import com.odim.simulator.http.ServerProvider
import com.odim.utils.ipAddressOf
import java.lang.Runtime.getRuntime

object SimulatorRunner {
    fun serveSimulatorByName(name: String, ip: String, port: Int): String {
        val simulator = getSimulatorByName(name)
        return serveSimulator(simulator, ip, port)
    }

    private fun serveSimulator(simulator: RequestProcessor, ip: String, port: Int): String {
        val provider = ServerProvider(port, ipAddressOf(ip)..ipAddressOf(ip))
        val (address) = provider.serve(simulator, port)
        registerShutDownHook { provider.stopAll() }
        return address
    }
}

fun registerShutDownHook(runnable: () -> Unit) {
    getRuntime().addShutdownHook(Thread(runnable))
}
