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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
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
		errMsg := "unable to parse the create request"
		log.Println(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errMsg, nil, nil)
	}

	err = validateElements(createRequest.Elements)
	if err != nil {
		errMsg := "invalid elements for create an aggregate" + err.Error()
		log.Println(errMsg)
		errArgs := []interface{}{"Elements", createRequest}
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errMsg, errArgs, nil)
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
		OdataContext: "/redfish/v1/$metadata#AggregationService.Aggregates",
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
func validateElements(elements []string) error {
	for _, element := range elements {
		if _, err := agmodel.GetSystem(element); err != nil {
			return err
		}
	}
	return nil
}
