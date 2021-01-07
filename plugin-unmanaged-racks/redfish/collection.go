/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package redfish

// Collection structs represents generic Redfish collection asset
type Collection struct {
	OdataContext string `json:"@odata.context"`
	Etag         string `json:"@odata.etag,omitempty"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	Description  string `json:"Description"`
	Name         string `json:"Name"`
	Members      []Link `json:"Members"`
	MembersCount int    `json:"Members@odata.count"`
}

// NewCollection created new instance of Collection
func NewCollection(odataID, odataType string, members ...Link) Collection {
	return Collection{
		OdataContext: "/ODIM/v1/$metadata" + odataType,
		OdataID:      odataID,
		OdataType:    odataType,
		Members:      append([]Link{}, members...),
		MembersCount: len(members),
	}
}
