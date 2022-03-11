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

from concurrent import futures
import grpc
import proto.compositionservice.composition_service_pb2_grpc as pb2_grpc
from rpc import CompositionService
import logging
from config.config import CONFIG_DATA, CERTIFICATES
from config.cli import CL_ARGS
import uuid
from utilities.connection import EtcdConnection


class Services():
    def __init__(self):
        self.server_address = CL_ARGS[
            "ServerAddress"] or "{host}:{port}".format(
                host=CONFIG_DATA["Host"], port=CONFIG_DATA["Port"])
        self.etcd_conn = EtcdConnection()

    def initialize_service(self, service_name):
        if CL_ARGS["FrameWork"] == "GRPC":
            # initialize etcd connection
            self.etcd_conn.etcd_connection()
            cs_service = "{sname}-{uuid}".format(sname=service_name,
                                                 uuid=str(uuid.uuid1()))

            self.etcd_conn.register_service(cs_service, self.server_address)

    def serve_secure(self):
        logging.info("initialize secure service")
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
        pb2_grpc.add_CompositionServicer_to_server(CompositionService(),
                                                   server)
        # Loading credentials
        server_credentials = grpc.ssl_server_credentials(((
            CERTIFICATES["server_private_key"],
            CERTIFICATES["server_certificate"],
        ), ))
        server.add_secure_port(self.server_address, server_credentials)
        server.start()
        logging.info("compositon_service is running")
        server.wait_for_termination()

