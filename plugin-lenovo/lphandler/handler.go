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

//Package lphandler ...
package lphandler

import (
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// MethodNotAllowed holds builds response for the unallowed http operation on Lenovo plugin URLs and returns 405 error.
func MethodNotAllowed(ctx iris.Context) {
	ctx.StatusCode(http.StatusMethodNotAllowed)
	var errorMessage string
	if strings.EqualFold(ctx.Request().Method, "POST") {
		errorMessage = "Use the CLI option to create the volume directly on the Lenovo server"
	} else {
		errorMessage = "Use the CLI option to delete the volume directly on the Lenovo server"
	}

	errArgs := &response.Args{
		Code: response.GeneralError,
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ActionNotSupported,
				MessageArgs:   []interface{}{ctx.Request().Method},
				ErrorMessage:  errorMessage,
			},
		},
	}
	ctx.JSON(errArgs.CreateGenericErrorResponse())
}
