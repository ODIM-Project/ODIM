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

package redfish

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
)

// HTTPClient is an client for communication with redfish service.
// In case of URP, it it used for calling ODIMRA's REST API.
type HTTPClient struct {
	baseURL            string
	httpc              *http.Client
	requestDecorators  []requestDecorator
	responseDecorators []responseDecorator
}

// HTTPClientOption is interface of configuration option for HTTPClient
type HTTPClientOption func(rc *HTTPClient)

// BaseURL configures HTTPClient by setting base URL
func BaseURL(baseURL string) HTTPClientOption {
	return func(rc *HTTPClient) {
		rc.baseURL = baseURL
	}
}

// InsecureSkipVerifyTransport configures HTTPClient with insecure transport(skips certificate verification)
// It is intended to be used only in tests!!!
func InsecureSkipVerifyTransport(c *HTTPClient) {
	c.httpc.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

// HTTPTransport configures HTTPClient by setting secure TLS transport
func HTTPTransport(c *config.PluginConfig) HTTPClientOption {
	return func(rc *HTTPClient) {
		caCert, err := ioutil.ReadFile(c.PKIRootCAPath)
		if err != nil {
			panic(err)
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(caCert)
		clientCert, err := tls.LoadX509KeyPair(c.PKICertificatePath, c.PKIPrivateKeyPath)
		if err != nil {
			panic(err)
		}
		tlsConf := tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{clientCert},
		}
		tlsTransport := http.Transport{
			TLSClientConfig: &tlsConf,
		}
		rc.httpc.Transport = &tlsTransport
	}
}

// NewHTTPClient creates new instance of HTTPClient.
// Returned clients is configured regarding to provided configuration options(opts)
func NewHTTPClient(opts ...HTTPClientOption) *HTTPClient {
	c := &HTTPClient{
		httpc: &http.Client{},
		requestDecorators: []requestDecorator{
			&basicAuth{
				u: "admin",
				p: "Od!m12$4",
			},
			&requestTranslator{},
		},
		responseDecorators: []responseDecorator{
			&responseBodyTranslator{},
		},
	}

	for _, o := range opts {
		o(c)
	}
	return c
}

type requestDecorator interface {
	decorate(r *http.Request) error
}

type responseDecorator interface {
	decorate(response *http.Response) error
}

func (h *HTTPClient) createURL(uri string) (*url.URL, error) {
	path, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	baseURL, err := url.Parse(h.baseURL)
	if err != nil {
		return nil, err
	}
	return baseURL.ResolveReference(path), nil
}

func (h *HTTPClient) decorateResponse(r *http.Response) error {
	for _, d := range h.responseDecorators {
		e := d.decorate(r)
		if e != nil {
			return e
		}
	}
	return nil
}

func (h *HTTPClient) decorateRequest(r *http.Request) error {
	for _, d := range h.requestDecorators {
		e := d.decorate(r)
		if e != nil {
			return e
		}
	}
	return nil
}

// Get executes GET operation against requested endpoint
func (h *HTTPClient) Get(uri string) (*http.Response, error) {
	requestedURL, err := h.createURL(uri)
	if err != nil {
		return nil, err
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    requestedURL,
		Header: http.Header{},
	}

	err = h.decorateRequest(&req)
	if err != nil {
		return nil, err
	}

	resp, err := h.httpc.Do(&req)
	if err != nil {
		return resp, err
	}

	err = h.decorateResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("cannot decorate response: %v", err)
	}
	return resp, nil
}

// Post executes POST operation against requested endpoint
// Executed POST request carries provided body.
func (h *HTTPClient) Post(uri string, bodyBytes []byte) (*http.Response, error) {
	requestedURL, err := h.createURL(uri)
	if err != nil {
		return nil, err
	}

	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	req := http.Request{
		Method: http.MethodPost,
		URL:    requestedURL,
		Body:   body,
		Header: http.Header{},
	}

	err = h.decorateRequest(&req)
	if err != nil {
		return nil, fmt.Errorf("cannot decorate request: %v", err)
	}

	resp, err := h.httpc.Do(&req)
	if err != nil {
		return resp, err
	}

	err = h.decorateResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("cannot decorate response: %v", err)
	}
	return resp, nil
}

type basicAuth struct {
	u, p string
}

func (b *basicAuth) decorate(r *http.Request) error {
	r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(([]byte)(b.u+":"+b.p)))
	return nil
}

type requestTranslator struct{}

func (r *requestTranslator) decorate(req *http.Request) error {
	for k, hvs := range req.Header {
		var translated []string
		for _, v := range hvs {
			translated = append(translated, Translator.ODIMToRedfish(v))
		}
		req.Header[k] = translated
	}

	if req.Body == nil {
		return nil
	}
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	translatedBody := Translator.ODIMToRedfish(string(bodyBytes))
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(translatedBody)))
	return nil
}

type responseBodyTranslator struct{}

func (r *responseBodyTranslator) decorate(resp *http.Response) error {
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	tempTarget := new(json.RawMessage)

	err := dec.Decode(tempTarget)
	if err != nil {
		return err
	}

	newBody := Translator.RedfishToODIM(string(*tempTarget))
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(newBody)))
	return nil
}

// ResponseWrappingClient is a wrapper for HTTPClient. Intentionally it wraps raw `http.Response` into Redfish entity.
type ResponseWrappingClient struct {
	c *HTTPClient
}

// NewResponseWrappingClient creates new instance of ResponseWrappingClient
func NewResponseWrappingClient(httpClient *HTTPClient) *ResponseWrappingClient {
	return &ResponseWrappingClient{httpClient}
}

// Get executes GET operation, wraps response into `target`, in case of any error CommonError is returned.
func (r *ResponseWrappingClient) Get(uri string, target interface{}) *CommonError {
	rsp, err := r.c.Get(uri)
	if err != nil {
		ce := CreateError(GeneralError, err.Error())
		return &ce
	}
	if rsp.StatusCode != http.StatusOK {
		ce := CreateError(GeneralError, fmt.Sprintf("GET %s operation finished with status != 200", rsp.Request.URL.String()))
		return &ce
	}

	dec := json.NewDecoder(rsp.Body)
	defer rsp.Body.Close()
	err = dec.Decode(target)
	if err != nil {
		ce := CreateError(GeneralError, fmt.Sprintf("Cannot decode body of GET %s operation: %s", rsp.Request.URL.String(), err))
		return &ce
	}
	return nil
}
