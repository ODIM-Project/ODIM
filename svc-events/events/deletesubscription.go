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

// Package events have the functionality of
// - Create Event Subscription
// - Delete Event Subscription
// - Get Event Subscription
// - Post Event Subscription to destination
// - Post TestEvent (SubmitTestEvent)
// and corresponding unit test cases
package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

// DeleteEventSubscriptions delete subscription data against given URL
func (p *PluginContact) DeleteEventSubscriptions(req *eventsproto.EventRequest) response.RPC {
	var resp response.RPC
	originResource := req.UUID
	uuid, err := getUUID(originResource)
	if err != nil {
		errorMessage := err.Error()
		msgArgs := []interface{}{"OriginResource", originResource}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		log.Error(err.Error())
		return resp
	}
	target, err := evmodel.GetTarget(uuid)
	if err != nil {
		log.Error("error while getting device details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"uuid", uuid}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		return resp
	}
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		log.Error(errorMessage)
		return resp
	}
	searchKey := evcommon.GetSearchKey(deviceIPAddress, evmodel.SubscriptionIndex)
	log.Info("Getting event subscription details of device: ", deviceIPAddress)
	subscriptionDetails, err := evmodel.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		log.Error("error while getting event subscription details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		return resp
	}
	if len(subscriptionDetails) < 1 {
		errorMessage := fmt.Sprintf("Subscription details not found for the requested device")
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		log.Error(errorMessage)
		return resp
	}
	log.Info("Number of subscription present :", strconv.Itoa(len(subscriptionDetails)))
	decryptedPasswordByte, err := common.DecryptWithPrivateKey(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		msgArgs := []interface{}{""}
		evcommon.GenErrorResponse(errorMessage, response.InternalError, http.StatusInternalServerError, msgArgs, &resp)
		log.Error(errorMessage)
		return resp
	}
	target.Password = decryptedPasswordByte

	// Delete Event Subscription from device also
	err = p.deleteSubscription(target, originResource)
	if err != nil {
		log.Error("error while deleting eventsubscription details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		return resp
	}

	searchKey = evcommon.GetSearchKey(deviceIPAddress, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err := evmodel.GetDeviceSubscriptions(searchKey)
	if err != nil {
		errorMessage := "Error while get subscription details of device : " + err.Error()
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		log.Error(errorMessage)
		return resp
	}
	originResource = deviceSubscription.OriginResources[0]
	log.Info("Device subcription information", deviceSubscription.EventHostIP)

	for _, evtSubscription := range subscriptionDetails {

		// Delete Event Subscription details from the Subscription(table) in DB

		// if there is only one host in Hosts entry then
		// delete the subscription from redis
		if len(evtSubscription.Hosts) == 1 {
			err = evmodel.DeleteEvtSubscription(evtSubscription.SubscriptionID)
			if err != nil {
				errorMessage := "Error while Updating event subscription : " + err.Error()
				msgArgs := []interface{}{"SubscriptionID", evtSubscription.SubscriptionID}
				evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
				log.Error(errorMessage)
				return resp
			}
		} else {
			// Delete the host and origin resource from the respective entry
			evtSubscription.Hosts = removeElement(evtSubscription.Hosts, target.ManagerAddress)
			evtSubscription.OriginResources = removeElement(evtSubscription.OriginResources, originResource)
			err = evmodel.UpdateEventSubscription(evtSubscription)
			if err != nil {
				errorMessage := "Error while Updating event subscription : " + err.Error()
				msgArgs := []interface{}{"SubscriptionID", evtSubscription.SubscriptionID}
				evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
				log.Error(errorMessage)
				return resp
			}
		}

	}
	err = evmodel.DeleteDeviceSubscription(searchKey)
	if err != nil {
		errorMessage := "Error while deleting device subscription : " + err.Error()
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		log.Error(errorMessage)
		return resp
	}

	resp.StatusCode = http.StatusNoContent
	resp.StatusMessage = response.ResourceRemoved
	return resp
}

// deleteSubscription to the Event Subsciption
func (p *PluginContact) deleteSubscription(target *evmodel.Target, originResource string) error {

	var plugin *evmodel.Plugin
	plugin, err := evmodel.GetPluginData(target.PluginID)
	if err != nil {
		return err
	}

	if _, errs := p.DeleteSubscriptions(originResource, "", plugin, target); errs != nil {
		return errs
	}
	return nil
}

// DeleteEventSubscriptionsDetails delete subscription data against given subscription id
func (p *PluginContact) DeleteEventSubscriptionsDetails(req *eventsproto.EventRequest) response.RPC {
	var resp response.RPC
	authResp := p.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authenticate session: status code: " + string(authResp.StatusCode) + ", status message: " + authResp.StatusMessage)
		return authResp
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
		"allow":             "POST,GET,DELETE",
	}

	subscriptionDetails, err := evmodel.GetEvtSubscriptions(req.EventSubscriptionID)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		log.Error("error while deleting eventsubscription details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"SubscriptionID", req.EventSubscriptionID}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		return resp
	}

	if len(subscriptionDetails) < 1 {
		errorMessage := fmt.Sprintf("Subscription details not found for subscription id: %s", req.EventSubscriptionID)
		log.Error(errorMessage)
		var msgArgs = []interface{}{"SubscrfiptionID", req.EventSubscriptionID}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		return resp
	}
	for _, evtSubscription := range subscriptionDetails {

		// Since we are searching subscription id with pattern search
		// we need to match the subscripton id
		if evtSubscription.SubscriptionID != req.EventSubscriptionID {
			errorMessage := fmt.Sprintf("Subscription details not found for subscription id: %s", req.EventSubscriptionID)
			log.Error(errorMessage)
			var msgArgs = []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
			return resp
		}

		// Delete and re subscrive Event Subscription
		err = p.deleteAndReSubscribetoEvents(evtSubscription)
		if err != nil {
			log.Error("error while deleting eventsubscription details : " + err.Error())
			errorMessage := err.Error()
			msgArgs := []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
			return resp
		}

		// Delete Event Subscription from the DB
		err = evmodel.DeleteEvtSubscription(evtSubscription.SubscriptionID)
		if err != nil {
			log.Error("error while deleting eventsubscription details : " + err.Error())
			errorMessage := err.Error()
			msgArgs := []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
			return resp
		}

	}

	commonResponse := response.Response{
		OdataType: "#EventDestination.v1_7_0.EventDestination",
		OdataID:   "/redfish/v1/EventService/Subscriptions/" + req.EventSubscriptionID,
		ID:        req.EventSubscriptionID,
		Name:      "Event Subscription",
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.ResourceRemoved

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	resp.Body = commonResponse
	return resp
}

// This function is to delete and re subscribe for Event Subscriptions
func (p *PluginContact) deleteAndReSubscribetoEvents(evtSubscription evmodel.Subscription) error {
	originResources := evtSubscription.OriginResources
	for _, origin := range originResources {
		// ignore if origin is empty
		if origin == "" {
			continue
		}
		subscriptionDetails, err := evmodel.GetEvtSubscriptions(origin)
		if err != nil {
			return err
		}

		// if origin contains fabrics then get all the collection and individual subscription details
		// for Systems need to add same later
		subscriptionDetails = getAllSubscriptions(origin, subscriptionDetails)
		// if deleteflag is true then only one document is there
		// so dont re subscribe again
		var deleteflag bool
		if len(subscriptionDetails) < 1 {
			return fmt.Errorf("Subscription details not found for subscription id: %s", origin)
		} else if len(subscriptionDetails) == 1 {
			deleteflag = true
		}

		var context, protocol, destination, name string
		var eventTypes, messageIDs, resourceTypes []string

		for index, evtSub := range subscriptionDetails {
			if evtSubscription.SubscriptionID != evtSub.SubscriptionID {
				if len(evtSub.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
					eventTypes = append(eventTypes, evtSub.EventTypes...)
				} else {
					eventTypes = []string{}
				}

				if len(evtSub.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
					messageIDs = append(messageIDs, evtSub.MessageIds...)
				} else {
					messageIDs = []string{}
				}

				if len(evtSub.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
					resourceTypes = append(resourceTypes, evtSub.ResourceTypes...)
				} else {
					resourceTypes = []string{}
				}
				name = evtSub.Name
				context = evtSub.Context
				protocol = evtSub.Protocol
				destination = evtSub.Destination
			}
		}

		eventTypesCount := len(eventTypes)
		messageIDsCount := len(messageIDs)
		resourceTypesCount := len(resourceTypes)
		removeDuplicatesFromSlice(&eventTypes, &eventTypesCount)
		removeDuplicatesFromSlice(&messageIDs, &messageIDsCount)
		removeDuplicatesFromSlice(&resourceTypes, &resourceTypesCount)
		var httpHeadersSlice = make([]evmodel.HTTPHeaders, 0)
		httpHeadersSlice = append(httpHeadersSlice, evmodel.HTTPHeaders{ContentType: "application/json"})
		subscriptionPost := evmodel.EvtSubPost{
			Name:          name,
			EventTypes:    eventTypes,
			MessageIds:    messageIDs,
			ResourceTypes: resourceTypes,
			HTTPHeaders:   httpHeadersSlice,
			Context:       context,
			Protocol:      protocol,
			Destination:   destination,
		}

		err = p.subscribe(subscriptionPost, origin, deleteflag)
		if err != nil {
			return err
		}
	}
	return nil
}

func isCollectionOriginResourceURI(origin string) bool {

	if origin == "" || !strings.HasPrefix(origin, "/") {
		return false
	}

	origin = strings.TrimSuffix(origin, "/")

	defaultCollectionURIs := []string{
		"/redfish/v1/Systems",
		"/redfish/v1/Chassis",
		"/redfish/v1/Fabrics",
		"/redfish/v1/Managers",
		"/redfish/v1/TaskService/Tasks",
	}

	front := 0
	rear := len(defaultCollectionURIs) - 1
	for front <= rear {
		if defaultCollectionURIs[front] == origin || defaultCollectionURIs[rear] == origin {
			return true
		}
		front++
		rear--
	}
	return false
}

// Subscribe to the Event Subsciption
func (p *PluginContact) subscribe(subscriptionPost evmodel.EvtSubPost, origin string, deleteflag bool) error {

	if strings.Contains(origin, "Fabrics") {
		return p.resubscribeFabricsSubscription(subscriptionPost, origin, deleteflag)
	}
	originResource := origin
	if isCollectionOriginResourceURI(originResource) {
		log.Error("Collection of origin resource:" + originResource)
		return nil
	}
	target, _, err := getTargetDetails(originResource)
	if err != nil {
		return err
	}

	plugin, errs := evmodel.GetPluginData(target.PluginID)
	if errs != nil {
		return errs
	}

	postBody, err := json.Marshal(subscriptionPost)
	if err != nil {
		return fmt.Errorf("Error while marshalling subscription details: %s", err)
	}
	target.PostBody = postBody
	_, err = p.DeleteSubscriptions(origin, "", plugin, target)
	if err != nil {
		return err
	}
	// if deleteflag is true then only one document is there
	// so dont re subscribe again
	if deleteflag {
		return nil
	}

	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := p.getPluginToken(plugin)
		if token == "" {
			return fmt.Errorf("error: Unable to create session with plugin " + plugin.ID)
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = target

	_, loc, _, err := p.PluginCall(contactRequest)
	if err != nil {
		return err
	}
	// Update Location to all destination of device if already subscribed to the device
	var resp response.RPC
	deviceIPAddress, errorMessage := evcommon.GetIPFromHostName(target.ManagerAddress)
	if errorMessage != "" {
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		log.Error(errorMessage)
	}
	searchKey := evcommon.GetSearchKey(deviceIPAddress, evmodel.DeviceSubscriptionIndex)
	devSub, err := evmodel.GetDeviceSubscriptions(searchKey)
	if err != nil {
		return err
	}
	deviceSub := evmodel.DeviceSubscription{
		EventHostIP:     devSub.EventHostIP,
		Location:        loc,
		OriginResources: devSub.OriginResources,
	}
	return evmodel.UpdateDeviceSubscriptionLocation(deviceSub)

}

// DeleteFabricsSubscription will delete fabric subscription
func (p *PluginContact) DeleteFabricsSubscription(originResource string, plugin *evmodel.Plugin) (response.RPC, error) {
	var resp response.RPC
	addr, errorMessage := evcommon.GetIPFromHostName(plugin.IP)
	if errorMessage != "" {
		var msgArgs = []interface{}{"ManagerAddress", plugin.IP}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		log.Error(errorMessage)
		return resp, fmt.Errorf(errorMessage)
	}
	searchKey := evcommon.GetSearchKey(addr, evmodel.DeviceSubscriptionIndex)
	devSub, err := evmodel.GetDeviceSubscriptions(searchKey)
	if err != nil {
		errorMessage := "Error while get device subscription details: " + err.Error()
		if strings.Contains(err.Error(), "No data found for the key") {
			var msgArgs = []interface{}{"CFM Plugin", addr}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
			log.Error(errorMessage)
			return resp, err
		}
		evcommon.GenErrorResponse(errorMessage, response.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		log.Error(errorMessage)
		return resp, err
	}

	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := p.getPluginToken(plugin)
		if token == "" {
			evcommon.GenErrorResponse("error: Unable to create session with plugin "+plugin.ID, response.NoValidSession, http.StatusUnauthorized,
				[]interface{}{}, &resp)
			return resp, fmt.Errorf("error: Unable to create session with plugin " + plugin.ID)
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	// Call to delete subscription to plugin
	contactRequest.URL = devSub.Location
	contactRequest.HTTPMethodType = http.MethodDelete
	contactRequest.PostBody = nil
	resp, _, _, err = p.PluginCall(contactRequest)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		resp, _, _, err = p.retryEventOperation(contactRequest)
		if err != nil {
			return resp, err
		}
	}
	return resp, nil
}

//  resubscribeFabricsSubscription updates subscription fabric subscription details  by forming the super set of MessageIDs,EventTypes and ResourceTypes
func (p *PluginContact) resubscribeFabricsSubscription(subscriptionPost evmodel.EvtSubPost, origin string, deleteflag bool) error {
	originResources := getSuboridanteResourcesFromCollection(origin)
	for _, origin := range originResources {
		originResource := origin
		fabricID := getFabricID(originResource)
		if fabricID == "" {
			return nil
		}
		// get Fabrics Details
		fabric, dberr := evmodel.GetFabricData(fabricID)
		if dberr != nil {
			errorMessage := "error while getting fabric data: " + dberr.Error()
			log.Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
		plugin, errs := evmodel.GetPluginData(fabric.PluginID)
		if errs != nil {
			errorMessage := "error while getting plugin data: " + errs.Error()
			log.Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
		// Deleting the fabric subscription
		_, err := p.DeleteFabricsSubscription(origin, plugin)
		if err != nil {
			return err
		}

		// if deleteflag is true then only one document is there
		// so dont re subscribe again
		if deleteflag {
			return nil
		}

		var contactRequest evcommon.PluginContactRequest

		contactRequest.Plugin = plugin
		if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			token := p.getPluginToken(plugin)
			if token == "" {
				return fmt.Errorf("error: Unable to create session with plugin " + plugin.ID)
			}
			contactRequest.Token = token
		} else {
			contactRequest.LoginCredential = map[string]string{
				"UserName": plugin.Username,
				"Password": string(plugin.Password),
			}

		}
		// filling origin resource
		subscriptionPost.OriginResources = []evmodel.OdataIDLink{
			evmodel.OdataIDLink{
				OdataID: originResource,
			},
		}
		postBody, _ := json.Marshal(subscriptionPost)
		var reqData string
		//replacing the reruest url with south bound translation URL
		for key, value := range config.Data.URLTranslation.SouthBoundURL {
			reqData = strings.Replace(string(postBody), key, value, -1)
		}

		// recreating the subscription
		contactRequest.URL = "/ODIM/v1/Subscriptions"
		contactRequest.HTTPMethodType = http.MethodPost
		err = json.Unmarshal([]byte(reqData), &contactRequest.PostBody)
		log.Info("Resubscribe request" + reqData)
		response, loc, _, err := p.PluginCall(contactRequest)
		if err != nil {
			return err
		}
		if response.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			_, _, _, err = p.retryEventOperation(contactRequest)
			if err != nil {
				return err
			}
		}
		log.Info("Resubscribe response status code: " + string(response.StatusCode))
		log.Info("Resubscribe response body: ", response.Body)
		addr, errorMessage := evcommon.GetIPFromHostName(plugin.IP)
		if errorMessage != "" {
			return fmt.Errorf(errorMessage)
		}
		searchKey := evcommon.GetSearchKey(addr, evmodel.DeviceSubscriptionIndex)
		// Update Location to all destination of device if already subscribed to the device
		devSub, err := evmodel.GetDeviceSubscriptions(searchKey)
		if err != nil {
			return err
		}
		deviceSub := evmodel.DeviceSubscription{
			EventHostIP:     devSub.EventHostIP,
			Location:        loc,
			OriginResources: devSub.OriginResources,
		}
		err = evmodel.UpdateDeviceSubscriptionLocation(deviceSub)
		if err != nil {
			return err
		}
	}

	return nil
}

func getSuboridanteResourcesFromCollection(originResources string) []string {
	data, _, collectionPresentflag, _ := checkCollection(originResources)
	if !collectionPresentflag {
		return []string{originResources}
	}
	return data
}

func getAllSubscriptions(origin string, subscriptionDetails []evmodel.Subscription) []evmodel.Subscription {
	if origin == "/redfish/v1/Fabrics" {
		return subscriptionDetails
	}

	searchKey := "/redfish/v1/Fabrics"
	subscriptions, err := evmodel.GetEvtSubscriptions(searchKey)
	if err != nil {
		return subscriptionDetails
	}
	// Checking the collection based subscription
	var collectionSubscription = make([]evmodel.Subscription, 0)
	for _, evtSubscription := range subscriptions {
		for _, originResource := range evtSubscription.OriginResources {
			if strings.Contains(origin, "Fabrics") && originResource == "/redfish/v1/Fabrics" {
				collectionSubscription = append(collectionSubscription, evtSubscription)
			}
		}
	}

	if len(collectionSubscription) < 1 {
		return subscriptionDetails
	}
	collectionSubscription = append(collectionSubscription, subscriptionDetails...)
	removeDuplicatesFromSubscription(collectionSubscription)
	return collectionSubscription
}

// remove duplicate elements in string slice.
// Takes string slice and length, and updates the same with new values
func removeDuplicatesFromSubscription(subscriptions []evmodel.Subscription) []evmodel.Subscription {
	uniqueElementsDs := make(map[string]bool)
	var uniqueElemenstsList []evmodel.Subscription
	for _, sub := range subscriptions {
		if exist := uniqueElementsDs[sub.SubscriptionID]; !exist {
			uniqueElemenstsList = append(uniqueElemenstsList, sub)
			uniqueElementsDs[sub.SubscriptionID] = true
		}
	}
	return uniqueElemenstsList
}
