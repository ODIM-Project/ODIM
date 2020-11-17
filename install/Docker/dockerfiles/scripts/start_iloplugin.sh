#!/bin/bash
declare PID=0

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
        trap 'echo "[$(date)] -- INFO  -- SIGTERM received for iloplugin, initiating shut down"; sigterm_handler' SIGTERM
}

# keep the script running till SIGTERM is received
run_forever()
{
        wait
}

start_iloplugin()
{
	cd /bin
	export PLUGIN_CONFIG_FILE_PATH=/etc/iloplugin_config/config.json
	nohup ./plugin-ilo --client_request_timeout=1m >> /var/log/iloplugin_logs/ilo-plugin.log 2>&1 &
	PID=$!
	sleep 2s
}

monitor_process()
{
        while true; do
                pid=$(pgrep plugin-ilo 2> /dev/null)
                if [[ $pid -eq 0 ]]; then
                        echo "plugin-ilo has exited" >> /var/log/iloplugin_logs/ilo-plugin.log 2>&1 &
                        exit 1
                fi
                sleep 5
        done &
}

##############################################
###############  MAIN  #######################
##############################################

start_iloplugin

create_signal_trap

monitor_process

run_forever

exit 0
