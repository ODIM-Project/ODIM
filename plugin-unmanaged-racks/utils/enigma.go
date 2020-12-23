/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
)

// Enigma offers encryption/decryption API which utilizes provided private/public key pair.
type Enigma struct {
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
}

// Decrypt decrypts provided toBeDecrypted string.
func (e *Enigma) Decrypt(toBeDecrypted string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(toBeDecrypted)
	if err != nil {
		logging.Fatal(err)
	}

	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, e.priv, decoded, nil)
	if err != nil {
		logging.Fatal(err)
	}
	return plaintext
}

// Encrypt encrypts provided toBeEncrypted string.
func (e *Enigma) Encrypt(toBeEncrypted []byte) string {
	hash := sha512.New()
	encrypted, err := rsa.EncryptOAEP(hash, rand.Reader, e.pub, toBeEncrypted, nil)
	if err != nil {
		logging.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(encrypted)
}

// NewEnigma constructs Enigma by loading private/public key pair from provided paths
func NewEnigma(privKeyPath, pubKeyPath string) *Enigma {
	privateKeyBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		logging.Fatalf("Cannot load PrivateKey from given path: '%s' because of  %s", privKeyPath, err)
	}

	publicKeyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		logging.Fatalf("Cannot load PublicKey from given path: '%s' because of %s", privKeyPath, err)
	}

	return &Enigma{
		priv: bytesToPrivateKey(privateKeyBytes),
		pub:  bytesToPublicKey(publicKeyBytes),
	}
}

// CreateEnigma constructs Enigma using provided private/public key pair
func CreateEnigma(privateKeyBytes, publicKeyBytes []byte) *Enigma {
	return &Enigma{
		priv: bytesToPrivateKey(privateKeyBytes),
		pub:  bytesToPublicKey(publicKeyBytes),
	}
}

func bytesToPublicKey(publicKey []byte) *rsa.PublicKey {
	block, _ := pem.Decode(publicKey)
	var b = block.Bytes
	if x509.IsEncryptedPEMBlock(block) {
		logging.Info("is encrypted pem block")
		dpb, err := x509.DecryptPEMBlock(block, nil)
		if err != nil {
			logging.Fatal(err)
		}
		b = dpb
	}
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		logging.Fatal(err)
	}
	key, ok := pub.(*rsa.PublicKey)
	if !ok {
		logging.Fatal("not ok")
	}
	return key
}

func bytesToPrivateKey(privateKey []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(privateKey)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		logging.Info("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			logging.Error(err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		logging.Fatal(err)
	}
	return key
}
