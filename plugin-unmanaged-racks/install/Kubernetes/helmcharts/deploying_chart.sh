#!/bin/bash
if [ -z "$1" ]
then
    echo "Please provide namespace to install the HELM Chart"
    exit -1
else
    echo "NameSpace provided is $1"
    echo
fi
if [ -z "$2" ]
then
    echo "HELM Chart will be installed for all ODIMRA Services"
else
    echo "HELM Chart will be installed for $2 Service"
    echo
fi

declare VALIDATION=false
declare DRYRUN=false
declare DEPLOY=false
LIST="urplugin"

if [ -z "$2" ]
then
	for i in $LIST ; do
		echo "Validating Chart for $i"
		helm lint ./$i/
		if [ $? -eq 0 ]; then
			echo "$i Chart Validation successful"
			echo "----------------------------------"
			echo
		else
			echo "$i Chart Validation failed"
			echo "----------------------------------"
			echo
			arr1+=$i\;
			VALIDATION=true
		fi
	done
else
	echo "Validating Chart for $2"
	helm lint ./$2/
	if [ $? -eq 0 ]; then
		echo "$2 Chart Validation successful"
		echo "----------------------------------"
		echo
	else
		echo "$2 Chart Validation failed"
		echo "----------------------------------"
		echo
		exit -1
	fi
fi

if $VALIDATION; then 
	echo "Validation failed for below Charts"
	echo $arr1
	exit -1; 
fi

if [ -z "$2" ]
then
	for i in $LIST ; do
        	echo " Performing install dry-run on $i Chart"
		helm install $i ./$i/ --dry-run -n $1
		if [ $? -eq 0 ]; then
			echo "$i Chart install dry-run  successful"
			echo "----------------------------------"
                	echo
        	else
                	echo "$i Chart install dry-run failed"
                	echo "----------------------------------"
                	echo
                	arr2+=$i\;
                	DRYRUN=true
        	fi
		sleep 2
	done
else
	echo " Performing install dry-run on $2 Chart"
	helm install $2 ./$2/ --dry-run -n $1
        if [ $? -eq 0 ]; then
		echo "$2 Chart install dry-run  successful"
                echo "----------------------------------"
		echo
	else
		echo "$2 Chart install dry-run failed"
		echo "----------------------------------"
		echo
		exit -1
	fi
fi

if $DRYRUN; then
        echo "Install DRY RUN failed for below Charts"
        echo $arr2
        exit -1;
fi

if [ -z "$2" ]
then
	for i in $LIST ; do
        	echo " Performing install on $i Chart"
        	helm install $i ./$i/ -n $1
        	if [ $? -eq 0 ]; then
                	echo "$i Chart Installation  successful"
                	echo "----------------------------------"
                	echo
        	else
                	echo "$i Chart Installation failed"
                	echo "----------------------------------"
                	echo
                	arr3+=$i\;
                	DEPLOY=true
        	fi
	done
else
	echo " Performing install on $2 Chart"
	helm install $2 ./$2/ -n $1
	if [ $? -eq 0 ]; then
		echo "$2 Chart Installation  successful"
		echo "----------------------------------"
		echo
	else
		echo "$2 Chart Installation failed"
		echo "----------------------------------"
		echo
		exit -1
	fi
fi

if $DEPLOY; then
        echo "Installation failed for below Charts"
        echo $arr3
        exit -1;
fi
