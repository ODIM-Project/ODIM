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

//Package rfphandler ...
package rfphandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	evtConfig "github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfpmodel"
	"github.com/ODIM-Project/ODIM/plugin-redfish/rfputilities"
	iris "github.com/kataras/iris/v12"
)

//CreateEventSubscription : Subscribes for events
func CreateEventSubscription(ctx iris.Context) {

	device, deviceDetails, err := getDeviceDetails(ctx)
	if err != nil {
		return
	}
	//First delete existing matching subscription(our subscription) from device
	deleteMatchingSubscriptions(device)

	var reqPostBody rfpmodel.EvtSubPost
	var reqData string

	//replacing the reruest  with south bound translation URL
	for key, value := range evtConfig.Data.URLTranslation.SouthBoundURL {
		reqData = strings.Replace(string(deviceDetails.PostBody), key, value, -1)
	}

	err = json.Unmarshal([]byte(reqData), &reqPostBody)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	// remove the mesaageids, resourcestypes and originresources from the request and post it to device
	// since some of device doesnt support these
	req := rfpmodel.EvtSubPost{
		Destination: "https://" + evtConfig.Data.LoadBalancerConf.Host + ":" + evtConfig.Data.LoadBalancerConf.Port + evtConfig.Data.EventConf.DestURI,
		EventTypes:  reqPostBody.EventTypes,
		Context:     reqPostBody.Context,
		HTTPHeaders: reqPostBody.HTTPHeaders,
		Protocol:    reqPostBody.Protocol,
	}
	device.PostBody, err = json.Marshal(req)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "error: internal processing error: " + err.Error()
		log.Println(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	var resp *http.Response
	//Subscribe to Events
	resp, err = redfishClient.SubscribeForEvents(device)
	if err != nil {
		log.Println(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	defer resp.Body.Close()
	if err := validateResponse(ctx, device, resp, http.MethodPost); err != nil {
		return
	}
}

// Delete match subscription from device
func deleteMatchingSubscriptions(device *rfputilities.RedfishDevice) {
	// get all subscriptions
	device.Location = "https://" + device.Host + "/redfish/v1/EventService/Subscriptions"
	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		log.Println("error: internal processing error:", err)
		return
	}

	//Get Subscription details to check if it is really ours
	resp, err := redfishClient.GetSubscriptionDetail(device)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("error: while getting subscription details for URI: %v got %v", device.Location, resp.StatusCode)
		log.Println(errorMessage)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var subscriptionCollectionBody interface{}
	err = json.Unmarshal(body, &subscriptionCollectionBody)
	if err != nil {
		log.Println(err.Error())
		return
	}
	members := subscriptionCollectionBody.(map[string]interface{})["Members"]
	for _, member := range members.([]interface{}) {
		device.Location = member.(map[string]interface{})["@odata.id"].(string)
		device.Location = "https://" + device.Host + device.Location
		if isOurSubscription(device) {
			resp, err = redfishClient.DeleteSubscriptionDetail(device)
			if err != nil {
				log.Println(err.Error())
				return
			}
			resp.Body.Close()
		}
	}
	return
}
func isOurSubscription(device *rfputilities.RedfishDevice) bool {

	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		log.Println("error: internal processing error:", err)
		return false
	}
	//Get Subscription details to check if it is really ours
	resp, err := redfishClient.GetSubscriptionDetail(device)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("error: while getting subscription details for URI: %v got %v", device.Location, resp.StatusCode)
		log.Println(errorMessage)
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	var subscriptionBody interface{}
	err = json.Unmarshal(body, &subscriptionBody)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	subscriptionDestinationFromDevice := subscriptionBody.(map[string]interface{})["Destination"].(string)
	// if the subscription is ours then the destination should match with LBHOST:LBPORT.
	//If it is not matching then retrun with MethodNotAllowed
	if !strings.Contains(subscriptionDestinationFromDevice, evtConfig.Data.LoadBalancerConf.Host+":"+evtConfig.Data.LoadBalancerConf.Port) {
		return false
	}
	return true
}

//DeleteEventSubscription : Delete subscription
func DeleteEventSubscription(ctx iris.Context) {
	device, _, err := getDeviceDetails(ctx)
	if err != nil {
		return
	}
	redfishClient, err := rfputilities.GetRedfishClient()
	if err != nil {
		errMsg := "error: internal processing error: " + err.Error()
		log.Println(errMsg)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(errMsg)
		return
	}

	//Delete Subscription details
	resp, err := redfishClient.DeleteSubscriptionDetail(device)
	if err != nil {
		log.Println(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	defer resp.Body.Close()
	if err := validateResponse(ctx, device, resp, http.MethodDelete); err != nil {
		return
	}
}

// getDeviceDetails will accepts iris context and it will extract device details from context
// then decrypt the password and return device details
func getDeviceDetails(ctx iris.Context) (*rfputilities.RedfishDevice, *rfpmodel.Device, error) {
	//Get token from Request
	token := ctx.GetHeader("X-Auth-Token")
	//Validating the token
	if token != "" {
		flag := TokenValidation(token)
		if !flag {
			log.Println("Invalid/Expired X-Auth-Token")
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.WriteString("Invalid/Expired X-Auth-Token")
			return nil, nil, fmt.Errorf("Invalid/Expired X-Auth-Token")
		}
	}

	var deviceDetails rfpmodel.Device

	//Get device details from request
	err := ctx.ReadJSON(&deviceDetails)
	if err != nil {
		log.Println("Error while trying to collect data from request: ", err)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Error: bad request.")
		return nil, nil, err
	}

	device := &rfputilities.RedfishDevice{
		Host:     deviceDetails.Host,
		Username: deviceDetails.Username,
		Password: string(deviceDetails.Password),
		Location: deviceDetails.Location,
	}
	/*
		plainText, err := descryptDevicePassword(deviceDetails.Password)
		if err != nil {
			log.Println("Error while trying decrypt data: ", err)
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.WriteString("Error while trying to decypt data")
			return nil, nil, err
		}
		device.Password = plainText
	*/
	device.Password = string(deviceDetails.Password)

	return device, &deviceDetails, nil
}

// validateResponse will accepts iris context to write status code and resopnse
// method is to return status created incase of create subscription
// otherwise return statusok
func validateResponse(ctx iris.Context, device *rfputilities.RedfishDevice, resp *http.Response, method string) error {
	var body []byte
	var err error
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return err
	}
	defer resp.Body.Close()
	if strings.EqualFold(method, http.MethodPost) {
		// if there was an error for message ids means device haven't support of MessageIds
		// So remove the MessageIds from the request and subscribe again.
		if resp.StatusCode != http.StatusCreated {
			removeMessageID(ctx, device)
			// Subscribe to Events
			redfishClient, err := rfputilities.GetRedfishClient()
			if err != nil {
				errMsg := "error: internal processing error: " + err.Error()
				log.Println(errMsg)
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(errMsg)
				return err
			}

			resp, err = redfishClient.SubscribeForEvents(device)
			if err != nil {
				log.Println(err.Error())
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return err
			}
			defer resp.Body.Close()
			body, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err.Error())
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString(err.Error())
				return err
			}
		}

	}
	header := map[string]string{
		"Location": resp.Header.Get("Location"),
	}

	if resp.StatusCode == 401 {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.WriteString("Authtication with the device failed")
		return errors.New("Authtication with the device failed")
	}
	if resp.StatusCode >= 300 {
		log.Printf("Subscription operation failed: \n%s\n\n", body)
	}
	common.SetResponseHeader(ctx, header)
	ctx.StatusCode(resp.StatusCode)
	log.Printf("Redfish plugin response body: %s \n", body)
	ctx.WriteString(string(body))
	return nil
}

func removeMessageID(ctx iris.Context, device *rfputilities.RedfishDevice) {
	var ReqPostBody rfpmodel.EvtSubPost
	err := json.Unmarshal(device.PostBody, &ReqPostBody)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	ReqPostBody.MessageIds = []string{}
	device.PostBody, err = json.Marshal(ReqPostBody)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	return
}
