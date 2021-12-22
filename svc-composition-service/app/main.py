from log.logging import logger
from config.config import set_configuraion
from config.cli import collect_cl_args
from config import constants
from rest.resource_blocks import ResourceBlocks
from rest.event_subscription import EventSubscription
from utilities.services import Services

if __name__ == "__main__":

    # initialize configuration
    set_configuraion()
    # initialize log format
    logger()
    # Initialize command line arguments
    collect_cl_args()

    # Initialize resource Block
    rb = ResourceBlocks()
    rb.initialize()

    # create subscription
    EventSubscription()

    services = Services()

    services.initialize_service(constants.COMPOSITION_SERVICE_NAME)

    services.serve_secure()
