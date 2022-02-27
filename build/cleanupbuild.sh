#!/bin/bash
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
# 
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

if [ -a build/docker-compose.yml ]; then
	cd build
	docker-compose down
	LIST=`docker image ls | grep -v REPOSITORY | grep -vwE 'golang|ubuntu|redis' | awk '{print $3}'`
	docker rmi $LIST
	rm -rf odimra/odimra
	rm -rf odimra/odimra_config/odimra_config.json
        rm -rf odimra/odimra_config/platformconfig.toml
	rm -rf odimra/odimra_config/schema.json
        rm -rf odimra/odimra_config/registrystore/*
        rm -rf RFPlugin/plugin_config/*
	rm -rf DellPlugin/dell_plugin_config/*
	rm -rf LenovoPlugin/lenovo_plugin_config/*
	sudo rm -rf Redis/redis-persistence/*
	sudo rm -rf /etc/kafka/data/*
	sudo rm -rf /etc/zookeeper/data/*
	sudo rm -rf /etc/etcd/data/*
	host=`whoami`
	echo "Cleanup done"
	cd ../
	exit 0
else
	echo "docker-compose.yml doesn't exist, are you in the odimra directory?"
	exit 1
fi
