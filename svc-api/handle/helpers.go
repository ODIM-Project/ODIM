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

//Package handle ...
package handle

import iris "github.com/kataras/iris/v12"

//SetResponseHeaders Method to Set response headers which takes iris.Context and a map as parameters
func SetResponseHeaders(ctx iris.Context, params map[string]string) iris.Context {
	ctx.ResponseWriter().Header().Set("OData-Version", "4.0")
	ctx.ResponseWriter().Header().Set("Connection", "keep-alive")
	ctx.ResponseWriter().Header().Set("X-Frame-Options", "sameorigin")
	for key, value := range params {
		ctx.ResponseWriter().Header().Set(key, value)
	}
	return ctx
}
