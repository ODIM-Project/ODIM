// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
// (C) Copyright 2020 Intel Corporation
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
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/chassis"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/stretchr/testify/assert"
)

func mockResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", table, err.Error())
	}
	return nil
}
func mockIsAuthorized(ctx context.Context, sessionToken string, privileges, oemPrivileges []string) (response.RPC, error) {
	if sessionToken != "validToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil), nil
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil), nil
}

func TestChassisRPC_GetChassisResource(t *testing.T) {
	config.SetUpMockConfig(t)
	common.SetUpMockConfig()
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
	reqData := []byte(`\"@odata.id\":\"/redfish/v1/Chassis/1/Power\"`)
	err := mockResourceData(reqData, "Power", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1/Power")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	cha := new(ChassisRPC)
	cha.IsAuthorizedRPC = mockIsAuthorized
	type args struct {
		ctx  context.Context
		req  *chassisproto.GetChassisRequest
		resp *chassisproto.GetChassisResponse
	}
	tests := []struct {
		name    string
		cha     *ChassisRPC
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1/Power",
					SessionToken: "validToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1/Power",
					SessionToken: "invalidToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid url",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					RequestParam: "6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1/Power1",
					SessionToken: "validToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.cha.GetChassisResource(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ChassisRPC.GetChassisResource() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChassis_GetAllChassis(t *testing.T) {
	cha := NewChassisRPC(
		mockIsAuthorized,
		nil,
		chassis.NewGetCollectionHandler(
			func(name string) (plugin.Client, *errors.Error) {
				return nil, errors.PackError(errors.DBKeyNotFound, "error")
			}, func(table string) ([]string, error) {
				return []string{}, nil
			}), nil, nil, nil)

	type args struct {
		ctx  context.Context
		req  *chassisproto.GetChassisRequest
		resp *chassisproto.GetChassisResponse
	}
	tests := []struct {
		name    string
		cha     *ChassisRPC
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					URL:          "/redfish/v1/Chassis",
					SessionToken: "validToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					URL:          "/redfish/v1/Chassis",
					SessionToken: "invalidToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.cha.GetChassisCollection(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ChassisRPC.GetChassisCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChassis_GetResourceInfo(t *testing.T) {
	common.SetUpMockConfig()
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
	reqData := []byte(`\"@odata.id\":\"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1\"`)
	err := mockResourceData(reqData, "chassis", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	cha := new(ChassisRPC)
	cha.IsAuthorizedRPC = mockIsAuthorized
	cha.GetHandler = chassis.NewGetHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.DBKeyNotFound, "urp os not registered")
		}, smodel.Find)
	type args struct {
		ctx  context.Context
		req  *chassisproto.GetChassisRequest
		resp *chassisproto.GetChassisResponse
	}
	tests := []struct {
		name    string
		cha     *ChassisRPC
		args    args
		wantErr bool
	}{
		{
			name: "Request with valid token",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
					SessionToken: "validToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
		{
			name: "Request with invalid token",
			cha:  cha,
			args: args{
				ctx: context.Background(),
				req: &chassisproto.GetChassisRequest{
					URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
					SessionToken: "invalidToken",
				},
				resp: &chassisproto.GetChassisResponse{},
			}, wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.cha.GetChassisInfo(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ChassisRPC.GetChassisInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChassisRPC_UpdateChassis(t *testing.T) {
	ctx := mockContext()
	common.SetUpMockConfig()
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
	reqData := []byte(`\"@odata.id\":\"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1\"`)
	err := mockResourceData(reqData, "chassis", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	cha := new(ChassisRPC)
	cha.IsAuthorizedRPC = mockIsAuthorized
	cha.UpdateHandler = chassis.NewUpdateHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.DBKeyNotFound, "urp os not registered")
		})

	req := chassisproto.UpdateChassisRequest{
		URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
		SessionToken: "validToken",
	}
	_, err = cha.UpdateChassis(ctx, &req)
	assert.Nil(t, err, "There should be no error")

}

func TestChassisRPC_DeleteChassis(t *testing.T) {
	ctx := mockContext()
	common.SetUpMockConfig()
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
	reqData := []byte(`\"@odata.id\":\"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1\"`)
	err := mockResourceData(reqData, "chassis", "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")
	if err != nil {
		t.Fatalf("Error in creating mock resource data :%v", err)
	}
	cha := new(ChassisRPC)
	cha.IsAuthorizedRPC = mockIsAuthorized
	cha.DeleteHandler = chassis.NewDeleteHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.DBKeyNotFound, "urp os not registered")
		}, smodel.Find)

	req := chassisproto.DeleteChassisRequest{
		URL:          "/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
		SessionToken: "validToken",
	}
	_, err = cha.DeleteChassis(ctx, &req)
	assert.Nil(t, err, "There should be no error")

}

func TestChassisRPC_CreateChassis(t *testing.T) {
	ctx := mockContext()
	common.SetUpMockConfig()
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
	reqData := []byte(`\"@odata.id\":\"/redfish/v1/Chassis/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1\"`)

	cha := new(ChassisRPC)
	cha.IsAuthorizedRPC = mockIsAuthorized
	cha.CreateHandler = chassis.NewCreateHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.DBKeyNotFound, "urp os not registered")
		})

	req := chassisproto.CreateChassisRequest{
		RequestBody:  reqData,
		SessionToken: "validToken",
	}
	_, err := cha.CreateChassis(ctx, &req)
	assert.Nil(t, err, "There should be no error")

}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}

func Test_jsonMarshal(t *testing.T) {
	ctx := mockContext()
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, &errors.Error{}
	}
	generateResponse(context.Background(), "dummy")
	jsonMarshal(ctx, "dummy")
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	}
}
