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

//Package tresponse ...
package tresponse

import (
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

//SubTask struct is used to display to the user
type SubTask struct {
	response.Response
	MembersCount    int          `json:"Members@odata.count"`
	Members         []ListMember `json:"Members"`
	TaskState       string       `json:"TaskState"`
	StartTime       time.Time    `json:"StartTime"`
	EndTime         time.Time    `json:"EndTime"`
	TaskStatus      string       `json:"TaskStatus"`
	SubTasks        string       `json:"SubTasks,omitempty"`
	TaskMonitor     string       `json:"TaskMonitor"`
	PercentComplete int32        `json:"PercentComplete,omitempty"`
	Payload         Payload      `json:"Payload,omitempty"`
	Messages        []Messages   `json:"Messages"`
}

//Task struct is used to display to the user
type Task struct {
	response.Response
	TaskState       string      `json:"TaskState"`
	StartTime       time.Time   `json:"StartTime"`
	EndTime         time.Time   `json:"EndTime,omitempty"`
	TaskStatus      string      `json:"TaskStatus"`
	SubTasks        *ListMember `json:"SubTasks,omitempty"`
	TaskMonitor     string      `json:"TaskMonitor"`
	PercentComplete int32       `json:"PercentComplete,omitempty"`
	HidePayload     bool        `json:"HidePayload,omitempty"`
	Payload         Payload     `json:"Payload,omitempty"`
	Messages        []Messages  `json:"Messages,omitempty"`
	Actions         *OemActions `json:"Actions,omitempty"`
	Oem             Oem         `json:"Oem,omitempty"`
}

//Messages struct is used to display to the user
type Messages struct {
	Message           string   `json:"Message"`
	MessageID         string   `json:"MessageId"`
	MessageArgs       []string `json:"MessageArgs"`
	Oem               Oem      `json:"Oem"`
	RelatedProperties []string `json:"RelatedProperties"`
	Resolution        string   `json:"Resolution"`
	Severity          string   `json:"Severity"`
	MessageSeverity   string   `json:"MessageSeverity,omitempty"`
}

// Oem Model
type Oem struct {
}

// Payload struct is used to give response to the user
type Payload struct {
	HTTPHeaders   []string `json:"HttpHeaders"`
	HTTPOperation string   `json:"HttpOperation"`
	JSONBody      string   `json:"JsonBody"`
	TargetURI     string   `json:"TargetUri"`
}

//TaskCollectionResponse is used to give back the response
type TaskCollectionResponse struct {
	response.Response
	MembersCount int          `json:"Members@odata.count"`
	Members      []ListMember `json:"Members"`
}

//TaskServiceResponse is used to give baxk the response
type TaskServiceResponse struct {
	response.Response
	CompletedTaskOverWritePolicy    string      `json:"CompletedTaskOverWritePolicy,omitempty"`
	DateTime                        time.Time   `json:"DateTime,omitempty"`
	LifeCycleEventOnTaskStateChange bool        `json:"LifeCycleEventOnTaskStateChange,omitempty"`
	ServiceEnabled                  bool        `json:"ServiceEnabled,omitempty"`
	Status                          Status      `json:"Status,omitempty"`
	Tasks                           Tasks       `json:"Tasks,omitempty"`
	TaskAutoDeleteTimeoutMinutes    int         `json:"TaskAutoDeleteTimeoutMinutes,omitempty"`
	Actions                         *OemActions `json:"Actions,omitempty"`
	Oem                             Oem         `json:"Oem,omitempty"`
}

//OemActions struct for oem actions
type OemActions struct {
	Oem *Oem `json:"Oem,omitempty"`
}

//Tasks struct for response
type Tasks struct {
	OdataID string `json:"@odata.id"`
}

//Status struct definition
type Status struct {
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
	Oem          Oem    `json:"Oem"`
	State        string `json:"State"`
}

// ListMember define the links for each account present in odimra
type ListMember struct {
	OdataID string `json:"@odata.id"`
}
