#!/usr/bin/python3

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

import argparse, yaml, logging, traceback
import os, sys, subprocess, grp, time
import glob, shutil, copy, getpass, socket

from yaml import SafeDumper
from Crypto.PublicKey import RSA
from os import path
from logging.handlers import RotatingFileHandler

# global variables
global logger, logger_u, logger_f, lock
logger = None
logger_u = None
logger_f = None
lock = None
CONTROLLER_CONF_DATA = None
CONTROLLER_CONF_FILE = ""
CONTROLLER_LOG_FILE = "odim-controller.log"
CONTROLLER_LOCK_FILE = "/tmp/odim-controller.lock"
DEPLOYMENT_SRC_DIR = ""
KUBESPRAY_SRC_PATH = ""
CONTROLLER_SRC_PATH = ""
CONTROLLER_BASE_PATH = ""
DRY_RUN_SET = False
NO_PROMPT_SET = False
IGNORE_ERRORS_SET = False
K8S_INVENTORY_DATA = None
K8S_INVENTORY_FILE = ""
ODIMRA_VAULT_KEY_FILE = ""
ANSIBLE_SUDO_PW_FILE = ""
ANSIBLE_BECOME_PASS = ""
DEPLOYMENT_ID = ""
ODIMRA_SRC_PATH = ""
ODIMRA_VAULT_BIN = ""
MIN_REPLICA_COUNT = 0
MAX_REPLICA_COUNT = 10
MAX_LOG_FILE_SIZE = 5*1024*1024

# write_node_details is used for creating hosts.yaml required
# for deploying kuberentes cluster using kubespray. hosts.yaml
# is prepared based on the parameters provided in odim-controller conf
def write_node_details():
	global CONTROLLER_CONF_DATA
	logger.debug("Preparing hosts file required for k8s cluster deployment")

	# initialize empty dict with mandatory keys of hosts.yaml
	node_details = {
		'all': {
			'hosts': {},
			'children': {
				'kube-master': {'hosts': {}},
				'kube-node': {'hosts': {}},
				'etcd': {'hosts': {}},
				'k8s-cluster': {'children': {'kube-master': None, 'kube-node': None}},
				'calico-rr': {'hosts': {}}
			}
		}	
	}

	# update node information in host.yamls as provided in odim-controller conf
	for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
		temp_dict = {node : {'ansible_host': attrs['ip'], 'ip':attrs['ip'], 'access_ip':attrs['ip']}}
		node_details['all']['hosts'].update(temp_dict)
		temp_dict = {node: None}
		if attrs["isMaster"]:
			logger.debug("%s(%s) is marked as master node", node, attrs['ip'])
			node_details['all']['children']['kube-master']['hosts'].update(temp_dict)
			node_details['all']['children']['kube-node']['hosts'].update(temp_dict)
			node_details['all']['children']['etcd']['hosts'].update(temp_dict)
		else:
			node_details['all']['children']['kube-node']['hosts'].update(temp_dict)

	# consider None as empty dictionary
	SafeDumper.add_representer(type(None),lambda dumper, value: dumper.represent_scalar(u'tag:yaml.org,2002:null', ''))
	with open('./kube_hosts_details.yaml', 'w') as f:
		yaml.safe_dump(node_details, f, default_flow_style=False)
	
	logger.debug("Hosts file prepared and stored at ./kube_hosts_details.yaml")

# read_conf is used for loading the odim-controller conf
def read_conf():
	global CONTROLLER_CONF_DATA

	if not os.path.isfile(CONTROLLER_CONF_FILE):
		logger.critical("invalid conf file %s passed, exiting!!!", CONTROLLER_CONF_FILE)
		exit(1)

	logger.debug("Reading config file %s", CONTROLLER_CONF_FILE)
	with open(CONTROLLER_CONF_FILE) as f:
		CONTROLLER_CONF_DATA = yaml.load(f, Loader=yaml.FullLoader)

# load existing hosts.yaml that created for the deployment_id
def load_k8s_host_conf():
	global K8S_INVENTORY_DATA, DEPLOYMENT_SRC_DIR, K8S_INVENTORY_FILE
	DEPLOYMENT_SRC_DIR = './inventory/k8s-cluster-' + DEPLOYMENT_ID
	K8S_INVENTORY_FILE = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')

	if not os.path.exists(K8S_INVENTORY_FILE):
		logger.critical("Previous deployment data not found for %s, not an existing deployment deployment", DEPLOYMENT_ID)
		exit(1)

	with open(K8S_INVENTORY_FILE) as f:
		K8S_INVENTORY_DATA = yaml.load(f, Loader=yaml.FullLoader)

# update_ansible_conf is used for updating kubespray's internal
# configuration which will be used for executing ansible-playbook
# commands.
# proxy related information will be updated in group_vars/all/all.yml
def update_ansible_conf():
	http_proxy = ""
	https_proxy = ""
	no_proxy = ""

	if 'httpProxy' in CONTROLLER_CONF_DATA and (CONTROLLER_CONF_DATA['httpProxy'] != "" or CONTROLLER_CONF_DATA['httpProxy'] != None):
		http_proxy = CONTROLLER_CONF_DATA['httpProxy']
	if 'httpsProxy' in CONTROLLER_CONF_DATA and (CONTROLLER_CONF_DATA['httpsProxy'] != "" or CONTROLLER_CONF_DATA['httpsProxy'] != None):
		https_proxy = CONTROLLER_CONF_DATA['httpsProxy']
	if 'noProxy' in CONTROLLER_CONF_DATA and (CONTROLLER_CONF_DATA['noProxy'] != "" or CONTROLLER_CONF_DATA['noProxy'] != None):
		no_proxy = CONTROLLER_CONF_DATA['noProxy']

	env_conf_filepath = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'group_vars/all/all.yml')
	fd = open(env_conf_filepath, "rt")
	fdata = fd.read()
	if http_proxy != "":
		fdata = fdata.replace('# http_proxy: ""', 'http_proxy: "'+http_proxy+'"')
	if https_proxy != "":
		fdata = fdata.replace('# https_proxy: ""', 'https_proxy: "'+https_proxy+'"')
	if no_proxy != "":
		fdata = fdata.replace('# no_proxy: ""', 'no_proxy: "'+no_proxy+'"')
	fd.close()
		
	if http_proxy != "" or https_proxy != "":
		fd = open(env_conf_filepath, "wt")
		fd.write(fdata)
		fd.close()

# perform_checks is used for validating the configuration
# parameters passed to odim-controller.
# For any operation KUBESPRAY_SRC_PATH, odim_controller_path, deployment_id
# are mandatory parameters and optional parameter checks can be skipped by
# passing skip_opt_param_check argument set to True
def perform_checks(skip_opt_param_check=False):
	global KUBESPRAY_SRC_PATH, CONTROLLER_SRC_PATH, CONTROLLER_CONF_DATA, DEPLOYMENT_ID
	global CONTROLLER_BASE_PATH, ANSIBLE_SUDO_PW_FILE, DEPLOYMENT_SRC_DIR, ODIMRA_SRC_PATH
	global ODIMRA_VAULT_BIN, ODIMRA_VAULT_KEY_FILE
	global KUBERNETES_IMAGE_PATH, ODIMRA_IMAGE_PATH

	if 'deploymentID' not in CONTROLLER_CONF_DATA or CONTROLLER_CONF_DATA['deploymentID'] == None or CONTROLLER_CONF_DATA['deploymentID'] == "":
		logger.critical("deployment ID not configured, exiting!!!")
		exit(1)
	DEPLOYMENT_ID = CONTROLLER_CONF_DATA['deploymentID']

	if not skip_opt_param_check:
		logger.debug("Checking if the local user matches with the configired nodes user")
		cur_user = os.getenv('USER')
		for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
			if cur_user != attrs['username']:
				logger.critical("User names of local host and all remote hosts should match")
				exit(1)

	if 'odimControllerSrcPath' not in CONTROLLER_CONF_DATA or \
		CONTROLLER_CONF_DATA['odimControllerSrcPath'] == None or \
		CONTROLLER_CONF_DATA['odimControllerSrcPath'] == "":
		logger.critical("odim-controller source path not configured, exiting!!!")
		exit(1)

	CONTROLLER_BASE_PATH = CONTROLLER_CONF_DATA['odimControllerSrcPath']
	if not os.path.isdir(CONTROLLER_BASE_PATH):
		logger.critical("invalid odim-controller source path configured, exiting!!!")
		exit(1)

	CONTROLLER_SRC_PATH = os.path.join(CONTROLLER_BASE_PATH, 'scripts')
	if not os.path.isdir(CONTROLLER_SRC_PATH):
		logger.critical("%s directory does not exist, exiting!!!", CONTROLLER_SRC_PATH)
		exit(1)

	KUBESPRAY_SRC_PATH = os.path.join(CONTROLLER_BASE_PATH, 'kubespray')
	if not os.path.isdir(KUBESPRAY_SRC_PATH):
		logger.critical("%s directory does not exist, exiting!!!", KUBESPRAY_SRC_PATH)
		exit(1)

	ODIMRA_SRC_PATH = os.path.join(CONTROLLER_BASE_PATH, 'odimra')
	if not os.path.isdir(ODIMRA_SRC_PATH):
		logger.critical("%s directory does not exist, exiting!!!", ODIMRA_SRC_PATH)
		exit(1)

	check_extract_kubespray_src()

	DEPLOYMENT_SRC_DIR = os.path.join(KUBESPRAY_SRC_PATH, 'inventory/k8s-cluster-' + DEPLOYMENT_ID)
	if not os.path.exists(DEPLOYMENT_SRC_DIR):
		os.mkdir(DEPLOYMENT_SRC_DIR, 0o755)

	ODIMRA_VAULT_BIN = os.path.join(CONTROLLER_SRC_PATH, 'odim-vault')
	if not os.path.exists(ODIMRA_VAULT_BIN):
		logger.critical("%s does not exist, exiting!!!", ODIMRA_VAULT_BIN)

	if 'odimVaultKeyFilePath' not in CONTROLLER_CONF_DATA or \
	CONTROLLER_CONF_DATA['odimVaultKeyFilePath'] == None or CONTROLLER_CONF_DATA['odimVaultKeyFilePath'] == "":
		store_vault_key()
	else:
		ODIMRA_VAULT_KEY_FILE = CONTROLLER_CONF_DATA['odimVaultKeyFilePath']

	if 'nodePasswordFilePath' not in CONTROLLER_CONF_DATA or \
	CONTROLLER_CONF_DATA['nodePasswordFilePath'] == None or CONTROLLER_CONF_DATA['nodePasswordFilePath'] == "":
		ANSIBLE_SUDO_PW_FILE = os.path.join(KUBESPRAY_SRC_PATH, 'inventory/k8s-cluster-' + DEPLOYMENT_ID, '.node_pw.dat')
		if not os.path.exists(ANSIBLE_SUDO_PW_FILE):
			store_password_in_vault()
	else:
		ANSIBLE_SUDO_PW_FILE = CONTROLLER_CONF_DATA['nodePasswordFilePath']
		if not os.path.exists(ANSIBLE_SUDO_PW_FILE):
			logger.critical("%s does not exist, exiting!!!", ANSIBLE_SUDO_PW_FILE)

	cert_dir = os.path.join(CONTROLLER_SRC_PATH, 'certs')
	if not os.path.exists(cert_dir):
		os.mkdir(cert_dir, 0o700)
		
	if 'kubernetesImagePath' not in  CONTROLLER_CONF_DATA or \
        CONTROLLER_CONF_DATA['kubernetesImagePath'] == None or CONTROLLER_CONF_DATA['kubernetesImagePath'] == "":
		logger.info(" Kubernetes Image directory not provided, required images will be downloaded!!!")
		KUBERNETES_IMAGE_PATH=""
	else:
		KUBERNETES_IMAGE_PATH =  CONTROLLER_CONF_DATA['kubernetesImagePath']
		if not os.path.exists(KUBERNETES_IMAGE_PATH):
                        logger.warning("%s does not exist, required images will be downloaded!!!", KUBERNETES_IMAGE_PATH)
        
	if 'odimraImagePath' not in CONTROLLER_CONF_DATA or \
                CONTROLLER_CONF_DATA['odimraImagePath'] == None or \
                CONTROLLER_CONF_DATA['odimraImagePath'] == "":
		logger.warning("odimra image source path not configured, expecting user to copy & load all the required odimra docker images on cluster nodes !!!")
		ODIMRA_IMAGE_PATH=""
	else:
		ODIMRA_IMAGE_PATH = CONTROLLER_CONF_DATA['odimraImagePath']
		if not os.path.isdir(ODIMRA_IMAGE_PATH):
			logger.critical("invalid odimra image source path configured, exiting!!!")
			exit(1)
		
# exec is used for executing shell commands.
# It accepts the command to be executed and environment
# variables to set in the form of dictionary.
# It returns command exit code of the command execution
def exec(cmd, set_env):
	cmd_env = os.environ.copy()
	cmd_env.update(set_env)

	execHdlr = subprocess.Popen(cmd,
			env=cmd_env,
			stdin=subprocess.PIPE,
			stdout=subprocess.PIPE,
			stderr=subprocess.STDOUT,
			shell=True,
			universal_newlines=True)

	for output in execHdlr.stdout:
		logger_u.info(output.strip())

	try:
		std_out, std_err = execHdlr.communicate()
	except TimeoutExpired:
		execHdlr.kill()

	return execHdlr.returncode

# copy_ssh_keys_remote_host is used for copying
# ssh keys to remote nodes to enable password-less
# login to those nodes provided in odim-controller conf
def copy_ssh_keys_remote_host():
	cur_user = os.getenv('USER')

	for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
		logger.debug("Enabling password-less login to %s(%s)", node, attrs['ip'])
		sync_cmd = '/usr/bin/sshpass -e /usr/bin/ssh-copy-id -o StrictHostKeyChecking=no -i {conf_path} {username}@{node_ip}'.format(
				conf_path=os.path.join(os.getenv('HOME'), '.ssh/id_rsa.pub'),
				username=attrs['username'],
				node_ip=attrs['ip'])

		ret = exec(sync_cmd, {'SSHPASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
			logger.critical("Enabling password-less login to %s(%s) failed", node, attrs['ip'])
			exit(1)

# gen_ssh_keys is used for generating ssh keys on the local node
# if not present, required for enabling password-less login
# to configured remote nodes.
def gen_ssh_keys():
	ssh_keys_dir = os.path.join(os.getenv('HOME'), '.ssh')
	ssh_priv_key_path = os.path.join(ssh_keys_dir, 'id_rsa')
	ssh_pub_key_path = os.path.join(ssh_keys_dir, 'id_rsa.pub')

	if not os.path.exists(os.path.join(ssh_keys_dir)):
		os.mkdir(ssh_keys_dir, mode = 0o700)

	privkey = RSA.generate(2048)
	with open(ssh_priv_key_path, 'wb') as f:
		os.chmod(ssh_priv_key_path, 0o600)
		f.write(privkey.exportKey('PEM'))
	
	pubkey = privkey.publickey()
	with open(ssh_pub_key_path, 'wb') as f:
		os.chmod(ssh_pub_key_path, 0o640)
		f.write(pubkey.exportKey('OpenSSH'))

# enable_passwordless_login is used for enabling password-less
# login from local node to configured remote nodes.
def enable_passwordless_login():
	if not os.path.exists(os.path.join(os.getenv('HOME'), '.ssh/id_rsa.pub')):
		logger.info("SSH keys does not exist, generating now")
		gen_ssh_keys()
	
	copy_ssh_keys_remote_host()

# dup_dir is used for duplicating the directory contents
def dup_dir(src, dest):
	# if source is a directory, create destination
	# directory and copy each files
	if os.path.isdir(src):
		if not os.path.isdir(dest):
			os.mkdir(dest, 0o755)
		file_list = glob.glob(src + '/*')
		for file in file_list:
			dup_dir(file, dest + '/' + file.split('/')[-1])
	else:
		shutil.copy(src, dest)

# helper_msg is used for logging any message
# to help the user with next steps or hints
def helper_msg():
	logger.info("Perform below steps to enable current user to use kubectl")
	print("""
--- mkdir -p $HOME/.kube
--- sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
--- sudo chown $(id -u):$(id -g) $HOME/.kube/config
""")

# check_time_sync is used for checking if time in all nodes
# provided in the configuration is in sync.
def check_time_sync():
	logger.info("Checking if time on all nodes provided are in sync")
	host_time_map = {}
	# fetch date and time from any one of the master node, if not new deployment
	if K8S_INVENTORY_DATA != None:
		for node, attrs in K8S_INVENTORY_DATA['all']['hosts'].items():
			cmd = '/usr/bin/ssh {username}@{ipaddr} date'.format(username=os.getenv('USER'), ipaddr=attrs['ip'])
			host_time_map[node] = subprocess.run(cmd, shell=True, stdout=subprocess.PIPE).stdout.decode('utf-8').rstrip('\n')
			break

	# fetch date and time from each of the node configured for k8s deployment
	for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
		cmd = '/usr/bin/ssh {username}@{ipaddr} date'.format(username=attrs['username'], ipaddr=attrs['ip'])
		host_time_map[node] = subprocess.run(cmd, shell=True, stdout=subprocess.PIPE).stdout.decode('utf-8').rstrip('\n')

	baseTimeInfo = []
	for host, time in host_time_map.items():
		logger.debug("Timestamp fetched from [%s] is [%s]", host, time)
		itemCount = len(baseTimeInfo)
		if itemCount == 0:
			baseTimeInfo = time.split(' ')
			continue

		#['Wed', 'Aug', '12', '09:47:14', 'UTC', '2020']
		timeToCompare = time.split(' ')
		i = 0
		while i < itemCount:
			# If next element to be parsed is time, will find
			# it by looking for ':' substring
			if (baseTimeInfo[i].find(':') == -1):
				if baseTimeInfo[i] != timeToCompare[i]:
					logger.critical("Time in %s(%s) is not in sync with other nodes", host, time)
					exit(1)
			else:
				# Compare time
				timeStr1 =  baseTimeInfo[i].split(':')
				timeStr2 =  timeToCompare[i].split(':')
				if len(timeStr1) != 3 or len(timeStr2) != 3:
					logger.critical("Timestamp fetched from %s(%s) is not in expected format", host, time)
					exit(1)

				# Compare time by converting hours into minutes and add the elasped minutes too,
				# the difference should not be greater than 1 minute
				baseTimeInMins = (int(timeStr1[0]) * 60) + int(timeStr1[1])
				timeToCompareInMins = (int(timeStr2[0]) * 60) + int(timeStr2[1])
				if timeToCompareInMins - baseTimeInMins > 1:
					logger.critical("Time in %s(%s) is not in sync with other nodes", host, time)
					exit(1)
			i += 1

		baseTimeInfo = time.split(' ')

# scale_in_k8s is for removing nodes from the exisitng k8s cluster
# Configuration passed will be parsed to find the nodes to be removed
# and the existing hosts.yaml created for the deployment_id will be updated.
# kubespray ansible command with remove-node.yaml will be invoked for removing
# the nodes.
def scale_in_k8s():
	cur_dir = os.getcwd()
	os.chdir(KUBESPRAY_SRC_PATH)

	no_nodes_to_remove = True
	existing_nodes = ""
	for node, attrs in K8S_INVENTORY_DATA['all']['hosts'].items():
		existing_nodes += '{hostname}\n'.format(hostname=node)

	# Conf data will be parsed to find the nodes to be removed from cluster,
	# and if found any will updatei the hosts.yaml file and also create a new
	# hosts_rm.yaml just for removal operation.
	rm_nodes = ""
	nodes_list = ""
	for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
		if node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts']:
			logger.warn("%s is master node, removing of which is not allowed, skipping!!!", node)
			continue
		if node in K8S_INVENTORY_DATA['all']['hosts'].keys():
			no_nodes_to_remove = False
			rm_nodes += '{hostname}\n'.format(hostname=node)
			nodes_list += '{hostname},'.format(hostname=node)
			K8S_INVENTORY_DATA['all']['hosts'].pop(node)
			K8S_INVENTORY_DATA['all']['children']['etcd']['hosts'].pop(node, 'No Key found')
			K8S_INVENTORY_DATA['all']['children']['kube-node']['hosts'].pop(node, 'No Key found')
		else:
			logger.info("%s node is not part of the existing cluster, skipped", node)

	if no_nodes_to_remove:
		logger.info("No nodes to remove from the cluster %s, no changes made", DEPLOYMENT_ID)
		return

	nodes_list = nodes_list.rstrip(',')
	logger.info("Current k8s deployment has below nodes\n%s" %(existing_nodes))
	logger.info("Nodes to be removed from the cluster are\n%s" %(rm_nodes))

	if not NO_PROMPT_SET:
		confirmation = input("Confirm to proceed with the scale-in action [y/n]: ")

		if confirmation != 'y' and confirmation != 'n':
			logger.critical("Invalid input, exiting!!!")
			exit(1)

		if confirmation == 'n':
			exit(0)

	if not DRY_RUN_SET:
		load_password_from_vault(cur_dir)
		logger.info("Starting k8s cluster scale-in")
		k8s_rm_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root --extra-vars "node={rm_node_list}" -e delete_nodes_confirmation=yes remove-node.yml'.format(host_conf_file=K8S_INVENTORY_FILE, rm_node_list=nodes_list)
		ret = exec(k8s_rm_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
			logger.critical("k8s cluster scale-in failed")
			os.chdir(cur_dir)
			exit(1)

		# copy controller config file
		helm_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/post-uninstall/files/odim_controller_config.yaml')
		odimra_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/odimra-delete-image/files/odimra_config_values.yaml')
		shutil.copyfile(CONTROLLER_CONF_FILE, helm_config_file)
		shutil.copyfile(CONTROLLER_CONF_FILE, odimra_config_file)

		os.chdir(ODIMRA_SRC_PATH)
		logger.info("Performing post-uninstall action on nodes %s", nodes_list)

		odimra_rm_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
				--extra-vars "host={nodes} helm_config_file={helm_config_file} ignore_err=True" post_uninstall.yaml'.format( \
						host_conf_file=K8S_INVENTORY_FILE, nodes=nodes_list, helm_config_file=CONTROLLER_CONF_FILE)
		ret = exec(odimra_rm_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
			logger.error("Post-uninstall action failed on nodes %s", nodes_list)
		else:
			logger.info("Post-uninstall action was successful on nodes %s", nodes_list)

		delete_k8_images(K8S_INVENTORY_FILE,nodes_list)
		# remove copy of controller config file created
		os.remove(helm_config_file)
		os.remove(odimra_config_file)

		SafeDumper.add_representer(type(None),lambda dumper, value: dumper.represent_scalar(u'tag:yaml.org,2002:null', ''))
		with open(K8S_INVENTORY_FILE, 'w') as f:
			yaml.safe_dump(K8S_INVENTORY_DATA, f, default_flow_style=False)

	os.chdir(cur_dir)
	logger.info("Completed k8s cluster scale-in")

# scale_out_k8s is for adding new nodes to the exisitng k8s cluster
# Configuration passed will be parsed to find the new nodes to be added
# and the existing hosts.yaml created for the deployment_id will be updated.
# kubespray ansible command with scale.yaml will be invoked for deploying
# the new nodes.
def scale_out_k8s():
	# if not HA deployment, restrict scaling out of nodes
	if 'haDeploymentEnabled' not in CONTROLLER_CONF_DATA['odimra'] or \
		not CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled'] or \
		len(K8S_INVENTORY_DATA['all']['hosts']) == 1:
		logger.critical("Scaling out of single node deployment is not allowed")
		return

	cur_dir = os.getcwd()
	os.chdir(KUBESPRAY_SRC_PATH)

	no_new_nodes_to_add = True
	existing_nodes = ""
	for node, attrs in K8S_INVENTORY_DATA['all']['hosts'].items():
		existing_nodes += '{hostname}\n'.format(hostname=node)

	# Conf data will be parsed to find the new nodes, and if found any
	# will update the hosts.yaml file.
	new_nodes = ""
	nodes_list = ""
	for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
		if node not in K8S_INVENTORY_DATA['all']['hosts'].keys():
			no_new_nodes_to_add = False
			new_nodes += '{hostname}\n'.format(hostname=node)
			nodes_list += '{hostname},'.format(hostname=node)
			temp_dict = {node : {'ansible_host': attrs['ip'], 'ip':attrs['ip'], 'access_ip':attrs['ip']}}
			K8S_INVENTORY_DATA['all']['hosts'].update(temp_dict)
			temp_dict = {node: None}
			K8S_INVENTORY_DATA['all']['children']['kube-node']['hosts'].update(temp_dict)

	if no_new_nodes_to_add:
		logger.info("No new nodes to add to cluster %s, no changes made", DEPLOYMENT_ID)
		return

	logger.info("Current k8s deployment has below nodes\n%s" %(existing_nodes))
	logger.info("New nodes to be added are\n%s" %(new_nodes))

	if not NO_PROMPT_SET:
		confirmation = input("Confirm to proceed with the scale-out action [y/n]: ")

		if confirmation != 'y' and confirmation != 'n':
			logger.critical("Invalid input, exiting!!!")
			exit(1)

		if confirmation == 'n':
			exit(0)

	if not DRY_RUN_SET:
		logger.info("Starting k8s cluster scale-out")

		load_password_from_vault(cur_dir)
		# Enable password-less login for the new node
		enable_passwordless_login()
		# Check if the new node time is in sync with other nodes
		check_time_sync()

		SafeDumper.add_representer(type(None),lambda dumper, value: dumper.represent_scalar(u'tag:yaml.org,2002:null', ''))
		with open(K8S_INVENTORY_FILE, 'w') as f:
			yaml.safe_dump(K8S_INVENTORY_DATA, f, default_flow_style=False)
		nodes_list = nodes_list.rstrip(',')

		update_ansible_conf()
                # copy k8 images if provided
		copy_k8_images(K8S_INVENTORY_FILE,nodes_list)
		k8s_add_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root scale.yml'.format( \
				host_conf_file=K8S_INVENTORY_FILE)
		ret = exec(k8s_add_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
			logger.critical("k8s cluster scale-out failed")
			os.chdir(cur_dir)
			exit(1)

		if os.path.exists(os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, '.odimra_deployed')):
			# copy controller config file
			helm_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/pre-install/files/helmcharts/helm_config_values.yaml')
			odimra_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/odimra-copy-image/files/odimra_config_values.yaml')
			shutil.copyfile(CONTROLLER_CONF_FILE, helm_config_file)
			shutil.copyfile(CONTROLLER_CONF_FILE, odimra_config_file)

			os.chdir(ODIMRA_SRC_PATH)
			logger.info("Performing ODIMRA pre-install action nodes %s", nodes_list)

			odimra_add_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
					--extra-vars "host={nodes}" pre_install.yaml'.format( \
							host_conf_file=K8S_INVENTORY_FILE, nodes=nodes_list)
			ret = exec(odimra_add_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
			if ret != 0:
				logger.critical("ODIMRA pre-install action failed on nodes %s", nodes_list)
				os.chdir(cur_dir)
				exit(1)
			else:
				logger.info("ODIMRA pre-install action was successful on nodes %s", nodes_list)

			# remove copy of controller config file created
			os.remove(helm_config_file)
			os.remove(odimra_config_file)

	os.chdir(cur_dir)
	logger.info("Completed k8s cluster scale-out")

# delete_k8_images is used for removing kubernetes
# in the cluster nodes when kubernetesImagePath
# config is set
def delete_k8_images(host_file,nodes_list):
	if KUBERNETES_IMAGE_PATH == "":
		return

	logger.info("Removing k8s images in cluster nodes")
	cur_dir = os.getcwd()
	os.chdir(ODIMRA_SRC_PATH)
	helm_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/k8-delete-image/files/helm_config_values.yaml')
	shutil.copyfile(CONTROLLER_CONF_FILE, helm_config_file)
	k8s_delete_deploy_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root --extra-vars "host={nodes} ignore_err={ignore_err}" k8_delete_image.yaml'.format(host_conf_file=host_file,nodes=nodes_list,ignore_err=IGNORE_ERRORS_SET)
	ret = exec(k8s_delete_deploy_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
	if ret != 0:
		logger.error("k8s image deletion failed, if needed delete images manually")
		os.chdir(cur_dir)
	os.remove(helm_config_file)
	os.chdir(KUBESPRAY_SRC_PATH)

# remove_k8s is used for removing k8s deployment
# from the nodes provided in the odim-controller conf
def remove_k8s():
	cur_dir = os.getcwd()

	os.chdir(KUBESPRAY_SRC_PATH)
	global DEPLOYMENT_SRC_DIR
	DEPLOYMENT_SRC_DIR = './inventory/k8s-cluster-' + DEPLOYMENT_ID
	host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')

	if not os.path.exists(host_file):
		logger.critical("Previous deployment data not found for %s, make sure deployment_id is correct", DEPLOYMENT_ID)
		exit(1)

	with open(host_file) as f:
		host_data = yaml.load(f, Loader=yaml.FullLoader)

	nodes = ""
	nodes_list = ""
	for node, attrs in host_data['all']['hosts'].items():
		nodes += '{hostname}\n'.format(hostname=node)
		nodes_list += '{hostname},'.format(hostname=node)

	logger.info("k8s deployment in below nodes will be reset\n%s" %(nodes))

	if not NO_PROMPT_SET:
		confirmation = input("Confirm to proceed with reset action [y/n]: ")

		if confirmation != 'y' and confirmation != 'n':
			logger.critical("Invalid input, exiting!!!")
			exit(1)

		if confirmation == 'n':
			exit(0)

	if not DRY_RUN_SET:
		load_password_from_vault(cur_dir)
		logger.info("Starting k8s cluster reset")
		k8s_reset_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root -e reset_confirmation=yes reset.yml'.format(host_conf_file=host_file)
		ret = exec(k8s_reset_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
			logger.critical("k8s cluster reset failed")
			os.chdir(cur_dir)
			exit(1)

		delete_k8_images(host_file,nodes_list)
		logger.debug("Clearing deployment specific data of %s cluster" %(DEPLOYMENT_ID))
		shutil.rmtree(DEPLOYMENT_SRC_DIR)

	os.chdir(cur_dir)
	logger.info("Completed k8s cluster reset")

# deploy docker , copy & load k8 images provided 
def copy_k8_images(host_file,nodes_list):
	if KUBERNETES_IMAGE_PATH!="":
		cur_dir = os.getcwd()
		os.chdir(ODIMRA_SRC_PATH)
		helm_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/k8-copy-image/files/helm_config_values.yaml')
		shutil.copyfile(CONTROLLER_CONF_FILE, helm_config_file)
		k8s_image_deploy_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root --extra-vars "host={nodes}" k8_copy_image.yaml'.format(host_conf_file=host_file, nodes=nodes_list)
		ret = exec(k8s_image_deploy_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
		    logger.critical("k8s image deployment failed")
		    os.chdir(cur_dir)
		    exit(1)
		os.remove(helm_config_file)
		os.chdir(KUBESPRAY_SRC_PATH)

# deploy_k8s is used for deploying k8s
# in the nodes provided in odim-controller conf
def deploy_k8s():
	cur_dir = os.getcwd()
	os.chdir(KUBESPRAY_SRC_PATH)

	host_file = os.path.join(DEPLOYMENT_SRC_DIR, 'hosts.yaml')
	if os.path.exists(host_file):
		logger.error("Cluster with deployment ID %s already exists" %(DEPLOYMENT_ID))
		exit(1)

	node_ip_list = ""
	nodes_list = ""
	for node, attrs in CONTROLLER_CONF_DATA['nodes'].items():
		node_ip_list += "%s,%s,%s " %(node, attrs['ip'], attrs['ip'])
		nodes_list += '{hostname},'.format(hostname=node)
	nodes_list = nodes_list.rstrip(',')

	if not DRY_RUN_SET:
		logger.info("Starting k8s cluster deployment")

		load_password_from_vault(cur_dir)
		# Enable password-less login for the new node
		enable_passwordless_login()
		# Check if the new node time is in sync with other nodes
		check_time_sync()

		# replicate the sample inventory data provided by
		# kubespray to create inventory for requested cluster
		dup_dir('./inventory/sample', DEPLOYMENT_SRC_DIR)

		logger.info("Generating hosts file required for k8s cluster deployment")
		host_file_gen_cmd = 'CONFIG_FILE={host_conf_file} python3 contrib/inventory_builder/inventory.py {node_details_list}'.format( \
				host_conf_file=host_file, node_details_list=node_ip_list)

		ret = exec(host_file_gen_cmd, {'KUBE_MASTERS_MASTERS': '3'})
		if ret != 0:
			logger.critical("k8s cluster hosts file generation failed")
			os.chdir(cur_dir)
			exit(1)

		# update proxy info in ansible conf
		update_ansible_conf()
                # Copy K8 images if absolute path for images is provided
		copy_k8_images(host_file,nodes_list)

		k8s_deploy_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root cluster.yml'.format(host_conf_file=host_file)
		ret = exec(k8s_deploy_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret != 0:
			logger.critical("k8s cluster deployment failed")
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)
	logger.info("Completed k8s cluster deployment")

# read_file is for reading a file with read
# only mode and returns the file descriptor
def read_file(filepath):
	return open(filepath, 'r').read()

# represent_yaml_multline_str formats multiline data
# by using '|' charater to denote the same.
# It accepts yaml writer and the data needs to be
# written and returns the formatted data.
def represent_yaml_multline_str(dumper, data):
	if '\n' in data:
		dumper.represent_scalar(u'tag:yaml.org,2002:null', '')
		return dumper.represent_scalar(u'tag:yaml.org,2002:str', data, style='|')
	return dumper.org_represent_str(data)

# reload_odimra_certs is for re-generating missing
# cerificates and keys and reloading the config file
# with latest content
def reload_odimra_certs():
	if 'odimCertsPath' not in CONTROLLER_CONF_DATA or \
	CONTROLLER_CONF_DATA['odimCertsPath'] == None or \
	CONTROLLER_CONF_DATA['odimCertsPath'] == "":
		logger.critical("ODIM-RA certificates path does not exist")
		exit(1)

	cert_dir = CONTROLLER_CONF_DATA['odimCertsPath']
	if os.path.exists(os.path.join(cert_dir, '.gen_odimra_certs.ok')):
		gen_cert_tool = os.path.join(CONTROLLER_SRC_PATH, 'gen_odimra_certs.sh')
		gen_cert_cmd = '/bin/bash {gen_cert_script} {cert_dir} {config_file}'.format(gen_cert_script=gen_cert_tool, cert_dir=cert_dir, config_file=CONTROLLER_CONF_FILE)
		ret = exec(gen_cert_cmd, {})
		if ret != 0:
			logger.critical("ODIM-RA certificate generation failed")
			exit(1)

	load_odimra_certs(True)

# load certificates present at configured path
# to be used for creating k8s secrets
def load_odimra_certs(isUpgrade):
	cert_dir = CONTROLLER_CONF_DATA['odimCertsPath']
	CONTROLLER_CONF_DATA['odimra']['rootCACert'] = read_file(os.path.join(cert_dir, 'rootCA.crt'))
	CONTROLLER_CONF_DATA['odimra']['odimraServerCert'] = read_file(os.path.join(cert_dir, 'odimra_server.crt'))
	CONTROLLER_CONF_DATA['odimra']['odimraServerKey'] = read_file(os.path.join(cert_dir, 'odimra_server.key'))
	CONTROLLER_CONF_DATA['odimra']['odimraKafkaClientCert'] = read_file(os.path.join(cert_dir, 'odimra_kafka_client.crt'))
	CONTROLLER_CONF_DATA['odimra']['odimraKafkaClientKey'] = read_file(os.path.join(cert_dir, 'odimra_kafka_client.key'))
	CONTROLLER_CONF_DATA['odimra']['odimraEtcdServerCert'] = read_file(os.path.join(cert_dir, 'odimra_etcd_server.crt'))
	CONTROLLER_CONF_DATA['odimra']['odimraEtcdServerKey'] = read_file(os.path.join(cert_dir, 'odimra_etcd_server.key'))

	# updating key pair once after deployment is not supported.
	if not isUpgrade:
		CONTROLLER_CONF_DATA['odimra']['odimraRSAPublicKey'] = read_file(os.path.join(cert_dir, 'odimra_rsa.public'))
		CONTROLLER_CONF_DATA['odimra']['odimraRSAPrivateKey'] = read_file(os.path.join(cert_dir, 'odimra_rsa.private'))

	# reload odim-controller conf with cert data
	yaml.SafeDumper.org_represent_str = yaml.SafeDumper.represent_str
	yaml.add_representer(str, represent_yaml_multline_str, Dumper=yaml.SafeDumper)
	with open(CONTROLLER_CONF_FILE, 'w') as f:
		yaml.safe_dump(CONTROLLER_CONF_DATA, f, default_flow_style=False)

# perform pre-requisites required for
# deploying ODIM-RA services
def perform_odimra_deploy_prereqs():
	if 'odimCertsPath' not in CONTROLLER_CONF_DATA or \
	CONTROLLER_CONF_DATA['odimCertsPath'] == None or \
	CONTROLLER_CONF_DATA['odimCertsPath'] == "":
		cert_dir = os.path.join(CONTROLLER_SRC_PATH, 'certs', DEPLOYMENT_ID)
		if not os.path.exists(cert_dir):
			os.mkdir(cert_dir, mode = 0o700)

		CONTROLLER_CONF_DATA['odimCertsPath'] = cert_dir
		gen_cert_tool = os.path.join(CONTROLLER_SRC_PATH, 'gen_odimra_certs.sh')
		gen_cert_cmd = '/bin/bash {gen_cert_script} {cert_dir} {config_file}'.format(gen_cert_script=gen_cert_tool, cert_dir=cert_dir, config_file=CONTROLLER_CONF_FILE)
		ret = exec(gen_cert_cmd, {})
		if ret != 0:
			logger.critical("ODIM-RA certificate generation failed")
			exit(1)
	else:
		if not os.path.isdir(CONTROLLER_CONF_DATA['odimCertsPath']):
			logger.critical("ODIM-RA certificates path does not exist")
			exit(1)

	load_odimra_certs(False)

# perform pre-requisites for HA deployment
def perform_check_ha_deploy():
        write_flag=1
        if 'haDeploymentEnabled' not in CONTROLLER_CONF_DATA['odimra'] or CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled'] == None:
                if len(CONTROLLER_CONF_DATA['nodes']) < 3:
                    logger.warning("Nodes provided for ODIMRA deployment is %s. \
ODIMRA-HA Deployment requires minimum 3 nodes for deployment." %(len(CONTROLLER_CONF_DATA['nodes'])))
                    logger.info("Setting HA Deployment to DISABLED")
                    CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled'] = False
                    HA_DEPLOYMENT = False
                    write_flag=0
                else:
                    logger.info("HA Deployment set to ENABLED")
                    CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled'] = True
                    HA_DEPLOYMENT = True
                    write_flag=0
        elif 'haDeploymentEnabled' in CONTROLLER_CONF_DATA['odimra'] and \
        CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled']:
                if len(CONTROLLER_CONF_DATA['nodes']) < 3:
                    logger.warning("Nodes provided for ODIMRA deployment is %s. \
ODIMRA-HA Deployment requires minimum 3 nodes for deployment." %(len(CONTROLLER_CONF_DATA['nodes'])))
                    logger.info("Setting HA Deployment to DISABLED")
                    CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled'] = False
                    HA_DEPLOYMENT = False
                    write_flag=0
                else:
                    logger.info("HA Deployment set to ENABLED")
        else:
                logger.info("HA Deployment set to DISABLED")
                HA_DEPLOYMENT = False
        if write_flag == 0:
                # reload odim-controller conf with haDeployment param
                yaml.SafeDumper.org_represent_str = yaml.SafeDumper.represent_str
                yaml.add_representer(str, represent_yaml_multline_str, Dumper=yaml.SafeDumper)
                with open(CONTROLLER_CONF_FILE, 'w') as f:
                    yaml.safe_dump(CONTROLLER_CONF_DATA, f, default_flow_style=False)

# operation_odimra is used for deploying/removing ODIMRA
# in the nodes provided in odim-controller conf based on the operation input
def operation_odimra(operation):
	cur_dir = os.getcwd()
	os.chdir(ODIMRA_SRC_PATH)

	host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
	if not os.path.exists(host_file):
		logger.error("Host file not found for deployment id %s" %(DEPLOYMENT_ID))
		exit(1)

	if not DRY_RUN_SET:
		load_password_from_vault(cur_dir)

		# set options based on the operation type
		helm_config_file = ""
		odimra_config_file = ""
		if operation == "install":
			helm_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/pre-install/files/helmcharts/helm_config_values.yaml')
			odimra_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/odimra-copy-image/files/odimra_config_values.yaml')
			perform_odimra_deploy_prereqs()
		elif operation == "uninstall":
			helm_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/post-uninstall/files/odim_controller_config.yaml')
			odimra_config_file = os.path.join(ODIMRA_SRC_PATH, 'roles/odimra-delete-image/files/odimra_config_values.yaml')
		shutil.copyfile(CONTROLLER_CONF_FILE, helm_config_file)
		shutil.copyfile(CONTROLLER_CONF_FILE, odimra_config_file)

		# as rollback of failed operation is not handled yet
		# will try on first master node and exit on failure
		master_node = list(K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].keys())[0]
		logger.info("Starting odimra %s on master node %s", operation, master_node)
		odimra_deploy_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
				    --extra-vars "host={master_node} helm_config_file={helm_config_file} ignore_err={ignore_err}" \
{operation_conf_file}.yaml'.format(host_conf_file=host_file, master_node=master_node, helm_config_file=CONTROLLER_CONF_FILE, \
		operation_conf_file=operation,ignore_err=IGNORE_ERRORS_SET)

		ret = exec(odimra_deploy_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		# remove copy of controller config file created
		os.remove(helm_config_file)
		os.remove(odimra_config_file)

		if ret != 0:
			logger.critical("ODIMRA %s failed on master node %s", operation, master_node)
			os.chdir(cur_dir)
			exit(1)

		if operation == "uninstall":
			if os.path.exists(os.path.join(CONTROLLER_CONF_DATA['odimCertsPath'], '.gen_odimra_certs.ok')):
				logger.info("Cleaning up certificates generated for the deployment")
				shutil.rmtree(CONTROLLER_CONF_DATA['odimCertsPath'])
			if os.path.exists(os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, '.odimra_deployed')):
				os.remove(os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, '.odimra_deployed'))
		if operation == "install":
			deployed_odimra_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, '.odimra_deployed')
			open(deployed_odimra_file, 'w').close()

		logger.info("Completed ODIMRA %s operation", operation)

	os.chdir(cur_dir)


def cleanUp():
	if DEPLOYMENT_SRC_DIR != "":
		path = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR)
		logger.info("Cleaning up temp directory : %s", path)
		shutil.rmtree(path)

# install_k8s is for performing all the necessary steps
# for deploying k8s cluster
def install_k8s():
	logger.info("Installing kubernetes")

	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
        # Check for HA deployment
	perform_check_ha_deploy()
	# Initiate k8s deployment
	deploy_k8s()
	exit(0)

# reset_k8s is for performing all the necessary steps
# for removing k8s from the deployed nodes
def reset_k8s():
	logger.info("Resetting kubernetes")

	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks(skip_opt_param_check=True)
	# Remove k8s from the deployed nodes
	remove_k8s()
	exit(0)

# install_odimra is for performing all the necessary steps for installing ODIMRA
def install_odimra():
	logger.info("Installing ODIMRA")
	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()
	# Initiate ODIMRA deployment
	operation_odimra("install")
	exit(0)

# uninstall_odimra is used for performing all the necessary steps for uninstalling ODIMRA
def uninstall_odimra():
	logger.info("Uninstalling ODIMRA")
	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()
	# Initiate ODIMRA removal
	operation_odimra("uninstall")
	exit(0)

# add_k8s_node is for performing all the necessary steps
# for adding a new node to existing k8s cluster
def add_k8s_node():
	logger.info("Adding new node to existing kubernetes cluster")

	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()
	# Initiate k8s deployment on new nodes
	scale_out_k8s()
	exit(0)

# rm_k8s_node is for performing all the necessary steps
# for removing a node from the existing k8s cluster
def rm_k8s_node():
	logger.info("Removing a node from the existing kubernetes cluster")

	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()
	# Initiate node removal from k8s deployment
	scale_in_k8s()
	exit(0)

# generateRandomAlphaNum geneartes generates a random
# string of requested length containing alphanumeric and
# special characters from the defined set
def generateRandomAlphaNum(length):
	random_char_set = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-{}<>+[]$?@:;()%,'
	return ''.join((random.choice(random_char_set) for i in range(length)))

# store_vault_key checks if vault password file exists,
# if not creates the file by asking user for the password
# else returns without performing any action.
def store_vault_key():
	global ODIMRA_VAULT_KEY_FILE

	user_home = os.getenv("HOME")
	odimra_vault_dir = os.path.join(user_home, '.odimra')
	if not os.path.exists(odimra_vault_dir):
		os.mkdir(odimra_vault_dir, mode = 0o700)

	ODIMRA_VAULT_KEY_FILE = os.path.join(odimra_vault_dir, '.key_dnd.dat')
	if not os.path.exists(ODIMRA_VAULT_KEY_FILE):
		print("\nProvide password for vault")
		pw_from_prompt = lambda: (getpass.getpass('Enter Password: '), getpass.getpass('Confirm Password: '))
		first_pw, second_pw = pw_from_prompt()
		if first_pw != second_pw:
			logger.critical("Passwords provided do not match")
			exit(1)

		fd = open(ODIMRA_VAULT_KEY_FILE, "wb")
		fd.write(first_pw.encode('utf-8'))
		fd.close()

		encode_cmd = '{vault_bin} -encode {key_file}'.format(vault_bin=ODIMRA_VAULT_BIN, key_file=ODIMRA_VAULT_KEY_FILE)
		ret = exec(encode_cmd, {})
		if ret != 0:
			logger.critical("storing vault key failed")
			exit(1)

	return

# store_password_in_vault stores the nodes sudo
# password securely by encrypting using odimra vault
def store_password_in_vault():
	global ANSIBLE_BECOME_PASS

	print("\nProvide sudo password of the nodes")
	pw_from_prompt = lambda: (getpass.getpass('Enter Password: '), getpass.getpass('Confirm Password: '))
	first_pw, second_pw = pw_from_prompt()
	if first_pw != second_pw:
		logger.critical("Passwords provided do not match")
		exit(1)

	fd = open(ANSIBLE_SUDO_PW_FILE, "wb")
	fd.write(first_pw.encode('utf-8'))
	fd.close()

	encrypt_cmd = '{vault_bin} -key {key_file} -encrypt {data_file}'.format(vault_bin=ODIMRA_VAULT_BIN,
			key_file=ODIMRA_VAULT_KEY_FILE, data_file=ANSIBLE_SUDO_PW_FILE)
	ret = exec(encrypt_cmd, {})
	if ret != 0:
		logger.critical("storing node password failed")
		exit(1)

	ANSIBLE_BECOME_PASS = first_pw

# load_password_from_vault loads the sudo password of nodes
# of present cluster securely stored usign ansible vault
def load_password_from_vault(cur_dir):
	global ANSIBLE_BECOME_PASS

	decrypt_cmd = '{vault_bin} -key {key_file} -decrypt {data_file}'.format(vault_bin=ODIMRA_VAULT_BIN,
			key_file=ODIMRA_VAULT_KEY_FILE, data_file=ANSIBLE_SUDO_PW_FILE)

	execHdlr = subprocess.Popen(decrypt_cmd,
			stdin=subprocess.PIPE,
			stdout=subprocess.PIPE,
			stderr=subprocess.STDOUT,
			shell=True,
			universal_newlines=True)

	try:
		std_out, std_err = execHdlr.communicate()
	except TimeoutExpired:
		execHdlr.kill()

	if execHdlr.returncode != 0 or std_out == "":
		print(std_out.strip())
		logger.critical("failed to read node password")
		os.chdir(cur_dir)
		exit(1)

	ANSIBLE_BECOME_PASS = std_out.rstrip('\n')

# check_extract_kubespray_src is used for invoking
# a script, after checking and if not exists, to extract
# kubespary source bundle
def check_extract_kubespray_src():
	if not os.path.isdir(os.path.join(KUBESPRAY_SRC_PATH, "inventory")):
		kubespray_extract_tool = os.path.join(KUBESPRAY_SRC_PATH, 'configure-kubespray.sh')
		kubespray_extract_cmd = '/bin/bash {kubespray_extract_tool} {kubespray_src_path}'.format( \
			kubespray_extract_tool=kubespray_extract_tool, kubespray_src_path=KUBESPRAY_SRC_PATH)
		ret = exec(kubespray_extract_cmd, {})
		if ret != 0:
			logger.critical("Extracting and configuring kubespray failed")
			exit(1)

def read_groupvar():
	global GROUP_VAR_DATA
	group_var_file = ODIMRA_SRC_PATH+'/group_vars/all/all.yaml'
	if not os.path.isfile(group_var_file):
		logger.critical("invalid group_var file %s passed, exiting!!!", group_var_file)
		exit(1)

	logger.debug("Reading group_var file %s", group_var_file)
	with open(group_var_file) as f:
		GROUP_VAR_DATA = yaml.load(f, Loader=yaml.FullLoader)

# upgrade_config_map update the config maps
def upgrade_config_map(config_map_name):
	logger.info("Upgrading config map"+config_map_name)
	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()
	#loading the  group_all yaml and finding helm chart full name
	read_groupvar()
	helm_chart_list=config_map_name.split(",")
	for data in helm_chart_list:
		if data == "all":
			odiraConfigHelmChartData= GROUP_VAR_DATA["odim_pv_pvc_secrets_helmcharts"]
			for helm_chart_name  in odiraConfigHelmChartData:
                                if 'pv-pvc' in helm_chart_name:
                                        continue
                                update_helm_charts(helm_chart_name)

			odimHelmChartData= GROUP_VAR_DATA["odim_svc_helmcharts"]
			for helm_chart_name  in odimHelmChartData:
				update_helm_charts(helm_chart_name)

			thirdPartyHelmCharts=GROUP_VAR_DATA["odim_third_party_helmcharts"]
			for helm_chart_name  in thirdPartyHelmCharts:
                                update_helm_charts(helm_chart_name)

			deploy_plugin('all')

		elif data == "odimra":
			odiraConfigHelmChartData= GROUP_VAR_DATA["odim_pv_pvc_secrets_helmcharts"]
			for helm_chart_name  in odiraConfigHelmChartData:
				if 'pv-pvc' in helm_chart_name:
					continue
				update_helm_charts(helm_chart_name)

			odimHelmChartData= GROUP_VAR_DATA["odim_svc_helmcharts"]
			for helm_chart_name  in odimHelmChartData:
				update_helm_charts(helm_chart_name)

			deploy_plugin('all')

		elif data == 'thirdparty':
			thirdPartyHelmCharts=GROUP_VAR_DATA["odim_third_party_helmcharts"]
			for helm_chart_name in thirdPartyHelmCharts:
				update_helm_charts(helm_chart_name)
		else:
			update_helm_charts(data)

# update_helm_charts is for upgrading the deployed
# helm releases
def update_helm_charts(config_map_name):
	
	optionHelmChartInfo = {
		"odimra-config":"odim_pv_pvc_secrets_helmcharts",
		"odimra-platformconfig":"odim_pv_pvc_secrets_helmcharts",
		"odimra-secret":"odim_pv_pvc_secrets_helmcharts",
		"kafka-secret":"odim_pv_pvc_secrets_helmcharts",
		"zookeeper-secret":"odim_pv_pvc_secrets_helmcharts",
		"configure-hosts":"odim_pv_pvc_secrets_helmcharts",
		"odimra-k8s-access-config":"odim_pv_pvc_secrets_helmcharts",
		"account-session":"odim_svc_helmcharts",
		"aggregation":"odim_svc_helmcharts",
		"api":"odim_svc_helmcharts",
		"events":"odim_svc_helmcharts",
		"fabrics":"odim_svc_helmcharts",
		"telemetry":"odim_svc_helmcharts",
		"managers":"odim_svc_helmcharts",
		"systems":"odim_svc_helmcharts",
                "task":"odim_svc_helmcharts",
		"update":"odim_svc_helmcharts",
		"kafka":"odim_third_party_helmcharts",
		"zookeeper":"odim_third_party_helmcharts",
		"redis":"odim_third_party_helmcharts",
		"etcd":"odim_third_party_helmcharts"
	}
	operationHelmChartInfo={
		"odimra-config":"upgrade-config",
		"odimra-platformconfig":"upgrade-config",
		"odimra-secret":"upgrade-config",
		"kafka-secret":"upgrade-config",
		"zookeeper-secret":"upgrade-config",
		"configure-hosts":"upgrade-config",
		"odimra-k8s-access-config":"upgrade-config",
		"account-session":"upgrade-config",
		"aggregation":"upgrade-config",
		"api":"upgrade-config",
		"events":"upgrade-config",
		"fabrics":"upgrade-config",
		"telemetry":"upgrade-config",
		"managers":"upgrade-config",
		"systems":"upgrade-config",
                "task":"upgrade-config",
		"update":"upgrade-config",
		"kafka":"upgrade_thirdparty",
		"zookeeper":"upgrade_thirdparty",
		"redis":"upgrade_thirdparty",
		"etcd":"upgrade_thirdparty"
	}

	if config_map_name not in optionHelmChartInfo:
		logger.critical("%s upgrade is not supported!!!", config_map_name)
		exit(1)

	helmCharatGroupName=optionHelmChartInfo[config_map_name]
	if 'haDeploymentEnabled' in CONTROLLER_CONF_DATA['odimra'] and \
		CONTROLLER_CONF_DATA['odimra']['haDeploymentEnabled'] and \
		helmCharatGroupName == 'odim_third_party_helmcharts':
		helmCharatGroupName = 'odim_third_party_ha_helmcharts'

	operationName=operationHelmChartInfo[config_map_name]
	helmchartData=GROUP_VAR_DATA[helmCharatGroupName]
	fullHelmChartName = helmchartData[config_map_name]
	if fullHelmChartName=='':
		logger.critical("%s upgrade is not supported!!!", config_map_name)
		exit(1)

	logger.info('Full helm chart name %s',fullHelmChartName)
	cur_dir = os.getcwd()
	os.chdir(ODIMRA_SRC_PATH)

	host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
	if not os.path.exists(host_file):
		logger.error("Host file not found for deployment id %s" %(DEPLOYMENT_ID))
		exit(1)

	if not DRY_RUN_SET:
		load_password_from_vault(cur_dir)

		# check if certs needs to be generated or loaded again
		if 'secret' in config_map_name:
			reload_odimra_certs()

		upgrade_flag = False
		if "third_party" in helmCharatGroupName or helmCharatGroupName =='odim_svc_helmcharts':
			if ODIMRA_IMAGE_PATH == "":
				logger.warning("odimra image source path not configured, expecting user to copy & load all the required odimra docker images on cluster nodes !!!")
			else:
				nodes_list = ""
				for node, attrs in K8S_INVENTORY_DATA['all']['hosts'].items():
					nodes_list += '{hostname},'.format(hostname=node)
				nodes_list = nodes_list.rstrip(',')
				dockerImageName=GROUP_VAR_DATA['odim_docker_images'][config_map_name]
				logger.info("Start copying of docker images for %s",config_map_name)
				docker_copy_image_command= 'ansible-playbook -i {host_conf_file} --become --become-user=root \
							   --extra-vars "docker_image_name={docker_image_name} helm_config_file={helm_config_file} host={nodes} ignore_err={ignore_err}" pre_upgrade.yaml'.format(\
									   host_conf_file=host_file,docker_image_name=dockerImageName,\
									   helm_config_file=CONTROLLER_CONF_FILE,\
									   nodes=nodes_list,\
								   ignore_err=IGNORE_ERRORS_SET)
				ret = exec(docker_copy_image_command, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
				if ret != 0:
					logger.critical("ODIMRA %s failed to copy docker image %s", operationName, dockerImageName)
					os.chdir(cur_dir)
					exit(1)
				else:
					logger.info("ODIMRA %s success copy docker image %s", operationName, dockerImageName)


		for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
			logger.info("Starting upgrade of  %s on master node %s", fullHelmChartName, master_node[0])
			odimra_upgrade_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
					     --extra-vars "host={master_node} helm_chart_name={helm_chart_name} helm_chart_name_version={helm_chart_name_version} helm_config_file={helm_config_file} ignore_err={ignore_err}" {operation_conf_file}.yaml'.format( \
							     host_conf_file=host_file, master_node=master_node[0], \
							     helm_chart_name=config_map_name, \
							     helm_chart_name_version=fullHelmChartName, \
							     helm_config_file=CONTROLLER_CONF_FILE, \
							     operation_conf_file=operationName,ignore_err=IGNORE_ERRORS_SET)
			ret = exec(odimra_upgrade_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
			if ret != 0:
				logger.critical("ODIMRA %s failed when tried on master node %s", operationName, master_node[0])
			else:
				logger.info("ODIMRA %s success on master node %s", operationName, master_node[0])
				upgrade_flag=True
				break

		if upgrade_flag:
			logger.info("Completed ODIMRA %s operation", operationName)
		else:
			logger.info("Could not %s ODIMRA on any master nodes", operationName)
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)

# list_deployments is for listing the
# helm deployed releases
def list_deployments():
	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	if 'namespace' not in CONTROLLER_CONF_DATA['odimra'] or \
		CONTROLLER_CONF_DATA['odimra']['namespace'] == None or \
	CONTROLLER_CONF_DATA['odimra']['namespace'] == "":
		logger.critical("namespace not configured, exiting!!!")
		exit(1)

	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()

	list_flag = False
	for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
		ip = K8S_INVENTORY_DATA['all']['hosts'][master_node[0]]['ip']
		list_deps_cmd = '/usr/bin/ssh {ip} helm list -n {namespace}'.format( \
				namespace=CONTROLLER_CONF_DATA['odimra']['namespace'], ip=ip)
		ret = exec(list_deps_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret == 0:
			list_flag = True
			break

	if not list_flag:
		exit(1)

# list_deployment_history is for listing the
# details of a particular helm deployed release
def list_deployment_history(depName):
	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	if 'namespace' not in CONTROLLER_CONF_DATA['odimra'] or \
		CONTROLLER_CONF_DATA['odimra']['namespace'] == None or \
	CONTROLLER_CONF_DATA['odimra']['namespace'] == "":
		logger.critical("namespace not configured, exiting!!!")
		exit(1)

	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()

	list_flag = False
	for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
		ip = K8S_INVENTORY_DATA['all']['hosts'][master_node[0]]['ip']
		list_history_cmd = '/usr/bin/ssh {ip} helm history {deployment} -n {namespace}'.format( \
				ip=ip, deployment=depName, \
				namespace=CONTROLLER_CONF_DATA['odimra']['namespace'])
		ret = exec(list_history_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
		if ret == 0:
			list_flag = True
			break

	if not list_flag:
		exit(1)

# rollback_deployment is for doing rollback of a
# particular helm deployed release
def rollback_deployment(depName, revision):
	logger.info("rollback %s deployment to revision %d", depName, revision)

	# Parse the conf file passed
	read_conf()

	# Validate conf parameters passed
	perform_checks()

	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()

	cur_dir = os.getcwd()
	if not DRY_RUN_SET:
		os.chdir(ODIMRA_SRC_PATH)
		load_password_from_vault(cur_dir) 
		rollback_flag = False

		host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
		for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
			logger.info("Starting rollback of %s deployment on master node %s", depName, master_node[0])
			rollback_dep_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
					   --extra-vars "host={master_node} release={depName} revision={revision}" rollback.yaml'.format( \
							   host_conf_file=host_file, master_node=master_node[0], \
							   depName=depName, revision=revision)
			ret = exec(rollback_dep_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
			if ret != 0:
				logger.critical("rollback of %s deployment failed on master node %s", depName, master_node[0])
			else:
				rollback_flag=True
				break

		if rollback_flag:
			logger.info("rollback of %s deployment to revision %d was successful", depName, revision)
		else:
			logger.info("rollback of %s deployment to revision %d failed", depName, revision)
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)

# scale_plugin is for scaling the helm deployed
# plugin release
def scale_plugin(plugin_name, replica_count):
	logger.info("scaling plugin %s deployment to replicas %d", plugin_name, replica_count)

	# Parse the conf file passed
	read_conf()

	# Validate conf parameters passed
	perform_checks()

	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()

	cur_dir = os.getcwd()
	if not DRY_RUN_SET:
		os.chdir(ODIMRA_SRC_PATH)
		load_password_from_vault(cur_dir) 
		scaling_flag = False

		pluginPackagePath = CONTROLLER_CONF_DATA['odimPluginPath'] + "/" + plugin_name
		if not(path.isdir(pluginPackagePath)):
			logger.error("%s plugin info not present in configured odimPluginPath, scaling not supported", plugin_name)
			exit(1)

		host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
		for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
			logger.info("Starting scaling of %s plugin on master node %s", plugin_name, master_node[0])
			scale_plugin_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
					   --extra-vars "host={master_node} helm_chart_name={helm_chart_name} helm_config_file={helm_config_file} replicas={replicas}" scale_plugin.yaml'.format( \
							   host_conf_file=host_file, master_node=master_node[0], \
							   helm_chart_name=plugin_name, helm_config_file=CONTROLLER_CONF_FILE, \
							   replicas=replica_count)
			ret = exec(scale_plugin_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
			if ret != 0:
				logger.critical("scaling %s plugin failed on master node %s", plugin_name, master_node[0])
			else:
				scaling_flag=True
				break

		if scaling_flag:
			logger.info("scaled %s plugin to %d replica(s)", plugin_name, replica_count)
		else:
			logger.info("failed to scale %s plugin to %d replica(s)", plugin_name, replica_count)
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)

# scale_svc is for scaling the helm deployed
# odim service release
def scale_svc(svc_uservice_name,replica_count):
	logger.info("scaling svc  %s deployment to replicas %d", svc_uservice_name, replica_count)
	# Parse the conf file passed
	read_conf()
	# Validate conf parameters passed
	perform_checks()
	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()
	#loading the  group_all yaml and finding helm chart full name
	read_groupvar()
	helmchartData=GROUP_VAR_DATA["odim_svc_helmcharts"]
	userviceList=svc_uservice_name.split(",")
	for data in userviceList:
		if data=="all":
			for helmChartName in helmchartData:
				scale_svc_helm_chart(helmChartName,replica_count,helmchartData)
		else:
			scale_svc_helm_chart(data,replica_count,helmchartData)
	
def scale_svc_helm_chart(svc_uservice_name,replica_count,helmchartData):
	if svc_uservice_name not in helmchartData:
		logger.critical("scaling of svc %s is not supported!!!", svc_uservice_name)
		exit(1)
	fullHelmChartName=helmchartData[svc_uservice_name]
	logger.info('Full helm chart name %s',fullHelmChartName)
	operationName="scale_svc"
	cur_dir = os.getcwd()
	os.chdir(ODIMRA_SRC_PATH)
	host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
	if not os.path.exists(host_file):
		logger.error("Host file not found for deployment id %s" %(DEPLOYMENT_ID))
		exit(1)

	if not DRY_RUN_SET:
		load_password_from_vault(cur_dir)
		scale_flag = False
		for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
			logger.info("Starting scaling of  %s on master node %s", fullHelmChartName, master_node[0])
			odimra_upgrade_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
					     --extra-vars "host={master_node} helm_chart_name={helm_chart_name} helm_chart_name_version={helm_chart_name_version} helm_config_file={helm_config_file} replicas={replicas} ignore_err={ignore_err}" {operation_conf_file}.yaml'.format( \
							     host_conf_file=host_file, master_node=master_node[0], \
							     helm_chart_name=svc_uservice_name, \
							     helm_chart_name_version=fullHelmChartName, \
							     helm_config_file=CONTROLLER_CONF_FILE, \
                                                             replicas=replica_count, \
							     operation_conf_file=operationName,ignore_err=IGNORE_ERRORS_SET)
			ret = exec(odimra_upgrade_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
			if ret != 0:
				logger.critical("ODIMRA %s failed when tried on master node %s", operationName, master_node[0])
			else:
				logger.info("ODIMRA %s success on master node %s", operationName, master_node[0])
				scale_flag=True
				break

		if scale_flag:
			logger.info("Completed ODIMRA %s operation", operationName)
		else:
			logger.info("Could not %s ODIMRA on any master nodes", operationName)
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)

def deploy_plugin(plugin_name):
	logger.info("Deploy %s plugin", plugin_name)

	# Parse the conf file passed
	read_conf()

	# Validate conf parameters passed
	perform_checks()

	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()

	plugin_list = []
	if plugin_name != 'all':
		pluginPackagePath = CONTROLLER_CONF_DATA['odimPluginPath'] + "/" + plugin_name
		if not(path.isdir(pluginPackagePath)):
			logger.error("%s plugin content not present in configured odimPluginPath, cannot deploy", plugin_name)
			exit(1)
		plugin_list.append(plugin_name)
	else:
		temp_list = []
		for (_, subDirName, _) in os.walk(CONTROLLER_CONF_DATA['odimPluginPath']):
			temp_list.append(subDirName)
			break
		if len(temp_list) <= 0 or len(temp_list[0]) <= 0:
			return

		for item in temp_list[0]:
			plugin_list.append(item)

	cur_dir = os.getcwd()
	host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
	if not DRY_RUN_SET:
		os.chdir(ODIMRA_SRC_PATH)
		load_password_from_vault(cur_dir) 
		plugin_count = 0

		for plugin in plugin_list:
			for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
				logger.info("Starting deployment of %s on master node %s", plugin, master_node[0])
				deploy_plugin_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
						    --extra-vars "host={master_node} release_name={plugin_name} helm_chart_name={helm_chart_name} helm_config_file={helm_config_file}" deploy_plugin.yaml'.format( \
								    host_conf_file=host_file, master_node=master_node[0], \
								    plugin_name=plugin, helm_chart_name=plugin, \
								    helm_config_file=CONTROLLER_CONF_FILE)
				ret = exec(deploy_plugin_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
				if ret != 0:
					logger.critical("deploying %s failed on master node %s", plugin, master_node)
				else:
					plugin_count += 1
					break

		upgrade_failed_count = len(plugin_list) - plugin_count
		if upgrade_failed_count == 0:
			logger.info("Successfully deployed %s", plugin_list)
		else:
			logger.info("Deployment of %d plugin(s) in %s failed", upgrade_failed_count, plugin_list)
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)

def remove_plugin(plugin_name):
	logger.info("remove %s plugin", plugin_name)

	# Parse the conf file passed
	read_conf()

	# Validate conf parameters passed
	perform_checks()

	# load existing hosts.yaml created for the deployment_id
	load_k8s_host_conf()

	cur_dir = os.getcwd()
	host_file = os.path.join(KUBESPRAY_SRC_PATH, DEPLOYMENT_SRC_DIR, 'hosts.yaml')
	pluginPackagePath = CONTROLLER_CONF_DATA['odimPluginPath'] + "/" + plugin_name
	
	if not(path.isdir(pluginPackagePath)):
		logger.info("%s was not deployed via odim controller", plugin_name)
		os.chdir(cur_dir)
		return

	if not DRY_RUN_SET:
		os.chdir(ODIMRA_SRC_PATH)
		load_password_from_vault(cur_dir) 
		upgrade_flag = False

		for master_node in K8S_INVENTORY_DATA['all']['children']['kube-master']['hosts'].items():
			logger.info("Starting removal of %s plugin on master node %s", plugin_name, master_node[0])
			remove_plugin_cmd = 'ansible-playbook -i {host_conf_file} --become --become-user=root \
					   --extra-vars "host={master_node} release_name={plugin_name} helm_chart_name={helm_chart_name} helm_config_file={helm_config_file}" remove_plugin.yaml'.format( \
							   host_conf_file=host_file, master_node=master_node[0], \
							   plugin_name=plugin_name, helm_chart_name=plugin_name, \
							   helm_config_file=CONTROLLER_CONF_FILE)
			ret = exec(remove_plugin_cmd, {'ANSIBLE_BECOME_PASS': ANSIBLE_BECOME_PASS})
			if ret != 0:
				logger.critical("removal of %s plugin failed on master node %s", plugin_name, master_node[0])
			else:
				upgrade_flag=True
				break

		if upgrade_flag:
			logger.info("Successfully removed %s plugin", plugin_name)
		else:
			logger.info("Failed to remove %s plugin", plugin_name)
			os.chdir(cur_dir)
			exit(1)

	os.chdir(cur_dir)

# init_log is for initializing logger for
# logging to console and odim-controller log file
def init_log():
	global logger, logger_u, logger_f

	# check if log path is set and use it for storing
	# log file else create log file in current dir
	oclp_env = os.getenv("ODIM_CONTROLLER_LOG_PATH", "./")
	logPath = os.path.join(oclp_env, CONTROLLER_LOG_FILE)

	# logger is for logging the with the configured log
	# format to both console and the log file
	logger = logging.getLogger('odim-controller')
	logger.setLevel(logging.DEBUG)

	# logger_u is for logging plain logs without log level,
	# timestamp or any other tags(unformatted) to both
	# console and the log file
	logger_u = logging.getLogger('logger_u')
	logger_u.setLevel(logging.DEBUG)

	# logger_f is for logging plain logs without log level,
	# timestamp or any other tags(unformatted) only to log file
	logger_f = logging.getLogger('logger_f')
	logger_f.setLevel(logging.DEBUG)

	# consoleHdlr is the log handler to log to stdout
	consoleHdlr = logging.StreamHandler()
	consoleHdlr.setLevel(logging.DEBUG)

	# fileHdlr is the log handler to log to file
	fileHdlr = RotatingFileHandler(logPath, mode = 'a', maxBytes=MAX_LOG_FILE_SIZE, backupCount=1, encoding=None, delay=0)
	fileHdlr.setLevel(logging.DEBUG)

	# consoleHdlr is the unformatted log handler to log to stdout
	consoleHdlr_u = logging.StreamHandler()
	consoleHdlr_u.setLevel(logging.DEBUG)

	# fileHdlr is the unformatetd log handler to log to file
	fileHdlr_u = RotatingFileHandler(logPath, mode = 'a', maxBytes=MAX_LOG_FILE_SIZE, backupCount=1, encoding=None, delay=0)
	fileHdlr_u.setLevel(logging.DEBUG)

	# logFormatter is for defining the log format, which will contain
	# timestamp, logger name, log level and the log
	logFormatter = logging.Formatter(fmt='%(asctime)s - %(name)s - %(levelname)-5s - %(message)s', datefmt="%Y-%m-%d %H:%M:%S")
	consoleHdlr.setFormatter(logFormatter)
	fileHdlr.setFormatter(logFormatter)

	logger.addHandler(consoleHdlr)
	logger.addHandler(fileHdlr)

	logger_u.addHandler(consoleHdlr_u)
	logger_u.addHandler(fileHdlr_u)

	logger_f.addHandler(fileHdlr_u)

# create_lock is for creating a lock to restrict
# multiple invocation of odim-controller CLI
def create_lock():
	global lock

	lock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
	lock.bind(CONTROLLER_LOCK_FILE)

# lockControllerInvocation is to lock and restrict
# multiple invocation of odim-controller CLI
def lockControllerInvocation():
	try:
		create_lock()
	except socket.error as e:
		# OSError : [Errno 98] Address already in use
		if e.errno == 98:
			logger.error("An instance of odim-controller is already active, another execution is not allowed")
		else:
			logger.error("failed to get lock on odim-controller invocation: %s", str(e))
			logger_f.error('%s', traceback.format_exc())
		sys.exit(1)
	except Exception as e:
		logger.error("failed to get lock on odim-controller invocation: %s", str(e))
		logger_f.error('%s', traceback.format_exc())
		sys.exit(1)

# unlockControllerInvocation is for removing the lock
# and to allow invocation of odim-controller CLI
def unlockControllerInvocation():
	try:
		lock.close()
		os.unlink(CONTROLLER_LOCK_FILE)
	except Exception:
		pass

# exit is for cleaning resouces and formal exit
def exit(code):
	unlockControllerInvocation()
	logger_f.info ("--------- %-7s %s ---------\n", "Ended", time.strftime("%d-%m-%Y %H:%M:%S"))
	sys.exit(code)

def main():
	init_log()

	user = getpass.getuser()
	groups = [grp.getgrgid(g).gr_name for g in os.getgroups()]

	logger_f.info ("--------- %-7s %s ---------", "Started", time.strftime("%d-%m-%Y %H:%M:%S"))
	logger_f.info("%s - %s invoked by user: \"%s\"\n\tgroups: %s\n\toptions: %s",
			time.strftime("%Y-%m-%d %H:%M:%S"), sys.argv[0], user, groups, str(sys.argv[1:]))
	lockControllerInvocation()

	parser = argparse.ArgumentParser(description='ODIM controller')
	parser.add_argument('--deploy', help='supported values: kubernetes, odimra')
	parser.add_argument('--reset', help='supported values: kubernetes, odimra')
	parser.add_argument('--addnode', help='supported values: kubernetes')
	parser.add_argument('--rmnode', help='supported values: kubernetes')
	parser.add_argument('--config', help='absolute path of the config file')
	parser.add_argument('--dryrun', action='store_true', help='only check for configurations without deploying k8s')
	parser.add_argument('--noprompt', action='store_true', help='do not prompt for confirmation')
	parser.add_argument('--ignore-errors', action='store_true', help='ignore errors during odimra reset')
	parser.add_argument("--upgrade", help='supported values:odimra-config,odimra-platformconfig,configure-hosts,odimra-k8s-access-config,odimra-secret,kafka-secret,zookeeper-secret,account-session,aggregation,api,events,fabrics,telemetry,managers,systems,task,update,kafka,zookeeper,redis,etcd,plugin,all,odimra,thirdparty')
	parser.add_argument("--scale", action='store_true', help='scale odimra services and plugins')
	parser.add_argument("--svc", help='supported values:account-session,aggregation,api,events,fabrics,telemetry,managers,systems,task,update,all')
	parser.add_argument("--plugin", help='release name of the plugin deployment to add,remove,upgrade or scale')
	parser.add_argument('--add', help='supported values: plugin')
	parser.add_argument('--remove', help='supported values: plugin')
	parser.add_argument("--replicas", help='replica count of the odimra services or plugins', type=int)
	parser.add_argument('--list', help='supported values:deployment, history')
	parser.add_argument('--dep', help='deployment name, should be used with --list=history, --rollback')
	parser.add_argument('--rollback', action='store_true', help='rollback deployment to particular revision')
	parser.add_argument('--revision', help='revision number of the deployment, should be used with --rollback', type=int)
	try:
		args = parser.parse_args()
	except SystemExit as e:
		exit(1)

	global CONTROLLER_CONF_FILE, DRY_RUN_SET, NO_PROMPT_SET, IGNORE_ERRORS_SET

	if args.deploy == None and args.reset == None and args.addnode == None and \
			args.rmnode == None and args.upgrade == None and args.scale == False and \
			args.list == None and args.add == None and args.remove == None and args.rollback == False:
		logger.critical("Atleast one mandatory option must be provided")
		parser.print_help()
		exit(1)

	if args.dryrun:
		DRY_RUN_SET = True

	if args.noprompt:
		NO_PROMPT_SET = True

	if args.config != None:
		CONTROLLER_CONF_FILE = args.config

	if args.deploy != None:
		if args.deploy == 'kubernetes':
			install_k8s()
		elif args.deploy == 'odimra':
			install_odimra()
		else:
			logger.critical("Unsupported value %s for deploy option", args.deploy)
			parser.print_help()
			exit(1)

	if args.reset != None:
		if args.reset == 'kubernetes':
			if args.ignore_errors:
				IGNORE_ERRORS_SET = True
			reset_k8s()
		elif args.reset == 'odimra':
			if args.ignore_errors:
				IGNORE_ERRORS_SET = True
			uninstall_odimra()
		else:
			logger.critical("Unsupported value %s for reset option", args.reset)
			parser.print_help()
			exit(1)

	if args.addnode != None:
		if args.addnode == 'kubernetes':
			add_k8s_node()
		else:
			logger.critical("Unsupported value %s for addnode option", args.addnode)
			parser.print_help()
			exit(1)

	if args.rmnode != None:
		if args.rmnode == 'kubernetes':
			rm_k8s_node()
		else:
			logger.critical("Unsupported value %s for rmnode option", args.rmnode)
			parser.print_help()
			exit(1)

	if args.upgrade != None:
		if args.upgrade == 'plugin':
			if args.plugin == None:
				logger.error("option --upgrade=plugin: expects --plugin argument")
				exit(1)
			deploy_plugin(args.plugin)
		else:
			upgrade_config_map(args.upgrade)

	if args.add != None:
		if args.add == 'plugin':
			if args.plugin == None:
				logger.error("option --add=plugin: expects --plugin argument")
				exit(1)
			deploy_plugin(args.plugin)
		else:
			logger.critical("Unsupported value %s for add option", args.add)
			exit(1)

	if args.remove != None:
		if args.remove == 'plugin':
			if args.plugin == None:
				logger.error("option --remove=plugin: expects --plugin argument")
				exit(1)
			remove_plugin(args.plugin)
		else:
			logger.critical("Unsupported value %s for remove option", args.remove)
			parser.print_help()
			exit(1)

	if args.scale:
		if args.replicas == None or args.replicas <= MIN_REPLICA_COUNT or args.replicas > MAX_REPLICA_COUNT:
			logger.critical("Unsupported value %d for replicas option", args.replicas)
			exit(1)
		if args.svc != None:
			scale_svc(args.svc, args.replicas)
		elif args.plugin != None:
			scale_plugin(args.plugin, args.replicas)
		else:
			logger.critical("option --scale: expects --svc or --plugin argument")
			parser.print_help()
			exit(1)

	if args.list != None:
		if args.list == 'deployment':
			list_deployments()
		elif args.list == 'history':
			if args.dep == None:
				logger.error("option --history: expects --dep argument")
				exit(1)
			list_deployment_history(args.dep)
		else:
			logger.error("Unsupported value %s for list option", args.list)
			exit(1)

	if args.rollback:
		if args.dep == None or args.revision == None:
			logger.error("option --rollback: expects both --dep and --revision arguments")
			exit(1)
		rollback_deployment(args.dep, args.revision)

if __name__=="__main__":
	try:
		main()
	except KeyboardInterrupt:
		logger.error("Caught interrupt from keyboard, exiting")
		exit(1)
	except Exception as e:
		logger.error("Caught an exception: %s", str(e))
		logger_f.error('%s', traceback.format_exc())
		exit(1)

	exit(0)
