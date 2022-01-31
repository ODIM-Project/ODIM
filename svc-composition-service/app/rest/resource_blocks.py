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
from db.persistant import RedisClient
from utilities.client import Client
from config import constants
from rest.resource_zones import ResourceZones
import copy
import uuid
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
from http import HTTPStatus


class ResourceBlocks():
    def __init__(self):
        self.redis = RedisClient()
        self.client = Client()
        self.resourcezone = ResourceZones()

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
            response = self.client.process_get_request(constants.SYSTEMS_URL)
            if response and response.get("Members"):
                for member in response["Members"]:
                    if member.get("@odata.id"):
                        # total_systems.append(member["@odata.id"])
                        system_uri = member["@odata.id"]
                        if system_uri in sys_urls:
                            continue
                        # get systems data from systems
                        system_res = self.client.process_get_request(
                            system_uri)
                        if system_res:
                            # create resource block and set data into database
                            resp = self.create_computer_sys_resource_block(
                                system_res)
                            logging.info(
                                "successfully Created Computer system Resource Block. {resp}"
                                .format(resp=resp))
        except Exception as err:
            logging.error(
                "unable to initialize the Resource Block. Error: {e}".format(
                    e=err))

    def create_computer_sys_resource_block(self, system_data=None):
        res = {}
        if system_data is None:
            return res
        try:
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

            self.redis.set(
                "{block}:{block_url}".format(block="ResourceBlocks",
                                             block_url=data['@odata.id']),
                str(json.dumps(data)))

            self.redis.set(
                "{block_system}:{block_url}".format(
                    block_system="ResourceBlocks-ComputerSystem",
                    block_url=data['@odata.id']), system_data['@odata.id'])

            self.redis.sadd("FreePool", data['@odata.id'])

            res = json.dumps(data)
            logging.debug(
                "New ResourceBlock data: {rb_data}".format(rb_data=data))
        except Exception as err:
            logging.error(
                "Unable to create Computer system Resource Block. Error: {e}".
                format(e=err))
        finally:
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
                "@odata.id": url,
                "Oem": {
                    "Ami": {
                        "Actions": {
                            "#ResourceBlock.Initialize": {
                                "target":
                                "/redfish/v1/CompositionService/ResourceBlocks/Actions/Oem/Ami/ResourceBlock.Initialize"
                            }
                        }
                    }
                }
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
        finally:
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
                return
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
        finally:
            return res, code

    def create_resource_block(self, request):
        res = {}
        code = HTTPStatus.OK
        try:
            logging.info("Initialising the creation of Resource Block")
            if request.get("ResourceBlockType") is None:
                logging.error(
                    "The property 'ResourceBlockType' is missing from post body"
                )
                res = {
                    "Error":
                    "The property 'ResourceBlockType' is missing from post body"
                }
                code = HTTPStatus.BAD_REQUEST
                return

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
                    logging.error(
                        "The property 'ComputerSystems' is missing from post body"
                    )
                    res = {
                        "Error":
                        "The property 'ComputerSystems' is missing from post body"
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return
                elif not len(request["ComputerSystems"]):
                    logging.error("The property 'ComputerSystems' is empty")
                    res = {"Error": "The property 'ComputerSystems' is empty"}
                    code = HTTPStatus.BAD_REQUEST
                    return
                elif len(request["ComputerSystems"]) > 1:
                    logging.debug(
                        "Request has more than one computer system, Resource Block will be created with only one comuter system"
                    )

                if request["ComputerSystems"][0]["@odata.id"] in sys_urls:
                    logging.error(
                        "The ComputerSystem {sys} is aready exist".format(
                            sys=request["ComputerSystems"][0]["@odata.id"]))
                    res = {
                        "Error":
                        "The ComputerSystem {sys} is aready exist".format(
                            sys=request["ComputerSystems"][0]["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return

                system_res = self.client.process_get_request(
                    request["ComputerSystems"][0]["@odata.id"])
                if system_res:
                    # create resource block and set data into database
                    res = self.create_computer_sys_resource_block(system_res)
                    logging.info(
                        "successfully Created Computer system Resource Block. {resp}"
                        .format(resp=res))
                    code = HTTPStatus.OK
                    return
                else:
                    logging.error("The System {uri} is not found".format(
                        uri=request["ComputerSystems"][0]["@odata.id"]))
                    res = {
                        "Error":
                        "The System {uri} is not found".format(
                            uri=request["ComputerSystems"][0]["@odata.id"])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return
        except Exception as err:
            logging.error(
                "Unable to Create Resource Block. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to Create Resource Block. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            return res, code

    def delete_resource_block(self, url):
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
                return
            rb_data = json.loads(data)

            if rb_data["CompositionStatus"]["CompositionState"] == "Composed":
                logging.error(
                    "The resource block {rb_id} delete is failed. The resource block is comoposed with system {sys_id}"
                    .format(rb_id=rb_data["Id"],
                            sys_id=rb_data["Links"]["ComputerSystems"][0]
                            ["@odata.id"]))
                res = {
                    "Error":
                    "The resouce block {rb_id} delete is failed because the resource is in 'composed' state"
                    .format(rb_id=rb_data["Id"])
                }
                code = HTTPStatus.CONFLICT
                return

            if rb_data.get("Links") and rb_data["Links"].get("Zones"):
                if len(rb_data["Links"]["Zones"]):
                    logging.error(
                        "The resource block {rb_id} delete is failed. The resource block is linked with resource zone"
                        .format(rb_id=rb_data["Id"]))
                    res = {
                        "Error":
                        "The resource block {rb_id} delete is failed because the resource is liked with resource zone"
                        .format(rb_id=rb_data["Id"])
                    }
                    code = HTTPStatus.CONFLICT
                    return

            if rb_data["Pool"] == "Active":
                self.redis.srem("ActivePool", rb_data["@odata.id"])
            elif rb_data["Pool"] == "Free":
                self.redis.srem("FreePool", rb_data["@odata.id"])

            self.redis.delete("ResourceBlocks:{block_uri}".format(
                block_uri=rb_data["@odata.id"]))
            self.redis.delete(
                "ResourceBlocks-ComputerSystem:{block_uri}".format(
                    block_uri=rb_data["@odata.id"]))
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
        finally:
            return res, code
