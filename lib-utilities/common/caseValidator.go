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

package common

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RequestParamsCaseValidator function validates input json w.r.t the structure defined, returns all the invalid properties sent in an input
func RequestParamsCaseValidator(rawRequestBody []byte, reqStruct interface{}) (string, error) {
	invalidProperties := ""
	// Converting structure into a map
	reqStructMap, err := convertStructToMap(reqStruct)
	if err != nil {
		return "", fmt.Errorf("error while converting struct to map in param case validation: %v", err)
	}

	//Unmarshalling request body into a map
	var rawReqBodyMap map[string]interface{}
	json.Unmarshal([]byte(rawRequestBody), &rawReqBodyMap)
	for key := range rawReqBodyMap {
		if _, found := reqStructMap[key]; !found {
			invalidProperties += getInvalidProperty(reqStructMap, key)
		} else {
			invalidCase, err := reqValidator(rawReqBodyMap[key], reqStructMap, key)
			if err != nil {
				return "", fmt.Errorf("error while trying to validate the request parameters case: %v", err)
			}
			invalidProperties += invalidCase
		}
	}
	return invalidProperties, nil
}

// Validates input request w.r.t. the structure, this is a recursive function,
// Based on the type, request parameters are compared and validates for invalid property
func reqValidator(rawReqBodyMap interface{}, reqStructMap map[string]interface{}, key string) (string, error) {
	var invalidProperties string
	switch v := rawReqBodyMap.(type) {
	case map[string]interface{}:
		if _, found := reqStructMap[key]; !found && key != "" {
			invalidProperties += getInvalidProperty(reqStructMap, key)
		} else {
			for nestedkey, nestedVal := range v {
				if key != "" {
					invalidCase, err := reqValidator(nestedVal, reqStructMap[key].(map[string]interface{}), nestedkey)
					if err != nil {
						return "", fmt.Errorf("error while trying to validate the request parameters case in map: %v", err)
					}
					invalidProperties += invalidCase
				} else {
					invalidProperties += getInvalidProperty(reqStructMap, nestedkey)
				}
			}
		}
	case []interface{}:
		if _, found := reqStructMap[key]; !found {
			invalidProperties += getInvalidProperty(reqStructMap, key)
		} else {
			invalidCase, err := reqSliceValidator(v, reqStructMap, key)
			if err != nil {
				return "", fmt.Errorf("error while trying to validate the request parameters case in slice: %v", err)
			}
			invalidProperties += invalidCase
		}
	default:
		invalidProperties += getInvalidProperty(reqStructMap, key)
	}
	return invalidProperties, nil
}

// If a type is slice, this function will be called
func reqSliceValidator(rawReqBodyMap []interface{}, reqStructMap map[string]interface{}, key string) (string, error) {
	var invalidProperties string
	mapData := reqStructMap
	for i, val := range rawReqBodyMap {
		if _, mapFound := val.(map[string]interface{}); mapFound {
			if reqStructMap[key] != nil {
				mapData = reqStructMap[key].([]interface{})[i].(map[string]interface{})
			}
			invalidCase, err := reqValidator(val, mapData, "")
			if err != nil {
				return "", fmt.Errorf("error while trying to validate the request parameters case in slice: %v", err)
			}
			invalidProperties += invalidCase
		}
	}
	return invalidProperties, nil
}

// Checks if a passed key is an invalid property or not on comparision with the structure passed
func getInvalidProperty(reqStructMap map[string]interface{}, key string) string {
	invalidProperties := ""
	if _, found := reqStructMap[key]; !found {
		keyFoundWithIgnoreCase := searchIgnoreCase(reqStructMap, key)
		if keyFoundWithIgnoreCase {
			invalidProperties += key + " "
		}
	}
	return invalidProperties
}

// Converts struct to map
func convertStructToMap(reqStructure interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(reqStructure)
	if err != nil {
		return nil, fmt.Errorf("error while trying to marshal the request structure: %v", err)
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, fmt.Errorf("error while trying to unmarshal the request structure: %v", err)
	}
	return mapData, nil
}

// Searches the key in map with ignore case
func searchIgnoreCase(reqMap map[string]interface{}, searchKey string) bool {
	for key := range reqMap {
		if strings.EqualFold(key, searchKey) {
			return true
		}
	}
	return false
}
