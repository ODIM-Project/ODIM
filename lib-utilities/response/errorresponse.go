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
	"log"
	"net/http"
	"reflect"
)

const (
	//ErrorHelperMessage holds helper error message sent in error response
	ErrorHelperMessage = "An error has occurred. See ExtendedInfo for more information."
	//ErrorMessageOdataType holds message registry version
	ErrorMessageOdataType               = "#Message.v1_0_8.Message"
	propertyMissingArgCount             = 1
	propertyValueNotInListArgCount      = 2
	propertyValueTypeErrorArgCount      = 2
	resourceNotFoundArgCount            = 2
	propertyValueFormatErrorArgCount    = 2
	resourceAtURIUnauthorizedArgCount   = 1
	couldNotEstablishConnectionArgCount = 1
	actionNotSupportedArgCount          = 1
	resourceAlreadyExistsArgCount       = 3
	actionParameterNotSupportedArgCount = 2
	propertyUnknownArgCount             = 1
	propertyValueConflict               = 2
)

//ValidateParamTypes will compare string slices and returns bool
func ValidateParamTypes(paramTypes []string, actualParamTypes []string) bool {
	if len(paramTypes) != len(actualParamTypes) {
		return false
	}
	for i := range actualParamTypes {
		if actualParamTypes[i] != paramTypes[i] {
			return false
		}
	}
	return true
}

// CreateGenericErrorResponse will fill the error response with respective data
func (a *Args) CreateGenericErrorResponse() CommonError {
	var e CommonError

	if a.Message == "" {
		a.Message = ErrorHelperMessage
	}
	e.Error = ErrorClass{
		Code:    a.Code,
		Message: a.Message,
	}

	for _, errArg := range a.ErrorArgs {
		switch errArg.StatusMessage {
		case Success:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "Successfully Completed Request",
					Severity:   "OK",
					Resolution: "None",
				})
		case ResourceRemoved:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "The resource has been removed successfully.",
					Severity:   "OK",
					Resolution: "None",
				})
		case InsufficientPrivilege:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "There are insufficient privileges for the account or credentials associated with the current session to perform the requested operation." + errArg.ErrorMessage,
					Severity:   "Critical",
					Resolution: "Either abandon the operation or change the associated access rights and resubmit the request if the operation failed.",
				})
		case InternalError:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "The request failed due to an internal service error.  The service is still operational." + errArg.ErrorMessage,
					Severity:   "Critical",
					Resolution: "Resubmit the request.  If the problem persists, consider resetting the service.",
				})
		case PropertyMissing:
			if len(errArg.MessageArgs) != propertyMissingArgCount {
				log.Println("warning: MessageArgs in PropertyMissing response is missing")
			}
			ParamTypes := []string{"string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in PropertyMissing response is missing")
			}

			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The property %v is a required property and must be included in the request. %v", errArg.MessageArgs[0], errArg.ErrorMessage),
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Ensure that the property is in the request body and has a valid value and resubmit the request if the operation failed.",
				})
		case PropertyUnknown:
			if len(errArg.MessageArgs) != propertyUnknownArgCount {
				log.Println("warning: MessageArgs in PropertyUnknown response is missing")
			}
			ParamTypes := []string{"string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in PropertyUnknown response is missing")
			}

			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The property %v is an unknown property and must not be included in the request. %v", errArg.MessageArgs[0], errArg.ErrorMessage),
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Ensure that the request body has valid properties with proper cases and resubmit the request.",
				})
		case PropertyValueNotInList:
			if len(errArg.MessageArgs) != propertyValueNotInListArgCount {
				log.Println("warning: MessageArgs in PropertyValueNotInList response is missing")
			}
			ParamTypes := []string{"string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in PropertyValueNotInList response is missing")
			}

			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The value %v for the property %v is not in the list of acceptable values. %v", errArg.MessageArgs[0], errArg.MessageArgs[1], errArg.ErrorMessage),
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Choose a value from the enumeration list that the implementation can support and resubmit the request if the operation failed.",
				})
		case PropertyValueTypeError:
			if len(errArg.MessageArgs) != propertyValueTypeErrorArgCount {
				log.Println("warning: MessageArgs in PropertyValueTypeError response is missing")
			}
			ParamTypes := []string{"string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in PropertyValueTypeError response is missing")
			}

			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The value %v for the property %v is of a different type than the property can accept. %v", errArg.MessageArgs[0], errArg.MessageArgs[1], errArg.ErrorMessage),
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Correct the value for the property in the request body and resubmit the request if the operation failed.",
				})
		case MalformedJSON:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "The request body submitted was malformed JSON and could not be parsed by the receiving service." + errArg.ErrorMessage,
					Severity:   "Critical",
					Resolution: "Ensure that the request body is valid JSON and resubmit the request.",
				})
		case ResourceNotFound:
			if len(errArg.MessageArgs) != resourceNotFoundArgCount {
				log.Println("warning: MessageArgs in ResourceNotFound response is missing")
			}
			ParamTypes := []string{"string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in ResourceNotFound response is missing")
			}

			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The requested resource of type %v named %v was not found. %v", errArg.MessageArgs[0], errArg.MessageArgs[1], errArg.ErrorMessage),
					Severity:    "Critical",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Provide a valid resource identifier and resubmit the request.",
				})
		case NoValidSession:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "There is no valid session established with the implementation." + errArg.ErrorMessage,
					Severity:   "Critical",
					Resolution: "Establish a session before attempting any operations.",
				})
		case ResourceInUse:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "The change to the requested resource failed because the resource is in use or in transition." + errArg.ErrorMessage,
					Severity:   "Warning",
					Resolution: "Remove the condition and resubmit the request if the operation failed.",
				})
		case PropertyValueFormatError:
			if len(errArg.MessageArgs) != propertyValueFormatErrorArgCount {
				log.Println("warning: MessageArgs in PropertyValueFormatError response is missing")
			}
			ParamTypes := []string{"string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in PropertyValueFormatError response is missing")
			}
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The value %v for the property %v is of a different format than the property can accept. %v", errArg.MessageArgs[0], errArg.MessageArgs[1], errArg.ErrorMessage),
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Correct the value for the property in the request body and resubmit the request if the operation failed.",
				})
		case ResourceAtURIUnauthorized:
			if len(errArg.MessageArgs) != resourceAtURIUnauthorizedArgCount {
				log.Println("warning: MessageArgs in ResourceAtURIUnauthorized response is missing")
			}
			ParamTypes := []string{"string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in RResourceAtURIUnauthorized response is missing")
			}
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("While accessing the resource at %v, the service received an authorization error. %v", errArg.MessageArgs[0], errArg.ErrorMessage),
					Severity:    "Critical",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Ensure that the appropriate access is provided for the service in order for it to access the URI.",
				})
		case CouldNotEstablishConnection:
			if len(errArg.MessageArgs) != couldNotEstablishConnectionArgCount {
				log.Println("warning: MessageArgs in CouldNotEstablishConnection response is missing")
			}
			ParamTypes := []string{"string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in CouldNotEstablishConnection response is missing")
			}
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The service failed to establish a connection with the URI %v. %v", errArg.MessageArgs[0], errArg.ErrorMessage),
					Severity:    "Critical",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Ensure that the URI contains a valid and reachable node name, protocol information and other URI components.",
				})
		case ActionNotSupported:
			errMsgArg := ""
			if len(errArg.MessageArgs) != actionNotSupportedArgCount {
				log.Println("warning: MessageArgs in ActionNotSupported response is missing")
			} else {
				errMsgArg = errArg.MessageArgs[0].(string)
			}
			ParamTypes := []string{"string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in ActionNotSupported response is missing")
			}
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The action %v is not supported by the resource. %v", errArg.MessageArgs[0], errMsgArg),
					Severity:    "Critical",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "The action supplied cannot be resubmitted to the implementation. Perhaps the action was invalid, the wrong resource was the target or the implementation documentation may be of assistance.",
				})
		case ResourceAlreadyExists:
			if len(errArg.MessageArgs) != resourceAlreadyExistsArgCount {
				log.Println("warning: MessageArgs in ResourceAlreadyExists response is missing")
			}
			ParamTypes := []string{"string", "string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in ResourceAlreadyExists response is missing")
			}
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The requested resource of type %v with the property %v with the value %v already exists. %v", errArg.MessageArgs[0], errArg.MessageArgs[1], errArg.MessageArgs[2], errArg.ErrorMessage),
					Severity:    "Critical",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Do not repeat the create operation as the resource has already been created.",
				})
		case QueryCombinationInvalid:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    "Two or more query parameters in the request cannot be used together." + errArg.ErrorMessage,
					Severity:   "Warning",
					Resolution: "Remove one or more of the query parameters and resubmit the request if the operation failed.",
				})
		case QueryNotSupported:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  ErrorMessageOdataType,
					MessageID:  errArg.StatusMessage,
					Message:    fmt.Sprintf("Querying is not supported by the implementation. %v", errArg.ErrorMessage),
					Severity:   "Warning",
					Resolution: "Remove the query parameters and resubmit the request if the operation failed.",
				})
		case ActionParameterNotSupported:
			if len(errArg.MessageArgs) != actionParameterNotSupportedArgCount {
				log.Println("warning: MessageArgs in ActionParameterNotSupported response is missing")
			}
			ParamTypes := []string{"string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in ActionParameterNotSupported response is missing")
			}
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     errArg.ErrorMessage,
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "Modify the parameter supplied and resubmit the request if the operation failed.",
				})
		case ResourceCannotBeDeleted:
			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:  "#Message.v1_0_8.Message",
					MessageID:  errArg.StatusMessage,
					Message:    "The delete request failed because the resource requested cannot be deleted." + errArg.ErrorMessage,
					Severity:   "Critical",
					Resolution: "Do not attempt to delete a non-deletable resource.",
				})
		case PropertyValueConflict:
			if len(errArg.MessageArgs) != propertyValueConflict {
				log.Println("warning: MessageArgs in PropertyValueConflict response is missing")
			}
			ParamTypes := []string{"string", "string"}
			actualParamTypes := []string{}
			for i := 0; i < len(errArg.MessageArgs); i++ {
				actualParamTypes = append(actualParamTypes, reflect.TypeOf(errArg.MessageArgs[i]).String())
			}
			if !ValidateParamTypes(ParamTypes, actualParamTypes) {
				log.Println("warning: Paramtypes in PropertyValueConflict response is missing")
			}

			e.Error.MessageExtendedInfo = append(e.Error.MessageExtendedInfo,
				Msg{
					OdataType:   ErrorMessageOdataType,
					MessageID:   errArg.StatusMessage,
					Message:     fmt.Sprintf("The property '%v' could not be written because its value would conflict with the value of the '%v' property, %v", errArg.MessageArgs[0], errArg.MessageArgs[1], errArg.ErrorMessage),
					Severity:    "Warning",
					MessageArgs: errArg.MessageArgs,
					Resolution:  "No resolution is required.",
				})
		}
	}
	return e
}

//CreateInternalErrorResponse is used to create internal server error response
func (resp *RPC) CreateInternalErrorResponse(errorMessage string) {
	resp.StatusCode = http.StatusInternalServerError
	resp.StatusMessage = InternalError
	messageArgs := []interface{}{}
	args := Args{
		Code:    GeneralError,
		Message: "",
		ErrorArgs: []ErrArgs{
			ErrArgs{
				StatusMessage: resp.StatusMessage,
				ErrorMessage:  errorMessage,
				MessageArgs:   messageArgs,
			},
		},
	}
	resp.Body = args.CreateGenericErrorResponse()
}
