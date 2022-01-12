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
package datacommunicator

import (
	"reflect"
	"testing"
)

func TestCommunicator(t *testing.T) {
	type args struct {
		bt               string
		mqConfigFilePath string
		pipe             string
	}
	tests := []struct {
		name    string
		args    args
		want    MQBus
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Communicator(tt.args.bt, tt.args.mqConfigFilePath, tt.args.pipe)
			if (err != nil) != tt.wantErr {
				t.Errorf("Communicator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Communicator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		d interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		d []byte
		a interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.d, tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
