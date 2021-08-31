#!/bin/bash
# (C) Copyright [2020] Hewlett Packard Enterprise Development LP
# 
# Licensed under the Apache License, Version 2.0 (the "License"); you may
# not use this file except in compliance with the License. You may obtain
# a copy of the License at
# 
#     http://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
# WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
# License for the specific language governing permissions and limitations
# under the License.

localhostFQDN=$1
if [ -z "$1" ]
then
    echo "Please provide localhostFQDN"
    exit -1
else
    echo "localhostFQDN is $1"
fi
 

openssl genrsa -out rootCA.key 4096
echo
echo
#generate root CA certificate
openssl req -new -key rootCA.key -sha512 -days 1024 -x509 -subj "/C=US/ST=CA/L=California/O=ODIMRA/CN=Root CA" -config <( cat <<EOF
[ req ]
prompt = no
distinguished_name = req_distinguished_name
x509_extensions = v3_ca
[ req_distinguished_name  ]
C=US
ST=CA
L=California
O=ACME Corp
CN=Root CA
[ v3_ca ]
basicConstraints        = critical, CA:true
keyUsage                = critical, keyCertSign
issuerAltName           = issuer:copy
subjectKeyIdentifier    = hash
authorityKeyIdentifier  = keyid:always, issuer:always
subjectAltName          = email:admin@telco.net
EOF
) -out rootCA.crt
echo
echo
 
#generate server private key
openssl genrsa -out odimra_server.key 4096
echo
echo
#generate server csr
openssl req -new -sha512 -key odimra_server.key -subj "/C=US/ST=CA/O=ODIMRA/CN=Server Cert" -config <(cat <<EOF
[ req ]
prompt = no
distinguished_name = subject
req_extensions    = req_ext
[ subject ]
commonName = Server Cert
[ req_ext ]
extendedKeyUsage=serverAuth, clientAuth
basicConstraints=critical,CA:false
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
EOF
) -out odimra_server.csr
echo
echo
#sign and gen server certificate
openssl x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt]
basicConstraints=CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName = @alternate_names
[ alternate_names ]
DNS.1 = $localhostFQDN
DNS.2 = URP
EOF
) -in  odimra_server.csr -CA  rootCA.crt -CAkey  rootCA.key  -CAcreateserial -out odimra_server.crt -days 500 -sha512
echo
echo
openssl genrsa -out odimra_rsa.private 4096
echo
echo
openssl rsa -in odimra_rsa.private -out odimra_rsa.public -pubout -outform PEM
echo
echo

# generate client certificates for odimra to connect with Kafka as client
/bin/bash generate_client_crt.sh ./rootCA.crt ./rootCA.key ${localhostFQDN} "odimra Kafka Client"
mv client.key odimra_kafka_client.key
mv client.crt odimra_kafka_client.crt

# generate etcd server certificates
/bin/bash generate_etcd_certs.sh

# cleanup temp files generated
rm -f odimra_server.csr rootCA.srl client.csr
