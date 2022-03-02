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

// Package errors ...
package errors

const (
	// BaseVersion defines the latest version of Base
	BaseVersion = "Base.1.11.0."
	// InsufficientPrivileges defines the status message at the time of Insufficient Privileges
	InsufficientPrivileges = BaseVersion + "InsufficientPrivilege"
	// InternalError defines the status message at the time of Internal Error
	InternalError = BaseVersion + "InternalError"
	// PropertyMissing defines the status message at the time of Property Missing
	PropertyMissing = BaseVersion + "PropertyMissing"
	// ResourceNotFound defines the status message at the time of Resource Not Found
	ResourceNotFound = BaseVersion + "ResourceNotFound"
	// MalformedJSON defines the status message at the time of Malformed JSON
	MalformedJSON = BaseVersion + "MalformedJSON"
	// ResourceCannotBeModified defines the status message at the time of Resource Cannot Be Modified
	ResourceCannotBeModified = BaseVersion + "ResourceCannotBeModified"
	// PropertyValueNotInList defines the status message at the time of Property Value Not In List
	PropertyValueNotInList = BaseVersion + "PropertyValueNotInList"
	// NoValidSession defines the status message at the time of No Valid Session
	NoValidSession = BaseVersion + "NoValidSession"
	// UnauthorizedLoginAttempt defines the status message at the time of Unauthorized Login Attempt
	UnauthorizedLoginAttempt = BaseVersion + "UnauthorizedLoginAttempt"
	// Unauthorized defines the status message at the time of Unauthorized service wants to perform task
	Unauthorized = BaseVersion + "Unauthorized"
	// ResourceInUse defines events aleady subscribed
	ResourceInUse = BaseVersion + "ResourceInUse"
	// PropertyValueFormatError defines the status message  given the correct value type but the value of that property was not supported
	PropertyValueFormatError = BaseVersion + "PropertyValueFormatError"
	// ResourceAtURIUnauthorized defines the authorization failure with plugin or other resources
	ResourceAtURIUnauthorized = BaseVersion + "ResourceAtUriUnauthorized"
	// CouldNotEstablishConnection defines the connection failure with plugin or other resources
	CouldNotEstablishConnection = BaseVersion + "CouldNotEstablishConnection"
	// ResourceCannotBeDeleted defines the  status message at the time of Resource deletion
	ResourceCannotBeDeleted = BaseVersion + "ResourceCannotBeDeleted"
)

// enum defined for error types
const (
	// UndefinedErrorType to be used when error type is of no significance to the caller
	UndefinedErrorType ErrType = iota + 1
	// DBKeyFetchFailed indicates DB read action failed
	DBKeyFetchFailed
	// DBKeyNotFound indicates DB read failure of non-existent key
	DBKeyNotFound
	// DBKeyAlreadyExist indicates failure for DB insert of already existing key
	DBKeyAlreadyExist
	// DBConnFailed indicates failure to establish connection with DB
	DBConnFailed
	// InvalidAuthToken indicates presented authentication token is invalid
	InvalidAuthToken
	// JSONUnmarshalFailed indicates JSON unmarshaling of data failed
	JSONUnmarshalFailed
	// DecryptionFailed indicates decryption of data failed
	DecryptionFailed
)

// constants defined for matching partial strings in error returned
const (
	// SystemNotSupportedErrString is used for matching unsupported server error
	SystemNotSupportedErrString = "computer system is not supported"
)
