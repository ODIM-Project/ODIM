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

// Package rpc ...
package rpc

import (
	"context"
	authproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/auth"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
)

// Auth struct helps to register service
type Auth struct{}

// IsAuthorized will accepts the request and send a request to Auth method
// from session package, if its authorized then respond with the status code.
func (a *Auth) IsAuthorized(ctx context.Context, req *authproto.AuthRequest) (*authproto.AuthResponse, error) {
	var resp authproto.AuthResponse
	statusCode, errorMessage := auth.Auth(req)
	resp.StatusCode = statusCode
	resp.StatusMessage = errorMessage
	return &resp, nil
}
