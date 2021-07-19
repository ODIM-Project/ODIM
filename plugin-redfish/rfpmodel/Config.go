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

//Package rfpmodel ...
package rfpmodel

// PluginPrivateKey will contains base64encoded private key of plugin
// this key will be used to decrypt the data.
var PluginPrivateKey []byte

// MetricPropertyData is map to store metric property
var MetricPropertyData = make(map[string]string)
