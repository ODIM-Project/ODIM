[![build_deploy_test Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_deploy_test/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)
[![build_unittest Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_unittest/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)

# Table of contents

- [Deploying ODIMRA](#deploying-odimra)
  * [1. Setting up OS and Docker environment](#1-setting-up-os-and-docker-environment)
  * [2. Installing the resource aggregator for ODIM and GRF plugin](#2-installing-the-resource-aggregator-for-odim-and-grf-plugin)
    + [Default user credentials for ODIMRA and GRF Plugin](#default-user-credentials-for-odimra-and-grf-plugin)
- [Modifying default configuration parameters for the resource aggregator](#modifying-default-configuration-parameters-for-the-resource-aggregator)
- [Configuring proxy for Docker](#configuring-proxy-for-docker)
   



# Deploying ODIMRA
## 1. Setting up OS and Docker environment

**Prerequisites**
------------------
Ensure that the Internet is available. If your system is behind a corporate proxy or firewall, set your proxy configuration. To know how to set proxy, see information provided at `https://www.serverlab.ca/tutorials/linux/administration-linux/how-to-set-the-proxy-for-apt-for-ubuntu-18-04/`.  


**Procedure**
--------------
1. Download and install `Ubuntu 18.04 LTS` on your system.
    >   **NOTE:**  Before installation, configure your system IP to access the data center network.
2. Install `Ubuntu Make` on your system.
To install `Ubuntu Make`, run the following command:
   ```
   $ sudo apt install make
   ```
3. Install `Java 11` on your system.
   To install `Java 11`, run the following command:
   ```
    $ sudo apt install openjdk-11-jre-headless -y
   ```
4. Set up Docker environment:
     > **IMPORTANT:** This procedure installs only the community edition of Docker.

   a. To install Docker, run the following commands:
   1.  ```
       $ sudo apt update
       ```
   2.  ```
       $ sudo apt install apt-transport-https ca-certificates curl software-properties-common -y
       ```
   3.  ```
       $ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
        ```
   4. ```
      $ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
      ```
   5. ```
      $ sudo apt update
      ```
   6. ```
      $ apt-cache policy docker-ce
       ```
      
        The following output is generated:
      ```
      docker-ce:
      Installed: (none)
      Candidate: 18.03.1~ce~3-0~ubuntu
      Version table:
      18.03.1~ce~3-0~ubuntu 500
      500 https://download.docker.com/linux/ubuntu bionic/stable
      amd64 Packages
      ```
     
      > **NOTE:** docker-ce is not installed, but the candidate for  installation is from the Docker repository for Ubuntu 18.04 (bionic).
     
    7. ```
       $ sudo apt install docker-ce -y
       ```
    8. ```
       $ sudo apt-get install docker-compose -y
       ```
        
       >  **NOTE:** To run the commands without sudo, add your username to the docker group using the following command:
        ```
         $ sudo usermod -aG docker ${USER}
        ```

    b. Check the status of Docker:
      ```
    $ sudo systemctl status docker
      ```
 
   If Docker is active and running, the following output is generated:
   ```
   docker.service - Docker Application Container Engine
   Loaded: loaded (/lib/systemd/system/docker.service; enabled; vendor
   preset: enabled)
   Active: active (running) since Thu 2018-07-05 15:08:39 UTC; 2min 55s
   ago
   Docs: https://docs.docker.com
   Main PID: 10096 (dockerd)
   Tasks: 16
   CGroup: /system.slice/docker.service
   +-10096 /usr/bin/dockerd -H fd://
   +-10113 docker-containerd --config /var/run/docker/containerd/
   containerd.toml
   ```
     
   >  **NOTE:** If your system is behind a corporate proxy, ensure to configure Docker to use proxy server and restart docker services. To know how to configure Docker proxy, see [Configuring Docker proxy](#configuring-proxy-for-docker).
						   
     
   c. Restart the server.
      ```
     $ sudo init 6
     ```
      
   >  **NOTE:** To enable Docker service to start on reboot, run the following command:
   
       `$ sudo systemctl enable docker`
  
   
	   


	   
## 2. Installing the resource aggregator for ODIM and GRF plugin
This section provides a step-by-step procedure for deploying the resource aggregator for ODIM (odimra) and GRF plugin.

  
  **NOTE:**
  - All configuration parameters are set to default values in the configuration files for odimra and GRF plugin. 
  - The following ports are used for deploying odimra and GRF plugin:
    45000, 45001, 45101-45110, 9092, 9082, 6380, 6379, 8500, 8300, 8302, 8301, 8600
    Ensure that the above ports are not in use.


**WARNING:** Do not run the commands provided in this section as root user unless mentioned.

**Procedure**
--------------
1. Clone the odimra repository form `https://github.com/ODIM-Project/ODIM.git` to the home directory of the user.
   ```
   $ git clone https://github.com/ODIM-Project/ODIM.git
   ```
2. Choose a Fully Qualified Domain Name (FQDN) for the resource aggregator server. 
   Example: odim.local.com.
3. Set FQDN to environment of the host machine using the following command:
    ```
    $ export FQDN=<user_preferred_fqdn_for_host>
    ```
4. Set the environment variable, `HOSTIP` to the IP address of your system.
   ```
   $ export HOSTIP=<ip_address_of_your_system>
   ```

5. Set below environment variables with user and group ID to be used for odimra
   ```
   $ export ODIMRA_USER_ID=1234
   $ export ODIMRA_GROUP_ID=1234
   ```

6. Set up FQDN in the `/etc/hosts` file (only if there is no DNS infrastructure):

    a. Open the `/etc/hosts` file for editing:
      ```
      $ sudo vim /etc/hosts
      ```
    b. Scroll to the end of the file, add the following line, and then save:
      ```
      <host_ipv4_address> <user_preferred_fqdn_for_host>
      ```
   Example:
`<host_ipv4_address> <fqdn>`

7. Generate certificates:

   
   **NOTE:**
   - Self-signed Root CA (Certificate Authority) certificate and key are generated with 4096 key length and sha512 digest algorithm.
   - Using the generated CA certificate, certificates and private keys for the resource aggregator services are also generated with 4096 key length and sha512 digest algorithm. They are valid for services matching the provided FQDN. You can use one-word description of the certificate as the common name.
   - Certificates are used by the resource aggregator services to communicate    internally (Remote Procedure Call) and with the plugin services.
   - If you are using an intermediate CA for signing certificates assigned to the resource aggregator and the plugin services, ensure to:
        - Append all the intermediate certificates to the server certificate file in   the order such that each certificate has signed the preceding one.
        - Append the Root CA used for signing the intermediate CA to the resource   aggregator CA file.


    **Procedure**
   
   a. Navigate to the path: `ODIM/build/cert_generator`
      ```
       $ cd ODIM/build/cert_generator
      ```

    > NOTE: `ODIM/build/cert_generator` contains the automated scripts to        generate the resource aggregator and GRF plugin TLS and Kafka TLS  certificates.

   b. Use the following command to generate the resource aggregator and the      GRF plugin certificates. Provide FQDN as a command-line argument.
      ```
      $ ./generate_odimra_cert.sh <FQDN>
      ```
   c. Use the following command to generate Kafka TLS certificate:
      ```
       $ ./generate_kafka_certs.sh kafka
      ```
   d. Use the following command to generate Zookeeper TLS certificate:
      ```
       $ ./generate_zookeeper_certs.sh zookeeper
      ```
   e. Use the following command to copy the resource aggregator, the GRF  plugin, the Kafka and Zookeeper TLS certificates:
     ```
      $ sudo ./copy_certificate.sh
     ```
     The following files are copied in the path: `/etc/odimracert/`
      - rootCA.crt
      - odimra_server.key
      - odimra_server.crt
      - odimra_rsa.public
      - odimra_rsa.private
      - odimra_kafka_client.key
      - odimra_kafka_client.crt

     The following files are copied in the path: `/etc/kafka/conf/`
      - kafka.keystore.jks
      - kafka.truststore.jks

     The following files are copied in the path: - /etc/zookeeper/conf
      - zookeeper.keystore.jks
      - zookeeper.trustore.jks
      
    The following files are copied in the path: `/etc/plugincert/`
      - rootCA.crt
      - odimra_server.key
      - odimra_server.crt
      - odimra_kafka_client.key
      - odimra_kafka_client.crt
      

8. Navigate to the odimra folder.
   ```
   $ cd ~/ODIM
   ```

9. Use the following command to deploy and start the containers:
   ```
   $ make all
   ```
    The following containers are loaded.
      - build_odimra_1
      - build_kafka_1
      - build_zookeeper_1
      - build_redis_1
      - build_consul_1
      - build_grf_plugin_1

10. Verify that the resource aggregator services are running successfully.
   ```
   $ ps -eaf | grep svc
   ```
   All the resource aggregator services are listed:
   ```
   root 8343 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-api --registry=consul --registry_address=consul:8500 --
   client_request_timeout=1m
   bruce 8346 8343 0 15:20 ? 00:00:00 ./svc-api --
   registry=consul --registry_address=consul:8500 --client_request_timeout=1m
   root 8406 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-account-session --registry=consul --registry_address=consul:8500 --
   server_address=odim:45101
   bruce 8408 8406 0 15:20 ? 00:00:00 ./svc-account-session --
   registry=consul --registry_address=consul:8500 --server_address=odim:45101
   root 8424 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-aggregation --registry=consul --registry_address=consul:8500 --
   server_address=odim:45102
   bruce 8426 8424 0 15:20 ? 00:00:00 ./svc-aggregation --
   registry=consul --registry_address=consul:8500 --server_address=odim:45102
   root 8441 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-events --registry=consul --registry_address=consul:8500 --
   server_address=odim:45103
   bruce 8443 8441 0 15:20 ? 00:00:00 ./svc-events --
   registry=consul --registry_address=consul:8500 --server_address=odim:45103
   root 8458 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-systems --registry=consul --registry_address=consul:8500 --
   server_address=odim:45104
   bruce 8460 8458 0 15:20 ? 00:00:00 ./svc-systems --
   registry=consul --registry_address=consul:8500 --server_address=odim:45104
   root 8474 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-task --registry=consul --registry_address=consul:8500 --
   server_address=odim:45105
   bruce 8476 8474 0 15:20 ? 00:00:00 ./svc-task --
   registry=consul --registry_address=consul:8500 --server_address=odim:45105
   root 8492 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-fabrics --registry=consul --registry_address=consul:8500 --
   server_address=odim:45106
   bruce 8494 8492 0 15:20 ? 00:00:00 ./svc-fabrics --
   registry=consul --registry_address=consul:8500 --server_address=odim:45106
   root 8519 8233 0 15:20 ? 00:00:00 sudo -E -u odim nohup ./
   svc-managers --registry=consul --registry_address=consul:8500 --
   server_address=odim:45107
   bruce 8521 8519 0 15:20 ? 00:00:00 ./svc-managers --
   registry=consul --registry_address=consul:8500 --server_address=odim:45107
   ```


    **NOTE:**
    - The resource aggregator configuration files are available at `/etc/odimra_config`.
    -  The GRF configuration files are available at `/etc/grf_plugin_config`.
    - The resource aggregator API service runs on the default port 45000.
    - The GRF plugin API service runs on default port 45001.
    - The resource aggregator logs are available at `/var/log/odimra`.
    - The GRF plugin logs are available at `/var/log/GRF_PLUGIN`.


11. To configure log rotation, do the following:

    a. Navigate to the `/etc/logrotate.d` directory.
    ```
    $ cd /etc/logrotate.d
    ```
    b. Edit a file by name odimra (`$ sudo vi odimra`) and add the following content:
     ```
    /var/log/GRF_PLUGIN/*.log
     /var/log/odimra/*.log {
    hourly
    missingok
    rotate 10
    notifempty
    maxsize 1M
    compress
    create 0644 <user> <group>
    shred
    copytruncate
    }
    ``` 
    c. Navigate to the `/etc/cron.hourly` directory.
       ```
       $ cd /etc/cron.hourly
       ```
    d. Edit a file by name logrotate (`$ sudo vi logrotate`) and add following content:
	   ```
       logrotate -s /var/lib/logrotate/status /etc/logrotate.d/odimra
	   ```
    e. To verify that the configuration is working, run the following command:
      ```
      $ sudo logrotate -v -f /etc/logrotate.d/odimra
  	  ```
  
12.  Use following commands to undeploy odimra solution and remove the docker images, persistent data, logs.

     ```
     $ make clean
     ```
     The above command will perform below tasks:
     - brings down all the deployed containers.
     - remove only those docker images which were created and deployed as containers.
     - removes data stored by Consul, Redis & Kafka.

     ```
     $ make deepclean
     ```
     The above command will perform below tasks: 
     - brings down the deployed conatiners
     - remove all the docker images, as well as intermediate/dependent images created as part of deployment.
     - removes configuration & data stored by Consul, Redis & Kafka
     - removes all generated certifictes.
     - removes logs files created for odimra services and grfplugin.

     [NOTE] Provide the sudo password when prompted by above commands.

     [CAUTION] The above commands are not encouraged to be executed in production envoirnment as this will erase important data.
               The action is irrecoverable and will wipe all the odimra completely. 

13. To add the Generic Redfish Plugin and servers to the resource aggregator for ODIM, refer to the following readme.  
    https://github.com/ODIM-Project/ODIM/blob/development/svc-aggregation/README.md
	
	
### Default user credentials for ODIMRA and GRF Plugin 


ODIMRA:

```
username: admin
password: Od!m12$4
```

GRF PLUGIN:

```
username: admin
password: GRFPlug!n12$4
``` 
 
 
  
#  Modifying default configuration parameters for the resource aggregator

1.   Navigate to the `build_odimra_1` container using the following command: 

      ```
     $ docker exec -it build_odimra_1/bin/bash
     ```

2.   Edit the parameters in the `odimra_config.json` file located in this path:   `/etc/odimra_config/odimra_config.json` and save. 

     The parameters that are configurable are listed in the following table.
      > **NOTE:** It is recommended not to modify parameters other than the ones listed in the following table.
     
     |Parameter|Type|Description|
     |---------|----|-----------|
     |RootServiceUUID|String|Static `UUID` used for the resource aggregator root service.  NOTE: Take a backup copy of `RootServiceUUID` as it is required during reinstallation.|
     |LocalhostFQDN|String|FQDN of the host.|
     |KeyCertConf{|Array| |
     |RootCACertificatePath|String|TLS Root CA file path (which can be a chain of CAs for verifying entities interacting with the resource aggregator services).|
     |RPCPrivateKeyPath|String|TLS private key file path for the microservice RPC communications.|
     |RPCCertificatePath}|String|TLS certificate file path for the microservice RPC communications.|
     |APIGatewayConf{|Array| |
     |Host|String|Host address for the resource aggregator API gateway.|
     |Port|String|Port for the resource aggregator API gateway.|
     |PrivateKeyPath|String|TLS private key file path for the API gateway.|
     |CertificatePath}|String|TLS certificate file path for the API gateway.|
     |TLSConf{|Array|TLS configuration parameters.<br> Note: It is not recommended to change these settings. |
     |MinVersion|String|Default value: `TLS1.2`<br> Supported values: `TLS1.0, TLS1.1, TLS1.2`<br> Recommended value: `TLS1.2`|
     |MaxVersion|String|Default value: `TLS1.2`<br> Supported values: `TLS1.0, TLS1.1, TLS1.2`<br>  Recommended value: `TLS1.2`<br>  NOTE: If `MinVersion` and `MaxVersion` are not specified, they will be set to default values.<br> If `MinVersion` and `MaxVersion` are set to unsupported values, the resource aggregator and the plugin services will exit with errors.|
     |VerifyPeer|Boolean| Default value: true<br>  Recommended value: true<br>  NOTE:<br>  - `VerifyPeer` is set to true, by default. For secure plugin interaction, add root CA certificate (that is used to sign the certificates of the southbound entities) to root CA certificate file.  If `VerifyPeer` is set to false, SSL communication will be insecure.  After setting `VerifyPeer` to false, restart the resource aggregator container (`odim_1`).<br>  - If `TLS1.2` is used, ensure that the entity certificate has `SAN` field for successful validation.  - Northbound entities interacting with resource aggregator `API` service can use root CA that signed odimra's certificate.|
     |PreferredCipherSuites}|List| Default and supported values: See "List of supported (default) cipher suites".<br> IMPORTANT:<br>  - If `PreferredCipherSuites` is not specified, it will be set to default cipher (secure) suites.<br>  - If `PreferredCipherSuites` is set to unsupported cipher suites, the resource aggregator and the plugin services will exit with errors.|

     List of supported (default) cipher suites:
     
     |||
     |----|----|
     |TLS_RSA_WITH_AES_128_GCM_SHA256|Supported in TLS1.2|
     |TLS_RSA_WITH_AES_256_GCM_SHA384|Supported in TLS1.2|
     |TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256|Supported in TLS1.2|
     |TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384|Supported in TLS1.2|
     |TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256|Supported in TLS1.2|
     |TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384|Supported in TLS1.2|

3.   Exit from Docker using the following command: 

     ```
     $ exit
     ```

4.   Restart Docker using the following command: 

     ```
     $ docker restart build_odimra_1
      ```
    
   
   
   
   
# Configuring proxy for Docker

<blockquote>
IMPORTANT:

During the course of this procedure, you will be required to create files and copy content into them. In the content to be copied, substitute the parameters listed in the following table with their original values:

|Parameter|Description|
|---------|-----------|
|`<Proxy_URL>` |Your company URL.|
|`<ODIM_server_VM_IP>` |The IP address of the system where the resource aggregator is installed.|
|`<FQDN>` |FQDN of the resource aggregator server.|

</blockquote>

**Procedure**
--------------

1.   In the home directory of odimra user, create a hidden directory called .docker, and then create a file called config.json inside it.

      ```
       mkdir .docker
      ```

      ```
       cd .docker
      ```

      ```
       vi config.json
      ```	   

2.   Add the following content in the config.json file and save: 

      ```
      {
         "proxies":
        {
           "default":
          {
             "httpProxy": "<Proxy_URL>",
             "httpsProxy": "<Proxy_URL>",
             "noProxy": "localhost,127.0.0.1, <ODIM_server_VM_IP>"
          }
        }
      }
    
     ```

3.   Update the `/etc/environment` file with the following content using sudo: 

      ```
      export http_proxy=<Proxy_URL>
      export https_proxy=<Proxy_URL>
      HOSTIP=<ODIM_server_VM_IP>
      FQDN=<FQDN>
      ```

4.   Do the following on the resource aggregator server: 

     1. Create a directory using the following command: 

         ```
         $ sudo mkdir /etc/systemd/system/docker.service.d/
         ```

     2. Create a file called http-proxy.conf in the `/etc/systemd/system/docker.service.d/` directory. 
     3. Add the following content in the `http-proxy.conf` file and save: 

         ```
         [Service]
         Environment="HTTP_PROXY=<Proxy_URL>"
         Environment="HTTPS_PROXY=<Proxy_URL>"
         Environment="NO_PROXY=localhost,127.0.0.1, <ODIM_server_VM_IP>"
        
         ```

     4. Run the following commands: 

         ```
         $ sudo systemctl daemon-reload
         ```

         ```
         $ sudo service docker restart
         ```

