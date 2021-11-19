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
    parser.add_argument("--registry_address", default="",
                        help="address of the registry")
    parser.add_argument("--server_address", default="",
                        help="address for the micro service")
    parser.add_argument("--client_request_timeout", default="",
                        help="maximum request time which client waits")
    parser.add_argument("--framework", default="GRPC",
                        help="framework used for micro service communication")

    args = vars(parser.parse_args())

    if args:
        CL_ARGS["Registry"] = args["registry"]
        CL_ARGS["RegistryAddress"] = args["registry_address"]
        CL_ARGS["ServerAddress"] = args["server_address"]
        CL_ARGS["ClientRequestTimeout"] = args["client_request_timeout"]
        CL_ARGS["FrameWork"] = args["framework"]
