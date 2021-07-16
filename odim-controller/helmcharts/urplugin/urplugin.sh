#!/bin/bash

#(C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
#under the License.

# This is an utility script to be executed while deploying plugins
# using odim-controller deployment tool. Which will invoke the 
# script during install/add, upgrade and uninstall/remove plugin.
# odim-controller passes just install/upgrade/uninstall argument
# to the script during corresponding actions, and the script
# should handle anything else required accordingly in these cases

declare ACTION

SCRIPT_DIR=$( dirname "${BASH_SOURCE[0]}" )
CONFIG_FILE_PATH=${SCRIPT_DIR}/urplugin-config.yaml
ODIMRA_USER_NAME=odimra
ODIMRA_GROUP_NAME=odimra

# install should handle any pre-reqs to be performed
# before installing/adding plugin using odim-controller
install()
{
	logpath=$(grep logPath ${CONFIG_FILE_PATH} | cut -d':' -f2 | xargs)
	if [[ $? -ne 0 ]] || [[ -z ${logpath} ]]; then
		echo "[$(date)] -- ERROR -- unable to get configured logPath from ${CONFIG_FILE_PATH}"
		echo "${logpath}"
		exit 1
	fi
	if [[ ! -d ${logpath} ]]; then
		mkdir -p ${logpath}
	else
		echo "[$(date)] -- INFO  -- ${logpath} already exists"
	fi
	chown ${ODIMRA_USER_NAME}:${ODIMRA_GROUP_NAME} ${logpath}
}

# install should handle any pre-reqs to be performed
# before upgrading plugin using odim-controller
upgrade()
{
	return
}

# install should handle any pre-reqs to be performed
# before uninstalling/removing plugin using odim-controller
uninstall()
{
	logpath=$(grep logPath ${CONFIG_FILE_PATH} | cut -d':' -f2 | xargs)
	if [[ $? -ne 0 ]] || [[ -z ${logpath} ]]; then
		echo "[$(date)] -- ERROR -- unable to get configured logPath from ${CONFIG_FILE_PATH}"
		echo "${logpath}"
	fi
	if [[ -d ${logpath} ]]; then
		rm -rf ${logpath}
	else
		echo "[$(date)] -- INFO  -- ${logpath} does not exist"
	fi
}

usage()
{
        echo "$(basename $BASH_SOURCE) install|upgrade|uninstall"
        exit 1
}

##############################################
###############  MAIN  #######################
##############################################
ACTION=$1

if [[ $ACTION =~ "help" ]] ; then
        usage
fi

if [[ ! -f ${CONFIG_FILE_PATH} ]]; then
	echo "[$(date)] -- ERROR -- plugin config file ${CONFIG_FILE_PATH} does not exist"
	exit 1
fi

case $ACTION in
        install)
                install
                ;;
        upgrade)
                upgrade
                ;;
        uninstall)
                uninstall
                ;;
        *)
                usage
                ;;
esac

exit 0
