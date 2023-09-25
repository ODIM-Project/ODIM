#!/bin/bash
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

#  @Description:
#    EntrypointSetup script for deploying redis HA via Sentinel in a kubernetes cluster
#    This script expects following environment variables to be set,
#    1. SENTINEL: true if this is sentinel instance, else false.
#    2. PRIMARY: true if this is primary instance, this is helpful when starting the cluster for the first time.
#    3. REDIS_HA_SENTINEL_SERVICE_HOST: this is service name of sentinel, check the yaml.
#    4. REDIS_HA_SENTINEL_SERVICE_PORT: this is service port of sentinel.
#    5. REDIS_HA_REDIS_SERVICE_HOST: this is primary's service name, this is needed when sentinel starts for the first time.
#    6. REDIS_HA_REDIS_SERVICE_PORT: this is primary's port, is needed when sentinel starts for the first time.
#    7. REDIS_DEFAULT_PASSWORD: default password for Redis instances.

#  This method launches redis instance which assumes it self as primary
function launchprimary() {
  echo "Starting Redis instance as Primary.."
  redis_password=$(openssl pkeyutl -decrypt -in cipher -inkey ${ODIMRA_RSA_PRIVATE_FILE} -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512)
  echo "while true; do   sleep 2;   export primary=\$(hostname -I | cut -d ' ' -f 1);   echo \"Primary IP is Me : \${primary}\";   echo \"Setting STARTUP_PRIMARY_IP in redis\";   redis-cli -a '${redis_password}' -h \${primary} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} set STARTUP_PRIMARY_IP \${primary};   if [ \$? == \"0\" ]; then     echo \"Successfully set STARTUP_PRIMARY_IP\"; if [ \${REDIS_ONDISK_DB} == \"true\" ]; then     bash \/createschema.sh; fi;   break;   fi;   echo \"Connecting to primary \${primary} failed.  Waiting...\";   sleep 5; done" > insert_primary_ip_and_default_entries.sh
  bash insert_primary_ip_and_default_entries.sh &
  sed -i "s/REDIS_DEFAULT_PASSWORD/${redis_password}/" /redis-primary/redis.conf
  hostname=$(hostname -f)
  sed -i "s/%replica-announce-ip%/${hostname}/" /redis-primary/redis.conf
  sed -i "s/%primary-port%/${REDIS_HA_REDIS_SERVICE_PORT}/" /redis-primary/redis.conf

  redis-server /redis-primary/redis.conf --protected-mode no
}

#  This method launches sentinels
function launchsentinel() {
  echo "Starting Sentinel.."
  sleep_for_rand_int=$(awk -v min=2 -v max=7 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
  sleep ${sleep_for_rand_int}

  echo -n "${REDIS_DEFAULT_PASSWORD}" | base64 --decode > cipher
  redis_password=$(openssl pkeyutl -decrypt -in cipher -inkey ${ODIMRA_RSA_PRIVATE_FILE} -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512)
  x=1
  while [ $x -le 5 ]
  do
    primary=$(redis-cli -a ${redis_password} -h ${REDIS_HA_SENTINEL_SERVICE_HOST} -p ${REDIS_HA_SENTINEL_SERVICE_PORT} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} --csv SENTINEL get-primary-addr-by-name ${REDIS_PRIMARY_SET} | tr ',' ' ' | cut -d' ' -f1)
    if [[ -n ${primary} ]]; then
      echo "Connected to Sentinel Service and retrieved Redis Primary hostname as ${primary}"
      primary="${primary//\"}"
      break
    else
      echo "Unable to connect to sentinel, retrying..."
      sleep 1
    fi
    x=$(( $x + 1 ))
  done

  if ! [[ -n ${primary} ]]; then
    echo "Unable to connect to Sentinel Service, probably because I am first Sentinel to start. I will use default primary hostname ${PRIMARY_HOST_NAME} to connect to sentinel"
    primary=${PRIMARY_HOST_NAME}
  fi

  while true; do
    redis-cli -a ${redis_password} -h ${primary} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} INFO
    if [[ "$?" == "0" ]]; then
      break
    fi
    echo "Connecting to primary failed.  Waiting..."
    sleep 10
  done

  sentinel_conf=sentinel.conf

  hostname=$(hostname -f)

  echo "sentinel resolve-hostnames yes" >> ${sentinel_conf}
  echo "sentinel announce-hostnames yes" >> ${sentinel_conf}
  echo "sentinel announce-ip ${hostname}" >> ${sentinel_conf}
  echo "sentinel announce-port ${REDIS_HA_SENTINEL_SERVICE_PORT}" >> ${sentinel_conf}
  echo "sentinel monitor ${REDIS_PRIMARY_SET} ${primary} ${REDIS_HA_REDIS_SERVICE_PORT} ${SENTINEL_QUORUM}" >> ${sentinel_conf}
  echo "sentinel auth-pass ${REDIS_PRIMARY_SET} ${redis_password}" >> ${sentinel_conf}
  echo "requirepass ${redis_password}" >> ${sentinel_conf}
  echo "sentinel down-after-milliseconds ${REDIS_PRIMARY_SET} ${DOWN_AFTER_MILLISECONDS}" >> ${sentinel_conf}
  echo "sentinel failover-timeout ${REDIS_PRIMARY_SET} ${FAILOVER_TIMEOUT}" >> ${sentinel_conf}
  echo "sentinel parallel-syncs ${REDIS_PRIMARY_SET} ${PARALLEL_SYNCS}" >> ${sentinel_conf}
  echo "bind 0.0.0.0" >> ${sentinel_conf}
  echo "port 0" >> ${sentinel_conf}
  echo "tls-port 26379" >> ${sentinel_conf}
  echo "tls-replication yes" >> ${sentinel_conf}
  echo "tls-cluster yes" >> ${sentinel_conf}
  echo "tls-cert-file /etc/odimra_certs/odimra_server.crt" >> ${sentinel_conf}
  echo "tls-key-file /etc/odimra_certs/odimra_server.key" >> ${sentinel_conf}
  echo "tls-ca-cert-file /etc/odimra_certs/rootCA.crt" >> ${sentinel_conf}

  redis-sentinel ${sentinel_conf} --protected-mode no
}

#  This method launches secondry instances
function launchsecondry() {
  echo "Starting Redis instance as Secondry , Primary IP $1"

  redis_password=$(openssl pkeyutl -decrypt -in cipher -inkey ${ODIMRA_RSA_PRIVATE_FILE} -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512)
  while true; do
    echo "Trying to retrieve the Primary IP again, in case of failover primary ip would have changed."
    Primary=$(redis-cli -a ${redis_password} -h ${REDIS_HA_SENTINEL_SERVICE_HOST} -p ${REDIS_HA_SENTINEL_SERVICE_PORT} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} --csv SENTINEL get-primary-addr-by-name ${REDIS_PRIMARY_SET} | tr ',' ' ' | cut -d' ' -f1)

    if [[ -n ${primary} ]]; then
      primary="${primary//\"}"
    else
      echo "Failed to find primary."
      sleep 60
      continue
    fi
    redis-cli -a ${redis_password} -h ${primary} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} INFO
    if [[ "$?" == "0" ]]; then
      break
    fi
    echo "Connecting to primary failed.  Waiting..."
    sleep 10
  done

  hostname=$(hostname -f)
  sed -i "s/%primary-ip%/${primary}/" /redis-secondry/redis.conf
  sed -i "s/%primary-port%/${REDIS_HA_REDIS_SERVICE_PORT}/" /redis-secondry/redis.conf
  sed -i "s/REDIS_DEFAULT_PASSWORD/${redis_password}/" /redis-secondry/redis.conf
  sed -i "s/%replica-announce-ip%/${hostname}/" /redis-secondry/redis.conf
  sed -i "s/%replicaof%/${primary}/" /redis-secondry/redis.conf

  redis-server /redis-secondry/redis.conf --protected-mode no
}


#  This method launches either secondry or primary based on some parameters
function launchredis() {
  echo "Launching Redis instance"

  hostname=$(hostname -f)
  echo -n "${REDIS_DEFAULT_PASSWORD}" | base64 --decode > cipher
  redis_password=$(openssl pkeyutl -decrypt -in cipher -inkey ${ODIMRA_RSA_PRIVATE_FILE} -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha512)

  # If it is sentinel restart, I am giving some time to sentinel for complete shutdown
  sentinel_down_time=10
  sleep ${sentinel_down_time}

  # Loop till I am able to launch secondry or primary
  while true; do
    # I will check if sentinel is up or not by connecting to it.
    echo "Trying to connect to sentinel, to retireve primary's ip"
      primary=$(redis-cli -a ${redis_password} -h ${REDIS_HA_SENTINEL_SERVICE_HOST} -p ${REDIS_HA_SENTINEL_SERVICE_PORT} --tls --cert ${TLS_CERT_FILE} --key ${TLS_KEY_FILE} --cacert ${TLS_CA_CERT_FILE} --csv SENTINEL get-primary-addr-by-name ${REDIS_PRIMARY_SET} | tr ',' ' ' | cut -d' ' -f1)
    # Is this instance marked as PRIMARY, it will matter only when the cluster is starting up for first time.
    if [[ "${PRIMARY}" == "true" ]]; then
      echo "PRIMARY is set to true"
      # If I am able get primary ip, then i will connect to the primary, else i will asume the role of primary
      if [[ -n ${primary} ]]; then
        echo "Connected to Sentinel, this means it is not first time start, hence will start as a secondry"
        currenthost=$(hostname -f | cut -d ' ' -f 1)
	      primary=`echo $primary |tr -d '"'`
        if [[ "${currenthost}" == "${primary}" ]]; then
           launchprimary
           exit 0
	      fi    
        launchsecondry ${primary}
        exit 0
      else
        launchprimary
        exit 0
      fi
    fi

    # If I am not primary, then i am definitely secondry.
    if [[ -n ${primary} ]]; then
      echo "Connected to Sentinel and Retrieved Primary IP ${primary}"
      launchsecondry ${primary}
      exit 0   
    else
      echo "Connecting to sentinel failed, Waiting..."
      sleep 10
    fi
  done
}

if [[ "${SENTINEL}" == "true" ]]; then
  launchsentinel
  exit 0
fi

launchredis
