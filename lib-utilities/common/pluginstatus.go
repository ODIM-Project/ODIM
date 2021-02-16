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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

// PluginStatus holds the data required for continuously checking the plugin health
type PluginStatus struct {
	// Method - Method for communicating with Plugin
	Method string
	// Token - plugin session token
	Token string
	// RequestBody - holds the data for request body
	RequestBody StatusRequest
	// ResponseWaitTime - the resposne wait time in seconds
	ResponseWaitTime int
	// Count - the number of times the status check should happen
	Count int
	// RetryInterval the wait time in minutes before initiating the next request
	RetryInterval int
	// PluginIP
	PluginIP string
	//PluginPort
	PluginPort string
	//PluginUsername
	PluginUsername string
	//PluginUserPassword
	PluginUserPassword string
	//PluginPrefferedAuthType
	PluginPrefferedAuthType string
	//CACertificate to use while making HTTP queries
	CACertificate *[]byte
}

// StatusRequest is the plugin request for status check
type StatusRequest struct {
	Comment string `json:"_comment"`
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

// PluginStats holds the actual status details
type PluginStats struct {
	Resources     string `json:"Resources"`
	Subscriptions string `json:"Subscriptions"`
}

// StatusResponse is the general reponse for status check
type StatusResponse struct {
	Comment         string                `json:"_comment"`
	Name            string                `json:"Name"`
	Version         string                `json:"Version"`
	Status          *PluginResponseStatus `json:"Status"`
	EventMessageBus *EventMessageBus      `json:"EventMessageBus"`
}

// PluginResponseStatus hold status data of Plugin
type PluginResponseStatus struct {
	Available string `json:"Available"`
	Uptime    string `json:"Uptime"`
	TimeStamp string `json:"TimeStamp"`
}

//EventMessageBus holds the  information of  EMB Broker type and EMBQueue information
type EventMessageBus struct {
	EmbType  string     `json:"EmbType"`
	EmbQueue []EmbQueue `json:"EmbQueue"`
}

//EmbQueue holds the  information of Queue Name and Queue Description
type EmbQueue struct {
	QueueName string `json:"EmbQueueName"`
	QueueDesc string `json:"EmbQueueDesc"`
}

// CheckStatus will check the for the status health of the plugin in a frequent interval
// The function will return the following
// bool: if true, plugin is alive
// int: number of tries on the plugin before concluding
// error: will contains the error details, if any happend during the tries
// queueList:will contains the list of queues to which events are published
func (p *PluginStatus) CheckStatus() (bool, int, []string, error) {

	jsonStr, err := json.Marshal(p.RequestBody)
	if err != nil {
		return false, 1, make([]string, 0), fmt.Errorf("error while unmarshaling the request: %v", err)
	}
	var queueList []string
	requestBody := bytes.NewBuffer(jsonStr)
	var statusLog string
	for i := 0; i < p.Count; i++ {
		statusChan := make(chan bool)
		errChan := make(chan error)
		queueListChan := make(chan []string)
		go p.getStatus(requestBody, statusChan, queueListChan, errChan)
		go responseTimer(p.ResponseWaitTime, statusChan, queueListChan, errChan)

		roundError := <-errChan
		if roundError != nil {
			statusLog = statusLog + " LOGS FROM TRY " + strconv.Itoa(i+1) + ":" + roundError.Error()
		}
		alive := <-statusChan
		queueList = <-queueListChan
		if alive {
			err = nil
			if statusLog != "" {
				err = fmt.Errorf("error logs: %v", statusLog)
			}
			return true, i + 1, queueList, err
		}
		time.Sleep(time.Duration(p.RetryInterval) * time.Minute)

	}

	return false, p.Count, queueList, fmt.Errorf("error: maximum retries are over. unable to contact the plugin: error logs: %v", statusLog)
}

// responseTimer helps the CheckStatus function to keep an eye on the response wait time
func responseTimer(waitTime int, statusChan chan bool, queueListChan chan []string, errChan chan error) {
	time.Sleep(time.Duration(waitTime) * time.Second)
	errChan <- fmt.Errorf("error: wait time exceeded and got no response from plugin")
	statusChan <- false
	queueListChan <- make([]string, 0)
}

// getStatus helps the CheckStatus by making a call to the plugin for the status
func (p *PluginStatus) getStatus(requestBody *bytes.Buffer, statusChan chan bool, queueListChan chan []string, errChan chan error) {
	url := fmt.Sprintf("https://%s:%s/ODIM/v1/Status", p.PluginIP, p.PluginPort)
	req, err := http.NewRequest(p.Method, url, requestBody)
	var queueList = make([]string, 0)
	if err != nil {
		errChan <- fmt.Errorf("error while trying to create the request: %v", err)
		statusChan <- false
		queueListChan <- queueList
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if strings.EqualFold(p.PluginPrefferedAuthType, "XAuthToken") {
		if err := p.login(); err != nil {
			errChan <- err
			statusChan <- false
			queueListChan <- queueList
			return
		}
		req.Header.Set("X-Auth-Token", p.Token)
	} else {
		req.SetBasicAuth(p.PluginUsername, p.PluginUserPassword)
	}
	httpConf := &config.HTTPConfig{
		CACertificate: p.CACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()

	config.TLSConfMutex.RLock()
	resp, err := httpClient.Do(req)
	config.TLSConfMutex.RUnlock()
	if err != nil {
		errChan <- fmt.Errorf("error while trying to make the request to the plugin: %v", err)
		statusChan <- false
		queueListChan <- queueList
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("error: expected response from plugin %v, but got %v", http.StatusOK, resp.StatusCode)
		statusChan <- false
		queueListChan <- queueList
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errChan <- fmt.Errorf("error while trying to read the response body: %v", err)
		statusChan <- false
		queueListChan <- queueList
		return
	}

	var bodyData StatusResponse
	err = json.Unmarshal(body, &bodyData)
	if err != nil {
		errChan <- fmt.Errorf("error while trying to unmarshal the response: %v", err)
		statusChan <- false
		queueListChan <- queueList
		return
	}

	if bodyData.EventMessageBus != nil {
		for i := 0; i < len(bodyData.EventMessageBus.EmbQueue); i++ {
			queueList = append(queueList, bodyData.EventMessageBus.EmbQueue[i].QueueName)
		}
	}
	if bodyData.Status != nil {
		if strings.EqualFold(bodyData.Status.Available, "yes") {
			errChan <- nil
			statusChan <- true
			queueListChan <- queueList
			return
		}
	}
	errChan <- fmt.Errorf("error: expected plugin status: yes, but got %v", bodyData.Status.Available)
	statusChan <- false
	queueListChan <- queueList
}

// login creates the session token with plugin
func (p *PluginStatus) login() error {
	jsonStr, err := json.Marshal(map[string]string{
		"Username": p.PluginUsername,
		"Password": p.PluginUserPassword,
	})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://%s:%s/ODIM/v1/Sessions", p.PluginIP, p.PluginPort)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	httpConf := &config.HTTPConfig{
		CACertificate: p.CACertificate,
	}
	httpClient, err := httpConf.GetHTTPClientObj()

	config.TLSConfMutex.RLock()
	resp, err := httpClient.Do(req)
	config.TLSConfMutex.RUnlock()
	if err != nil {
		return fmt.Errorf("error while trying to make the request to the plugin: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error :Unable to create session with plugin %v", resp.StatusCode)
	}
	p.Token = resp.Header.Get("X-Auth-Token")
	return nil
}
