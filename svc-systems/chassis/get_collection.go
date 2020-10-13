package chassis

import (
	"log"
	"net/http"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
)

func NewGetCollectionHandler(
	pdp func(pluginID string) (smodel.Plugin, *errors.Error),
	imkp func(table string) ([]string, error)) *GetCollectionHandler {

	return &GetCollectionHandler{
		&sourceProviderImpl{
			getPluginConfig: pdp,
			getAllKeys:      imkp,
		},
	}
}

type GetCollectionHandler struct {
	sourcesProvider sourceProvider
}

func (h *GetCollectionHandler) Handle() (r response.RPC) {
	sources, e := h.sourcesProvider.findSources()
	if e != nil {
		return e.AsRPCResponse()
	}

	allChassisCollection := sresponse.NewChassisCollection()
	for _, s := range sources {
		r, e := s.read()
		if e != nil {
			return e.AsRPCResponse()
		}
		for _, m := range r {
			allChassisCollection.AddMember(m)
		}
	}

	initializeRPCResponse(&r, allChassisCollection)
	return
}

type sourceProvider interface {
	findSources() ([]source, sresponse.Error)
}

type sourceProviderImpl struct {
	getPluginConfig func(pluginID string) (smodel.Plugin, *errors.Error)
	getAllKeys      func(table string) ([]string, error)
}

func (c *sourceProviderImpl) findSources() ([]source, sresponse.Error) {
	sources := []source{&managedChassisProvider{c.getAllKeys}}

	pluginConf, dberr := c.getPluginConfig("URP")
	if dberr != nil {
		if dberr.ErrNo() == errors.DBKeyNotFound {
			return sources, nil
		}
		return nil, &sresponse.RPCErrorWrapper{
			RPC: common.GeneralError(http.StatusInternalServerError, response.InternalError, dberr.Error(), nil, nil),
		}
	}

	sources = append(sources, &unmanagedChassisProvider{pluginConf: &pluginConf})
	return sources, nil
}

type source interface {
	read() ([]dmtf.Link, sresponse.Error)
}

type managedChassisProvider struct {
	inMemoryKeysProvider func(table string) ([]string, error)
}

func (m *managedChassisProvider) read() ([]dmtf.Link, sresponse.Error) {
	keys, e := m.inMemoryKeysProvider("Chassis")
	if e != nil {
		log.Printf("error getting all keys of ChassisCollection table : %v", e)
		return nil, &sresponse.RPCErrorWrapper{
			RPC: common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil),
		}
	}
	var r []dmtf.Link
	for _, key := range keys {
		r = append(r, dmtf.Link{Oid: key})
	}
	return r, nil

}

type unmanagedChassisProvider struct {
	pluginConf *smodel.Plugin
}

func (u unmanagedChassisProvider) read() ([]dmtf.Link, sresponse.Error) {
	r, e := plugin.NewClient(*u.pluginConf).Get("/ODIM/v1/Chassis")
	if e != nil {
		return nil, e
	}

	c := new(sresponse.Collection)
	if e := r.JSON(c); e != nil {
		return nil, &sresponse.RPCErrorWrapper{
			RPC: common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil),
		}
	}
	return c.Members, nil
}

func initializeRPCResponse(target *response.RPC, body sresponse.Collection) {
	target.StatusMessage = response.Success
	target.Body = body
	target.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	target.StatusCode = http.StatusOK
}
