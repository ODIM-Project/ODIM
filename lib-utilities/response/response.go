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

package response

const (
	// Created is the message for successful creation
	Created = "Base.1.10.0.Created"
	// ExtendedInfo message
	ExtendedInfo = "Base.1.10.0.ExtendedInfo"
	// AccountRemoved is the message for successful removal of account
	AccountRemoved = "Base.1.10.0.AccountRemoved"
	// Success is the message for successful completion
	Success = "Base.1.10.0.Success"
	// AccountModified is the message for successful account modification
	AccountModified = "Base.1.10.0.AccountModified"
	// GeneralError defines the code at the time of General Error
	GeneralError = "Base.1.10.0.GeneralError"
	// Failure code for failed message
	Failure = "Base.1.10.0.Failed"
	// InsufficientPrivilege defines the status message at the time of Insufficient Privileges
	InsufficientPrivilege = "Base.1.10.0.InsufficientPrivilege"
	// InternalError defines the status message at the time of Internal Error
	InternalError = "Base.1.10.0.InternalError"
	// PropertyMissing defines the status message at the time of Property Missing
	PropertyMissing = "Base.1.10.0.PropertyMissing"
	// PropertyUnknown defines the status message at the time of Property Unknown
	PropertyUnknown = "Base.1.10.0.PropertyUnknown"
	// ResourceNotFound defines the status message at the time of Resource Not Found
	ResourceNotFound = "Base.1.10.0.ResourceNotFound"
	// MalformedJSON defines the status message at the time of Malformed JSON
	MalformedJSON = "Base.1.10.0.MalformedJSON"
	// PropertyValueNotInList defines the status message at the time of Property Value Not In List
	PropertyValueNotInList = "Base.1.10.0.PropertyValueNotInList"
	// NoValidSession defines the status message at the time of No Valid Session
	NoValidSession = "Base.1.10.0.NoValidSession"
	// ResourceInUse defines events aleady subscribed
	ResourceInUse = "Base.1.10.0.ResourceInUse"
	// PropertyValueFormatError defines the status message  given the correct value type but the value of that property was not supported
	PropertyValueFormatError = "Base.1.10.0.PropertyValueFormatError"
	// PropertyValueTypeError defines the message that the property is value given is having a different format
	PropertyValueTypeError = "Base.1.10.0.PropertyValueTypeError"
	// ResourceAtURIUnauthorized defines the authorization failure with plugin or other resources
	ResourceAtURIUnauthorized = "Base.1.10.0.ResourceAtUriUnauthorized"
	// CouldNotEstablishConnection defines the connection failure with plugin or other resources
	CouldNotEstablishConnection = "Base.1.10.0.CouldNotEstablishConnection"
	// QueryCombinationInvalid defines the status message at the time of invalid query
	QueryCombinationInvalid = "Base.1.10.0.QueryCombinationInvalid"
	// QueryNotSupported defines the status message at the time of not supported query
	QueryNotSupported = "Base.1.10.0.QueryNotSupported"
	// ResourceRemoved is the message for successful removal of resource
	ResourceRemoved = "ResourceEvent.1.0.3.ResourceRemoved"
	// ResourceCreated is the message for successful creation of resource
	ResourceCreated = "ResourceEvent.1.0.3.ResourceCreated"
	// TaskStarted is the message for denoting the starting of the task
	TaskStarted = "TaskEvent.1.0.3.TaskStarted"
	// ActionNotSupported defines requested POST operation is not supported by the resource
	ActionNotSupported = "Base.1.10.0.ActionNotSupported"
	// ResourceAlreadyExists indicates the request is for creation of a resource, which already exists
	ResourceAlreadyExists = "Base.1.10.0.ResourceAlreadyExists"
	// ActionParameterNotSupported indicates that the parameter supplied for the action is not supported on the resource.
	ActionParameterNotSupported = "Base.1.10.0.ActionParameterNotSupported"
	// ResourceCannotBeDeleted indicates the requested delete operation cannot be performed
	ResourceCannotBeDeleted = "Base.1.10.0.ResourceCannotBeDeleted"
	// PropertyValueConflict indicates that the requested write of a property value could not be completed, because of a conflict with another property value.
	PropertyValueConflict = "Base.1.10.0.PropertyValueConflict"
	// NoOperation  defines the status message at the time of of there is no opeartion need to be performed.
	NoOperation = "Base.1.10.0.NoOperation"
)

// Response holds the generic response from odimra
type Response struct {
	OdataType    string   `json:"@odata.type"`
	OdataID      string   `json:"@odata.id"`
	OdataContext string   `json:"@odata.context,omitempty"`
	Description  string   `json:"Description,omitempty"`
	ID           string   `json:"Id"`
	Name         string   `json:"Name"`
	Message      string   `json:"Message,omitempty"`
	MessageID    string   `json:"MessageId,omitempty"`
	MessageArgs  []string `json:"MessageArgs,omitempty"`
	NumberOfArgs int      `json:"NumberOfArgs,omitempty"`
	Severity     string   `json:"Severity,omitempty"`
	Resolution   string   `json:"Resolution,omitempty"`
}

// RPC defines the reponse which odimra service returns back as
// part of the RPC call.
//
// StatusCode defines the status code of the requested service operation.
// StatusMessage defines the message regarding the status of the requested operation.
// Header defines the headers required to create a proper response from the api gateway.
// Body defines the actual response of the requested service operation.
type RPC struct {
	StatusCode    int32
	StatusMessage string
	Header        map[string]string
	Body          interface{}
}

//CommonError holds the error response from odimra
type CommonError struct {
	Error ErrorClass `json:"error"`
}

//ErrorClass holds the properties that describe error from odimra
//
//Code indicates a specific MessageId from a Message Registry.
//Message contains error message corresponding to the message in a Message Registry.
//MessageExtendedInfo is an message objects describing one or more error messages.
type ErrorClass struct {
	Code                string `json:"code"`
	Message             string `json:"message"`
	MessageExtendedInfo []Msg  `json:"@Message.ExtendedInfo,omitempty"`
}

// Msg holds the properties of Message object
type Msg struct {
	OdataType   string        `json:"@odata.type,omitempty"`
	MessageID   string        `json:"MessageId,omitempty"`
	Message     string        `json:"Message,omitempty"`
	Severity    string        `json:"Severity,omitempty"`
	MessageArgs []interface{} `json:"MessageArgs,omitempty"`
	Resolution  string        `json:"Resolution,omitempty"`
}

//Args holds the slice of ErrArgs
type Args struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	ErrorArgs []ErrArgs
}

// ErrArgs holds the parameters to build error response
type ErrArgs struct {
	StatusMessage string
	ErrorMessage  string
	MessageArgs   []interface{}
}
