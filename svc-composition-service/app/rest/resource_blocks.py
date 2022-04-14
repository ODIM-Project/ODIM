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

from db.persistant import RedisClient, RedisDb
from config import constants
from config.config import CONFIG_DATA
import copy
import uuid
from http import HTTPStatus


class ResourceBlocks():
    def __init__(self):
        self.redis = RedisClient()  # redis-ondisk
        self.redis_inmemory = RedisDb(
            CONFIG_DATA["RedisInMemoryAddress"])  # redis-inmemory

    """
    def initialize(self):
        logging.info("Initialize Resource Blocks")

        try:

            sys_urls = []
            rb_computer_keys = self.redis.keys(
                "ResourceBlocks-ComputerSystem:*")
            if rb_computer_keys:
                for key in rb_computer_keys:
                    sys_uri = self.redis.get(key)
                    if sys_uri:
                        sys_urls.append(sys_uri)
            # get systems collection data
            system_keys = self.redis_inmemory.keys("ComputerSystem:*")
            for system_key in system_keys:
                system_uri = system_key.replace("ComputerSystem:", "")
                if system_uri in sys_urls:
                    continue
                # get systems data from database
                system_res = self.redis_inmemory.get(system_key)
                if system_res:
                    system_res = json.loads(system_res)
                    # create resource block and set data into database
                    self.create_computer_sys_resource_block(system_res)
                    logging.info(
                        "successfully Created Computer system Resource Block.")

        except Exception as err:
            logging.error(
                "unable to initialize the Resource Block. Error: {e}".format(
                    e=err))
    """

    def create_computer_sys_resource_block(self, system_data=None):
        res = {}
        if system_data is None:
            return res
        pipe = self.redis.pipeline()
        try:
            if isinstance(system_data, str):
                system_data = json.loads(system_data)

            logging.info("Initialize for creation of new Resource Block")
            data = copy.deepcopy(constants.RESOURCE_BLOCK_TEMP)

            data['ResourceBlockType'].append("ComputerSystem")

            if system_data.get("Status") is not None:
                data['Status']['State'] = system_data['Status'].get('State')
                data['Status']['Health'] = system_data['Status'].get('Health')

            data['Id'] = str(uuid.uuid1())
            data['@odata.id'] = "{url}/{id}".format(url=data['@odata.id'],
                                                    id=data['Id'])

            data['ComputerSystems'] = [{'@odata.id': system_data['@odata.id']}]

            # if ResourceBlockType is ComputerSystem then created resource block becomes Active and Composed
            data["Pool"] = "Active"
            data["CompositionStatus"]["CompositionState"] = "Composed"
            data["CompositionStatus"]["NumberOfCompositions"] += 1

            data["Links"]= {"ComputerSystems": data['ComputerSystems']}

            pipe.set(
                "{block}:{block_url}".format(block="ResourceBlocks",
                                             block_url=data['@odata.id']),
                str(json.dumps(data)))

            pipe.set(
                "{block_system}:{block_url}".format(
                    block_system="ResourceBlocks-ComputerSystem",
                    block_url=data['@odata.id']), system_data['@odata.id'])

            pipe.sadd("ActivePool", data['@odata.id'])

            if system_data.get("Links"):
                if system_data["Links"].get("ResourceBlocks"):
                    system_data["Links"]["ResourceBlocks"].append(
                        {"@odata.id": data['@odata.id']})
                else:
                    system_data["Links"]["ResourceBlocks"] = [{
                        "@odata.id":
                        data['@odata.id']
                    }]
            else:
                system_data["Links"] = {
                    "ResourceBlocks": [{
                        "@odata.id": data['@odata.id']
                    }]
                }

            pipe.execute()
            self.redis_inmemory.set(
                "ComputerSystem:{sys_uri}".format(
                    sys_uri=system_data["@odata.id"]),
                json.dumps(json.dumps(system_data)))

            res = data
            logging.debug(
                "Successfully created Computer system resource Block. Data: {rb_data}"
                .format(rb_data=data))
        except Exception as err:
            logging.error(
                "Unable to create Computer system Resource Block. Error: {e}".
                format(e=err))
        pipe.reset()
        return res

    def get_resource_block_collection(self, url):

        res = {}
        code = HTTPStatus.OK
        if url is None:
            res["Error"] = "The Resource Block Collection url is empty"
            return res, HTTPStatus.NOT_FOUND
        try:

            res = {
                "@odata.type":
                "#ResourceBlockCollection.ResourceBlockCollection",
                "Name": "Resource Block Collection",
                "Members@odata.count": 0,
                "Members": [],
                "@odata.id": url
            }

            rb_keys = self.redis.keys("ResourceBlocks:*")
            for rb_key in rb_keys:
                res["Members"].append({
                    "@odata.id":
                    "{uri}".format(uri=rb_key.replace("ResourceBlocks:", ""))
                })

            res["Members@odata.count"] = len(rb_keys)
            code = HTTPStatus.OK
            logging.debug("ResourceBlocks collection: {rb_collection}".format(
                rb_collection=res))

        except Exception as err:
            logging.error(
                "Unable to create Resource Block Collection. Error: {e}".
                format(e=err))
            res = {
                "error":
                "Unable to get Resource Block Collection. Error: {e}".format(
                    e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        return res, code

    def get_resource_block(self, url):
        res = {}
        code = HTTPStatus.OK
        if url is None:
            return res, HTTPStatus.NOT_FOUND

        try:
            data = self.redis.get(
                "ResourceBlocks:{block_uri}".format(block_uri=url))
            if not data:
                res["Error"] = "The URI {uri} is not found".format(uri=url)
                code = HTTPStatus.NOT_FOUND
                return res, code
            res = json.loads(data)
            code = HTTPStatus.OK
            logging.debug("Get Resource Block: {rb_data}".format(rb_data=res))

        except Exception as err:
            logging.error(
                "Unable to get Resource Block. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to get Resource Block. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        return res, code

    def create_resource_block(self, request):
        res = {}
        code = HTTPStatus.OK
        try:
            logging.info("Initializing the creation of Resource Block")
            if request.get("ResourceBlockType") is None:
                res = {
                    "Error":
                    "The property 'ResourceBlockType' is missing from post body"
                }
                code = HTTPStatus.BAD_REQUEST
                return res, code

            if "ComputerSystem" in request["ResourceBlockType"]:
                sys_urls = []
                rb_computer_keys = self.redis.keys(
                    "ResourceBlocks-ComputerSystem:*")
                if rb_computer_keys:
                    for key in rb_computer_keys:
                        sys_uri = self.redis.get(key)
                        if sys_uri:
                            sys_urls.append(sys_uri)
                if request.get("ComputerSystems") is None:
                    res = {
                        "Error":
                        "The property 'ComputerSystems' is missing from post body"
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code
                elif not len(request["ComputerSystems"]):
                    res = {"Error": "The property 'ComputerSystems' is empty"}
                    code = HTTPStatus.BAD_REQUEST
                    return res, code
                elif len(request["ComputerSystems"]) > 1:
                    res = {
                        "Error":
                        "Request body has more than one computer system, Resource Block will be created with only one computer system"
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code

                if request["ComputerSystems"][0]["@odata.id"] in sys_urls:
                    res = {
                        "Error":
                        "The ComputerSystem {sys} is aready exist".format(
                            sys=request["ComputerSystems"][0]["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code

                system_res = self.redis_inmemory.get(
                    "ComputerSystem:{sys_uri}".format(
                        sys_uri=request["ComputerSystems"][0]["@odata.id"]))
                if system_res:
                    system_res = json.loads(system_res)
                    # create resource block and set data into database
                    res = self.create_computer_sys_resource_block(system_res)
                    logging.info(
                        "successfully Created Computer system Resource Block.")
                    code = HTTPStatus.CREATED
                    return res, code
                else:
                    res = {
                        "Error":
                        "The System {uri} is not found".format(
                            uri=request["ComputerSystems"][0]["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return res, code
        except Exception as err:
            logging.error(
                "Unable to Create Resource Block. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to Create Resource Block. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        return res, code

    def delete_resource_block(self, url):
        res = {}
        code = HTTPStatus.OK
        if url is None:
            return res, HTTPStatus.NOT_FOUND
        pipe = self.redis.pipeline()
        inmemory_pipe = self.redis_inmemory.pipeline()

        try:
            data = self.redis.get(
                "ResourceBlocks:{block_uri}".format(block_uri=url))
            if not data:
                res["Error"] = "The URI {uri} is not found".format(uri=url)
                code = HTTPStatus.NOT_FOUND
                return res, code
            rb_data = json.loads(data)
            if "ComputerSystem" not in rb_data["ResourceBlockType"]:
                if rb_data["CompositionStatus"][
                        "CompositionState"] != "Unused" and rb_data[""] != "Free":
                    res = {
                        "Error":
                        "The resouce block {rb_id} deletion is failed because the resource block is in 'active/composed' state"
                        .format(rb_id=rb_data["Id"])
                    }
                    code = HTTPStatus.CONFLICT
                    return res, code

            system_link = self.redis.get("{rb_cs}:{block_url}".format(
                rb_cs="ResourceBlocks-ComputerSystem", block_url=url))

            if system_link:
                system_key = "ComputerSystem:{sys_uri}".format(
                    sys_uri=system_link)
                system_res = self.redis_inmemory.get(system_key)
                if system_res:
                    system_res = json.loads(system_res)
                    if isinstance(system_res, str):
                        system_res = json.loads(system_res)
                    if system_res.get("Links") and system_res["Links"].get(
                            "ResourceBlocks"):
                        for rb in system_res["Links"]["ResourceBlocks"]:
                            if rb.get("@odata.id") == url:
                                system_res["Links"]["ResourceBlocks"].remove(
                                    rb)
                                if not system_res["Links"]["ResourceBlocks"]:
                                    del system_res["Links"]["ResourceBlocks"]
                                inmemory_pipe.set(
                                    system_key,
                                    json.dumps(json.dumps(system_res)))
                                break

            if rb_data.get("Links") and rb_data["Links"].get("Zones"):
                for zone_uri in rb_data["Links"]["Zones"]:
                    if zone_uri.get("@odata.id"):
                        zone_key = "ResourceZones:{zuri}".format(
                            zuri=zone_uri["@odata.id"])
                        zone_block_key = "{zone_block}:{zone_uri}".format(
                            zone_block="ResourceZone-ResourceBlock",
                            zone_uri=zone_uri["@odata.id"])
                        zone_block_link = self.redis.smembers(zone_block_key)
                        if len(zone_block_link) <= 1:
                            # if resource block linked zone is less than 1 then we can remove zone
                            pipe.delete(zone_block_key, zone_key)
                        else:
                            pipe.srem(zone_block_key, rb_data["@odata.id"])
                            # delete resource block link in ResourceZones data
                            zone_data = self.redis.get(zone_key)
                            if zone_data:
                                zone_data = json.loads(zone_data)
                                if zone_data.get("Links") and zone_data[
                                        "Links"].get("ResourceBlocks"):
                                    for rb_uri in zone_data["Links"][
                                            "ResourceBlocks"]:
                                        if rb_uri.get("@odata.id") and rb_uri[
                                                "@odata.id"] == rb_data[
                                                    "@odata.id"]:
                                            zone_data["Links"][
                                                "ResourceBlocks"].remove(
                                                    rb_uri)
                                            if not zone_data["Links"][
                                                    "ResourceBlocks"]:
                                                del zone_data["Links"][
                                                    "ResourceBlocks"]
                                            pipe.set(
                                                zone_key,
                                                str(json.dumps(zone_data)))
                                            break

            if rb_data["Pool"] == "Free":
                pipe.srem("FreePool", rb_data["@odata.id"])
            elif rb_data["Pool"] == "Active":
                pipe.srem("ActivePool", rb_data["@odata.id"])

            pipe.delete("ResourceBlocks:{block_uri}".format(
                block_uri=rb_data["@odata.id"]))
            pipe.delete("ResourceBlocks-ComputerSystem:{block_uri}".format(
                block_uri=rb_data["@odata.id"]))

            pipe.execute()
            inmemory_pipe.execute()
            logging.info(
                "The Resource Block {rb_uri} is deleted successfully".format(
                    rb_uri=rb_data["@odata.id"]))
            code = HTTPStatus.NO_CONTENT

        except Exception as err:
            logging.error(
                "Unable to get Resource Block. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to get Resource Block. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        pipe.reset()
        inmemory_pipe.reset()
        return res, code
