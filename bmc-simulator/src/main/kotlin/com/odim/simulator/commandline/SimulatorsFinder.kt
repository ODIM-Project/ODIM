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
import com.odim.simulator.SimulatorName
import com.odim.utils.Timing.logExecutionTime
import org.springframework.data.util.AnnotatedTypeScanner
import kotlin.reflect.KClass
import kotlin.reflect.full.createInstance
import kotlin.reflect.full.findAnnotation

object SimulatorsFinder {

    private const val REPO_PACKAGE_PATH = "com.odim.simulator.repo"

    fun getAllNames(): List<String> {
        return mutableListOf<String>().also { result ->
            getAllSimulatorClassesFromRepo().forEach {
                it.findAnnotation<SimulatorName>()?.run { result.add(name) }
            }
        }
    }

    fun getSimulatorByName(name: String): RequestProcessor {
        val repoClasses = getAllSimulatorClassesFromRepo()
        val clazz = firstClassWithName(repoClasses, name)

        requireNotNull(clazz, { "No class with $name found." })
        return logExecutionTime("Simulator init time: {} ms") {
            clazz.createInstance() as RequestProcessor
        }
    }

    private fun getAllSimulatorClassesFromRepo(): List<KClass<out Any>> {
        val annotatedTypeScanner = AnnotatedTypeScanner(SimulatorName::class.java)
        return annotatedTypeScanner.findTypes(REPO_PACKAGE_PATH).map { it.kotlin }
    }

    private fun firstClassWithName(repoClasses: List<KClass<out Any>>, name: String) =
            repoClasses.firstOrNull { isClassAnnotatedWithName(it, name) }

    private fun isClassAnnotatedWithName(clazz: KClass<out Any>, name: String): Boolean {
        val nameAnnotation = clazz.findAnnotation<SimulatorName>()
        return name.equals(nameAnnotation?.name, ignoreCase = true)
    }
}
