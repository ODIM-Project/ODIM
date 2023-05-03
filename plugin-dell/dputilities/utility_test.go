//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

//Package //(C) Copyright [2020] Hewlett Packard Enterprise Development LP
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

// Package dputilities ...
package dputilities

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/stretchr/testify/assert"
)

func TestTrackConfigFileChanges(t *testing.T) {
	config.SetUpMockConfig(t)
	// Create a temporary config file for testing
	configFile, err := ioutil.TempFile("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(configFile.Name())

	// Start watching the config file
	errChan := make(chan error)

	go TrackIPConfigListener(configFile.Name(), errChan)

	// Write to the config file and wait for the changes to be detected
	configFile.Write([]byte("test config"))
	time.Sleep(1 * time.Second)

	// Check that the configuration was updated
	err = config.SetConfiguration()
	assert.NotNil(t, err)

	// Invalid Path
	go TrackIPConfigListener("", errChan)

}

func TestGetPlainText(t *testing.T) {
	config.SetUpMockConfig(t)

	password := []byte("password123")

	// Generate a new RSA key pair for testing
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error generating RSA key pair: %v", err)
	}

	// Encrypt the password using RSA-OAEP
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, key.Public().(*rsa.PublicKey), password, nil)
	if err != nil {
		t.Fatalf("error encrypting password: %v", err)
	}

	// Convert the private key to PEM format
	priv := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	dpmodel.PluginPrivateKey = priv
	// Decrypt the password using the GetPlainText function
	ctxt := mockContext()
	plaintext, err := GetPlainText(ctxt, ciphertext)
	if err != nil {
		t.Fatalf("error decrypting password: %v", err)
	}

	// Check that the decrypted plaintext matches the original password
	if !bytes.Equal(plaintext, password) {
		t.Fatalf("plaintext does not match original password")
	}
}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}
