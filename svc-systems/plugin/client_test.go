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

package plugin

import (
	"fmt"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"

	"github.com/stretchr/testify/assert"
)

func Test_uriTranslator(t *testing.T) {

	sut := uriTranslator{&config.URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}}

	tests := []struct {
		name           string
		translate      func(toBeTranslated string) (translated string)
		toBeTranslated string
		expected       string
	}{
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "", expected: ""},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "redfish", expected: "ODIM"},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "redfish redfish", expected: "ODIM ODIM"},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "Redfish", expected: "Redfish"},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: `{"@odata.id":"/redfish/v1"}`, expected: `{"@odata.id":"/ODIM/v1"}`},

		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "", expected: ""},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "ODIM", expected: "redfish"},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "ODIM ODIM", expected: "redfish redfish"},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "Redfish", expected: "Redfish"},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: `{"@odata.id":"/ODIM/v1"}`, expected: `{"@odata.id":"/redfish/v1"}`},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("translating %s(%s)", test.name, test.toBeTranslated), func(t *testing.T) {
			assert.Equal(t, test.expected, test.translate(test.toBeTranslated))
		})
	}
}
