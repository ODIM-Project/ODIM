#!/bin/bash
if [ -z "$1" ]
then
    echo "Please provide namespace to uninstall the HELM Chart from"
    exit -1
else
    echo "NameSpace provided is $1"
    echo
fi
if [ -z "$2" ]
then
    echo "HELM Chart will be uninstalled for all ODIMRA Services"
else
    echo "HELM Chart will be uninstalled for $2 Service"
    echo
fi

declare UNDEPLOY=false
LIST="urplugin-pv-pvc urplugin-config urplugin"

if [ -z "$2" ]
then
	for i in $LIST ; do
        	echo "Uninstalling Chart for $i"
        	helm delete $i -n $1
        	if [ $? -eq 0 ]; then
                	echo "$i Chart Uninstalled successful"
                	echo "----------------------------------"
                	echo
       	 	else
                	echo "$i Chart Uninstallation failed"
                	echo "----------------------------------"
                	echo
                	arr1+=$i\;
                	UNDEPLOY=true
        	fi
		sleep 2
	done
else
	echo "Uninstalling Chart for $2"
	helm delete $2 -n $1
	if [ $? -eq 0 ]; then
		echo "$2 Chart Uninstalled successful"
		echo "----------------------------------"
		echo
	else
		echo "$2 Chart Uninstallation failed"
		echo "----------------------------------"
		echo
		exit -1
	fi
fi

if $UNDEPLOY; then
        echo "Uninstallation failed for below Charts"
        echo $arr1
        exit -1;
fi
