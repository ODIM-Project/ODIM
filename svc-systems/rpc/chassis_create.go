package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	dberrors "github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-plugin-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

func createChassis(req *chassisproto.CreateChassisRequest) (r response.RPC) {
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

	plugin, dbErr := smodel.GetPluginData(nameCarrier.Name)
	if dbErr != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, dbErr.Error(), nil, nil)
	}

	body := new(json.RawMessage)
	//todo: use translators provided by config
	e = json.Unmarshal([]byte(strings.Replace(string(req.RequestBody), "/redfish/", "/ODIM/", -1)), body)
	if e != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
	}

	pluginResponse, pluginCallErr := pmbhandle.ContactPlugin(
		"https://"+plugin.IP+":"+plugin.Port+"/ODIM/v1/Chassis",
		http.MethodPost,
		"",
		"",
		body,
		map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		},
	)

	if pluginCallErr != nil {
		return common.GeneralError(
			http.StatusInternalServerError,
			response.InternalError,
			fmt.Sprintf("Error occurred during communication with plugin(%s): %v", plugin.PluginType, pluginCallErr),
			nil,
			nil,
		)
	}

	r.StatusCode, r.Body, r.Header = jsonResponseWriter(*pluginResponse, func(toBeTransformed string) string {
		//todo: use translators provided by config
		return strings.Replace(toBeTransformed, "/ODIM/", "/redfish/", -1)
	})

	return
}

func findOrNull(conn *persistencemgr.ConnPool, table, key string) (string, error) {
	r, e := conn.Read(table, key)
	if e != nil {
		switch e.ErrNo() {
		case dberrors.DBKeyNotFound:
			return "", nil
		default:
			return "", e
		}
	}
	return r, nil
}

type linksManagedByCollection struct {
	Name  string
	Links struct {
		ManagedBy []struct {
			Oid string `json:"@odata.id"`
		}
	}
}

type nameCarrier struct {
	Name string
}
