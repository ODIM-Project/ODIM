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
"""
from config.config import CONFIG_DATA
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
            "Initialize Event Service Subscription for Composition Service Event"
        )
        response = None
        data = {
            "Name":
            "EventSubscription",
            "Destination":
            "{protocol}://{host}:{port}{url}".format(
                protocol="https",
                host=CONFIG_DATA["Host"],
                port=CONFIG_DATA["Port"],
                url=constants.EVENT_DESTINATION_URL),
            "EventTypes":
            constants.EVENT_TYPES,
            "Context":
            "Subscribed_by_Composition_Service",
            "Protocol":
            "Redfish",
            "SubscriptionType":
            "RedfishEvent",
            "EventFormatType":
            "Event",
            "SubordinateResources":
            False,
            "ResourceTypes": ["ComputerSystem"],
            "OriginResources": [{
                "@odata.id": "/redfish/v1/Systems"
            }]
        }

        try:
            # get systems collection data
            response = self.client.process_post_request(
                constants.EVENT_SUBSCRIPTION_URL, json.dumps(data))
            if response is not None:
                if response.status_code in [
                        HTTPStatus.OK, HTTPStatus.ACCEPTED,
                        HTTPStatus.NO_CONTENT
                ]:
                    logging.info(
                        "Composition Service Event Subscription is Successfully Registered"
                    )
                else:
                    logging.info(
                        "Composition Service Event Subscription has been rejected with status code {code}"
                        .format(code=response.status_code))
            else:
                logging.info("Composition Service Event Subscription Failed")
        except Exception as err:
            logging.error(
                "unable to initialize the Event Subscription for Composition Service. Error: {e}"
                .format(e=err))
"""
