package redfish

import "fmt"

const (
	// Created is the message for successful creation
	Created = "Base.1.6.1.Created"
	// AccountRemoved is the message for successful removal of account
	AccountRemoved = "Base.1.6.1.AccountRemoved"
	// Success is the message for successful completion
	Success = "Base.1.6.1.Success"
	// AccountModified is the message for successful account modification
	AccountModified = "Base.1.6.1.AccountModified"
	// GeneralError defines the code at the time of General Error
	GeneralError = "Base.1.6.1.GeneralError"
	// InsufficientPrivilege defines the status message at the time of Insufficient Privileges
	InsufficientPrivilege = "Base.1.6.1.InsufficientPrivilege"
	// InternalError defines the status message at the time of Internal Error
	InternalError = "Base.1.6.1.InternalError"
	// PropertyMissing defines the status message at the time of Property Missing
	PropertyMissing = "Base.1.6.1.PropertyMissing"
	// PropertyUnknown defines the status message at the time of Property Unknown
	PropertyUnknown = "Base.1.6.1.PropertyUnknown"
	// ResourceNotFound defines the status message at the time of Resource Not Found
	ResourceNotFound = "Base.1.6.1.ResourceNotFound"
	// MalformedJSON defines the status message at the time of Malformed JSON
	MalformedJSON = "Base.1.6.1.MalformedJSON"
	// PropertyValueNotInList defines the status message at the time of Property Value Not In List
	PropertyValueNotInList = "Base.1.6.1.PropertyValueNotInList"
	// NoValidSession defines the status message at the time of No Valid Session
	NoValidSession = "Base.1.6.1.NoValidSession"
	// ResourceInUse defines events aleady subscribed
	ResourceInUse = "Base.1.6.1.ResourceInUse"
	// PropertyValueFormatError defines the status message  given the correct value type but the value of that property was not supported
	PropertyValueFormatError = "Base.1.6.1.PropertyValueFormatError"
	// PropertyValueTypeError defines the message that the property is value given is having a different format
	PropertyValueTypeError = "Base.1.6.1.PropertyValueTypeError"
	// ResourceAtURIUnauthorized defines the authorization failure with plugin or other resources
	ResourceAtURIUnauthorized = "Base.1.6.1.ResourceAtUriUnauthorized"
	// CouldNotEstablishConnection defines the connection failure with plugin or other resources
	CouldNotEstablishConnection = "Base.1.6.1.CouldNotEstablishConnection"
	// QueryCombinationInvalid defines the status message at the time of invalid query
	QueryCombinationInvalid = "Base.1.6.1.QueryCombinationInvalid"
	// QueryNotSupported defines the status message at the time of not supported query
	QueryNotSupported = "Base.1.6.1.QueryNotSupported"
	// ResourceRemoved is the message for successful removal of resource
	ResourceRemoved = "ResourceEvent.1.0.2.ResourceRemoved"
	// ResourceCreated is the message for successful creation of resource
	ResourceCreated = "ResourceEvent.1.0.2.ResourceCreated"
	// TaskStarted is the message for denoting the starting of the task
	TaskStarted = "TaskEvent.1.0.1.TaskStarted"
	// ActionNotSupported defines requested POST operation is not supported by the resource
	ActionNotSupported = "Base.1.6.1.ActionNotSupported"
	// ResourceAlreadyExists indicates the request is for creation of a resource, which already exists
	ResourceAlreadyExists = "Base.1.6.1.ResourceAlreadyExists"
	// ActionParameterNotSupported indicates that the parameter supplied for the action is not supported on the resource.
	ActionParameterNotSupported = "Base.1.6.1.ActionParameterNotSupported"
	// ResourceCannotBeDeleted indicates the requested delete operation cannot be performed
	ResourceCannotBeDeleted = "Base.1.6.1.ResourceCannotBeDeleted"
	// PropertyValueConflict indicates that the requested write of a property value could not be completed, because of a conflict with another property value.
	PropertyValueConflict = "Base.1.6.1.PropertyValueConflict"
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

func CreateError(code string, errorMessage string) CommonError {
	return CommonError{
		Error: ErrorClass{
			Code:    code,
			Message: errorMessage,
		},
	}
}

func NewError() *CommonError {
	return &CommonError{
		Error: ErrorClass{
			Code:                GeneralError,
			Message:             "An error has occurred. See ExtendedInfo for more information.",
			MessageExtendedInfo: []MsgExtendedInfo{},
		},
	}
}

func (e *CommonError) AddExtendedInfo(ei MsgExtendedInfo) *CommonError {
	e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo, ei)
	return e
}

func NewMalformedJsonMsg(errMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:  "#Message.v1_0_8.Message",
		MessageID:  MalformedJSON,
		Message:    "The request body submitted was malformed JSON and could not be parsed by the receiving service: " + errMsg,
		Severity:   "Critical",
		Resolution: "Ensure that the request body is valid JSON and resubmit the request.",
	}
}

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

func NewResourceInUseMsg(errorMsg string) MsgExtendedInfo {
	return MsgExtendedInfo{
		OdataType:  "#Message.v1_0_8.Message",
		MessageID:  ResourceInUse,
		Message:    "The change to the requested resource failed because the resource is in use or in transition: " + errorMsg,
		Severity:   "Warning",
		Resolution: "Remove the condition and resubmit the request if the operation failed.",
	}
}

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
