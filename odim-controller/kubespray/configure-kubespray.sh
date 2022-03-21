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

declare KUBESPRAY_SRC_PATH

unpack_kubespray_bundle()
{
	bundle_list=$(ls -r ${KUBESPRAY_SRC_PATH}/kubespray-*.tar.gz)
	if [[ ${#bundle_list[@]} -eq 0 ]]; then
		echo "[$(date)] -- ERROR -- kubespray src bundle not found, exiting"
		exit 1
	elif [[ ${#bundle_list[@]} -ne 1 ]]; then
		echo "[$(date)] -- ERROR -- More than one version of kubespray exists, exiting"
		exit 1
	fi

	output=$(/bin/tar -xzf ${KUBESPRAY_SRC_PATH}/kubespray-*.tar.gz --strip-components 1 -C ${KUBESPRAY_SRC_PATH} > /dev/null 2>&1)
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- kubespray bundle extraction failed"
		echo ${output}
		exit 1
	fi
}

configure_kubespray()
{
	# Disable dashboard and enable helm deployment
	k8s_add_ons_conf_file=${KUBESPRAY_SRC_PATH}/inventory/sample/group_vars/k8s_cluster/addons.yml
	sed -i "s/dashboard_enabled: true/dashboard_enabled: false/; s/helm_enabled: false/helm_enabled: true/" ${k8s_add_ons_conf_file}
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- configuring kubespray features failed"
		exit 1
	fi

	# Change default nginx path
	kubespray_defaults_file=${KUBESPRAY_SRC_PATH}/roles/kubespray-defaults/defaults/main.yaml
	sed -i "s:/etc/nginx:/etc/k8s-nginx:" ${kubespray_defaults_file}
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- changing nginx path used by kubespray failed"
		exit 1
	fi

	# Modify pod_eviction_rate
	k8s_node_defaults_file=${KUBESPRAY_SRC_PATH}/roles/kubernetes/node/defaults/main.yml
	sed -i "s/kubelet_status_update_frequency:.*/kubelet_status_update_frequency: 3s/" ${k8s_node_defaults_file}

	k8s_master_defaults_file=${KUBESPRAY_SRC_PATH}/roles/kubernetes/control-plane/defaults/main/main.yml
	sed -i "s/kube_controller_node_monitor_period:.*/kube_controller_node_monitor_period: 2s/;\
	s/kube_controller_node_monitor_grace_period:.*/kube_controller_node_monitor_grace_period: 15s/;\
	s/kube_kubeadm_apiserver_extra_args:.*/kube_kubeadm_apiserver_extra_args: {default-not-ready-toleration-seconds: '5', default-unreachable-toleration-seconds: '5'}/" ${k8s_master_defaults_file}
}

enable_dualStack()
{
        kubespray_defaults_file=${KUBESPRAY_SRC_PATH}/roles/kubespray-defaults/defaults/main.yaml
        sed -i "s/enable_dual_stack_networks: false/enable_dual_stack_networks: true/" ${kubespray_defaults_file}

	kubespray_k8s_cluster=${KUBESPRAY_SRC_PATH}/inventory/sample/group_vars/k8s_cluster/k8s-cluster.yml
	sed -i "s/enable_dual_stack_networks: false/enable_dual_stack_networks: true/" ${kubespray_k8s_cluster}
}

usage()
{
        echo -e "$(basename $BASH_SOURCE) <kubespray_src_path>"
        exit 1
}

##############################################
###############  MAIN  #######################
##############################################

if [[ $# -lt 1 ]]; then
        usage
fi

KUBESPRAY_SRC_PATH=$1

if [[ -z ${KUBESPRAY_SRC_PATH} ]] || [[ ! -d ${KUBESPRAY_SRC_PATH} ]]; then
        echo "[$(date)] -- ERROR -- invalid kubespray src path [${KUBESPRAY_SRC_PATH}] passed, exiting!!!"
        exit 1
fi

unpack_kubespray_bundle

configure_kubespray

ENABLE_DUALSTACK=$2

if [[ ${ENABLE_DUALSTACK} == "dualStack" ]]; then
	enable_dualStack
fi


exit 0
