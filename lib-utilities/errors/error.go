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

import (
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"net/http"
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
	OdataType  string `json:"@odata.type,omitempty"`
	MessageID  string `json:"MessageId,omitempty"`
	Message    string `json:"Message,omitempty"`
	Severity   string `json:"Severity,omitempty"`
	Resolution string `json:"Resolution,omitempty"`
}

// ErrType defines the error type
type ErrType int

// Error contains the error number and error string
// Error number to be used for comparision of error type
// instead of error string
type Error struct {
	errNum ErrType
	errMsg string
}

const messageOdataTypeStr string = "#Message.v1_1_2.Message"

// CreateErrorResponse defines the creation of the error message response body.
// As input the fuction takes status message and the error message and will
// create the response according to statusMessage
func CreateErrorResponse(statusMessage string, errorMessage string) interface{} {
	response := CommonError{
		Error: ErrorClass{
			Code:    statusMessage,
			Message: "An error has occurred. See ExtendedInfo for more information.",
		},
	}

	switch statusMessage {
	case InsufficientPrivileges:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "There are insufficient privileges for the account or credentials associated with the current session to perform the requested operation." + errorMessage,
				Severity:   "Critical",
				Resolution: "Either abandon the operation or change the associated access rights and resubmit the request if the operation failed.",
			},
		}
	case InternalError:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "The request failed due to an internal service error.  The service is still operational." + errorMessage,
				Severity:   "Critical",
				Resolution: "Resubmit the request.  If the problem persists, consider resetting the service.",
			},
		}
	case PropertyMissing:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "The property is a required property and must be included in the request." + errorMessage,
				Severity:   "Warning",
				Resolution: "Ensure that the property is in the request body and has a valid value and resubmit the request if the operation failed.",
			},
		}
	case PropertyValueNotInList:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "The value for the property is not in the list of acceptable values." + errorMessage,
				Severity:   "Warning",
				Resolution: "Choose a value from the enumeration list that the implementation can support and resubmit the request if the operation failed.",
			},
		}
	case MalformedJSON:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "The request body submitted was malformed JSON and could not be parsed by the receiving service." + errorMessage,
				Severity:   "Critical",
				Resolution: "Ensure that the request body is valid JSON and resubmit the request.",
			},
		}
	case ResourceNotFound:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "The requested resource was not found." + errorMessage,
				Severity:   "Critical",
				Resolution: "Provide a valid resource identifier and resubmit the request.",
			},
		}
	case ResourceCannotBeModified:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "The requested resource Cannot be modified." + errorMessage,
				Severity:   "Critical",
				Resolution: "Do not attempt to modify a non-modifiable resource.",
			},
		}
	case NoValidSession:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "No valid session: " + errorMessage,
				Severity:   "Critical",
				Resolution: "Do not attempt to with an invalid session.",
			},
		}
	case UnauthorizedLoginAttempt, Unauthorized:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "Please retry with valid login credentials or session token.",
				Severity:   "Critical",
				Resolution: "If the issue persists contact system admin for further support.",
			},
		}
	case ResourceInUse:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    errorMessage,
				Severity:   "OK",
				Resolution: "None",
			},
		}
	case PropertyValueFormatError:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    errorMessage,
				Severity:   "Warning",
				Resolution: "Correct the value for the property in the request body and resubmit the request if the operation failed.",
			},
		}
	case ResourceCannotBeDeleted:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    errorMessage,
				Severity:   "Critical",
				Resolution: "Do not attempt to delete a non-deletable resource.",
			},
		}
	default:
		response.Error.MessageExtendedInfo = []MsgExtendedInfo{
			MsgExtendedInfo{
				OdataType:  messageOdataTypeStr,
				MessageID:  statusMessage,
				Message:    "default error message" + errorMessage,
				Severity:   "Warning",
				Resolution: "default resolution.",
			},
		}
	}

	return response
}

// PackError defines creation of error response, which packs
// relevant error number and the error message
func PackError(errno ErrType, errmsg ...interface{}) *Error {
	return &Error{
		errNum: errno,
		errMsg: fmt.Sprint(errmsg...)}
}

// ErrNo to be used for obtaining the error number
func (err *Error) ErrNo() ErrType {
	return err.errNum
}

// Error to be used for obtaining the error message as string type
func (err *Error) Error() string {
	return err.errMsg
}

// String to be used for obtaining error message as error type
func (err *Error) String() error {
	return fmt.Errorf(err.errMsg)
}

// GetAuthStatusCodeAndMessage will return suitable status code
// status message for failures
func (err *Error) GetAuthStatusCodeAndMessage() (int32, string) {
	switch err.ErrNo() {
	case DBConnFailed:
		return http.StatusServiceUnavailable, response.CouldNotEstablishConnection
	case InvalidAuthToken:
		return http.StatusUnauthorized, response.NoValidSession
	}
	return http.StatusUnauthorized, response.NoValidSession
}

// CreateErrResp defines the creation of the error message response body without exteneded info
func CreateErrResp(code string, errorMessage string) interface{} {
	return CommonError{
		Error: ErrorClass{
			Code:    code,
			Message: errorMessage,
		},
	}
}
