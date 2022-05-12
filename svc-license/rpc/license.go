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

package rpc

import (
	"context"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/license"
)

func (l *License) GetLicenseService(ctx context.Context, req *licenseproto.GetLicenseServiceRequest) (*licenseproto.GetLicenseResponse, error) {
	resp := &licenseproto.GetLicenseResponse{}
	authResp := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, l.connector.GetLicenseService(req))
	return resp, nil
}

func (l *License) GetLicenseCollection(ctx context.Context, req *licenseproto.GetLicenseRequest) (*licenseproto.GetLicenseResponse, error) {
	resp := &licenseproto.GetLicenseResponse{}
	authResp := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, l.connector.GetLicenseCollection(req))
	return resp, nil
}

func (l *License) GetLicenseResource(ctx context.Context, req *licenseproto.GetLicenseResourceRequest) (*licenseproto.GetLicenseResponse, error) {
	resp := &licenseproto.GetLicenseResponse{}
	authResp := l.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		fillProtoResponse(resp, authResp)
		return resp, nil
	}
	fillProtoResponse(resp, l.connector.GetLicenseResource(req))
	return resp, nil
}
