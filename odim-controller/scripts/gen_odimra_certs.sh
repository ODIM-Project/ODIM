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

declare ODIMRA_CERT_DIR
declare ODIM_CONTROLLER_CONFIG
declare ODIMRA_FQDN
declare KAFKA_JKS_PASSWORD
declare ZOOKEEPER_JKS_PASSWORD
declare ODIM_SERVER_FQDN_SAN
declare ODIM_SERVER_IP_SAN
declare ODIM_KAFKA_CLIENT_FQDN_SAN
declare ODIM_KAFKA_CLIENT_IP_SAN
declare ODIMRA_ROOTCA_CRT_PATH
declare ODIMRA_ROOTCA_KEY_PATH
declare ODIMRA_SERVER_CSR_PATH
declare ODIMRA_SERVER_CRT_PATH
declare ODIMRA_SERVER_KEY_PATH
declare ODIMRA_RSA_PUBLIC_KEY
declare ODIMRA_RSA_PRIVATE_KEY
declare ODIMRA_KAFKA_CLIENT_CSR_PATH
declare ODIMRA_KAFKA_CLIENT_CRT_PATH
declare ODIMRA_KAFKA_CLIENT_KEY_PATH
declare KAFKA_JKS_PATH
declare KAFKA_JTS_PATH
declare ZOOKEEPER_JKS_PATH
declare ZOOKEEPER_JTS_PATH
declare ODIMRA_NAMESPACE
declare ODIMRA_ETCD_SERVER_CSR_PATH
declare ODIMRA_ETCD_SERVER_CRT_PATH
declare ODIMRA_ETCD_SERVER_KEY_PATH
declare ODIMRA_HA_DEPLOYMENT
declare ODIMRA_VIRTUAL_IP
declare NGINX_SERVER_CSR_PATH
declare NGINX_SERVER_CRT_PATH
declare NGINX_SERVER_KEY_PATH

OPENSSL_BIN_PATH="/usr/bin/openssl"
KEYTOOL_BIN_PATH="/usr/bin/keytool"
CERT_COUNTRY_NAME="US"
CERT_STATE_NAME="CA"
CERT_LOCALITY_NAME="California"
CERT_ORGANIZATION_NAME="HPE"
CERT_ORGANIZATION_UNIT_NAME="Telco Solutions"
CERT_VALIDITY_PERIOD=3650
KEY_LENGTH=4096
ROOT_CA_KEY_NOT_PROVIDED=false

pre_reqs()
{
	if [[ ! -e ${OPENSSL_BIN_PATH} ]]; then
		echo "[$(date)] -- ERROR -- ${OPENSSL_BIN_PATH} does not exist"
		exit 1
	fi

	if [[ ! -e ${KEYTOOL_BIN_PATH} ]]; then
		echo "[$(date)] -- ERROR -- ${KEYTOOL_BIN_PATH} does not exist"
		exit 1
	fi
}

eval_cmd_exec()
{
        if [[ $# -lt 2 ]]; then
                echo "[$(date)] -- ERROR -- eval_cmd_exec syntax error $2"
		exit 1
        fi
        if [[ $1 -ne 0 ]]; then
                echo "[$(date)] -- ERROR -- $2"
		exit 1
        fi
}

# generate_ca_certs is for generating rootCA key pair
# required for signing certificates of ODIM-RA services
generate_ca_certs()
{
	# check if cert and key exists and is not empty
	if [[ -s ${ODIMRA_ROOTCA_CRT_PATH} ]]; then
		if [[ -s ${ODIMRA_ROOTCA_KEY_PATH} ]]; then
			echo "[$(date)] -- INFO  -- rootCA crt and key exists, not generating again"
		else
			echo "[$(date)] -- INFO  -- rootCA crt present but not the key, expects every required crt to be present"
			ROOT_CA_KEY_NOT_PROVIDED=true
		fi
		return
	fi

	# generate rootCA private key
	${OPENSSL_BIN_PATH} genrsa -out ${ODIMRA_ROOTCA_KEY_PATH} ${KEY_LENGTH}
	eval_cmd_exec $? "${ODIMRA_ROOTCA_KEY_PATH} generation failed"

	# Have rootCA validity longer than the validity of the certs signed by it
	rootCA_Validity_Period=$((${CERT_VALIDITY_PERIOD} * 2))
	# generate rootCA certificate
	${OPENSSL_BIN_PATH} req -new -key ${ODIMRA_ROOTCA_KEY_PATH} -days ${rootCA_Validity_Period} -x509 -out ${ODIMRA_ROOTCA_CRT_PATH} -config <( cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
x509_extensions = v3_req_ca

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_ROOT_CA

[v3_req_ca]
basicConstraints        = critical,CA:true
subjectKeyIdentifier    = hash
keyUsage                = critical, digitalSignature, nonRepudiation, keyEncipherment, cRLSign, keyCertSign
authorityKeyIdentifier  = keyid:always,issuer
EOF
)
	eval_cmd_exec $? "${ODIMRA_ROOTCA_CRT_PATH} generation failed"
	chmod 0600 ${ODIMRA_ROOTCA_KEY_PATH} ${ODIMRA_ROOTCA_CRT_PATH}
}

# generate_odim_server_certs is for generating ODIM-RA
# services certificate and private key
generate_odim_server_certs()
{
	# check if cert and key exists and is not empty
	if [[ -s ${ODIMRA_SERVER_KEY_PATH} ]] && [[ -s ${ODIMRA_SERVER_CRT_PATH} ]]; then
		echo "[$(date)] -- INFO  -- odimra server crt and key already exists"
		# verify crt was signed by the rootCA present
		${OPENSSL_BIN_PATH} verify -CAfile ${ODIMRA_ROOTCA_CRT_PATH} ${ODIMRA_SERVER_CRT_PATH} > /dev/null
		eval_cmd_exec $? "${ODIMRA_SERVER_CRT_PATH} is not signed by ${ODIMRA_ROOTCA_CRT_PATH}"
		return
	fi

	# check if rootCA key was made available, if it was not generated this script.
	if ${ROOT_CA_KEY_NOT_PROVIDED}; then
		echo "[$(date)] -- ERROR -- rootCA key was not provided, odim server crt generation not possible"
		exit 1
	fi

	# generate private key
	${OPENSSL_BIN_PATH} genrsa -out ${ODIMRA_SERVER_KEY_PATH} ${KEY_LENGTH}
	eval_cmd_exec $? "${ODIMRA_SERVER_KEY_PATH} generation failed"
	
	#generate CSR
	${OPENSSL_BIN_PATH} req -new -key ${ODIMRA_SERVER_KEY_PATH} -out ${ODIMRA_SERVER_CSR_PATH} -config <(cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_SVC_CRT

[v3_req]
subjectKeyIdentifier = hash
keyUsage             = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage     = clientAuth, serverAuth
subjectAltName       = @alt_names

[alt_names]
$ODIM_SERVER_FQDN_SAN
$ODIM_SERVER_IP_SAN
EOF
)
	eval_cmd_exec $? "${ODIMRA_SERVER_CSR_PATH} generation failed"

	# obtain certificate
	${OPENSSL_BIN_PATH} x509 -req -days ${CERT_VALIDITY_PERIOD} -in ${ODIMRA_SERVER_CSR_PATH} -CA ${ODIMRA_ROOTCA_CRT_PATH} -CAkey ${ODIMRA_ROOTCA_KEY_PATH} -CAcreateserial -out ${ODIMRA_SERVER_CRT_PATH} -extensions v3_req -extfile <( cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_SVC_CRT

[v3_req]
subjectKeyIdentifier    = hash
authorityKeyIdentifier  = keyid:always,issuer:always
keyUsage                = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage        = clientAuth, serverAuth
subjectAltName          = @alt_names

[alt_names]
$ODIM_SERVER_FQDN_SAN
$ODIM_SERVER_IP_SAN
EOF
)
	eval_cmd_exec $? "${ODIMRA_SERVER_CRT_PATH} generation failed"

	# remove temp files
	rm -f ${ODIMRA_SERVER_CSR_PATH} ${ODIMRA_CERT_DIR}/rootCA.srl

	if [[ -s ${ODIMRA_RSA_PRIVATE_KEY} ]] && [[ -s ${ODIMRA_RSA_PUBLIC_KEY} ]]; then
		echo "[$(date)] -- INFO  -- odim RSA private and public keys already exists"
		return
	fi

	# generate RSA private key
	${OPENSSL_BIN_PATH} genrsa -out ${ODIMRA_RSA_PRIVATE_KEY} ${KEY_LENGTH}
	eval_cmd_exec $? "${ODIMRA_RSA_PRIVATE_KEY} generation failed"

	# generate RSA public key
	${OPENSSL_BIN_PATH} rsa -in ${ODIMRA_RSA_PRIVATE_KEY} -out ${ODIMRA_RSA_PUBLIC_KEY} -pubout -outform PEM
	eval_cmd_exec $? "${ODIMRA_RSA_PUBLIC_KEY} generation failed"

	chmod 0600  ${ODIMRA_SERVER_KEY_PATH} ${ODIMRA_SERVER_CRT_PATH} ${ODIMRA_RSA_PRIVATE_KEY} ${ODIMRA_RSA_PUBLIC_KEY}
}

# generate_odim_kafka_client_certs is for generating
# client certificate and private key to be used by
# ODIM-RA services for interacting with kafka
generate_odim_kafka_client_certs()
{
	# check if cert and key exists and is not empty
	if [[ -s ${ODIMRA_KAFKA_CLIENT_KEY_PATH} ]] && [[ -s ${ODIMRA_KAFKA_CLIENT_CRT_PATH} ]]; then
		echo "[$(date)] -- INFO  -- odimra kafka client crt and key already exists"
		# verify crt was signed by the rootCA present
		${OPENSSL_BIN_PATH} verify -CAfile ${ODIMRA_ROOTCA_CRT_PATH} ${ODIMRA_KAFKA_CLIENT_CRT_PATH} > /dev/null
		eval_cmd_exec $? "${ODIMRA_KAFKA_CLIENT_CRT_PATH} is not signed by ${ODIMRA_ROOTCA_CRT_PATH}"
		return
	fi

	# check if rootCA key was made available, if it was not generated this script.
	if ${ROOT_CA_KEY_NOT_PROVIDED}; then
		echo "[$(date)] -- ERROR -- rootCA key was not provided, odim kafka client crt generation not possible"
		exit 1
	fi

	# generate private key
	${OPENSSL_BIN_PATH} genrsa -out ${ODIMRA_KAFKA_CLIENT_KEY_PATH} ${KEY_LENGTH}
	eval_cmd_exec $? "${ODIMRA_KAFKA_CLIENT_KEY_PATH} generation failed"
	
	#generate CSR
	${OPENSSL_BIN_PATH} req -new -key ${ODIMRA_KAFKA_CLIENT_KEY_PATH} -out ${ODIMRA_KAFKA_CLIENT_CSR_PATH} -config <(cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_KAFKA_CLIENT_CRT

[v3_req]
subjectKeyIdentifier = hash
keyUsage             = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage     = clientAuth
subjectAltName       = @alt_names

[alt_names]
$ODIM_KAFKA_CLIENT_FQDN_SAN
$ODIM_KAFKA_CLIENT_IP_SAN
EOF
)
	eval_cmd_exec $? "${ODIMRA_KAFKA_CLIENT_CSR_PATH} generation failed"

	# obtain certificate
	${OPENSSL_BIN_PATH} x509 -req -days ${CERT_VALIDITY_PERIOD} -in ${ODIMRA_KAFKA_CLIENT_CSR_PATH} -CA ${ODIMRA_ROOTCA_CRT_PATH} -CAkey ${ODIMRA_ROOTCA_KEY_PATH} -CAcreateserial -out ${ODIMRA_KAFKA_CLIENT_CRT_PATH} -extensions v3_req -extfile <( cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_KAFKA_CLIENT_CRT

[v3_req]
subjectKeyIdentifier    = hash
authorityKeyIdentifier  = keyid:always,issuer:always
keyUsage                = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage        = clientAuth
subjectAltName          = @alt_names

[alt_names]
$ODIM_KAFKA_CLIENT_FQDN_SAN
$ODIM_KAFKA_CLIENT_IP_SAN
EOF
)
	eval_cmd_exec $? "${ODIMRA_KAFKA_CLIENT_CRT_PATH} generation failed"

	# remove temp files
	rm -f ${ODIMRA_KAFKA_CLIENT_CSR_PATH} ${ODIMRA_CERT_DIR}/rootCA.srl

	chmod 0600 ${ODIMRA_KAFKA_CLIENT_KEY_PATH} ${ODIMRA_KAFKA_CLIENT_CRT_PATH}
}

# generate_kafka_certs is for generating kafka
# certificate and private key, and storing the same in
# java keystore
generate_kafka_certs()
{
	# check if cert and key exists and is not empty
	if [[ -s ${KAFKA_JKS_PATH} ]] && [[ -s ${KAFKA_JTS_PATH} ]]; then
		echo "[$(date)] -- INFO  -- kafka keystore and truststore already exists"
		return
	fi

	# check if rootCA key was made available, if it was not generated this script.
	if ${ROOT_CA_KEY_NOT_PROVIDED}; then
		echo "[$(date)] -- ERROR -- rootCA key was not provided, kafka keystore generation not possible"
		exit 1
	fi

	# generate keystore
	${KEYTOOL_BIN_PATH} -keystore ${KAFKA_JKS_PATH} -storetype pkcs12 -alias kafka -validity ${CERT_VALIDITY_PERIOD} -keyalg RSA -genkey <<HERE
${KAFKA_JKS_PASSWORD}
${KAFKA_JKS_PASSWORD}
kafka
${CERT_ORGANIZATION_UNIT_NAME}
${CERT_ORGANIZATION_NAME}
${CERT_LOCALITY_NAME}
${CERT_STATE_NAME}
${CERT_COUNTRY_NAME}
yes
.
HERE
	eval_cmd_exec $? "${KAFKA_JKS_PATH} generation failed"

	# add kafka rootCA to kafka server truststore
	${KEYTOOL_BIN_PATH} -keystore ${KAFKA_JTS_PATH} -storetype pkcs12 -alias rootCA -import -file ${ODIMRA_ROOTCA_CRT_PATH} <<HERE
${KAFKA_JKS_PASSWORD}
${KAFKA_JKS_PASSWORD}
yes
HERE
	eval_cmd_exec $? "${KAFKA_JTS_PATH} generation failed"

	# generate keystore CSR
	${KEYTOOL_BIN_PATH} -keystore ${KAFKA_JKS_PATH} -alias kafka -certreq -file ${ODIMRA_CERT_DIR}/kafka.csr <<HERE
${KAFKA_JKS_PASSWORD}
HERE
	eval_cmd_exec $? "${ODIMRA_CERT_DIR}/kafka.csr generation failed"

	#generate the keystore certificates
	${OPENSSL_BIN_PATH} x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt] 
basicConstraints=CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName = @alternate_names

[ alternate_names ]
DNS.1 = kafka
DNS.2 = kafka1.kafka.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.3 = kafka2.kafka.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.4 = kafka3.kafka.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.5 = kafka-ext
DNS.6 = kafka1-ext
DNS.7 = kafka2-ext
DNS.8 = kafka3-ext
EOF
) -in ${ODIMRA_CERT_DIR}/kafka.csr -CA ${ODIMRA_ROOTCA_CRT_PATH} -CAkey ${ODIMRA_ROOTCA_KEY_PATH} -CAcreateserial -out ${ODIMRA_CERT_DIR}/kafka.crt -days ${CERT_VALIDITY_PERIOD} -sha512
	eval_cmd_exec $? "${ODIMRA_CERT_DIR}/kafka.crt generation failed"

	# add kafka rootCA to kafka server keystore
	${KEYTOOL_BIN_PATH} -keystore ${KAFKA_JKS_PATH} -alias rootCA -import -file ${ODIMRA_ROOTCA_CRT_PATH} <<HERE
${KAFKA_JKS_PASSWORD}
yes
HERE
	eval_cmd_exec $? "CA file import to keystore failed"

	# adding server certificate to server key store
	${KEYTOOL_BIN_PATH} -keystore ${KAFKA_JKS_PATH} -alias kafka -import -file ${ODIMRA_CERT_DIR}/kafka.crt <<HERE
${KAFKA_JKS_PASSWORD}
HERE
	eval_cmd_exec $? "Kafka certificate import to keystore failed"

	# clean up temp files generated
	rm -f ${ODIMRA_CERT_DIR}/kafka.csr ${ODIMRA_CERT_DIR}/kafka.crt ${ODIMRA_CERT_DIR}/rootCA.srl

	chmod 0600 ${KAFKA_JKS_PATH} ${KAFKA_JTS_PATH}
}

# generate_zookeeper_certs is for generating zookeeper
# certificate and private key, and storing the same in
# java keystore
generate_zookeeper_certs()
{
	# check if cert and key exists and is not empty
	if [[ -s ${ZOOKEEPER_JKS_PATH} ]] && [[ -s ${ZOOKEEPER_JTS_PATH} ]]; then
		echo "[$(date)] -- INFO  -- zookeeper keystore and truststore already exists"
		return
	fi

	# check if rootCA key was made available, if it was not generated this script.
	if ${ROOT_CA_KEY_NOT_PROVIDED}; then
		echo "[$(date)] -- ERROR -- rootCA key was not provided, zookeeper keystore generation not possible"
		exit 1
	fi

	# generate keystore
	${KEYTOOL_BIN_PATH} -keystore ${ZOOKEEPER_JKS_PATH} -storetype pkcs12 -alias zookeeper -validity ${CERT_VALIDITY_PERIOD} -keyalg RSA -genkey <<HERE
${ZOOKEEPER_JKS_PASSWORD}
${ZOOKEEPER_JKS_PASSWORD}
zookeeper
${CERT_ORGANIZATION_UNIT_NAME}
${CERT_ORGANIZATION_NAME}
${CERT_LOCALITY_NAME}
${CERT_STATE_NAME}
${CERT_COUNTRY_NAME}
yes
.
HERE
	eval_cmd_exec $? "${ZOOKEEPER_JKS_PATH} generation failed"

	# add zookeeper rootCA to zookeeper server truststore
	${KEYTOOL_BIN_PATH} -keystore ${ZOOKEEPER_JTS_PATH} -storetype pkcs12 -alias rootCA -import -file ${ODIMRA_ROOTCA_CRT_PATH} <<HERE
${ZOOKEEPER_JKS_PASSWORD}
${ZOOKEEPER_JKS_PASSWORD}
yes
HERE
	eval_cmd_exec $? "${ZOOKEEPER_JTS_PATH} generation failed"

	# generate keystore CSR
	${KEYTOOL_BIN_PATH} -keystore ${ZOOKEEPER_JKS_PATH} -alias zookeeper -certreq -file ${ODIMRA_CERT_DIR}/zookeeper.csr <<HERE
${ZOOKEEPER_JKS_PASSWORD}
HERE
	eval_cmd_exec $? "${ODIMRA_CERT_DIR}/zookeeper.csr generation failed"

	#generate the keystore certificates
	${OPENSSL_BIN_PATH} x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt]
basicConstraints=CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName = @alternate_names

[ alternate_names ]
DNS.1 = zookeeper
DNS.2 = zookeeper1.zookeeper.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.3 = zookeeper2.zookeeper.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.4 = zookeeper3.zookeeper.${ODIMRA_NAMESPACE}.svc.cluster.local
EOF
) -in ${ODIMRA_CERT_DIR}/zookeeper.csr -CA ${ODIMRA_ROOTCA_CRT_PATH} -CAkey ${ODIMRA_ROOTCA_KEY_PATH} -CAcreateserial -out ${ODIMRA_CERT_DIR}/zookeeper.crt -days ${CERT_VALIDITY_PERIOD} -sha512
	eval_cmd_exec $? "${ODIMRA_CERT_DIR}/zookeeper.crt generation failed"

	# add zookeeper rootCA to zookeeper server keystore
	${KEYTOOL_BIN_PATH} -keystore ${ZOOKEEPER_JKS_PATH} -alias rootCA -import -file ${ODIMRA_ROOTCA_CRT_PATH} <<HERE
${ZOOKEEPER_JKS_PASSWORD}
yes
HERE
	eval_cmd_exec $? "CA file import to keystore failed"

	# adding server certificate to server key store
	${KEYTOOL_BIN_PATH} -keystore ${ZOOKEEPER_JKS_PATH} -alias zookeeper -import -file ${ODIMRA_CERT_DIR}/zookeeper.crt <<HERE
${ZOOKEEPER_JKS_PASSWORD}
HERE
	eval_cmd_exec $? "Zookeeper certificate import to keystore failed"

	# clean up temp files generated
	rm -f ${ODIMRA_CERT_DIR}/zookeeper.csr ${ODIMRA_CERT_DIR}/zookeeper.crt ${ODIMRA_CERT_DIR}/rootCA.srl

	chmod 0600 ${ZOOKEEPER_JKS_PATH} ${ZOOKEEPER_JTS_PATH}
}

# generate_nginx_certs is for generating certificate
# and private key required for nginx server
generate_nginx_certs()
{
	# check HA deployment is enabled
	if [[ ${ODIMRA_HA_DEPLOYMENT,,} == false ]]; then
		echo "[$(date)] -- INFO  -- HA deployment not enabled, not generating cert and key required for nginx"
		return
	fi

	# check if cert and key exists and is not empty
	if [[ -s ${NGINX_SERVER_KEY_PATH} ]] && [[ -s ${NGINX_SERVER_CRT_PATH} ]]; then
		echo "[$(date)] -- INFO  -- nginx server crt and key already exists"
		# verify crt was signed by the rootCA present
		${OPENSSL_BIN_PATH} verify -CAfile ${ODIMRA_ROOTCA_CRT_PATH} ${NGINX_SERVER_CRT_PATH} > /dev/null
		eval_cmd_exec $? "${NGINX_SERVER_CRT_PATH} is not signed by ${ODIMRA_ROOTCA_CRT_PATH}"
		return
	fi

	# check if rootCA key was made available, if it was not generated this script.
	if ${ROOT_CA_KEY_NOT_PROVIDED}; then
		echo "[$(date)] -- ERROR -- rootCA key was not provided, nginx server crt generation not possible"
		exit 1
	fi

	# generate private key
	${OPENSSL_BIN_PATH} genrsa -out ${NGINX_SERVER_KEY_PATH} ${KEY_LENGTH}
	eval_cmd_exec $? "${NGINX_SERVER_KEY_PATH} generation failed"

	#generate CSR
	${OPENSSL_BIN_PATH} req -new -key ${NGINX_SERVER_KEY_PATH} -out ${NGINX_SERVER_CSR_PATH} -config <(cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_PROXY_CRT

[v3_req]
subjectKeyIdentifier = hash
keyUsage             = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage     = clientAuth, serverAuth
subjectAltName       = @alt_names

[alt_names]
DNS.0 = odimra.proxy.net
IP.0 = ${ODIMRA_VIRTUAL_IP}
EOF
)
	eval_cmd_exec $? "${NGINX_SERVER_CSR_PATH} generation failed"

	# obtain certificate
	${OPENSSL_BIN_PATH} x509 -req -days ${CERT_VALIDITY_PERIOD} -in ${NGINX_SERVER_CSR_PATH} -CA ${ODIMRA_ROOTCA_CRT_PATH} -CAkey ${ODIMRA_ROOTCA_KEY_PATH} -CAcreateserial -out ${NGINX_SERVER_CRT_PATH} -extensions v3_req -extfile <( cat <<EOF
[req]
default_bits = ${KEY_LENGTH}
encrypt_key  = no
default_md   = sha512
prompt       = no
utf8         = yes
distinguished_name = req_distinguished_name
req_extensions = v3_req

[req_distinguished_name]
C  = ${CERT_COUNTRY_NAME}
ST = ${CERT_STATE_NAME}
L  = ${CERT_LOCALITY_NAME}
O  = ${CERT_ORGANIZATION_NAME}
OU = ${CERT_ORGANIZATION_UNIT_NAME}
CN = ODIMRA_PROXY_CRT

[v3_req]
subjectKeyIdentifier    = hash
authorityKeyIdentifier  = keyid:always,issuer:always
keyUsage                = critical, nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage        = clientAuth, serverAuth
subjectAltName          = @alt_names

[alt_names]
DNS.0 = odimra.proxy.net
IP.0 = ${ODIMRA_VIRTUAL_IP}
EOF
)
	eval_cmd_exec $? "${NGINX_SERVER_CRT_PATH} generation failed"

	# remove temp files
	rm -f ${NGINX_SERVER_CSR_PATH} ${ODIMRA_CERT_DIR}/rootCA.srl

	chmod 0600 ${NGINX_SERVER_KEY_PATH} ${NGINX_SERVER_CRT_PATH}
}

# read_config_value parses the config file
# and fetches the value of the parameter passed
read_config_value()
{
	param_name=$1
	echo "$(grep -w "${param_name}" ${ODIM_CONTROLLER_CONFIG} | cut -d":" -f2 | awk '{$1=$1};1' | sed -e 's/^"//' -e 's/"$//')"
}

parse_cert_san()
{
	ODIM_SERVER_FQDN_SAN="DNS.0 = ${ODIMRA_FQDN}"
	ODIM_SERVER_FQDN_COUNT=1
	ODIM_SERVER_IP_SAN=""
	ODIM_SERVER_IP_COUNT=0
	ODIM_KAFKA_CLIENT_FQDN_SAN="DNS.0 = ${ODIMRA_FQDN}"
	ODIM_KAFKA_CLIENT_FQDN_COUNT=1
	ODIM_KAFKA_CLIENT_IP_SAN=""
	ODIM_KAFKA_CLIENT_IP_COUNT=0

	odimra_server_fqdn_san=$(read_config_value "odimraServerCertFQDNSan")
	odimra_server_ip_san=$(read_config_value "odimraServerCertIPSan")
	odimra_kafka_client_fqdn_san=$(read_config_value "odimraKafkaClientCertFQDNSan")
	odimra_kafka_client_ip_san=$(read_config_value "odimraKafkaClientCertIPSan")

	OLDIFS=$OIFS
	IFS=','
	if [[ -n $odimra_server_fqdn_san ]]; then
		tempArray=($odimra_server_fqdn_san)
		count=${ODIM_SERVER_FQDN_COUNT}
		for fqdn in "${tempArray[@]}"; do
			ODIM_SERVER_FQDN_SAN+="\nDNS.$count = $fqdn"
			((count++))
		done
	fi
	if [[ -n $odimra_server_ip_san ]]; then
		tempArray=($odimra_server_ip_san)
		count=${ODIM_SERVER_IP_COUNT}
		for ip in "${tempArray[@]}"; do
			ODIM_SERVER_IP_SAN+="\nIP.$count = $ip"
			((count++))
		done
	fi
	if [[ -n $odimra_kafka_client_fqdn_san ]]; then
		tempArray=($odimra_kafka_client_fqdn_san)
		count=${ODIM_KAFKA_CLIENT_FQDN_COUNT}
		for fqdn in "${tempArray[@]}"; do
			ODIM_KAFKA_CLIENT_FQDN_SAN+="\nDNS.$count = $fqdn"
			((count++))
		done
	fi
	if [[ -n $odimra_kafka_client_ip_san ]]; then
		tempArray=($odimra_kafka_client_ip_san)
		count=${ODIM_KAFKA_CLIENT_IP_COUNT}
		for ip in "${tempArray[@]}"; do
			ODIM_KAFKA_CLIENT_IP_SAN+="\nIP.$count = $ip"
			((count++))
		done
	fi
	IFS=$OLDIFS

	ODIM_SERVER_FQDN_SAN=$(echo -e $ODIM_SERVER_FQDN_SAN)
	ODIM_SERVER_IP_SAN=$(echo -e $ODIM_SERVER_IP_SAN)
	ODIM_KAFKA_CLIENT_FQDN_SAN=$(echo -e $ODIM_KAFKA_CLIENT_FQDN_SAN)
	ODIM_KAFKA_CLIENT_IP_SAN=$(echo -e $ODIM_KAFKA_CLIENT_IP_SAN)
}

# read_config_file parses the config file
# and assigns the values to the pre-defined
# global variables
read_config_file()
{
	# parse config value and assign values
	# to global variables
	ODIMRA_FQDN=$(read_config_value "fqdn")
	KAFKA_JKS_PASSWORD=$(read_config_value "kafkaJKSPassword")
	ZOOKEEPER_JKS_PASSWORD=$(read_config_value "zookeeperJKSPassword")
	ODIMRA_NAMESPACE=$(read_config_value "namespace")
	ODIMRA_HA_DEPLOYMENT=$(read_config_value "haDeploymentEnabled")
	ODIMRA_VIRTUAL_IP=$(read_config_value "virtualIP")

	parse_cert_san

	# use global variable and derive filepaths
	ODIMRA_ROOTCA_CRT_PATH=${ODIMRA_CERT_DIR}/rootCA.crt
	ODIMRA_ROOTCA_KEY_PATH=${ODIMRA_CERT_DIR}/rootCA.key
	ODIMRA_SERVER_CSR_PATH=${ODIMRA_CERT_DIR}/odimra_server.csr
	ODIMRA_SERVER_CRT_PATH=${ODIMRA_CERT_DIR}/odimra_server.crt
	ODIMRA_SERVER_KEY_PATH=${ODIMRA_CERT_DIR}/odimra_server.key
	ODIMRA_RSA_PUBLIC_KEY=${ODIMRA_CERT_DIR}/odimra_rsa.public
	ODIMRA_RSA_PRIVATE_KEY=${ODIMRA_CERT_DIR}/odimra_rsa.private
	ODIMRA_KAFKA_CLIENT_CSR_PATH=${ODIMRA_CERT_DIR}/odimra_kafka_client.csr
	ODIMRA_KAFKA_CLIENT_CRT_PATH=${ODIMRA_CERT_DIR}/odimra_kafka_client.crt
	ODIMRA_KAFKA_CLIENT_KEY_PATH=${ODIMRA_CERT_DIR}/odimra_kafka_client.key
	KAFKA_JKS_PATH=${ODIMRA_CERT_DIR}/kafka.keystore.jks
	KAFKA_JTS_PATH=${ODIMRA_CERT_DIR}/kafka.truststore.jks
	ZOOKEEPER_JKS_PATH=${ODIMRA_CERT_DIR}/zookeeper.keystore.jks
	ZOOKEEPER_JTS_PATH=${ODIMRA_CERT_DIR}/zookeeper.truststore.jks
	ODIMRA_ETCD_SERVER_CSR_PATH=${ODIMRA_CERT_DIR}/odimra_etcd_server.csr
	ODIMRA_ETCD_SERVER_CRT_PATH=${ODIMRA_CERT_DIR}/odimra_etcd_server.crt
	ODIMRA_ETCD_SERVER_KEY_PATH=${ODIMRA_CERT_DIR}/odimra_etcd_server.key
	NGINX_SERVER_CSR_PATH=${ODIMRA_CERT_DIR}/nginx_server.csr
	NGINX_SERVER_CRT_PATH=${ODIMRA_CERT_DIR}/nginx_server.crt
	NGINX_SERVER_KEY_PATH=${ODIMRA_CERT_DIR}/nginx_server.key
}

# validate_config_params is for validating
# all the mandatory config params in the
# passed config file
validate_config_params()
{
	count=0
	if [[ -z ${ODIMRA_FQDN} ]]; then
		echo "[$(date)] -- ERROR -- mandatory param fqdn cannot be empty"
		((count++))
	fi
	if [[ -z ${KAFKA_JKS_PASSWORD} ]]; then
		echo "[$(date)] -- ERROR -- mandatory param kafkaJKSPassword cannot be empty"
		((count++))
	fi
	if [[ -z ${ZOOKEEPER_JKS_PASSWORD} ]]; then
		echo "[$(date)] -- ERROR -- mandatory param zookeeperJKSPassword cannot be empty"
		((count++))
	fi
	if [[ -z ${ODIMRA_NAMESPACE} ]]; then
		echo "[$(date)] -- ERROR -- mandatory param namespace cannot be empty"
		((count++))
	fi
	if [[ -z ${ODIMRA_HA_DEPLOYMENT} ]]; then
		echo "[$(date)] -- INFO  -- haDeploymentEnabled param not found or value not assigned, default value considered"
		ODIMRA_HA_DEPLOYMENT=false
	fi
	if [[ ${ODIMRA_HA_DEPLOYMENT,,} == true ]] && [[ -z ${ODIMRA_VIRTUAL_IP} ]]; then
		echo "[$(date)] -- ERROR -- mandatory param virtualIP cannot be empty"
		((count++))
	fi

	if [[ $count -ne 0 ]]; then
		echo "[$(date)] -- ERROR -- $count parameter(s) have invalid value configured, exiting"
		exit 1
	fi
}

generate_etcd_certs()
{
	# check if cert and key exists and is not empty
	if [[ -s ${ODIMRA_ETCD_SERVER_KEY_PATH} ]] && [[ -s ${ODIMRA_ETCD_SERVER_CRT_PATH} ]]; then
		echo "[$(date)] -- INFO  -- odimra etcd server crt and key already exists"
		# verify crt was signed by the rootCA present
		${OPENSSL_BIN_PATH} verify -CAfile ${ODIMRA_ROOTCA_CRT_PATH} ${ODIMRA_ETCD_SERVER_CRT_PATH} > /dev/null
		eval_cmd_exec $? "${ODIMRA_ETCD_SERVER_CRT_PATH} is not signed by ${ODIMRA_ROOTCA_CRT_PATH}"
		return
	fi

	# check if rootCA key was made available, if it was not generated by this script.
	if ${ROOT_CA_KEY_NOT_PROVIDED}; then
		echo "[$(date)] -- ERROR -- rootCA key was not provided, odim etcd server crt generation not possible"
		exit 1
	fi

	#generate etcd private key
	${OPENSSL_BIN_PATH} genrsa -out ${ODIMRA_ETCD_SERVER_KEY_PATH} ${KEY_LENGTH}
	eval_cmd_exec $? "${ODIMRA_ETCD_SERVER_KEY_PATH} generation failed"

	#generate etcd csr
	${OPENSSL_BIN_PATH} req -new -sha512 -key ${ODIMRA_ETCD_SERVER_KEY_PATH} -subj "/C=${CERT_COUNTRY_NAME}/ST=${CERT_STATE_NAME}/L=${CERT_LOCALITY_NAME}/O=${CERT_ORGANIZATION_NAME}/OU=${CERT_ORGANIZATION_UNIT_NAME}/CN=etcd" -config <(cat <<EOF
[ req ]
prompt=no
distinguished_name=subject
req_extensions=req_ext
[ subject ]
commonName=etcd
[ req_ext ]
extendedKeyUsage=serverAuth, clientAuth
basicConstraints=critical,CA:false
keyUsage=nonRepudiation, digitalSignature, keyEncipherment
EOF
) -out ${ODIMRA_ETCD_SERVER_CSR_PATH}
	eval_cmd_exec $? "${ODIMRA_ETCD_SERVER_CSR_PATH} generation failed"

	#sign and gen etcd certificate
	${OPENSSL_BIN_PATH} x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt]
basicConstraints=CA:FALSE
keyUsage=nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName=@alternate_names
[ alternate_names ]
DNS.0 = etcd
DNS.1 = etcd1
DNS.2 = etcd2
DNS.3 = etcd3
DNS.4 = etcd1.etcd.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.5 = etcd2.etcd.${ODIMRA_NAMESPACE}.svc.cluster.local
DNS.6 = etcd3.etcd.${ODIMRA_NAMESPACE}.svc.cluster.local
EOF
) -in  ${ODIMRA_ETCD_SERVER_CSR_PATH} -CA  ${ODIMRA_ROOTCA_CRT_PATH} -CAkey  ${ODIMRA_ROOTCA_KEY_PATH}  -CAcreateserial -out ${ODIMRA_ETCD_SERVER_CRT_PATH} -days ${CERT_VALIDITY_PERIOD} -sha512
	eval_cmd_exec $? "${ODIMRA_ETCD_SERVER_CRT_PATH} generation failed"

	# clean up temp files generated
	rm -f ${ODIMRA_ETCD_SERVER_CSR_PATH} ${ODIMRA_CERT_DIR}/rootCA.srl
}

# generate_certs is for generating the
# certificates and private keys required
# by ODIM-RA, kafka and zookeeper services
generate_certs()
{
	generate_ca_certs
	generate_odim_server_certs
	generate_odim_kafka_client_certs
	generate_kafka_certs
	generate_zookeeper_certs
	generate_etcd_certs
	generate_nginx_certs

	#create a temp file, to indicate certs were generated by this script
	touch ${ODIMRA_CERT_DIR}/.gen_odimra_certs.ok
}

usage()
{
        echo -e "$(basename $BASH_SOURCE) <dir_path_to_store_gen_certs> <odim-controller_config_file_path>"
        exit 1
}

##############################################
###############  MAIN  #######################
##############################################

if [[ $# -ne 2 ]]; then
        usage
fi

ODIMRA_CERT_DIR=$1
ODIM_CONTROLLER_CONFIG=$2

if [[ -z ${ODIMRA_CERT_DIR} ]] || [[ ! -d ${ODIMRA_CERT_DIR} ]]; then
	echo "[$(date)] -- ERROR -- invalid directory path [${ODIMRA_CERT_DIR}] passed, exiting!!!"
	exit 1
fi

if [[ -z ${ODIM_CONTROLLER_CONFIG} ]] || [[ ! -f ${ODIM_CONTROLLER_CONFIG} ]]; then
	echo "[$(date)] -- ERROR -- invalid odim-controller config file [${ODIM_CONTROLLER_CONFIG}] passed, exiting!!!"
	exit 1
fi

pre_reqs

read_config_file

validate_config_params

generate_certs

exit 0
