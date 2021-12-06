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

// Package common ...
package common

import (
	iris "github.com/kataras/iris/v12"
	"net/http"
)

// SetResponseHeader will add the params to the response header
func SetResponseHeader(ctx iris.Context, params map[string]string) {
    SetCommonHeaders(ctx.ResponseWriter())
	for key, value := range params {
		ctx.ResponseWriter().Header().Set(key, value)
	}
}

// SetCommonHeaders will add the common headers to the response writer
func SetCommonHeaders(w http.ResponseWriter){
    w.Header().Set("Connection", "keep-alive")
	w.Header().Set("OData-Version", "4.0")
	w.Header().Set("X-Frame-Options", "sameorigin")
	w.Header().Set("X-Content-Type-Options","nosniff")
	w.Header().Set("Content-type","application/json; charset=utf-8")
	w.Header().Set("Cache-Control","no-cache")
	w.Header().Set("Transfer-Encoding","chunked")
}