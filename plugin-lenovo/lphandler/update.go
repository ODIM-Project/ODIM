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

//Package lphandler ...
package lphandler

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpmodel"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lpresponse"
	"github.com/ODIM-Project/ODIM/plugin-lenovo/lputilities"
	iris "github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// SimpleUpdate updates the BMC resources
func SimpleUpdate(ctx iris.Context) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Error("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return
		}
	}
	var deviceDetails lpmodel.Device
	uri := ctx.Request().RequestURI
	uri = translateToSouthBoundURL(uri)
	log.Debug("Incoming uri in SimpleUpdate : ", uri)
	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		errMsg := "Unable to collect data from request: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(errMsg)
		return
	}
	var requestBody map[string]interface{}
	err = json.Unmarshal(deviceDetails.PostBody, &requestBody)
	if err != nil {
		log.Error("While trying to unmarshal request, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return
	}
	log.Debug("Incoming request body in SimpleUpdate: ", requestBody)
	delete(requestBody, "Targets")
	urlList := strings.Split(uri, ".")

	var operationApplyTime string
	if requestBody["@Redfish.OperationApplyTime"] != nil {
		operationApplyTime = requestBody["@Redfish.OperationApplyTime"].(string)
	}
	if operationApplyTime == "OnStartUpdateRequest" && urlList[1] == "SimpleUpdate" {
		body := `{"error":{"code": Base.1.4.Success,"message": "See @Message.ExtendedInfo for more information.","@Message.ExtendedInfo":[{"MessageId": "Base.1.4.Success"}]}}`
		ctx.StatusCode(http.StatusOK)
		ctx.WriteString(body)
		return
	} else if operationApplyTime != "" && urlList[1] == "SimpleUpdate" {
		errMsg := "The @Redfish.OperationApplyTime property value supplied is not supported by Lenovo, Lenovo supports only OnStartUpdateRequest"
		msgArgs := []string{operationApplyTime}
		body, _ := createUpdateActionResponse(response.PropertyValueNotInList, errMsg, msgArgs)
		log.Error(string(body))
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString(string(body))
		return
	}
	if urlList[1] == "StartUpdate" {
		uri = strings.Replace(uri, "StartUpdate", "SimpleUpdate", -1)
		delete(requestBody, "@Redfish.OperationApplyTime")
	}

	marshalBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Error("While trying to marshal request, got: " + err.Error())
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
	}
	device := &lputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		PostBody: marshalBody,
	}
	redfishClient, err := lputilities.GetRedfishClient()
	if err != nil {
		errMsg := "While trying to get redfish client, got: " + err.Error()
		log.Error(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}
	log.Debug("URI for device in SimpleUpdate:", uri)
	log.Debug("Payload for device in SimpleUpdate:", string(marshalBody))

	//Update BMC resource
	resp, err := redfishClient.DeviceCall(device, uri, http.MethodPost)
	if err != nil {
		errorMessage := "While trying to update BMC resource, got: " + err.Error()
		log.Error(errorMessage)
		if resp == nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString(errorMessage)
			return
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorMessage := "While trying to update BMC resource, got: " + err.Error()
		log.Error(errorMessage)
		ctx.WriteString(errorMessage)
	}
	log.Debug("Device response status : ", resp.StatusCode)
	log.Debug("Device response body in SimpleUpdate: ", string(body))
	ctx.StatusCode(resp.StatusCode)
	ctx.Write(body)
}

// createUpdateActionResponse is used for creating a final response for update action
func createUpdateActionResponse(messageID, message string, msgArgs []string) ([]byte, error) {
	resp := lpresponse.ErrorResopnse{
		Error: lpresponse.Error{
			Code:    response.ExtendedInfo,
			Message: "See @Message.ExtendedInfo for more information.",
			MessageExtendedInfo: []lpresponse.MsgExtendedInfo{
				lpresponse.MsgExtendedInfo{
					MessageID:   messageID,
					Message:     message,
					MessageArgs: msgArgs,
				},
			},
		},
	}
	body, err := json.Marshal(resp)
	return body, err
}
