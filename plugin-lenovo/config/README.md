# LENOVO PLUGIN CONFIG VARIABLES

|   Variables   |   Type    |   Subcategory|Subcategory Type|    Description
|   ---------   |   ------- | -------   |  -------  |  ----
|RootServiceUUID|string |||Static uuid used for plugin root service
|PluginConf||ID|string|Identifier used by ODIMRA for identifying the plugin
|PluginConf||Host|string|plugin host address for ODIMRA to contact plugin
|PluginConf||Port|string|plugin port for ODIMRA to contact plugin
|PluginConf||UserName|string|plugin user name for ODIMRA to interact with plugin
|PluginConf||Password|string|plugin password for ODIMRA to interact with plugin
|EventConf||DestinationURI|string|URI that will be posted on the resource as destination for events
|EventConf||ListenerHost|string|Host address that will be posted on the resource as destination for events
|EventConf||ListenerPort|string|Host address port that will be posted on the resource as destination for events
|KeyCertCon||RootCACertificatePath|string|TLS root certificate
|KeyCertCon||PrivateKeyPath|string|Plugin private key path for ODIMRA and plugin interaction 
|KeyCertCon||CertificatePath|string|Plugin certificate path for ODIMRA and plugin interaction
|FirmwareVersion|string|||version information of the plugin
|SessionTimeoutInMinutes|integer|||Plugin session time out in minutes
|LoadBalancerConf||LBHost|string|Load Balancer host address for plugin
|LoadBalancerConf||LBPort|string|Load Balancer host address port for plugin
|MessageBusConf||MessageQueueConfigFilePath|string|||File path to the config file which having required configuration details regarding supported message queues 
|MessageBusConf||MessageBusType|string|This holds information Event Message Bus Type
|MessageBusConf||MessageBusQueue|list of strings|This holds name of all message bus Queues
|URLTranslation|collection|||This holds the north bound and south bound urls
|URLTranslation||NorthBoundURL.ODIM|collection of strings| This the north bound urls
|URLTranslation||SouthBoundURL.redfish|collection of strings| This holds the south bound urls
|TLSConf||MinVersion|string|Minimum TLS version
|TLSConf||MaxVersion|string|Maximum TLS version
|TLSConf||VerifyPeer|boolean|If server validation is required
|TLSConf||PreferredCipherSuites |list of string|Preferred list of cipher suites
