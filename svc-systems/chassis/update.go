package chassis

import (
	"encoding/json"
	"github.com/prometheus/common/log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
)

func (h *Update) Handle(req *chassis.UpdateChassisRequest) response.RPC {
	pc, e := h.createPluginClient("URP")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "", []interface{}{"PATCH"}, nil)
	}
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	body := new(json.RawMessage)
	ue := json.Unmarshal(req.RequestBody, body)
	if ue != nil {
		log.Error(ue.Error())
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, "Cannot deserialize request body", nil, nil)
	}

	pr, pe := pc.Patch(req.URL, body)
	if pe != nil {
		return pe.AsRPCResponse()
	}
	return pr.AsRPCResponse()
}

type Update struct {
	createPluginClient plugin.ClientFactory
}

func NewUpdateHandler(pluginClientFactory plugin.ClientFactory) *Update {
	return &Update{
		createPluginClient: pluginClientFactory,
	}
}
