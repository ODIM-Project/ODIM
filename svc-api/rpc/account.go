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
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// DoGetAccountServiceRequest defines the RPC call function for
// the GetAccountService from account-session micro service
func DoGetAccountServiceRequest(req accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	account := accountproto.NewAccountClient(conn)

	resp, err := account.GetAccountServices(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// DoAccountCreationRequest defines the RPC call function for
// the AccountCreation from account-session micro service
func DoAccountCreationRequest(req accountproto.CreateAccountRequest) (*accountproto.AccountResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	account := accountproto.NewAccountClient(conn)

	resp, err := account.Create(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// DoGetAllAccountRequest defines the RPC call function for
// the GetAllAccount from account-session micro service
func DoGetAllAccountRequest(req accountproto.AccountRequest) (*accountproto.AccountResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	account := accountproto.NewAccountClient(conn)

	resp, err := account.GetAllAccounts(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// DoGetAccountRequest defines the RPC call function for
// the GetAccount from account-session micro service
func DoGetAccountRequest(req accountproto.GetAccountRequest) (*accountproto.AccountResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	account := accountproto.NewAccountClient(conn)

	resp, err := account.GetAccount(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// DoUpdateAccountRequest defines the RPC call function for
// the UpdateAccount from account-session micro service
func DoUpdateAccountRequest(req accountproto.UpdateAccountRequest) (*accountproto.AccountResponse, error) {
	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	account := accountproto.NewAccountClient(conn)

	resp, err := account.Update(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}

// DoAccountDeleteRequest defines the RPC call function for
// the AccountDelete from account-session micro service
func DoAccountDeleteRequest(req accountproto.DeleteAccountRequest) (*accountproto.AccountResponse, error) {

	conn, err := services.ODIMService.Client(services.AccountSession)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	account := accountproto.NewAccountClient(conn)

	resp, err := account.Delete(context.TODO(), &req)
	if err != nil && resp == nil {
		return nil, fmt.Errorf("error: something went wrong with rpc call: %v", err)
	}

	return resp, err
}
