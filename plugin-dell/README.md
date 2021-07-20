# Dell-plugin  

Dell-plugin communicates with redfish compliant BMC.  
This is an independent module which provides two primary communication channels:  
1.  An API mechanism that is used to exchange control data  
2.  An Event Message Bus (EMB) that is used to exchange event and notifications.


This guide provides a set of guidelines for developing API and EMB functions to work within the Resource Aggregator for ODIM™ environment. It ensures consistency around API semantics for all plugins.

To ensure continued adoption of open technologies, the APIs for the plugins are based on the [OpenAPI specification](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md). Messaging is based on the now-evolving [OpenMessaging specifications](https://github.com/openmessaging/specification/blob/master/domain_architecture.md) under Linux Foundation.



## Dell plugin deployment instructions

For deploying the Dell plugin and adding the plugin to the Resource Aggregator for ODIM framework, refer to the "Deploying the Dell plugin" section in the [Resource Aggregator for Open Distributed Infrastructure Management™ Readme](https://github.com/ODIM-Project/ODIM/blob/main/README.md).

