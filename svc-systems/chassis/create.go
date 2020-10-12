package chassis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

func Create(req *chassisproto.CreateChassisRequest) response.RPC {
	mbc := new(linksManagedByCollection)
	e := json.Unmarshal(req.RequestBody, mbc)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("cannot deserialize request: %v", e), nil, nil)
	}

	if len(mbc.Links.ManagedBy) == 0 {
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, "", []interface{}{"Links.ManagedBy[0]"}, nil)
	}

	inMemoryConn, dbErr := common.GetDBConnection(common.InMemory)
	if dbErr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("cannot acquire database connection: %v", dbErr), nil, nil)
	}

	m, e := findOrNull(inMemoryConn, "Managers", mbc.Links.ManagedBy[0].Oid)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("error occured during database access: %v", e), nil, nil)
	}

	if m == "" {
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, "", []interface{}{"Manager", mbc.Links.ManagedBy[0].Oid}, nil)
	}

	//todo: not sure why manager in redis is quoted
	m, e = strconv.Unquote(m)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	nameCarrier := new(nameCarrier)
	e = json.Unmarshal([]byte(m), nameCarrier)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pluginConfig, dbErr := smodel.GetPluginData(nameCarrier.Name)
	if dbErr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, dbErr.Error(), nil, nil)
	}

	body := new(json.RawMessage)
	//todo: use translators provided by config
	e = json.Unmarshal([]byte(strings.Replace(string(req.RequestBody), "/redfish/", "/ODIM/", -1)), body)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pr, pe := plugin.NewClient(pluginConfig).Post("/ODIM/v1/Chassis", body)
	if pe != nil {
		return pe.AsRPCResponse()
	}

	return pr.AsRPCResponse()
}

//{
//	"Links" : {
//		"ManagedBy": [
//			"@odata.id": "/redfish/v1/Managers/1"
//		]
//	}
//}
type linksManagedByCollection struct {
	Links struct {
		ManagedBy []struct {
			Oid string `json:"@odata.id"`
		}
	}
}

//{
//	"Name" : "name"
//}
type nameCarrier struct {
	Name string
}
