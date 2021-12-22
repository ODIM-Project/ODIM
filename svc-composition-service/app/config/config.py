import os
import json

CONF_FILE = os.getenv("PLUGIN_CONFIG_FILE_PATH")

PLUGIN_CONFIG = {
    "Host": "csplugin",
    "Port": "45100",
    "UserName": "",
    "Password": "",
    "RootServiceUUID": "",
    "OdimURL": "",
    "OdimUserName": "",
    "OdimPassword": "",
    "LogLevel": "debug",
    "RedisAddress": "",
    "Db": 0,
    "SocketTimeout": 10,
    "LogPath": "",
    "RedisAddress": "",
    "PrivateKeyPath": "",
    "CertificatePath": "",
    "RootCAPath": ""
}

CERTIFICATES = {
    "server_certificate": "",
    "server_private_key": "",
    "root_ca_certificate": ""
}


def set_configuraion():
    config_data = {}
    if os.path.exists(CONF_FILE):
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
    if os.path.exists(config_data["PrivateKeyPath"]):
        CERTIFICATES["server_private_key"] = _load_credential_from_file(
            config_data["PrivateKeyPath"])
    # get server certificate data from CertificatePath
    if os.path.exists(config_data["CertificatePath"]):
        CERTIFICATES["server_certificate"] = _load_credential_from_file(
            config_data["CertificatePath"])
    # get root ca certificate data from RootCAPath
    if os.path.exists(config_data["RootCAPath"]):
        CERTIFICATES["root_ca_certificate"] = _load_credential_from_file(
            config_data["RootCAPath"])


def _load_credential_from_file(filepath):
    with open(filepath, 'rb') as f:
        return f.read()
