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

import grpc

from config.config import CONFIG_DATA, CERTIFICATES


class ClientService():
    def __init__(self):
        self._channel = None

    def client_channel(self, server_address):

        channel_credentials = grpc.ssl_channel_credentials(
            root_certificates=CERTIFICATES["root_ca_certificate"])
        
        options = (('grpc.ssl_target_name_override', CONFIG_DATA["LocalhostFQDN"],),)

        _channel = grpc.secure_channel(server_address,
                                       credentials=channel_credentials,options=options)

        return _channel

    def __exit__(self):
        if self._channel:
            self._channel.close()

