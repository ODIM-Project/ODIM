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

package com.odim.simulator.repo.behaviors

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.dsl.merger.Merger
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.created
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.EmbeddedArray
import com.odim.simulator.tree.structure.EmbeddedObjectType.METRIC_VALUE
import com.odim.simulator.tree.structure.EmbeddedResourceArray
import com.odim.simulator.tree.structure.GeneratedValue
import com.odim.simulator.tree.structure.Item
import com.odim.simulator.tree.structure.LinkableResource
import com.odim.simulator.tree.structure.ObjectArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT_DEFINITION
import com.odim.utils.doLaunch
import kotlinx.coroutines.InternalCoroutinesApi
import kotlinx.coroutines.Job
import kotlinx.coroutines.NonCancellable.isActive
import java.lang.Thread.sleep

@InternalCoroutinesApi
@Suppress("ReturnCount", "LongMethod")
class PostOnMetricReportDefinitionCollection : Behavior {
    var metricReportValuesThreads = mutableListOf<Job>()

    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val metricReportDefinition = tree.create(METRIC_REPORT_DEFINITION)
        Merger.merge(tree, metricReportDefinition, request)

        val metricReport = tree.create(METRIC_REPORT)

        val telemetryService = tree.root.traverse<Resource>("TelemetryService")
        telemetryService(metricReport, metricReportDefinition)

        metricReport.traverse<LinkableResource>("MetricReportDefinition").addLink(metricReportDefinition, metricReport)
        metricReportDefinition.traverse<LinkableResource>("MetricReport").addLink(metricReport, metricReportDefinition)

        metricReportValuesThreads.add(doLaunch {
            while (isActive) {
                val repeat = metricReportDefinition.traverse<Int>("Schedule/RecurrenceInterval")

                val report = metricReportDefinition.traverse<LinkableResource>("MetricReport").getElement() as Resource
                val metrics = metricReportDefinition.traverse<ObjectArray>("Metrics")

                val metricValues = ObjectArray(METRIC_VALUE)
                metrics.forEach { metric ->
                    @Suppress("UNCHECKED_CAST")
                    (metric["MetricProperties"] as EmbeddedArray<String>).forEach { metricProperty ->
                        val property = metricProperty.split("#/")
                        val search = tree.search(property[0]) as Resource
                        metricValues.add(tree.createEmbeddedObject(METRIC_VALUE) {
                            "MetricId" to metric["MetricId"]
                            "MetricValue" to if (Regex("[a-zA-Z]+/[0-9]+/.*").matches(property[1])) {
                                val embeddedSearch = property[1].split("/")
                                search.traverse<EmbeddedResourceArray>(embeddedSearch[0])
                                        .data[embeddedSearch[1].toInt()]
                                        .traverse<GeneratedValue<*>>(embeddedSearch[2])
                                        .getValue()
                            } else {
                                search.traverse<GeneratedValue<*>>(property[1]).getValue()
                            }
                            "MetricProperty" to metricProperty
                        })
                    }
                }
                report.data["MetricValues"] = metricValues
                report.data["ReportSequence"] = report.traverse<String>("ReportSequence").toInt().inc().toString()
                sleep(repeat * 1000L)
            }
        })
        return terminal(created(metricReportDefinition))
    }
}
