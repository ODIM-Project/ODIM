# Introduction

Organizations today are highly dependent on the manageability of their converged infrastructure, especially as they move towards increasingly complex environments that include multiple remote servers, data storage devices, networking equipment, third-party applications and so on. 

Resource Aggregator for Open Distributed Infrastructure Management \(ODIM™\) is a modular, open framework for centralized management and simplified orchestration of your distributed physical IT infrastructure.

## About the document

This document helps you troubleshoot any common issues you might experience while deploying or using Resource Aggregator for ODIM. Along with this document, Resource Aggregator for ODIM is shipped with the following comprehensive set of electronic documentation:

- **Resource Aggregator for ODIM Getting Started Readme** — This document mainly provides instructions for deploying Resource Aggregator for ODIM and the supported plugins, and covers few typical product use cases.
- **Resource Aggregator for ODIM API Readme** — This document provides detailed information on all the supported APIs of Resource Aggregator for ODIM.

## Conventions

The troubleshooting information is listed in the form of Questions and Answers. You can also find some of the Frequently Asked Questions in this document.

Questions and the associated error messages are in **bold** font. Solutions are in the regular font.

# Troubleshooting Information

This section covers issues you might experience while deploying or using Resource Aggregator for ODIM, and provides suggestions for resolving these issues. You will also find answers to some of the Frequently Asked Questions.


------

**1. The docker start fails with the following error during Kubernetes deployment:<br />`Error log found in journalctl -u docker:`<br />`unable to configure the Docker daemon with file /etc/docker/daemon.json: the following directives are specified both as a flag and in the configuration file: log-opts: (from flag: map[max-file:5 [max-size:50m](http://max-size:50m/)], from file: map[[max-size:100m](http://max-size:100m/)])`**

**Solution**:

   1. Create a file `/etc/systemd/system/docker.service.d/docker.conf `and add the following content in it:
      `[Service]`
      `ExecStart=`
      `ExecStart=/usr/bin/dockerd`
   2. Reset and deploy Kubernetes again.

      Reference links to the issue:
      https://docs.docker.com/config/daemon/#troubleshoot-conflicts-between-the-daemonjson-and-startup-scripts
      https://docs.docker.com/config/daemon/#use-the-hosts-key-in-daemonjson-with-systemd

------

**2. I get `500 Internal Server Error` or `503 Service Unavailable Error` upon sending HTTP requests.**

**Solution**:

Your server encounters unexpected conditions that can prevent it from fulfilling requests due to temporary overloading, session timeout, or any unforeseen reasons.

   1. Run the following command on the master node to verify all deployed services are running successfully:
      `kubectl get pods -n odim -o wide`
      
   2. Navigate to the configured ODIMRA log path and the plugin log path for each server and check the latest logs for any errors.


------

**3. Resetting Resource Aggregator for ODIM deployment or removing the resource aggregator services fails when you use the following command:<br />`python3 odim-controller.py --reset odimra --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml`**  

**Solution**:

Use the following command to reset Resource Aggregator for ODIM deployment:

   ```
python3 odim-controller.py --reset odimra --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --ignore-errors
   ```

------

**4. Resource Aggregator for ODIM redeployment fails after reset because of invalid odimCertsPath with the following error:**

   **`2021-08-17 09:54:48,613 - odim_controller - INFO  - Installing ODIMRA`**

   **`2021-08-17 09:54:48,613 - odim_controller - DEBUG - Reading config file /home/odim/kube_deploy_nodes.yaml`**

   **`2021-08-17 09:54:48,644 - odim_controller - DEBUG - Checking if the local user matches with the configired nodes user`**

   **`2021-08-17 09:54:48,764 - odim_controller - CRITICAL - ODIM-RA certificates path does not exist`**

**Solution**:

The error is displayed if `odimCertsPath` has invalid or null value. Perform the following steps:

   1. Navigate to the `kube_deploy_nodes.yaml` file.
   
   2. Specify a valid value for `odimCertsPath`, else specify its value as `""` (empty double quotations).


------

**5. Certificate generation fails with the following error:**

   **`2021-08-17 09:20:26,814 - odim_controller - INFO - Installing ODIMRA`**
   **`2021-08-17 09:20:26,819 - odim_controller - DEBUG - Reading config file /home/odim/kube_deploy_nodes.yaml`**
   **`2021-08-17 09:20:26,826 - odim_controller - DEBUG - Checking if the local user matches with the configired nodes user`**
   **`2021-08-17 09:20:29,072 - odim_controller - CRITICAL - ODIM-RA certificate generation failed`**
   **`Generating RSA private key, 4096 bit long modulus (2 primes)`**
   **`......................................................................................++++`**
   **`....................................................................................................................................................................................................++++`**
   **`e is 65537 (0x010001)`**
   **`Generating RSA private key, 4096 bit long modulus (2 primes)`**
   **`.................................++++`**
   **`...................................++++`**
   **`e is 65537 (0x010001)`**
   **`Error Loading request extension section v3_req`**
   **`140414011879872:error:220A4076:X509 V3 routines:a2i_GENERAL_NAME:bad ip address:../crypto/x509v3/v3_alt.c:457:value=null`**
   **`140414011879872:error:22098080:X509 V3 routines:X509V3_EXT_nconf:error in extension:../crypto/x509v3/v3_conf.c:47:name=subjectAltName, value=@alt_names`**
**`[Tue Aug 17 09:20:29 UTC 2021] -- ERROR -- /home/odim/ODIM/odim-controller/scripts/certs/OneNodeDeployment/odimra_server.csr generation failed`**

**Solution**:

   1. Navigate to the `kube_deploy_nodes.yaml` file.
   
   2. Specify valid values for the following parameters, else specify their values as `""` (empty double quotations):
   
      `odimraServerCertFQDNSan`<br />`odimraServerCertIPSan`<br />`odimraKafkaClientCertFQDNSan`<br />`odimraKafkaClientCertIPSan`<br />
------

**6. The plugin pod is in “CrashLoopBackOff” state after adding the plugin using the command  `python3 odim-controller.py --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --add plugin --plugin <plugin name>`** 

**Solution**:

You must be the owner of the plugin tar file.

   1. Run `kubectl describe pod <pod name> -n odim`.

   2. If "permission denied" error log is displayed at the end of the output, go to the location of stored plugin images and ensure you are the owner of the plugin tar file. If required, use the `sudo chown` command for changing permissions.

   3. Remove the plugin and add the plugin again using the command:

      `python3 odim-controller.py --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --add plugin --plugin <plugin name>`.

------

**7. When I restart Nginx, the service fails to start with the following error:** 

**`nginx: [emerg] bind() to <VIP>:<nginx_port> failed (99: Cannot assign requested address)`**

**Solution**:

You must restart Nginx systemd service ONLY on the leader node \(cluster node where Keepalived priority is set to the highest number). Restarting Nginx systemd service on the follower nodes (cluster nodes having lower Keepalived priority numbers) will result in the above error.`

------

**8. When the execution of odim-controller is active and I try to run another command, I get the following error:**

**`$ python3 odim-controller.py --config ~/kube_deploy_nodes.yaml --upgrade
aggregation 2021-09-16 05:53:35 - odim-controller - ERROR - An instance of
odim-controller is already active, another execution is not allowed`**

**Solution**:

The error is displayed because multiple executions of odim-controller is not allowed simultaneously.
To verify the active instance of odim-controller:

1. Run `$ ps -eaf | grep odim-controller`. The following output is displayed listing the active instance.

   ```
   odim 31291 4196 0 05:53 pts/1 00:00:00 python3 odim-controller.py --config
   ~/kube_deploy_nodes.yaml --upgrade account-session
   odim 31310 2734 0 05:53 pts/0 00:00:00 grep odim-controller
   ```

   Wait until the current execution completes and then retry the operation.

2. If odim-controller is not active and the error is displayed, run `$ ps -eaf | grep odim-controller`. The following output is displayed without listing any active instance.

   ```
     odim 31326 2734 0 05:53 pts/0 00:00:00 grep odim-controller
   ```

    In such cases, run `$ ls -ltr /tmp/odim-controller.lock`, remove the lock file by running`$ unlink /tmp/odim-controller.lock` and retry the operation.

------

**9. The following sample error is displayed:**

**`$ TASK [odimra-copy-image : Get list of plugins under plugin helm charts path]`**

**`An exception occurred during task execution. To see the full traceback, use -
vvv. The error was: TypeError: 'NoneType' object is not iterable
fatal: [lenovo2 -> localhost]: FAILED! => {"changed": false,
"module_stderr": "Traceback (most recent call last):\n File \"<stdin>\",
line 102, in <module>\n File \"<stdin>\", line 94, in _ansiballz_main\n
File \"<stdin>\", line 40, in invoke_module\n File \"/usr/lib/python3.8/
runpy.py\", line 207, in run_module\n return _run_module_code(code,
init_globals, run_name, mod_spec)\n File \"/usr/lib/python3.8/runpy.py\",
line 97, in _run_module_code\n _run_code(code, mod_globals,
init_globals,\n File \"/usr/lib/python3.8/runpy.py\", line 87, in _run_code
\n exec(code, run_globals)\n File \"/tmp/ansible_find_payload_o1n23qsr/
ansible_find_payload.zip/ansible/modules/files/find.py\", line 475, in
<module>\n File \"/tmp/ansible_find_payload_o1n23qsr/
ansible_find_payload.zip/ansible/modules/files/find.py\", line 409, in main
\nTypeError: 'NoneType' object is not iterable\n", "module_stdout": "",
"msg": "MODULE FAILURE\nSee stdout/stderr for the exact error", "rc": 1}`**

**Solution**:

The error is displayed if the `odimPluginPath` parameter has invalid or null value. Perform the following steps:

1. Navigate to the kube_deploy_nodes.yaml file.
2. Specify a valid value for `odimPluginPath`, else specify its value as "" (empty double quotation marks).

------

**10. The following sample error is displayed:
` 2021-09-22 13:36:05 - odim-controller - ERROR - Caught an exception: can only concatenate str (not "NoneType") to str`**

**Solution**:

The error is displayed if the `httpProxy`, `httpsProxy` and `noProxy` parameters have invalid or null values. Perform the following steps:

1. Navigate to the kube_deploy_nodes.yaml file.
2. Specify valid values for `httpProxy`, `httpsProxy` and `noProxy`, else specify their respective values as `""` (empty double quotation marks).

------

**11. I see the following error when I upgrade odimra-secret:**

**`bruce@deploy:~/R4H60-11008/odim-controller/scripts$ python3 odimcontroller.
py --config /home/${USER}/R4H60-11008/odim-controller/scripts/
kube_deploy_nodes.yaml --upgrade odimra-secret
2021-09-23 07:37:39 - odim-controller - INFO - Upgrading config mapodimrasecret
2021-09-23 07:37:39 - odim-controller - DEBUG - Reading config file /home/
bruce/R4H60-11008/odim-controller/scripts/kube_deploy_nodes.yaml
2021-09-23 07:37:39 - odim-controller - DEBUG - Checking if the local user
matches with the configired nodes user
2021-09-23 07:37:39 - odim-controller - DEBUG - Reading group_var file /home/
bruce/R4H60-11008/odim-controller/odimra/group_vars/all/all.yaml
2021-09-23 07:37:39 - odim-controller - INFO - Full helm chart name odimrasecret-
2.0.0
[Thu 23 Sep 2021 07:37:39 AM MDT] -- INFO -- rootCA crt and key exists, not
generating again
Generating RSA private key, 4096 bit long modulus (2 primes)
.............................................................................
........................++++
.............................................................................
......................++++
e is 65537 (0x010001)
Error Loading request extension section v3_req
140428313810240:error:220A4076:X509 V3 routines:a2i_GENERAL_NAME:bad ip
address:../crypto/x509v3/v3_alt.c:477:value=
140428313810240:error:22098080:X509 V3 routines:X509V3_EXT_nconf:error in
extension:../crypto/x509v3/v3_conf.c:47:name=subjectAltName, value=@alt_names
[Thu 23 Sep 2021 07:37:40 AM MDT] -- ERROR -- /home/bruce/R4H60-11008/odim-controller/scripts/certs/snap4/odimra_server.csr generation failed
2021-09-23 07:37:40 - odim-controller - CRITICAL - ODIM-RA certificate
generation failed`**

**Solution**:

The error is displayed if the `odimraServerCertIPSan` and `odimraKafkaClientCertIPSan` parameters have invalid or null values. Perform the following steps:

1. Navigate to the kube_deploy_nodes.yaml file.
2. Specify valid values for `odimraServerCertIPSan` and `odimraKafkaClientCertIPSan`, else specify
   both their values as `""` (empty double quotation marks).

------

**12. When I deploy Resource Aggregator for ODIM or when I scale up a node, nginx installation fails with the following error:**

**`TASK [pre-install : Install packages required by nginx] ************************`**

**`changed: [containerd5]`**

**`Tuesday 20 December 2022 17:29:31 +0000 (0:00:11.971)    0:05:19.037 ******`**

**Workaround**:

1. Identify the process by running the following command on the affected server and note down the process id:

   ```
   ps -ef|grep "apt install curl" 
   ```

   The following sample output is displayed:

   ```
   root       <process-id>   21151  0 Dec20 pts/1    00:00:00 apt install curl gnupg2 ca-certificates lsb-release ubuntu-keyring
   ```

2. Kill the process by running the following command:

   ```
   sudo kill -9 <process-id>
   ```



## Other Frequently Asked Questions

**1. How do I know if the User ID in the configuration file already exists?**

1. Run `getent passwd <User ID>` on all cluster nodes. If the Used ID already exists, you get the following output:

   `odimra: x:<User ID>:<User ID>::/home/odimra:/bin/bash`

2. Set a unique User ID of your own in the configuration file.

------

**2. How do I know if the Group ID in the configuration file already exists?**

1. Run `getent group <Group ID>` on all cluster nodes. If the Group ID already exists, you get the following output:

   `odimra:x:2021:`

2. Set a unique Group ID of your own in the configuration file.

------

**3. Which parameters in the `kube_deploy_nodes.yaml ` file are immutable after deployment?**

<blockquote> Caution: Do NOT modify the private key "odimra_rsa.private" and the public key "odimra_rsa.public". If modified, it will result in service non-recoverable data loss, unless backup of the keys are present. </blockquote>

- appsLogPath
- etcdConfPath
- etcdDataPath
- groupID
- haDeploymentEnabled
- hostIP
- hostname
- kafkaConfPath
- kafkaDataPath
- namespace
- pluginHelmChartsPath
- redisInmemoryDataPath
- redisOndiskDataPath
- rootServiceUUID
- userID
- zookeeperConfPath
- zookeeperDataPath

------

**4. What is the recommended value for replica count for each Resource Aggregator for ODIM service?**

For a single node deployment, recommended value is 1. 
For multi-node cluster deployment, count can be set to 1 or desired value (verified value is 3).

------

**5. How do I change the node's sudo password, which is persisted often?**

1. Delete the following file:
   `<odim-controller-cloned-path>/kubespray/inventory/k8s-cluster-<deploymentID>/.sudo_pw`. 
   
2. After the deletion, when you invoke odim-controller for any operation on the cluster, you are prompted to type a new password. 
   

Alternatively, you can perform the following steps:

1. Enter the password of the default non-root user \(that was set across all cluster nodes initially\) in plain text in a file called `nodePasswordFile`. Save the file. 

   ```
   vi nodePasswordFile
   ```

2. To encrypt the entered password, run the following command: 

   ```
   ./odim-vault -key ~/ODIM/odim-controller/\
   scripts/odimVaultKeyFile -encrypt /home/${USER}/ODIM/odim-controller/\
   scripts/nodePasswordFile
   ```

   **Result**: `nodePasswordFile` contains the encrypted node password.

3. Change the file permissions of` nodePasswordFile`.

   ```
   chmod 0400 /home/${USER}/ODIM/odim-controller/\
   scripts/nodePasswordFile
   ```

------

**6. How do I enable `kubectl` usage without the `sudo` command?**

Run the following commands on controller nodes:

```
-p $HOME/.kube
```

```
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
```

```
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

------

**7. How do I check the logs of Resource Aggregator for ODIM services?**

Navigate to the default log path `/var/log/odimra` or to the path specified in the `appsLogPath` parameter in `kube_deploy_nodes.yaml` file on that node for the service log file.

------

**8. How do I check the logs of plugin services?**

Navigate to the default log path `/var/log/<plugin>` or to the path specified in the `logPath` parameter in  `<plugin>-config.yaml` file.

------

**9. How do I check the logs of third party services (Kafka, Zookeeper, Redis, etcd)?**

1. On the master node, run the following command to get the name of the pod:

   ```
   kubectl get pods -n odim -o wide
   ```

2. Run the following command to tail the logs:
   
   ```
   kubectl logs -n odim -f <pod_name>
   ```
   

------

**10. Resource Aggregator for ODIM deployment fails. What are some of the possible ways of addressing this issue?**

   ​Resource Aggregator for ODIM deployment can fail due to multiple reasons. Perform the following tasks:

1. Analyze the probable errors and fix them.

2. Reset Resource Aggregator for ODIM deployment using the following command:
   `python3 odim-controller.py --reset odimra --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --ignore-errors`

3. Retry deploying Resource Aggregator for ODIM services using the following command:

   `python3 odim-controller.py --deploy \<br/> odimra --config /home/${USER}/ODIM/odim-controller/\<br/>scripts/kube_deploy_nodes.yaml`

<blockquote>NOTE: Verify the content and the formatting of the content in the `kube_deploy_nodes.yaml` configuration file. Formatting in the `kube_deploy_nodes.yaml.tmpl` file provided must be copied and retained.</blockquote>

