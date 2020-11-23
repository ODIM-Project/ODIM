package chassis

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"net/http"
)

func (d *Delete) Handle(req *chassisproto.DeleteChassisRequest) response.RPC {
	e := d.findInMemory("Chassis", req.URL, new(json.RawMessage))
	if e == nil {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "Managed Chassis cannot be deleted", []interface{}{"DELETE"}, nil)
	}

	if e.ErrNo() != errors.DBKeyNotFound {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	c, e := d.createPluginClient("URP_v1.0.0")
	if e != nil && e.ErrNo() == errors.DBKeyNotFound {
		return common.GeneralError(http.StatusMethodNotAllowed, response.ActionNotSupported, "", []interface{}{"DELETE"}, nil)
	}
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	return c.Delete(req.URL)
}

func NewDeleteHandler(createPluginClient plugin.ClientFactory, finder func(Table string, key string, r interface{}) *errors.Error) *Delete {
	return &Delete{
		createPluginClient: createPluginClient,
		findInMemory:       finder,
	}
}

type Delete struct {
	createPluginClient plugin.ClientFactory
	findInMemory       func(Table string, key string, r interface{}) *errors.Error
}
