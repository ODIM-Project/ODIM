#!/bin/bash
cp -r /var/ur_plugin_config/ /etc/ && rm -rf /var/ur_plugin_config/* && /ur-plugin/start_plugin.sh
