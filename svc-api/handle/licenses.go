//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

// Package handle ...
package handle

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// LicenseRPCs defines all the RPC methods in license service
type LicenseRPCs struct {
	GetLicenseServiceRPC     func(ctx context.Context, req licenseproto.GetLicenseServiceRequest) (*licenseproto.GetLicenseResponse, error)
	GetLicenseCollectionRPC  func(ctx context.Context, req licenseproto.GetLicenseRequest) (*licenseproto.GetLicenseResponse, error)
	GetLicenseResourceRPC    func(ctx context.Context, req licenseproto.GetLicenseResourceRequest) (*licenseproto.GetLicenseResponse, error)
	InstallLicenseServiceRPC func(ctx context.Context, req licenseproto.InstallLicenseRequest) (*licenseproto.GetLicenseResponse, error)
}

func (lcns *LicenseRPCs) GetLicenseService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := licenseproto.GetLicenseServiceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := lcns.GetLicenseServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetLicenseCollection fetches all licenses
func (lcns *LicenseRPCs) GetLicenseCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := licenseproto.GetLicenseRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := lcns.GetLicenseCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetLicenseResource fetches license resource
func (lcns *LicenseRPCs) GetLicenseResource(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := licenseproto.GetLicenseResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := lcns.GetLicenseResourceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// InstallLicenseService installs license
func (lcns *LicenseRPCs) InstallLicenseService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var reqIn interface{}
	err := ctx.ReadJSON(&reqIn)
	if err != nil {
		errorMessage := "Error while trying to get JSON body from request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(reqIn)
	if err != nil {
		errorMessage := "while trying to create JSON request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}
	req := licenseproto.InstallLicenseRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		RequestBody:  request,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := lcns.InstallLicenseServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
