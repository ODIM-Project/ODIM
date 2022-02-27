//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

// Packahe lphandler ...
package lphandler

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	testhttp "net/http/httptest"
	"strings"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-lenovo/config"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpresponse"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

var (
	hostCA = []byte(`-----BEGIN CERTIFICATE-----
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
6ATq
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFsjCCA5qgAwIBAgIUNN0LnDt13/3nM1GipqN9Pku+ZdEwDQYJKoZIhvcNAQEN
BQAwUjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQ8wDQYDVQQKDAZPRElNUkExEDAOBgNVBAMMB1Jvb3QgQ0EwHhcNMjExMDI2
MDkzNDE5WhcNMjMwMzEwMDkzNDE5WjBBMQswCQYDVQQGEwJVUzELMAkGA1UECAwC
Q0ExDzANBgNVBAoMBk9ESU1SQTEUMBIGA1UEAwwLU2VydmVyIENlcnQwggIiMA0G
CSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQC1vI9gNVOQnpoSsm2Bonk41xbcJ1CR
5tj6wswEOBaDNXHtYOq85vMT29ou4r5DYanN05Fhd+uJC22Gf4nbGoLkkt1AuIIb
d2U3PY9YJvNfQ2R2GAvTSBTuB5wZ6Co7JEX9kRlVbCcQrF1TJjTI0MgSsHf5U4rg
pr70YN2WecpUieQ3GG+x94NpupwwGJWeh8JIgDxAdh8m/cebOT1VO5xxtPpcd8n4
TxlroPLbgz3UFyMiElmih3+w+L89pOuuPgivBdIVKNsYLYICFw9s09cgwlsi6qhy
Esa0yMjLARuaxz2VFVPwDH9+JmGHWWIUfjXnQ3WixGTrpAI1iboN24DrFtjdIRfP
RgVpuDP7HSN2/CNIXxQ2RQF0fB/z4g1KJS56UgGFCvgDl2BG3aOj+tlQofRvgxT9
jdBBkeiCLgixkVs4nu73rB4Bq4RHRosrpt4kT4nmioMuUcLJRsGB1zt0+xQ+DR6L
vsnzbD4bS3fYADuTzDzXjBem9Z2RYUX9wUtXeiOqqNXepb/AeHmZLzUlVlSu5h6c
PxR9adSVw8NBSR1JXn0hqLFjOVEDk3wMicQnzRXYUrpJgh0tUqu6gHs7DHUsYQPj
qzXxMZbpJyMZSJirqu43VkyVhyIswnzGErwDs6gEM9qUlrQ7L4Zcpy+Af9SqpZPG
aB6p6/m1a4T7XwIDAQABo4GQMIGNMAkGA1UdEwQCMAAwCwYDVR0PBAQDAgXgMB0G
A1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4EFgQUdWXGZkyYvxcX
A8WCUJ9Df/KA4w0wHwYDVR0jBBgwFoAU2x5lPlMmcsAOk5p96W+fpo2kPeYwFAYD
VR0RBA0wC4IJbG9jYWxob3N0MA0GCSqGSIb3DQEBDQUAA4ICAQBpMuSlw82RHEFJ
YOFy8iPM8bxTThrZRzCpubePKdRPOKEor8WXuHyiKNhbaaPdvgHgPvTBYjFGOsJn
gjWbUc8WwhyCj/1U+m/xSiOTzgbs2sE5Qa4/v4s7Xai3GRAeTC3Y+QDTpaZQWkRT
mSoZQbELzjPYtyxcblUUzLjjupjjld5vlNKKoPyl6WdtRrnkkvgABEFAysgWyW1l
yKO/xfe6sBS7H/UWqZOZZQQIy8svmCzNEPXsr1UwjQiJB74KgZrRuomm/FaA+fVW
jpAQ+LsN6Dhe9dv8J4e4MgHA9uaiFRwI1z+xeQnUHXJUbWHQNWpbI62ZTH3LftFT
+f54Obh3DdVyt+zugBZoIHUaAB+eKcTcruVEwDcDRFS0ObyqBkIej2ED+1iTFw3l
VoASse8XMdn734a3Uso65yUoTCj/Z54/N09m24ZRXyoyahgWwl0Kc+2lINcEEj7l
zAAjO94sWIWvFnA1c2rgx9zFQd7PU2zIHOP/A3tse3yyoszvS38eiTZkvP2Tj0mu
AzjfunVr6C7YjfYWnAa3tGYTpSlKnKVnCl7vIfNtthC2i3EcXACQcCecQOBpp6A2
7pCccBHHUvPMuETjEnCC5wSrSYA0MSHOreS0to7w+NGAc5/YjINN+7Zo6JooALNk
KoQRsFQ0udOGCu+Z5uZfMr0nre618w==
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
)

type mockHandlerFunc func(string, string, string, http.ResponseWriter)

func startTestServer(handler mockHandlerFunc) *testhttp.Server {
	// create a listener with the desired port.
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	ts := testhttp.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/redfish/v1" {
				handler("", "", r.URL.Path, w)
			} else {
				auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
				payload, _ := base64.StdEncoding.DecodeString(auth[1])
				pair := strings.SplitN(string(payload), ":", 2)
				handler(pair[0], pair[1], r.URL.Path, w)
			}
		}))

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener = l
	tlsConfig := &tls.Config{}

	cert, err := tls.X509KeyPair(hostCert, hostPrivKey)
	if err != nil {
		log.Fatal("Failed to load key pair: " + err.Error())
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.BuildNameToCertificate()

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(hostCA) {
		log.Fatal("Failed to load CA certificate")
	}
	tlsConfig.RootCAs = capool
	tlsConfig.ClientCAs = capool

	ts.TLS = tlsConfig
	ts.Config.TLSConfig = tlsConfig
	return ts
}

func mockDeviceHandler(username, password, url string, w http.ResponseWriter) {
	resp, err := mockChangeBiosSettings(username, url)
	if err != nil && resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(resp.StatusCode)
	return
}

func mockChangeBiosSettings(username, url string) (*http.Response, error) {
	if url == "/ODIM/v1/Systems/1" && username == "admin" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Success")),
		}, nil
	}
	if url == "/ODIM/v1/Systems/1" && username != "admin" {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Failed")),
		}, fmt.Errorf("Error")
	}
	if url == "/ODIM/v1/Systems/1/bios/settings" && username == "admin" {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Success")),
		}, nil
	}
	if url == "/ODIM/v1/Systems/1/bios/settings" && username != "admin" {
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Failed")),
		}, fmt.Errorf("Error")
	}
	return nil, fmt.Errorf("Error")
}

func TestChangeBootOrderSettings(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockDeviceHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Patch("/Systems/{id}", ChangeSettings)

	lpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}
	//Unit Test for success scenario
	e.PATCH("/redfish/v1/Systems/1").WithJSON(requestBody).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.PATCH("/redfish/v1/Systems/1").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	//unittest for bad request scenario: given device details are wrong
	requestBody1 := "requestbody"
	e.PATCH("/redfish/v1/Systems/1").WithJSON(requestBody1).Expect().Status(http.StatusBadRequest)
}

func TestChangeBiosSettings(t *testing.T) {
	config.SetUpMockConfig(t)

	deviceHost := "localhost"
	devicePort := "1234"
	ts := startTestServer(mockDeviceHandler)
	// Start the server.
	ts.StartTLS()
	defer ts.Close()

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")

	redfishRoutes.Patch("/Systems/{id}/bios/settings", ChangeSettings)

	lpresponse.PluginToken = "token"

	e := httptest.New(t, mockApp)

	requestBody := map[string]interface{}{
		"ManagerAddress": fmt.Sprintf("%s:%s", deviceHost, devicePort),
		"UserName":       "admin",
		"Password":       []byte("P@$$w0rd"),
	}
	//Unit Test for success scenario
	e.PATCH("/redfish/v1/Systems/1/bios/settings").WithJSON(requestBody).Expect().Status(http.StatusOK)

	//Case for invalid token
	e.PATCH("/redfish/v1/Systems/1/bios/settings").WithHeader("X-Auth-Token", "token").WithJSON(requestBody).Expect().Status(http.StatusUnauthorized)

	//unittest for bad request scenario: given device details are wrong
	requestBody1 := "requestbody"
	e.PATCH("/redfish/v1/Systems/1/bios/settings").WithJSON(requestBody1).Expect().Status(http.StatusBadRequest)
}
