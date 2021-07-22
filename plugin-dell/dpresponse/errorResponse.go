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

//Package dpresponse ...
package dpresponse

import (
	iris "github.com/kataras/iris/v12"
)

//ErrorResopnse struct is response Error struct
type ErrorResopnse struct {
	Error Error `json:"Error"`
}

//Error struct is standard response struct
type Error struct {
	Code                string            `json:"Code"`
	Message             string            `json:"Message"`
	MessageExtendedInfo []MsgExtendedInfo `json:"@Message.ExtendedInfo"`
}

//MsgExtendedInfo struct definition
type MsgExtendedInfo struct {
	MessageID   string   `json:"MessageId"`
	Message     string   `json:"Message"`
	MessageArgs []string `json:"MessageArgs"`
}

// CreateErrorResponse will accepts the error string and create standard error response
func CreateErrorResponse(errs string) ErrorResopnse {
	var err = ErrorResopnse{
		Error{
			Code:    "Base.1.6.1.GeneralError",
			Message: "See @Message.ExtendedInfo for more information.",
			MessageExtendedInfo: []MsgExtendedInfo{
				MsgExtendedInfo{
					MessageID: "Base.1.6.1.GeneralError",
					Message:   errs,
				},
			},
		},
	}
	return err
}

// SetErrorResponse will accepts the iris context, error string and status code
// it will set error response to ctx
func SetErrorResponse(ctx iris.Context, err string, statusCode int32) {
	ctx.StatusCode(int(statusCode))
	ctx.JSON(CreateErrorResponse(err))
}
