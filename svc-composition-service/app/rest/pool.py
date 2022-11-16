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
from db.persistant import RedisClient
from http import HTTPStatus


class ResourcePool():
    def __init__(self):
        self.redis = RedisClient()

    def get_free_pool_collection(self, url):

        res = {}
        code = HTTPStatus.OK
        try:

            res = {
                "@odata.type": "#FreePoolCollection.FreePoolCollection",
                "Name": "Free Pool Collection",
                "Members@odata.count": 0,
                "Members": [],
                "@odata.id": url
            }

            fp_data = self.redis.smembers("FreePool")
            if fp_data:
                for fp in fp_data:
                    res["Members"].append(
                        {"@odata.id": "{uri}".format(uri=fp)})

                res["Members@odata.count"] = len(fp_data)

            code = HTTPStatus.OK
            logging.debug("FreePool collection: {fp_collection}".format(
                fp_collection=res))

        except Exception as err:
            logging.error(
                "Unable to get Free Pool Collection. Error: {e}".format(e=err))
            res = {
                "error":
                "Unable to get Free Pool Collection. Error: {e}".format(e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        return res, code

    def get_active_pool_collection(self, url):

        res = {}
        code = HTTPStatus.OK
        try:

            res = {
                "@odata.type": "#ActivePoolCollection.ActivePoolCollection",
                "Name": "Active Pool Collection",
                "Members@odata.count": 0,
                "Members": [],
                "@odata.id": url
            }

            ap_data = self.redis.smembers("ActivePool")
            if ap_data:
                for ap in ap_data:
                    res["Members"].append(
                        {"@odata.id": "{uri}".format(uri=ap)})

                res["Members@odata.count"] = len(ap_data)

            code = HTTPStatus.OK
            logging.debug("ActivePool collection: {ap_collection}".format(
                ap_collection=res))

        except Exception as err:
            logging.error(
                "Unable to get Active Pool Collection. Error: {e}".format(
                    e=err))
            res = {
                "error":
                "Unable to get Active Pool Collection. Error: {e}".format(
                    e=err)
            }
            code = HTTPStatus.INTERNAL_SERVER_ERROR
        return res, code
