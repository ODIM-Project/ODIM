#(C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
#Licensed under the Apache License, Version 2.0 (the "License"); you may
#not use this file except in compliance with the License. You may obtain
#a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
#WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
#License for the specific language governing permissions and limitations
# under the License.
.PHONY: dep copy down-containers

build/odimra/odimra:
	mkdir build/odimra/odimra

COPY = svc-account-session svc-aggregation svc-api svc-events svc-fabrics svc-managers svc-systems svc-task svc-update lib-dmtf lib-messagebus lib-persistence-manager lib-utilities plugin-redfish

copy: build/odimra/odimra
	$(foreach var,$(COPY),cp -a $(var) build/odimra/odimra/;)
	cp -f lib-utilities/config/odimra_config.json build/odimra/odimra_config/odimra_config.json
	cp -f plugin-redfish/config/config.json build/RFPlugin/plugin_config/config_redfish_plugin.json
	cp -f lib-messagebus/platforms/platformconfig.toml build/odimra/odimra_config/
	cp -f lib-messagebus/platforms/platformconfig.toml build/RFPlugin/plugin_config/platformconfig.toml
	cp -f lib-utilities/config/schema.json build/odimra/odimra_config/
	cp -f lib-utilities/etc/* build/odimra/odimra_config/registrystore

dep: copy
	build/odimra/makedep.sh

build-containers: dep
	cd build && docker-compose build

standup-containers: build-containers
	cd build && docker-compose up -d  && docker exec -d build_odimra_1 /bin/command.sh && docker restart build_odimra_1 && docker exec -d build_grf_plugin_1 /bin/command.sh && docker restart build_grf_plugin_1

down-containers:
	cd build && docker-compose down

all: standup-containers

clean: 
	build/cleanupbuild.sh
deepclean: 
	build/deepcleanupbuild.sh
