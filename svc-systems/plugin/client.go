package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-plugin-rest-client/pmbhandle"
	smodel "github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
)

type ClientFactory func(name string) (Client, *errors.Error)

func ClientCreator(name string) (Client, *errors.Error) {
	c, e := smodel.GetPluginData(name)
	if e != nil {
		return nil, e
	}
	return NewClient(c), nil
}

func NewClient(plugin smodel.Plugin) *client {
	return &client{plugin: plugin}
}

type Response interface {
	JSON(t interface{}) error
	AsRPCResponse() (r response.RPC)
}

type pluginResponse struct {
	*http.Response
}

func (p *pluginResponse) JSON(t interface{}) error {
	bb := new(bytes.Buffer)
	io.Copy(bb, p.Body)
	bodyAsString := strings.Replace(bb.String(), "/ODIM/", "/redfish/", -1)
	return json.Unmarshal([]byte(bodyAsString), t)
}

func (p *pluginResponse) AsRPCResponse() (r response.RPC) {
	r.StatusCode, r.Body, r.Header = jsonResponseWriter(p.Response, func(toBeTransformed string) string {
		return strings.Replace(toBeTransformed, "/ODIM/", "/redfish/", -1)
	})
	return
}

func jsonResponseWriter(res *http.Response, bodyTransformers ...func(string) string) (statusCode int32, body []byte, headers map[string]string) {
	jsonBodyBytes := new(json.RawMessage)
	err := json.NewDecoder(res.Body).Decode(jsonBodyBytes)
	if err != nil {
		ge := common.GeneralError(
			http.StatusInternalServerError,
			response.GeneralError,
			fmt.Sprintf("Cannot read response: %v", err), nil, nil)

		return ge.StatusCode, jsonMarshal(ge.Body), ge.Header
	}

	bodyToBeTransformed := string(*jsonBodyBytes)
	for _, t := range bodyTransformers {
		bodyToBeTransformed = t(bodyToBeTransformed)
	}

	headersToBeReturned := map[string]string{}
	for k := range res.Header {
		if hasToBeSkipped(k) {
			continue
		}
		for _, t := range bodyTransformers {
			headersToBeReturned[k] = t(res.Header.Get(k))
		}
	}

	return int32(res.StatusCode), []byte(bodyToBeTransformed), headersToBeReturned
}

func jsonMarshal(input interface{}) []byte {
	if bytes, alreadyBytes := input.([]byte); alreadyBytes {
		return bytes
	}
	bytes, err := json.Marshal(input)
	if err != nil {
		log.Println("error in unmarshalling response object from util-libs", err.Error())
	}
	return bytes
}

func hasToBeSkipped(header string) bool {
	return header == "Content-Length"
}

type Client interface {
	Get(uri string) (Response, sresponse.Error)
	Post(uri string, body interface{}) (Response, sresponse.Error)
}

type client struct {
	plugin smodel.Plugin
}

func (c *client) Post(uri string, body interface{}) (Response, sresponse.Error) {
	url := fmt.Sprintf("https://%s:%s%s", c.plugin.IP, c.plugin.Port, uri)
	resp, err := pmbhandle.ContactPlugin(url, http.MethodPost, "", "", body, map[string]string{
		"UserName": c.plugin.Username,
		"Password": string(c.plugin.Password),
	})
	if err != nil {
		return nil, &sresponse.UnknownErrorWrapper{Error: err, StatusCode: resp.StatusCode}
	}
	if !is2xx(resp.StatusCode) {
		return nil, extractError(resp)
	}
	return &pluginResponse{resp}, nil
}

func (c *client) Get(uri string) (Response, sresponse.Error) {
	url := fmt.Sprintf("https://%s:%s%s", c.plugin.IP, c.plugin.Port, uri)
	resp, err := pmbhandle.ContactPlugin(url, http.MethodGet, "", "", nil, map[string]string{
		"UserName": c.plugin.Username,
		"Password": string(c.plugin.Password),
	})
	if err != nil {
		if resp != nil {
			return nil, &sresponse.UnknownErrorWrapper{Error: err, StatusCode: resp.StatusCode}
		}
		return nil, &sresponse.UnknownErrorWrapper{Error: err, StatusCode: http.StatusInternalServerError}
	}
	if !is2xx(resp.StatusCode) {
		return nil, extractError(resp)
	}
	return &pluginResponse{resp}, nil
}

func extractError(resp *http.Response) sresponse.Error {
	ce := new(response.CommonError)
	if e := (&pluginResponse{resp}).JSON(ce); e != nil {
		return &sresponse.UnknownErrorWrapper{StatusCode: resp.StatusCode, Error: e}
	}

	r := response.RPC{}
	r.StatusCode, r.Body, r.Header = jsonResponseWriter(resp, func(toBeTransformed string) string {
		return strings.Replace(toBeTransformed, "/ODIM/", "/redfish/", -1)
	})

	return &sresponse.RPCErrorWrapper{RPC: r}
}

func is2xx(status int) bool {
	return status/100 == 2
}
