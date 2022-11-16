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

from db.persistant import RedisClient
import logging


class TestResourcePool(unittest.TestCase):

    def test_update_free_pool(self):

        self.redis = RedisClient()

        self.redis.sadd(
            "FreePool", "/redfish/v1/CompositionService/ResourceBlocks/7cf65cda-a143-11ec-a8a0-be78894f3ea6")

    def test_update_active_pool(self):

        self.redis = RedisClient()

        
        self.redis.sadd(
            "ActivePool", "/redfish/v1/CompositionService/ResourceBlocks/7cf65cda-a143-11ec-a8a0-be78894f3ea6")
