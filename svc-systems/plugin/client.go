/*
 * Copyright (c) 2020 Intel Corporation
 * (C) Copyright [2020] Hewlett Packard Enterprise Development LP
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

package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type ClientFactory func(name string) (Client, *errors.Error)

func NewClientFactory(t *config.URLTranslation) ClientFactory {

	pluginClientCreator := func(plugin *smodel.Plugin) Client {
		return &client{plugin: *plugin, translator: &uriTranslator{t}}
	}

	return func(name string) (Client, *errors.Error) {
		if strings.HasSuffix(name, "*") {
			plugins, err := findAllPlugins(name)
			if err != nil {
				return nil, errors.PackError(errors.UndefinedErrorType, err)
			}

			return &multiTargetClient{
				call:         ReturnFirstResponse(&CallConfig{}),
				createClient: pluginClientCreator,
				targets:      plugins,
			}, nil
		} else {
			plugin, e := smodel.GetPluginData(name)
			if e != nil {
				return nil, e
			}
			return pluginClientCreator(&plugin), nil
		}
	}
}

type multiTargetClient struct {
	call         *CallConfig
	createClient func(plugin *smodel.Plugin) Client
	targets      []*smodel.Plugin
}

type CallConfig struct {
	collector Collector
}

type CallOption func(*CallConfig)

type Collector interface {
	Collect(response.RPC) error
	GetResult() response.RPC
}

type returnFirst struct {
	ReqURI string
	resp   *response.RPC
}

func (c *returnFirst) Collect(r response.RPC) error {
	c.resp = &r
	return nil
}

func (c *returnFirst) GetResult() response.RPC {
	if c.resp == nil {
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", c.ReqURI}, nil)
	}
	return *c.resp
}

func ReturnFirstResponse(c *CallConfig) *CallConfig {
	c.collector = new(returnFirst)
	return c
}

func AggregateResults(c *CallConfig) {
	c.collector = new(collectCollectionMembers)
}

type collectCollectionMembers struct {
	collection sresponse.Collection
	errored    bool
}

func (c *collectCollectionMembers) Collect(r response.RPC) error {
	if is2xx(int(r.StatusCode)) {
		collection := new(sresponse.Collection)
		err := json.Unmarshal(r.Body.([]byte), collection)
		if err != nil {
			return err
		}
		for _, m := range collection.Members {
			c.collection.AddMember(m)
		}
	} else {
		return fmt.Errorf("provided response has unexpecte error code")
	}
	return nil
}

func (c *collectCollectionMembers) GetResult() response.RPC {
	collectionAsBytes, err := json.Marshal(c.collection)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("Unexpected error: %v", err), nil, nil)
	}
	return response.RPC{
		StatusCode: http.StatusOK,
		Body:       collectionAsBytes,
	}
}

func (m *multiTargetClient) Get(uri string, opts ...CallOption) response.RPC {
	for _, opt := range opts {
		opt(m.call)
	}
	for _, target := range m.targets {
		client := m.createClient(target)
		resp := client.Get(uri)
		err := m.call.collector.Collect(resp)
		if err != nil {
			log.Warn("execution of GET " + uri + " on " + target.ID + " plugin returned non 2xx status code; " + convertToString(resp.Body))
		}
	}
	// Checking whether the struct passed as the interface has a ReqURI field.
	// If it is a collection request then the ReqURI field won't be present.
	// Only struct associating with individual resource collection requires the ReqURI field.
	field := reflect.ValueOf(m.call.collector).Elem().FieldByName("ReqURI")
	if field.IsValid() {
		chassisID := getChassisID(uri)
		field.SetString(chassisID)
	}
	return m.call.collector.GetResult()
}

func getChassisID(uri string) string {
	parts := strings.Split(uri, "/")
	return parts[len(parts)-1]
}

func (m *multiTargetClient) Post(uri string, body *json.RawMessage) response.RPC {
	// TODO: Implement this
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}

func (m *multiTargetClient) Patch(uri string, body *json.RawMessage) response.RPC {
	for _, target := range m.targets {
		client := m.createClient(target)
		log.Info("Request received to patch chassis to rack, URI: ", uri)
		resp := client.Patch(uri, body)
		switch {
		case resp.StatusCode == http.StatusNotFound:
			continue
		case is2xx(int(resp.StatusCode)):
			return resp
		case is4xx(int(resp.StatusCode)):
			return resp
		default:
			log.Warn("execution of PATCH " + uri + " on " + target.ID + " plugin returned non 2xx status code; " + convertToString(resp.Body))
		}
	}
	return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", uri}, nil)
}

func (m *multiTargetClient) Delete(uri string) response.RPC {
	for _, target := range m.targets {
		client := m.createClient(target)
		resp := client.Delete(uri)
		switch {
		case resp.StatusCode == http.StatusNotFound:
			continue
		case is2xx(int(resp.StatusCode)):
			return resp
		case is4xx(int(resp.StatusCode)):
			return resp
		default:
			log.Warn("execution of DELETE " + uri + " on " + target.ID + " plugin returned non 2xx status code; " + convertToString(resp.Body))
		}
	}
	return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, "", []interface{}{"Chassis", uri}, nil)
}

type Client interface {
	Get(uri string, opts ...CallOption) response.RPC
	Post(uri string, body *json.RawMessage) response.RPC
	Patch(uri string, body *json.RawMessage) response.RPC
	Delete(uri string) response.RPC
}

type client struct {
	translator *uriTranslator
	plugin     smodel.Plugin
}

func (c *client) Delete(uri string) response.RPC {
	url := fmt.Sprintf("https://%s:%s%s", c.plugin.IP, c.plugin.Port, uri)
	url = c.translator.toSouthbound(url)
	resp, err := pmbhandle.ContactPlugin(url, http.MethodDelete, "", "", nil, map[string]string{
		"UserName": c.plugin.Username,
		"Password": string(c.plugin.Password),
	})
	return c.extractResp(resp, err)
}

func (c *client) Post(uri string, body *json.RawMessage) response.RPC {
	url := fmt.Sprintf("https://%s:%s%s", c.plugin.IP, c.plugin.Port, uri)
	url = c.translator.toSouthbound(url)
	*body = json.RawMessage(c.translator.toSouthbound(string(*body)))
	resp, err := pmbhandle.ContactPlugin(url, http.MethodPost, "", "", body, map[string]string{
		"UserName": c.plugin.Username,
		"Password": string(c.plugin.Password),
	})
	return c.extractResp(resp, err)
}

func (c *client) Patch(uri string, body *json.RawMessage) response.RPC {
	url := fmt.Sprintf("https://%s:%s%s", c.plugin.IP, c.plugin.Port, uri)
	url = c.translator.toSouthbound(url)
	*body = json.RawMessage(c.translator.toSouthbound(string(*body)))
	resp, err := pmbhandle.ContactPlugin(url, http.MethodPatch, "", "", body, map[string]string{
		"UserName": c.plugin.Username,
		"Password": string(c.plugin.Password),
	})
	return c.extractResp(resp, err)
}

func (c *client) Get(uri string, _ ...CallOption) response.RPC {
	url := fmt.Sprintf("https://%s:%s%s", c.plugin.IP, c.plugin.Port, uri)
	url = c.translator.toSouthbound(url)
	resp, err := pmbhandle.ContactPlugin(url, http.MethodGet, "", "", nil, map[string]string{
		"UserName": c.plugin.Username,
		"Password": string(c.plugin.Password),
	})
	return c.extractResp(resp, err)
}

func (c *client) extractResp(httpResponse *http.Response, err error) response.RPC {
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
	}

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, fmt.Sprintf("Cannot read response body: %v", err), nil, nil)
	}

	if !is2xx(httpResponse.StatusCode) {
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()

		ce := new(response.CommonError)
		err := dec.Decode(ce)
		if err != nil {
			log.Error("Cannot decode CommonError: " + err.Error())
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, string(body), nil, nil)
		}
	}
	return createRPCResponse(httpResponse.StatusCode, httpResponse.Header, body, c.translator.toNorthbound)
}

func createRPCResponse(s int, h http.Header, b []byte, translator func(string) string) response.RPC {
	r := response.RPC{
		StatusCode: int32(s),
		Header:     map[string]string{},
	}

	for k := range h {
		if hasToBeSkipped(k) {
			continue
		}
		r.Header[k] = translator(h.Get(k))
	}

	if b != nil && len(b) != 0 {
		bodyToBeTransformed := string(b)
		bodyToBeTransformed = translator(bodyToBeTransformed)
		r.Body = []byte(bodyToBeTransformed)
	}
	return r
}

func is2xx(status int) bool {
	return status/100 == 2
}

func is4xx(status int) bool {
	return status/100 == 4
}

func hasToBeSkipped(header string) bool {
	return header == "Content-Length"
}

type uriTranslator struct {
	dictionaries *config.URLTranslation
}

func (u *uriTranslator) toSouthbound(data string) string {
	translated := data
	for k, v := range u.dictionaries.SouthBoundURL {
		translated = strings.Replace(translated, k, v, -1)
	}
	return translated
}

func (u *uriTranslator) toNorthbound(data string) string {
	translated := data
	for k, v := range u.dictionaries.NorthBoundURL {
		translated = strings.Replace(translated, k, v, -1)
	}
	return translated
}

func findAllPlugins(key string) (res []*smodel.Plugin, err error) {
	pluginsAsBytesSlice, err := smodel.FindAll("Plugin", key)
	if err != nil {
		return
	}

	for _, bytes := range pluginsAsBytesSlice {
		plugin := new(smodel.Plugin)
		err = json.Unmarshal(bytes, plugin)
		if err != nil {
			return nil, err
		}
		decryptedPass, err := common.DecryptWithPrivateKey(plugin.Password)
		if err != nil {
			return nil, errors.PackError(
				errors.DecryptionFailed,
				fmt.Sprintf("error: %s plugin password decryption failed: ", plugin.ID),
				err.Error(),
			)
		}
		plugin.Password = decryptedPass
		res = append(res, plugin)
	}
	return
}

func convertToString(data interface{}) string {
	byteData, err := json.Marshal(data)
	if err != nil {
		log.Error("converting interface to string type failed: " + err.Error())
		return ""
	}

	return string(byteData)
}
