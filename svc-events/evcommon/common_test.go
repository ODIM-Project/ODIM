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

// Package evcommon ...
package evcommon

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/stretchr/testify/assert"
)

var (
	hostCA = []byte(`-----BEGIN CERTIFICATE-----
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
9AjGBBg=
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIG1jCCBL6gAwIBAgIUYttkxU3D5hOJyBWaZHYhHUgRcfwwDQYJKoZIhvcNAQEL
BQAwcDELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRMwEQYDVQQHDApDYWxpZm9y
bmlhMQwwCgYDVQQKDANIUEUxGDAWBgNVBAsMD1RlbGNvIFNvbHV0aW9uczEXMBUG
A1UEAwwOT0RJTVJBX1JPT1RfQ0EwHhcNMjIwNTExMDcxMDEzWhcNMzIwNTA4MDcx
MDEzWjBwMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExEzARBgNVBAcMCkNhbGlm
b3JuaWExDDAKBgNVBAoMA0hQRTEYMBYGA1UECwwPVGVsY28gU29sdXRpb25zMRcw
FQYDVQQDDA5PRElNUkFfU1ZDX0NSVDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBALvx70x1niMxDMRHxUSJZDfQYtNQuCySjYZ/CnCwuJTAYQ+2+ogC0elU
ZBGzbB3m+4ca4MttEilinzsVdHATvOx8zqjqAJNRZWd+JsBpp6Y1sOMrQHuAHCFx
eLK1EkYLjXifq7ScJwQrm6MNxDL/Wa9pa5f3j2sxKnqf7SrPBqzrXKsaXMhOtoa+
vOTZkCRCeNn+XjK7MJRi0OCEjhKCppzSt1lp6xURm/K8ELqowoKSqLkEZIxcgGxD
u7U9lrf6CL4+TODaZ9qlQ7xUmggffOONnvlxRYNUPEINSLgWuFU+F94x8DAtNvtm
6iOOXD+AzJ1IwhChQnbx5iKCtyDcJoMWh0xEiapMBJPeDAyD6D2p7IdbqphIqa/9
gFDN7jw9GBmWaOlT8z8W3mojcEdekRpwCua9hKKwFbUYs7FItyjn3hPmBFA5VqLc
Cov8cN7PJJNPrcOhxTL9UWN2AmMlpaiQE320wXwjduu74OsuZomDAgDiecKiZUrm
pCqv8fJ20W2RkSvzWM754m+1mgqOAfPzggOOHKiIA4oFc2o4R+ZF4pmLtuAEcNvx
iIMPETIP0qGot2oZp0tlNe13xITxj9KkrBQdlNs8+aAVdxi/n75rfrDEArcaDue7
+hXbHJ0ZBowRCYfp7YVwwnvSnva4db729QOE4FDDx9xNyyl8PxdnAgMBAAGjggFm
MIIBYjAdBgNVHQ4EFgQUrOVT19Um6aIJPXDFIKocnyrhALkwga0GA1UdIwSBpTCB
ooAUO6rgMJaE7Ug1L12DTP7Zh9TX8pOhdKRyMHAxCzAJBgNVBAYTAlVTMQswCQYD
VQQIDAJDQTETMBEGA1UEBwwKQ2FsaWZvcm5pYTEMMAoGA1UECgwDSFBFMRgwFgYD
VQQLDA9UZWxjbyBTb2x1dGlvbnMxFzAVBgNVBAMMDk9ESU1SQV9ST09UX0NBghRw
rQR8cLxuujIPUdZRqxJKyZGPKzAOBgNVHQ8BAf8EBAMCBeAwHQYDVR0lBBYwFAYI
KwYBBQUHAwIGCCsGAQUFBwMBMGIGA1UdEQRbMFmCC2NsdXN0ZXJub2Rlgg5yZWRp
cy1pbm1lbW9yeYIMcmVkaXMtb25kaXNrgglpbG9wbHVnaW6CEGlsb3BsdWdpbi1l
dmVudHOCCWxvY2FsaG9zdIcEfwAAATANBgkqhkiG9w0BAQsFAAOCAgEAIKkIlcZg
tC2RiNJnvOnwDpLin0Ygy5BZbHVizo82RFAhHI2UPSpSNaRi6/c9gVGKy9RX7sDR
w3a7SAIcD0NDgLddvemfFJ/yLmQ4OJ8J9+1R4+PszwmzYXFBEKWr5WzTVNOSBvi3
INItVaWeI81m/dXVNQ7PHiVkFhpEqW/HsXuG/VSKff4e1jnU/6a7Zc2qnrZvFRha
Q/HtIu42eIMFTtNgFEPQkeD3OIsFLcRSIP+uPu3V/GmOWPAPJhzgrJBk2y82g9j9
gmofhYWiL1DWwV0P/2LeAIQct9txqfsxNX3LqtVCHSZmfeTjN02KHFpLiTsOI28r
LJBi+6auCn5oLIEhEhuD3o7Wg/UAxsbekmXVCJmCwl3ez1HIXdbuuytbiSK2pylp
HiBkXiXywsqUOqQRx1l4uRamnpwZRY3Ox5PFKGa5UnXgaeQMu0KtWxd7H+LsC0f3
LmkskCGEfQ2TDA3zoJMy2+adyOua300JzG74AoAAnWnGZ8CS2ONdkA6E0AlyYuS5
PBYAdaXyoT21GdSgSOl/oTa4I4zM/St3mp1AJSVlCqLg0mpPmIT1/g7my0La1C/a
xsO9Vj1rJ+m8HoLuTecLuyR/z/zfagyiyODBfb2mSMFM2b5XQE7WFd5x3BTYUP5Y
zHq9dIL4UN8D74DLasF9SlPYC99+DnPz8Mk=
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
)

func startTestServer() *httptest.Server {

	// create a listener with the desired port.
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err.Error())
	}

	respBody := make(map[string]string)
	respBody["100.100.100.100"] = "/redfish/v1/EventService/Subscriptions/2"
	body, _ := json.Marshal(respBody)
	pluginStatusRespBody := common.StatusResponse{
		Status: &common.PluginResponseStatus{
			Available: "yes",
		},
		EventMessageBus: &common.EventMessageBus{
			EmbType: "RedisStreams",
			EmbQueue: []common.EmbQueue{
				common.EmbQueue{
					QueueName: "EVENTS",
				},
			},
		},
	}
	plguinStatusBody, _ := json.Marshal(pluginStatusRespBody)
	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/ODIM/v1/Sessions" {
				w.WriteHeader(http.StatusOK)
				w.Header().Set("X-Auth-Token", "token")
				w.Write([]byte("OK"))
			} else if r.URL.Path == "/ODIM/v1/Startup" {

				w.WriteHeader(http.StatusOK)
				w.Write(body)
			} else if r.URL.Path == "/ODIM/v1/Status" {

				w.WriteHeader(http.StatusOK)
				w.Write(plguinStatusBody)
			}
		}))

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener = l
	tlsConfig := &tls.Config{}

	cert, err := tls.X509KeyPair(hostCert, hostPrivKey)
	if err != nil {
		log.Fatal("error: failed to load key pair: " + err.Error())
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.BuildNameToCertificate()

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(hostCA) {
		log.Fatal("error: failed to load CA certificate")
	}
	tlsConfig.RootCAs = capool
	tlsConfig.ClientCAs = capool

	ts.TLS = tlsConfig
	ts.Config.TLSConfig = tlsConfig
	return ts
}

func TestGenErrorResponse(t *testing.T) {
	config.SetUpMockConfig(t)

	errArgs := response.Args{
		Code:    response.GeneralError,
		Message: "",
		ErrorArgs: []response.ErrArgs{
			response.ErrArgs{
				StatusMessage: response.ResourceNotFound,
				ErrorMessage:  "Invalid request",
				MessageArgs:   []interface{}{"request", "bad request"},
			},
		},
	}
	type args struct {
		errorMessage   string
		statusMessage  string
		httpStatusCode int32
		msgArgs        []interface{}
		respPtr        *response.RPC
	}
	tests := []struct {
		name string
		args args
		want *response.RPC
	}{
		{
			name: "error response for bad request",
			args: args{
				errorMessage:   "Invalid request",
				statusMessage:  response.ResourceNotFound,
				httpStatusCode: http.StatusForbidden,
				msgArgs:        []interface{}{"request", "bad request"},
				respPtr:        &response.RPC{},
			},
			want: &response.RPC{
				StatusCode:    http.StatusForbidden,
				StatusMessage: response.ResourceNotFound,
				Body:          errArgs.CreateGenericErrorResponse(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenErrorResponse(tt.args.errorMessage, tt.args.statusMessage, tt.args.httpStatusCode, tt.args.msgArgs, tt.args.respPtr)
			if !reflect.DeepEqual(tt.args.respPtr, tt.want) {
				t.Errorf("GenErrorResposne() = %v, want %v", tt.args.respPtr, tt.want)
			}
		})
	}
}

func TestGetAllServers(t *testing.T) {
	config.SetUpMockConfig(t)
	st := StartUpInteraface{
		DecryptPassword: stubDevicePassword,
		GetAllSystems:   MockGetAllSystems,
		GetSingleSystem: MockGetSingleSystem,
	}
	servers, err := st.getAllServers("ILO")
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 2, len(servers), "there should be 2 server")

}

func TestCallPluginStartUp(t *testing.T) {
	config.SetUpMockConfig(t)
	ts := startTestServer()
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	servers := []SavedSystems{
		SavedSystems{
			ManagerAddress: "100.100.100.100",
			Password:       []byte("password"),
			UserName:       "admin",
			DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			PluginID:       "ILO",
		},
	}
	st := StartUpInteraface{
		DecryptPassword:                  stubDevicePassword,
		EMBConsume:                       stubEMBConsume,
		GetAllSystems:                    MockGetAllSystems,
		GetSingleSystem:                  MockGetSingleSystem,
		GetPluginData:                    MockGetPluginData,
		GetEvtSubscriptions:              MockGetEvtSubscriptions,
		GetDeviceSubscriptions:           MockGetDeviceSubscriptions,
		UpdateDeviceSubscriptionLocation: MockUpdateDeviceSubscriptionLocation,
	}
	err := st.callPluginStartUp(servers, "ILO")
	assert.Nil(t, err, "Error Should be nil")

	err = st.callPluginStartUp(servers, "pluginBadData")
	assert.NotNil(t, err, "error should not be nil")
}

func TestGetandStoreToken(t *testing.T) {
	var result = &PluginToken{
		Tokens: make(map[string]string),
	}
	plguinID := "GRF"
	token := "token"
	result.StoreToken(plguinID, token)
	respToken := result.GetToken(plguinID)
	assert.Equal(t, token, respToken, "respToken should be token")
}

func TestGetPluginStatus(t *testing.T) {
	config.SetUpMockConfig(t)
	ts := startTestServer()
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	password, _ := GetEncryptedKey([]byte("Password"))
	result := GetPluginStatus(&evmodel.Plugin{
		IP:                "localhost",
		Port:              "1234",
		Password:          password,
		Username:          "admin",
		ID:                "GRF",
		PreferredAuthType: "BasicAuth",
		PluginType:        "RF-GENERIC",
	},
	)
	assert.Equal(t, true, result, "Should be same")

}

func TestGenEventErrorResponse(t *testing.T) {
	resp := evresponse.EventResponse{}
	GenEventErrorResponse("error: Unable to create session with plugin GRF", response.NoValidSession, http.StatusUnauthorized,
		&resp, []interface{}{})
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should be same")
}

func TestGetPluginStatusandStartUP(t *testing.T) {
	ts := startTestServer()
	// Start the server.
	ts.StartTLS()
	defer ts.Close()
	// Intializing the TopicsList
	EMBTopics.TopicsList = make(map[string]bool)
	st := StartUpInteraface{
		DecryptPassword:                  stubDevicePassword,
		EMBConsume:                       stubEMBConsume,
		GetAllSystems:                    MockGetAllSystems,
		GetSingleSystem:                  MockGetSingleSystem,
		GetPluginData:                    MockGetPluginData,
		GetEvtSubscriptions:              MockGetEvtSubscriptions,
		GetDeviceSubscriptions:           MockGetDeviceSubscriptions,
		UpdateDeviceSubscriptionLocation: MockUpdateDeviceSubscriptionLocation,
	}
	password, _ := GetEncryptedKey([]byte("Password"))
	st.getPluginStatus(evmodel.Plugin{
		IP:                "localhost",
		Port:              "1234",
		Password:          password,
		Username:          "admin",
		ID:                "ILO",
		PreferredAuthType: "BasicAuth",
		PluginType:        "RF-GENERIC",
	})
	searchKey := GetSearchKey("100.100.100.100", evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err := st.GetDeviceSubscriptions(searchKey)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	assert.Equal(t, "https://odim.2.com/EventService/Subscriptions/1", deviceSubscription.Location, "should be same")
}
