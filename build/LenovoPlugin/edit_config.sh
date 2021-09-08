#!/bin/bash
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

# Script is for generating certificate and private key
# for Client mode connection usage only

t=/etc/lenovo_plugin_certs
ip=`echo $HOSTIP`

#########changes in lenovo_plugin.json ######
RootServiceUUID=$(uuidgen)
sed -i "s#\"RootServiceUUID\".*#\"RootServiceUUID\": \"${RootServiceUUID}\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"ID\".*#\"ID\": \"LENOVO\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"Host\".*#\"Host\": \"lenovo_plugin\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"Port\".*#\"Port\": \"45009\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"ListenerHost\".*#\"ListenerHost\": \"lenovo_plugin\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"ListenerPort\".*#\"ListenerPort\": \"45010\"#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"RootCACertificatePath\".*#\"RootCACertificatePath\": \"$t/rootCA.crt\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"PrivateKeyPath\".*#\"PrivateKeyPath\": \"$t/odimra_server.key\",#"  /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"CertificatePath\".*#\"CertificatePath\": \"$t/odimra_server.crt\"#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"LBHost\".*#\"LBHost\": \"$ip\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"LBPort\".*#\"LBPort\": \"45010\"#" /etc/lenovo_plugin_config/config_lenovo_plugin.json
sed -i "s#\"MessageQueueConfigFilePath\".*#\"MessageQueueConfigFilePath\": \"/etc/lenovo_plugin_config/platformconfig.toml\",#" /etc/lenovo_plugin_config/config_lenovo_plugin.json

########changes in platformconfig.toml file ######
sed -i "s#.*KServersInfo.*#KServersInfo      = [\"kafka:9092\"]#" /etc/lenovo_plugin_config/platformconfig.toml
sed -i "s#.*KAFKACertFile.*#KAFKACertFile      = \"$t/odimra_kafka_client.crt\"#" /etc/lenovo_plugin_config/platformconfig.toml
sed -i "s#.*KAFKAKeyFile.*#KAFKAKeyFile      = \"$t/odimra_kafka_client.key\"#" /etc/lenovo_plugin_config/platformconfig.toml
sed -i "s#.*KAFKACAFile.*#KAFKACAFile      = \"$t/rootCA.crt\"#" /etc/lenovo_plugin_config/platformconfig.toml

