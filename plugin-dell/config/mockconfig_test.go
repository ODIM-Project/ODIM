// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
// package ...
package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetUpMockConfig(t *testing.T) {
	if err := SetUpMockConfig(t); err != nil {
		t.Error("error: SetUpMockConfig failed with", err)
	}
	if Data.FirmwareVersion == "" {
		t.Error("error: Data.FirmwareVersion is not initialized")
	}
	if Data.RootServiceUUID == "" {
		t.Error("error: Data.RootServiceUUID is not initialized")
	}
	if Data.SessionTimeoutInMinutes == 0 {
		t.Error("error: Data.SessionTimeoutInMinutes is not initialized")
	}
	if Data.PluginConf == nil {
		t.Error("error: Data.PluginConf is not initialized")
	}
	if Data.KeyCertConf == nil {
		t.Error("error: Data.KeyCertConf is not initialized")
	}
	if Data.LoadBalancerConf == nil {
		t.Error("error: Data.LoadBalancerConf is not initialized")
	}
	if Data.EventConf == nil {
		t.Error("error: Data.EventConf is not initialized")
	}
	if Data.MessageBusConf == nil {
		t.Error("error: Data.MessageBusConfis not initialized")
	}
	if Data.URLTranslation == nil {
		t.Error("error: Data.URLTranslation is not initialized")
	}
	if Data.TLSConf == nil {
		t.Error("error: Data.TLSConf is not initialized")
	}
}

func TestGetPublicKey(t *testing.T) {
	SetUpMockConfig(t)
	key := GetPublicKey()
	hostPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuZwJSkQK5blNhxu+Fo5c
xeUcMX9rpUcB8us4BwZCsGq5DDpY8iunwmxLjtZb/fLFiz6iAfWx1vqLOcXPYbeY
LjF8jIqJaWuYryrV9WRctw5p7OdiYmtJK8ILqe08VIZLfs8qr0KZZP5zzoMNEntu
Elbs3Id2HUTrj7uJbSTZVMS32oJUEqtDzNK9pDl+cQIKFiV7Do+KPyMamKeiiari
zDKiyYsNxtBS+53Cp1MPctqKwcr85u5aN1MXZnDSVoB6HewwuPlrLzf/f1d0H7Hf
LJzAjxA9ikizJPL90oQiA94Ra1ZcTSMKZxbcErPJoOEWqMwTAzmYfd7KDinu64vL
NF+CEQEJlLFdMIf3zIDQKY9UI8SD9JqM1NYfzH6a8GGK3rqEUBDrLkvUbOZs8DV6
3YzY7ZB0lDxxtV/BVoSoVqONNYFyn7/vz+HXCaFGuO5x3ddPb5Gt0ckUWV9h3AsL
CPe1s3VnWVys/lJyLuTGRs1QdR77gXQbv29g6QfB4fIrqaOit4DguTV1xmyWjIhj
BMaLcqDJJ1bPJiyhMl5fQvnRgyk/HbKejW7wli59OnZW9stYxrhrPTqVfJOWvnJE
Bq4VYWoMrcs2G3NGfgwBABsMEYbm2Nn558Nv8OkXuYd2ENFndoSxRa5Crk3HZ5mE
Fy7PCcRO16uhaVrY97PbthcCAwEAAQ==
-----END PUBLIC KEY-----`)
	assert.Equal(t, hostPubKey, key, "Valid Host key ")
}

func TestGetRandomPort(t *testing.T) {
	port := GetRandomPort()
	assert.NotEmpty(t, port, "Port should be not empty")
}
