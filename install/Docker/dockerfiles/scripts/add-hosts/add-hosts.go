//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http:#www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	// NULL is a constant for empty string
	NULL          = ""
	hostsFilePath = "/etc/hosts"
	contentHeader = "# --- User configured entries --- BEGIN"
	contentFooter = "# --- User configured entries --- END"
)

var log = logrus.New()

func main() {

	var inputFile string

	flag.StringVar(&inputFile, "file", "", "Path of the file which contains hosts info")
	flag.Parse()

	if inputFile == NULL {
		flag.PrintDefaults()
		os.Exit(1)
	}
	go trackConfigFileChanges(inputFile)
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Fatal("Failed to read ", err.Error())
	}

	fd, err := os.OpenFile(hostsFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open file ", err.Error())
	}
	defer fd.Close()

	hostsData := fmt.Sprintf("\n%s\n%s\n%s\n", contentHeader, string(data), contentFooter)
	if _, err := fd.Write([]byte(hostsData)); err != nil {
		log.Fatal("Failed to write to file ", err.Error())
	}

	// to prevent exiting the program
	go forever()
	select {} // block forever
}

func forever() {
	for {
		time.Sleep(time.Second)
	}
}

// trackConfigFileChanges monitors the host file using fsnotfiy
func trackConfigFileChanges(configFilePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Info(err.Error())
	}
	err = watcher.Add(configFilePath)
	if err != nil {
		log.Info(err.Error())
	}
	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					continue
				}
				log.Info("event:" + fileEvent.String())
				if fileEvent.Op&fsnotify.Write == fsnotify.Write || fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
					log.Info("modified file:" + fileEvent.Name)
					// update the host file
					addHost(configFilePath)
				}
				//Reading file to continue the watch
				watcher.Add(configFilePath)
			case err, _ := <-watcher.Errors:
				if err != nil {
					log.Info(err.Error())
					defer watcher.Close()
				}
			}
		}
	}()
}

func addHost(inputFile string) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		log.Info("Failed to read file ", err.Error())
	}

	etcHostData, err := ioutil.ReadFile(hostsFilePath)
	if err != nil {
		log.Info(err.Error())
	}

	str := strings.Split(string(etcHostData), contentHeader)
	hostsData := fmt.Sprintf("%s\n%s\n%s\n%s", str[0], contentHeader, data, contentFooter)
	err = ioutil.WriteFile(hostsFilePath, []byte(hostsData), 0644)
	if err != nil {
		log.Info(err.Error())
	}
	return
}
