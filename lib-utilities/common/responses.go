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
	"context"
	"net/http"

	iris "github.com/kataras/iris/v12"
	"google.golang.org/grpc/metadata"
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

// GetContextData is used to fetch data from metadata and add it to context
func GetContextData(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewIncomingContext(ctx, md)
	if len(md[TransactionID]) > 0 {
		ctx = context.WithValue(ctx, ProcessName, md[ProcessName][0])
		ctx = context.WithValue(ctx, TransactionID, md[TransactionID][0])
		ctx = context.WithValue(ctx, ActionID, md[ActionID][0])
		ctx = context.WithValue(ctx, ActionName, md[ActionName][0])
		ctx = context.WithValue(ctx, ThreadID, md[ThreadID][0])
		ctx = context.WithValue(ctx, ThreadName, md[ThreadName][0])
	}

	return ctx
}

// CreateMetadata is used to add metadata values in context to be used in grpc calls
func CreateMetadata(ctx context.Context) context.Context {
	if ctx.Value(TransactionID) != nil {
		md := metadata.New(map[string]string{
			ProcessName:   ctx.Value(ProcessName).(string),
			TransactionID: ctx.Value(TransactionID).(string),
			ActionName:    ctx.Value(ActionName).(string),
			ActionID:      ctx.Value(ActionID).(string),
			ThreadID:      ctx.Value(ThreadID).(string),
			ThreadName:    ctx.Value(ThreadName).(string),
		})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return ctx
}

func ModifyContext(ctx context.Context, threadName, podName string) context.Context {
	ctx = context.WithValue(ctx, ThreadName, threadName)
	ctx = context.WithValue(ctx, ProcessName, podName)
	return ctx
}
