import logging
import json
from db.persistant import RedisClient
from utilities.client import Client
from config import constants
import copy
import uuid
from http import HTTPStatus
import urllib3

urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


class ResourceZones():
    def __init__(self):
        self.redis = RedisClient()
        self.client = Client()

    def create_resource_zone(self, request_data=None):

        res = {}
        code = HTTPStatus.CREATED
        pipe = self.redis.pipeline()

        if request_data is None:
            res["Error"] = "Request Body is empty"
            return res, HTTPStatus.BAD_REQUEST

        try:
            logging.info("Initialize for creation of Resource Zone")

            # request_data = {
            #   "Name": "CS Reosurce Zone",
            #   "Links":{
            #       "ResourceBlocks": [
            #       {
            #           "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/1"
            #       },
            #       {
            #           "@odata.id": "/redfish/v1/CompositionService/ResourceBlocks/2"
            #       },
            #       ]
            #    }
            # }
            # Required parameters Name, Links, and ResourceBlocks.
            logging.debug(
                "Request payload is {body}".format(body=request_data))
            zone = copy.deepcopy(constants.RESOURCE_ZONE_TEMP)
            zone['Id'] = str(uuid.uuid1())
            zone['@odata.id'] = "{zone_uri}/{id}".format(
                zone_uri=zone['@odata.id'], id=zone['Id'])
            if request_data.get('Name') is None or request_data['Name'] == "":
                # logging.error("The Property Name is missing")
                res = {"Error": "The Property 'Name' is missing."}
                code = HTTPStatus.BAD_REQUEST
                return

            zone['Name'] = request_data['Name']

            if request_data.get('Links') is None:
                res = {"Error": "The Property 'Links' is missing."}
                code = HTTPStatus.BAD_REQUEST
                return
            if request_data['Links'].get('ResourceBlocks') is None:
                res = {
                    "Error": "The Property 'Links.ResourceBlocks' is missing."
                }
                code = HTTPStatus.BAD_REQUEST
                return

            resource_block_list = request_data['Links']['ResourceBlocks']
            for resource_block in resource_block_list:
                zone["Links"]["ResourceBlocks"].append(
                    {"@odata.id": resource_block['@odata.id']})
                data = self.redis.get("ResourceBlocks:{block_uri}".format(
                    block_uri=resource_block['@odata.id']))
                if data is not None:
                    data = json.loads(data)
                    zone['Status']['State'] = data['Status'].get('State')
                    zone['Status']['Health'] = data['Status'].get('Health')
                    if data.get('Links') is not None:
                        if data['Links'].get('Zones') is None:
                            data['Links']['Zones'] = []
                    else:
                        data['Links'] = {"Zones": []}

                    data['Links']['Zones'].append(
                        {"@odata.id": zone['@odata.id']})

                    pipe.set(
                        "{block}:{block_url}".format(
                            block="ResourceBlocks",
                            block_url=data['@odata.id']),
                        str(json.dumps(data)))
                    logging.info(
                        "Resource Block linked to Resource Zone is successfully updated"
                    )

                else:
                    logging.debug(
                        "Getting resource block data from redis is failed for this resource: {uri}"
                        .format(uri=resource_block['@odata.id']))
                    res = {
                        "Error":
                        "The Resource Block {rs_block} is not found. Create ResourceZone is failed"
                        .format(rs_block=resource_block['@odata.id'])
                    }
                    code = HTTPStatus.BAD_REQUEST
                    return

            pipe.set(
                "{zones}:{zone_uri}".format(zones="ResourceZones",
                                            zone_uri=zone['@odata.id']),
                str(json.dumps(zone)))
            pipe.set(
                "{zone_block}:{zone_url}".format(
                    zone_block="ResourceZone-ResourceBlock",
                    zone_url=zone['@odata.id']), resource_block['@odata.id'])

            pipe.execute()

            logging.info("Created a Resource Zone successfully")
            res = {"Id": zone["Id"]}
            code = HTTPStatus.CREATED

        except Exception as err:
            logging.error(
                "Unable to create Resource Zone. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to create Resource Zone. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            pipe.reset()
            return res, code

    def get_resource_zone_collection(self, url):

        res = {}
        code = HTTPStatus.OK
        if url is None:
            res["Error"] = "The Resource Zone Collection url is empty"
            return res, HTTPStatus.NOT_FOUND

        try:
            res = {
                "@odata.type": "#ZoneCollection.ZoneCollection",
                "Name": "Resource Zone Collection",
                "Members@odata.count": 0,
                "Members": [],
                "@odata.id": url
            }

            rz_keys = self.redis.keys("ResourceZones:*")

            for rz_key in rz_keys:
                res["Members"].append({
                    "@odata.id":
                    "{uri}".format(uri=rz_key.replace("ResourceZones:", ""))
                })

            res["Members@odata.count"] = len(rz_keys)

            code = HTTPStatus.OK

        except Exception as err:
            logging.error(
                "Unable to Get Resource Zone Collection. Error: {e}".format(
                    e=err))
            res = {
                "Error":
                "Unable to Get Resource Zone Collection. Error: {e}".format(
                    e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            return res, code

    def get_resource_zone(self, url):

        res = {}
        code = HTTPStatus.OK
        if url is None:
            res["Error"] = "The Resource Zone Collection url is empty"
            return res, HTTPStatus.NOT_FOUND

        try:
            data = self.redis.get(
                "ResourceZones:{zone_uri}".format(zone_uri=url))
            if not data:
                res["Error"] = "The URI {uri} is not found".format(uri=url)
                code = HTTPStatus.NOT_FOUND
                return
            res = json.loads(str(data))
            code = HTTPStatus.OK

        except Exception as err:
            logging.error(
                "Unable to Get Resource Zone. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to Get Resource Zone. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            return res, code

    def delete_resource_zone(self, url):

        res = {}
        code = HTTPStatus.NO_CONTENT
        pipe = self.redis.pipeline()

        try:
            logging.info("Initialize delete Resource Zone")
            data = self.redis.get(
                "ResourceZones:{zone_uri}".format(zone_uri=url))
            if not data:
                res["Error"] = "The URI {uri} is not found".format(uri=url)
                code = HTTPStatus.NOT_FOUND
                return

            data = json.loads(data)
            logging.debug("ResourceZone data: {data}".format(data=data))

            for property, value in data['Links'].items():
                if type(value) is list:
                    for obj in value:
                        if obj.get("@odata.id"):
                            resource_data = self.redis.get(
                                "{resource}:{resource_uri}".format(
                                    resource=property,
                                    resource_uri=obj["@odata.id"]))
                            if not resource_data:
                                logging.error(
                                    "Unable to get {rb_uri} redis db".format(
                                        rb_uri=obj["@odata.id"]))
                                continue
                            resource_data = json.loads(resource_data)
                            if resource_data and resource_data.get(
                                    "Links") and resource_data["Links"].get(
                                        "Zones"):
                                done = False
                                for zone_id in resource_data["Links"]["Zones"]:
                                    if zone_id["@odata.id"] == data[
                                            "@odata.id"]:
                                        logging.info(
                                            "Updating the ResourceBlock {rb_uri}"
                                            .format(rb_uri=obj["@odata.id"]))
                                        resource_data["Links"]["Zones"].remove(
                                            zone_id)
                                        done = True
                                        break

                                if done:
                                    pipe.set(
                                        "{resource}:{resource_uri}".format(
                                            resource=property,
                                            resource_uri=obj["@odata.id"]),
                                        json.dumps(resource_data))
                                    logging.info(
                                        "{resource}:{resource_uri} is updateded"
                                        .format(resource=property,
                                                resource_uri=obj["@odata.id"]))
                                else:
                                    logging.info(
                                        "{resource}:{resource_uri} updated is failed"
                                        .format(resource=property,
                                                resource_uri=obj["@odata.id"]))

            pipe.delete("ResourceZones:{zone_uri}".format(zone_uri=url))
            pipe.delete("{zone_block}:{zone_url}".format(
                zone_block="ResourceZone-ResourceBlock",
                zone_url=data['@odata.id']))
            logging.info("{resource}:{resource_uri} is deleted".format(
                resource="ResourceZones", resource_uri=url))

            pipe.execute()
            res = {"Id": data["Id"]}
            code = HTTPStatus.OK

        except Exception as err:
            logging.error(
                "Unable to delete the Resource Zone. Error: {e}".format(e=err))
            res = {
                "Error":
                "Unable to delete the Resource Zone. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        finally:
            pipe.reset()
            return res, code
