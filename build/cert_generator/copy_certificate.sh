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

logerr()
{
        echo "[$(date)] -- ERROR -- $1"
        _exit 1
}

eval_cmd_exec()
{
        if [[ $# -lt 3 ]]; then
                logerr "eval_cmd_exec syntax error $2"
        fi
        if [[ $1 -ne 0 ]]; then
                echo "$3"
                logerr $2
        fi
}

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

#create kafka required directories
if [[ ! -d /etc/kafka ]]; then
        kafka_output=$(sudo mkdir -p /etc/kafka/conf /etc/kafka/data 2>&1)
        eval_cmd_exec $? "failed to create kafka directories" "$kafka_output"
fi

kafka_owner_output=$(sudo chown -R odimra:odimra /etc/kafka* && sudo chmod 0755 /etc/kafka* 2>&1)
eval_cmd_exec $? "failed to modify kafka directories permission" "$kafka_owner_output"

#create zookeeper required directories
if [[ ! -d /etc/zookeeper ]]; then
        zookeeper_output=$(sudo mkdir -p /etc/zookeeper/conf /etc/zookeeper/data /etc/zookeeper/data/log 2>&1)
        eval_cmd_exec $? "failed to create zookeeper directories" "$zookeeper_output"
fi

zookeeper_owner_output=$(sudo chown -R odimra:odimra /etc/zookeeper* && sudo chmod 0755 /etc/zookeeper* 2>&1)
eval_cmd_exec $?  "failed to modify zookeeper directories permission" "$zookeeper_owner_output"

#create odimra required directories
if [[ ! -d /etc/odimracert ]]; then
        odimra_output=$(sudo mkdir -p /etc/odimracert 2>&1)
        eval_cmd_exec $? "failed to create odimra directories" "$odimra_output"
fi

odimra_owner_output=$(sudo chown -R odimra:odimra /etc/odimracert* && sudo chmod 0755 /etc/odimracert* 2>&1)
eval_cmd_exec $? "failed to modify odimra directories permission" "$odimra_owner_output"

#create plugin required directories
if [[ ! -d /etc/plugincert ]]; then
        plugin_output=$(sudo mkdir -p /etc/plugincert 2>&1)
        eval_cmd_exec $? "failed to create plugin directories" "$plugin_output"
fi

plugin_owner_output=$(sudo chown -R plugin:plugin /etc/plugincert* && sudo chmod 0755 /etc/plugincert* 2>&1)
eval_cmd_exec $? "failed to modify plugin directories permission" "$plugin_owner_output"

#create etcd required directories
if [[ ! -d /etc/etcd/data ]] || [[ ! -d /etc/etcd/conf ]]; then
        etcd_output=$(sudo mkdir -p /etc/etcd/data /etc/etcd/conf 2>&1)
        eval_cmd_exec $? "failed to create etcd directories" "$etcd_output"
fi
etcd_owner_output=$(sudo chown -R odimra:odimra /etc/etcd* && sudo chmod 0755 /etc/etcd* 2>&1)
eval_cmd_exec $?  "failed to modify etcd directories permission" "$etcd_owner_output"

# copy certificates and keys
sudo cp rootCA.crt odimra_server.crt odimra_server.key odimra_rsa.public odimra_rsa.private odimra_kafka_client.crt odimra_kafka_client.key /etc/odimracert/
sudo cp kafka.keystore.jks kafka.truststore.jks /etc/kafka/conf/
sudo cp zookeeper.keystore.jks zookeeper.truststore.jks /etc/zookeeper/conf/
sudo cp rootCA.crt odimra_server.crt odimra_server.key odimra_kafka_client.crt odimra_kafka_client.key odimra_rsa.public odimra_rsa.private /etc/plugincert/
sudo cp rootCA.crt odimra_etcd_server.crt odimra_etcd_server.key /etc/etcd/conf/

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
	if [ $a -eq 7 ];
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

cd /etc/etcd/conf
if [ $? -eq 0 ];
then
	a=`echo \`ls | wc -l\` `
	if [ $a -eq 3 ];
	then
		echo "etcd Certificates copied successfully"
	else
		echo "Copying of etcd Certificates failed"
		exit -1
	fi
else
	echo "Copying of etcd Certificates failed"
	exit -1
fi

sudo chown -R odimra:odimra /etc/odimracert/
sudo chown -R plugin:plugin /etc/plugincert/
sudo chown -R odimra:odimra /etc/kafka/conf/
sudo chown -R odimra:odimra /etc/zookeeper/conf/
sudo chown -R odimra:odimra /etc/etcd/conf/*
