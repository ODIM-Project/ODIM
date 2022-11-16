/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chassis

import (
	"encoding/json"
	"net/http"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
)

const collectionURL = "/redfish/v1/Chassis"

// NewGetCollectionHandler returns an instance of GetCollection struct
func NewGetCollectionHandler(
	pcf plugin.ClientFactory,
	imkp func(table string) ([]string, error)) *GetCollection {

	return &GetCollection{
		&sourceProviderImpl{
			pluginClientFactory: pcf,
			getAllKeys:          imkp,
			getFabricFactory:    getFabricFactory,
		},
	}
}

// GetCollection struct helps to get chassis collection information
type GetCollection struct {
	sourcesProvider sourceProvider
}

// Handle defines the operations which handle the RPC request-response for getting chassis collection information
func (h *GetCollection) Handle() (r response.RPC) {
	sources, e := h.sourcesProvider.findSources()
	if e != nil {
		return *e
	}

	allChassisCollection := sresponse.NewChassisCollection()
	for _, s := range sources {
		r, e := s.read()
		if e != nil {
			return *e
		}
		for _, m := range r {
			allChassisCollection.AddMember(m)
		}
	}

	h.sourcesProvider.findFabricChassis(&allChassisCollection)

	initializeRPCResponse(&r, allChassisCollection)
	return
}

type sourceProvider interface {
	findSources() ([]source, *response.RPC)
	findFabricChassis(c *sresponse.Collection)
}

type sourceProviderImpl struct {
	pluginClientFactory plugin.ClientFactory
	getAllKeys          func(table string) ([]string, error)
	getFabricFactory    func(collection *sresponse.Collection) *fabricFactory
}

func (c *sourceProviderImpl) findSources() ([]source, *response.RPC) {
	sources := []source{&managedChassisProvider{c.getAllKeys}}

	pc, dberr := c.pluginClientFactory("URP*")
	if dberr != nil {
		if dberr.ErrNo() == errors.DBKeyNotFound {
			return sources, nil
		}
		ge := common.GeneralError(http.StatusInternalServerError, response.InternalError, dberr.Error(), nil, nil)
		return nil, &ge
	}

	sources = append(sources, &unmanagedChassisProvider{c: pc})
	return sources, nil
}

type source interface {
	read() ([]dmtf.Link, *response.RPC)
}

type managedChassisProvider struct {
	inMemoryKeysProvider func(table string) ([]string, error)
}

func (m *managedChassisProvider) read() ([]dmtf.Link, *response.RPC) {
	keys, e := m.inMemoryKeysProvider("Chassis")
	if e != nil {
		l.Log.Error("while getting all keys of ChassisCollection table, got " + e.Error())
		ge := common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
		return nil, &ge
	}
	var r []dmtf.Link
	for _, key := range keys {
		r = append(r, dmtf.Link{Oid: key})
	}
	return r, nil

}

type unmanagedChassisProvider struct {
	c plugin.Client
}

func (u unmanagedChassisProvider) read() ([]dmtf.Link, *response.RPC) {
	r := u.c.Get(collectionURL, plugin.AggregateResults)
	if r.StatusCode != http.StatusOK {
		return nil, &r
	}

	c := new(sresponse.Collection)
	if e := json.Unmarshal(r.Body.([]byte), c); e != nil {
		ge := common.GeneralError(http.StatusInternalServerError, response.InternalError, e.Error(), nil, nil)
		return nil, &ge
	}
	return c.Members, nil
}

func initializeRPCResponse(target *response.RPC, body interface{}) {
	target.StatusMessage = response.Success
	target.Body = body
	target.StatusCode = http.StatusOK
}
