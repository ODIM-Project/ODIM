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

//Package models ...
package models

import "encoding/xml"

//Metadata struct definition
type Metadata struct {
	XMLName      xml.Name `xml:"edmx:Edmx"`
	Xmlnsedmx    string   `xml:"xmlns:edmx,attr"`
	Version      string   `xml:"Version,attr"`
	TopReference []Reference
}

//Reference strcut definition
type Reference struct {
	XMLName    xml.Name  `xml:"edmx:Reference"`
	URI        string    `xml:"Uri,attr"`
	TopInclude []Include `xml:"edmx:Include"`
}

//Include struct definition
type Include struct {
	Namespace string `xml:"Namespace,attr"`
	Alias     string `xml:"Alias,attr,omitempty"`
}
