package chassis

import (
	"net/http"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
)

// GetChassisInfo is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting chassis resource information, parameters need to be passed Request .
// Request holds the Uuid and Url ,
// Url will be parsed from that search key will created
// There will be two return values for the function. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (h *Get) Handle(req *chassisproto.GetChassisRequest) response.RPC {
	//managed chassis lookup
	managedChassis := new(dmtf.Chassis)
	e := h.findInMemoryDB("Chassis", req.URL, managedChassis)
	if e == nil {
		return response.RPC{
			StatusMessage: response.Success,
			StatusCode:    http.StatusOK,
			Body:          *managedChassis,
		}
	}

	if e.ErrNo() != errors.DBKeyNotFound {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pluginClient, e := h.createPluginClient("URP")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		//urp plugin is not registered, requested chassis unknown -> status not found
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", req.URL}, nil)
	}

	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pr, pe := pluginClient.Get("/ODIM/v1/Chassis/" + req.RequestParam)
	if pe != nil {
		return pe.AsRPCResponse()
	}

	return pr.AsRPCResponse()
}

type Get struct {
	findInMemoryDB     func(table, key string, r interface{}) *errors.Error
	createPluginClient plugin.ClientFactory
}

func NewGetHandler(
	pluginClientCreator plugin.ClientFactory,
	inMemoryDBFinder func(table, key string, r interface{}) *errors.Error) *Get {

	return &Get{
		createPluginClient: pluginClientCreator,
		findInMemoryDB:     inMemoryDBFinder,
	}
}
