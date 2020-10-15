package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-plugin-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
)

type ClientFactory func(name string) (Client, *errors.Error)

func NewClientFactory(t *config.URLTranslation) ClientFactory {
	return func(name string) (Client, *errors.Error) {
		pc, e := smodel.GetPluginData(name)
		if e != nil {
			return nil, e
		}
		return &client{plugin: pc, translator: &uriTranslator{t}}, nil
	}
}

type Client interface {
	Get(uri string) response.RPC
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

func (c *client) Get(uri string) response.RPC {
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
			log.Println("WARNING: ", "Cannot decode CommonError: ", err)
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
