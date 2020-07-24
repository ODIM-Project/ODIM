# lib-utilities

lib-utilities is a library containing all the common models and functions. It also contains the common configuration file for ODIMRA.

Packages in this library are as follows:
1. config:   contains the ODIMRA configuration file, functions, and models for processing and validating configuration file.
2. common:   contains miscellaneous models and functions which can be used across multiple ODIMRA services and other libraries.
3. errors:   contains the definition and operations of ODIMRA custom error especially for database related operations.
4. proto:    holds the proto files and the auto-generated supporting files for enabling RPC communication between ODIMRA services.
5. response: this defines the common responses of ODIMRA. This also contains functions for creating error responses that matches the redfish standards.
6. services: acts as the backbone for the ODIMRA micro services. The package contains the functions for creating micro service clients and servers.
