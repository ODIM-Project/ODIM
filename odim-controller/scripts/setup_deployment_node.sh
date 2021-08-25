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

_if_ne_0_exit_()
{
	if [[ $1 -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- $2"
		exit 1
	fi
}

install_pypkgs()
{
        echo "[$(date)] -- INFO  -- installing required python packages"
        python3 -m pip install -U --proxy=${http_proxy} pip &&
	python3 -m pip install -U --proxy=${http_proxy} setuptools
	_if_ne_0_exit_ $? "failed to install python pip tool"

	sudo -H pip3 install -I --proxy=${http_proxy} pycryptodome==3.9.8 \
	ansible==2.9.6 \
	jinja2==2.11.1 \
	netaddr==0.7.19 \
	pbr==5.4.4 \
	hvac==0.10.0 \
	jmespath==0.9.5 \
	ruamel.yaml==0.16.10 \
	pyyaml==5.3.1

	_if_ne_0_exit_ $? "failed to install required python packages"
}

install_go()
{
	ODIMRA_GO_VERSION="go1.13.7"
	go_bin_path=$(which go 2>&1)
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- INFO  -- installing and configuring ${ODIMRA_GO_VERSION}"

		go_already_installed=false
		if [[ -x /usr/local/go/bin/go ]]; then
			/usr/local/go/bin/go version > /dev/null 2>&1
			if [[ $? -eq 0 ]]; then
				go_already_installed=true
			fi
		fi

		if ! $go_already_installed; then
			wget https://dl.google.com/go/${ODIMRA_GO_VERSION}.linux-amd64.tar.gz -P /var/tmp &&
			sudo tar -C /usr/local -xzf /var/tmp/${ODIMRA_GO_VERSION}.linux-amd64.tar.gz &&
			rm -rf /var/tmp/${ODIMRA_GO_VERSION}.linux-amd64.tar.gz
			_if_ne_0_exit_ $? "failed to install and configure ${ODIMRA_GO_VERSION}"
		fi

		export GO111MODULE=on
		export GOROOT=/usr/local/go
		go env -w GO111MODULE=on
		go env -w GOROOT=/usr/local/go

		echo $PATH | grep "${GOROOT}/bin"
		if [[ $? -ne 0 ]]; then
			export PATH=$PATH:${GOROOT}/bin
			sed "s#.*PATH.*#PATH=\"${PATH}\"#" /etc/environment
		fi
		source /etc/environment

		go_version=$(go version 2>&1)
		echo "[$(date)] -- INFO  -- successfully installed \"${go_version}\""
		return
	fi

	go_version=$(go version 2>&1)
	if [[ $go_version =~ ${ODIMRA_GO_VERSION} ]]; then
		echo "[$(date)] -- INFO  -- go($go_version) already installed at ${go_bin_path}"
	else
		echo "[$(date)] -- ERROR -- go version expected:[go version ${ODIMRA_GO_VERSION} linux/amd64], existing: [${go_version}]"
		echo "[$(date)] -- INFO  -- uninstall exisitng version or install required go version to proceed"
		exit 1
	fi
}

configure_docker_proxy()
{
	sudo sed -i '/no_proxy/d; /http_proxy/d; /https_proxy/d' /etc/environment
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- WARN  -- failed to remove existing proxy entries in /etc/environment"
	fi

	cat << EOF | sudo tee -a /etc/environment
http_proxy=${http_proxy}
https_proxy=${https_proxy}
no_proxy=${no_proxy}
EOF
	_if_ne_0_exit_ $? "failed to update proxy details in /etc/environment"
	source /etc/environment

	sudo mkdir -p /etc/systemd/system/docker.service.d &&
	cat << EOF | sudo tee /etc/systemd/system/docker.service.d/http-proxy.conf
[Service]
Environment="HTTP_PROXY=${http_proxy}"
Environment="HTTPS_PROXY=${https_proxy}"
Environment="NO_PROXY=${no_proxy}"
EOF
	_if_ne_0_exit_ $? "failed to update proxy details in /etc/systemd/system/docker.service.d/http-proxy.conf"

	mkdir -p ~/.docker &&
	sudo chown ${USER}:${USER} ~/.docker -R &&
	sudo chmod 0700 ~/.docker &&
	cat > ~/.docker/config.json <<EOF
{
	"proxies":
	{
		"default":
		{
			"httpProxy": "${http_proxy}",
			"httpsProxy": "${https_proxy}",
			"noProxy": "${no_proxy}"
		}
	}
}
EOF
	_if_ne_0_exit_ $? "failed to update proxy details in ~/.docker/config.json"
}

install_docker()
{
	if [[ -n ${http_proxy} ]] || [[ -n ${https_proxy} ]]; then
		echo "[$(date)] -- INFO  -- http_proxy or https_proxy set, configuring proxy for docker"
		configure_docker_proxy
	fi

	echo "[$(date)] -- INFO  -- installing and configuring docker 5.19 version"

	sudo apt-get update

	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	_if_ne_0_exit_ $? "failed to update docker repo key"

	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" -y
	_if_ne_0_exit_ $? "failed to fetch docker repo"

	sudo apt-get install -y docker-ce=5:19.03.12~3-0~ubuntu-bionic \
	docker-ce-cli=5:19.03.12~3-0~ubuntu-bionic containerd.io --allow-downgrades
	_if_ne_0_exit_ $? "failed to install docker packages"

	cat << EOF | sudo tee /etc/docker/daemon.json
{
	"exec-opts": ["native.cgroupdriver=systemd"],
	"log-driver": "json-file",
	"storage-driver": "overlay2"
}
EOF
	_if_ne_0_exit_ $? "failed to update /etc/docker/daemon.json"

	echo "[$(date)] -- INFO  -- check if docker group exists"
	getent group docker > /dev/null 2>&1
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- INFO  -- creating docker group"
		sudo groupadd docker
		_if_ne_0_exit_ $? "failed to create group docker"
	fi

	echo "[$(date)] -- INFO  -- adding user $USER to docker group"
	sudo usermod -aG docker $USER
	_if_ne_0_exit_ $? "failed to add user $USER to group docker"

	echo "[$(date)] -- INFO  -- enabling and restarting docker systemd service"
	sudo systemctl daemon-reload
	sudo systemctl enable docker && sudo systemctl restart docker
	_if_ne_0_exit_ $? "failed to restart docker systemd service"
}

install_helm()
{
	echo "[$(date)] -- INFO  -- installing and configuring helm3 tool"
	mkdir -p ~/helm &&
	curl -fsSL -o ~/helm/get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | /bin/bash &&
	chmod 0700 ~/helm/get_helm.sh &&
	/bin/bash ~/helm/get_helm.sh
	_if_ne_0_exit_ $? "failed to install helm tool"
}

docker_save()
{
	image_name=$1
	image_version=$2
	output_file=$3

	echo "[$(date)] -- INFO  -- check if ${image_name}:${image_version} image exists"
	docker save ${image_name}:${image_version} -o ${output_file} > /dev/null 2>&1
	if [[ $? -ne 0 ]]; then
		echo "[$(date)] -- INFO  -- ${image_name}:${image_version} does not exist, fetching online"
		docker pull ${image_name}:${image_version}
		_if_ne_0_exit_ $? "failed to pull ${image_name}:${image_version} docker image"
		docker save ${image_name}:${image_version} -o ${output_file}
		_if_ne_0_exit_ $? "failed to save ${image_name}:${image_version} docker image as ${output_file}"
	fi
}

pull_k8s_images()
{
	if [[ -z ${ODIMRA_K8S_IMAGES_PATH} ]]; then
		echo "[$(date)] -- INFO  -- ODIMRA_K8S_IMAGES_PATH not set, not fetching kubernetes images"
		return
	fi

	if [[ ! -d ${ODIMRA_K8S_IMAGES_PATH} ]]; then
		echo "[$(date)] -- ERROR -- ${ODIMRA_K8S_IMAGES_PATH} is not a valid directory"
		exit 1
	fi

	echo "[$(date)] -- INFO  -- ODIMRA_K8S_IMAGES_PATH set, fetching required kubernetes images"
	docker_save "calico/kube-controllers" "v3.15.1"	"${ODIMRA_K8S_IMAGES_PATH}/calico_kube-controllers.tar"
	docker_save "calico/node" "v3.15.1" "${ODIMRA_K8S_IMAGES_PATH}/calico_node.tar"
	docker_save "coredns/coredns" "1.6.7" "${ODIMRA_K8S_IMAGES_PATH}/coredns_coredns.tar"
	docker_save "k8s.gcr.io/cluster-proportional-autoscaler-amd64" \
	"1.8.1" "${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_cluster-proportional-autoscaler-amd64.tar"
	docker_save "k8s.gcr.io/k8s-dns-node-cache" "1.15.13" "${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_k8s-dns-node-cache.tar"
	docker_save "k8s.gcr.io/kube-apiserver" "v1.18.5" "${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_kube-apiserver.tar"
	docker_save "k8s.gcr.io/kube-controller-manager" "v1.18.5" \
	"${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_kube-controller-manager.tar"
	docker_save "k8s.gcr.io/kube-proxy" "v1.18.5" "${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_kube-proxy.tar"
	docker_save "k8s.gcr.io/kube-scheduler" "v1.18.5" "${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_kube-scheduler.tar"
	docker_save "k8s.gcr.io/pause" "3.2" "${ODIMRA_K8S_IMAGES_PATH}/k8s.gcr.io_pause.tar"
	docker_save "lachlanevenson/k8s-helm" "v3.2.3" "${ODIMRA_K8S_IMAGES_PATH}/lachlanevenson_k8s-helm.tar"
	docker_save "nginx" "1.19" "${ODIMRA_K8S_IMAGES_PATH}/nginx.tar"
	docker_save "quay.io/coreos/etcd" "v3.4.3" "${ODIMRA_K8S_IMAGES_PATH}/quay.io_coreos_etcd.tar"
}

pull_odimra_pre_req_images()
{
	if [[ -z ${ODIMRA_IMAGES_PATH} ]]; then
		echo "[$(date)] -- INFO  -- ODIMRA_IMAGES_PATH not set, not fetching ODIMRA deployment required images"
		return
	fi

	if [[ ! -d ${ODIMRA_IMAGES_PATH} ]]; then
		echo "[$(date)] -- ERROR -- ${ODIMRA_IMAGES_PATH} is not a valid directory"
		exit 1
	fi

	echo "[$(date)] -- INFO  -- ODIMRA_IMAGES_PATH set, fetching docker images required for ODIM-RA deployment"
	docker_save "stakater/reloader" "v0.0.76" "${ODIMRA_IMAGES_PATH}/stakater_reloader.tar"
	docker_save "busybox" "1.33" "${ODIMRA_IMAGES_PATH}/busybox.tar"
}

build_odim_images()
{
	if [[ -z ${ODIMRA_SRC_PATH} ]]; then
		echo "[$(date)] -- INFO  -- ODIMRA_SRC_PATH not set, not building ODIMRA images"
		return
	fi

	if [[ -z ${ODIMRA_USER_ID} ]] || [[ -z ${ODIMRA_GROUP_ID} ]]; then
		echo "[$(date)] -- ERROR -- either ODIMRA_USER_ID or ODIMRA_GROUP_ID env var is not set"
		return
	fi

	export ODIMRA_USER_ID=${ODIMRA_USER_ID}
	export ODIMRA_GROUP_ID=${ODIMRA_GROUP_ID}
	cur_dir=$(pwd 2>/dev/null)
	cd ${ODIMRA_SRC_PATH}
	/bin/bash ${ODIMRA_SRC_PATH}/build_images.sh
	eCode=$?
	cd ${cur_dir}
	_if_ne_0_exit_ $eCode "failed to build ODIMRA images"

	if [[ -z ${ODIMRA_IMAGES_PATH} ]]; then
		echo "[$(date)] -- INFO  -- ODIMRA_IMAGES_PATH not set, not storing ODIMRA images"
		return
	fi

	if [[ ! -d ${ODIMRA_IMAGES_PATH} ]]; then
		echo "[$(date)] -- ERROR -- ${ODIMRA_IMAGES_PATH} is not a valid directory"
		exit 1
	fi

	/bin/bash ${ODIMRA_SRC_PATH}/docker-images.sh save ${ODIMRA_IMAGES_PATH}
	_if_ne_0_exit_ $? "failed to store ODIMRA images at ${ODIMRA_IMAGES_PATH}"
}

build_odim_vault()
{
	if [[ -z ${ODIMRA_SRC_PATH} ]]; then
		echo "[$(date)] -- INFO  -- ODIMRA_SRC_PATH not set, not building odim-vault"
		return
	fi

	if [[ ! -d ${ODIMRA_SRC_PATH} ]]; then
		echo "[$(date)] -- ERROR -- ${ODIMRA_SRC_PATH} is not a valid directory"
		exit 1
	fi

	if [[ ! -d ${ODIMRA_SRC_PATH}/odim-controller/scripts ]]; then
		echo "[$(date)] -- ERROR -- ${ODIMRA_SRC_PATH}/odim-controller/scripts does not exist"
		exit 1
	fi

	echo "[$(date)] -- INFO  -- compile and build odim-vault tool"
	cur_dir=$(pwd 2>/dev/null)
	cd ${ODIMRA_SRC_PATH}/odim-controller/scripts
	go build -ldflags "-s -w" -o odim-vault odim-vault.go
	eCode=$?
	cd ${cur_dir}
	_if_ne_0_exit_ $eCode "failed to build odim-vault binary"
}

create_vault_key()
{
	if [[ -z ${ODIMRA_VAULT_KEY_PATH} ]]; then
		echo "[$(date)] -- INFO  -- ODIMRA_VAULT_KEY_PATH not set, not creating odim vault key"
		return
	fi

	if [[ ! -f ${ODIMRA_SRC_PATH}/odim-controller/scripts/odim-vault ]]; then
		echo "[$(date)] -- INFO  -- ${ODIMRA_SRC_PATH}/odim-controller/scripts/odim-vault does not exists, not creating odim vault key"
		return
	fi

	echo "[$(date)] -- INFO  -- encode ${ODIMRA_VAULT_KEY_PATH} odim-vault key"
	${ODIMRA_SRC_PATH}/odim-controller/scripts/odim-vault -encode ${ODIMRA_VAULT_KEY_PATH} &&
	chmod 0600 ${ODIMRA_VAULT_KEY_PATH}
	_if_ne_0_exit_ $? "failed to create odim vault key"
}

usage()
{
        echo -e "$(basename $BASH_SOURCE)"
        exit 1
}

##############################################
###############  MAIN  #######################
##############################################

install_pypkgs

install_go

install_docker

install_helm

pull_k8s_images

pull_odimra_pre_req_images

build_odim_images

build_odim_vault

create_vault_key
