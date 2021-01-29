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
        trap 'echo "[$(date)] -- INFO  -- SIGTERM received for grfplugin, initiating shut down"; sigterm_handler' SIGTERM
}

# keep the script running till SIGTERM is received
run_forever()
{
        wait
}

start_grfplugin()
{
	cd /bin
	export PLUGIN_CONFIG_FILE_PATH=/etc/grfplugin_config/config.json
	nohup ./plugin-redfish >> /var/log/grfplugin_logs/grfplugin.log 2>&1 &
	PID=$!
	sleep 2s
  nohup /bin/add-hosts -file /tmp/host.append >> /var/log/grfplugin_logs/add-hosts.log 2>&1 &
}

monitor_process()
{
        while true; do
                pid=$(pgrep plugin-redfish 2> /dev/null)
                if [[ $pid -eq 0 ]]; then
                        echo "plugin-redfish has exited" >> /var/log/grfplugin_logs/grfplugin.log 2>&1 &
                        exit 1
                fi
                sleep 5
        done &
}

##############################################
###############  MAIN  #######################
##############################################

start_grfplugin

create_signal_trap

monitor_process

run_forever

exit 0
