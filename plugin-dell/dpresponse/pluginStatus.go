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

//Package dpresponse ...
package dpresponse

//PluginStatusResponse holds the information of response of PluginStatus
type PluginStatusResponse struct {
	Comment         string          `json:"_comment"`
	Name            string          `json:"Name"`
	Version         string          `json:"Version"`
	Status          Status          `json:"Status"`
	EventMessageBus EventMessageBus `json:"EventMessageBus"`
}

//Status holds information of Plugin Status
type Status struct {
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
