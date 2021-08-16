//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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

//Package sresponse ...
package sresponse

import (
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
)

//Collection holds response  of System and ChassisCollection
type Collection struct {
	OdataContext string      `json:"@odata.context"`
	Etag         string      `json:"@odata.etag,omitempty"`
	OdataID      string      `json:"@odata.id"`
	OdataType    string      `json:"@odata.type"`
	Description  string      `json:"Description"`
	Name         string      `json:"Name"`
	Members      []dmtf.Link `json:"Members"`
	MembersCount int         `json:"Members@odata.count"`
}

func (c *Collection) AddMember(m dmtf.Link) {
	c.Members = append(c.Members, m)
	c.MembersCount = len(c.Members)
}

func NewChassisCollection() Collection {
	return Collection{
		OdataContext: "/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
		OdataID:      "/redfish/v1/Chassis",
		OdataType:    "#ChassisCollection.ChassisCollection",
		Description:  "Chassis Collection",
		Name:         "Chassis Collection",
		Members:      []dmtf.Link{},
	}
}
