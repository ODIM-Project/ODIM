from requests.auth import HTTPBasicAuth
from config.config import PLUGIN_CONFIG
from http import HTTPStatus
import requests
import logging


class Client():

    def __init__(self):
        self.auth = HTTPBasicAuth(
            PLUGIN_CONFIG["UserName"], PLUGIN_CONFIG["Password"])
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
            target_url = "{burl}{url}".format(
                burl=PLUGIN_CONFIG["OdimURL"], url=uri)
            response = requests.get(
                target_url, auth=self.auth, headers=self.headers, verify=self.verify)
            if response.status_code == HTTPStatus.OK:
                res = response.json()
            logging.debug("GET Response for the url {url}: {resp}".format(
                url=uri, resp=response.__dict__))
        except Exception as err:
            logging.error(
                "Unable to Process GET Request for uri {url}. Error: {e}".format(url=uri, e=err))
        finally:
            return res

    def process_post_request(self, uri, payload):
        response = None
        try:
            target_url = "{burl}{url}".format(
                burl=PLUGIN_CONFIG["OdimURL"], url=uri)
            response = requests.post(
                target_url, auth=self.auth, headers=self.headers, verify=self.verify, data=payload)
            logging.debug("POST Response for the url {url}: {resp}".format(
                url=uri, resp=response.__dict__))
        except Exception as err:
            logging.error(
                "Unable to Process POST Request for uri {url}. Error: {e}".format(url=uri, e=err))
        finally:
            return response
