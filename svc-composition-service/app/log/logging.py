import logging
from config.config import PLUGIN_CONFIG


def logger():
    logging.basicConfig(filename=PLUGIN_CONFIG["LogPath"],
                        format='%(asctime)s | %(levelname)s | %(message)s',
                        datefmt='%m/%d/%Y %I:%M:%S %p',
                        level=PLUGIN_CONFIG["LogLevel"].upper())
