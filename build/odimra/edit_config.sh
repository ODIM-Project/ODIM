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

fqdn=`echo $FQDN`
hostip=`echo $HOSTIP`
t=/etc/odimra_config
c=/etc/odimra_certs
d=/etc/registrystore
e=/etc
############changes in odimra_json.json #######
sed -i "s#\"LocalhostFQDN\".*#\"LocalhostFQDN\": \"$fqdn\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"MessageQueueConfigFilePath\".*#\"MessageQueueConfigFilePath\": \"$t/platformconfig.toml\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"SearchAndFilterSchemaPath\".*#\"SearchAndFilterSchemaPath\": \"$e/schema.json\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"RegistryStorePath\".*#\"RegistryStorePath\": \"$d\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"RootCACertificatePath\".*#\"RootCACertificatePath\": \"$c/rootCA.crt\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"RPCPrivateKeyPath\".*#\"RPCPrivateKeyPath\": \"$c/odimra_server.key\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"RPCCertificatePath\".*#\"RPCCertificatePath\": \"$c/odimra_server.crt\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"RSAPublicKeyPath\".*#\"RSAPublicKeyPath\": \"$c/odimra_rsa.public\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"RSAPrivateKeyPath\".*#\"RSAPrivateKeyPath\": \"$c/odimra_rsa.private\"#" /etc/odimra_config/odimra_config.json
sed -i "s#\"InMemoryHost\".*#\"InMemoryHost\": \"redis\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"OnDiskHost\".*#\"OnDiskHost\": \"redis\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"Host\".*#\"Host\": \"odimra\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"PrivateKeyPath\".*#\"PrivateKeyPath\": \"$c/odimra_server.key\",#" /etc/odimra_config/odimra_config.json
sed -i "s#\"CertificatePath\".*#\"CertificatePath\": \"$c/odimra_server.crt\"#" /etc/odimra_config/odimra_config.json

########changes in platformconfig.toml file ######
sed -i "s#.*KServersInfo.*#KServersInfo      = [\"kafka:9092\"]#" /etc/odimra_config/platformconfig.toml
sed -i "s#.*KAFKACertFile.*#KAFKACertFile      = \"$c/odimra_kafka_client.crt\"#" /etc/odimra_config/platformconfig.toml
sed -i "s#.*KAFKAKeyFile.*#KAFKAKeyFile      = \"$c/odimra_kafka_client.key\"#" /etc/odimra_config/platformconfig.toml
sed -i "s#.*KAFKACAFile.*#KAFKACAFile      = \"$c/rootCA.crt\"#" /etc/odimra_config/platformconfig.toml

