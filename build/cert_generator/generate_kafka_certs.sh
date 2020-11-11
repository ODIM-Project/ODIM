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
keytool -keystore kafka.keystore.jks -alias $1 -validity 365 -keyalg RSA -genkey <<HERE
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

# add kafka rootCA to kafka server truststore
keytool -keystore kafka.truststore.jks -alias CARoot -import -file rootCA.crt <<HERE
K@fk@_store1
K@fk@_store1
yes
HERE
echo
echo

# generate keystore CSR
keytool -keystore kafka.keystore.jks -alias $1 -certreq -file cert-file <<HERE
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
DNS.1 = kafka
DNS.2 = kafka1
DNS.3 = kafka2
DNS.4 = kafka3
EOF
) -in cert-file -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out cert-signed -days 500 -sha512
echo
echo

# add kafka rootCA to kafka server keystore
keytool -keystore kafka.keystore.jks -alias CARoot -import -file rootCA.crt <<HERE
K@fk@_store1
yes
HERE
echo
echo

# adding server certificate to server key store
keytool -keystore kafka.keystore.jks -alias $1 -import -file cert-signed <<HERE
K@fk@_store1
HERE
echo
echo

# clean up temp files generated
rm -f cert-file cert-signed kafka.p12 rootCA.srl kafka.crt kafka.key
