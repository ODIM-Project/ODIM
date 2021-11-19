import json
from logging import log
import logging
import grpc
from http import HTTPStatus
import proto.composition_service_pb2 as pb2
import proto.composition_service_pb2_grpc as pb2_grpc
from rest.resource_zones import ResourceZones
from rest.resource_blocks import ResourceBlocks
from rest.composition_service import CompositonService


class CompositionServiceRpc(pb2_grpc.CompositionServicer):
    """The listener function implemests the rpc call as described in the .proto file"""

    def __init__(self):
        self.cs = CompositonService()
        self.resourcezone = ResourceZones()
        self.resourceblock = ResourceBlocks()

    def __str__(self):
        return self.__class__.__name__

    def GetCompositionService(self, request: pb2.GetCompositionServiceRequest, context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response, code = self.cs.get_cs()
        return pb2.CompositionServiceResponse(statusCode=code, body=bytes(json.dumps(response)))
        # return pb2.CompositionServiceResponse(statusCode=code, body=response)

    def GetCompositionResource(self, request: pb2.GetCompositionResourceRequest, context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.OK
        logging.info("GetCompositionResource request: {}".format(request))
        if request.URL:
            segments = request.URL.split("/")
            if request.resourceID:
                # ResourceBlocks Instance
                if segments[-2] == "ResourceBlocks":
                    logging.info("In Resource Blocks")
                    response, code = self.resourceblock.get_resource_block(
                        request.URL)
                # ResourceZones Instance
                elif segments[-2] == "ResourceZones":
                    response, code = self.resourcezone.get_resource_zone(
                        request.URL)
            else:
                # ResourceBlocks Collection
                if segments[-1] == "ResourceBlocks":
                    response, code = self.resourceblock.get_resource_block_collection(
                        request.URL)
                # ResourceZones Collection
                elif segments[-1] == "ResourceZones":
                    response, code = self.resourcezone.get_resource_zone_collection(
                        request.URL)

        return pb2.CompositionServiceResponse(statusCode=code, body=bytes(json.dumps(response), 'utf-8'))

    def CreateCompositionResource(self, request: pb2.CreateCompositionResourceRequest, context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.CREATED
        if request.URL:
            segments = request.URL.split("/")
            # ResourceBlocks Collection
            if segments[-1] == "ResourceBlocks":
                # create resource Block
                pass
            # ResourceZones Collection
            elif segments[-1] == "ResourceZones":
                # create Resource Zone
                response, code = self.resourcezone.create_resource_zone(
                    json.loads(str(request.RequestBody.decode("utf-8"))))

        return pb2.CompositionServiceResponse(statusCode=code, body=bytes(json.dumps(response), 'utf-8'))

    def DeleteCompositionResource(self, request: pb2.DeleteCompositionResourceRequest, context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.NO_CONTENT
        if request.URL:
            segments = request.URL.split("/")
            # ResourceBlock Instance
            if segments[-2] == "ResourceBlocks":
                # Delete resource Block
                pass
            # ResourceZone Instance
            elif segments[-2] == "ResourceZones":
                # Delete Resource Zone
                response, code = self.resourcezone.delete_resource_zone(
                    request.URL)

        return pb2.CompositionServiceResponse(statusCode=code, body=bytes(json.dumps(response), 'utf-8'))

    def Compose(self, request: pb2.ComposeRequest, context: grpc.ServicerContext) -> pb2.CompositionServiceResponse:
        response = {}
        code = HTTPStatus.CREATED
        try:
            response, code = self.cs.compose_action(
                json.loads(str(request.RequestBody.decode("utf-8"))))
            return pb2.CompositionServiceResponse(StatusCode=code, Body=bytes(json.dumps(response), 'utf-8'))
        except Exception as err:
            logging.error("Compose Error: {}".format(err))
