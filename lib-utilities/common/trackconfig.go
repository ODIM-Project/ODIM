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
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

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
					// update the odim config
					config.TLSConfMutex.Lock()
					if _, err := config.SetConfiguration(); err != nil {
						errChan <- fmt.Errorf("error while trying to set configuration: %s", err.Error())
					}
					config.TLSConfMutex.Unlock()
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
