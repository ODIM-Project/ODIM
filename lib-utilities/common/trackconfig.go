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

// Package common ...
package common

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

// TrackConfigFileChanges monitors the config changes using fsnotfiy
func TrackConfigFileChanges(configFilePath string, eventChan chan<- interface{}) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err.Error())
	}
	err = watcher.Add(configFilePath)
	if err != nil {
		log.Error(err.Error())
	}
	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if fileEvent.Op&fsnotify.Write == fsnotify.Write || fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
					log.Info("modified file:" + fileEvent.Name)
					// update the odim config
					config.TLSConfMutex.Lock()
					if err := config.SetConfiguration(); err != nil {
						log.Error("error while trying to set configuration: " + err.Error())
					}
					config.TLSConfMutex.Unlock()
					eventChan <- "config file modified"
				}
				//Reading file to continue the watch
				watcher.Add(configFilePath)
			case err, _ := <-watcher.Errors:
				if err != nil {
					log.Error(err.Error())
					defer watcher.Close()
				}
			}
		}
	}()
}
