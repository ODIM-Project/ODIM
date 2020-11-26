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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func NewRedfishClient(baseURL string) RedfishClient {
	return &redfishClient{
		NewClient(baseURL),
	}
}

type RedfishClient interface {
	Get(uri string, target interface{}) *CommonError
}

type redfishClient struct {
	c Client
}

func (r *redfishClient) Get(uri string, target interface{}) *CommonError {
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

func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		httpc: &http.Client{
			//todo:  configure tls transport
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		requestDecorators: []requestDecorator{
			&basicAuth{
				u: "admin",
				p: "Od!m12$4",
			},
			&requestTranslator{
				old: "/ODIM/",
				new: "/redfish/",
			},
		},
		responseDecorators: []responseDecorator{
			&responseBodyTranslator{
				old: "/redfish/",
				new: "/ODIM/",
			},
		},
	}
}

type requestDecorator interface {
	decorate(r *http.Request) error
}

type responseDecorator interface {
	decorate(response *http.Response) error
}

type Client interface {
	Get(uri string) (*http.Response, error)
	Post(uri string, body []byte) (*http.Response, error)
}

type httpClient struct {
	baseURL            string
	httpc              *http.Client
	requestDecorators  []requestDecorator
	responseDecorators []responseDecorator
}

func (h *httpClient) createURL(uri string) (*url.URL, error) {
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

func (h *httpClient) decorateResponse(r *http.Response) error {
	for _, d := range h.responseDecorators {
		e := d.decorate(r)
		if e != nil {
			return e
		}
	}
	return nil
}

func (h *httpClient) decorateRequest(r *http.Request) error {
	for _, d := range h.requestDecorators {
		e := d.decorate(r)
		if e != nil {
			return e
		}
	}
	return nil
}

func (h *httpClient) Get(uri string) (*http.Response, error) {
	requestUrl, err := h.createURL(uri)
	if err != nil {
		return nil, err
	}
	req := http.Request{
		Method: http.MethodGet,
		URL:    requestUrl,
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

func (h *httpClient) Post(uri string, bodyBytes []byte) (*http.Response, error) {
	requestUrl, err := h.createURL(uri)
	if err != nil {
		return nil, err
	}

	body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	req := http.Request{
		Method: http.MethodPost,
		URL:    requestUrl,
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

type requestTranslator struct {
	old, new string
}

func (r *requestTranslator) decorate(req *http.Request) error {
	for k, hvs := range req.Header {
		var translated []string
		for _, v := range hvs {
			translated = append(translated, strings.Replace(v, r.old, r.new, -1))
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
	translatedBody := strings.Replace(string(bodyBytes), r.old, r.new, -1)
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(translatedBody)))
	return nil
}

type responseBodyTranslator struct {
	old, new string
}

func (r *responseBodyTranslator) decorate(resp *http.Response) error {
	dec := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	tempTarget := new(json.RawMessage)

	err := dec.Decode(tempTarget)
	if err != nil {
		return err
	}

	newBody := strings.Replace(string(*tempTarget), r.old, r.new, -1)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(newBody)))
	return nil
}
