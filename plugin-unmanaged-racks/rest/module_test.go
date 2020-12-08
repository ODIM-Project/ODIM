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

package rest

import (
	stdContext "context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/db"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/alicebob/miniredis/v2"
	"github.com/gavv/httpexpect"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/require"
)

var TEST_CONFIG = config.PluginConfig{
	RootServiceUUID: "99999999-9999-9999-9999-999999999999",
	UserName:        "admin",
	Password:        "O01bKrP7Tzs7YoO3YvQt4pRa2J_R6HI34ZfP4MxbqNIYAVQVt2ewGXmhjvBfzMifM7bHFccXKGmdHvj3hY44Hw==",
	FirmwareVersion: "0.0.0",
	OdimNBUrl:       "https://localhost:45000",
	URLTranslation: &config.URLTranslation{
		NorthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
		SouthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
	},
}

func Test_secured_endpoints_return_401_when_unauthorized(t *testing.T) {
	tests := []struct {
		method string
		uri    string
	}{
		{http.MethodPost, "/ODIM/v1/Startup"},

		{http.MethodGet, "/ODIM/v1/Chassis"},
		{http.MethodPost, "/ODIM/v1/Chassis"},
		{http.MethodGet, "/ODIM/v1/Chassis/1"},
		{http.MethodDelete, "/ODIM/v1/Chassis/1"},
		{http.MethodPatch, "/ODIM/v1/Chassis/1"},

		{http.MethodGet, "/ODIM/v1/Managers"},
		{http.MethodGet, "/ODIM/v1/Managers/1"},
	}

	testApp, _ := createTestApplication()
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s", test.method, test.uri), func(t *testing.T) {
			he := httptest.New(t, testApp)
			resp := he.Request(test.method, test.uri).Expect()
			resp.Status(http.StatusUnauthorized)
		})
	}
}

func Test_unsecured_endpoints_return_NON_401(t *testing.T) {
	tests := []struct {
		method string
		uri    string
	}{
		{http.MethodGet, "/ODIM/v1/Status"},
		{http.MethodPost, "/EventService/Events"},
	}

	testApp, _ := createTestApplication()
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s %s", test.method, test.uri), func(t *testing.T) {
			he := httptest.New(t, testApp)
			resp := he.Request(test.method, test.uri).Expect()
			require.NotEqual(t, http.StatusUnauthorized, resp.Raw().Status)
		})
	}
}

func Test_empty_chassis_collection_is_exposed(t *testing.T) {
	testApp, _ := createTestApplication()
	httptest.New(t, testApp).
		GET("/ODIM/v1/Chassis").
		WithBasicAuth("admin", "Od!m12$4").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		ValueEqual("Members", []redfish.Link{}).
		ValueEqual("Members@odata.count", 0)
}

func Test_invalid_chassis_creation_request_should_be_rejected(t *testing.T) {
	testApp, _ := createTestApplication()

	t.Run("no body", func(t *testing.T) {
		httptest.New(t, testApp).
			POST("/ODIM/v1/Chassis").
			WithBasicAuth("admin", "Od!m12$4").
			Expect().
			Status(http.StatusBadRequest).
			ContentType("application/json", "UTF-8")
	})

	t.Run("empty body", func(t *testing.T) {
		httptest.New(t, testApp).
			POST("/ODIM/v1/Chassis").
			WithBasicAuth("admin", "Od!m12$4").
			WithBytes([]byte("{}")).
			Expect().
			Status(http.StatusBadRequest).
			ContentType("application/json", "UTF-8")
	})
}

func Test_creation_of_rack_pointing_to_not_existing_rack_group_should_be_rejected(t *testing.T) {
	testApp, _ := createTestApplication()
	httptest.New(t, testApp).
		POST("/ODIM/v1/Chassis").
		WithBasicAuth("admin", "Od!m12$4").
		WithBytes([]byte(`
						{
							"Name": "Rack#1",
							"ChassisType": "Rack",
							"Links": {
								"ManagedBy": [
									{"@odata.id": "/ODIM/v1/Managers/99999999-9999-9999-9999-999999999999"}
								],
								"ContainedBy": [
									{"@odata.id": "/not/existing/rack-group"}
								]
							}
						}`),
		).
		Expect().
		Status(http.StatusBadRequest).
		ContentType("application/json", "UTF-8")
}

func Test_creation_of_chassis_with_previously_used_name_should_be_rejected(t *testing.T) {
	testApp, _ := createTestApplication()

	t.Run("creating rack-group", func(t *testing.T) {
		httptest.New(t, testApp).
			POST("/ODIM/v1/Chassis").
			WithBasicAuth("admin", "Od!m12$4").
			WithBytes([]byte(`
				{
					"Name": "RackGroup#1",
					"ChassisType": "RackGroup",
					"Links": {
						"ManagedBy": [
							{"@odata.id": "/ODIM/v1/Managers/99999999-9999-9999-9999-999999999999"}
						]
					}
				}
			`)).
			Expect().
			Status(http.StatusCreated).
			ContentType("application/json", "UTF-8").
			Headers().Value("Location").NotNull().Array().NotEmpty().Element(0).String().NotEmpty().Raw()

		t.Run("checking chassis collection size", func(t *testing.T) {
			getChassisCollection(httptest.New(t, testApp), 1)
		})
	})

	t.Run("creating rack-group with same name again", func(t *testing.T) {
		httptest.New(t, testApp).
			POST("/ODIM/v1/Chassis").
			WithBasicAuth("admin", "Od!m12$4").
			WithBytes([]byte(`
				{
					"Name": "RackGroup#1",
					"ChassisType": "RackGroup",
					"Links": {
						"ManagedBy": [
							{"@odata.id": "/ODIM/v1/Managers/99999999-9999-9999-9999-999999999999"}
						]
					}
				}
			`)).
			Expect().
			Status(http.StatusConflict).
			ContentType("application/json", "UTF-8")

		t.Run("checking chassis collection size", func(t *testing.T) {
			getChassisCollection(httptest.New(t, testApp), 1)
		})

	})
}

func getChassisCollection(e *httpexpect.Expect, expectedCollectionSize int) {
	e.GET("/ODIM/v1/Chassis").
		WithBasicAuth("admin", "Od!m12$4").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		ValueEqual("Members@odata.count", expectedCollectionSize).
		Path("$.Members").Array().Length().Equal(expectedCollectionSize)
}

func Test_unmanaged_chassis_chain(t *testing.T) {

	testApp, _ := createTestApplication()

	t.Run("create rack group", func(t *testing.T) {
		rackGroupURI := httptest.New(t, testApp).
			POST("/ODIM/v1/Chassis").
			WithBasicAuth("admin", "Od!m12$4").
			WithBytes([]byte(`
				{
					"Name": "RackGroup#1",
					"ChassisType": "RackGroup",
					"Links": {
						"ManagedBy": [
							{"@odata.id": "/ODIM/v1/Managers/99999999-9999-9999-9999-999999999999"}
						]
					}
				}
			`)).
			Expect().
			Status(http.StatusCreated).
			ContentType("application/json", "UTF-8").
			Headers().Value("Location").NotNull().Array().NotEmpty().Element(0).String().NotEmpty().Raw()

		t.Run("get created rack group", func(t *testing.T) {
			httptest.New(t, testApp).
				GET(rackGroupURI).
				WithBasicAuth("admin", "Od!m12$4").
				Expect().
				Status(http.StatusOK).
				ContentType("application/json", "UTF-8").
				JSON().Object().
				ValueEqual("ChassisType", "RackGroup").
				ValueEqual("Name", "RackGroup#1")
		})

		t.Run("create rack", func(t *testing.T) {
			rackURI := httptest.New(t, testApp).
				POST("/ODIM/v1/Chassis").
				WithBasicAuth("admin", "Od!m12$4").
				WithBytes([]byte(`
						{
							"Name": "Rack#1",
							"ChassisType": "Rack",
							"Links": {
								"ManagedBy": [
									{"@odata.id": "/ODIM/v1/Managers/99999999-9999-9999-9999-999999999999"}
								],
								"ContainedBy": [
									{"@odata.id": "`+rackGroupURI+`"}
								]
							}
						}`),
				).
				Expect().
				Status(http.StatusCreated).
				ContentType("application/json", "UTF-8").
				Headers().Value("Location").NotNull().Array().NotEmpty().Element(0).String().NotEmpty().Raw()

			t.Run("verify if rack-group contains rack", func(t *testing.T) {
				httptest.New(t, testApp).
					GET(rackGroupURI).
					WithBasicAuth("admin", "Od!m12$4").
					Expect().
					ContentType("application/json", "UTF-8").
					Status(http.StatusOK).
					JSON().Path(`$.Links.Contains[0]["@odata.id"]`).String().Equal(rackURI)
			})

			t.Run("delete occupied rack-group", func(t *testing.T) {
				httptest.New(t, testApp).
					DELETE(rackGroupURI).
					WithBasicAuth("admin", "Od!m12$4").
					Expect().
					Status(http.StatusConflict).
					ContentType("application/json", "UTF-8")
			})

			t.Run("try to attach system under rack", func(t *testing.T) {
				odimstubapp := odimstub{iris.New()}
				odimstubapp.Run()
				defer odimstubapp.Stop()

				httptest.New(t, testApp).
					PATCH(rackURI).
					WithBasicAuth("admin", "Od!m12$4").
					WithBytes([]byte(`
						{
							"Links":{
								"Contains": [
									{"@odata.id":"/ODIM/v1/Systems/1"}
								]
							}
						}
					`)).
					Expect().
					Status(http.StatusBadRequest)
			})

			t.Run("attach chassis under rack", func(t *testing.T) {
				odimstubapp := odimstub{iris.New()}
				odimstubapp.Run()
				defer odimstubapp.Stop()

				httptest.New(t, testApp).
					PATCH(rackURI).
					WithBasicAuth("admin", "Od!m12$4").
					WithBytes([]byte(`
						{
							"Links":{
								"Contains": [
									{"@odata.id":"/ODIM/v1/Chassis/1"}
								]
							}
						}
					`)).
					Expect().
					Status(http.StatusOK).
					ContentType("application/json", "UTF-8").
					JSON().Path(`$.Links.Contains`).Array().Length().Equal(1)
			})

			t.Run("attach another chassis under rack", func(t *testing.T) {
				odimstubapp := odimstub{iris.New()}
				odimstubapp.Run()
				defer odimstubapp.Stop()

				httptest.New(t, testApp).
					PATCH(rackURI).
					WithBasicAuth("admin", "Od!m12$4").
					WithBytes([]byte(`
						{
							"Links":{
								"Contains": [
									{"@odata.id":"/ODIM/v1/Chassis/1"},
									{"@odata.id":"/ODIM/v1/Chassis/2"}
								]
							}
						}
					`)).
					Expect().
					Status(http.StatusOK).
					ContentType("application/json", "UTF-8").
					JSON().Path(`$.Links.Contains`).Array().Length().Equal(2)
			})

			t.Run("try to delete occupied rack", func(t *testing.T) {
				httptest.New(t, testApp).
					DELETE(rackURI).
					WithBasicAuth("admin", "Od!m12$4").
					Expect().
					Status(http.StatusConflict).
					ContentType("application/json", "UTF-8")
			})

			t.Run("detach Chassis/1", func(t *testing.T) {
				odimstubapp := odimstub{iris.New()}
				odimstubapp.Run()
				defer odimstubapp.Stop()

				httptest.New(t, testApp).
					PATCH(rackURI).
					WithBasicAuth("admin", "Od!m12$4").
					WithBytes([]byte(`
						{
							"Links":{
								"Contains": [
									{"@odata.id":"/ODIM/v1/Chassis/2"}
								]
							}
						}
					`)).
					Expect().
					Status(http.StatusOK).
					ContentType("application/json", "UTF-8").
					JSON().Path(`$.Links.Contains`).Array().Length().Equal(1)
			})

			t.Run("resource removed event detaches Chassis/2 existing under rack", func(t *testing.T) {
				odimstubapp := odimstub{iris.New()}
				odimstubapp.Run()
				defer odimstubapp.Stop()

				httptest.New(t, testApp).
					POST("/EventService/Events").
					WithBytes([]byte(`
						{
							"Events": [
								{
									"OriginOfCondition": {"@odata.id": "/redfish/v1/Chassis/2"}
								}
							]
						}
					`)).
					Expect().
					Status(http.StatusOK).
					NoContent()
			})

			t.Run("delete rack", func(t *testing.T) {
				httptest.New(t, testApp).
					DELETE(rackURI).
					WithBasicAuth("admin", "Od!m12$4").
					Expect().
					Status(http.StatusNoContent).
					Body().Empty()
			})

			t.Run("delete rack-group", func(t *testing.T) {
				httptest.New(t, testApp).
					DELETE(rackGroupURI).
					WithBasicAuth("admin", "Od!m12$4").
					Expect().
					Status(http.StatusNoContent).
					NoContent()
			})
		})
	})
}

type odimstub struct {
	app *iris.Application
}

func (o *odimstub) Run() {
	o.app.Get("/redfish/v1/Systems", func(context context.Context) {
		context.JSON(redfish.NewCollection(
			"/redfish/v1/Systems",
			"#ComputerSystemCollection.ComputerSystemCollection",
			[]redfish.Link{
				{Oid: "/redfish/v1/Systems/1"},
			}...,
		))
	})

	o.app.Get("/redfish/v1/Chassis", func(context context.Context) {
		context.JSON(redfish.NewCollection(
			"/redfish/v1/Chassis",
			"#ChassisCollection.ChassisCollection",
			[]redfish.Link{
				{Oid: "/redfish/v1/Chassis/1"},
				{Oid: "/redfish/v1/Chassis/2"},
			}...,
		))
	})

	odimurl, err := url.Parse(TEST_CONFIG.OdimNBUrl)
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", odimurl.Host)
	if err != nil {
		panic(err)
	}
	go o.app.Run(iris.Listener(httptest.NewLocalTLSListener(l)), iris.WithoutStartupLog)
}

func (o *odimstub) Stop() {
	o.app.Shutdown(stdContext.TODO())
}

func createTestApplication() (*iris.Application, *miniredis.Miniredis) {
	r, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	cm := db.NewConnectionManager(r.Addr(), "")
	return createApplication(&TEST_CONFIG, cm), r
}
