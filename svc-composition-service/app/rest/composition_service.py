#(C) Copyright [2022] American Megatrends International LLC
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.

import logging
import json
import copy
import uuid
from http import HTTPStatus
from utilities.client import Client
from db.persistant import RedisClient


class CompositonService():
    def __init__(self):
        self.client = Client()
        self.redis = RedisClient()

    def get_cs(self):
        res = {
            "@odata.context": "/redfish/v1/$metadata#CompositionService",
            "@odata.type": "#CompositionService.v1_2_0.CompositionService",
            "@odata.id": "/redfish/v1/CompositionService",
            "Id": "CompositionService",
            "Name": "Composition Service",
            "Status": {
                "State": "Enabled",
                "Health": "OK"
            },
            "ServiceEnabled": True,
            "AllowOverprovisioning": True,
            "AllowZoneAffinity": True,
            "ResourceBlocks": {
                "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks"
            },
            "ResourceZones": {
                "@odata.id": "/redfish/v1/CompositionService/ResourceZones"
            },
            "ActivePool": {
                "@odata.id": "/redfish/v1/CompositionService/ActivePool"
            },
            "CompositionReservations": {
                "@odata.id":
                "/redfish/v1/CompositionService/CompositionReservations"
            },
            "FreePool": {
                "@odata.id": "/redfish/v1/CompositionService/FreePool"
            },
            "Actions": {
                "#CompositionService.Compose": {
                    "target":
                    "/redfish/v1/CompostionService/Actions/CompositionService.Compose",
                }
            },
            "ReservationDuration": None,
            "Oem": {}
        }

        return res, HTTPStatus.OK

    def compose_action(self, req):
        res = {}
        code = HTTPStatus.OK

        if req is None:
            res["Error"] = "Request Body is empty"
            return res, HTTPStatus.BAD_REQUEST

        try:

            stanzas = req["Manifest"]["Stanzas"]

            for stanza in stanzas:
                if stanza["StanzaType"] == "ComposeSystem":
                    compose_res, compose_code = self.create_compose_system(
                        stanza["Request"])
                    stanza["Response"] = compose_res
                elif stanza["StanzaType"] == "DecomposeSystem":
                    decompose_res, decompose_code = self.decompose_system(
                        stanza["Request"])
                    stanza["Response"] = decompose_res
            res = req
            code = HTTPStatus.OK
        except Exception as err:
            logging.error(
                "Unable to process the compose action. Error: {e}".format(
                    e=err))
            res = {
                "Error":
                "Unable to process the compose action. Error: {e}".format(
                    e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            return res, code

    def create_compose_system(self, req):

        res = {}
        compose_sys = {}
        code = HTTPStatus.CREATED
        system_data = {}
        pipe = self.redis.pipeline()
        try:
            logging.info("Initialize create compose system")

            if not (req.get("Links") and req["Links"].get("ResourceBlocks")):
                logging.error("Unable to find ResourceBlocks Links")
                res = {"Error": "Unable to find ResourceBlocks Links"}
                code = HTTPStatus.BAD_REQUEST
                return

            logging.debug(
                "Create Compose system request body is: {req}".format(req=req))

            for block_uri in req["Links"]["ResourceBlocks"]:

                rs_block = self.redis.get("ResourceBlocks:{block_uri}".format(
                    block_uri=block_uri["@odata.id"]))
                if rs_block is None:
                    logging.error(
                        "The Resource Block {id} is not found in the database".
                        format(id=block_uri["@odata.id"]))
                    res = {
                        "Error":
                        "The Resource Block {id} is not valid".format(
                            id=block_uri["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return
                rs_block = json.loads(rs_block)
                logging.debug("ResourceBlock {id} data is: {data}".format(
                    id=block_uri["@odata.id"], data=rs_block))

                if "ComputerSystem" in rs_block["ResourceBlockType"]:

                    system_data = self.client.process_get_request(
                        rs_block["ComputerSystems"][0]["@odata.id"])
                    if system_data is None:
                        logging.error(
                            "The ComputerSystem {sys_id} from Resource Block {id} is not found valid"
                            .format(sys_id=rs_block["ComputerSystems"][0]
                                    ["@odata.id"],
                                    id=block_uri["@odata.id"]))
                        res = {"Error": "Get ComputerSystem is failed"}
                        code = HTTPStatus.BAD_REQUEST
                        return

                    compose_sys["Id"] = str(uuid.uuid1())
                    compose_sys["@odata.id"] = system_data[
                        "@odata.id"].replace(system_data["Id"],
                                             compose_sys["Id"])
                    compose_sys[
                        "Name"] = "Computer system composed from physical system"
                    compose_sys["@odata.type"] = system_data["@odata.type"]
                    logging.debug("New Compose System Id is: {id}".format(
                        id=compose_sys["Id"]))

                    compose_sys["Processors"] = copy.deepcopy(system_data["Processors"])
                    compose_sys["Memory"] = copy.deepcopy(system_data["Memory"])
                    compose_sys["Storage"] = copy.deepcopy(system_data["Storage"])

                    if system_data.get("Links"):
                        compose_sys["Links"] = copy.deepcopy(
                            system_data["Links"])

                    if not compose_sys.get("Links"):
                        compose_sys["Links"] = {
                            "ResourceBlocks": [],
                            "SupplyingComputerSystems": []
                        }
                    else:
                        if not compose_sys["Links"].get("ResourceBlocks"):
                            compose_sys["Links"]["ResourceBlocks"] = []
                        if not compose_sys["Links"].get(
                                "SupplyingComputerSystems"):
                            compose_sys["Links"][
                                "SupplyingComputerSystems"] = []

                    compose_sys["Links"]["ResourceBlocks"].append(
                        {"@odata.id": rs_block["@odata.id"]})
                    compose_sys["Links"]["SupplyingComputerSystems"].append(
                        {"@odata.id": system_data["@odata.id"]})

                    """
                    compose_sys["Actions"] = {
                        "#ComputerSystem.AddResourceBlock": {
                            "target":
                            "{compose_sys_uri}/{add_resource}".format(
                                compose_sys_uri=compose_sys["@odata.id"],
                                add_resource=
                                "Actions/ComputerSystem.AddResourceBlock")
                        },
                        "#ComputerSystem.RemoveResourceBlock": {
                            "target":
                            "{compose_sys_uri}/{remove_resource}".format(
                                compose_sys_uri=compose_sys["@odata.id"],
                                remove_resource=
                                "Actions/ComputerSystem.RemoveResourceBlock")
                        }
                    }
                    """
                if (rs_block["CompositionStatus"]["MaxCompositions"] <=
                        rs_block["CompositionStatus"]["NumberOfCompositions"]):
                    logging.error(
                        "NumberOfCompositions are excided to MaxCompositions for this Resource Block {id}"
                        .format(id=rs_block["Id"]))
                    res = {
                        "Error":
                        "NumberOfCompositions are excided to MaxCompositions for this Resource Block {id}"
                        .format(id=rs_block["Id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return

                rs_block["CompositionStatus"]["NumberOfCompositions"] += 1

                if (rs_block["Pool"] != "Free") or (
                        rs_block["CompositionStatus"]["CompositionState"]
                        == "Composed"):
                    logging.error(
                        "The Resource Block {rb_uri} is already used by other composed system"
                        .format(rb_uri=block_uri["@odata.id"]))
                    res = {
                        "Error":
                        "The Resource Block {rb_uri} is already used by other composed system"
                        .format(rb_uri=block_uri["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return

                rs_block["Pool"] = "Active"
                rs_block["CompositionStatus"]["CompositionState"] = "Composed"
                logging.info(
                    "The properties 'Pool' and 'CompositionState' of Resource Block {rb_uri}"
                    .format(rb_uri=block_uri["@odata.id"]))
                if not rs_block.get("Links"):
                    rs_block["Links"] = {"ComputerSystems": []}
                elif not rs_block["Links"].get("ComputerSystems"):
                    rs_block["Links"]["ComputerSystems"] = []
                rs_block["Links"]["ComputerSystems"].append(
                    {"@odata.id": compose_sys["@odata.id"]})

                pipe.set(
                    "ResourceBlocks:{rb_uri}".format(
                        rb_uri=block_uri["@odata.id"]), json.dumps(rs_block))
                pipe.srem("FreePool", block_uri["@odata.id"])
                pipe.sadd("ActivePool", block_uri["@odata.id"])

            compose_sys["SystemType"] = "Composed"
            pipe.set(
                "ComputerSystem:{compose_uri}".format(
                    compose_uri=compose_sys["@odata.id"]),
                json.dumps(json.dumps(compose_sys)))
            pipe.execute()
            logging.info("Compose System {id} is created successfully".format(
                id=compose_sys["Id"]))
            res = compose_sys
            code = HTTPStatus.CREATED

        except Exception as err:
            logging.error(
                "Unable to create composed system. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to create composed system. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            pipe.reset()
            return res, code

    def decompose_system(self, req):
        res = {}
        code = HTTPStatus.OK
        pipe = self.redis.pipeline()

        try:
            logging.info("Initialize Decompose System")

            if not (req.get("Links") and req["Links"].get("ComputerSystems")):
                res = {
                    "Error":
                    "ComputerSystems Links is empty. Provide the ComputerSystem link and resubmit the request"
                }
                code = HTTPStatus.BAD_REQUEST
                return

            logging.debug(
                "DecomposeSystem request body: {req}".format(req=req))

            for system_id in req["Links"]["ComputerSystems"]:
                system_data = self.redis.get("ComputerSystem:{}".format(
                    system_id["@odata.id"]))
                if system_data is None:
                    res = {
                        "Error":
                        "The System id {sys_id} is not available".format(
                            sys_id=system_id["@odata.id"].split('/')[-1])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return
                system_data = json.loads(json.loads(system_data))
                if not (system_data.get("SystemType") == "Composed"):
                    logging.error(
                        "The system {sys} provided in links is not a Composed system."
                        .format(sys=system_id["@odata.id"]))
                    res = {
                        "Error":
                        "The system {sys} provided in links is not a Composed system. Please provide composed systems for decompose"
                        .format(sys=system_id["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return
                logging.debug("ComposeSystem data: {sys_data}".format(
                    sys_data=system_data))

                for property, value in system_data["Links"].items():
                    if property == "ResourceBlocks":
                        for obj in value:
                            if obj.get("@odata.id"):
                                resource_data = self.redis.get(
                                    "{resource}:{resource_uri}".format(
                                        resource=property,
                                        resource_uri=obj["@odata.id"]))
                                if not resource_data:
                                    logging.error(
                                        "The Resource {rs_uri} is not found in db"
                                        .format(rs_uri=obj["@odata.id"]))
                                    continue
                                resource_data = json.loads(resource_data)
                                if resource_data and resource_data.get(
                                        "Links") and resource_data[
                                            "Links"].get("ComputerSystems"):
                                    done = False
                                    for sys_id in resource_data["Links"][
                                            "ComputerSystems"]:
                                        if sys_id["@odata.id"] == system_id[
                                                "@odata.id"]:
                                            logging.info(
                                                "Removing Composed system from Resource {uri}"
                                                .format(uri=obj["@odata.id"]))
                                            resource_data["Links"][
                                                "ComputerSystems"].remove(
                                                    sys_id)
                                            resource_data["Pool"] = "Free"
                                            resource_data["CompositionStatus"][
                                                "CompositionState"] = "Unused"
                                            if resource_data["CompositionStatus"][
                                                    "NumberOfCompositions"] > 0:
                                                resource_data[
                                                    "CompositionStatus"][
                                                        "NumberOfCompositions"] -= 1
                                            done = True
                                            break

                                    if done:
                                        pipe.set(
                                            "{resource}:{resource_uri}".format(
                                                resource=property,
                                                resource_uri=obj["@odata.id"]),
                                            json.dumps(resource_data))
                                        pipe.srem("ActivePool",
                                                  obj["@odata.id"])
                                        pipe.sadd("FreePool", obj["@odata.id"])
                                        logging.info(
                                            "{resource}:{resource_uri} is updateded"
                                            .format(
                                                resource=property,
                                                resource_uri=obj["@odata.id"]))
                                    else:
                                        logging.info(
                                            "{resource}:{resource_uri} updated is failed"
                                            .format(
                                                resource=property,
                                                resource_uri=obj["@odata.id"]))

                res["@odata.id"] = system_data["@odata.id"]
                res["@odata.type"] = system_data["@odata.type"]
                res["Id"] = system_data["Id"]
                res["Name"] = "Computer system decomposed"
                res["Links"] = {"ResourceBlocks": []}
                res["Links"]["ResourceBlocks"] = system_data["Links"][
                    "ResourceBlocks"]
                code = HTTPStatus.OK

                pipe.delete("ComputerSystem:{system_uri}".format(
                    system_uri=system_id["@odata.id"]))
                logging.info(
                    "ComputerSystem:{system_uri} is Decomposed Successfully".
                    format(system_uri=system_data))
                pipe.execute()

        except Exception as err:
            logging.error(
                "Unable to decompose the composed system. Error: {e}".format(
                    e=err))
            res = {
                "Error":
                "Unable to decompose the composed system. Error: {e}".format(
                    e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            pipe.reset()
            return res, code

    def get_composition_reservations_collection(self, url):

        res = {}
        code = HTTPStatus.OK
        try:

            res = {
                "@odata.type":
                "#CompositionReservationsCollection.CompositionReservationsCollection",
                "Name": "Composition Reservations Collection",
                "Members@odata.count": 0,
                "Members": [],
                "@odata.id": url
            }

            cr_data = self.redis.keys("CompositionReservations:*")
            if cr_data:
                for cr in cr_data:
                    res["Members"].append(
                        {"@odata.id": "{uri}".format(uri=cr)})

                res["Members@odata.count"] = len(cr_data)

            code = HTTPStatus.OK
            logging.debug(
                "CompositionReservations collection: {cr_collection}".format(
                    cr_collection=res))

        except Exception as err:
            logging.error(
                "Unable to get Composition Reservations Collection. Error: {e}"
                .format(e=err))
            res = {
                "error":
                "Unable to get Composition Reservations Collection. Error: {e}"
                .format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            return res, code
