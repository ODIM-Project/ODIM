//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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
package handle

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetChassisResource(chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	return &chassisproto.GetChassisResponse{
		StatusCode: http.StatusOK,
	}, nil
}
func mockGetChassisResourceWithRPCError(chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	return &chassisproto.GetChassisResponse{}, errors.New("Unable to RPC Call")
}

func TestChassisRPCs_GetChassisResource(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisResourceRPC = mockGetChassisResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/{id}/Power", cha.GetChassisResource)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK)
}
func TestChassisRPCs_GetChassisResourceWithRPCError(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisResourceRPC = mockGetChassisResourceWithRPCError
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/{id}/Power", cha.GetChassisResource)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestChassisRPCs_GetChassisResourceWithoutToken(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisResourceRPC = mockGetChassisResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/{id}/Power", cha.GetChassisResource)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1/Power",
	).Expect().Status(http.StatusUnauthorized)
}

func TestChassisRPCs_GetChassisCollection(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisCollectionRPC = mockGetChassisResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/", cha.GetChassisCollection)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK)
	e.GET(
		"/redfish/v1/Chassis/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}
func TestChassisRPCs_GetChassisCollectionWithRPCError(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisCollectionRPC = mockGetChassisResourceWithRPCError
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/", cha.GetChassisCollection)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestChassisRPCs_GetChassis(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisRPC = mockGetChassisResource
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/", cha.GetChassis)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusOK)
	e.GET(
		"/redfish/v1/Chassis/",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
}
func TestChassisRPCs_GetChassisWithRPCError(t *testing.T) {
	var cha ChassisRPCs
	cha.GetChassisRPC = mockGetChassisResourceWithRPCError
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/Chassis")
	redfishRoutes.Get("/", cha.GetChassis)

	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/Chassis/",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestChassisRPCs_CreateChassisWithNoInputBody(t *testing.T) {
	sut := ChassisRPCs{}
	app := iris.New()
	app.Any("/", sut.CreateChassis)

	resp := httptest.New(t, app).POST("/").
		WithBytes(nil).
		Expect()

	resp.Status(http.StatusBadRequest)
	resp.JSON().Schema(redfishErrorSchema)
	resp.Headers().Equal(map[string][]string{
		"Connection":      {"keep-alive"},
		"Odata-Version":   {"4.0"},
		"X-Frame-Options": {"sameorigin"},
		"Content-Type":    {"application/json; charset=utf-8"},
	})
}

func TestChassisRPCs_CreateChassisWithRPCError(t *testing.T) {
	sut := ChassisRPCs{
		CreateChassisRPC: func(req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error) {
			return nil, fmt.Errorf("RPC ERROR")
		},
	}

	app := iris.New()
	app.Any("/", sut.CreateChassis)

	resp := httptest.New(t, app).POST("/").
		WithBytes([]byte(`{"chassis":"creationRequest"}`)).
		Expect()

	resp.Status(http.StatusInternalServerError)
	resp.JSON().Schema(redfishErrorSchema)
	resp.Headers().Equal(map[string][]string{
		"Connection":      {"keep-alive"},
		"Odata-Version":   {"4.0"},
		"X-Frame-Options": {"sameorigin"},
		"Content-Type":    {"application/json; charset=utf-8"},
	})
}

func TestChassisRPCs_CreateChassis(t *testing.T) {
	expectedRPCResponse := chassisproto.GetChassisResponse{
		StatusCode: http.StatusOK,
		Header: map[string]string{
			"Location": "/redfish/odim/blebleble",
		},
		Body: []byte(`{"chassis":"body"}`),
	}

	sut := ChassisRPCs{
		CreateChassisRPC: func(req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error) {
			return &expectedRPCResponse, nil
		},
	}

	app := iris.New()
	app.Any("/", sut.CreateChassis)

	resp := httptest.New(t, app).POST("/").
		WithBytes([]byte(`{"chassis":"creationRequest"}`)).
		Expect()

	resp.Status(http.StatusOK)
	resp.Body().Contains(string(expectedRPCResponse.Body))
	resp.Headers().Equal(map[string][]string{
		"Connection":      {"keep-alive"},
		"Odata-Version":   {"4.0"},
		"X-Frame-Options": {"sameorigin"},
		"Location":        {"/redfish/odim/blebleble"},
	})
}
func TestChassisRPCs_CreateChassisWithMalformedBody(t *testing.T) {
	expectedRPCResponse := chassisproto.GetChassisResponse{
		StatusCode: http.StatusBadRequest,
		Body: []byte(`{
  "error": {
    "code": "Base.1.6.1.GeneralError",
    "message": "An error has occurred. See ExtendedInfo for more information.",
    "@Message.ExtendedInfo": [
      {
        "@odata.type": "#Message.v1_0_8.Message",
        "MessageId": "Base.1.6.1.MalformedJSON",
        "Message": "The request body submitted was malformed JSON and could not be parsed by the receiving service.error while trying to read obligatory json body: invalid character '[' looking for beginning of object key string",
        "Severity": "Critical",
        "Resolution": "Ensure that the request body is valid JSON and resubmit the request."
      }
    ]
  }
}
`),
	}

	sut := ChassisRPCs{
		CreateChassisRPC: func(req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error) {
			return &expectedRPCResponse, nil
		},
	}

	app := iris.New()
	app.Any("/", sut.CreateChassis)

	resp := httptest.New(t, app).POST("/").
		WithBytes([]byte(`{"Sample":"Body","Links":{[],},}`)).
		Expect()
	resp.Status(http.StatusBadRequest)
	resp.Body().Contains(string(expectedRPCResponse.Body))

}

var redfishErrorSchema = `
{
   "$schema": "http://json-schema.org/draft-04/schema#",
   "type": "object",
   "properties": {
      "error": {
		"type":"object",
		"properties": {
          "@Message.ExtendedInfo": {
			 "type": "array",
			 "items": {
				"type": "object",
				"properties": {
					"@odata.type": {
						"type": "string"
					},
					"Message": {
						"type": "string"
					},
					"MessageId": {
						"type": "string"
					},
					"Resolution": {
						"type": "string"
					},
					"Severity": {
						"type": "string",
						"enum": ["Critical", "Warning"]
					}
				},
				"required": ["@odata.type", "Message", "MessageId", "Resolution", "Severity"]
			 }
		  },
			
		  "code": {
			 "type": "string"
		  },
			
		  "message": {
			 "type": "string"
		  }
        },
		"required":["code", "message"]
      }
   },
   "required": ["error"]
}`
