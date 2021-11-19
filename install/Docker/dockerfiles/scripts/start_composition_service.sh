#!/bin/bash
declare PID=0
declare OWN_PID=$$

sigterm_handler()
{
        if [[ $PID -ne 0 ]]; then
                sleep 1
                kill -9 $PID
                wait "$PID" 2>/dev/null
        fi
        exit 0
}

# create a signal trap
create_signal_trap()
{
        trap 'echo "[$(date)] -- INFO  -- SIGTERM received for composition service, initiating shut down"; sigterm_handler' SIGTERM
}

# keep the script running till SIGTERM is received
run_forever()
{
        wait
}

start_composition_service()
{
        registry_address="etcd:2379"
	export PLUGIN_CONFIG_FILE_PATH=/etc/csplugin_config/config.json
	nohup /bin/svc-composition-service --registry=etcd --registry_address=${registry_address} --server_address=csplugin:45100 --client_request_timeout=`expr $(cat $CONFIG_FILE_PATH | grep SouthBoundRequestTimeoutInSecs | cut -d : -f2 | cut -d , -f1 | tr -d " ")`s  >> /var/log/csplugin_logs/csplugin.log 2>&1 &

	PID=$!
	sleep 3

	nohup /bin/add-hosts -file /tmp/host.append >> /var/log/csplugin_logs/csplugin.log 2>&1 &
}


##############################################
###############  MAIN  #######################
##############################################

start_composition_service

create_signal_trap

run_forever

exit 0
