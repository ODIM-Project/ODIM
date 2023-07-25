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
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	// NULL is a constant for empty string
	NULL              = ""
	hostsFilePath     = "/etc/hosts"
	contentHeader     = "# --- User configured entries --- BEGIN"
	contentFooter     = "# --- User configured entries --- END"
	AddHostActionID   = "1000"
	AddHostActionName = "AddHosts"
	AddHostService    = "AddHosts"
)

var log = logrus.New()

func main() {

	var inputFile string
	hostName := os.Getenv("HOST_NAME")
	podName := os.Getenv("POD_NAME")
	pid := os.Getpid()
	log := logs.Log
	logs.Adorn(logrus.Fields{
		"host":       hostName,
		"process_id": podName + fmt.Sprintf("_%d", pid),
	})
	if _, err := config.SetConfiguration(); err != nil {
		log.Logger.SetFormatter(&logs.SysLogFormatter{})
		log.Fatal("Error while trying set up configuration: " + err.Error())
	}
	logs.SetFormatter(config.Data.LogFormat)
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(config.Data.LogLevel)
	transactionID := uuid.New()
	ctx := createContext(transactionID.String(), AddHostActionID, AddHostActionName, "1", AddHostService, podName)
	flag.StringVar(&inputFile, "file", "", "Path of the file which contains hosts info")
	flag.Parse()

	if inputFile == NULL {
		flag.PrintDefaults()
		os.Exit(1)
	}
	go trackConfigFileChanges(ctx, inputFile)
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
func trackConfigFileChanges(ctx context.Context, configFilePath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logs.LogWithFields(ctx).Error(err.Error())
	}
	err = watcher.Add(configFilePath)
	if err != nil {
		logs.LogWithFields(ctx).Error(err.Error())
	}
	go func() {
		for {
			select {
			case fileEvent, ok := <-watcher.Events:
				if !ok {
					continue
				}
				logs.LogWithFields(ctx).Info("event:" + fileEvent.String())
				if fileEvent.Op&fsnotify.Write == fsnotify.Write || fileEvent.Op&fsnotify.Remove == fsnotify.Remove {
					logs.LogWithFields(ctx).Info("modified file:" + fileEvent.Name)
					// update the host file
					addHost(ctx, configFilePath)
				}
				//Reading file to continue the watch
				watcher.Add(configFilePath)
			case err, _ := <-watcher.Errors:
				if err != nil {
					logs.LogWithFields(ctx).Error(err.Error())
					defer watcher.Close()
				}
			}
		}
	}()
}

func addHost(ctx context.Context, inputFile string) {
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		logs.LogWithFields(ctx).Error("Failed to read file ", err.Error())
	}

	etcHostData, err := ioutil.ReadFile(hostsFilePath)
	if err != nil {
		logs.LogWithFields(ctx).Error(err.Error())
	}

	str := strings.Split(string(etcHostData), contentHeader)
	hostsData := fmt.Sprintf("%s\n%s\n%s\n%s", str[0], contentHeader, data, contentFooter)
	err = ioutil.WriteFile(hostsFilePath, []byte(hostsData), 0644)
	if err != nil {
		logs.LogWithFields(ctx).Error(err.Error())
	}
	return
}

// createContext creates a new context based on transactionId, actionId, actionName, threadId, threadName, ProcessName
func createContext(transactionID, actionID, actionName, threadID, threadName, ProcessName string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, transactionID)
	ctx = context.WithValue(ctx, common.ActionID, actionID)
	ctx = context.WithValue(ctx, common.ActionName, actionName)
	ctx = context.WithValue(ctx, common.ThreadID, threadID)
	ctx = context.WithValue(ctx, common.ThreadName, threadName)
	ctx = context.WithValue(ctx, common.ProcessName, ProcessName)
	return ctx
}
