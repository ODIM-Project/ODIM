# lib-utilities

lib-utilities is a library containing all the common models and functions. It also contains the common configuration file for Resource Aggregator for ODIM (ODIMRA).

Packages in this library are as follows:
1. config—Contains the ODIMRA configuration file, functions, and models for processing and validating configuration file.
2. common—Contains miscellaneous models and functions which can be used across multiple ODIMRA services and other libraries.
3. errors—Contains the definition and operations of ODIMRA custom errors especially for database related operations.
4. proto—Holds the proto files and the auto-generated supporting files for enabling RPC communication between ODIMRA services.
5. response—Defines the common responses of ODIMRA. This also contains functions for creating error responses that matches the Redfish standards.
6. services—Acts as the backbone for the ODIMRA micro services. The package contains the functions for creating microservice clients and servers.
