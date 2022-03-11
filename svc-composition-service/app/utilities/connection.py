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

import etcd3
import logging
from config.config import CONFIG_DATA
from config.cli import CL_ARGS


class EtcdConnection():
    def __init__(self):
        self.etcd_host = None
        self.etcd_port = None
        if ":" in CL_ARGS["RegistryAddress"]:
            self.etcd_host, self.etcd_port = CL_ARGS["RegistryAddress"].split(
                ":")

        self.time_out = 1  # one second
        self.etcd_client = None

    def etcd_connection(self):
        logging.info("initialize etcd connection")
        if self.etcd_host and self.etcd_port:
            self.etcd_client = etcd3.client(
                host=self.etcd_host,
                port=self.etcd_port,
                ca_cert=CONFIG_DATA["RootCAPath"],
                cert_key=CONFIG_DATA["PrivateKeyPath"],
                cert_cert=CONFIG_DATA["CertificatePath"],
                timeout=self.time_out)
        else:
            # default etcd host and port
            self.etcd_client = etcd3.client(
                ca_cert=CONFIG_DATA["RootCAPath"],
                cert_key=CONFIG_DATA["PrivateKeyPath"],
                cert_cert=CONFIG_DATA["CertificatePath"],
                timeout=self.time_out)

    def register_service(self, service, server_address):
        if service:
            # putting server address into etcd
            self.put(service, server_address)

    def get_service_address(self, service_name):
        if service_name:
            service_address = self.get(service_name)
            if not service_address:
                logging.warning(
                    "No service with {name} found in the service registry".
                    format(name=service_name))

            return service_address

    def put(self, key, value):
        logging.info("putting data into etcd key {key}".format(key=key))
        try:
            if self.etcd_client:
                self.etcd_client.put(key=key, value=value)
                logging.info(
                    "Successfully registered {key} into etcd server".format(
                        key=key))
        except Exception as err:
            logging.error(
                "Unable to put data into etcd server. Error: {e}".format(
                    e=err))

    def get(self, key):
        logging.info("Getting etcd key {key} from etcd server".format(key=key))
        try:
            if self.etcd_client:
                resp = self.etcd_client.get_prefix_response(key).kvs
                if resp:
                    return resp[0].value.decode('utf-8')
        except Exception as err:
            logging.error(
                "Unable to get key {key} from etcd server. Error: {e}".format(
                    key=key, e=err))
