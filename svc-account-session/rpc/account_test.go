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
package rpc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	accountproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/account"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	"golang.org/x/crypto/sha3"
)

func mockSession(token string, roleID string) error {
	currentTime := time.Now()
	session := asmodel.Session{
		ID:       "id",
		Token:    token,
		UserName: "admin",
		RoleID:   roleID,
		Privileges: map[string]bool{
			common.PrivilegeConfigureUsers: true,
		},
		CreatedTime:  currentTime,
		LastUsedTime: currentTime,
	}
	if err := session.Persist(); err != nil {
		return err
	}
	return nil
}

func createMockUser(username, roleID string) error {
	hash := sha3.New512()
	hash.Write([]byte("P@$$w0rd"))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)
	user := asmodel.User{
		UserName: username,
		Password: hashedPassword,
		RoleID:   roleID,
	}
	if err := asmodel.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func truncateDB(t *testing.T) {
	err := common.TruncateDB(common.OnDisk)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = common.TruncateDB(common.InMemory)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestAccount_Create(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()

	token := "token"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	reqBodyCreateUser, _ := json.Marshal(asmodel.Account{
		UserName: "testUser",
		Password: "Password@123",
		RoleID:   "admin",
	})

	type args struct {
		ctx  context.Context
		req  *accountproto.CreateAccountRequest
		resp *accountproto.AccountResponse
	}
	tests := []struct {
		name    string
		a       *Account
		args    args
		wantErr bool
	}{
		{
			name: "create user",
			a:    &Account{},
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody:  reqBodyCreateUser,
					SessionToken: token,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
		{
			name: "create user with invalid session",
			a:    &Account{},
			args: args{
				req: &accountproto.CreateAccountRequest{
					RequestBody:  reqBodyCreateUser,
					SessionToken: "invalidSession",
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestAccount_GetAllAccounts(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token := "token"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *accountproto.AccountRequest
		resp *accountproto.AccountResponse
	}
	tests := []struct {
		name    string
		a       *Account
		args    args
		wantErr bool
	}{
		{
			name: "get all accounts",
			a:    &Account{},
			args: args{
				req: &accountproto.AccountRequest{
					SessionToken: token,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
		{
			name: "get all accounts with invalid session",
			a:    &Account{},
			args: args{
				req: &accountproto.AccountRequest{
					SessionToken: "invalidSession",
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.GetAllAccounts(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.GetAllAccounts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_GetAccount(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token, accountID := "token", "testUser"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = createMockUser(accountID, common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *accountproto.GetAccountRequest
		resp *accountproto.AccountResponse
	}
	tests := []struct {
		name    string
		a       *Account
		args    args
		wantErr bool
	}{
		{
			name: "get account",
			a:    &Account{},
			args: args{
				req: &accountproto.GetAccountRequest{
					SessionToken: token,
					AccountID:    accountID,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
		{
			name: "get account with invalid session",
			a:    &Account{},
			args: args{
				req: &accountproto.GetAccountRequest{
					SessionToken: "invalidSession",
					AccountID:    accountID,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.GetAccount(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_GetAccountServices(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token := "token"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *accountproto.AccountRequest
		resp *accountproto.AccountResponse
	}
	tests := []struct {
		name    string
		a       *Account
		args    args
		wantErr bool
	}{
		{
			name: "get account service",
			a:    &Account{},
			args: args{
				req: &accountproto.AccountRequest{
					SessionToken: token,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
		{
			name: "get account service with invalid session",
			a:    &Account{},
			args: args{
				req: &accountproto.AccountRequest{
					SessionToken: "invalidSession",
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.GetAccountServices(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.GetAccountServices() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_Update(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token, accountID := "token", "testUser"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = createMockUser(accountID, common.RoleMonitor)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}

	reqBodyRoleIDReadOnly, _ := json.Marshal(asmodel.Account{
		RoleID: common.RoleClient,
	})

	type args struct {
		ctx  context.Context
		req  *accountproto.UpdateAccountRequest
		resp *accountproto.AccountResponse
	}
	tests := []struct {
		name    string
		a       *Account
		args    args
		wantErr bool
	}{
		{
			name: "update account",
			a:    &Account{},
			args: args{
				req: &accountproto.UpdateAccountRequest{
					SessionToken: token,
					AccountID:    accountID,
					RequestBody:  reqBodyRoleIDReadOnly,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
		{
			name: "update account with invalid session",
			a:    &Account{},
			args: args{
				req: &accountproto.UpdateAccountRequest{
					SessionToken: "invalidSession",
					AccountID:    accountID,
					RequestBody:  reqBodyRoleIDReadOnly,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.Update(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_Delete(t *testing.T) {
	defer truncateDB(t)
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()
	token, accountID := "token", "testUser"
	err := mockSession(token, common.RoleClient)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = createMockUser(accountID, common.RoleMonitor)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	type args struct {
		ctx  context.Context
		req  *accountproto.DeleteAccountRequest
		resp *accountproto.AccountResponse
	}
	tests := []struct {
		name    string
		a       *Account
		args    args
		wantErr bool
	}{
		{
			name: "delete account",
			a:    &Account{},
			args: args{
				req: &accountproto.DeleteAccountRequest{
					SessionToken: token,
					AccountID:    accountID,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
		{
			name: "delete account with invalid session",
			a:    &Account{},
			args: args{
				req: &accountproto.DeleteAccountRequest{
					SessionToken: "invalidSession",
					AccountID:    accountID,
				},
				resp: &accountproto.AccountResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.a.Delete(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Account.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
