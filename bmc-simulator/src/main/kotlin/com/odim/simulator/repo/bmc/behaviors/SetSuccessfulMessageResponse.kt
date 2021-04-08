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

package com.odim.simulator.repo.bmc.behaviors

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.dsl.makeJson
import com.odim.simulator.http.HttpStatusCode
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item

class SetSuccessfulMessageResponse : Behavior {
    override fun run(tree: ResourceTree, item: Item, request: Request, response: Response, dataStore: BehaviorDataStore): BehaviorResponse {
        val messageResponseBody = makeJson {
            "@odata.type" to "#Message.v1_0_5.Message"
            "MessageId" to "Base.1.1.0.Success"
            "Message" to "Successfully Completed Request"
            "Severity" to "OK"
            "Resolution" to "None"
        }.toString()

        return nonTerminal(Response(HttpStatusCode.OK, messageResponseBody))
    }
}
