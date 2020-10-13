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

//Package rpc ...
package rpc

import (
	"context"
	"fmt"

	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

//GetChassisCollection will do the rpc call to collect all chassis
func GetChassisCollection(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	asService := chassisproto.NewChassisService(services.Systems, services.Service.Client())
	resp, err := asService.GetChassisCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

//GetChassisResource will do the rpc call to collect Chassis Resource
func GetChassisResource(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	asService := chassisproto.NewChassisService(services.Systems, services.Service.Client())
	resp, err := asService.GetChassisResource(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	return resp, nil
}

//GetChassis will do the rpc call to collect System Resource
func GetChassis(req chassisproto.GetChassisRequest) (*chassisproto.GetChassisResponse, error) {
	asService := chassisproto.NewChassisService(services.Systems, services.Service.Client())
	resp, err := asService.GetChassisInfo(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

func CreateChassis(req chassisproto.CreateChassisRequest) (*chassisproto.GetChassisResponse, error) {
	service := chassisproto.NewChassisService(services.Systems, services.Service.Client())
	resp, err := service.CreateChassis(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}

func DeleteChassis(req chassisproto.DeleteChassisRequest) (*chassisproto.GetChassisResponse, error) {
	service := chassisproto.NewChassisService(services.Systems, services.Service.Client())
	resp, err := service.DeleteChassis(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	return resp, nil
}
