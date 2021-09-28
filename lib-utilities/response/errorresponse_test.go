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

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCreateErrorResponse(t *testing.T) {
	errMsg := "error message"
	tests := []struct {
		name string
		args Args
		want interface{}
	}{
		{
			name: Success,
			args: Args{
				Code:    Success,
				Message: Success,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: Success,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    Success,
					Message: Success,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  Success,
							Message:    "Successfully Completed Request",
							Severity:   "OK",
							Resolution: "None",
						},
					},
				},
			},
		},
		{
			name: ResourceRemoved,
			args: Args{
				Code:    ResourceRemoved,
				Message: ResourceRemoved,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ResourceRemoved,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceRemoved,
					Message: ResourceRemoved,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  ResourceRemoved,
							Message:    "The resource has been removed successfully.",
							Severity:   "OK",
							Resolution: "None",
						},
					},
				},
			},
		},
		{
			name: InsufficientPrivilege,
			args: Args{
				Code:    InsufficientPrivilege,
				Message: InsufficientPrivilege,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: InsufficientPrivilege,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    InsufficientPrivilege,
					Message: InsufficientPrivilege,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  InsufficientPrivilege,
							Message:    "There are insufficient privileges for the account or credentials associated with the current session to perform the requested operation." + errMsg,
							Severity:   "Critical",
							Resolution: "Either abandon the operation or change the associated access rights and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: InternalError,
			args: Args{
				Code:    InternalError,
				Message: InternalError,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: InternalError,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    InternalError,
					Message: InternalError,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  InternalError,
							Message:    "The request failed due to an internal service error.  The service is still operational." + errMsg,
							Severity:   "Critical",
							Resolution: "Resubmit the request.  If the problem persists, consider resetting the service.",
						},
					},
				},
			},
		},
		{
			name: PropertyMissing,
			args: Args{
				Code:    PropertyMissing,
				Message: PropertyMissing,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: PropertyMissing,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyMissing,
					Message: PropertyMissing,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   PropertyMissing,
							Message:     fmt.Sprintf("The property %v is a required property and must be included in the request. %v", "test", errMsg),
							Severity:    "Warning",
							MessageArgs: []interface{}{"test"},
							Resolution:  "Ensure that the property is in the request body and has a valid value and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: PropertyValueNotInList,
			args: Args{
				Code:    PropertyValueNotInList,
				Message: PropertyValueNotInList,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: PropertyValueNotInList,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1", "test2"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyValueNotInList,
					Message: PropertyValueNotInList,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   PropertyValueNotInList,
							Message:     fmt.Sprintf("The value %v for the property %v is not in the list of acceptable values. %v", "test1", "test2", errMsg),
							Severity:    "Warning",
							MessageArgs: []interface{}{"test1", "test2"},
							Resolution:  "Choose a value from the enumeration list that the implementation can support and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: MalformedJSON,
			args: Args{
				Code:    MalformedJSON,
				Message: MalformedJSON,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: MalformedJSON,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    MalformedJSON,
					Message: MalformedJSON,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  MalformedJSON,
							Message:    "The request body submitted was malformed JSON and could not be parsed by the receiving service." + errMsg,
							Severity:   "Critical",
							Resolution: "Ensure that the request body is valid JSON and resubmit the request.",
						},
					},
				},
			},
		},
		{
			name: ResourceNotFound,
			args: Args{
				Code:    ResourceNotFound,
				Message: ResourceNotFound,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ResourceNotFound,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1", "test2"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceNotFound,
					Message: ResourceNotFound,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   ResourceNotFound,
							Message:     fmt.Sprintf("The requested resource of type %v named %v was not found. %v", "test1", "test2", errMsg),
							Severity:    "Critical",
							MessageArgs: []interface{}{"test1", "test2"},
							Resolution:  "Provide a valid resource identifier and resubmit the request.",
						},
					},
				},
			},
		},
		{
			name: NoValidSession,
			args: Args{
				Code:    NoValidSession,
				Message: NoValidSession,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: NoValidSession,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    NoValidSession,
					Message: NoValidSession,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  NoValidSession,
							Message:    "There is no valid session established with the implementation." + errMsg,
							Severity:   "Critical",
							Resolution: "Establish a session before attempting any operations.",
						},
					},
				},
			},
		},
		{
			name: ResourceInUse,
			args: Args{
				Code:    ResourceInUse,
				Message: ResourceInUse,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ResourceInUse,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceInUse,
					Message: ResourceInUse,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  ResourceInUse,
							Message:    "The change to the requested resource failed because the resource is in use or in transition." + errMsg,
							Severity:   "Warning",
							Resolution: "Remove the condition and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: PropertyValueFormatError,
			args: Args{
				Code:    PropertyValueFormatError,
				Message: PropertyValueFormatError,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: PropertyValueFormatError,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1", "test2"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyValueFormatError,
					Message: PropertyValueFormatError,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   PropertyValueFormatError,
							Message:     fmt.Sprintf("The value %v for the property %v is of a different format than the property can accept. %v", "test1", "test2", errMsg),
							Severity:    "Warning",
							MessageArgs: []interface{}{"test1", "test2"},
							Resolution:  "Correct the value for the property in the request body and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: ResourceAtURIUnauthorized,
			args: Args{
				Code:    ResourceAtURIUnauthorized,
				Message: ResourceAtURIUnauthorized,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ResourceAtURIUnauthorized,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceAtURIUnauthorized,
					Message: ResourceAtURIUnauthorized,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   ResourceAtURIUnauthorized,
							Message:     fmt.Sprintf("While accessing the resource at %v, the service received an authorization error. %v", "test1", errMsg),
							Severity:    "Critical",
							MessageArgs: []interface{}{"test1"},
							Resolution:  "Ensure that the appropriate access is provided for the service in order for it to access the URI.",
						},
					},
				},
			},
		},
		{
			name: CouldNotEstablishConnection,
			args: Args{
				Code:    CouldNotEstablishConnection,
				Message: CouldNotEstablishConnection,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: CouldNotEstablishConnection,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    CouldNotEstablishConnection,
					Message: CouldNotEstablishConnection,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   CouldNotEstablishConnection,
							Message:     fmt.Sprintf("The service failed to establish a connection with the URI %v. %v", "test1", errMsg),
							Severity:    "Critical",
							MessageArgs: []interface{}{"test1"},
							Resolution:  "Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.",
						},
					},
				},
			},
		},
		{
			name: ActionNotSupported,
			args: Args{
				Code:    ActionNotSupported,
				Message: ActionNotSupported,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ActionNotSupported,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"POST"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ActionNotSupported,
					Message: ActionNotSupported,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   ActionNotSupported,
							Message:     "The action POST is not supported by the resource. " + errMsg,
							Severity:    "Critical",
							MessageArgs: []interface{}{"POST"},
							Resolution:  "The action supplied cannot be resubmitted to the implementation. Perhaps the action was invalid, the wrong resource was the target or the implementation documentation may be of assistance.",
						},
					},
				},
			},
		},
		{
			name: ResourceAlreadyExists,
			args: Args{
				Code:    ResourceAlreadyExists,
				Message: ResourceAlreadyExists,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ResourceAlreadyExists,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1", "test2", "test3"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ResourceAlreadyExists,
					Message: ResourceAlreadyExists,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   ResourceAlreadyExists,
							Message:     "The requested resource of type test1 with the property test2 with the value test3 already exists. " + errMsg,
							Severity:    "Critical",
							MessageArgs: []interface{}{"test1", "test2", "test3"},
							Resolution:  "Do not repeat the create operation as the resource has already been created.",
						},
					},
				},
			},
		},
		{
			name: QueryCombinationInvalid,
			args: Args{
				Code:    QueryCombinationInvalid,
				Message: QueryCombinationInvalid,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: QueryCombinationInvalid,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    QueryCombinationInvalid,
					Message: QueryCombinationInvalid,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  QueryCombinationInvalid,
							Message:    "Two or more query parameters in the request cannot be used together." + errMsg,
							Severity:   "Warning",
							Resolution: "Remove one or more of the query parameters and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: ActionParameterNotSupported,
			args: Args{
				Code:    ActionParameterNotSupported,
				Message: ActionParameterNotSupported,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: ActionParameterNotSupported,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1", "test2"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    ActionParameterNotSupported,
					Message: ActionParameterNotSupported,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   ActionParameterNotSupported,
							Message:     errMsg,
							Severity:    "Warning",
							MessageArgs: []interface{}{"test1", "test2"},
							Resolution:  "Modify the parameter supplied and resubmit the request if the operation failed.",
						},
					},
				},
			},
		},
		{
			name: PropertyValueConflict,
			args: Args{
				Code:    PropertyValueConflict,
				Message: PropertyValueConflict,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: PropertyValueConflict,
						ErrorMessage:  errMsg,
						MessageArgs:   []interface{}{"test1", "test2"},
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    PropertyValueConflict,
					Message: PropertyValueConflict,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:   ErrorMessageOdataType,
							MessageID:   PropertyValueConflict,
							Message:     fmt.Sprintf("The property '%v' could not be written because its value would conflict with the value of the '%v' property, %v", "test1", "test2", errMsg),
							Severity:    "Warning",
							MessageArgs: []interface{}{"test1", "test2"},
							Resolution:  "No resolution is required.",
						},
					},
				},
			},
		},
		{
			name: NoOperation,
			args: Args{
				Code:    NoOperation,
				Message: NoOperation,
				ErrorArgs: []ErrArgs{
					ErrArgs{
						StatusMessage: NoOperation,
						ErrorMessage:  errMsg,
					},
				},
			},
			want: CommonError{
				Error: ErrorClass{
					Code:    NoOperation,
					Message: NoOperation,
					MessageExtendedInfo: []Msg{
						Msg{
							OdataType:  ErrorMessageOdataType,
							MessageID:  NoOperation,
							Message:    "The request body submitted contain no data to act upon and no changes to the resource took place.",
							Severity:   "Warning",
							Resolution: "Add properties in the JSON object and resubmit the request.",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.CreateGenericErrorResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateErrorResponse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
