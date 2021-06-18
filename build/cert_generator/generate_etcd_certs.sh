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
# for etcd server

CommonName=etcd
SAN_DNS_LIST="DNS.0 = etcd"

if [[ ! -e ./rootCA.crt ]] || [[ ! -e ./rootCA.key ]]; then
	echo "[$(date)] -- ERROR -- root CA files not present in current directory"
	exit 1
fi

#generate odimra_etcd_server private key
openssl genrsa -out odimra_etcd_server.key 4096
if [[ $? -ne 0 ]]; then
	echo "[$(date)] -- ERROR -- private key generation failed"
	exit 1
fi

#generate odimra_etcd_server csr
openssl req -new -sha512 -key odimra_etcd_server.key -subj "/C=US/ST=CA/L=California/O=ODIMRA/CN=${CommonName}" -config <(cat <<EOF
[ req ]
prompt=no
distinguished_name=subject
req_extensions=req_ext
[ subject ]
commonName=${CommonName}
[ req_ext ]
extendedKeyUsage=serverAuth, clientAuth
basicConstraints=critical,CA:false
keyUsage=nonRepudiation, digitalSignature, keyEncipherment
EOF
) -out odimra_etcd_server.csr
if [[ $? -ne 0 ]]; then
	echo "[$(date)] -- ERROR -- csr generation failed"
	exit 1
fi

#sign and gen odimra_etcd_server certificate
openssl x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt]
basicConstraints=CA:FALSE
keyUsage=nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=serverAuth, clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName=@alternate_names
[ alternate_names ]
${SAN_DNS_LIST}
EOF
) -in  odimra_etcd_server.csr -CA  rootCA.crt -CAkey  rootCA.key  -CAcreateserial -out odimra_etcd_server.crt -days 500 -sha512
if [[ $? -ne 0 ]]; then
	echo "[$(date)] -- ERROR -- certificate generation failed"
	exit 1
fi

# clean up temp files generated
rm -f odimra_etcd_server.csr rootCA.srl
