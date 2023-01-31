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

// Package rpc ...
package rpc

import (
	"context"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	ClientFunc           = services.ODIMService.Client
	NewAccountClientFunc = accountproto.NewAccountClient
)

// DoGetAccountServiceRequest defines the RPC call function for
// the GetAccountService from account-session micro service
func DoGetAccountServiceRequest(ctx context.Context, req accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	account := NewAccountClientFunc(conn)

	resp, err := account.GetAccountServices(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoAccountCreationRequest defines the RPC call function for
// the AccountCreation from account-session micro service
func DoAccountCreationRequest(ctx context.Context, req accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	account := NewAccountClientFunc(conn)

	resp, err := account.Create(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetAllAccountRequest defines the RPC call function for
// the GetAllAccount from account-session micro service
func DoGetAllAccountRequest(ctx context.Context, req accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	account := NewAccountClientFunc(conn)

	resp, err := account.GetAllAccounts(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetAccountRequest defines the RPC call function for
// the GetAccount from account-session micro service
func DoGetAccountRequest(ctx context.Context, req accountproto.GetAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	account := NewAccountClientFunc(conn)

	resp, err := account.GetAccount(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoUpdateAccountRequest defines the RPC call function for
// the UpdateAccount from account-session micro service
func DoUpdateAccountRequest(ctx context.Context, req accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	account := NewAccountClientFunc(conn)

	resp, err := account.Update(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoAccountDeleteRequest defines the RPC call function for
// the AccountDelete from account-session micro service
func DoAccountDeleteRequest(ctx context.Context, req accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	account := NewAccountClientFunc(conn)

	resp, err := account.Delete(ctx, &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}
	defer conn.Close()
	return resp, err
}
