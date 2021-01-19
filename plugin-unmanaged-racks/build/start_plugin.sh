#!/bin/bash
mkdir -p /var/log/ur_plugin_logs
export logFolder="/var/log/ur_plugin_logs"
ip=`echo $HOSTIP`
RootServiceUUID=$(uuidgen)
sed -i "s#.*RootServiceUUID\": \"\",# \"RootServiceUUID\": \"${RootServiceUUID}\",#" /etc/ur_plugin_config/ur_config.json
sed -i 's/"Host":\s*".*",/"Host": "plugin",/g' /etc/ur_plugin_config/ur_config.json
sed -i 's/"ListenerHost":\s*".*",/"ListenerHost": "plugin",/g' /etc/ur_plugin_config/ur_config.json
sed -i 's@"RootCACertificatePath":\s*".*",@"RootCACertificatePath": "/etc/plugin_certs/rootCA.crt",@g' /etc/ur_plugin_config/ur_config.json
sed -i 's@"PrivateKeyPath":\s*".*",@"PrivateKeyPath": "/etc/plugin_certs/odimra_server.key",@g' /etc/ur_plugin_config/ur_config.json
sed -i 's@"CertificatePath":\s*".*"@"CertificatePath": "/etc/plugin_certs/odimra_server.crt"@g' /etc/ur_plugin_config/ur_config.json
sed -i "s#.*LBHost.*# \"LBHost\": \"${ip}\",#" /etc/ur_plugin_config/ur_config.json
sed -i "s#.*LBPort.*# \"LBPort\": \"45006\"#" /etc/ur_plugin_config/ur_config.json
sed -i 's@"MessageQueueConfigFilePath":\s*".*",@"MessageQueueConfigFilePath": "/etc/ur_plugin_config/platformconfig.toml",@g' /etc/ur_plugin_config/ur_config.json

sed -i "s#.*KServersInfo.*#KServersInfo      = [\"kafka:9092\"]#" /etc/ur_plugin_config/platformconfig.toml
sed -i "s#.*KAFKACertFile.*#KAFKACertFile      = \"/etc/plugin_certs/odimra_kafka_client.crt\"#" /etc/ur_plugin_config/platformconfig.toml
sed -i "s#.*KAFKAKeyFile.*#KAFKAKeyFile      = \"/etc/plugin_certs/odimra_kafka_client.key\"#" /etc/ur_plugin_config/platformconfig.toml
sed -i "s#.*KAFKACAFile.*#KAFKACAFile      = \"/etc/plugin_certs/rootCA.crt\"#" /etc/ur_plugin_config/platformconfig.toml

systemctl enable plugin

#export PLUGIN_CONFIG_FILE_PATH=/etc/ur_plugin_config/ur_config.json
#sudo -E -u plugin nohup ./plugin-ur --registry=consul --registry_address=consul:8500 --client_request_timeout=1m >> ${logFolder}/ur-plugin.log 2>&1 &
#sleep 2s
while true; do
        sleep 5s
done
