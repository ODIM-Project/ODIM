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

type ValidationResult []MsgExtendedInfo

func (ve *ValidationResult) Error() *CommonError {
	e := NewError()
	for _, info := range *ve {
		e.AddExtendedInfo(info)
	}
	return e
}

func (ve *ValidationResult) HasErrors() bool {
	return len(*ve) > 0
}

type CompositeValidator []Validator

func (v *CompositeValidator) Validate() ValidationResult {
	vr := ValidationResult{}
	for _, validationStep := range *v {
		if validationStep.ValidationRule() {
			vr = append(vr, validationStep.ErrorGenerator())
		}
	}
	return vr
}

type ValidationRule func() bool
type ErrorGenerator func() MsgExtendedInfo

type Validator struct {
	ValidationRule
	ErrorGenerator
}
