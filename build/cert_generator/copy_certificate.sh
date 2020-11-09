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

if id odimra >/dev/null 2>&1; then
        echo "Continue"
else
        echo "Creating odimra user"
	sudo groupadd -r -g 1234 odimra
	sudo useradd -s /bin/bash -u 1234 -m -d /home/odimra -r -g odimra odimra
fi
if id plugin >/dev/null 2>&1; then
        echo "Continue"
else
        echo "Creating plugin user"
	sudo groupadd -r -g 1235 plugin
	sudo useradd -s /bin/bash -u 1235 -m -d /home/plugin -r -g plugin plugin
fi

sudo mkdir /etc/odimracert /etc/kafkacert /etc/plugincert
sudo cp rootCA.crt odimra_server.crt odimra_server.key odimra_rsa.public odimra_rsa.private odimra_kafka_client.crt odimra_kafka_client.key /etc/odimracert/
sudo cp kafka.keystore.jks kafka.truststore.jks /etc/kafka/conf/
sudo cp zookeeper.keystore.jks zookeeper.truststore.jks /etc/zookeeper/conf/
sudo cp rootCA.crt odimra_server.crt odimra_server.key odimra_kafka_client.crt odimra_kafka_client.key /etc/plugincert/

cd /etc/odimracert/

if [ $? -eq 0 ];
then
	a=`echo \`ls | wc -l\` `
	if [ $a -eq 7 ];
	then 
		echo "odimra Certificates copied successfully"
	else
		echo "Copying of odimra Certificates failed"
		exit -1
	fi
else
	echo "Copying of odimra Certificates failed"
	exit -1
fi

cd /etc/kafka/conf

if [ $? -eq 0 ];
then
        a=`echo \`ls | wc -l\` `
        if [ $a -eq 2 ];
        then
                echo "Kafka Certificates copied successfully"
        else
                echo "Copying of Kafka Certificates failed"
                exit -1
        fi
else
        echo "Copying of Kafka Certificates failed"
        exit -1
fi

cd /etc/zookeeper/conf

if [ $? -eq 0 ];
then
        a=`echo \`ls | wc -l\` `
        if [ $a -eq 2 ];
        then
                echo "Zookeeper Certificates copied successfully"
        else
                echo "Copying of Zookeeper Certificates failed"
                exit -1
        fi
else
        echo "Copying of Zookeeper Certificates failed"
        exit -1
fi

cd /etc/plugincert/

if [ $? -eq 0 ];
then
	a=`echo \`ls | wc -l\` `
	if [ $a -eq 5 ];
	then 
		echo "Plugin Certificates copied successfully"
	else
		echo "Copying of Plugin Certificates failed"
		exit -1
	fi
else
	echo "Copying of Plugin Certificates failed"
	exit -1
fi

sudo chown -R odimra:odimra /etc/odimracert/
sudo chown -R plugin:plugin /etc/plugincert/
sudo chown -R odimra:odimra /etc/kafkacert/

