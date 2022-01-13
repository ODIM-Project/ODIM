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

package model

// Collection is the redfish resource collection  model according to the 2020.3 release
type Collection struct {
	ODataContext         string      `json:"@odata.context,omitempty"`
	ODataEtag            string      `json:"@odata.etag,omitempty"`
	ODataID              string      `json:"@odata.id"`
	ODataType            string      `json:"@odata.type"`
	Description          string      `json:"Description,omitempty"`
	Name                 string      `json:"Name"`
	Members              []*Link     `json:"Members"`
	MembersCount         int         `json:"Members@odata.count"`
	MemberNavigationLink string      `json:"Members@odata.navigationLink,omitempty"`
	Oem                  interface{} `json:"Oem,omitempty"`
	MembersNextLink      string      `json:"Members@odata.nextLink,omitempty"`
}
