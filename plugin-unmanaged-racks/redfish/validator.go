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

// Validator is an interface of validator. When `Validate()` is called and
// validator detects problem, non empty array containing error message(s) is returned.
type Validator interface {
	Validate() (violation *MsgExtendedInfo, errorCode *int)
}

type isViolated func() bool
type generateError func() (violation MsgExtendedInfo, errorCode int)

// NewValidator constructs new instance of `Validator`.
// `isViolated` function returns true if validation rule is violated.
// `generateError` function will be called only if `isViolated' will return true.
func NewValidator(isViolated isViolated, generateError generateError) Validator {
	return &validator{
		isViolated:    isViolated,
		generateError: generateError,
	}
}

type validator struct {
	isViolated
	generateError
}

func (v *validator) Validate() (*MsgExtendedInfo, *int) {
	if v.isViolated() {
		violation, errCode := v.generateError()
		return &violation, &errCode
	}
	return nil, nil
}

// CompositeValidator constructs multiple validators validator.
// Execution of 'Validate` function verifies all provided validation rules and returns all reported problems.
func CompositeValidator(v ...Validator) Validator {
	cv := compositeValidator(v)
	return &cv
}

type compositeValidator []Validator

func (v *compositeValidator) Validate() (*MsgExtendedInfo, *int) {
	for _, validator := range *v {
		if violation, errCode := validator.Validate(); violation != nil {
			return violation, errCode
		}
	}
	return nil, nil
}
