//(C) Copyright [2023] Hewlett Packard Enterprise Development LP
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

package model

// TaskState - This property shall indicate the state of the task.
type TaskState string

// Health - This property shall represent the health state of the resource
// without considering its dependent resources.  The values shall conform
// to those defined in the Redfish Specification.
type Health string

// OverWritePolicy - This property shall contain the overwrite policy for
// completed tasks.  This property shall indicate if the task service
// overwrites completed task information.
type OverWritePolicy string

// State - The known state of the resource, such as, enabled.
// This property shall indicate whether and why this component is available.
type State string

const (
	// This value shall represent that the task is newly created,
	// but has not started.
	TaskStateNew TaskState = "New"
	// This value shall represent that the task is starting.
	TaskStateStarting TaskState = "Starting"
	// This value shall represent that the task is executing.
	TaskStateRunning TaskState = "Running"
	// This value shall represent that the task has been suspended
	// but is expected to restart and is therefore not complete.
	TaskStateSuspended TaskState = "Suspended"
	// This value shall represent that the task has been interrupted
	// but is expected to restart and is therefore not complete.
	TaskStateInterrupted TaskState = "Interrupted"
	// This value shall represent that the task is pending some condition
	// and has not yet begun to execute.
	TaskStatePending TaskState = "Pending"
	// This value shall represent that the task is stopping
	// but is not yet complete.
	TaskStateStopping TaskState = "Stopping"
	// This value shall represent that the task completed successfully
	// or with warnings.
	TaskStateCompleted TaskState = "Completed"
	// This value shall represent that the task is complete
	// because an operator killed it.
	TaskStateKilled TaskState = "Killed"
	// This value shall represent that the task completed with errors.
	TaskStateException TaskState = "Exception"
	// This value shall represent that the task is now running as a service
	// and expected to continue operation until stopped or killed.
	TaskStateService TaskState = "Service"
	// "This value shall represent that the task is in the process of
	// being cancelled.
	TaskStateCancelling TaskState = "Cancelling"
	// This value shall represent that either a DELETE operation
	// on a task monitor or Task resource or by an internal process
	// cancelled the task.
	TaskStateCancelled TaskState = "Cancelled"
)

const (
	// Health is Normal.
	HealthOK Health = "OK"
	// A condition requires attention.
	HealthWarning Health = "Warning"
	// A critical condition requires immediate attention
	HealthCritical Health = "Critical"
)

const (
	// Completed tasks are not automatically overwritten.
	OverWritePolicyManual OverWritePolicy = "Manual"
	// Oldest completed tasks are overwritten.
	OverWritePolicyOldest OverWritePolicy = "Oldest"
)

const (
	// This function or resource is enabled.
	StateEnabled State = "Enabled"
	// This function or resource is disabled.
	StateDisabled State = "Disabled"
	// This function or resource is enabled but awaits an external action
	// to activate it.
	StateStandbyOffline State = "StandbyOffline"
	// This function or resource is part of a redundancy set and awaits a
	// failover or other external action to activate it.
	StateStandbySpare State = "StandbySpare"
	// This function or resource is undergoing testing, or is in the process
	// of capturing information for debugging.
	StateInTest State = "InTest"
	//This function or resource is starting.
	StateStarting State = "Starting"
	// This function or resource is either not present or detected.
	StateAbsent State = "Absent"
	// This function or resource is present but cannot be used.
	StateUnavailableOffline State = "UnavailableOffline"
	// The element does not process any commands but queues new requests.
	StateDeferring State = "Deferring"
	// The element is enabled but only processes a restricted set of commands.
	StateQuiesced State = "Quiesced"
	// The element is updating and might be unavailable or degraded.
	StateUpdating State = "Updating"
	//The element quality is within the acceptable range of operation.
	StateQualified State = "Qualified"
)

// Task - The Task schema contains information about a task that the Redfish
// task service schedules or executes.  Tasks represent operations that take
// more time than a client typically wants to wait.
// Reference	                : Task.v1_6_1.json
type Task struct {
	ODataContext      string      `json:"@odata.context,omitempty"`
	ODataEtag         string      `json:"@odata.etag,omitempty"`
	ODataID           string      `json:"@odata.id"`
	ODataType         string      `json:"@odata.type"`
	Actions           *OemActions `json:"Actions,omitempty"`
	Description       string      `json:"Description,omitempty"`
	EndTime           string      `json:"EndTime,omitempty"`
	EstimatedDuration string      `json:"EstimatedDuration,omitempty"`
	HidePayload       bool        `json:"HidePayload,omitempty"`
	ID                string      `json:"Id"`
	Messages          []*Message  `json:"Messages,omitempty"`
	Name              string      `json:"Name"`
	Payload           *Payload    `json:"Payload,omitempty"`
	PercentComplete   int         `json:"PercentComplete,omitempty"`
	StartTime         string      `json:"StartTime,omitempty"`
	SubTasks          *Link       `json:"TaskCollection,omitempty"`
	TaskMonitor       string      `json:"TaskMonitor,omitempty"`
	TaskState         TaskState   `json:"TaskState,omitempty"`
	TaskStatus        Health      `json:"TaskStatus,omitempty"`
}

// Message - This type shall contain a message that the Redfish service returns,
// as described in the Redfish Specification.
type Message struct {
	Message           string   `json:"Message,omitempty"`
	MessageArgs       []string `json:"MessageArgs,omitempty"`
	MessageID         string   `json:"MessageId"`
	Oem               Oem      `json:"Oem,omitempty"`
	RelatedProperties []string `json:"RelatedProperties,omitempty"`
	Resolution        string   `json:"Resolution,omitempty"`
	Severity          string   `json:"Severity,omitempty"`
}

// Payload - This type shall contain information detailing the HTTP and JSON
// payload information for executing this task.  This property shall not be
// included in the response if the HidePayload property is `true`.
type Payload struct {
	HTTPHeaders   []string `json:"HttpHeaders,omitempty"`
	HTTPOperation string   `json:"HttpOperation,omitempty"`
	JSONBody      string   `json:"JsonBody,omitempty"`
	TargetURI     string   `json:"TargetUri,omitempty"`
}

// TaskCollection - This Resource shall represent a Resource Collection of
// Task instances for a Redfish implementation.
// Reference	                : TTaskCollection.json
type TaskCollection struct {
	ODataContext    string  `json:"@odata.context,omitempty"`
	ODataID         string  `json:"@odata.id"`
	ODataType       string  `json:"@odata.type"`
	ODataETag       string  `json:"@odata.etag"`
	Description     string  `json:"Description,omitempty"`
	Members         []*Link `json:"Members,omitempty"`
	MembersCount    int     `json:"Members@odata.count"`
	MembersNextLink string  `json:"Members@odata.nextLink,omitempty"`
	Name            string  `json:"Name"`
	Oem             Oem     `json:"Oem,omitempty"`
}

// TaskService schema describes a task service that enables management of
// long-duration operations, includes the properties for the task service
// itself, and has links to the resource collection of tasks.
type TaskService struct {
	ODataContext                    string             `json:"@odata.context,omitempty"`
	ODataEtag                       string             `json:"@odata.etag,omitempty"`
	ODataID                         string             `json:"@odata.id"`
	ODataType                       string             `json:"@odata.type"`
	Actions                         *OemActions        `json:"Actions,omitempty"`
	CompletedTaskOverWritePolicy    OverWritePolicy    `json:"CompletedTaskOverWritePolicy,omitempty"`
	DateTime                        string             `json:"DateTime,omitempty"`
	Description                     string             `json:"Description,omitempty"`
	ID                              string             `json:"Id"`
	LifeCycleEventOnTaskStateChange bool               `json:"LifeCycleEventOnTaskStateChange,omitempty"`
	Name                            string             `json:"Name"`
	Oem                             Oem                `json:"Oem,omitempty"`
	ServiceEnabled                  bool               `json:"ServiceEnabled,omitempty"`
	Status                          *TaskServiceStatus `json:"Status,omitempty"`
	TaskAutoDeleteTimeoutMinutes    bool               `json:"TaskAutoDeleteTimeoutMinutes,omitempty"`
	Tasks                           *Link              `json:"Tasks,omitempty"`
}

// TaskServiceStatus - This type shall contain any status or health
// properties of Task service
type TaskServiceStatus struct {
	Conditions   []*Conditions `json:"Conditions,omitempty"`
	Health       Health        `json:"Health,omitempty"`
	HealthRollup Health        `json:"HealthRollup,omitempty"`
	Oem          Oem           `json:"Oem,omitempty"`
	State        State         `json:"State,omitempty"`
}

// Conditions - This type shall contain the description and details of a
// condition that exists within this resource or a related resource that
// requires attention.
type Conditions struct {
	LogEntry          *Link    `json:"LogEntry,omitempty"`
	Message           string   `json:"Message,omitempty"`
	MessageArgs       []string `json:"MessageArgs,omitempty"`
	MessageId         string   `json:"MessageId"`
	OriginOfCondition *Link    `json:"OriginOfCondition,omitempty"`
	Resolution        string   `json:"Resolution,omitempty"`
	Severity          Health   `json:"Severity,omitempty"`
	Timestamp         string   `json:"Timestamp,omitempty"`
}
