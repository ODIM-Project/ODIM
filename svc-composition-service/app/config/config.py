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

CONFIG_DATA = {
    "Host": "",
    "Port": "",
    "UserName": "",
    "Password": "",
    "RootServiceUUID": "",
    "OdimURL": "",
    "OdimUserName": "",
    "OdimPassword": "",
    "LogLevel": "info",
    "RedisOnDiskAddress": "",
    "RedisInMemoryAddress": "",
    "Db": 0,
    "SocketTimeout": 10,
    "LogPath": "",
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
        if config_data.get("RootServiceUUID"):
            CONFIG_DATA["RootServiceUUID"] = config_data["RootServiceUUID"]

        if config_data.get("LocalhostFQDN"):
            CONFIG_DATA["LocalhostFQDN"] = config_data["LocalhostFQDN"]

        if config_data.get("APIGatewayConf"):
            if config_data["APIGatewayConf"].get(
                    "Host") and config_data["APIGatewayConf"].get("Port"):
                CONFIG_DATA["OdimURL"] = "{protocol}://{host}:{port}".format(
                    protocol="https",
                    host=config_data["APIGatewayConf"]["Host"],
                    port=config_data["APIGatewayConf"]["Port"])

        if config_data.get("DBConf"):
            if config_data["DBConf"].get("InMemoryHost") and config_data[
                    "DBConf"].get("InMemoryPort"):
                CONFIG_DATA["RedisInMemoryAddress"] = "{host}:{port}".format(
                    host=config_data["DBConf"]["InMemoryHost"],
                    port=config_data["DBConf"]["InMemoryPort"])

            if config_data["DBConf"].get(
                    "OnDiskHost") and config_data["DBConf"].get("OnDiskPort"):
                CONFIG_DATA["RedisOnDiskAddress"] = "{host}:{port}".format(
                    host=config_data["DBConf"]["OnDiskHost"],
                    port=config_data["DBConf"]["OnDiskPort"])

        if config_data.get("KeyCertConf"):
            if config_data["KeyCertConf"].get("RootCACertificatePath"):
                CONFIG_DATA["RootCAPath"] = config_data["KeyCertConf"][
                    "RootCACertificatePath"]
            if config_data["KeyCertConf"].get("RPCPrivateKeyPath"):
                CONFIG_DATA["PrivateKeyPath"] = config_data["KeyCertConf"][
                    "RPCPrivateKeyPath"]
            if config_data["KeyCertConf"].get("RPCCertificatePath"):
                CONFIG_DATA["CertificatePath"] = config_data["KeyCertConf"][
                    "RPCCertificatePath"]
            if config_data["KeyCertConf"].get("RSAPublicKeyPath"):
                CONFIG_DATA["RSAPublicKeyPath"] = config_data["KeyCertConf"][
                    "RSAPublicKeyPath"]
            if config_data["KeyCertConf"].get("RSAPrivateKeyPath"):
                CONFIG_DATA["RSAPrivateKeyPath"] = config_data["KeyCertConf"][
                    "RSAPrivateKeyPath"]

    # get server private key data from PrivateKeyPath
    if os.path.exists(CONFIG_DATA["PrivateKeyPath"]):
        CERTIFICATES["server_private_key"] = _load_credential_from_file(
            CONFIG_DATA["PrivateKeyPath"])
    # get server certificate data from CertificatePath
    if os.path.exists(CONFIG_DATA["CertificatePath"]):
        CERTIFICATES["server_certificate"] = _load_credential_from_file(
            CONFIG_DATA["CertificatePath"])
    # get root ca certificate data from RootCAPath
    if os.path.exists(CONFIG_DATA["RootCAPath"]):
        CERTIFICATES["root_ca_certificate"] = _load_credential_from_file(
            CONFIG_DATA["RootCAPath"])


def _load_credential_from_file(filepath):
    with open(filepath, 'rb') as f:
        return f.read()
