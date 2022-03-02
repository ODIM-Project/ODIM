import grpc

from config.config import PLUGIN_CONFIG, CERTIFICATES


class ClientService():
    def __init__(self):
        self._channel = None

    def client_channel(self, server_address):

        channel_credentials = grpc.ssl_channel_credentials(
            root_certificates=CERTIFICATES["root_ca_certificate"])
        
        options = (('grpc.ssl_target_name_override', PLUGIN_CONFIG["LocalhostFQDN"],),)

        _channel = grpc.secure_channel(server_address,
                                       credentials=channel_credentials,options=options)

        return _channel

    def __exit__(self):
        if self._channel:
            self._channel.close()

