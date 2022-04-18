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
from http import HTTPStatus
from db.persistant import RedisClient, RedisDb
from config.config import CONFIG_DATA


class CompositonService():
    def __init__(self):
        self.redis = RedisClient()  # redis-ondisk
        self.redis_inmemory = RedisDb(
            CONFIG_DATA["RedisInMemoryAddress"])  # redis-inmemory

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
                    compose_res, code = self.create_compose_system(
                        stanza["Request"])
                    stanza["Response"] = compose_res
                elif stanza["StanzaType"] == "DecomposeSystem":
                    decompose_res, code = self.decompose_system(
                        stanza["Request"])
                    stanza["Response"] = decompose_res

            if code != HTTPStatus.OK:
                res = stanzas[0]["Response"]
            else:
                res = req
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
        return res, code

    def create_compose_system(self, req):
        res = {}
        code = HTTPStatus.OK
        system_data = {}
        pipe = self.redis.pipeline()
        inmemory_pipe = self.redis_inmemory.pipeline()
        system_uri = None
        resource_block_uri = None
        try:
            logging.info("Initialize create compose system")

            if not (req.get("Links") and req["Links"].get("ResourceBlocks")
                    and any(req["Links"]["ResourceBlocks"])):
                res = {"Error": "Unable to find ResourceBlocks Links"}
                code = HTTPStatus.BAD_REQUEST
                return res, code

            logging.debug(
                "Create Compose system request body is: {req}".format(req=req))

            rb_list = []

            for block_uri in req["Links"]["ResourceBlocks"]:
                if block_uri.get("@odata.id"):
                    rb_list.append(block_uri["@odata.id"])

            if len(rb_list) <= 1:
                res = {
                    "Error":
                    "Compose system accepts two or more resource blocks. Resubmit the request with two or more resource Blocks"
                }
                code = HTTPStatus.BAD_REQUEST
                return res, code

            # find the computer system from resource blocks

            for rb_uri in rb_list:
                system_uri = self.redis.get("{rb_cs}:{rb}".format(
                    rb_cs="ResourceBlocks-ComputerSystem", rb=rb_uri))
                if system_uri:
                    resource_block_uri = rb_uri
                    break
            if not system_uri:
                res = {
                    "Error":
                    "Atleast one ComputerSystem Type Resource Block must present. Resubmit with valid request"
                }
                code = HTTPStatus.BAD_REQUEST
                return res, code
            system_key = "ComputerSystem:{sys_uri}".format(sys_uri=system_uri)
            system_data = self.redis_inmemory.get(system_key)
            if system_data is not None:
                system_data = json.loads(system_data)
                if isinstance(system_data, str):
                    system_data = json.loads(system_data)
            else:
                res = {
                    "Error":
                    "The system {sys_uri} which is linked to resource block {ruri} is not found"
                    .format(sys_uri=system_uri, ruri=resource_block_uri)
                }
                code = HTTPStatus.BAD_REQUEST
                return res, code

            for rb_uri in rb_list:
                if resource_block_uri == rb_uri:
                    continue
                rb_data = self.redis.get(
                    "ResourceBlocks:{block_uri}".format(block_uri=rb_uri))
                if not rb_data:
                    res = {
                        "Error":
                        "The Resource at the URI {uri} is not found".format(
                            uri=rb_uri)
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code
                rb_data = json.loads(rb_data)
                logging.debug("ResourceBlock {id} data is: {data}".format(
                    id=rb_uri, data=rb_data))
                if "ComputerSystem" in rb_data["ResourceBlockType"]:
                    continue

                if rb_data["CompositionStatus"]["MaxCompositions"] <= rb_data[
                        "CompositionStatus"]["NumberOfCompositions"]:
                    res = {
                        "Error":
                        "NumberOfCompositions are excided to MaxCompositions for the Resource Block {uri}"
                        .format(uri=rb_uri)
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code

                if (rb_data["Pool"] != "Free") or (
                        rb_data["CompositionStatus"]["CompositionState"] !=
                        "Unused"):
                    res = {
                        "Error":
                        "The Resource Block {uri} is already used by other composed system"
                        .format(uri=rb_uri)
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code

                rb_data["Pool"] = "Active"
                rb_data["CompositionStatus"]["CompositionState"] = "Composed"
                rb_data["CompositionStatus"]["NumberOfCompositions"] += 1

                if not rb_data.get("Links"):
                    rb_data["Links"] = {"ComputerSystems": []}
                elif not rb_data["Links"].get("ComputerSystems"):
                    rb_data["Links"]["ComputerSystems"] = []
                rb_data["Links"]["ComputerSystems"].append(
                    {"@odata.id": system_uri})

                pipe.set("ResourceBlocks:{rb_uri}".format(rb_uri=rb_uri),
                         json.dumps(rb_data))
                pipe.srem("FreePool", rb_uri)
                pipe.sadd("ActivePool", rb_uri)

                if {"@odata.id": rb_uri} not in system_data["Links"]["ResourceBlocks"]:
                    system_data["Links"]["ResourceBlocks"].append(
                        {"@odata.id": rb_uri})

            inmemory_pipe.set(system_key, json.dumps(json.dumps(system_data)))

            res["@odata.id"] = system_data["@odata.id"]
            res["@odata.type"] = system_data["@odata.type"]
            res["Id"] = system_data["Id"]
            res["Name"] = system_data["Name"]
            res["SystemType"] = system_data["SystemType"]
            res["Links"] = system_data["Links"]

            pipe.execute()
            inmemory_pipe.execute()
            code = HTTPStatus.OK
            logging.info("Successfully composed system")
        except Exception as err:
            logging.error(
                "Unable to create composed system. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to create composed system. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
            
        pipe.reset()
        inmemory_pipe.reset()
        return res, code

    def decompose_system(self, req):
        res = {}
        code = HTTPStatus.OK
        pipe = self.redis.pipeline()
        inmemory_pipe = self.redis_inmemory.pipeline()
        system_data = {}

        try:
            logging.info("Initialize Decompose System")

            if not (req.get("Links") and req["Links"].get("ComputerSystems")):
                res = {
                    "Error":
                    "ComputerSystems Links is empty. Provide the ComputerSystem link and resubmit the request"
                }
                code = HTTPStatus.BAD_REQUEST
                return res, code

            logging.debug(
                "DecomposeSystem request body: {req}".format(req=req))

            for system_id in req["Links"]["ComputerSystems"]:
                system_key = "ComputerSystem:{sys_uri}".format(
                    sys_uri=system_id["@odata.id"])
                system_data = self.redis_inmemory.get(system_key)
                if system_data is None:
                    res = {
                        "Error":
                        "The Resource at the URI {sys_id} is not found".format(
                            sys_id=system_id["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code
                system_data = json.loads(system_data)
                if isinstance(system_data, str):
                    system_data = json.loads(system_data)

                if not (system_data.get("Links")
                        and system_data["Links"].get("ResourceBlocks")
                        and any(system_data["Links"]["ResourceBlocks"])):
                    res = {
                        "Error":
                        "Decompose System failed because No Resource Blocks linked to system"
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code
                system_rb_links = copy.deepcopy(
                    system_data["Links"]["ResourceBlocks"])
                for rb_uri in system_data["Links"]["ResourceBlocks"]:
                    if not rb_uri.get("@odata.id"):
                        continue
                    rb_key = "{resource}:{resource_uri}".format(
                        resource="ResourceBlocks",
                        resource_uri=rb_uri["@odata.id"])
                    rb_data = self.redis.get(rb_key)
                    if rb_data is None:
                        res = {
                            "Error":
                            "The Resource at the URI {uri} is not found".
                            format(uri=rb_uri["@odata.id"])
                        }
                        code = HTTPStatus.BAD_REQUEST
                        return res, code
                    rb_data = json.loads(rb_data)
                    if "ComputerSystem" in rb_data["ResourceBlockType"]:
                        continue

                    system_uri = self.redis.get("{rb_cs}:{rb}".format(
                        rb_cs="ResourceBlocks-ComputerSystem",
                        rb=rb_uri["@odata.id"]))
                    if system_uri != system_id["@odata.id"]:
                        system_rb_links.remove(rb_uri)

                    rb_data["Pool"] = "Free"
                    rb_data["CompositionStatus"]["CompositionState"] = "Unused"

                    if rb_data["CompositionStatus"]["NumberOfCompositions"] > 0:
                        rb_data["CompositionStatus"][
                            "NumberOfCompositions"] -= 1

                    if rb_data.get("Links") and rb_data["Links"].get(
                            "ComputerSystems"):
                        if system_id in rb_data["Links"]["ComputerSystems"]:
                            rb_data["Links"]["ComputerSystems"].remove(
                                system_id)
                            if not rb_data["Links"]["ComputerSystems"]:
                                del rb_data["Links"]["ComputerSystems"]

                    pipe.set(rb_key, json.dumps(rb_data))

                if system_rb_links:
                    system_data["Links"]["ResourceBlocks"] = system_rb_links
                else:
                    del system_data["Links"]["ResourceBlocks"]
                inmemory_pipe.set(system_key,
                                  json.dumps(json.dumps(system_data)))

            res["@odata.id"] = system_data["@odata.id"]
            res["@odata.type"] = system_data["@odata.type"]
            res["Id"] = system_data["Id"]
            res["Name"] = system_data["Name"]
            res["SystemType"] = system_data["SystemType"]
            res["Links"] = system_data["Links"]
            pipe.execute()
            inmemory_pipe.execute()
            logging.info("Successfully Decomposed system")
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
        
        pipe.reset()
        inmemory_pipe.reset()
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
        return res, code
