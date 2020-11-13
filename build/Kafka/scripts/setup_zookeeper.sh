#!/bin/bash

ZOOKEEPER_CONF_FILE=/opt/zookeeper/config/zookeeper.properties

add_modify_header()
{
	grep -q "^#Customised configuration -- BEGIN" $ZOOKEEPER_CONF_FILE \
	&& sed -i "/^#Customised configuration -- END/d" $ZOOKEEPER_CONF_FILE \
	|| echo -e "\n#Customised configuration -- BEGIN" >> $ZOOKEEPER_CONF_FILE
}

add_modify_footer()
{
	echo "#Customised configuration -- END" >> $ZOOKEEPER_CONF_FILE
}

update_param()
{
	config_param=$1
	config_value=$2

	if [[ -z $config_value ]]; then
		return
	fi

	grep -q "^${config_param}=${config_value}" $ZOOKEEPER_CONF_FILE
	if [[ $? -eq 0 ]]; then
		return
	fi

	sed -i "s%^${config_param}=%#&%" $ZOOKEEPER_CONF_FILE \
	&& echo "${config_param}=${config_value}" >> $ZOOKEEPER_CONF_FILE
}

configure_zookeeper_properties()
{
	add_modify_header

	## comment non-ssl client port
	sed -i "s%^clientPort=%#&%" $ZOOKEEPER_CONF_FILE

	#update_param "maxClientCnxns" "10"
	#update_param "syncEnabled" "true"
	#update_param "reconfigEnabled" "false"
	update_param "serverCnxnFactory" "org.apache.zookeeper.server.NettyServerCnxnFactory"
	update_param "clientCnxnSocket" "org.apache.zookeeper.ClientCnxnSocketNetty"
	update_param "secureClientPort" "${ZOOKEEPER_SSL_CLIENT_PORT}"
	update_param "dataDir" "${ZOOKEEPER_DATA_DIR}"
	update_param "dataLogDir" "${ZOOKEEPER_DATA_LOG_DIR}"
	update_param "ssl.protocol" "TLSv1.2"
	update_param "ssl.keyStore.location" "${ZOOKEEPER_KEYSTORE_PATH}"
	update_param "ssl.keyStore.password" "${ZOOKEEPER_KEYSTORE_PASSWORD}"
	update_param "ssl.keyStore.type" "JKS"
	update_param "ssl.trustStore.location" "${ZOOKEEPER_TRUSTSTORE_PATH}"
	update_param "ssl.trustStore.password" "${ZOOKEEPER_TRUSTSTORE_PASSWORD}"
	update_param "ssl.trustStore.type" "JKS"
	#update_param "ssl.hostnameVerification" "false"
	if ${IS_ZOOKEEPER_CLUSTER}; then
		update_param "server.1" "${ZOOKEEPER_SERVER1_NAME}:2888:3888"
		update_param "server.2" "${ZOOKEEPER_SERVER2_NAME}:2888:3888"
		update_param "server.3" "${ZOOKEEPER_SERVER3_NAME}:2888:3888"
		update_param "initLimit" "5"
		update_param "syncLimit" "2"
		#update_param "electionAlg" "3"
		#update_param "standaloneEnabled" "false"
		update_param "sslQuorum" "true"
		update_param "ssl.quorum.keyStore.location" "${ZOOKEEPER_KEYSTORE_PATH}"
		update_param "ssl.quorum.keyStore.password" "${ZOOKEEPER_KEYSTORE_PASSWORD}"
		update_param "ssl.quorum.keyStore.type" "JKS"
		update_param "ssl.quorum.trustStore.location" "${ZOOKEEPER_TRUSTSTORE_PATH}"
		update_param "ssl.quorum.trustStore.password" "${ZOOKEEPER_TRUSTSTORE_PASSWORD}"
		update_param "ssl.quorum.trustStore.type" "JKS"
		update_param "ssl.quorum.protocol" "TLSv1.2"
		update_param "ssl.quorum.hostnameVerification" "false"

		if [[ -z ${ZOOKEEPER_SERVER_ID} ]]; then
			echo "[$(date)] -- ERROR -- Mandatory cluster config param ZOOKEEPER_SERVER_ID not set, exiting"
			exit 1
		fi

		echo "${ZOOKEEPER_SERVER_ID}" > ${ZOOKEEPER_DATA_DIR}/myid
	fi

	add_modify_footer
}

##############################################
###############  MAIN  #######################
##############################################

configure_zookeeper_properties

exit 0
