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


from config.config import CONFIG_DATA
from http import HTTPStatus
import requests
from requests.auth import HTTPBasicAuth
import logging
from utilities.crypt import Crypt


class Client():
    def __init__(self):
        crypt = Crypt(CONFIG_DATA["RSAPublicKeyPath"],
                      CONFIG_DATA["RSAPrivateKeyPath"])
        self.auth = HTTPBasicAuth(CONFIG_DATA["OdimUserName"],
                                  crypt.decrypt(CONFIG_DATA["OdimPassword"]))
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
            target_url = "{burl}{url}".format(burl=CONFIG_DATA["OdimURL"],
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
            target_url = "{burl}{url}".format(burl=CONFIG_DATA["OdimURL"],
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
