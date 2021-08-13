/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redfish

import "fmt"

const (
	// Created is the message for successful creation
	Created = "Base.1.10.0.Created"
	// AccountRemoved is the message for successful removal of account
	AccountRemoved = "Base.1.10.0.AccountRemoved"
	// Success is the message for successful completion
	Success = "Base.1.10.0.Success"
	// AccountModified is the message for successful account modification
	AccountModified = "Base.1.10.0.AccountModified"
	// GeneralError defines the code at the time of General Error
	GeneralError = "Base.1.10.0.GeneralError"
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
)

//CommonError struct definition
type CommonError struct {
	Error ErrorClass `json:"error"`
}

//ErrorClass struct definition
type ErrorClass struct {
	Code                string            `json:"code"`
	Message             string            `json:"message"`
	MessageExtendedInfo []MsgExtendedInfo `json:"@Message.ExtendedInfo,omitempty"`
}

//MsgExtendedInfo struct definition
type MsgExtendedInfo struct {
	OdataType   string        `json:"@odata.type,omitempty"`
	MessageID   string        `json:"MessageId,omitempty"`
	Message     string        `json:"Message,omitempty"`
	MessageArgs []interface{} `json:"MessageArgs,omitempty"`
	Severity    string        `json:"Severity,omitempty"`
	Resolution  string        `json:"Resolution,omitempty"`
}

// CreateError creates new instance of CommonError
func CreateError(code string, errorMessage string) CommonError {
	return CommonError{
		Error: ErrorClass{
			Code:    code,
			Message: errorMessage,
		},
	}
}

// NewError creates new instance of CommonError
func NewError(msgs ...MsgExtendedInfo) *CommonError {
	return &CommonError{
		Error: ErrorClass{
			Code:                GeneralError,
			Message:             "An error has occurred. See ExtendedInfo for more information.",
			MessageExtendedInfo: msgs,
		},
	}
}

// AddExtendedInfo adds instance of MsgExtendedInfo to `CommonError.Error.MessageExtendedInfo` collection
func (e *CommonError) AddExtendedInfo(ei MsgExtendedInfo) *CommonError {
	e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo, ei)
	return e
}

// NewMalformedJSONMsg constructs instance of MsgExtendedInfo representing `MalformedJSON` error message
func NewMalformedJSONMsg(errMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:  "#Message.v1_0_8.Message",
		MessageID:  MalformedJSON,
		Message:    "The request body submitted was malformed JSON and could not be parsed by the receiving service: " + errMsg,
		Severity:   "Critical",
		Resolution: "Ensure that the request body is valid JSON and resubmit the request.",
	}
}

// NewPropertyMissingMsg constructs instance of MsgExtendedInfo representing `PropertyMissing` error message
func NewPropertyMissingMsg(missingProperty, errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:   "#Message.v1_0_8.Message",
		MessageID:   PropertyMissing,
		Message:     fmt.Sprintf("The property %v is a required property and must be included in the request: %s", missingProperty, errorMsg),
		Severity:    "Warning",
		Resolution:  "Ensure that the property is in the request body and has a valid value and resubmit the request if the operation failed.",
		MessageArgs: []interface{}{missingProperty},
	}
}

// NewPropertyValueNotInListMsg constructs instance of MsgExtendedInfo representing `PropertyValueNotInList` error message
func NewPropertyValueNotInListMsg(currentValue, propertyName, errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:   "#Message.v1_0_8.Message",
		MessageID:   PropertyValueNotInList,
		Message:     fmt.Sprintf("The value %v for the property %v is not in the list of acceptable values. %v", currentValue, propertyName, errorMsg),
		Severity:    "Warning",
		Resolution:  "Choose a value from the enumeration list that the implementation can support and resubmit the request if the operation failed.",
		MessageArgs: []interface{}{currentValue, propertyName},
	}
}

// NewPropertyValueConflictMsg constructs instance of MsgExtendedInfo representing `PropertyValueConflict` error message
func NewPropertyValueConflictMsg(propertyName, conflictingPropertyName, errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:   "#Message.v1_0_8.Message",
		MessageID:   PropertyValueConflict,
		Message:     fmt.Sprintf("The property '%v' could not be written because its value would conflict with the value of the '%v' property: %v", propertyName, conflictingPropertyName, errorMsg),
		Severity:    "Warning",
		Resolution:  "No resolution is required.",
		MessageArgs: []interface{}{propertyName, conflictingPropertyName},
	}
}

// NewResourceNotFoundMsg constructs instance of MsgExtendedInfo representing `ResourceNotFound` error message
func NewResourceNotFoundMsg(resourceType, resourceName, errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:   "#Message.v1_0_8.Message",
		MessageID:   ResourceNotFound,
		Message:     fmt.Sprintf("The requested resource of type %v named %v was not found: %v", resourceType, resourceName, errorMsg),
		Severity:    "Critical",
		Resolution:  "Provide a valid resource identifier and resubmit the request.",
		MessageArgs: []interface{}{resourceType, resourceName},
	}
}

// NewResourceInUseMsg constructs instance of MsgExtendedInfo representing `ResourceInUse` error message
func NewResourceInUseMsg(errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:  "#Message.v1_0_8.Message",
		MessageID:  ResourceInUse,
		Message:    "The change to the requested resource failed because the resource is in use or in transition: " + errorMsg,
		Severity:   "Warning",
		Resolution: "Remove the condition and resubmit the request if the operation failed.",
	}
}

// NewResourceAlreadyExistsMsg constructs instance of MsgExtendedInfo representing `ResourceAlreadyExists` error message
func NewResourceAlreadyExistsMsg(resourceType, propertyName, propertyValue, errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:   "#Message.v1_0_8.Message",
		MessageID:   ResourceAlreadyExists,
		Message:     fmt.Sprintf("The requested resource of type %v with the property %v with the value %v already exists. %v", resourceType, propertyName, propertyValue, errorMsg),
		Severity:    "Critical",
		Resolution:  "Do not repeat the create operation as the resource has already been created.",
		MessageArgs: []interface{}{resourceType, propertyName, propertyValue},
	}
}

// NewResourceAtURIUnauthorizedMsg constructs instance of MsgExtendedInfo representing `ResourceAtURIUnauthorized` error message
func NewResourceAtURIUnauthorizedMsg(uri, errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:   "#Message.v1_0_8.Message",
		MessageID:   ResourceAtURIUnauthorized,
		Message:     fmt.Sprintf("While accessing the resource at %v, the service received an authorization error. %v", uri, errorMsg),
		Severity:    "Critical",
		Resolution:  "Ensure that the appropriate access is provided for the service in order for it to access the URI.",
		MessageArgs: []interface{}{uri},
	}
}
