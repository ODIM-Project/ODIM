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

KAFKA_CONF_FILE=/opt/kafka/config/server.properties
PID=0

# start kafka server and push it to background
start_kafka()
{
	sleep 5
	/bin/bash /opt/kafka/bin/kafka-server-start.sh /opt/kafka/config/server.properties &
	PID="$!"
}

# handler for SIGTERM signal
# on receiving the signal, will initiate graceful shutdown of kafka
sigterm_handler()
{
	if [[ $PID -ne 0 ]]; then
		# will wait for other instances to gracefully announce quorum exit
		sleep 5

		/bin/bash /opt/kafka/bin/kafka-server-stop.sh
    		wait "$PID" 2>/dev/null
  	fi
  	exit 0
}

# create a signal trap
create_signal_trap()
{
	trap 'echo "[$(date)] -- INFO  -- SIGTERM received for kafka, initiating shut down"; sigterm_handler' SIGTERM
}

# keep the script running till SIGTERM is received
run_forever()
{
	wait
}

add_modify_header()
{
	grep -q "^#Customised configuration -- BEGIN" $KAFKA_CONF_FILE \
	&& sed -i "/^#Customised configuration -- END/d" $KAFKA_CONF_FILE \
	|| echo -e "\n#Customised configuration -- BEGIN" >> $KAFKA_CONF_FILE
}

add_modify_footer()
{
	echo "#Customised configuration -- END" >> $KAFKA_CONF_FILE
}

update_param()
{
	config_param=$1
	config_value=$2

	if [[ -z $config_value ]]; then
		return
	fi

	grep -q "^${config_param}=${config_value}" $KAFKA_CONF_FILE
	if [[ $? -eq 0 ]]; then
		return
	fi

	sed -i "s%^${config_param}=%#&%" $KAFKA_CONF_FILE \
	&& echo "${config_param}=${config_value}" >> $KAFKA_CONF_FILE
}

configure_kafka_properties()
{
	add_modify_header

	update_param "broker.id" "-1"
	update_param "auto.create.topics.enable" "true"
	update_param "num.partitions" "50"
	update_param "log.dirs" ${KAFKA_LOG_DIRS}
	update_param "zookeeper.connect" ${KAFKA_ZOOKEEPER_CONNECT}
	update_param "listeners" ${KAFKA_LISTENERS}
	update_param "advertised.listeners" ${KAFKA_ADV_LISTENERS}
	update_param "listener.security.protocol.map" "SSL:SSL,EXTERNAL:SSL"
	update_param "security.inter.broker.protocol" "SSL"
	update_param "ssl.enabled.protocols" "TLSv1.2"
	update_param "ssl.keystore.location" ${KAFKA_KEYSTORE_PATH}
	update_param "ssl.keystore.password" ${KAFKA_KEYSTORE_PASSWORD}
	update_param "ssl.keystore.type" "JKS"
	update_param "ssl.truststore.location" ${KAFKA_TRUSTSTORE_PATH}
	update_param "ssl.truststore.password" ${KAFKA_TRUSTSTORE_PASSWORD}
	update_param "ssl.truststore.type" "JKS"
	update_param "ssl.client.auth" ${KAFKA_CLIENT_AUTH}
	update_param "zookeeper.ssl.client.enable" "true"
	update_param "zookeeper.clientCnxnSocket" "org.apache.zookeeper.ClientCnxnSocketNetty"
	update_param "zookeeper.ssl.protocol" "TLSv1.2"
	update_param "zookeeper.ssl.keystore.location" ${KAFKA_KEYSTORE_PATH}
	update_param "zookeeper.ssl.keystore.password" ${KAFKA_KEYSTORE_PASSWORD}
	update_param "zookeeper.ssl.keystore.type" "JKS"
	update_param "zookeeper.ssl.truststore.location" ${KAFKA_TRUSTSTORE_PATH}
	update_param "zookeeper.ssl.truststore.password" ${KAFKA_TRUSTSTORE_PASSWORD}
	update_param "zookeeper.ssl.truststore.type" "JKS"
	if $IS_KAFKA_CLUSTER; then
		update_param "min.insync.replicas" "2"
		update_param "default.replication.factor" "3"
		update_param "offsets.topic.replication.factor" "3"
		update_param "transaction.state.log.replication.factor" "3"
	fi

	add_modify_footer
}

##############################################
###############  MAIN  #######################
##############################################

configure_kafka_properties

start_kafka

create_signal_trap

run_forever

exit 0
