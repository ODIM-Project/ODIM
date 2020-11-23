#!/bin/bash

KAFKA_CONF_FILE=/opt/kafka/config/server.properties

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
	update_param "log.dirs" ${KAFKA_LOG_DIRS}
	update_param "zookeeper.connect" ${KAFKA_ZOOKEEPER_CONNECT}
	update_param "listeners" ${KAFKA_LISTENERS}
	update_param "advertised.listeners" ${KAFKA_ADV_LISTENERS}
       	update_param "listener.security.protocol.map" "SSL:SSL"
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
		update_param "offsets.topic.replication.factor" "1"
		update_param "transaction.state.log.replication.factor" "1"
		update_param "default.replication.factor" "3"
	fi

	add_modify_footer
}

##############################################
###############  MAIN  #######################
##############################################

configure_kafka_properties

exit 0
