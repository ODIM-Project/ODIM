package com.odim.simulator.repo.bmc.behaviors

import com.odim.simulator.behaviors.Behavior
import com.odim.simulator.behaviors.BehaviorDataStore
import com.odim.simulator.behaviors.BehaviorResponse
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.http.Request
import com.odim.simulator.http.Response
import com.odim.simulator.http.Response.Companion.badRequest
import com.odim.simulator.http.Response.Companion.noContent
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.structure.Item
import com.odim.utils.isNull
import com.odim.utils.toMap

class PostOnSimpleUpdate : Behavior {

    private val requestJsonProperties = mutableListOf(
        "ImageURI",
        "Password",
        "Targets",
        "TransferProtocol",
        "Username",
    )

    override fun run(
        tree: ResourceTree,
        item: Item,
        request: Request,
        response: Response,
        dataStore: BehaviorDataStore
    ): BehaviorResponse {
        if (!validateRequest(request)) return nonTerminal(badRequest())
        return nonTerminal(noContent())
    }

    private fun validateRequest(request: Request): Boolean {
        return request.json!!.toMap().keys.stream().allMatch {
            requestJsonProperties.contains(it)
        } && !request.json!!.isNull("ImageURI")
    }
}
