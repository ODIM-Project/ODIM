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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"sync"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

var (
	// MuxLock is used for avoiding race conditions
	MuxLock = &sync.Mutex{}
)

// DecryptWithPrivateKey is used to decrypt ciphered text to device password
// with the private key whose path is available in the config file
func DecryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
	MuxLock.Lock()
	defer MuxLock.Unlock()
	var err error
	block, _ := pem.Decode(config.Data.KeyCertConf.RSAPrivateKey)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, fmt.Errorf("error while trying to decrypt pem block: %v", err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, fmt.Errorf("error while parsing private key: %v", err)
	}
	hash := sha512.New()
	plainText, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		key,
		ciphertext,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error while trying to decrypt password: %v", err)
	}
	return plainText, nil
}

// EncryptWithPublicKey is used to encrypt device password using odimra public key
func EncryptWithPublicKey(password []byte) ([]byte, error) {
	var err error
	block, _ := pem.Decode(config.Data.KeyCertConf.RSAPublicKey)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, fmt.Errorf("error while trying to decrypt pem block: %v", err)
		}
	}

	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, fmt.Errorf("error while parsing private key: %v", err)
	}

	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, err
	}

	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, key, password, nil)
	if err != nil {
		return nil, fmt.Errorf("error while trying to encrypt password: %v", err)
	}
	return ciphertext, nil
}
