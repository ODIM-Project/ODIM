#!/bin/bash -x
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

declare PID=0
declare ETCD_CONF_FILE=/opt/etcd/bin/etcd.conf

update_param()
{
	config_param=$1
	config_value=$2
	config_param_actual_name=$3

	if [[ -z $config_value ]]; then
		return
	fi

	grep -q "^${config_param}=${config_value}" $ETCD_CONF_FILE
	if [[ $? -eq 0 ]]; then
		return
	fi

	if [[ -z ${config_param_actual_name} ]]; then
		sed -i "s%${config_param}:.*%${config_param}: ${config_value}%" $ETCD_CONF_FILE
	else
		sed -i "s%${config_param}:.*%${config_param_actual_name}: ${config_value}%" $ETCD_CONF_FILE
	fi
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- Failed to update ${config_param}"
		exit 1
	fi
}

configure_etcd_properties()
{
	update_param "name" "${MEMBER_NAME}"
	update_param "data-dir" "${DATA_DIR_PATH}"
	update_param "listen-peer-urls" "${LISTEN_PEER_ADDR}"
	update_param "listen-client-urls" "${LISTEN_CLIENT_ADDR}"
	update_param "initial-advertise-peer-urls" "${INITIAL_ADV_ADDR}"
	update_param "initial-cluster" "${INITIAL_CLUSTER}"
	update_param "initial-cluster-state" "${INITIAL_CLUSTER_STATE}"
	update_param "initial-cluster-token" "${INITIAL_CLUSTER_TOKEN}"
	update_param "advertise-client-urls" "${ADV_CLIENT_ADDR}"
	update_param "log-level" "debug"
	if [[ ${ENABLE_TLS} == "true" ]]; then
		update_param "client-cert-file" "${CLIENT_CERT_FILE}" "cert-file"
		update_param "client-key-file" "${CLIENT_KEY_FILE}" "key-file"
		update_param "client-trusted-ca-file" "${CA_FILE}" "trusted-ca-file"
		update_param "client-cert-auth" "true"
		#update_param "auto-tls" "true"

		update_param "peer-cert-file" "${SERVER_CERT_FILE}" "cert-file"
		update_param "peer-key-file" "${SERVER_KEY_FILE}" "key-file"
		update_param "peer-trusted-ca-file" "${CA_FILE}" "trusted-ca-file"
		update_param "peer-cert-auth" "true" "client-cert-auth"
		#update_param "peer-auto-tls" "true" "auto-tls"
	fi
}

# start etcd server and push it to background
start_etcd()
{
	sleep 5
	/opt/etcd/bin/etcd --config-file /opt/etcd/bin/etcd.conf &
	PID="$!"
}

# handler for SIGTERM signal
# on receiving the signal, will initiate graceful shutdown of etcd
sigterm_handler()
{
	if [[ $PID -ne 0 ]]; then
		# will wait for other instances to gracefully announce quorum exit
		sleep 5

		kill -SIGTERM $PID
    		wait "$PID" 2>/dev/null
  	fi
  	exit 0
}

# create a signal trap
create_signal_trap()
{
	trap 'echo "[$(date)] -- INFO  -- SIGTERM received for etcd, initiating shut down"; sigterm_handler' SIGTERM
}

# monitor etcd process, and if has exited, let us exit here
monitor_etcd()
{
	count=0
	while true; do
		count=$(ps -eaf | grep ${PID} | grep -v grep | wc -l 2>/dev/null)
		if [[ $count -eq 0 ]]; then
			echo "[$(date)] -- ERROR -- etcd has exited or restarted, stopping container"
			exit 1
		fi
		sleep 5
	done &
}

# keep the script running till SIGTERM is received
run_forever()
{
	wait
}

##############################################
###############  MAIN  #######################
##############################################

configure_etcd_properties

start_etcd

create_signal_trap

monitor_etcd

run_forever

exit 0
