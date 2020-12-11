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

import (
	"strings"
)

// Translator provides ODIMRA|URP API related keywords.
// URP plugin exposes its endpoints(/ODIM/v1/*) in different domain than ODIMRA(/redfish/v1/*),
// this static translator instance translates keywords between domains.
var Translator = &translator{
	r20: map[string]string{
		"redfish": "ODIM",
	},
	o2r: map[string]string{
		"ODIM": "redfish",
	},
}

type translator struct {
	o2r map[string]string
	r20 map[string]string
}

func (u *translator) ODIMToRedfish(data string) string {
	translated := data
	for k, v := range u.o2r {
		translated = strings.Replace(translated, k, v, -1)
	}
	return translated
}

func (u *translator) RedfishToODIM(data string) string {
	translated := data
	for k, v := range u.r20 {
		translated = strings.Replace(translated, k, v, -1)
	}
	return translated
}
