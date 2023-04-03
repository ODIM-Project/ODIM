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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

var domainIP = "odim.controller.com"
var destinationIP = "odim.destination.com"

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func stubEMBConsume(context.Context, string) {

}

// MockIsAuthorized is for mocking up of authorization
func MockIsAuthorized(ctx context.Context, sessionToken string, privileges, oemPrivileges []string) (response.RPC, error) {
	if sessionToken != "validToken" && sessionToken != "token" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "", nil, nil), nil
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil), nil
}

// MockGetSessionUserName is for mocking up of session user name
func MockGetSessionUserName(ctx context.Context, sessionToken string) (string, error) {
	if sessionToken == "validToken" {
		return "admin", nil
	} else if sessionToken == "token" {
		return "non-admin", nil
	}
	return "", fmt.Errorf("user not found")
}

// MockCreateTask is for mocking up of crete task
func MockCreateTask(ctx context.Context, sessionusername string) (string, error) {
	if sessionusername == "non-admin" {
		return "", fmt.Errorf("no task details")
	}
	return "/redfish/v1/tasks/123", nil
}

// GetEncryptedKey is for mocking up of getting encrypted key
func GetEncryptedKey(key []byte) ([]byte, error) {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		return cryptedKey, fmt.Errorf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey, nil
}

// MockContactClient is for mocking up of contacting client
func MockContactClient(ctx context.Context, url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "https://localhost:1234/ODIM/v1/Subscriptions" {
		if method == http.MethodDelete {
			body := `{"MessageId": "` + response.Success + `"}`
			response := &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			}
			return response, nil
		}
		b := body.(*common.Target)
		if b.DeviceUUID == "d72dade0-c35a-984c-4859-1108132d72da" {
			body := `{"MessageId": "` + response.Failure + `"}`
			return &http.Response{
				StatusCode: http.StatusBadRequest,
				Header: map[string][]string{
					"location": {"/ODIM/v1/Subscriptions/12"},
				},
				Body: ioutil.NopCloser(bytes.NewBufferString(body)),
			}, nil
		}
		body := `{"MessageId": "` + response.Success + `"}`
		response := &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"location": {"/ODIM/v1/Subscriptions/12"},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		response.Header.Set("location", "/ODIM/v1/Subscriptions/12")
		return response, nil
	} else if url == "https://localhost:1234/ODIM/v1/Sessions" {
		body := `{"MessageId": "` + response.Success + `"}`

		r := &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"X-Auth-Token": {"token"},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		return r, nil
	} else if url == "https://10.10.10.23:4321/ODIM/v1/Sessions" || url == "https://10.10.1.6:4321/ODIM/v1/Sessions" {
		body := `{"MessageId": "` + response.Success + `"}`

		r := &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"X-Auth-Token": {"token"},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		return r, nil
	} else if url == "https://10.10.10.23:4321/ODIM/v1/Subscriptions" {
		body := `{"MessageId": "` + response.Failure + `"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://odim.controller.com:1234/ODIM/v1/Subscriptions/123" {
		body := `{"MessageId": "` + response.Success + `"}`
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		return response, nil
	} else if url == "https://localhost:1234/ODIM/v1/Subscriptions/12345" {
		body := `{"MessageId": "` + response.Success + `"}`
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		return response, nil
	} else if url == "https://10.10.1.6:4321/ODIM/v1/Subscriptions" {
		body := `{"MessageId": "` + response.Success + `"}`
		response := &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"location": {"/ODIM/v1/Subscriptions/12345"},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		response.Header.Set("location", "/ODIM/v1/Subscriptions/12345")
		return response, nil
	} else if url == "https://10.10.1.6:4321/ODIM/v1/Subscriptions/12345" {
		body := `{"MessageId": "` + response.Success + `"}`
		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		return response, nil
	}

	return &http.Response{StatusCode: 0}, fmt.Errorf("InvalidRequest")
}

// MockCreateChildTask is for mocking up of crete child task
func MockCreateChildTask(ctx context.Context, sessionID, taskid string) (string, error) {
	return "123456", nil
}

// MockUpdateTask is for mocking up of update task
func MockUpdateTask(context context.Context, task common.TaskData) error {
	return nil
}

// MockGetTarget is for mocking up of getting target info
func MockGetTarget(uuid string) (*common.Target, error) {
	var target *common.Target
	encryptedData, keyErr := GetEncryptedKey([]byte("testData"))
	if keyErr != nil {
		return target, keyErr
	}
	switch uuid {
	case "6d4a0a66-7efa-578e-83cf-44dc68d2874e":
		target = &common.Target{
			ManagerAddress: "100.100.100.100",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			PluginID:       "GRF",
		}
	case "11081de0-4859-984c-c35a-6c50732d72da":
		target = &common.Target{
			ManagerAddress: "10.10.1.3",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "11081de0-4859-984c-c35a-6c50732d72da",
			PluginID:       "ILO",
		}
	case "d72dade0-c35a-984c-4859-1108132d72da":
		target = &common.Target{
			ManagerAddress: "odim.test1.com",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "d72dade0-c35a-984c-4859-1108132d72da",
			PluginID:       "ILO",
		}
	case "110813e0-4859-984c-984c-d72da32d72da":
		target = &common.Target{
			ManagerAddress: domainIP,
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "110813e0-4859-984c-984c-d72da32d72da",
			PluginID:       "ILO",
		}
	case "abab09db-e7a9-4352-8df0-5e41315a2a4c":
		target = &common.Target{
			ManagerAddress: "localhost",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "abab09db-e7a9-4352-8df0-5e41315a2a4c",
			PluginID:       "ILO",
		}
	case "6d4a0a66-7efa-578e-83cf-44dc68d2874d":
		target = &common.Target{
			ManagerAddress: "100.100.100.101",
			Password:       encryptedData,
			UserName:       "admin",
			DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874d",
			PluginID:       "GRF",
		}
	default:
		return target, fmt.Errorf("UUID not found")
	}
	return target, nil
}

// MockGetPluginData is for mocking up of get plugin data of particural plugin id
func MockGetPluginData(pluginID string) (*common.Plugin, *errors.Error) {
	var plugin *common.Plugin
	password, keyErr := GetEncryptedKey([]byte("Password"))
	if keyErr != nil {
		return plugin, errors.PackError(errors.UndefinedErrorType, keyErr.Error())
	}
	switch pluginID {
	case "GRF":
		plugin = &common.Plugin{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "GRF",
		}
	case "ILO":
		plugin = &common.Plugin{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "ILO",
		}
	case "CFM":
		plugin = &common.Plugin{
			IP:                "10.10.1.6",
			Port:              "4321",
			Password:          password,
			Username:          "admin",
			ID:                "CFM",
			PreferredAuthType: "XAuthToken",
			PluginType:        "CFM",
		}
	case "CFMPlugin":
		plugin = &common.Plugin{
			IP:                "10.10.10.23",
			Port:              "4321",
			Password:          password,
			Username:          "admin",
			ID:                "CFMPlugin",
			PreferredAuthType: "XAuthToken",
			PluginType:        "CFMPlugin",
		}
	default:
		return plugin, errors.PackError(errors.UndefinedErrorType, "No data found for the key")
	}
	return plugin, nil
}

// MockGetAllSystems is for mocking up of get all system info
func MockGetAllSystems() ([]string, error) {
	return []string{
		"6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		"11081de0-4859-984c-c35a-6c50732d72da",
		"d72dade0-c35a-984c-4859-1108132d72da",
	}, nil
}

// MockGetSingleSystem is for mocking up of get system info
func MockGetSingleSystem(id string) (string, error) {
	var systemData SavedSystems
	switch id {
	case "6d4a0a66-7efa-578e-83cf-44dc68d2874e":
		systemData = SavedSystems{
			ManagerAddress: "100.100.100.100",
			Password:       []byte("Password"),
			UserName:       "admin",
			DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			PluginID:       "ILO",
		}
	case "11081de0-4859-984c-c35a-6c50732d72da":
		systemData = SavedSystems{
			ManagerAddress: "10.10.1.3",
			Password:       []byte("Password"),
			UserName:       "admin",
			DeviceUUID:     "11081de0-4859-984c-c35a-6c50732d72da",
			PluginID:       "ILO",
		}
	case "d72dade0-c35a-984c-4859-1108132d72da":
		systemData = SavedSystems{
			ManagerAddress: "odim.system.com",
			Password:       []byte("Password"),
			UserName:       "admin",
			DeviceUUID:     "d72dade0-c35a-984c-4859-1108132d72da",
			PluginID:       "GRF",
		}
	default:
		return "", fmt.Errorf("No Data found for the id")
	}
	marshalData, _ := json.Marshal(systemData)
	return string(marshalData), nil
}

// MockGetFabricData is for mocking up of get fabric data against the fabric id
func MockGetFabricData(fabricID string) (evmodel.Fabric, error) {
	var fabric evmodel.Fabric
	if fabricID == "123456" {
		fabric = evmodel.Fabric{
			FabricUUID: "123456",
			PluginID:   "CFM",
		}
	} else if fabricID == "6d4a0a66-7efa-578e-83cf-44dc68d2874e" {
		fabric = evmodel.Fabric{
			FabricUUID: "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			PluginID:   "CFM",
		}
	} else if fabricID == "11081de0-4859-984c-c35a-6c50732d72da" {
		fabric = evmodel.Fabric{
			FabricUUID: "11081de0-4859-984c-c35a-6c50732d72da",
			PluginID:   "CFM1",
		}
	} else if fabricID == "48591de0-4859-1108-c35a-6c50110872da" {
		fabric = evmodel.Fabric{
			FabricUUID: "48591de0-4859-1108-c35a-6c50110872da",
			PluginID:   "CFMPlugin",
		}
	} else {
		return fabric, fmt.Errorf("No data found for the key")
	}
	return fabric, nil
}

// MockGetAggregateDatacData is for mocking up of get aggregate data against the aggregate id
func MockGetAggregateDatacData(aggregateID string) (evmodel.Aggregate, error) {
	aggregate := evmodel.Aggregate{
		Elements: []model.Link{{
			Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
		}},
	}

	return aggregate, nil
}

// MockGetEvtSubscriptions is for mocking up of get event  subscription
func MockGetEvtSubscriptions(searchKey string) ([]evmodel.SubscriptionResource, error) {
	var subarr []evmodel.SubscriptionResource
	switch searchKey {
	case "81de0110-c35a-4859-984c-072d6c5a32d7", "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1", "[^0-9]100.100.100.100[^0-9]":
		subarr = []evmodel.SubscriptionResource{
			{
				UserName:       "admin",
				SubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
				EventDestination: &model.EventDestination{
					Destination:          "https://odim.destination.com:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert", "ResourceAdded"},
					MessageIds:           []string{"IndicatorChanged"},
					ResourceTypes:        []string{"ComputerSystem"},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"}},
					SubordinateResources: true,
				},
				Hosts:    []string{"100.100.100.100"},
				Location: "https://odim.2.com/EventService/Subscriptions/1",
			},
		}
	case "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874d.1", "[^0-9]100.100.100.101[^0-9]":
		subarr = []evmodel.SubscriptionResource{
			{
				UserName:       "admin",
				SubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
				EventDestination: &model.EventDestination{
					Destination:          "https://odim.destination.com:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert", "ResourceAdded"},
					MessageIds:           []string{"IndicatorChanged"},
					ResourceTypes:        []string{"ComputerSystem"},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"}},
					SubordinateResources: true,
				},
				Location: "https://odim.2.com/EventService/Subscriptions/1",

				Hosts: []string{"100.100.100.100", "100.100.100.101"},
			},
		}
	case "11081de0-4859-984c-c35a-6c50732d72da", "/redfish/v1/Systems", "https://odim.destination.com:9090/events", "*", "":
		subarr = []evmodel.SubscriptionResource{
			{
				UserName:       "admin",
				SubscriptionID: "11081de0-4859-984c-c35a-6c50732d72da",
				EventDestination: &model.EventDestination{
					Destination:          "https://odim.destination.com:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert", "StatusChange"},
					MessageIds:           []string{"IndicatorChanged", "StateChanged"},
					ResourceTypes:        []string{"ComputerSystem"},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"}, {Oid: "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1"}},
					SubordinateResources: true,
				},
				Location: "https://odim.2.com/EventService/Subscriptions/1",
				Hosts:    []string{"100.100.100.100", "10.10.1.3"},
			},
		}
	case "71de0110-c35a-4859-984c-072d6c5a32d8", "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1", "[^0-9]10.10.1.3[^0-9]":
		subarr = []evmodel.SubscriptionResource{
			{
				UserName:       "admin",
				SubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d8",
				EventDestination: &model.EventDestination{
					Destination:          "https://localhost:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert", "ResourceAdded"},
					MessageIds:           []string{},
					ResourceTypes:        []string{},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1"}},
					SubordinateResources: false,
				},
				Location: "https://10.10.10.3/EventService/Subscriptions/1",
				Hosts:    []string{"10.10.1.3"},
			},
		}
	case "71de0110-c35a-4859-984c-072d6c5a32d9", "/redfish/v1/Fabrics/123456":
		subarr = []evmodel.SubscriptionResource{
			{
				SubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d9",
				EventDestination: &model.EventDestination{
					Destination:          "https://localhost:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert"},
					MessageIds:           []string{},
					ResourceTypes:        []string{},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Fabrics/123456"}},
					SubordinateResources: true,
				},
				Location: "/ODIM/v1/Subscriptions/12345",
				Hosts:    []string{"localhost"},
			},
		}
	case "5a321010-c35a-4859-984c-072d6c":
		subarr = []evmodel.SubscriptionResource{
			{
				SubscriptionID: "5a321010-c35a-4859-984c-072d6c",
				EventDestination: &model.EventDestination{
					Destination:          "https://localhost:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert", "ResourceAdded"},
					MessageIds:           []string{},
					ResourceTypes:        []string{},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Fabrics/123"}},
					SubordinateResources: true,
				},
				Location: "/ODIM/v1/Subscriptions/123",

				Hosts: []string{domainIP},
			},
		}
	case "71de0110-c35a-4859-984c-072d6c5a3211", "/redfish/v1/Fabrics", "/redfish/v1/Managers", "/redfish/v1/TaskService/Tasks", "/redfish/v1/Chassis":
		subarr = []evmodel.SubscriptionResource{
			{
				SubscriptionID: "71de0110-c35a-4859-984c-072d6c5a3211",
				EventDestination: &model.EventDestination{
					Destination:   "https://localhost:1234/eventsListener",
					Name:          "Subscription",
					Context:       "context",
					EventTypes:    []string{"Alert"},
					MessageIds:    []string{},
					ResourceTypes: []string{},
					OriginResources: []model.Link{
						{Oid: "/redfish/v1/Systems"},
						{Oid: "/redfish/v1/Chassis"},
						{Oid: "/redfish/v1/Fabrics"},
						{Oid: "/redfish/v1/Managers"},
						{Oid: "/redfish/v1/TaskService/Tasks"}},
					SubordinateResources: true,
				},
				Location: "/ODIM/v1/Subscriptions/12345",
				Hosts:    []string{"localhost"},
			},
		}
	case "81de0110-c35a-4859-984c-072d6c5a32d8", "admin":
		subarr = []evmodel.SubscriptionResource{
			{
				UserName:       "admin",
				SubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d8",
				EventDestination: &model.EventDestination{
					Destination:          "https://odim.t.com:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert"},
					MessageIds:           []string{"IndicatorChanged"},
					ResourceTypes:        []string{"ComputerSystem"},
					OriginResources:      []model.Link{{Oid: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"}},
					SubordinateResources: true,
				},
				Location: "https://odim.2.com/EventService/Subscriptions/1",
				Hosts:    []string{"100.100.100.100"},
			},
		}
	default:
		return subarr, fmt.Errorf("No data found for the key")
	}
	return subarr, nil
}

// MockGetDeviceSubscriptions is for mocking up of get device subscription
func MockGetDeviceSubscriptions(hostIP string) (*common.DeviceSubscription, error) {
	var deviceSub *common.DeviceSubscription
	if strings.Contains(hostIP, "100.100.100.100") || hostIP == "*" {
		deviceSub = &common.DeviceSubscription{
			Location:        "https://odim.2.com/EventService/Subscriptions/1",
			EventHostIP:     "100.100.100.100",
			OriginResources: []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"},
		}
	} else if strings.Contains(hostIP, "100.100.100.101") || hostIP == "*" {
		deviceSub = &common.DeviceSubscription{
			Location:        "https://odim.2.com/EventService/Subscriptions/1",
			EventHostIP:     "100.100.100.100",
			OriginResources: []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"},
		}
	} else if strings.Contains(hostIP, "10.10.1.3") {
		deviceSub = &common.DeviceSubscription{
			Location:        "https://10.10.10.3/EventService/Subscriptions/1",
			EventHostIP:     "10.10.1.3",
			OriginResources: []string{"/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1"},
		}
	} else if strings.Contains(hostIP, "odim.ip.com") {
		deviceSub = &common.DeviceSubscription{
			Location:        "/ODIM/v1/Subscriptions/12345",
			EventHostIP:     "odim.ip.com",
			OriginResources: []string{"/redfish/v1/Fabrics/123456"},
		}
	} else if strings.Contains(hostIP, domainIP) {
		deviceSub = &common.DeviceSubscription{
			Location:        "/ODIM/v1/Subscriptions/123",
			EventHostIP:     domainIP,
			OriginResources: []string{"/redfish/v1/Fabrics/123"},
		}
	} else if strings.Contains(hostIP, "localhost") {
		deviceSub = &common.DeviceSubscription{
			Location:        "/ODIM/v1/Subscriptions/12345",
			EventHostIP:     "localhost",
			OriginResources: []string{""},
		}
	} else if strings.Contains(hostIP, "odim.t.com") {
		deviceSub = &common.DeviceSubscription{
			Location:        "https://odim.t.com/EventService/Subscriptions/1",
			EventHostIP:     "odim.t.com",
			OriginResources: []string{"/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72ea.1"},
		}
	} else if hostIP == "*/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1" {
		deviceSub = &common.DeviceSubscription{
			Location:        "https://100.100.100.100/ODIM/v1/Subscriptions/1",
			EventHostIP:     "100.100.100.100",
			OriginResources: []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"},
		}
	} else {
		return deviceSub, fmt.Errorf("No data found for the key")
	}

	return deviceSub, nil
}

// MockSaveEventSubscription is for mocking up of save event subscription
func MockSaveEventSubscription(evtSubscription evmodel.SubscriptionResource) error {
	return nil
}

// MockUpdateEventSubscription is for mocking up of update event subscription
func MockUpdateEventSubscription(evtSubscription evmodel.SubscriptionResource) error {
	return nil
}

// MockDeleteEvtSubscription is for mocking up of delete event subscription
func MockDeleteEvtSubscription(key string) error {
	return nil
}

// MockDeleteDeviceSubscription is for mocking up of delete device subscription
func MockDeleteDeviceSubscription(hostIP string) error {
	return nil
}

// MockUpdateDeviceSubscriptionLocation is for mocking up of updating device subscription based on location
func MockUpdateDeviceSubscriptionLocation(devSubscription common.DeviceSubscription) error {
	return nil
}

// MockGetAllKeysFromTable is for mocking up of get all keys from the given table
func MockGetAllKeysFromTable(table string) ([]string, error) {
	return []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1", "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1"}, nil
}

// MockGetAllFabrics is for mocking up of get all fabric details
func MockGetAllFabrics() ([]string, error) {
	return []string{}, nil
}

// MockGetAllMatchingDetails is for mocking up of get all matching details from the given table
func MockGetAllMatchingDetails(table, pattern string, dbtype common.DbType) ([]string, *errors.Error) {
	return []string{}, nil
}

// MockGetAggregateHosts is for mocking up of get all matching details from the given table
func MockGetAggregateHosts(aggregateIP string) ([]string, error) {
	return []string{}, nil
}

// MockSaveAggregateSubscription is for mocking up of get all matching details from the given table
func MockSaveAggregateSubscription(aggregateID string, hostIP []string) error {
	return nil
}

// MockSaveUndeliveredEvents is for mocking up of save undelivered events
func MockSaveUndeliveredEvents(key string, event []byte) error {
	return nil
}

// MockSaveDeviceSubscription is for mocking up of save undelivered events
func MockSaveDeviceSubscription(common.DeviceSubscription) error {
	return nil
}

// MockGetUndeliveredEvents is for mocking up of get undelivered events
func MockGetUndeliveredEvents(destination string) (string, error) {
	return "", nil
}

// MockGetUndeliveredEventsFlag is for mocking up of getting undelivered events flag
func MockGetUndeliveredEventsFlag(destination string) (bool, error) {
	return true, nil
}

// MockSetUndeliveredEventsFlag is for mocking up of setting undelivered events flag
func MockSetUndeliveredEventsFlag(destination string) error {
	return nil
}

// MockDeleteUndeliveredEventsFlag is for mocking up of deleting undelivered events flag
func MockDeleteUndeliveredEventsFlag(destination string) error {
	return nil
}

// MockDeleteUndeliveredEvents is for mocking up of deleting undelivered events
func MockDeleteUndeliveredEvents(destination string) error {
	return nil
}

// MockContext creates a context for unit test
func MockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}

// MockGetUndeliveredEventsKeyList is for mocking up of get all matching details from the given table
func MockGetUndeliveredEventsKeyList(table, pattern string, dbtype common.DbType, nextCursor int) ([]string, int, *errors.Error) {
	return []string{}, 0, nil
}
