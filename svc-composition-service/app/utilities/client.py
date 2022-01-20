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

from requests.auth import HTTPBasicAuth
from config.config import PLUGIN_CONFIG
from http import HTTPStatus
import requests
import logging
from utilities.crypt import Crypt


class Client():
    def __init__(self):
        crypt = Crypt(PLUGIN_CONFIG["RSAPublicKeyPath"],
                      PLUGIN_CONFIG["RSAPrivateKeyPath"])
        self.auth = HTTPBasicAuth(PLUGIN_CONFIG["OdimUserName"],
                                  crypt.decrypt(PLUGIN_CONFIG["OdimPassword"]))
        self.headers = {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
            'If-Match': '*'
        }
        self.verify = False

    def process_get_request(self, uri):
        res = {}
        if not uri:
            return res
        try:
            target_url = "{burl}{url}".format(burl=PLUGIN_CONFIG["OdimURL"],
                                              url=uri)
            response = requests.get(target_url,
                                    auth=self.auth,
                                    headers=self.headers,
                                    verify=self.verify)
            if response.status_code == HTTPStatus.OK:
                res = response.json()
            logging.debug("GET Response for the url {url}: {resp}".format(
                url=uri, resp=response.__dict__))
        except Exception as err:
            logging.error(
                "Unable to Process GET Request for uri {url}. Error: {e}".
                format(url=uri, e=err))
        finally:
            return res

    def process_post_request(self, uri, payload):
        response = None
        try:
            target_url = "{burl}{url}".format(burl=PLUGIN_CONFIG["OdimURL"],
                                              url=uri)
            response = requests.post(target_url,
                                     auth=self.auth,
                                     headers=self.headers,
                                     verify=self.verify,
                                     data=payload)
            logging.debug("POST Response for the url {url}: {resp}".format(
                url=uri, resp=response.__dict__))
        except Exception as err:
            logging.error(
                "Unable to Process POST Request for uri {url}. Error: {e}".
                format(url=uri, e=err))
        finally:
            return response
