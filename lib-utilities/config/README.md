# ODIMRA CONFIG VARIABLES

|   Variables   |   Type    |   Subcategory|Subcategory Type|    Description
|   ------      |   ------- | -------	|  -------  |  ----
|RootServiceUUID|string |||Static uuid used for ODIMRA root service
|LocalhostFQDN|string|||common name of the certificate used for ODIMRA services
|MessageQueueConfigFilePath|string|||File path to the config file which having required configuration details regarding supported message queues
|SearchAndFilterSchemaPath|string|||File path to the search and filter schema file
|RegistryStorePath|string|||Location for storing registry data
|KeyCertConf||RootCACertificatePath|string|TLS root CA file path, which can be a chain of CAs for verifying entities interacting with ODIMRA services
|KeyCertConf||RPCPrivateKeyPath|string|TLS private key file path for the micro service rpc communications
|KeyCertConf||RPCCertificatePath|string|TLS certificate file path for the micro service rpc communications
|KeyCertConf||RSAPublicKeyPath|string|RSA public key file path
|KeyCertConf||RSAPrivateKeyPath|string|RSA private key file path
|APIGatewayConf||Host|string|Host address for the ODIMRA api gateway
|APIGatewayConf||Port|string|Port for the ODIMRA api gateway
|APIGatewayConf||CertificatePath|string|TLS certificate file path for the api gateway
|APIGatewayConf||PrivateKeyPath|string|TLS private key file path for the api gateway
|DBConf||Protocol|string |Redis DB dialing protocol
|DBConf||InMemoryHost|string|Redis DB host for in-memory storage
|DBConf||InMemoryPort|string|Redis DB port for in-memory storage
|DBConf||OnDiskHost|string|Redis DB host for on-disk storage
|DBConf||OnDiskPort|string|Redis DB port for on-disk storage
|DBConf||MaxIdleConns|integer|Maximum number of idle connections allowed in the Redis DB pool
|DBConf||MaxActiveConns|integer|Maximum number of active connections allowed in the Redis DB pool
|FirmwareVersion|string|||version information of the ODIMRA
|SouthBoundRequestTimeoutInSecs|integer|||Timeout for request towards south bound
|ServerRediscoveryBatchSize|integer|||Number of servers can be rediscovered at a time
|AuthConf||SessionTimeOutInMins|integer|Session validity time after each session usage
|AuthConf||ExpiredSessionCleanUpTimeInMins|integer|Duration in minute to clean expired session data from DB
|PasswordRules||MinPasswordLength|integer|This holds the value of min password length
|PasswordRules||MaxPasswordLength|integer|This holds the value of max password length
|PasswordRules||AllowedSpecialCharcters|string|This holds all value of all sppecial charcters
|AddComputeSkipResources|collection|||This stores all resource which need to igonered while adding Computer System
|AddComputeSkipResources||SkipResourceListUnderSystem|list of strings|This holds the value of system resource which need to be ignored
|AddComputeSkipResources||SkipResourceListUnderChassis|list of strings|This holds the value of chassis resource which need to be ignored
|AddComputeSkipResources||SkipResourceListUnderOthers|list of strings|This holds the value resource name for which next level retrieval to be ignored
|URLTranslation|collection|||This holds the north bound and south bound urls
|URLTranslation||NorthBoundURL.ODIM|collection of strings| This the north bound urls
|URLTranslation||SouthBoundURL.redfish|collection of strings| This holds the south bound urls
|PluginStatusPolling||PollingFrequencyInMins|integer|Frequency at which plugin status will be polled
|PluginStatusPolling||MaxRetryAttempt|integer|Max status polling retries
|PluginStatusPolling||RetryIntervalInMins|integer|Interval between status polling retries
|PluginStatusPolling||ResponseTimeoutInSecs|integer|Timeout for status polling requests
|PluginStatusPolling||StartUpResouceBatchSize|integer|Number of resources to retrieve in batch
|ExecPriorityDelayConf||MinResetPriority|integer|Minimum priority for a serverreset action
|ExecPriorityDelayConf||MaxResetPriority|integer|Maximum priority for a server reset action
|ExecPriorityDelayConf||MaxResetDelayInSecs|integer|Maximum delay before executing server reset action
|EnabledServices|list of strings|||List of services enabled
|TLSConf||MinVersion|string|Minimum TLS version
|TLSConf||MaxVersion|string|Maximum TLS version
|TLSConf||VerifyPeer|boolean|If server validation is required
|TLSConf||PreferredCipherSuites |list of string|Preferred list of cipher suites
