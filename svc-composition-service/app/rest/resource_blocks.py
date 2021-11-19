from http import HTTPStatus
import logging
import requests
import json
import random
from db.persistant import RedisClient
from utilities.client import Client
from config import constants
from rest.resource_zones import ResourceZones
import copy
import uuid
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


class ResourceBlocks():

    def __init__(self):
        self.redis = RedisClient()
        self.client = Client()
        self.resourcezone = ResourceZones()

    def initialize(self):
        logging.info("Initialize Resource Blocks")

        try:
            # get systems collection data
            response = self.client.process_get_request(constants.SYSTEMS_URL)
            if response and response.get("Members"):
                for member in response["Members"]:
                    if member.get("@odata.id"):
                        # total_systems.append(member["@odata.id"])
                        system_uri = member["@odata.id"]
                        # get systems data from systems
                        system_res = self.client.process_get_request(
                            system_uri)
                        if system_res:
                            # create resource block and set data into database
                            resp = self.create_resource_block(system_res)
                            logging.info(
                                "successfully Created Resource Block. {resp}".format(resp=resp))
        except Exception as err:
            logging.error("unable to initialize the Resource Block")

    def create_resource_block(self, system_data=None):
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

            data['id'] = str(uuid.uuid1())
            data['@odata.id'] = "{url}/{id}".format(
                url=data['@odata.id'], id=data['id'])

            data['ComputerSystems'] = [
                {'@odata.id': system_data['@odata.id']}
            ]

            self.redis.set("{block}:{block_url}".format(
                block="ResourceBlocks", block_url=data['@odata.id']), str(json.dumps(data)))

            self.redis.set("{block_system}:{block_url}".format(
                block_system="ResourceBlocks-ComputerSystem", block_url=data['@odata.id']), system_data['@odata.id'])

            res = {"id": data['id']}
            logging.debug(
                "New ResourceBlock data: {rb_data}".format(rb_data=data))
        except Exception as err:
            logging.error(
                "Unable to create Resource Block. Error: {e}".format(e=err))
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
                "@odata.type": "#ResourceBlockCollection.ResourceBlockCollection",
                "Name": "Resource Block Collection",
                "Members@odata.count": 0,
                "Members": [],
                "@odata.id": url
            }

            rb_keys = self.redis.keys("ResourceBlocks:*")
            for rb_key in rb_keys:
                res["Members"].append({"@odata.id": "{uri}".format(
                    uri=rb_key.replace("ResourceBlocks:", ""))})

            res["Members@odata.count"] = len(rb_keys)
            code = HTTPStatus.OK
            logging.debug(
                "ResourceBlocks collection: {rb_collection}".format(rb_collection=res))

        except Exception as err:
            logging.error(
                "Unable to create Resource Block Collection. Error: {e}".format(e=err))
            res = {
                "error": "Unable to get Resource Block Collection. Error: {e}".format(e=err)
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
                "Error": "Unable to get Resource Block. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            return res, code
