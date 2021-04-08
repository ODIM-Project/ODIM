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

import com.odim.simulator.CoreConfig.SERVE_IP
import com.odim.simulator.CoreConfig.SERVE_PORT
import com.odim.simulator.SimulatorConfig.Config.appendValuesFromExternalConfig
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.commandline.SimulatorRunner.serveSimulatorByName
import org.slf4j.LoggerFactory.getLogger
import picocli.CommandLine
import picocli.CommandLine.Command
import picocli.CommandLine.Option
import picocli.CommandLine.Parameters
import kotlin.reflect.full.findAnnotation

private const val COPYRIGHTS = "\nCopyright (c) Intel Corporation"

fun main(args: Array<String>) {
    CommandLine(SimulatorsCmd()).execute(*args)
}

@Command(
        name = "simulator-runner.jar",
        description = ["\nRedfish simulators runner"],
        subcommands = [ListSimulatorsCmd::class, RunSimulatorsCmd::class],
        footer = [
            "Example of run BMC simulator:",
            "  simulator-runner.jar run BMC --port 1234",
            COPYRIGHTS
        ])
class SimulatorsCmd : Runnable {
    override fun run() {
        println("Usage: ${this::class.findAnnotation<Command>()!!.name} [run|list]\n$COPYRIGHTS")
    }
}

@Command(name = "list", aliases = ["l"], description = ["List available simulators", COPYRIGHTS])
class ListSimulatorsCmd : Runnable {
    override fun run() {
        println("Available simulators: ${SimulatorsFinder.getAllNames()}")
    }
}

@Command(name = "run", aliases = ["r"], description = ["Start simulators with specified name"], footer = [COPYRIGHTS])
class RunSimulatorsCmd : Runnable {
    private val logger = getLogger(this.javaClass)

    @Parameters(paramLabel = "SIMULATOR", description = ["Simulator to be run - " +
            "command 'list' will return list of simulators"])
    private var name: String = ""

    @Option(names = ["-i", "--ip"], description = ["Bind to ip address"])
    private var ip: String? = null

    @Option(names = ["-p", "--port"], description = ["Bind to port"])
    private var port: Int? = null

    @Option(names = ["-c", "--config"], description = ["External config file path"])
    private var externalConfigFilePath: String? = "simulator-config.json"

    override fun run() {
        logger.info("Running simulator with name: '{}' on port {}", name, port)
        appendValuesFromExternalConfig(externalConfigFilePath)
        val serveIp = ip ?: getConfigProperty(SERVE_IP)
        val servePort = port ?: getConfigProperty(SERVE_PORT)
        serveSimulatorByName(name, serveIp, servePort)
    }
}
