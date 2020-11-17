#!/bin/bash
declare PID=0

add_host()
{
        /bin/add-hosts -file /tmp/host.append
        if [ $? -ne 0 ]; then
                echo "Appending host entry to /etc/hosts file Failed"
                exit 0
        fi
}

sigterm_handler()
{
        if [[ $PID -ne 0 ]]; then
                # will wait for other instances to gracefully announce quorum exit
                sleep 5
                kill -9 $PID
                wait "$PID" 2>/dev/null
        fi
        exit 0
}

# create a signal trap
create_signal_trap()
{
        trap 'echo "[$(date)] -- INFO  -- SIGTERM received for aggregation, initiating shut down"; sigterm_handler' SIGTERM
}

# keep the script running till SIGTERM is received
run_forever()
{
        wait
}

start_aggregation()
{
	cd /bin
	export CONFIG_FILE_PATH=/etc/odimra_config/odimra_config.json
	nohup ./svc-aggregation --registry=consul --registry_address=consul:8500 --server_address=aggregation:45102 --client_request_timeout=`expr $(cat $CONFIG_FILE_PATH | grep SouthBoundRequestTimeoutInSecs | cut -d : -f2 | cut -d , -f1 | tr -d " ")`s >> /var/log/odimra_logs/aggregation.log 2>&1 &
	PID=$!
	sleep 2s
}

monitor_process()
{
        while true; do
                pid=$(pgrep svc-aggregation 2> /dev/null)
                if [[ $pid -eq 0 ]]; then
                        echo "svc-aggregation has exited"
                        exit 1
                fi
                sleep 5
        done &
}

##############################################
###############  MAIN  #######################
##############################################

add_host

start_aggregation

create_signal_trap

monitor_process

run_forever

exit 0
