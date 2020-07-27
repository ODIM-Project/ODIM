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
package session

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/svc-account-session/auth"
	"testing"
)

func TestUpdateLastUsedTime(t *testing.T) {
	auth.Lock.Lock()
	common.SetUpMockConfig()
	auth.Lock.Unlock()

	_, token := createSession(t, common.RoleAdmin, "admin", []string{common.PrivilegeConfigureUsers, common.PrivilegeLogin})
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully updated session last used time",
			args: args{
				token: token,
			},
			wantErr: false,
		},
		{
			name: "update session last used time with invalid token",
			args: args{
				token: "invalid token",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateLastUsedTime(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("UpdateLastUsedTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	err := common.TruncateDB(common.OnDisk)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	err = common.TruncateDB(common.InMemory)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
