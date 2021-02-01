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

package agresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// ListMember define the links for each account present in odimra
type ListMember struct {
	OdataID string `json:"@odata.id"`
}

// List defines the collection of resources like accounts, sessions, roles etc in svc-aggregation
type List struct {
	response.Response
	MembersCount int          `json:"Members@odata.count"`
	Members      []ListMember `json:"Members"`
}
