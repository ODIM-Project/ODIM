//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.
package common

import (
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

var (
	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA8dbQu0GbdVU5TNykAmKp
014wyjmUEJ7oiKOJKWBjfUb6UoH7/iutDr5wu0E0H6Rgup2rlDt+quEI/MgZ7Fdm
4Jzp7n0Xc2Xs8Dc1Au7n0z+k70huGZqJcrB4giBnp5gIb1e1/gbPvCHiOYzhCVAS
cIirp6KKyRt2nREKd8WUECzFahKOnw6/yEEfPVjZJAtrxm4cGlMTRIEp3Nq3V0l+
4tllA4xOHXPEBuQm9fYgO+8WZvBXdEFFovgjOKEShYu6czrt+/1Lld7vgy0X2GUc
odP+YB3X40S8vnswtzBkvScQV1yg7u+rmK80c98LSLteAPpIukzqHYHBDxo35g6e
5QIDAQAB
-----END PUBLIC KEY-----`
	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA8dbQu0GbdVU5TNykAmKp014wyjmUEJ7oiKOJKWBjfUb6UoH7
/iutDr5wu0E0H6Rgup2rlDt+quEI/MgZ7Fdm4Jzp7n0Xc2Xs8Dc1Au7n0z+k70hu
GZqJcrB4giBnp5gIb1e1/gbPvCHiOYzhCVAScIirp6KKyRt2nREKd8WUECzFahKO
nw6/yEEfPVjZJAtrxm4cGlMTRIEp3Nq3V0l+4tllA4xOHXPEBuQm9fYgO+8WZvBX
dEFFovgjOKEShYu6czrt+/1Lld7vgy0X2GUcodP+YB3X40S8vnswtzBkvScQV1yg
7u+rmK80c98LSLteAPpIukzqHYHBDxo35g6e5QIDAQABAoIBAQDnsAd0/puywx0M
N+2go2lTqE9Rzeu+KJ9aGGJVk5R89rzmwsTqcmlvUJ+rpgILtm09G8S/VGg7yS/V
DNdZBzr2QR4Ubx9CXQmr8RgGYV8TkUuwOlHQka7Qg6RP9j+X3h1mnj8qyNfHwyZ6
QC9vvpiL20OobB5OINN4ElVW/aCmBIGxg/ooWXB1Ck3HqzEFq/Zg4Uy1HE6b8BV5
uQ+Jc4B1gBppvYpf/Ie0QKmqccAzd4RM5efEohrJO5eZE8FjL5xwEt7vXfARbo6c
gF2UMkZuusxIsirPdWvw/Xw/zqFpur25cEYqdHEsiH1e6pMnevk3kUgPkpMq+yBP
VILtLl5BAoGBAPoyM03v5cexLeq3c4WE7XM/CApdOr7hhDkBYwXmUXP8kp+MziMo
hS4RihVZUf17D4KxWlBl9wmZOGyUb4Nq+/bwluUELnnaYBzWwU4OkJE6fO2qJcEs
S6sla5FlJ61VH0goxaVDGO2X9ToouXAxXW4uhK0SJI3WNltW6udsyZXtAoGBAPdy
/JpLcjo9LcL6MabUE+ZKpcqrPYQhoBw9slpQhn0s2ER/ho9/b4M5UdMiY3VUuAtF
iBzqlMVM6MkC2xW3qmgW5U8lcBzyFI8fu4IS5ZXUVKGMWIMAl6vwnsXUamu+6OcD
nVZS6u0yu3yZg32L8AiFO4YTepba+5k25akCp43ZAoGBAPVNSM3uGmzKg4lwahwL
sz9eGkUHGTTTKO83M94x7cR5a0xxIh6IeOMtISRDWcbb494wgqr2/dl0V0Tl19uS
hg2b32YUznh8KeW8jPQ6BXXOUXQ3cSLPijT30FpSQi+ImM4H45hfi85PQYjPKtkc
HU2M4FpLwnkqAEtXkaJrH84VAoGAFpMz/nOhoTSRpzciLoEsq5bl1z6WJybWL51l
Vx3/lw3vURh9UzwiFUu2ble10S+Adu7KAzFXj0R7/FK5YBrYfhSQqQ7WUp23SHNx
rOVCcs/jRLXEIXd9Xt9d7Nh7OQc6wlCvGwAHlMpLFov+1gZdSLm2+31tcrPZvlmm
zCuE08kCgYEAxgxjPCIoysB2yRHuScQjtKTfz5LdRmZ2uHOhxJMWxhvYmzqVBSm8
TkuKsMlbg2EbYyhFAYb0GucRV/jTz6tWNLreW97dn5NkO4IU/I0FbnKWB0zgAqCp
BLbLNzuwoImcammGEA3pXDNVzvYilBct3Mjr1no2wCQhUQbqp33vBls=
-----END RSA PRIVATE KEY-----`
)

func TestEncryptWithPublicKey(t *testing.T) {
	var err error

	config.Data.KeyCertConf = &config.KeyCertConf{
		RSAPublicKey:  []byte(publicKey),
		RSAPrivateKey: []byte(privateKey),
	}

	encryptedData, err := EncryptWithPublicKey([]byte("testData"))
	if err != nil {
		t.Errorf("EncryptWithPublicKey failed with %v", err)
	}
	decryptedData, err := DecryptWithPrivateKey(encryptedData)
	if err != nil {
		t.Errorf("DecryptWithPrivateKey failed with %v", err)
	}
	if string(decryptedData) != "testData" {
		t.Errorf("Mismatch in data encrypted and the decrypted data, want = testData got = %v", string(decryptedData))
	}

}
