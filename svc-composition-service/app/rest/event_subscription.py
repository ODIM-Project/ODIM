from config.config import PLUGIN_CONFIG
import logging
import json
from utilities.client import Client
from config import constants
from http import HTTPStatus
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)


class EventSubscription():
    def __init__(self):
        self.client = Client()
        self.initialize()

    def initialize(self):
        logging.info(
            "Initialize Event Service Subscription for Composition Service Event")
        response = None
        data = {
            "Name": "EventSubscription",
            "Destination": "{protocol}://{host}:{port}{url}".format(protocol="http", host=PLUGIN_CONFIG["Host"], port=PLUGIN_CONFIG["Port"], url=constants.EVENT_DESTINATION_URL),
            "EventTypes": constants.EVENT_TYPES,
            "Context": "Subscribed_by_Composition_Service",
            "Protocol": "Redfish",
            "SubscriptionType": "RedfishEvent",
            "EventFormatType": "Event",
            "SubordinateResources": False,
            "ResourceTypes": ["ComputerSystem"],
            "OriginResources": [
                {
                    "@odata.id": "/redfish/v1/Systems"
                }
            ]
        }

        try:
            # get systems collection data
            response = self.client.process_post_request(
                constants.EVENT_SUBSCRIPTION_URL, json.dumps(data))
            if response is not None:
                if response.status_code in [HTTPStatus.OK, HTTPStatus.ACCEPTED, HTTPStatus.NO_CONTENT]:
                    logging.info(
                        "Composition Service Event Subscription is Successfully Registered")
                else:
                    logging.info("Composition Service Event Subscription has been rejected with status code {code}".format(
                        code=response.status_code))
            else:
                logging.info("Composition Service Event Subscription Failed")
        except Exception as err:
            logging.error(
                "unable to initialize the Event Subscription for Composition Service")
