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
import logging
from http import HTTPStatus
from rest.composition_service import CompositonService
from rest.resource_blocks import ResourceBlocks
from rest.pool import ResourcePool
from rest.resource_zones import ResourceZones
import proto.compositionservice.composition_service_pb2 as pb2
from config import constants


class TestCompositionService(unittest.TestCase):

    def mock_is_authorized(self, session_token, privileges, oem_privileges):
        resp = {}
        if session_token != "validToken":
            resp["status_code"] = HTTPStatus.UNAUTHORIZED
            resp["status_message"] = {
                "error": "error while trying to authenticate session"
            }
        else:
            resp["status_code"] = HTTPStatus.OK
            resp["status_message"] = {
                "error": "success"
            }

        return resp

    def test_get_composition_service(self):

        self.cs = CompositonService()

        req_body = [
            {
                "req": pb2.GetCompositionServiceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService"),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionServiceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService"),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            }
        ]

        for request in req_body:
            auth_resp = self.mock_is_authorized(
                request["req"].SessionToken, [constants.PrivilegeLogin], [""])
            if auth_resp["status_code"] == HTTPStatus.OK:
                response, code = self.cs.get_cs()
            else:
                response = auth_resp["status_message"]
                code = auth_resp["status_code"]
                logging.error("error while trying to get error")

            self.assertEqual(
                code, request["resp"].StatusCode, msg=str(response))

    def test_get_composition_resource(self):

        self.resourcezone = ResourceZones()
        self.resourceblock = ResourceBlocks()
        self.pool = ResourcePool()

        req_body = [
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService/ResourceBlocks", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService/ResourceBlocks", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService/ResourceBlocks/7cf65cda-a143-11ec-a8a0-be78894f3ea6", ResourceID="7cf65cda-a143-11ec-a8a0-be78894f3ea6"),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService/ResourceBlocks/7cf65cda-a143-11ec-a8a0-be78894f3ea6", ResourceID="7cf65cda-a143-11ec-a8a0-be78894f3ea6"),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService/ResourceZones", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService/ResourceZones", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService/ResourceZones/d1bf2f54-6c8b-11ec-b071-ce46db785eee", ResourceID="d1bf2f54-6c8b-11ec-b071-ce46db785eee"),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService/ResourceZones/d1bf2f54-6c8b-11ec-b071-ce46db785eee", ResourceID="d1bf2f54-6c8b-11ec-b071-ce46db785eee"),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService/ActivePool", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService/ActivePool", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="validToken", URL="/redfish/v1/CompositionService/FreePool", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.OK)
            },
            {
                "req": pb2.GetCompositionResourceRequest(SessionToken="InvalidToken", URL="/redfish/v1/CompositionService/FreePool", ResourceID=''),
                "resp": pb2.CompositionServiceResponse(StatusCode=HTTPStatus.UNAUTHORIZED)
            }
        ]

        for request in req_body:
            auth_resp = self.mock_is_authorized(
                request["req"].SessionToken, [constants.PrivilegeLogin], [""])

            if auth_resp["status_code"] == HTTPStatus.OK:
                if request["req"].URL:
                    segments = request["req"].URL.split("/")
                    if request["req"].ResourceID:
                        # ResourceBlocks Instance
                        if segments[-2] == "ResourceBlocks":
                            response, code = self.resourceblock.get_resource_block(
                                request["req"].URL)
                        # ResourceZones Instance
                        elif segments[-2] == "ResourceZones":
                            response, code = self.resourcezone.get_resource_zone(
                                request["req"].URL)
                    else:
                        # ResourceBlocks Collection
                        if segments[-1] == "ResourceBlocks":
                            response, code = self.resourceblock.get_resource_block_collection(
                                request["req"].URL)
                        # ResourceZones Collection
                        elif segments[-1] == "ResourceZones":
                            response, code = self.resourcezone.get_resource_zone_collection(
                                request["req"].URL)
                        # ActivePool Collection
                        elif segments[-1] == "ActivePool":
                            response, code = self.pool.get_active_pool_collection(
                                request["req"].URL)
                        # FreePool Collection
                        elif segments[-1] == "FreePool":
                            response, code = self.pool.get_free_pool_collection(
                                request["req"].URL)
            else:
                response = auth_resp["status_message"]
                code = auth_resp["status_code"]

            self.assertEqual(
                code, request["resp"].StatusCode, msg=str(response))


if __name__ == '__main__':

    unittest.main()
