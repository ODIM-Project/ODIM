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

import json
import grpc
from http import HTTPStatus
import proto.composition_service_pb2 as pb2
import proto.composition_service_pb2_grpc as pb2_grpc
from rest.resource_zones import ResourceZones
from rest.resource_blocks import ResourceBlocks
from rest.pool import ResourcePool
from rest.composition_service import CompositonService
from config import constants
from utilities.auth import Auth


class CompositionServiceRpc(pb2_grpc.CompositionServicer):
    """The listener function implements the rpc call as described in the .proto file"""
    def __init__(self):
        self.cs = CompositonService()
        self.resourcezone = ResourceZones()
        self.resourceblock = ResourceBlocks()
        self.pool = ResourcePool()
        self.auth = Auth()
        self.headers = {
            "Allow": '"GET"',
            "Content-type": "application/json; charset=utf-8"
        }

    def __str__(self):
        return self.__class__.__name__

    def GetCompositionService(
            self, request: pb2.GetCompositionServiceRequest,
            context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        auth_resp = self.auth.is_authorized_rpc(request.SessionToken,
                                                [constants.PrivilegeLogin],
                                                [""])
        if auth_resp["status_code"] == HTTPStatus.OK:
            response, code = self.cs.get_cs()
        else:
            response = auth_resp["status_message"]
            code = auth_resp["status_code"]
        self.headers["Allow"] = '"GET"'
        return pb2.CompositionServiceResponse(StatusCode=code,
                                              Body=bytes(
                                                  json.dumps(response),
                                                  'utf-8'),
                                              Header=self.headers)

    def GetCompositionResource(
            self, request: pb2.GetCompositionResourceRequest,
            context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.OK
        auth_resp = self.auth.is_authorized_rpc(request.SessionToken,
                                                [constants.PrivilegeLogin],
                                                [""])
        if auth_resp["status_code"] == HTTPStatus.OK:
            if request.URL:
                segments = request.URL.split("/")
                if request.ResourceID:
                    # ResourceBlocks Instance
                    if segments[-2] == "ResourceBlocks":
                        response, code = self.resourceblock.get_resource_block(
                            request.URL)
                        self.headers["Allow"] = '"GET", "DELETE"'
                    # ResourceZones Instance
                    elif segments[-2] == "ResourceZones":
                        response, code = self.resourcezone.get_resource_zone(
                            request.URL)
                        self.headers["Allow"] = '"GET", "DELETE"'
                else:
                    # ResourceBlocks Collection
                    if segments[-1] == "ResourceBlocks":
                        response, code = self.resourceblock.get_resource_block_collection(
                            request.URL)
                        self.headers["Allow"] = '"GET", "POST"'
                    # ResourceZones Collection
                    elif segments[-1] == "ResourceZones":
                        response, code = self.resourcezone.get_resource_zone_collection(
                            request.URL)
                        self.headers["Allow"] = '"GET", "POST"'
                    # ActivePool Collection
                    elif segments[-1] == "ActivePool":
                        response, code = self.pool.get_active_pool_collection(
                            request.URL)
                        self.headers["Allow"] = '"GET"'
                    # FreePool Collection
                    elif segments[-1] == "FreePool":
                        response, code = self.pool.get_free_pool_collection(
                            request.URL)
                        self.headers["Allow"] = '"GET"'
                    elif segments[-1] == "CompositionReservations":
                        response, code = self.cs.get_composition_reservations_collection(
                            request.URL)
                        self.headers["Allow"] = '"GET"'
        else:
            response = auth_resp["status_message"]
            code = auth_resp["status_code"]

        return pb2.CompositionServiceResponse(StatusCode=code,
                                              Body=bytes(
                                                  json.dumps(response),
                                                  'utf-8'),
                                              Header=self.headers)

    def CreateCompositionResource(
            self, request: pb2.CreateCompositionResourceRequest,
            context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.CREATED
        headers = {}

        auth_resp = self.auth.is_authorized_rpc(
            request.SessionToken, [constants.PrivilegeConfigureComponents],
            [""])
        if auth_resp["status_code"] == HTTPStatus.OK:
            if request.URL:
                segments = request.URL.split("/")
                # ResourceBlock
                if segments[-1] == "ResourceBlocks":
                    # create resource Block
                    response, code = self.resourceblock.create_resource_block(
                        json.loads(str(request.RequestBody.decode("utf-8"))))
                    headers["Allow"] = '"GET", "POST"'
                    headers["Location"] = response["@odata.id"]
                    headers["Content-type"] = "application/json; charset=utf-8"
                # ResourceZone
                elif segments[-1] == "ResourceZones":
                    # create Resource Zone
                    response, code = self.resourcezone.create_resource_zone(
                        json.loads(str(request.RequestBody.decode("utf-8"))))
                    headers["Allow"] = '"GET", "POST"'
                    headers["Location"] = response["@odata.id"]
                    headers["Content-type"] = "application/json; charset=utf-8"
                # Initialize all Resource Blocks
                elif segments[-1] == "ResourceBlock.Initialize":
                    self.resourceblock.initialize()
                    code = HTTPStatus.NO_CONTENT
                    headers["Allow"] = '"POST"'
        else:
            response = auth_resp["status_message"]
            code = auth_resp["status_code"]

        return pb2.CompositionServiceResponse(StatusCode=code,
                                              Body=bytes(
                                                  json.dumps(response),
                                                  'utf-8'),
                                              Header=headers)

    def DeleteCompositionResource(
            self, request: pb2.DeleteCompositionResourceRequest,
            context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.NO_CONTENT
        auth_resp = self.auth.is_authorized_rpc(
            request.SessionToken, [constants.PrivilegeConfigureComponents],
            [""])
        if auth_resp["status_code"] == HTTPStatus.OK:
            if request.URL:
                segments = request.URL.split("/")
                # ResourceBlock Instance
                if segments[-2] == "ResourceBlocks":
                    # Delete resource Block
                    response, code = self.resourceblock.delete_resource_block(
                        request.URL)
                    self.headers["Allow"] = '"GET", "DELETE"'
                # ResourceZone Instance
                elif segments[-2] == "ResourceZones":
                    # Delete Resource Zone
                    response, code = self.resourcezone.delete_resource_zone(
                        request.URL)
                    self.headers["Allow"] = '"GET", "DELETE"'
        else:
            response = auth_resp["status_message"]
            code = auth_resp["status_code"]

        return pb2.CompositionServiceResponse(StatusCode=code,
                                              Body=bytes(
                                                  json.dumps(response),
                                                  'utf-8'),
                                              Header=self.headers)

    def Compose(
            self, request: pb2.ComposeRequest,
            context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.CREATED
        auth_resp = self.auth.is_authorized_rpc(
            request.SessionToken, [constants.PrivilegeConfigureManager], [""])
        if auth_resp["status_code"] == HTTPStatus.OK:
            response, code = self.cs.compose_action(
                json.loads(str(request.RequestBody.decode("utf-8"))))
            self.headers["Allow"] = '"POST"'
        else:
            response = auth_resp["status_message"]
            code = auth_resp["status_code"]

        return pb2.CompositionServiceResponse(StatusCode=code,
                                              Body=bytes(
                                                  json.dumps(response),
                                                  'utf-8'),
                                              Header=self.headers)
