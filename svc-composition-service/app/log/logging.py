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

import logging
from config.config import PLUGIN_CONFIG


def logger():
    if PLUGIN_CONFIG["LogPath"]:
        logging.basicConfig(
            filename=PLUGIN_CONFIG["LogPath"],
            format='%(asctime)s | %(levelname)s | %(message)s',
            datefmt='%m/%d/%Y %I:%M:%S %p',
            level=PLUGIN_CONFIG["LogLevel"].upper() if isinstance(
                PLUGIN_CONFIG["LogLevel"], str) else None)

    else:
        logging.basicConfig(
            format='%(asctime)s | %(levelname)s | %(message)s',
            datefmt='%m/%d/%Y %I:%M:%S %p',
            level=PLUGIN_CONFIG["LogLevel"].upper() if isinstance(
                PLUGIN_CONFIG["LogLevel"], str) else None)
