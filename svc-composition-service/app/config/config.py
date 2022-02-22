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

import os
import json

CONF_FILE = os.getenv("CONFIG_FILE_PATH")
ODIM_CONF_FILE = os.getenv("ODIM_CONFIG_FILE_PATH")

PLUGIN_CONFIG = {
    "Host": "",
    "Port": "",
    "UserName": "",
    "Password": "",
    "RootServiceUUID": "",
    "OdimURL": "",
    "OdimUserName": "",
    "OdimPassword": "",
    "LogLevel": "",
    "RedisOnDiskAddress": "",
    "RedisInMemoryAddress": "",
    "Db": 0,
    "SocketTimeout": 10,
    "LogPath": "",
    "RedisAddress": "",
    "PrivateKeyPath": "",
    "CertificatePath": "",
    "RootCAPath": "",
    "RSAPrivateKeyPath": "",
    "RSAPublicKeyPath": "",
    "LocalhostFQDN": ""
}

CERTIFICATES = {
    "server_certificate": "",
    "server_private_key": "",
    "root_ca_certificate": ""
}


def set_configuraion():
    config_data = {}
    if CONF_FILE and os.path.exists(CONF_FILE):
        with open(CONF_FILE) as f:
            try:
                config_data = json.load(f)
            except Exception:
                pass
    if config_data:
        for key in PLUGIN_CONFIG.keys():
            if config_data.get(key):
                PLUGIN_CONFIG[key] = config_data[key]

    # get server private key data from PrivateKeyPath
    if os.path.exists(PLUGIN_CONFIG["PrivateKeyPath"]):
        CERTIFICATES["server_private_key"] = _load_credential_from_file(
            PLUGIN_CONFIG["PrivateKeyPath"])
    # get server certificate data from CertificatePath
    if os.path.exists(PLUGIN_CONFIG["CertificatePath"]):
        CERTIFICATES["server_certificate"] = _load_credential_from_file(
            PLUGIN_CONFIG["CertificatePath"])
    # get root ca certificate data from RootCAPath
    if os.path.exists(PLUGIN_CONFIG["RootCAPath"]):
        CERTIFICATES["root_ca_certificate"] = _load_credential_from_file(
            PLUGIN_CONFIG["RootCAPath"])

    if not PLUGIN_CONFIG["LocalhostFQDN"] and ODIM_CONF_FILE and os.path.exists(ODIM_CONF_FILE):
        with open(ODIM_CONF_FILE) as f:
            try:
                odim_config_data = json.load(f)
                if odim_config_data.get("LocalhostFQDN"):
                    PLUGIN_CONFIG["LocalhostFQDN"] = odim_config_data["LocalhostFQDN"]
            except Exception:
                pass

def _load_credential_from_file(filepath):
    with open(filepath, 'rb') as f:
        return f.read()
