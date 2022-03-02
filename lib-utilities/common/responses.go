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
	"net/http"

	iris "github.com/kataras/iris/v12"
)

// commonHeaders holds the common response headers
var commonHeaders = map[string]string{
	"Connection":             "keep-alive",
	"OData-Version":          "4.0",
	"X-Frame-Options":        "sameorigin",
	"X-Content-Type-Options": "nosniff",
	"Content-type":           "application/json; charset=utf-8",
	"Cache-Control":          "no-cache, no-store, must-revalidate",
	"Transfer-Encoding":      "chunked",
}

// SetResponseHeader will add the params to the response header
func SetResponseHeader(ctx iris.Context, params map[string]string) {
	SetCommonHeaders(ctx.ResponseWriter())
	for key, value := range params {
		ctx.ResponseWriter().Header().Set(key, value)
	}
}

// SetCommonHeaders will add the common headers to the response writer
func SetCommonHeaders(w http.ResponseWriter) {
	for key, value := range commonHeaders {
		w.Header().Set(key, value)
	}
}
