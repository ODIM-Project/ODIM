.PHONY: dep copy down-containers

build/odimra/odimra:
	mkdir build/odimra/odimra

COPY = svc-account-session svc-aggregation svc-api svc-events svc-fabrics svc-managers svc-systems svc-task lib-dmtf lib-messagebus lib-persistence-manager lib-utilities plugin-redfish

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
