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
package errors

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func TestCreateErrorResponse(t *testing.T) {
	errorMessage := "testing"
	type args struct {
		statusMessage string
		errorMessage  string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: InsufficientPrivileges,
			args: args{
				statusMessage: InsufficientPrivileges,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    InsufficientPrivileges,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  InsufficientPrivileges,
							Message:    "There are insufficient privileges for the account or credentials associated with the current session to perform the requested operation." + errorMessage,
							Severity:   "Critical",
							Resolution: "Either abandon the operation or change the associated access rights and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: InternalError,
			args: args{
				statusMessage: InternalError,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    InternalError,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  InternalError,
							Message:    "The request failed due to an internal service error.  The service is still operational." + errorMessage,
							Severity:   "Critical",
							Resolution: "Resubmit the request.  If the problem persists, consider resetting the service.",
						},
					},
				},
			},
		},
		{
			name: PropertyMissing,
			args: args{
				statusMessage: PropertyMissing,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyMissing,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  PropertyMissing,
							Message:    "The property is a required property and must be included in the request." + errorMessage,
							Severity:   "Warning",
							Resolution: "Ensure that the property is in the request body and has a valid value and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: PropertyValueNotInList,
			args: args{
				statusMessage: PropertyValueNotInList,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyValueNotInList,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  PropertyValueNotInList,
							Message:    "The value for the property is not in the list of acceptable values." + errorMessage,
							Severity:   "Warning",
							Resolution: "Choose a value from the enumeration list that the implementation can support and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: MalformedJSON,
			args: args{
				statusMessage: MalformedJSON,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    MalformedJSON,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  MalformedJSON,
							Message:    "The request body submitted was malformed JSON and could not be parsed by the receiving service." + errorMessage,
							Severity:   "Critical",
							Resolution: "Ensure that the request body is valid JSON and resubmit the request.",
						},
					},
				},
			},
		},
		{
			name: ResourceNotFound,
			args: args{
				statusMessage: ResourceNotFound,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceNotFound,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  ResourceNotFound,
							Message:    "The requested resource was not found." + errorMessage,
							Severity:   "Critical",
							Resolution: "Provide a valid resource identifier and resubmit the request.",
						},
					},
				},
			},
		},
		{
			name: ResourceCannotBeModified,
			args: args{
				statusMessage: ResourceCannotBeModified,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceCannotBeModified,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  ResourceCannotBeModified,
							Message:    "The requested resource Cannot be modified." + errorMessage,
							Severity:   "Critical",
							Resolution: "Do not attempt to modify a non-modifiable resource.",
						},
					},
				},
			},
		},
		{
			name: NoValidSession,
			args: args{
				statusMessage: NoValidSession,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    NoValidSession,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  NoValidSession,
							Message:    "No valid session: " + errorMessage,
							Severity:   "Critical",
							Resolution: "Do not attempt to with an invalid session.",
						},
					},
				},
			},
		},
		{
			name: UnauthorizedLoginAttempt,
			args: args{
				statusMessage: UnauthorizedLoginAttempt,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    UnauthorizedLoginAttempt,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  UnauthorizedLoginAttempt,
							Message:    "Please retry with valid login credentials or session token.",
							Severity:   "Critical",
							Resolution: "If the issue persists contact system admin for further support.",
						},
					},
				},
			},
		},
		{
			name: Unauthorized,
			args: args{
				statusMessage: Unauthorized,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    Unauthorized,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  Unauthorized,
							Message:    "Please retry with valid login credentials or session token.",
							Severity:   "Critical",
							Resolution: "If the issue persists contact system admin for further support.",
						},
					},
				},
			},
		},
		{
			name: ResourceInUse,
			args: args{
				statusMessage: ResourceInUse,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceInUse,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  ResourceInUse,
							Message:    errorMessage,
							Severity:   "OK",
							Resolution: "None",
						},
					},
				},
			},
		},
		{
			name: PropertyValueFormatError,
			args: args{
				statusMessage: PropertyValueFormatError,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyValueFormatError,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  PropertyValueFormatError,
							Message:    errorMessage,
							Severity:   "Warning",
							Resolution: "Correct the value for the property in the request body and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: ResourceCannotBeDeleted,
			args: args{
				statusMessage: ResourceCannotBeDeleted,
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceCannotBeDeleted,
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  ResourceCannotBeDeleted,
							Message:    errorMessage,
							Severity:   "Critical",
							Resolution: "Do not attempt to delete a non-deletable resource.",
						},
					},
				},
			},
		},
		{
			name: "default",
			args: args{
				statusMessage: "default",
				errorMessage:  errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    "default",
					Message: response.ErrorHelperMessage,
					MessageExtendedInfo: []MsgExtendedInfo{
						MsgExtendedInfo{
							OdataType:  messageOdataTypeStr,
							MessageID:  "default",
							Message:    "default error message" + errorMessage,
							Severity:   "Warning",
							Resolution: "default resolution.",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateErrorResponse(tt.args.statusMessage, tt.args.errorMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateErrorResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateErrResp(t *testing.T) {
	errorMessage := "testing"
	type args struct {
		code         string
		errorMessage string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: response.ActionNotSupported,
			args: args{
				code:         response.ActionNotSupported,
				errorMessage: errorMessage,
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    response.ActionNotSupported,
					Message: errorMessage,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateErrResp(tt.args.code, tt.args.errorMessage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateErrResp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackError(t *testing.T) {
	errorMessage := "testing"
	type args struct {
		errno        ErrType
		errorMessage string
	}
	tests := []struct {
		name  string
		args  args
		want  *Error
		want1 int32
		want2 string
	}{
		{
			name: "1. Postive case",
			args: args{
				errno:        DBConnFailed,
				errorMessage: errorMessage,
			},
			want: &Error{
				errNum: DBConnFailed,
				errMsg: errorMessage,
			},
			want1: http.StatusServiceUnavailable,
			want2: response.CouldNotEstablishConnection,
		},
		{
			name: "2. Postive case",
			args: args{
				errno:        InvalidAuthToken,
				errorMessage: errorMessage,
			},
			want: &Error{
				errNum: InvalidAuthToken,
				errMsg: errorMessage,
			},
			want1: http.StatusUnauthorized,
			want2: response.NoValidSession,
		},
		{
			name: "3. Postive case",
			args: args{
				errno:        UndefinedErrorType,
				errorMessage: errorMessage,
			},
			want: &Error{
				errNum: UndefinedErrorType,
				errMsg: errorMessage,
			},
			want1: http.StatusUnauthorized,
			want2: response.NoValidSession,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PackError(tt.args.errno, tt.args.errorMessage)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("1. PackError() = %v, want %v", got, tt.want)
			}
			if got.ErrNo() != tt.args.errno {
				t.Errorf("2. PackError() = %v, want %v", got.ErrNo(), tt.args.errno)
			}
			if got.Error() != errorMessage {
				t.Errorf("3. PackError() = %v, want %v", got.Error(), errorMessage)
			}
			if !reflect.DeepEqual(got.String(), fmt.Errorf("%v", errorMessage)) {
				t.Errorf("4. PackError() = %v, want %v", got.String(), fmt.Errorf("%v", errorMessage))
			}

			ret1, ret2 := got.GetAuthStatusCodeAndMessage()
			if ret1 != tt.want1 {
				t.Errorf("1. GetAuthStatusCodeAndMessage() = %v, want %v", ret1, tt.want1)
			}
			if ret2 != tt.want2 {
				t.Errorf("2. GetAuthStatusCodeAndMessage() = %v, want %v", ret2, tt.want2)
			}
		})
	}
}
