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

# Script is for generating certificate and private key
# for Client mode connection usage only

if [ -z "$1" ]
then
    echo "Please provide localhostFQDN"
    exit -1
else
    echo "localhostFQDN is $1"
fi

echo

# generate keystore
keytool -keystore zookeeper.keystore.jks -alias $1 -validity 365 -keyalg RSA -genkey <<HERE
K@fk@_store1
K@fk@_store1
$1
Telco Solutions
HPE
California
CA
US
yes
.
HERE
echo
echo

# add zookeeper rootCA to zookeeper server truststore
keytool -keystore zookeeper.truststore.jks -alias CARoot -import -file rootCA.crt <<HERE
K@fk@_store1
K@fk@_store1
yes
HERE
echo
echo

# generate keystore CSR
keytool -keystore zookeeper.keystore.jks -alias $1 -certreq -file cert-file <<HERE
K@fk@_store1
HERE
echo
echo

#generate the keystore certificates
openssl x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt] basicConstraints=CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName = @alternate_names
[ alternate_names ]
DNS.1 = zookeeper
DNS.2 = zookeeper1
DNS.3 = zookeeper2
DNS.4 = zookeeper3
DNS.5 = zookeeper1.odim.svc.cluster.local
DNS.6 = zookeeper2.odim.svc.cluster.local
DNS.7 = zookeeper3.odim.svc.cluster.local
EOF
) -in cert-file -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out cert-signed -days 500 -sha512
echo
echo

# add zookeeper rootCA to zookeeper server keystore
keytool -keystore zookeeper.keystore.jks -alias CARoot -import -file rootCA.crt <<HERE
K@fk@_store1
yes
HERE
echo
echo

# adding server certificate to server key store
keytool -keystore zookeeper.keystore.jks -alias $1 -import -file cert-signed <<HERE
K@fk@_store1
HERE
echo
echo

# generating zookeeper certs
keytool -importkeystore -srckeystore zookeeper.keystore.jks -destkeystore zookeeper.p12 -srcstoretype JKS -deststoretype PKCS12 <<HERE
K@fk@_store1
K@fk@_store1
K@fk@_store1
HERE
echo
echo

# clean up temp files generated
rm -f cert-file cert-signed zookeeper.p12 rootCA.srl zookeeper.key zookeeper.crt
