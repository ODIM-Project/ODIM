# URP PLUGIN CONFIG VARIABLES

|   Variables   |   Type    |   Description|
|   ---------   |   ------- | -------   |
|ID|string|Identifier used by ODIMRA for identifying the plugin
|Host|string|plugin host address for ODIMRA to contact plugin
|Port|string|plugin port for ODIMRA to contact plugin
|UserName|string|plugin user name for ODIMRA to interact with plugin
|Password|string|plugin password for ODIMRA to interact with plugin
|RootServiceUUID|string|Static uuid used for plugin root service
|OdimraNBUrl|string|URL pointing to ODIMRA's NB interface
|FirmwareVersion|string|version information of the plugin
|SessionTimeoutInMinutes|integer|Plugin session time out in minutes
|**EventConf:**||
|&nbsp; DestinationURI|string|URI that will be posted on the resource as destination for events
|&nbsp; ListenerHost|string|Host address that will be posted on the resource as destination for events
|&nbsp; ListenerPort|string|Host address port that will be posted on the resource as destination for events
|**KeyCertCon:**||
|&nbsp; RootCACertificatePath|string|TLS root certificate
|&nbsp; PrivateKeyPath|string|Plugin private key path for ODIMRA and plugin interaction 
|&nbsp; CertificatePath|string|Plugin certificate path for ODIMRA and plugin interaction
|**URLTranslation:**||
|&nbsp; NorthBoundURL|map[string]|This holds the north bound urls
|&nbsp; SouthBoundURL|map[string]|This holds the south bound urls
|**TLSConf:**||
|&nbsp; MinVersion|string|Minimum TLS version
|&nbsp; MaxVersion|string|Maximum TLS version
|&nbsp; VerifyPeer|boolean|If server validation is required
|&nbsp; PreferredCipherSuites|[]string|Preferred list of cipher suites
