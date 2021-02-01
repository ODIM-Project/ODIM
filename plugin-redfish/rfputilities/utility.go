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

//Package rfputilities ...
package rfputilities

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpresponse"
	"github.com/fsnotify/fsnotify"
)

// GetPlainText ...
func GetPlainText(password []byte) ([]byte, error) {
	priv := []byte(rfpmodel.PluginPrivateKey)
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Info("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Error(err.Error())
			return []byte{}, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Info(err.Error())
		return []byte{}, err
	}

	hash := sha512.New()

	return rsa.DecryptOAEP(
		hash,
		rand.Reader,
		key,
		password,
		nil,
	)
}

//Status holds the Status of plugin it will be intizaied during startup time
var Status rfpresponse.Status

// TrackConfigFileChanges monitors the odim config changes using fsnotfiy
func TrackConfigFileChanges(configFilePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err)
	}
	err = watcher.Add(configFilePath)
	if err != nil {
		log.Error(err)
	}
	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					continue
				}
				log.Info("event:", fileEvent)
				if fileEvent.Op&fsnotify.Write == fsnotify.Write || fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
					log.Info("modified file:", fileEvent.Name)
					// update the odim config
					if err := config.SetConfiguration(); err != nil {
						log.Error("error while trying to set configuration" + err.Error())
					}
				}
				//Reading file to continue the watch
				watcher.Add(configFilePath)
			case err, _ := <-watcher.Errors:
				if err != nil {
					log.Error(err)
					defer watcher.Close()
				}
			}
		}
	}()
}
