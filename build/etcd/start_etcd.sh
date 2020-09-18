#!/bin/bash -x

declare PID=0

# prepare config file required for starting etcd server
prepare_config()
{
        . /opt/etcd/scripts/setup_etcd.sh
	configure_etcd_properties
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

prepare_config

start_etcd

create_signal_trap

monitor_etcd

run_forever

exit 0
