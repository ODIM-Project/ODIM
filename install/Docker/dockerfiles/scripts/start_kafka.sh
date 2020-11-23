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

# prepare config file required for starting kafka server
prepare_config()
{
        /bin/bash /opt/kafka/scripts/setup_kafka.sh
}

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

##############################################
###############  MAIN  #######################
##############################################

prepare_config

start_kafka

create_signal_trap

run_forever

exit 0
