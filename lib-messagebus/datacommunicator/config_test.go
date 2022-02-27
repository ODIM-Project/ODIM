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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var cwdDir, _ = os.Getwd()

const sampleFileContent = `[KAFKA]
KServersInfo   = ["someServer:0"]
KTimeout   = 0
KAFKACertFile       = "someCertFile"
KAFKAKeyFile        = "someKeyFlie"
KAFKACAFile         = "someCAFile"`

func createFile(t *testing.T, fName, fContent string) {
	if err := ioutil.WriteFile(fName, []byte(fContent), 0644); err != nil {
		t.Fatal("error :failed to create a sample file for tests:", err)
	}
}

func TestSetConfiguration(t *testing.T) {
	sampleConfigFile := filepath.Join(cwdDir, "sample.toml")
	createFile(t, sampleConfigFile, sampleFileContent)
	defer os.Remove(sampleConfigFile)

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{filePath: sampleConfigFile},
			wantErr: false,
		},
		{
			name:    "no file",
			args:    args{filePath: "some/random/file/path"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetConfiguration(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("SetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
