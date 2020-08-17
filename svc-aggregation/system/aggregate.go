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

// Package system ...
package system

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agresponse"
	uuid "github.com/satori/go.uuid"
)

// CreateAggregate is the handler for creating an aggregate
// check if the elelments/resources added into odimra if not then return an error.
// else add an entry of an aggregayte in db
func (e *ExternalInterface) CreateAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	// parsing the aggregate request
	var createRequest agmodel.Aggregate
	err := json.Unmarshal(req.RequestBody, &createRequest)
	if err != nil {
		errMsg := "unable to parse the create request" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}
	//empty request check
	if reflect.DeepEqual(agmodel.Aggregate{}, createRequest) {
		errMsg := "empty request can not be processed"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.PropertyMissing, errMsg, []interface{}{"Elements"}, nil)
	}

	statuscode, err := validateElements(createRequest.Elements)
	if err != nil {
		errMsg := "invalid elements for create an aggregate" + err.Error()
		log.Println(errMsg)
		errArgs := []interface{}{"Elements", createRequest}
		return common.GeneralError(statuscode, response.ResourceNotFound, errMsg, errArgs, nil)
	}
	targetURI := "/redfish/v1/AggregationService/Aggregates"
	aggregateUUID := uuid.NewV4().String()
	var aggregateURI = fmt.Sprintf("%s/%s", targetURI, aggregateUUID)

	dbErr := agmodel.CreateAggregate(createRequest, aggregateURI)
	if dbErr != nil {
		errMsg := dbErr.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	commonResponse := response.Response{
		OdataType:    "#Aggregate.v1_0_0.Aggregate",
		OdataID:      aggregateURI,
		OdataContext: "/redfish/v1/$metadata#Aggregate.Aggregate",
		ID:           aggregateUUID,
		Name:         "Aggregate",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Link":              "<" + aggregateURI + "/>; rel=describedby",
		"Location":          aggregateURI,
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Created)
	resp.Body = agresponse.AggregateResponse{
		Response: commonResponse,
		Elements: createRequest.Elements,
	}
	resp.StatusCode = http.StatusCreated
	return resp
}

// check if the resource is exist in odim
func validateElements(elements []string) (int32, error) {
	if checkDuplicateElements(elements) {
		return http.StatusBadRequest, errors.PackError(errors.UndefinedErrorType, fmt.Errorf("Duplicate elements present"))
	}
	for _, element := range elements {
		if _, err := agmodel.GetSystem(element); err != nil {
			return http.StatusNotFound, err
		}
	}
	return http.StatusOK, nil
}

//check if the elements have duplicate element
func checkDuplicateElements(elelments []string) bool {
	duplicate := make(map[string]int)
	for _, element := range elelments {
		// check if the item/element exist in the duplicate map
		_, exist := duplicate[element]
		if exist {
			return true
		}
		duplicate[element] = 1

	}
	return false
}

// GetAllAggregates is the handler for getting collection of aggregates
func (e *ExternalInterface) GetAllAggregates(req *aggregatorproto.AggregatorRequest) response.RPC {
	aggregateKeys, err := agmodel.GetAllKeysFromTable("Aggregate")
	if err != nil {
		log.Printf("error getting aggregate : %v", err.Error())
		errorMessage := err.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errorMessage, []interface{}{config.Data.DBConf.OnDiskHost + ":" + config.Data.DBConf.OnDiskPort}, nil)
	}
	var members = make([]agresponse.ListMember, 0)
	for i := 0; i < len(aggregateKeys); i++ {
		members = append(members, agresponse.ListMember{
			OdataID: aggregateKeys[i],
		})
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	commonResponse := response.Response{
		OdataType:    "#AggregateCollection.v1_0_0.AggregateCollection",
		OdataID:      "/redfish/v1/AggregationService/Aggregates",
		OdataContext: "/redfish/v1/$metadata#AggregateCollection.AggregateCollection",
		ID:           "Aggregate",
		Name:         "Aggregate",
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
	resp.Body = agresponse.List{
		Response:     commonResponse,
		MembersCount: len(members),
		Members:      members,
	}
	return resp
}

// GetAggregate is the handler for getting an aggregate
//if the aggregate id is present then return aggregate details else return an error.
func (e *ExternalInterface) GetAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	aggregate, err := agmodel.GetAggregate(req.URL)
	if err != nil {
		log.Printf("error getting  Aggregate : %v", err)
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	var data = strings.Split(req.URL, "/redfish/v1/AggregationService/Aggregates/")
	commonResponse := response.Response{
		OdataType:    "#Aggregate.v1_0_0.Aggregate",
		OdataID:      req.URL,
		OdataContext: "/redfish/v1/$metadata#Aggregate.Aggregate",
		ID:           data[1],
		Name:         "Aggregate",
	}
	var resp = response.RPC{
		StatusCode:    http.StatusOK,
		StatusMessage: response.Success,
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	commonResponse.CreateGenericResponse(response.Success)
	resp.Body = agresponse.AggregateResponse{
		Response: commonResponse,
		Elements: aggregate.Elements,
	}
	return resp
}

// DeleteAggregate is the handler for deleting an aggregate
// if the aggregate id is present then delete from the db else return an error.
func (e *ExternalInterface) DeleteAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	var resp response.RPC
	_, err := agmodel.GetAggregate(req.URL)
	if err != nil {
		log.Printf("error getting  Aggregate : %v", err)
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	err = agmodel.DeleteAggregate(req.URL)
	if err != nil {
		log.Printf("error while deleting an aggregate : %v", err)
		errorMessage := err.Error()
		if errors.DBKeyNotFound == err.ErrNo() {
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), []interface{}{"Aggregate", req.URL}, nil)
		}
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errorMessage, nil, nil)
	}
	resp.StatusCode = http.StatusNoContent
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	return resp
}

// AddElementsToAggregate is the handler for adding elements to an aggregate
func (e *ExternalInterface) AddElementsToAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	// TODO add functionality to create an aggregate
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}

// RemoveElementsFromAggregate is the handler for removing elements from an aggregate
func (e *ExternalInterface) RemoveElementsFromAggregate(req *aggregatorproto.AggregatorRequest) response.RPC {
	// TODO add functionality to create an aggregate
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}
