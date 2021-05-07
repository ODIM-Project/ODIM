[![build_deploy_test Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_deploy_test/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)
[![build_unittest Actions Status](https://github.com/ODIM-Project/ODIM/workflows/build_unittest/badge.svg)](https://github.com/ODIM-Project/ODIM/actions)

# Table of contents

- [Introduction](#introduction)
  * [Resource Aggregator for ODIM logical architecture](#resource-aggregator-for-odim-logical-architecture)
- [Deploying ODIMRA](#deploying-odimra)
  * [1. Setting up OS and Docker environment](#1-setting-up-os-and-docker-environment)
  * [2. Installing the resource aggregator for ODIM, the Generic redfish plugin (GRF), and the Unmanaged Rack Plugin (URP)](#2-installing-the-resource-aggregator-for-odim-the-generic-redfish-plugin-grf-and-the-unmanaged-rack-plugin-urp)
    + [Default user credentials for ODIMRA, the GRF Plugin, and URP](#default-user-credentials-for-odimra-the-grf-plugin-and-urp)
- [Modifying default configuration parameters for the resource aggregator](#modifying-default-configuration-parameters-for-the-resource-aggregator)
- [Configuring proxy for Docker](#configuring-proxy-for-docker)
- [Uninstalling ODIMRA](#uninstalling-odimra)
- [CI Process](#ci-process)
  * [GitHub action workflow details](#github-action-workflow-details)



# Introduction 

 Welcome to Resource Aggregator for Open Distributed Infrastructure Management!

Resource Aggregator for Open Distributed Infrastructure Management (ODIMRA) is a modular, open framework for
simplified management and orchestration of distributed physical infrastructure. It provides a unified management platform for
converging multivendor hardware equipment. By exposing a standards-based programming interface, it enables easy and
secure management of wide range of multivendor IT infrastructure distributed across multiple data centers.

ODIMRA framework comprises the following two components.

- The resource aggregation function (the resource aggregator)

  The resource aggregation function is the single point of contact between the northbound clients and the
  southbound infrastructure. Its primary function is to build and maintain a central resource inventory. It exposes
  Redfish-compliant APIs to allow northbound infrastructure management systems to:
    - Get a unified view of the southbound compute, local storage, and Ethernet switch fabrics available in the
      resource inventory.
    - Gather crucial configuration information about southbound resources.
    - Manipulate groups of resources in a single action.
    - Listen to similar events from multiple southbound resources.
	
- One or more plugins

  The plugins abstract, translate, and expose southbound resource information to the resource aggregator through
  RESTful APIs. HPE Resource Aggregator for ODIM supports:
 
    - Generic Redfish plugin for ODIM (GRF): Generic Redfish plugin that can be used as a plugin for any Redfishcompliant
      device.
	- Plugin for unmanaged racks (URP): Plugin that acts as a resource manager for unmanaged racks. 
    - Integration of additional third-party plugins.  




##  Resource Aggregator for ODIM logical architecture

Resource Aggregator for ODIM framework adopts a layered architecture and has many functional layers.


The following figure shows these functional layers of Resource Aggregator for ODIM deployed in a data center.


![ODIM_architecture](docs/images/arch.png)


- **API layer**


This layer hosts a REST server which is open-source and secure. It learns about the southbound resources from
the plugin layer and exposes the corresponding Redfish data model payloads to the northbound clients. The
northbound clients communicate with this layer through a REST-based protocol that is fully compliant with
DMTF's Redfish® specifications (Schema 2020.3 and Specification 1.11.1).
The API layer sends user requests to the plugins through the aggregation, the event, and the fabric services.

- **Services layer**


The services layer is where all the services are hosted. This layer implements service logic for all use cases
through an extensible domain model (Redfish Data Model). Requests coming from the API layer and the
responses coming from the plugin layer are mapped to the actual end resources in this layer. It maintains the state
for event subscriptions, credentials, and tasks. It also hosts a message bus called the Plug-in Message Bus (PMB).


![Redfish_data_model](docs/images/redfish_data_model.png)


- **Event message bus layer**


This layer hosts a message broker which acts as a communication channel between the upper layers and the
plugin layer. It supports common messaging architecture and real-time streaming. Resource Aggregator for
ODIM uses Kafka as the event message bus.
The services and the event message bus layers host Redis data store.

- **Plugin layer**


This layer connects the actual managed resources to the aggregator layers and is de-coupled from the upper
layers. A plugin abstracts vendor-specific access protocols to a common interface which the aggregator layers use
to communicate with the resources. It uses REST-based communication which is based on OpenAPI Specification
v3.0 to interact with the other layers. It collects events to be exposed to fault management systems and uses the
event message bus to publish events. The messaging mechanism is based on OpenMessaging Specification.
The plugin layer allows developers to create plugins on any tool set of their choice without enforcing any strict
language binding. To know how to develop plugins, refer to [Resource Aggregator for Open Distributed
Infrastructure Management Plugin Developer's Guide](https://github.com/ODIM-Project/ODIM/blob/development/plugin-redfish/README.md).

   

# Introduction 

 Welcome to Resource Aggregator for Open Distributed Infrastructure Management!


Resource Aggregator for Open Distributed Infrastructure Management (ODIMRA) is a modular, open framework for
simplified management and orchestration of distributed physical infrastructure. It provides a unified management platform for
converging multivendor hardware equipment. By exposing a standards-based programming interface, it enables easy and
secure management of wide range of multivendor IT infrastructure distributed across multiple data centers.

ODIMRA framework comprises the following two components.

- The resource aggregation function (the resource aggregator)

  The resource aggregation function is the single point of contact between the northbound clients and the
  southbound infrastructure. Its primary function is to build and maintain a central resource inventory. It exposes
  Redfish-compliant APIs to allow northbound infrastructure management systems to:
    - Get a unified view of the southbound compute, local storage, and Ethernet switch fabrics available in the
      resource inventory.
    - Gather crucial configuration information about southbound resources.
    - Manipulate groups of resources in a single action.
    - Listen to similar events from multiple southbound resources.
	
- One or more plugins

  The plugins abstract, translate, and expose southbound resource information to the resource aggregator through
  RESTful APIs. HPE Resource Aggregator for ODIM supports:
 
    - Generic Redfish plugin for ODIM (GRF): Generic Redfish plugin that can be used as a plugin for any Redfishcompliant
      device.
	- Plugin for unmanaged racks (URP): Plugin that acts as a resource manager for unmanaged racks. 
    - Integration of additional third-party plugins.  




##  Resource Aggregator for ODIM logical architecture

Resource Aggregator for ODIM framework adopts a layered architecture and has many functional layers.


The following figure shows these functional layers of Resource Aggregator for ODIM deployed in a data center.

![ODIM_architecture](images/arch.png)

- **API layer**


This layer hosts a REST server which is open-source and secure. It learns about the southbound resources from
the plugin layer and exposes the corresponding Redfish data model payloads to the northbound clients. The
northbound clients communicate with this layer through a REST-based protocol that is fully compliant with
DMTF's Redfish® specifications (Schema 2020.3 and Specification 1.11.1).
The API layer sends user requests to the plugins through the aggregation, the event, and the fabric services.

- **Services layer**


The services layer is where all the services are hosted. This layer implements service logic for all use cases
through an extensible domain model (Redfish Data Model). Requests coming from the API layer and the
responses coming from the plugin layer are mapped to the actual end resources in this layer. It maintains the state
for event subscriptions, credentials, and tasks. It also hosts a message bus called the Plug-in Message Bus (PMB).

![Redfish_data_model](images/redfish_data_model.png)

- **Event message bus layer**


This layer hosts a message broker which acts as a communication channel between the upper layers and the
plugin layer. It supports common messaging architecture and real-time streaming. Resource Aggregator for
ODIM uses Kafka as the event message bus.
The services and the event message bus layers host Redis data store.

- **Plugin layer**


This layer connects the actual managed resources to the aggregator layers and is de-coupled from the upper
layers. A plugin abstracts vendor-specific access protocols to a common interface which the aggregator layers use
to communicate with the resources. It uses REST-based communication which is based on OpenAPI Specification
v3.0 to interact with the other layers. It collects events to be exposed to fault management systems and uses the
event message bus to publish events. The messaging mechanism is based on OpenMessaging Specification.
The plugin layer allows developers to create plugins on any tool set of their choice without enforcing any strict
language binding. To know how to develop plugins, refer to [Resource Aggregator for Open Distributed
Infrastructure Management Plugin Developer's Guide](https://github.com/ODIM-Project/ODIM/blob/development/plugin-redfish/README.md).


# Deploying ODIMRA

## 1. Setting up OS and Docker environment

**Prerequisites**
------------------
- Ensure that the Internet is available. If your system is behind a corporate proxy or firewall, set your proxy configuration. To know how to set proxy, see information provided at https://www.serverlab.ca/tutorials/linux/administration-linux/how-to-set-the-proxy-for-apt-for-ubuntu-18-04/.  

- Ensure not to create `odimra` user on the system where you want to deploy ODIMRA.

**Procedure**
--------------
Perform the following steps on the system where you want to deploy ODIMRA:
1. Download and install `Ubuntu 18.04 LTS`.
    >   **NOTE:**  Before installation, configure your system IP to access the data center network.
2. Install `Ubuntu Make`. It is required to build ODIMRA services.

   To install `Ubuntu Make`, run the following command:
   ```
   $ sudo apt install make
   ```
3. Install `Java 11`. It is required to create keystores for Kafka and Zookeeper services.

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
     
      > **NOTE:** docker-ce is not installed, but the candidate for installation is from the Docker repository for Ubuntu 18.04 (bionic).
     
    7. ```
       $ sudo apt install docker-ce -y
       ```
    8. ```
       $ wget https://github.com/docker/compose/releases/download/1.25.5/docker-compose-Linux-x86_64
       ```
	   
	   ```
       $ sudo cp docker-compose-Linux-x86_64 /usr/bin/docker-compose
       ```
        
       >  **NOTE:** To run the commands without sudo, add your username to the Docker group using the following command:
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
     
   >  **NOTE:** If your system is behind a corporate proxy, ensure to configure Docker to use proxy server and restart the Docker services. To know how to configure Docker proxy, see [Configuring Docker proxy](#configuring-proxy-for-docker).
						   
     
   c. To enable Docker service to start on reboot, run the following command:
      ```
      $ sudo systemctl enable docker
	  ```
      
	
   d. Restart the system.
      ```
      $ sudo reboot
      ```
   

  
   
	   


	   
## 2. Installing the resource aggregator for ODIM, the Generic redfish plugin (GRF), and the Unmanaged Rack Plugin (URP)
This section provides a step-by-step procedure for deploying ODIMRA, the GRF plugin, and URP.

  
Following are some important points to consider before starting the deployment:
  - All configuration parameters are set to default values in the configuration files for ODIMRA and the GRF plugin. 
  - The following ports are used for deploying ODIMRA, the GRF plugin, and URP:
    45000, 45001, 45003, 45101-45110, 9092, 9082, 6380, 6379, 8500, 8300, 8302, 8301, 8600
    > **NOTE:** Ensure that the above ports are not in use.
  - The following users are created and added to group Ids automatically when the certificates are generated during the deployment. 
  
    |User Id| Group Id|
	-----|---------|
	|`odimra`|1234 |
	|`plugin`|1235 |
	

     `odimra` is created on both the system and the container for the resource aggregator.
	
	 `plugin` is created  on both the system and the container for the GRF plugin and URP.
	
	 > **NOTE:** Ensure that these user Ids and group Ids are not present on the system prior to deployment.




**Procedure**
--------------

**WARNING:** Do not run the commands provided in this section as a root user unless mentioned.


1. Clone the ODIMRA repository form `https://github.com/ODIM-Project/ODIM.git` to the home directory of the user:
   ```
   $ git clone https://github.com/ODIM-Project/ODIM.git
   ```
2. Choose a Fully Qualified Domain Name (FQDN) for the resource aggregator server. 
   Example: odim.local.com.
3. Set FQDN to the environment of the host machine using the following command:
    ```
    $ export FQDN=<user_preferred_fqdn_for_host>
    ```
4. Set the environment variable, `HOSTIP` to the IP address of your system:
   ```
   $ export HOSTIP=<ip_address_of_your_system>
   ```

5. Set the following environment variables for the ODIMRA user and group IDs:
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
   - Certificates are used by the resource aggregator services to communicate internally (Remote Procedure Call) and with the plugin services.
   - If you are using an intermediate CA for signing certificates assigned to the resource aggregator and the plugin services, ensure to:
        - Append all the intermediate certificates to the server certificate file in the order such that each certificate has signed the preceding one.
        - Append the Root CA used for signing the intermediate CA to the resource aggregator CA file.


    **Procedure**
   
   a. Navigate to the path: `ODIM/build/cert_generator`:
      ```
       $ cd ODIM/build/cert_generator
      ```

    > **NOTE:** `ODIM/build/cert_generator` contains the automated scripts to generate the TLS certificates for the resource aggregator, the GRF plugin, URP, and Kafka.

   b. Use the following command to generate certificates for the resource aggregator, the GRF plugin, and URP. Provide FQDN as a command-line argument.
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
   e. Use the following command to copy the TLS certificates of the resource aggregator, the GRF  plugin, URP, Kafka, and Zookeeper:
     ```
      $ sudo ./copy_certificate.sh
     ```
     The following files are copied in the path - `/etc/odimracert/`:
      - rootCA.crt
      - odimra_server.key
      - odimra_server.crt
      - odimra_rsa.public
      - odimra_rsa.private
      - odimra_kafka_client.key
      - odimra_kafka_client.crt

     The following files are copied in the path - `/etc/kafka/conf/`:
      - kafka.keystore.jks
      - kafka.truststore.jks

     The following files are copied in the path - `/etc/zookeeper/conf`:
      - zookeeper.keystore.jks
      - zookeeper.trustore.jks
      
    The following files are copied in the path - `/etc/plugincert/`:
      - rootCA.crt
      - odimra_server.key
      - odimra_server.crt
      - odimra_kafka_client.key
      - odimra_kafka_client.crt
      

8. Navigate to the ODIMRA folder.
   ```
   $ cd ~/ODIM
   ```

9. Use the following command to deploy and start the containers:
   ```
   $ make all
   ```
    The following containers are loaded:
      - build_odimra_1
      - build_kafka_1
      - build_zookeeper_1
      - build_redis_1
      - build_consul_1
      - build_grf_plugin_1
	  - build_urp_1

10. Verify that the resource aggregator services are running successfully.
   ```
   $ ps -eaf | grep svc
   ```
   
   All the resource aggregator services are listed:
   ```
   root     26491 30077  0 Jan10 ?        00:00:00 sudo -E -u odimra nohup ./svc-events --registry=consul --registry_address=consul:8500 --server_address=odimra:45103 --client_request_timeout=300s
   odimra   26499 26491  0 Jan10 ?        00:09:51 ./svc-events --registry=consul --registry_address=consul:8500 --server_address=odimra:45103 --client_request_timeout=300s
   root     30291 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-task --registry=consul --registry_address=consul:8500 --server_address=odimra:45105 --client_request_timeout=300s
   root     30301 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-update --registry=consul --registry_address=consul:8500 --server_address=odimra:45108 --client_request_timeout=300s
   root     30304 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-account-session --registry=consul --registry_address=consul:8500 --server_address=odimra:45101
   root     30306 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-fabrics --registry=consul --registry_address=consul:8500 --server_address=odimra:45106 --client_request_timeout=300s
   root     30314 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-systems --registry=consul --registry_address=consul:8500 --server_address=odimra:45104 --client_request_timeout=300s
   root     30321 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-managers --registry=consul --registry_address=consul:8500 --server_address=odimra:45107 --client_request_timeout=300s
   odimra   30326 30304  0 Jan07 ?        00:41:25 ./svc-account-session --registry=consul --registry_address=consul:8500 --server_address=odimra:45101
   root     30344 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-aggregation --registry=consul --registry_address=consul:8500 --server_address=odimra:45102 --client_request_timeout=300s
   root     30366 30077  0 Jan07 ?        00:00:00 sudo -E -u odimra nohup ./svc-api --registry=consul --registry_address=consul:8500 --client_request_timeout=302s
   odimra   30374 30301  0 Jan07 ?        00:01:42 ./svc-update --registry=consul --registry_address=consul:8500 --server_address=odimra:45108 --client_request_timeout=300s
   odimra   30375 30291  0 Jan07 ?        00:31:24 ./svc-task --registry=consul --registry_address=consul:8500 --server_address=odimra:45105 --client_request_timeout=300s
   odimra   30381 30344  0 Jan07 ?        00:04:27 ./svc-aggregation --registry=consul --registry_address=consul:8500 --server_address=odimra:45102 --client_request_timeout=300s
   odimra   30398 30306  0 Jan07 ?        00:01:25 ./svc-fabrics --registry=consul --registry_address=consul:8500 --server_address=odimra:45106 --client_request_timeout=300s
   odimra   30399 30314  0 Jan07 ?        00:03:13 ./svc-systems --registry=consul --registry_address=consul:8500 --server_address=odimra:45104 --client_request_timeout=300s
   odimra   30414 30321  0 Jan07 ?        00:01:14 ./svc-managers --registry=consul --registry_address=consul:8500 --server_address=odimra:45107 --client_request_timeout=300s
   odimra   30426 30366  0 Jan07 ?        00:17:52 ./svc-api --registry=consul --registry_address=consul:8500 --client_request_timeout=302s
   ```


  **NOTE:**
  - The configuration files of the resource aggregator are available at `/etc/odimra_config`.
  - The GRF configuration files are available at `/etc/grf_plugin_config`.
  - The URP configuration files are available at `/etc/urp_plugin_config`.
  - The resource aggregator API service runs on the default port 45000.
  - The GRF plugin API service runs on the default port 45001.
  - The URP API service runs on the default port 45003.
  - The resource aggregator logs are available at `/var/log/odimra`.
  - The GRF plugin logs are available at `/var/log/GRF_PLUGIN`.
  - To see the URP logs, run the following command:
    ```
    $ docker logs build_urp_1
	```


11. To configure log rotation, do the following:

    a. Navigate to the `/etc/logrotate.d` directory:
    ```
    $ cd /etc/logrotate.d
    ```
    b. Open the `odimra` file to edit:
	```
    $ sudo vi odimra
    ```	
	c. Add the following content and save:
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
    d. Navigate to the `/etc/cron.hourly` directory:
       ```
       $ cd /etc/cron.hourly
       ```
    e. Open the `logrotate` file:
	   ```
	   $ sudo vi logrotate
	   ```
	f. Add following content and save:
	   ```
       logrotate -s /var/lib/logrotate/status /etc/logrotate.d/odimra
	   ```
    e. To verify that the configuration is working, run the following command:
      ```
      $ sudo logrotate -v -f /etc/logrotate.d/odimra
  	  ```
  


12. To add the Generic Redfish Plugin, URP, and servers to the resource aggregator for ODIM, see "*Adding a plugin as an aggregation source*" and "*Adding a server as an aggregation source*" in the following readme:  
    https://github.com/ODIM-Project/ODIM/blob/development/svc-aggregation/README.md
	
	
### Default user credentials for ODIMRA, the GRF Plugin, and URP


ODIMRA:

```
Username: admin
Password: Od!m12$4
```

The GRF plugin:

```
Username: admin
Password: GRFPlug!n12$4
``` 

URP:

```
Username: admin
Password: Od!m12$4
``` 

 
 
  
#  Modifying default configuration parameters for the resource aggregator

1.   Navigate to the `build_odimra_1` container using the following command: 

      ```
     $ docker exec -it build_odimra_1 /bin/bash
     ```

2.   Edit the parameters in the `odimra_config.json` file located in this path: `/etc/odimra_config/odimra_config.json` and save. 

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
     |VerifyPeer|Boolean| Default value: true<br>  Recommended value: true<br>  NOTE:<br>  - `VerifyPeer` is set to true, by default. For secure plugin interaction, add root CA certificate (that is used to sign the certificates of the southbound entities) to root CA certificate file.  If `VerifyPeer` is set to false, SSL communication will be insecure.  After setting `VerifyPeer` to false, restart the resource aggregator container (`odim_1`).<br><br>  - If `TLS1.2` is used, ensure that the entity certificate has `SAN` field for successful validation.  - Northbound entities interacting with resource aggregator `API` service can use root CA that signed ODIMRA's certificate.|
     |PreferredCipherSuites}|List| Default and supported values: See "List of supported (default) cipher suites".<br> IMPORTANT:<br>  - If `PreferredCipherSuites` is not specified, it will be set to default cipher (secure) suites.<br><br>  - If `PreferredCipherSuites` is set to unsupported cipher suites, the resource aggregator and the plugin services will exit with errors.|

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
|`<ODIM_server_IP>` |The IP address of the system where the resource aggregator is installed.|
|`<FQDN>` |FQDN of the resource aggregator server.|

</blockquote>

**Procedure**
--------------

1.   In the home directory of the logged in user, create a hidden directory called *.docker*, and then create a file called *config.json* inside it.

      ```
       mkdir .docker
      ```

      ```
       cd .docker
      ```

      ```
       vi config.json
      ```	   

2.   Add the following content in the `config.json` file and save: 

      ```
      {
         "proxies":
        {
           "default":
          {
             "httpProxy": "<Proxy_URL>",
             "httpsProxy": "<Proxy_URL>",
             "noProxy": "localhost,127.0.0.1, <ODIM_server_IP>"
          }
        }
      }
    
     ```

3.   Update the `/etc/environment` file with the following content using sudo: 

      ```
      export http_proxy=<Proxy_URL>
      export https_proxy=<Proxy_URL>
      HOSTIP=<ODIM_server_IP>
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
         Environment="NO_PROXY=localhost,127.0.0.1, <ODIM_server_IP>"
        
         ```

     4. Run the following commands: 

         ```
         $ sudo systemctl daemon-reload
         ```

         ```
         $ sudo service docker restart

         ```
		 
# Uninstalling ODIMRA

  To uninstall ODIMRA, use either of the two commands listed in this section.
  
1. ```
   $ make clean
   ```

   Use the following command to:
   - Remove all the deployed Docker containers.
   - Remove only those Docker images which were created and deployed as containers.
   - Remove data stored by Consul, Redis, and Kafka.	
	
    When prompted for password, enter the sudo password.
	 
  
2. ```
   $ make deepclean
   ``` 
  
   Use the following command to:
   - Remove all the deployed Docker containers.
   - Remove all the Docker images including the intermediate or dependent images created during the deployment.
   - Remove configuration information and data stored by Consul, Redis, and Kafka.
   - Remove all generated certificates.
   - Remove log files created for the ODIMRA services, the GRF plugin, and URP.
   
    When prompted for password, enter the sudo password.

>**CAUTION**:
 Running these commands will unistall ODIMRA and all the data related to it completely. It is best not to run these commands unless absolutely necessary.
               

# CI Process

GitHub action workflows, also called as checks, are added to the ODIM repository. They are triggered whenever a Pull Request(PR) is raised against the development branch.
The result from the workflow execution is then updated to the PR.
 
>**Note:** You can review and merge PRs only if the checks are passed.

Following checks are added as part of the CI process:

|Sl No.|Workflow Name|Description|
|---------|-----------|----------|
|1|`build_unittest.yml` |Builds and runs Unit Tests with code coverage enabled.|
|2|`build_deploy_test.yml` |Builds, Deploys, runs sanity tests, and uploads build artifacts (like odimra logs).|
|3|`LGTM analysis` |Semantic code analyzer and query tool which finds security vulnerabilities in codebases.| 

These checks run in parallel and take approximately nine minutes to complete.

GitHub action workflow details
------------------------------

1. build_unittest.yml
   - Brings up a Ubuntu 18.04 VM hosted on GitHub infrastructure with preinstalled packages mentioned in the link: https://github.com/actions/virtual-environments/blob/master/images/linux/Ubuntu1804-README.md.
   - Installs Go 1.13.8 package.
   - Installs and configures Redis 5.0.8 with two instances running on ports 6379 and 6380.
   - Checks out the PR code into the Go module directory.
   - Builds the code.
   - Runs the unit tests.

2. build_deploy_test.yml
   - Brings up a Ubuntu 18.04 VM hosted on GitHub infrastructure with preinstalled packages mentioned in the link: https://github.com/actions/virtual-environments/blob/master/images/linux/Ubuntu1804-README.md.
   - Checks out the PR code.
   - Builds and deploys the following docker containers
     ODIMRA
     Generic redfish plugin
     Unmanaged Rack Plugin
     Kakfa
     Zookeeper
     Consul
     Redisdb
   - Runs the sanity tests.
   - Uploads the build artifacts.

> **NOTE:** Build status notifications having a link to the GitHub Actions' build job page will be sent to the developer’s email.
