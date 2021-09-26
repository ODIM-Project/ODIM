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

//Package systems ...
package systems

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"

	"regexp"
	"strconv"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	systemsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/systems"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-systems/scommon"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
)

func setRegexFlag(val string) bool {
	var re = regexp.MustCompile(`(?m)[\[\]!@#$%^&*(),.?":{}|<>]`)

	for i, match := range re.FindAllString(val, -1) {
		log.Info("Matched entry no.: " + string(i) + " match=" + match)
		return true
	}
	return false
}

func errorResp(each string, resp response.RPC) (response.RPC, error) {
	if each == "" {
		errorMessage := " not a valid search/filter expression"
		return common.GeneralError(http.StatusBadRequest, response.QueryCombinationInvalid, errorMessage, []interface{}{"ComputerSystem", ""}, nil), fmt.Errorf(errorMessage)
	}
	return resp, nil
}

func ifMatches(strPara string, operator string, resp response.RPC) (response.RPC, error) {
	strSplit := strings.Split(strPara, operator)
	for _, e := range strSplit {

		resp, err := errorResp(e, resp)
		if err != nil {
			log.Error(err.Error())
			return resp, err
		}
	}
	return resp, nil
}

func validate(strPara string, resp response.RPC) (response.RPC, error) {
	if strings.Contains(strPara, "and") {
		resp, err := ifMatches(strPara, "and", resp)
		if err != nil {
			log.Error(err.Error())
			return resp, err
		}
	}
	if strings.Contains(strPara, "or") {
		resp, err := ifMatches(strPara, "or", resp)
		if err != nil {
			log.Error(err.Error())
			return resp, err
		}
	}
	return resp, nil
}

// validateLastParameter checks whether last parameter in the expression
// is an operator or not. It throughs an error if the last parameter is operator.
// It also checks whether the expression is empty eg: /redfish/v1/Systems?$filter=%20
func validateLastParameter(expression []string) error {
	var lastParam string
	for i := len(expression) - 1; i >= 0; i-- {
		if expression[i] != "" {
			lastParam = expression[i]
			break
		}
	}
	if lastParam == "" {
		return fmt.Errorf("no valid expression found")
	}
	operatorSet := []string{"eq", "ne", "gt", "ge", "lt", "le"}
	for _, op := range operatorSet {
		if lastParam == op {
			return fmt.Errorf("expression ends with an operator %v", op)
		}
	}
	return nil
}

//GetMembers will fetch the resource members based on the filter expression
func GetMembers(allowed map[string]map[string]bool, expression []string, resp response.RPC) ([]dmtf.Link, response.RPC, error) {
	err := validateLastParameter(expression)
	if err != nil {
		return nil, common.GeneralError(http.StatusBadRequest, response.QueryNotSupported, err.Error(), nil, nil), err
	}

	allowed["searchKeys"] = make(map[string]bool)
	allowed["conditionKeys"] = make(map[string]bool)
	for _, value := range scommon.SF.SearchKeys {
		for k := range value {
			allowed["searchKeys"][k] = true
		}
	}
	for _, value := range scommon.SF.ConditionKeys {
		allowed["conditionKeys"][value] = true
	}
	var members []dmtf.Link
	var regexFlag, typeFlag, arrayFlag bool
	for amp, pam := range expression {
		for _, value := range scommon.SF.SearchKeys {
			for k, v := range value {
				if k == pam {
					if v["type"] != "string" && v["type"] != "[]string" {
						typeFlag = true
					}
					if v["type"] == "[]float64" || v["type"] == "[]int" {
						arrayFlag = true
					}
				}
			}
		}
		if allowed["searchKeys"][pam] {
			key := strings.Replace(pam, "/", "\\/", -1)
			new := amp + 1
			if new < len(expression) {
				var val, regex string
				if allowed["conditionKeys"][expression[new]] {
					val = expression[new+1]
					regexFlag = setRegexFlag(val)
					key = strings.Replace(pam, "\\/", "/", -1)
					if regexFlag {
						regex = val
						// regular expression flag is true then get all data for key depending on the type
						var list []string
						var err error
						if arrayFlag {
							list, err = smodel.GetStorageList(key, "ne", 0, true)
							if err != nil {
								return nil, common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(), []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil), err
							}
						} else {
							if !typeFlag {
								list, err = getStringData(key, "", "eq", true)
								if err != nil {
									return nil, common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(), []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil), err
								}
							} else {
								list, err = getRangeData(key, "ge", 0, true)
								if err != nil {
									return nil, common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(), []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil), err
								}
							}
						}
						// parse the data with the regex
						parsedList, err := parseRegexData(list, regex)
						if err != nil {
							errorMessage := " not a valid search/filter expression"
							return nil, common.GeneralError(http.StatusBadRequest, response.QueryCombinationInvalid, errorMessage, []interface{}{"ComputerSystem", ""}, nil), fmt.Errorf(errorMessage)
						}
						for i := 0; i < len(parsedList); i++ {
							members = append(members, dmtf.Link{Oid: parsedList[i]})
						}
					} else if arrayFlag {
						searchValue, err := strconv.ParseFloat(val, 64)
						if err != nil {
							return nil, common.GeneralError(http.StatusBadRequest, response.PropertyValueFormatError, err.Error(), []interface{}{key, "Invalida value"}, nil), err
						}
						list, err := smodel.GetStorageList(key, expression[new], searchValue, false)
						if err != nil {
							return nil, common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(), []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil), err

						}
						for i := 0; i < len(list); i++ {
							members = append(members, dmtf.Link{Oid: list[i]})
						}
					} else {
						// to get the members for type string
						if !typeFlag {
							// rejecting the request if expression is not eq or ne
							if !(expression[new] == "eq" || expression[new] == "ne") {
								return nil, common.GeneralError(http.StatusBadRequest, response.QueryCombinationInvalid, "error:invalid expression", []interface{}{expression[new], "Invalid Expression for " + key}, nil), fmt.Errorf("error:invalid expression")
							}
							list, err := getStringData(key, val, expression[new], false)
							if err != nil {
								return nil, common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(), []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil), err
							}
							for i := 0; i < len(list); i++ {
								members = append(members, dmtf.Link{Oid: list[i]})
							}
						} else {
							//validate the value
							searchValue, err := strconv.Atoi(val)
							if err != nil {
								return nil, common.GeneralError(http.StatusBadRequest, response.PropertyValueFormatError, err.Error(), []interface{}{key, "Invalida value"}, nil), err
							}
							list, err := getRangeData(key, expression[new], searchValue, false)
							if err != nil {
								return nil, common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, err.Error(), []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil), err
							}
							for i := 0; i < len(list); i++ {
								members = append(members, dmtf.Link{Oid: list[i]})
							}
						}

					}

				}

			}

		}
	}
	return members, resp, nil
}

//getAllSystemIDs will fetch all the document ID's present in the DB
func getAllSystemIDs(resp response.RPC) ([]dmtf.Link, response.RPC, error) {
	var mems []dmtf.Link
	systemKeys, err := smodel.GetAllKeysFromTable("ComputerSystem")
	if err != nil {
		log.Error("error getting all keys of systemcollection table : " + err.Error())
		errorMessage := err.Error()
		if errorMessage == "error while trying to get resource details: no data with the with table name SystemCollection found" {
			return nil, common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"ComputerSystem", ""}, nil), err
		}
		return nil, common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil), err
	}
	for _, key := range systemKeys {
		mems = append(mems, dmtf.Link{Oid: key})
	}

	return mems, resp, nil

}

//LogicalOperation will take slice of slices of member collection and logical operation to perform
//return 1 slice of member collection with logical operation performed on them
func LogicalOperation(interm [][]dmtf.Link, lo string) []dmtf.Link {
	var respMembers []dmtf.Link
	if lo == "and" {
		if len(interm[0]) < len(interm[1]) {
			for i := 0; i < len(interm[0]); i++ {
				for _, v := range interm[1] {
					if interm[0][i] == v {
						respMembers = append(respMembers, v)
					}
				}
			}
		} else {
			for i := 0; i < len(interm[1]); i++ {
				for _, v := range interm[0] {
					if interm[1][i] == v {
						respMembers = append(respMembers, v)
					}
				}
			}
		}

	} else if lo == "or" {
		encountered := map[dmtf.Link]bool{}
		var tmp []dmtf.Link
		tmp = append(interm[0], interm[1]...)
		for v := range tmp {
			encountered[tmp[v]] = true
		}
		for key := range encountered {
			respMembers = append(respMembers, key)
		}
	}

	return respMembers
}

//SearchAndFilter take the url as input and return the search result based on the filter expression
func SearchAndFilter(paramStr []string, resp response.RPC) (response.RPC, error) {
	allowed := make(map[string]map[string]bool)
	allowed["queryKeys"] = make(map[string]bool)
	allowed["logicalOperators"] = make(map[string]bool)

	for _, value := range scommon.SF.QueryKeys {
		value = "$" + value
		allowed["queryKeys"][value] = true
	}
	query := strings.Split(paramStr[1], "=")
	if len(query) < 2 {
		errorMessage := " not a valid search/filter expression"
		return common.GeneralError(http.StatusBadRequest, response.QueryCombinationInvalid, errorMessage, []interface{}{"ComputerSystem", ""}, nil), fmt.Errorf(errorMessage)
	}
	if !allowed["queryKeys"][query[0]] {
		errorMessage := " not a valid search/filter expression"
		return common.GeneralError(http.StatusBadRequest, response.QueryCombinationInvalid, errorMessage, []interface{}{"ComputerSystem", ""}, nil), fmt.Errorf(errorMessage)
	}
	percentSpl := strings.Split(query[1], "%")
	var strPara string
	if len(percentSpl) > 1 {
		for i, j := range percentSpl {
			if i == 0 {
				strPara += j
			} else {
				bl, _ := hex.DecodeString(j[:2])
				strung := string(bl)
				strPara += strung
				strPara += j[2:]
			}
		}
	} else {
		strPara = paramStr[1]
	}
	var respMembers []dmtf.Link
	var err error

	if checkParentheses(strPara) {
		if strings.Count(strPara, "(") != strings.Count(strPara, ")") {
			errorMessage := " not a valid search/filter expression"
			return common.GeneralError(http.StatusBadRequest, response.QueryCombinationInvalid, errorMessage, []interface{}{"ComputerSystem", ""}, nil), fmt.Errorf(errorMessage)
		}

		var sa []dmtf.Link
		re := regexp.MustCompile(`\((.*?)\)`)
		submatchall := re.FindAllString(strPara, -1)
		for _, element := range submatchall {
			element = strings.Trim(element, "(")
			element = strings.Trim(element, ")")
			fmt.Println(element)
			if strings.Contains(element, " and ") || strings.Contains(element, " or ") {
				var x []string
				if strings.Contains(element, " and ") {
					x = strings.Split(element, " and ")
				} else {
					x = strings.Split(element, " or ")
				}
				for _, each := range x {
					sa, resp, err = GetMembers(allowed, strings.Split(each, " "), resp)
					if err != nil {
						return resp, err
					}
					s, err := json.Marshal(sa)
					if err != nil {
						errorMessage := " error while marshalling database data"
						return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil), fmt.Errorf(errorMessage)
					}
					strPara = strings.Replace(strPara, each, string(s), -1)
				}
			} else {
				sa, resp, err = GetMembers(allowed, strings.Split(element, " "), resp)
				if err != nil {
					return resp, err
				}
				s, err := json.Marshal(sa)
				if err != nil {
					errorMessage := " error while marshalling database data"
					return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil), fmt.Errorf(errorMessage)
				}
				strPara = strings.Replace(strPara, element, string(s), -1)
			}
		}
		expressionVals := re.FindAllString(strPara, -1)
		for _, eles := range expressionVals {
			var ww []dmtf.Link
			ele := strings.Trim(eles, "(")
			ele = strings.Trim(ele, ")")
			if strings.Contains(ele, " and ") {
				var inter [][]dmtf.Link
				for _, e := range strings.Split(ele, " and ") {
					_ = json.Unmarshal([]byte(e), &ww)
					inter = append(inter, ww)
				}
				interResp := LogicalOperation(inter, "and")
				interSam, _ := json.Marshal(interResp)
				strPara = strings.Replace(strPara, eles, string(interSam), 1)
			} else if strings.Contains(ele, " or ") {
				var inter [][]dmtf.Link
				for _, e := range strings.Split(ele, " or ") {
					_ = json.Unmarshal([]byte(e), &ww)
					inter = append(inter, ww)
				}
				interResp := LogicalOperation(inter, "or")
				interSam, _ := json.Marshal(interResp)
				strPara = strings.Replace(strPara, eles, string(interSam), 1)
			} else {
				strPara = strings.Replace(strPara, eles, ele, 1)
			}
		}

		if strings.Contains(strPara, "and ") || strings.Contains(strPara, " and") || strings.Contains(strPara, "or ") || strings.Contains(strPara, " or") {
			var ww []dmtf.Link
			var inter [][]dmtf.Link
			resp, err = validate(strPara, resp)
			if err != nil {
				log.Error(err.Error())
				return resp, err
			}
			if strings.Contains(strPara, " and ") {
				for _, e := range strings.Split(strPara, " and ") {
					resp, err = errorResp(e, resp)
					if err != nil {
						log.Error(err.Error())
						return resp, err
					}
					_ = json.Unmarshal([]byte(e), &ww)
					inter = append(inter, ww)
				}
				respMembers = LogicalOperation(inter, "and")
			}
			if strings.Contains(strPara, " or ") {
				for _, e := range strings.Split(strPara, " or ") {
					resp, err = errorResp(e, resp)
					if err != nil {
						log.Error(err.Error())
						return resp, err
					}
					_ = json.Unmarshal([]byte(e), &ww)
					inter = append(inter, ww)
				}
				respMembers = LogicalOperation(inter, "or")

			}
		}
		if !strings.Contains(strPara, " and ") && !strings.Contains(strPara, " or ") && !strings.Contains(strPara, " not ") {
			var mems []dmtf.Link
			var mem dmtf.Link
			para := strings.Trim(strPara, "[")
			para = strings.Trim(para, "]")
			for _, e := range strings.Split(para, ",") {
				_ = json.Unmarshal([]byte(e), &mem)
				mems = append(mems, mem)
			}
			respMembers = mems
		}

	} else if strings.Contains(strPara, "and ") || strings.Contains(strPara, " and") || strings.Contains(strPara, "or ") || strings.Contains(strPara, " or") {
		var ww []dmtf.Link
		var inter [][]dmtf.Link
		resp, err = validate(strPara, resp)
		if err != nil {
			log.Error(err.Error())
			return resp, err
		}
		if strings.Contains(strPara, " and ") {
			for _, each := range strings.Split(strPara, " and ") {
				resp, err = errorResp(each, resp)
				if err != nil {
					log.Error(err.Error())
					return resp, err
				}
				ww, resp, err = GetMembers(allowed, strings.Split(each, " "), resp)
				if err != nil {
					return resp, err
				}
				inter = append(inter, ww)
			}
			respMembers = LogicalOperation(inter, "and")
		}
		if strings.Contains(strPara, " or ") {
			for _, each := range strings.Split(strPara, " or ") {
				resp, err = errorResp(each, resp)
				if err != nil {
					log.Error(err.Error())
					return resp, err
				}
				ww, resp, err = GetMembers(allowed, strings.Split(each, " "), resp)
				if err != nil {
					return resp, err
				}
				inter = append(inter, ww)
			}
			respMembers = LogicalOperation(inter, "or")

		}
	} else if strings.Contains(strPara, " not ") || strings.Contains(strPara, "not ") {
		var all, exp []dmtf.Link
		all, resp, err = getAllSystemIDs(resp)
		if err != nil {
			return resp, err
		}
		expression := strings.Split(strPara, "not ")
		exp, resp, err = GetMembers(allowed, strings.Split(expression[1], " "), resp)
		if err != nil {
			return resp, err
		}
		for _, eachA := range exp {
			for j, eachB := range all {
				if eachA == eachB {
					all = append(all[:j], all[j+1:]...)
				}
			}
		}
		respMembers = all
	} else {
		respMembers, resp, err = GetMembers(allowed, strings.Split(strPara, " "), resp)
		if err != nil {
			return resp, err
		}
	}
	systemCollection := sresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
		OdataID:      "/redfish/v1/Systems",
		OdataType:    "#ComputerSystemCollection.ComputerSystemCollection",
		Description:  "Computer Systems view",
		Name:         "Computer Systems",
	}
	if len(respMembers) == 0 {
		respMembers = []dmtf.Link{}
	}
	systemCollection.Members = respMembers
	systemCollection.MembersCount = len(respMembers)
	resp.Body = systemCollection
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp, nil
}

// GetSystemResource is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting system resource information,  parameters need to be passed GetSystemsRequest .
// GetSystemsRequest holds the  Uuid,Url and Resourceid ,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (p *PluginContact) GetSystemResource(req *systemsproto.GetSystemsRequest) response.RPC {
	log.Debug("Entering the GetSystemResource with URL : ", req.URL)
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	// Splitting the SystemID to get UUID
	requestData := strings.Split(req.RequestParam, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"ComputerSystem", req.RequestParam}, nil)
	}
	uuid := requestData[0]

	var respData string
	var saveRequired bool
	// Getting the reset flag details for the requested URL
	deviceLoadFlag := getDeviceLoadInfo(req.URL, req.RequestParam)
	// deviceLoadFlag is true means flag is set for requested URL or the SystemID URL, load from device
	// deviceLoadFlag is false indicates flag is not set, load from DB
	if deviceLoadFlag {
		log.Debug("SystemReset flag is found for the URL ", req.URL)
		var getDeviceInfoRequest = scommon.ResourceInfoRequest{
			URL:             req.URL,
			UUID:            uuid,
			SystemID:        requestData[1],
			ContactClient:   p.ContactClient,
			DevicePassword:  p.DevicePassword,
			GetPluginStatus: p.GetPluginStatus,
		}
		log.Debug("Getting resource data from device for URL ", req.URL)
		var err error
		if respData, err = scommon.GetResourceInfoFromDevice(getDeviceInfoRequest, saveRequired); err != nil {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ComputerSystem", req.URL}, nil)
		}
	} else {
		saveRequired = true
		urlData := strings.Split(req.URL, "/")
		//generating search URL which will be a part of key and also used in formatting response
		var tableName string
		if req.ResourceID == "" {
			resourceName := urlData[len(urlData)-1]
			tableName = common.SystemResource[resourceName]
		} else {
			tableName = urlData[len(urlData)-2]
		}

		log.Debug("Getting the details from DB for URL ", req.URL)
		data, err := smodel.GetResource(tableName, req.URL)
		if err != nil {
			log.Error("getting system details from DB: " + err.Error())
			errorMessage := err.Error()
			if errors.DBKeyNotFound == err.ErrNo() {
				var getDeviceInfoRequest = scommon.ResourceInfoRequest{
					URL:             req.URL,
					UUID:            uuid,
					SystemID:        requestData[1],
					ContactClient:   p.ContactClient,
					DevicePassword:  p.DevicePassword,
					GetPluginStatus: p.GetPluginStatus,
				}
				var err error
				log.Debug("Getting the details from device for URL ", req.URL)
				if data, err = scommon.GetResourceInfoFromDevice(getDeviceInfoRequest, saveRequired); err != nil {
					return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ComputerSystem", req.URL}, nil)
				}
				if saveRequired && strings.Contains(req.URL, "/Storage") {
					rediscoverStorageInventory(uuid, "/redfish/v1/Systems/"+requestData[1]+"/Storage")
				}
			} else {
				return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
			}
		}
		respData = data
	}
	var resource map[string]interface{}
	json.Unmarshal([]byte(respData), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	log.Debug("Exiting the GetSystemResource with response ", resp)
	return resp
}

// getDeviceLoadInfo accepts URL and System ID as parameters and returns int
// returns true if exact URL entry is found in the DeviceLoad
// returns true if System ID entry is found in the DeviceLoad
// returns false if no entry is found in the DeviceLoad for the requested URL and System ID
func getDeviceLoadInfo(URL, systemID string) bool {
	systemURL := "/redfish/v1/Systems/" + systemID
	if _, err := smodel.GetSystemResetInfo(URL); err == nil {
		return true
	} else if _, err := smodel.GetSystemResetInfo(systemURL); err == nil {
		return true
	} else {
		return false
	}
}

// rediscoverSystemInventory will be triggered when ever the a valid storage URI or underneath URI's
// are requested which does not exist in DB. It will create a rpc for aggregation which will delete all storage inventory //
// and rediscover all of them
func rediscoverStorageInventory(systemID, systemURL string) {
	systemURL = strings.TrimSuffix(systemURL, "/")

	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		log.Error("failed to get client connection object for aggregator service")
		return
	}
	defer conn.Close()
	aggregator := aggregatorproto.NewAggregatorClient(conn)

	_, err = aggregator.RediscoverSystemInventory(context.TODO(), &aggregatorproto.RediscoverSystemInventoryRequest{
		SystemID:  systemID,
		SystemURL: systemURL,
	})
	if err != nil {
		log.Error("Error while rediscoverStorageInventroy")
		return
	}
	log.Info("rediscovery of system storage started.")
	return
}

// GetSystemsCollection is to fetch all the Systems uri's and retruns with created collection
// of systems data from odimra
func GetSystemsCollection(req *systemsproto.GetSystemsRequest) response.RPC {
	allowed := make(map[string]map[string]bool)
	allowed["searchKeys"] = make(map[string]bool)
	allowed["conditionKeys"] = make(map[string]bool)
	allowed["queryKeys"] = make(map[string]bool)
	for _, value := range scommon.SF.SearchKeys {
		for k := range value {
			allowed["searchKeys"][k] = true
		}
	}
	for _, value := range scommon.SF.ConditionKeys {
		allowed["conditionKeys"][value] = true
	}
	for _, value := range scommon.SF.QueryKeys {
		allowed["queryKeys"][value] = true
	}
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	paramStr := strings.SplitN(req.URL, "?", 2)
	if len(paramStr) > 1 {
		resp, retError := SearchAndFilter(paramStr, resp)
		if retError != nil {
			return resp
		}
		return resp
	}
	systemKeys, err := smodel.GetAllKeysFromTable("ComputerSystem")
	if err != nil {
		log.Error("error getting all keys of systemcollection table : " + err.Error())
		errorMessage := err.Error()
		if errorMessage == "error while trying to get resource details: no data with the with table name SystemCollection found" {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"ComputerSystem", ""}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	systemCollection := sresponse.Collection{
		OdataContext: "/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
		OdataID:      "/redfish/v1/Systems",
		OdataType:    "#ComputerSystemCollection.ComputerSystemCollection",
		Description:  "Computer Systems view",
		Name:         "Computer Systems",
	}
	members := []dmtf.Link{}
	for _, key := range systemKeys {
		members = append(members, dmtf.Link{Oid: key})
	}
	if len(members) == 0 {
		members = []dmtf.Link{}
	}
	systemCollection.Members = members
	systemCollection.MembersCount = len(members)
	resp.Body = systemCollection
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// GetSystems is used to fetch resource data. The function is supposed to be used as part of RPC
// For getting system resource information,  parameters need to be passed GetSystemsRequest .
// GetSystemsRequest holds the  Uuid,Url,
// Url will be parsed from that search key will created
// There will be two return values for the fuction. One is the RPC response, which contains the
// status code, status message, headers and body and the second value is error.
func (p *PluginContact) GetSystems(req *systemsproto.GetSystemsRequest) response.RPC {
	var resp response.RPC
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}

	requestData := strings.Split(req.RequestParam, ":")
	if len(requestData) <= 1 {
		errorMessage := "error: SystemUUID not found"
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"ComputerSystem", req.RequestParam}, nil)
	}
	uuid := requestData[0]
	var data string
	var err *errors.Error
	// check the whether SystemResetInfo available in db. If it is available, then get the data from device
	_, err = smodel.GetSystemResetInfo(req.URL)
	if err == nil {
		var getDeviceInfoRequest = scommon.ResourceInfoRequest{
			URL:             req.URL,
			UUID:            uuid,
			SystemID:        requestData[1],
			ContactClient:   p.ContactClient,
			DevicePassword:  p.DevicePassword,
			GetPluginStatus: p.GetPluginStatus,
			ResourceName:    "ComputerSystem",
		}
		var err error
		if data, err = scommon.GetResourceInfoFromDevice(getDeviceInfoRequest, true); err != nil {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"ComputerSystem", req.URL}, nil)
		}
	} else {
		data, err = smodel.GetSystemByUUID(req.URL)
		if err != nil {
			log.Error("error getting system details : " + err.Error())
			errorMessage := err.Error()
			if errors.DBKeyNotFound == err.ErrNo() {
				return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"ComputerSystem", req.RequestParam}, nil)
			}
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
		}
	}
	data = strings.Replace(data, `"Id":"`, `"Id":"`+uuid+`:`, -1)
	var resource map[string]interface{}
	json.Unmarshal([]byte(data), &resource)
	resp.Body = resource
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	return resp
}

// getStringData supports the eq and ne only for  expression
func getStringData(key, match, expr string, regexFlag bool) ([]string, error) {
	if expr == "eq" {
		return smodel.GetString(key, match, regexFlag)
	}
	// get all data
	allKeys, err := smodel.GetString(key, "", regexFlag)
	if err != nil {
		return []string{}, err
	}
	// get matching data
	matchedKeys, err := smodel.GetString(key, match, regexFlag)
	if err != nil {
		return []string{}, err
	}
	if len(allKeys) == len(matchedKeys) {
		return []string{}, nil
	}
	if len(matchedKeys) == 0 {
		return allKeys, nil
	}
	for i := 0; i < len(matchedKeys); i++ {
		for j := 0; j < len(allKeys); j++ {
			if matchedKeys[i] == allKeys[j] {
				allKeys = append(allKeys[:j], allKeys[j+1:]...)
				break
			}
		}
	}
	return allKeys, nil
}

func getRangeData(key, expr string, match int, regexFlag bool) ([]string, error) {
	if (expr == "le" || expr == "lt") && match == 0 {
		return []string{}, nil
	}
	switch expr {
	case "eq":
		return smodel.GetRange(key, match, match, regexFlag)
	case "gt":
		return smodel.GetRange(key, match+1, 100000000, regexFlag)
	case "ge":
		return smodel.GetRange(key, match, 100000000, regexFlag)
	case "lt":
		return smodel.GetRange(key, 0, match-1, regexFlag)
	case "le":
		return smodel.GetRange(key, 0, match, regexFlag)
	case "ne":
		if match == 0 {
			return smodel.GetRange(key, match+1, 100000000, regexFlag)

		}
		lowerboundKeys, err := smodel.GetRange(key, 0, match-1, regexFlag)
		if err != nil {
			return []string{}, err
		}
		upperboundKeys, err := smodel.GetRange(key, match+1, 100000000, regexFlag)
		if err != nil {
			return []string{}, err
		}
		var allKeys = make([]string, 0)
		allKeys = append(lowerboundKeys, upperboundKeys...)
		return allKeys, nil
	}
	return []string{}, nil
}

func parseRegexData(data []string, regex string) ([]string, error) {
	regex = strings.Replace(regex, "(", "\\(", -1)
	regex = strings.Replace(regex, ")", "\\)", -1)
	regex = "(?i)" + regex
	var list = make([]string, 0)
	for i := 0; i < len(data); i++ {
		values := strings.Split(data[i], "::")
		found, err := regexp.MatchString(regex, values[0])
		if err != nil {
			log.Error("regular expression error: " + err.Error())
			return list, err
		}
		if found {
			list = append(list, values[1])
		}

	}
	return list, nil
}

// this function checks query has open bracket["("] as prefix to ignore the brackets inside the string
// for e.g if query is ProcessorSummary/Model eq Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz
// here Inter(R) has bracket inbetween the string, so ignore this string for the first if search criteria
func checkParentheses(strPara string) bool {
	for _, val := range strings.Split(strPara, " ") {
		if strings.HasPrefix(val, "(") {
			return true
		}
	}
	return false
}
