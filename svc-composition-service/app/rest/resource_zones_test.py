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

import unittest
import json

from db.persistant import RedisClient


class TestResourcePool(unittest.TestCase):

    def test_update_resource_block(self):

        self.redis = RedisClient()

        zone_uri = "/redfish/v1/CompositionService/ResourceZones/d1bf2f54-6c8b-11ec-b071-ce46db785eee"

        request_body = {
            "@odata.id": zone_uri
        }

        self.redis.set("ResourceZones:{url}".format(
            url=zone_uri), str(json.dumps(request_body)))
