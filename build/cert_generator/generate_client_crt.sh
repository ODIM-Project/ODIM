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

CACertPath=$1
CAKeyPath=$2
SanDNS=$3
CommonName=$4
FileName=client

usage()
{
	echo -e "${BC}$(basename $BASH_SOURCE) <CACertPath> <CAKeyPath> <FQDN> <CommonName>"
	exit 1
}

gen_crt() {

	#generate server private key
	openssl genrsa -out ${FileName}.key 4096

	echo
	echo

	#generate server csr
	openssl req -new -sha512 -key ${FileName}.key -subj "/C=US/ST=CA/O=ODIMRA/CN=${CommonName}" -config <(cat <<EOF
[ req ]
prompt=no
distinguished_name=subject
req_extensions=req_ext
[ subject ]
commonName=Server Cert
[ req_ext ]
extendedKeyUsage=clientAuth
basicConstraints=critical,CA:false
keyUsage=nonRepudiation, digitalSignature, keyEncipherment
EOF
) -out ${FileName}.csr

	echo
	echo

	#sign and gen server certificate
	openssl x509 -req -extensions server_crt -extfile <( cat <<EOF
[server_crt]
basicConstraints=CA:FALSE
keyUsage=nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage=clientAuth
subjectKeyIdentifier=hash
authorityKeyIdentifier=keyid,issuer
subjectAltName=@alternate_names
[ alternate_names ]
DNS.0=${SanDNS}
EOF
) -in  ${FileName}.csr -CA  ${CACertPath} -CAkey  ${CAKeyPath}  -CAcreateserial -out ${FileName}.crt -days 500 -sha512

	echo "Generated files : [${FileName}.key ${FileName}.csr ${FileName}.crt]"
}

##############################################
###############  MAIN  #######################
##############################################
if [[ $# -ne 4 ]]; then
	usage
fi

if [[ -z ${CACertPath} ]]; then
	echo "-- ERROR -- CA certificate path cannot be empty"
	usage
fi

if [[ -z ${CAKeyPath} ]]; then
	echo "-- ERROR -- CA key path cannot be empty"
	usage
fi

if [[ -z ${SanDNS} ]]; then
	echo "-- ERROR -- FQDN value cannot be empty"
	usage
fi

if [[ -z ${CommonName} ]]; then
	echo "-- ERROR -- Common Name value cannot be empty"
	usage
fi

gen_crt
