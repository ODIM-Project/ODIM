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

package com.odim.tools.coverage

import org.gradle.api.DefaultTask
import org.gradle.api.GradleException
import org.gradle.api.tasks.TaskAction

class CoverageLevelVerificationTask extends DefaultTask {

    @TaskAction
    void verifyCoverage() {
        Node doc = parseXmlReport()
        def counter = project.extensions.getByType(CoverageLevelVerification).counter
        def globalCoverage = doc.counter.find { it.@type == counter }
        def covered = globalCoverage.@covered.toInteger()
        def missed = globalCoverage.@missed.toInteger()
        def coverageRatio = (covered / (covered + missed)).toFloat().trunc(3)
        checkCoverageLevel(coverageRatio)
    }

    private Node parseXmlReport() {
        File xmlReport = project.file(project.extensions.getByType(CoverageLevelVerification).xmlReport)
        def parser = new XmlParser()
        parser.with {
            setFeature("http://apache.org/xml/features/disallow-doctype-decl", false)
            setFeature("http://apache.org/xml/features/nonvalidating/load-external-dtd", false);
        }
        def doc = parser.parse(xmlReport)
        doc
    }

    private def checkCoverageLevel(Float coverageRatio) {
        CoverageLevelVerification setCoverageLevel = project.extensions.getByType(CoverageLevelVerification)
        setCoverageLevel.with {
            project.logger.quiet("""verifyCoverage: {} coverage bounds = [{}, {}], computed coverage = {}""",
                    counter, coverageMinThreshold, coverageMinThreshold + coverageExcessMax, coverageRatio)
            if (coverageRatio < coverageMinThreshold) {
                throw new GradleException("$counter coverage level below threshold ($coverageRatio < $coverageMinThreshold)")
            }
            if (coverageRatio > coverageMinThreshold + coverageExcessMax) {
                throw new GradleException("$counter coverage minimal threshold set too low ($coverageRatio > $coverageMinThreshold + $coverageExcessMax); increase minThreshold")
            }
        }
    }

}
