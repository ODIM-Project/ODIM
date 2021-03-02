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

// Package events have the functionality of
// - Create Event Subscription
// - Delete Event Subscription
// - Get Event Subscription
// - Post Event Subscription to destination
// and corresponding unit test cases
package events

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
)

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return err
	}
	return nil
}

func mockFabricsData(t *testing.T) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	fabrics := evmodel.Fabric{
		FabricUUID: "123456",
		PluginID:   "CFM",
	}
	if err = connPool.Create("Fabric", "123456", &fabrics); err != nil {
		return err
	}
	return nil
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockTargetandPlugin(t *testing.T) error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	encryptedData := getEncryptedKey(t, []byte("testData"))
	password := getEncryptedKey(t, []byte("Password"))

	targetArr := []evmodel.Target{
		{
			ManagerAddress: "10.4.1.2",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			PluginID:       "GRF",
		},
		{
			ManagerAddress: "10.4.1.3",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "11081de0-4859-984c-c35a-6c50732d72da",
			PluginID:       "ILO",
		},
		{
			ManagerAddress: "10.4.1.4",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "d72dade0-c35a-984c-4859-1108132d72da",
			PluginID:       "ILO",
		},
		{
			ManagerAddress: "10.4.1.5",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "110813e0-4859-984c-984c-d72da32d72da",
			PluginID:       "ILO",
		},
		{
			ManagerAddress: "localhost",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "abab09db-e7a9-4352-8df0-5e41315a2a4c",
			PluginID:       "ILO",
		},
	}
	for _, target := range targetArr {
		const table string = "System"
		//Save data into Database
		if err = connPool.Create(table, target.DeviceUUID, &target); err != nil {
			return err
		}
	}
	pluginArr := []evmodel.Plugin{
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "GRF",
		},
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "ILO",
		},
		{
			IP:                "10.4.1.6",
			Port:              "4321",
			Password:          password,
			Username:          "admin",
			ID:                "CFM",
			PreferredAuthType: "XAuthToken",
			PluginType:        "CFM",
		},
		{
			IP:                "10.24.1.23",
			Port:              "4321",
			Password:          password,
			Username:          "admin",
			ID:                "CFMPlugin",
			PreferredAuthType: "XAuthToken",
			PluginType:        "CFMPlugin",
		},
	}
	for _, plugin := range pluginArr {
		pl := "Plugin"
		//Save data into Database
		if err = connPool.Create(pl, plugin.ID, &plugin); err != nil {
			return err
		}
	}
	var reqData = `{"@odata.id":"/redfish/v1/Systems/1"}`
	mockSystemResourceData([]byte(reqData), "ComputerSystem", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")
	mockSystemResourceData([]byte(reqData), "ComputerSystem", "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da:1")
	mockSystemResourceData([]byte(reqData), "ComputerSystem", "/redfish/v1/Systems/d72dade0-c35a-984c-4859-1108132d72da:1")
	mockSystemResourceData([]byte(reqData), "ComputerSystem", "/redfish/v1/Systems/abab09db-e7a9-4352-8df0-5e41315a2a4c:1")
	return nil
}

var (
	hostCA = []byte(`-----BEGIN CERTIFICATE-----
MIIFFDCCAvygAwIBAgIUSRrUi66SkjZFWwLVDhdyrFX65pYwDQYJKoZIhvcNAQEN
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIwMDQxODEyMjYwN1oXDTMwMDQx
NjEyMjYwN1owFDESMBAGA1UEAwwJbG9jYWxob3N0MIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEAtl+m9W7LvxhV9ghNu8Rr2VMUCha3LB+k4XOiey0Wt/DF
e1D36tJBbfNHfe+6m4/GHM4e4Tp+h8inCyBCCaHLPXfcA0jY3bHOC0cYjRjhx4rP
TG06ok2YIhovfJIZcANse8nSCenAFSSJqdAAKlsNH9z5lT+JhPr9kG+48WnvDPeX
DNWNf/VUqkYjon7q5hLQL9ImtGSFmjPMKFjNACSDZfIyb/x6TUuTVHBmtKT25+rp
4FHa3mxBYNLisavyNlg59U189FISAGYjC+fOlYQBL7RXb2Xrrd4DgGJC1N8Hr0To
z9zmFj0u4sJ5RFdOabDCNZgEx3tGD1fYqjEdSuspRzMu6CUl8PjN/C0Ml0NOHyHn
7x/0VKxOwGUW0eoGE9BNGfSONrX6Wt3Ej3uKkyDPxIX0PhbgHewQF/jF+698hv10
uGoeAZLhyZlFhCxPCZe4EtBOQ7LXWcEroneoVj4CWY0tm1NVvDY7PVw4MsnloHW/
/qfjgi9yQbaJ7IllSNGaZhqYvtc8bXH5xwd4/F9wOkKC6p63RqB0Pom7lDOBeiha
uhE8gsqg0oCJCkYyKdQ8Ye6xiF+cfmnV5VT+S9XLyTza3utuuJPG8qg6xY5aVBkg
P5JrVaiDTQg19FopUFHfUVY8K121s4WMS02XjySpLyY11ZGrqp3hqIRL+kUJX7sC
AwEAAaNeMFwwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUE3wuzAFexqxv1/4o
HH2HNVut8q4wCwYDVR0PBAQDAgHmMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEF
BQcDATANBgkqhkiG9w0BAQ0FAAOCAgEATgew4Rh3VEgOm6OF5od83AUm/R2nomMe
FXNAw70duP3miYsQbeE9jiLL5vF3Rwv9X3EQ7yBo+EAnKaaMha5Ee9IT8Yu9LnFj
qYKkPlK/fslFROFASZCyLltwJ3KBwfNcUBh+55IrYjvEzh1NUOK/9r+Lt7fn6KVq
a0vVqF9HZjnmr7I4yOw142QMS7nPmQ/PF7ql/Q8IuAuFhJA3pFPK5YO0p0cdaWmS
X0RSmxr0ANkbPWp82lonX4aYz+2uqXDHWBv+kQBMdhU400s9MC+q1vjMD+C+qG3g
BO1bt8/p/DDFnZPwwUyHvU8VIxL99oYDi4nk02YRJrXpsWhHogeD9m8Kv3BcVFx2
7kvuhV+2UiX6IMyricN4EQzG+XS+Pu6UkxiHrxIl1BNVjc3pfRjPUOIU8wt1nBHS
v3jiYtR7fAqWHdfiqz9LDjfxam1QqOt3l6CVMnuQvaZF+nzC/Aa7l09ElrsSWrF5
sp5MSJplugY7r3qG0/evQ/HzkThhVD+0HMwUGx9xa+Z96kRzQ5qYEaj03IU8qZve
wMEvS6jE4aNQ5PFZ/9HKIJ7LkcubuBpaMGTeNzO6vn+9XvxxHyTlMWLpfm04ukBI
3ylqIAA630pHjCUJGPHBiYXw9DD7MRGczqkOLVPE3sOHuTiyL/DNQ9kSReyMG8Ss
8Uv+gESYkqk=
-----END CERTIFICATE-----`)
	hostCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFVjCCAz6gAwIBAgIUKcTsAJKTQSMIKWDWYdZf/Y0uxLUwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MB4XDTIwMDQxODEyMjk1OFoXDTMwMDQx
NTEyMjk1OFowFDESMBAGA1UEAwwJbG9jYWxob3N0MIICIjANBgkqhkiG9w0BAQEF
AAOCAg8AMIICCgKCAgEA5eYgv0jdi0GY+x+fs9bFwEX6Knam7cWhxUjAMUs00ii8
8n6o5XE0FXbfLs44xf/FpYkpZtf+qXqPXLQCohfeUUiOXGRTGUkQA5hVVxj6gE+U
NqrOSocxwep3KpACPYIniAB2hSz0Dt6fxAdM29SONtNHI9xihuGm9deN+5Dx+ZgW
gcaMfio62wlwTdMcD6VZT3oaYog/+dvS7tqU55tROC335GEGYRg2E7RP+elsvESK
wm1x8US5cfPWgcY2P/hGfbrnhsqW2XVaVrA7KQzWS9/1WRrq7j41sIbG5WHydMLI
kmaXnMuZ+xz4tcEsnt2MOlL/izXJs4H+pUpL+iTF/i9MSRfw2Bi18fVH1CitNRdo
M5SrAmt+YFeDw8N6SuXtIu8l2fQecCeieCSbHraW+VGqY+n9Em4MNgw+tSxl3lsi
j4ijvbptGjbvs2Vkr37MwYqfwIWfLEOdYJVcNgJI72T+mhpoHtKTZ50Ugk4MOkNL
WWuX8wX7UHswwuSTv2jORBXdoyq3mz6PqNrnXOGzmujhVOyp9PxV/AT66Px5GIWy
VpC4MMHiJCNSGP/ocsKepAg23fmboQ3OzjB9Y+AHUnp8ML/h2HhCKMSREj/BilRO
K+Om3ZA/baCngtEvna2pXKuHFKHiHILpyaGrlhRm55j094VEHrNfNHfCJ2+Qf10C
AwEAAaOBnzCBnDAdBgNVHQ4EFgQUqNezn0zjq7KHDGV0IZy94danHiwwTwYDVR0j
BEgwRoAUE3wuzAFexqxv1/4oHH2HNVut8q6hGKQWMBQxEjAQBgNVBAMMCWxvY2Fs
aG9zdIIUSRrUi66SkjZFWwLVDhdyrFX65pYwCwYDVR0PBAQDAgWgMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATANBgkqhkiG9w0BAQsFAAOCAgEAiqyKAGvP
HWt8boQe6KkiH5tmPjrd9U6YSsB3qzzfac20BIuY7ZSdpNf2REt/pAUKUyO6mYyV
uuIHtwd58Yuvdop3WD+K70qTKqXB8wNeyfQNeWFhEDumJikbjcRRNn9P+GhCTATI
wp+U7zxxXOxWWlm1sDsZmAuouaFVDYG0F4AUdHP16gLrZEkzup9r7VnUFpw/0tkB
k1PR8RKgpz4r/x5QtCft+FFcyPYqRZNiA/EJ9o+ObtZe2ig5QTpzJhGSwTRm9emt
PxLKv672deyjyOypIgwa94dZRtaiZWK8apDMmzx/ab1Q0rMyZzaVp/EHmEyDQNDd
zqiapUOu2t4vu3NhqgEjLTfSgQR7Kj5BguzQfZ3oJr4zn8JGKfkPdW0OV0LFG9o3
C0HvinWBueGW8dpiHGbOVH57iOx8CU+AesDK5T2rD3kryJ+fos0fSO0fWwmkNYMZ
WkH7hmtVup2jCQy97T0Jlpwi9/gkhnJW4XV1HbAtdkX5QiHC0xs7oTlTpSjm409d
8Zqevusshp2FLkZC4NypW34dJhkrJl/vFJN+zRlwNQUu4fEVFJUb3tTkPWGqYWwz
Ne4ovENCT7e8oBqFDJrATPbHs06b1GNk1nxO1BtAgnba37/7WyQxXK9aD5j8w3pK
dIRctXGgrIroKS/k4eFQa99kbmFl5Khb8Wo=
-----END CERTIFICATE-----`)
	hostPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
)

func startTestServer() *httptest.Server {

	// create a listener with the desired port.
	l, err := net.Listen("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.URL.Path)
			if r.URL.Path == "/ODIM/v1/Sessions" {

				w.WriteHeader(http.StatusOK)
				w.Header().Set("X-Auth-Token", "token")
				w.Write([]byte("OK"))
			} else if r.URL.Path == "/ODIM/v1/Subscriptions" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			} else if r.URL.Path == "/ODIM/v1/Subscriptions/12345" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			}
		}))

	// NewUnstartedServer creates a listener. Close that listener and replace
	// with the one we created.
	ts.Listener = l

	tlsConfig := &tls.Config{}

	cert, err := tls.X509KeyPair(hostCert, hostPrivKey)
	if err != nil {
		log.Fatalf("error: failed to load key pair: %v", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.BuildNameToCertificate()

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(hostCA) {
		log.Fatalf("error: failed to load CA certificate")
	}
	tlsConfig.RootCAs = capool
	tlsConfig.ClientCAs = capool

	ts.TLS = tlsConfig
	ts.Config.TLSConfig = tlsConfig
	return ts
}

func TestDeleteEventSubscription(t *testing.T) {
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mockTargetandPlugin(t)

	storeTestEventDetails(t)

	pc := PluginContact{
		Auth:          mockIsAuthorized,
		ContactClient: mockContactClient,
	}

	// positive test case with basic auth type
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
	}
	resp := pc.DeleteEventSubscriptionsDetails(req)
	data := resp.Body.(response.Response)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	assert.Equal(t, "81de0110-c35a-4859-984c-072d6c5a32d7", data.ID, "ID should be 81de0110-c35a-4859-984c-072d6c5a32d7")

	// positive test case with basic auth type
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d8",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// positive test case deletion of collection subscription
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a3211",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// Negative test cases
	// if subscription id is bot present
	req1 := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "de018110-4859-984c-c35a-0a32d772d6c5",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req1)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")

	// Invalid token
	req2 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req2)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")
}

func TestDeleteEventSubscriptionOnDeletServer(t *testing.T) {
	config.SetUpMockConfig(t)
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mockTargetandPlugin(t)

	storeTestEventDetails(t)

	pc := PluginContact{
		ContactClient: mockContactClient,
	}

	// positive test case
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	}
	resp := pc.DeleteEventSubscriptions(req)
	assert.Equal(t, http.StatusNoContent, int(resp.StatusCode), "Status Code should be StatusNoContent")

	// Negative test cases
	// if UUID is invalid
	req1 := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "de018110-4859-984c-c35a-0a32d772d6c5",
	}
	resp = pc.DeleteEventSubscriptions(req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	// if UUID is is not present in DB
	req1 = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/de018110-4859-984c-c35a-0a32d772d6c5:1",
	}

	resp = pc.DeleteEventSubscriptions(req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	//storeTestEventDetails(t)
	req = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	}
	resp = pc.DeleteEventSubscriptions(req)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")

}

func TestDeleteEventSubscriptionOnFabrics(t *testing.T) {
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"CFM": "token",
	}
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mockTargetandPlugin(t)
	mockFabricsData(t)
	storeTestEventDetails(t)

	pc := PluginContact{
		Auth:          mockIsAuthorized,
		ContactClient: mockContactClient,
	}

	// positive test case with basic auth type
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d9",
	}
	resp := pc.DeleteEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
}

func TestIsCollectionOriginResourceURI(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "Positive: First element in list",
			arg:  "/redfish/v1/Systems",
			want: true,
		},
		{
			name: "Positive: Last element in list",
			arg:  "/redfish/v1/TaskService/Tasks",
			want: true,
		},
		{
			name: "Positive: Middle element in list",
			arg:  "/redfish/v1/Fabrics/",
			want: true,
		},
		{
			name: "Negative: Empty string",
			arg:  "",
			want: false,
		},
		{
			name: "Negative: Non-existent element",
			arg:  "non-existent string",
			want: false,
		},
		{
			name: "Negative: Non-existent element -2",
			arg:  "/redfish/v1/Fabrics1",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCollectionOriginResourceURI(tt.arg); got != tt.want {
				t.Errorf("isCollectionOriginResourceURI got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestDeleteFabricsSubscription(t *testing.T) {
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	storeTestEventDetails(t)

	pc := PluginContact{
		ContactClient: mockContactClient,
	}

	password := getEncryptedKey(t, []byte("Password"))
	// positive test case with basic auth type
	plugin := &evmodel.Plugin{
		IP:                "10.4.1.5",
		Port:              "1234",
		Password:          password,
		Username:          "admin",
		ID:                "ILO",
		PreferredAuthType: "BasicAuth",
		PluginType:        "ILO",
	}
	resp, err := pc.DeleteFabricsSubscription("", plugin)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// Negative test cases
	// if subscription id is not present
	plugin.IP = "10.10.10.10"
	resp, err = pc.DeleteFabricsSubscription("", plugin)
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")
}

func TestDeleteSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	mockTargetandPlugin(t)

	storeTestEventDetails(t)

	pc := PluginContact{
		Auth:          mockIsAuthorized,
		ContactClient: mockContactClient,
	}

	encryptedData := getEncryptedKey(t, []byte("testData"))
	target := evmodel.Target{
		ManagerAddress: "10.4.1.2",
		Password:       encryptedData,
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}

	err := pc.deleteSubscription(&target, "/redfish/v1/Systems")
	assert.Nil(t, err, "error should be nil")

	target.PluginID = "non-existent"
	err = pc.deleteSubscription(&target, "/redfish/v1/Systems")
	assert.NotNil(t, err, "error should not be nil")
}
