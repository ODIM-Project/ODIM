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
package auth

import (
	"encoding/base64"
	"golang.org/x/crypto/sha3"
	"reflect"
	"testing"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
)

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

func TestCheckSessionCreationCredentials(t *testing.T) {

	Lock.Lock()
	common.SetUpMockConfig()
	Lock.Unlock()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := createMockUser("testUser1", common.RoleAdmin)
	if err != nil {
		t.Fatalf("Error in creating mock admin user %v", err)
	}
	type args struct {
		userName string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *asmodel.User
		wantErr bool
	}{
		{
			name: "successfully checked the credentials",
			args: args{
				userName: "testUser1",
				password: "P@$$w0rd",
			},
			want: &asmodel.User{
				UserName: "testUser1",
				RoleID:   common.RoleAdmin,
				Password: "",
			},
			wantErr: false,
		},
		{
			name: "incorrect credentials",
			args: args{
				userName: "ssss",
				password: "ssss",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckSessionCreationCredentials(tt.args.userName, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSessionCreationCredentials() error = %v, wantErr %v", err.Error(), tt.wantErr)
				return
			}
			if got != nil {
				got.Password = ""
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckSessionCreationCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSessionTimeOut(t *testing.T) {
	Lock.Lock()
	common.SetUpMockConfig()
	Lock.Unlock()
	config.Data.AuthConf.SessionTimeOutInMins = 0.0333333
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sessionToken, id := "validToken", "someID"
	privilegesMap := map[string]bool{
		common.PrivilegeConfigureUsers: true,
	}
	currentTime := time.Now()
	err := createSession(sessionToken, id, privilegesMap, currentTime, currentTime)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	type args struct {
		sessionToken string
	}
	tests := []struct {
		name    string
		args    args
		want    *asmodel.Session
		wantErr bool
	}{
		{
			name: "successfull validation",
			args: args{
				sessionToken: sessionToken,
			},
			want: &asmodel.Session{
				ID:           id,
				Token:        sessionToken,
				Privileges:   privilegesMap,
				CreatedTime:  currentTime,
				LastUsedTime: currentTime,
			},
			wantErr: false,
		},
		{
			name: "expired token",
			args: args{
				sessionToken: sessionToken,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckSessionTimeOut(tt.args.sessionToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckSessionTimeOut() error = %+v, wantErr %v", err, tt.wantErr)
			}
			if got != nil {
				got.CreatedTime = currentTime
				got.LastUsedTime = currentTime
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckSessionTimeOut() = %#v, want %#v", got, tt.want)
			}
		})
		time.Sleep(4 * time.Second)
	}
}
