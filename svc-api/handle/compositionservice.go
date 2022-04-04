//(C) Copyright [2022] American Megatrends International LLC
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

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"

	compositionserviceproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/compositionservice"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// CompositionServiceRPCs defines all the RPC methods in compositon service
type CompositionServiceRPCs struct {
	GetCompositionServiceRPC      func(req compositionserviceproto.GetCompositionServiceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetResourceBlockCollectionRPC func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetResourceBlockRPC           func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	CreateResourceBlockRPC        func(req compositionserviceproto.CreateCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	DeleteResourceBlockRPC        func(req compositionserviceproto.DeleteCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetResourceZoneCollectionRPC  func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetResourceZoneRPC            func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	CreateResourceZoneRPC         func(req compositionserviceproto.CreateCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	DeleteResourceZoneRPC         func(req compositionserviceproto.DeleteCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	ComposeRPC                    func(req compositionserviceproto.ComposeRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetActivePoolRPC              func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetFreePoolRPC                func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetCompositionReservationsRPC func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
}

//GetCompositionService fetches all composition service
func (cs *CompositionServiceRPCs) GetCompositionService(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionServiceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetCompositionServiceRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetResourceBlockCollection fetch the Resource Blocks Instance collection
func (cs *CompositionServiceRPCs) GetResourceBlockCollection(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceBlockCollectionRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetResourceBlock get the Resource Block Instance
func (cs *CompositionServiceRPCs) GetResourceBlock(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		ResourceID:   ctx.Params().Get("id"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceBlockRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// CreateResourceBlock Create the Resource Block Instance
func (cs *CompositionServiceRPCs) CreateResourceBlock(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the create Resource Block request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	blockReq := compositionserviceproto.CreateCompositionResourceRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.CreateResourceBlockRPC(blockReq)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteResourceBlock Remove Resource Block Instance
func (cs *CompositionServiceRPCs) DeleteResourceBlock(ctx iris.Context) {
	defer ctx.Next()
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	req := compositionserviceproto.DeleteCompositionResourceRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.DeleteResourceBlockRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetResourceZoneCollection fetch the Resource zones Instance collection
func (cs *CompositionServiceRPCs) GetResourceZoneCollection(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceZoneCollectionRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetResourceZone get the Resource zone Instance
func (cs *CompositionServiceRPCs) GetResourceZone(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		ResourceID:   ctx.Params().Get("id"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceZoneRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// CreateResourceZone create Resource zone Instance
func (cs *CompositionServiceRPCs) CreateResourceZone(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the create Resource zone request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}
	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	zoneReq := compositionserviceproto.CreateCompositionResourceRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.CreateResourceZoneRPC(zoneReq)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteResourceZone remove Resource zone Instance
func (cs *CompositionServiceRPCs) DeleteResourceZone(ctx iris.Context) {
	defer ctx.Next()
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	req := compositionserviceproto.DeleteCompositionResourceRequest{
		SessionToken: sessionToken,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.DeleteResourceZoneRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// Compose Action for compose system and decompose system
func (cs *CompositionServiceRPCs) Compose(ctx iris.Context) {
	defer ctx.Next()
	var req interface{}
	err := ctx.ReadJSON(&req)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the compose request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	request, err := json.Marshal(req)
	if err != nil {
		errorMessage := "error while trying to create JSON request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	composeReq := compositionserviceproto.ComposeRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.ComposeRPC(composeReq)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetActivePool Active Resource Block Instance collection
func (cs *CompositionServiceRPCs) GetActivePool(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := cs.GetActivePoolRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetFreePool Free Resource Block instance collection
func (cs *CompositionServiceRPCs) GetFreePool(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := cs.GetFreePoolRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetCompositionReservations Compose action reservation collection
func (cs *CompositionServiceRPCs) GetCompositionReservations(ctx iris.Context) {
	defer ctx.Next()
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := cs.GetCompositionReservationsRPC(req)
	if err != nil {
		errorMessage := "RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
