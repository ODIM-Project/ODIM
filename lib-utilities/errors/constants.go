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
	// InsufficientPrivileges defines the status message at the time of Insufficient Privileges
	InsufficientPrivileges = "Base.1.10.0.InsufficientPrivilege"
	// InternalError defines the status message at the time of Internal Error
	InternalError = "Base.1.10.0.InternalError"
	// PropertyMissing defines the status message at the time of Property Missing
	PropertyMissing = "Base.1.10.0.PropertyMissing"
	// ResourceNotFound defines the status message at the time of Resource Not Found
	ResourceNotFound = "Base.1.10.0.ResourceNotFound"
	// MalformedJSON defines the status message at the time of Malformed JSON
	MalformedJSON = "Base.1.10.0.MalformedJSON"
	// ResourceCannotBeModified defines the status message at the time of Resource Cannot Be Modified
	ResourceCannotBeModified = "Base.1.10.0.ResourceCannotBeModified"
	// PropertyValueNotInList defines the status message at the time of Property Value Not In List
	PropertyValueNotInList = "Base.1.10.0.PropertyValueNotInList"
	// NoValidSession defines the status message at the time of No Valid Session
	NoValidSession = "Base.1.10.0.NoValidSession"
	// UnauthorizedLoginAttempt defines the status message at the time of Unauthorized Login Attempt
	UnauthorizedLoginAttempt = "Base.1.10.0.UnauthorizedLoginAttempt"
	// Unauthorized defines the status message at the time of Unauthorized service wants to perform task
	Unauthorized = "Base.1.10.0.Unauthorized"
	// ResourceInUse defines events aleady subscribed
	ResourceInUse = "Base.1.10.0.ResourceInUse"
	// PropertyValueFormatError defines the status message  given the correct value type but the value of that property was not supported
	PropertyValueFormatError = "Base.1.10.0.PropertyValueFormatError"
	// ResourceAtURIUnauthorized defines the authorization failure with plugin or other resources
	ResourceAtURIUnauthorized = "Base.1.10.0.ResourceAtUriUnauthorized"
	// CouldNotEstablishConnection defines the connection failure with plugin or other resources
	CouldNotEstablishConnection = "Base.1.10.0.CouldNotEstablishConnection"
	// ResourceCannotBeDeleted defines the  status message at the time of Resource deletion
	ResourceCannotBeDeleted = "Base.1.10.0.ResourceCannotBeDeleted"
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
