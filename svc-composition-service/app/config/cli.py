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

# parse command line arguments
import argparse

CL_ARGS = {
    "Registry": "",
    "RegistryAddress": "",
    "ServerAddress": "",
    "ClientRequestTimeout": "",
    "FrameWork": ""
}


def collect_cl_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--registry", default="", help="service registry")
    parser.add_argument("--registry_address",
                        default="",
                        help="address of the registry")
    parser.add_argument("--server_address",
                        default="",
                        help="address for the micro service")
    parser.add_argument("--client_request_timeout",
                        default="",
                        help="maximum request time which client waits")
    parser.add_argument("--framework",
                        default="GRPC",
                        help="framework used for micro service communication")

    args = vars(parser.parse_args())

    if args:
        CL_ARGS["Registry"] = args["registry"]
        CL_ARGS["RegistryAddress"] = args["registry_address"]
        CL_ARGS["ServerAddress"] = args["server_address"]
        CL_ARGS["ClientRequestTimeout"] = args["client_request_timeout"]
        CL_ARGS["FrameWork"] = args["framework"]
