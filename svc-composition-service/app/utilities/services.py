from concurrent import futures
import grpc
import proto.composition_service_pb2_grpc as pb2_grpc
from rpc import CompositionService
import etcd3
import sys
import logging
from config.config import PLUGIN_CONFIG, CERTIFICATES
from config.cli import CL_ARGS
from config import constants
import uuid


class Services():

    def __init__(self):
        self.etcd_host = None
        self.etcd_port = None
        if ":" in CL_ARGS["RegistryAddress"]:
            self.etcd_host, self.etcd_port = CL_ARGS["RegistryAddress"].split(
                ":")

        self.time_out = 1  # one second
        self.etcd_client = None
        self.server_address = CL_ARGS["ServerAddress"] or "{host}:{port}".format(
            host=PLUGIN_CONFIG["Host"], port=PLUGIN_CONFIG["Port"])

    def initialize_service(self, service_name):
        if CL_ARGS["FrameWork"] == "GRPC":
            # initialize etcd connection
            self.etcd_connection()
            cs_service = "{sname}-{uuid}".format(
                sname=service_name, uuid=str(uuid.uuid1()))

            self.register_service(cs_service)

    def register_service(self, service):
        if service:
            # putting server address into etcd
            self.put(service, self.server_address)

    def serve_secure(self):
        logging.info("initialize secure service")
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
        pb2_grpc.add_CompositionServicer_to_server(
            CompositionService(), server)
        # Loading credentials
        server_credentials = grpc.ssl_server_credentials(((
            CERTIFICATES["server_private_key"],
            CERTIFICATES["server_certificate"],
        ),))
        server.add_secure_port(self.server_address, server_credentials)
        server.start()
        logging.info("in wait")
        server.wait_for_termination()

    def etcd_connection(self):
        logging.info("initialize etcd connection")
        if self.etcd_host and self.etcd_port:
            self.etcd_client = etcd3.client(host=self.etcd_host, port=self.etcd_port, ca_cert=PLUGIN_CONFIG[
                                            "RootCAPath"], cert_key=PLUGIN_CONFIG["PrivateKeyPath"], cert_cert=PLUGIN_CONFIG["CertificatePath"], timeout=self.time_out)
        else:
            # default etcd host and port
            self.etcd_client = etcd3.client(
                ca_cert=PLUGIN_CONFIG["RootCAPath"], cert_key=PLUGIN_CONFIG["PrivateKeyPath"], cert_cert=PLUGIN_CONFIG["CertificatePath"], timeout=self.time_out)

    def put(self, key, value):
        logging.info("putting data {val} into etcd key {key}".format(
            val=value, key=key))
        try:
            if self.etcd_client:
                self.etcd_client.put(key=key, value=value)
                logging.info(
                    "Successfully registered {key} into etcd server".format(key=key))
        except Exception as err:
            logging.error(
                "Unable to put data into etcd server. Error: {e}".format(e=err))

    def get(self, key):
        logging.info("Getting etcd key {key} from etcd server".format(key=key))
        try:
            if self.etcd_client:
                return self.etcd_client.get(key)
        except Exception as err:
            logging.error(
                "Unable to get key {key} from etcd server. Error: {e}".format(key=key, e=err))
