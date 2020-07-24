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

//Package response ...
package response

//MessageRegistryFileID defines the message registry file id
type MessageRegistryFileID struct {
	ID           string     `json:"Id"`
	OdataContext string     `json:"@odata.context"`
	Etag         string     `json:"@odata.etag,omitempty"`
	OdataID      string     `json:"@odata.id"`
	OdataType    string     `json:"@odata.type"`
	Name         string     `json:"Name"`
	Description  string     `json:"Description"`
	Languages    []string   `json:"Languages"`
	Location     []Location `json:"Location"`
	Registry     string     `json:"Registry"`
}

//Location defines the locations/Paths of the file to retrieve.
type Location struct {
	ArchiveFile    string `json:"ArchiveFile,omitempty"`
	ArchiveURI     string `json:"ArchiveUri,omitempty"`
	Language       string `json:"Language"`
	PublicationURI string `json:"PublicationUri,omitempty"`
	URI            string `json:"Uri"`
}
