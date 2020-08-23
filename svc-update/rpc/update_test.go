// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	updateproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/update"
	"github.com/ODIM-Project/ODIM/svc-update/update"
)

func TestUpdate_GetUpdateService(t *testing.T) {
	config.SetUpMockConfig(t)
	config.Data.EnabledServices = append(config.Data.EnabledServices, "UpdateService")
	type args struct {
		ctx  context.Context
		req  *updateproto.UpdateRequest
		resp *updateproto.UpdateResponse
	}
	tests := []struct {
		name    string
		a       *Updater
		args    args
		wantErr bool
	}{
		{
			name: "positive GetAggregationService",
			a:    &Updater{connector: connector},
			args: args{
				req:  &updateproto.UpdateRequest{SessionToken: "validToken"},
				resp: &updateproto.UpdateResponse{},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    &Update{connector: connector},
			args: args{
				req:  &updateproto.UpdateRequest{SessionToken: "invalidToken"},
				resp: &updateproto.UpdateResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.GetUpdateService(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Update.GetUpdateService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
