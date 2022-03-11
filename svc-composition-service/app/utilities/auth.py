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

import grpc

from utilities.connection import EtcdConnection
from config import constants

import proto.auth.odim_auth_pb2_grpc as auth_pb2_grpc
import proto.auth.odim_auth_pb2 as auth_pb2

from utilities.client_service import ClientService
from http import HTTPStatus


class Auth():
    def __init__(self):
        self.etcd_conn = EtcdConnection()
        self.etcd_conn.etcd_connection()
        self.client_service = ClientService()

    def is_authorized_rpc(self, session_token, privileges, oem_privileges):
        resp = {}
        logging.info("Initialize authorizing rpc")
        # connect accoutsession service
        service_address = self.etcd_conn.get_service_address(
            constants.ACCOUNT_SESSION_NAME)

        # connect odim auth proto stubs
        channel = self.client_service.client_channel(service_address)

        stub = auth_pb2_grpc.AuthorizationStub(channel)

        auth_req = auth_pb2.AuthRequest(sessionToken=session_token,
                                        privileges=privileges,
                                        oemprivileges=oem_privileges)

        try:
            response = stub.IsAuthorized(auth_req)
            logging.debug(
                "Authorized rpc response from server is {res}".format(
                    res=response))
            resp["status_code"] = response.statusCode
            resp["status_message"] = {"Error": response.statusMessage}
        except grpc.RpcError as rpc_error:
            logging.error(
                "Unable to authorize rpc call. Error: {e}".format(e=rpc_error))
            resp["status_code"] = HTTPStatus.BAD_REQUEST
            resp["status_message"] = {
                "Error": "Failed to connect to gRPC server"
            }
        finally:
            return resp
