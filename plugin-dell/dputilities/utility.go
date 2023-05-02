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

// Package dputilities ...
package dputilities

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/plugin-dell/config"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpmodel"
	"github.com/ODIM-Project/ODIM/plugin-dell/dpresponse"
	log "github.com/sirupsen/logrus"
)

var (
	confMutex = &sync.RWMutex{}
	podName   = os.Getenv("POD_NAME")
)

// GetPlainText ...
func GetPlainText(ctx context.Context, password []byte) ([]byte, error) {
	priv := []byte(dpmodel.PluginPrivateKey)
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			l.LogWithFields(ctx).Error(err.Error())
			return []byte{}, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		l.LogWithFields(ctx).Error(err.Error())
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

// Status holds the Status of plugin it will be intizaied during startup time
var Status dpresponse.Status

// PluginStartTime hold the time from which plugin started
var PluginStartTime time.Time

// TrackIPConfigListener listes to the chanel on the config file changes of Plugin
func TrackIPConfigListener(configFilePath string, errChan chan error) {
	eventChan := make(chan interface{})
	format := config.Data.LogFormat
	transactionID := uuid.New()
	ctx := CreateContext(transactionID.String(), common.PluginTrackFileConfigActionID, common.PluginTrackFileConfigActionName, "1", common.PluginTrackFileConfigActionName)
	go TrackConfigFileChanges(configFilePath, eventChan, errChan)
	for {
		select {
		case info := <-eventChan:
			l.LogWithFields(ctx).Info(info) // new data arrives through eventChan channel
			if l.Log.Level != config.Data.LogLevel {
				l.LogWithFields(ctx).Info("Log level is updated, new log level is ", config.Data.LogLevel)
				l.Log.Logger.SetLevel(config.Data.LogLevel)
			}
			if format != config.Data.LogFormat {
				l.SetFormatter(config.Data.LogFormat)
				format = config.Data.LogFormat
				l.LogWithFields(ctx).Info("Log format is updated, new log format is ", config.Data.LogFormat)
			}
		case err := <-errChan:
			l.LogWithFields(ctx).Error(err)
		}
	}
}

// TrackConfigFileChanges monitors the config changes using fsnotfiy
func TrackConfigFileChanges(configFilePath string, eventChan chan<- interface{}, errChan chan<- error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		errChan <- err
	}
	err = watcher.Add(configFilePath)
	if err != nil {
		errChan <- err
	}

	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if fileEvent.Op&fsnotify.Write == fsnotify.Write || fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
					log.Debug("Modified file: " + fileEvent.Name)
					confMutex.Lock()

					// update the plugin config
					if err := config.SetConfiguration(); err != nil {
						log.Error("While trying to set configuration, got: " + err.Error())
					}
					confMutex.Unlock()
					eventChan <- "config file modified" + fileEvent.Name
				}
				//Reading file to continue the watch
				watcher.Add(configFilePath)
			case err, _ := <-watcher.Errors:
				if err != nil {
					errChan <- err
					defer watcher.Close()
				}
			}
		}
	}()
}

// CreateContext creates a custom context for logging
func CreateContext(transactionID, actionID, actionName, threadID, threadName string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, transactionID)
	ctx = context.WithValue(ctx, common.ActionID, actionID)
	ctx = context.WithValue(ctx, common.ActionName, actionName)
	ctx = context.WithValue(ctx, common.ThreadID, threadID)
	ctx = context.WithValue(ctx, common.ThreadName, threadName)
	ctx = context.WithValue(ctx, common.ProcessName, podName)
	return ctx
}
