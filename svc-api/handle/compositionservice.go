//Package handle ...
package handle

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"

	compositionserviceproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/compositionservice"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-api/rpc"
	iris "github.com/kataras/iris/v12"
)

// CompositionServiceRPCs defines all the RPC methods in compositon service
type CompositionServiceRPCs struct {
	GetCompositionServiceRPC      func(req compositionserviceproto.GetCompositionServiceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetResourceBlockCollectionRPC func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
	GetResourceBlockRPC           func(req compositionserviceproto.GetCompositionResourceRequest) (*compositionserviceproto.CompositionServiceResponse, error)
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
	var res interface{}
	req := compositionserviceproto.GetCompositionServiceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetCompositionServiceRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetResourceBlockCollection(ctx iris.Context) {

	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := rpc.GetResourceBlockCollection(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetResourceBlock(ctx iris.Context) {
	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		ResourceID:   ctx.Params().Get("id"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceBlockRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) DeleteResourceBlock(ctx iris.Context) {

	var res interface{}
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
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetResourceZoneCollection(ctx iris.Context) {
	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceZoneCollectionRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetResourceZone(ctx iris.Context) {
	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
		ResourceID:   ctx.Params().Get("id"),
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}
	resp, err := cs.GetResourceZoneRPC(req)
	if err != nil {
		errorMessage := "error:  RPC error:" + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError) // TODO: add error headers
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) CreateResourceZone(ctx iris.Context) {

	var req interface{}
	var res interface{}
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
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
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

	zone_req := compositionserviceproto.CreateCompositionResourceRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.CreateResourceZoneRPC(zone_req)
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) DeleteResourceZone(ctx iris.Context) {

	var res interface{}
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
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) Compose(ctx iris.Context) {

	var req interface{}
	var res interface{}
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
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
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

	compose_req := compositionserviceproto.ComposeRequest{
		SessionToken: sessionToken,
		RequestBody:  request,
		URL:          ctx.Request().RequestURI,
	}

	resp, err := cs.ComposeRPC(compose_req)
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetActivePool(ctx iris.Context) {
	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetFreePool(ctx iris.Context) {
	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}

func (cs *CompositionServiceRPCs) GetCompositionReservations(ctx iris.Context) {
	var res interface{}
	req := compositionserviceproto.GetCompositionResourceRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error headers
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
	json.Unmarshal(resp.Body, &res)
	ctx.JSON(res)
}
