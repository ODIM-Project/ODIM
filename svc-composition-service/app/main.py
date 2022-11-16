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

from log.logging import logger
from config.config import set_configuraion
from config.cli import collect_cl_args
from config import constants
# from rest.resource_blocks import ResourceBlocks
# from rest.event_subscription import EventSubscription
from utilities.services import Services

if __name__ == "__main__":

    # initialize configuration
    set_configuraion()
    # initialize log format
    logger()
    # Initialize command line arguments
    collect_cl_args()

    # Initialize resource Block
    """
    rb = ResourceBlocks()
    rb.initialize()
    """

    # create subscription
    #EventSubscription()

    services = Services()

    services.initialize_service(constants.COMPOSITION_SERVICE_NAME)

    services.serve_secure()
