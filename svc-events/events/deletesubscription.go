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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

var (
	//GetIPFromHostNameFunc ...
	GetIPFromHostNameFunc = common.GetIPFromHostName
	//DecryptWithPrivateKeyFunc ...
	DecryptWithPrivateKeyFunc = common.DecryptWithPrivateKey
)

// DeleteEventSubscriptions delete subscription data against given URL
func (e *ExternalInterfaces) DeleteEventSubscriptions(ctx context.Context, req *eventsproto.EventRequest, taskId string) response.RPC {
	var resp response.RPC
	originResource := req.UUID
	uuid, err := getUUID(originResource)
	var (
		percentComplete int32 = 100
		targetURI             = "/redfish/v1/EventService/Subscriptions"
	)
	if err != nil {
		msgArgs := []interface{}{"OriginResource", originResource}
		evcommon.GenErrorResponse(err.Error(), response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		l.LogWithFields(ctx).Error(err.Error())
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	target, err := e.GetTarget(uuid)
	if err != nil {
		l.LogWithFields(ctx).Error("error while getting device details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"uuid", uuid}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	deviceIPAddress, err := GetIPFromHostNameFunc(target.ManagerAddress)
	if err != nil {
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(err.Error(), response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		l.LogWithFields(ctx).Error(err.Error())
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	searchKey := evcommon.GetSearchKey(deviceIPAddress, evmodel.SubscriptionIndex)
	subscriptionDetails, err := e.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		l.LogWithFields(ctx).Error("error while getting event subscription details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	l.LogWithFields(ctx).Debug("Number of subscription present :", strconv.Itoa(len(subscriptionDetails)))
	decryptedPasswordByte, err := DecryptWithPrivateKeyFunc(target.Password)
	if err != nil {
		// Frame the RPC response body and response Header below
		errorMessage := "error while trying to decrypt device password: " + err.Error()
		msgArgs := []interface{}{""}
		evcommon.GenErrorResponse(errorMessage, response.InternalError, http.StatusInternalServerError, msgArgs, &resp)
		l.LogWithFields(ctx).Error(errorMessage)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	target.Password = decryptedPasswordByte

	// Delete Event Subscription from device also
	err = e.deleteSubscription(ctx, target, originResource)
	if err != nil {
		l.LogWithFields(ctx).Error("error while deleting event subscription details : " + err.Error())
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(err.Error(), response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	searchKey = evcommon.GetSearchKey(deviceIPAddress, evmodel.DeviceSubscriptionIndex)
	deviceSubscription, err := e.GetDeviceSubscriptions(searchKey)
	if err != nil {
		errorMessage := "Error while get subscription details of device : " + err.Error()
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		l.LogWithFields(ctx).Error(errorMessage)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return resp
	}
	originResource = deviceSubscription.OriginResources[0]
	l.LogWithFields(ctx).Info("Device subscription information ", deviceSubscription.EventHostIP)

	for _, evtSubscription := range subscriptionDetails {

		// Delete Event Subscription details from the Subscription(table) in DB

		// if there is only one host in Hosts entry then
		// delete the subscription from redis
		if len(evtSubscription.Hosts) == 1 {
			err = e.DeleteEvtSubscription(evtSubscription.SubscriptionID)
			if err != nil {
				errorMessage := "Error while Updating event subscription : " + err.Error()
				msgArgs := []interface{}{"SubscriptionID", evtSubscription.SubscriptionID}
				evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
				l.LogWithFields(ctx).Error(errorMessage)
				e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
					common.Critical, percentComplete, http.MethodDelete))
				return resp
			}
		} else {
			// Delete the host and origin resource from the respective entry
			evtSubscription.Hosts = removeElement(evtSubscription.Hosts, target.ManagerAddress)
			evtSubscription.EventDestination.OriginResources = removeLinks(evtSubscription.EventDestination.OriginResources, originResource)
			err = e.UpdateEventSubscription(evtSubscription)
			if err != nil {
				errorMessage := "Error while Updating event subscription : " + err.Error()
				msgArgs := []interface{}{"SubscriptionID", evtSubscription.SubscriptionID}
				evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
				l.LogWithFields(ctx).Error(errorMessage)
				e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
					common.Critical, percentComplete, http.MethodDelete))
				return resp
			}
		}

	}
	err = e.DeleteDeviceSubscription(searchKey)
	if err != nil {
		errorMessage := "Error while deleting device subscription : " + err.Error()
		l.LogWithFields(ctx).Error(errorMessage)
	}
	e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.OK, common.Completed,
		percentComplete, http.MethodDelete))

	resp.StatusCode = http.StatusNoContent
	resp.StatusMessage = response.ResourceRemoved
	return resp
}

// deleteSubscription to the Event Subscription
func (e *ExternalInterfaces) deleteSubscription(ctx context.Context, target *common.Target, originResource string) error {
	var plugin *common.Plugin
	plugin, err := e.GetPluginData(target.PluginID)
	if err != nil {
		return err
	}

	if _, errs := e.DeleteSubscriptions(ctx, originResource, "", plugin, target); errs != nil {
		return errs
	}
	return nil
}

// DeleteEventSubscriptionsDetails delete subscription data against given subscription id
func (e *ExternalInterfaces) DeleteEventSubscriptionsDetails(ctx context.Context, req *eventsproto.EventRequest,
	sessionUserName string, taskId string) response.RPC {
	var resp response.RPC
	var (
		percentComplete int32 = 100
		targetURI             = "/redfish/v1/EventService/Subscriptions"
	)
	reqCtx := common.CreateNewRequestContext(ctx)
	reqCtx = common.CreateMetadata(reqCtx)
	authResp, err := e.Auth(reqCtx, req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("error while trying to authenticate session: status code: %v, status message: %v", authResp.StatusCode, authResp.StatusMessage)
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		l.LogWithFields(ctx).Error(errMsg)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))
		return authResp
	}
	subscriptionDetails, err := e.GetEvtSubscriptions(req.EventSubscriptionID)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		l.LogWithFields(ctx).Error("error while deleting event subscription details : " + err.Error())
		errorMessage := err.Error()
		msgArgs := []interface{}{"SubscriptionID", req.EventSubscriptionID}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))

		return resp
	}
	if len(subscriptionDetails) < 1 {
		errorMessage := fmt.Sprintf("Subscription details not found for subscription id: %s", req.EventSubscriptionID)
		l.LogWithFields(ctx).Error(errorMessage)
		var msgArgs = []interface{}{"SubscriptionID", req.EventSubscriptionID}
		evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
			common.Critical, percentComplete, http.MethodDelete))

		return resp
	}
	for _, evtSubscription := range subscriptionDetails {
		// Since we are searching subscription id with pattern search
		// we need to match the subscription id
		if evtSubscription.SubscriptionID != req.EventSubscriptionID {
			errorMessage := fmt.Sprintf("Subscription details not found for subscription id: %s", req.EventSubscriptionID)
			l.LogWithFields(ctx).Error(errorMessage)
			var msgArgs = []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
			e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
				common.Critical, percentComplete, http.MethodDelete))

			return resp
		}

		// Delete and re subscribe Event Subscription
		isStatusAccepted, err := e.deleteAndReSubscribeToEvents(ctx, evtSubscription, req.SessionToken, sessionUserName, taskId)
		if err != nil {
			errorMessage := err.Error()
			msgArgs := []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
			e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
				common.Critical, percentComplete, http.MethodDelete))

			return resp
		}
		if isStatusAccepted {
			services.SaveEventSubscriptionID(ctx, taskId, evtSubscription.SubscriptionID)
			resp.StatusCode = http.StatusAccepted
			return resp
		}

		// Delete Event Subscription from the DB
		err = e.DeleteEvtSubscription(evtSubscription.SubscriptionID)
		if err != nil {
			l.LogWithFields(ctx).Error("error while deleting event subscription details : " + err.Error())
			errorMessage := err.Error()
			msgArgs := []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusBadRequest, msgArgs, &resp)
			e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Exception,
				common.Critical, percentComplete, http.MethodDelete))
			return resp
		}
	}

	commonResponse := response.Response{
		OdataType: common.EventDestinationType,
		OdataID:   "/redfish/v1/EventService/Subscriptions/" + req.EventSubscriptionID,
		ID:        req.EventSubscriptionID,
		Name:      "Event Subscription",
	}

	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.ResourceRemoved

	commonResponse.CreateGenericResponse(resp.StatusMessage)
	resp.Body = commonResponse
	err = e.UpdateTask(ctx, fillTaskData(taskId, targetURI, string(req.EventSubscriptionID), resp, common.Completed,
		common.OK, percentComplete, http.MethodDelete))

	return resp
}

// This function is to delete and re subscribe for Event Subscriptions
func (e *ExternalInterfaces) deleteAndReSubscribeToEvents(ctx context.Context, evtSubscription evmodel.SubscriptionResource,
	sessionToken, sessionUserName string, taskID string) (bool, error) {

	var isStatusAccepted bool
	originResources := evtSubscription.EventDestination.OriginResources
	for _, origin := range originResources {
		// ignore if origin is empty
		if origin.Oid == "" {
			continue
		}
		subscriptionDetails, err := e.GetEvtSubscriptions(origin.Oid)
		if err != nil {
			return isStatusAccepted, err
		}
		subscriptionDetails = append([]evmodel.SubscriptionResource{{
			EventDestination: &model.EventDestination{
				EventTypes: []string{"Alert"},
				Protocol:   "Redfish",
				Context:    "Creating the Default Event Subscription",
			},
		},
		}, subscriptionDetails...)
		// if origin contains fabrics then get all the collection and individual subscription details
		// for Systems need to add same later
		subscriptionDetails = e.getAllSubscriptions(origin.Oid, subscriptionDetails)
		// if delete flag is true then only one document is there
		// so don't re subscribe again

		var deleteFlag bool
		if len(subscriptionDetails) < 1 {
			return isStatusAccepted, fmt.Errorf("subscription details not found for subscription id: %s", origin)
		} else if len(subscriptionDetails) == 1 {
			deleteFlag = true
		}

		var context, protocol, destination, name string
		var eventTypes, messageIDs, resourceTypes []string

		for index, evtSub := range subscriptionDetails {
			if evtSubscription.SubscriptionID != evtSub.SubscriptionID {
				if len(evtSub.EventDestination.EventTypes) > 0 && (index == 0 || len(eventTypes) > 0) {
					eventTypes = append(eventTypes, evtSub.EventDestination.EventTypes...)
				} else {
					eventTypes = []string{}
				}

				if len(evtSub.EventDestination.MessageIds) > 0 && (index == 0 || len(messageIDs) > 0) {
					messageIDs = append(messageIDs, evtSub.EventDestination.MessageIds...)
				} else {
					messageIDs = []string{}
				}

				if len(evtSub.EventDestination.ResourceTypes) > 0 && (index == 0 || len(resourceTypes) > 0) {
					resourceTypes = append(resourceTypes, evtSub.EventDestination.ResourceTypes...)
				} else {
					resourceTypes = []string{}
				}
				name = evtSub.EventDestination.Name
				context = evtSub.EventDestination.Context
				protocol = evtSub.EventDestination.Protocol
				destination = evtSub.EventDestination.Destination
			}
		}

		removeDuplicatesFromSlice(&eventTypes)
		removeDuplicatesFromSlice(&messageIDs)
		removeDuplicatesFromSlice(&resourceTypes)

		subscriptionPost := model.EventDestination{
			Name:          name,
			EventTypes:    eventTypes,
			MessageIds:    messageIDs,
			ResourceTypes: resourceTypes,
			Context:       context,
			Protocol:      protocol,
			Destination:   destination,
		}
		isStatusAccepted, err = e.subscribe(ctx, subscriptionPost, origin.Oid, deleteFlag, sessionToken, sessionUserName, taskID, evtSubscription.SubscriptionID)
		if err != nil {
			return isStatusAccepted, err
		}
	}
	return isStatusAccepted, nil
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

// Subscribe to the Event Subscription
func (e *ExternalInterfaces) subscribe(ctx context.Context, subscriptionPost model.EventDestination, origin string,
	deleteFlag bool, sessionToken string, sessionUserName string, taskId, subscriptionId string) (bool, error) {
	var isStatusAccepted bool
	if strings.Contains(origin, "Fabrics") {
		return isStatusAccepted, e.resubscribeFabricsSubscription(ctx, subscriptionPost, origin, sessionUserName, taskId, deleteFlag)
	}
	if strings.Contains(origin, "/redfish/v1/AggregationService/Aggregates") {
		return isStatusAccepted, e.resubscribeAggregateSubscription(ctx, subscriptionPost, origin, deleteFlag, sessionToken, sessionUserName, taskId, subscriptionId)
	}
	originResource := origin
	if isCollectionOriginResourceURI(originResource) {
		isStatusAccepted = true
		l.LogWithFields(ctx).Error("Collection of origin resource:" + originResource)
		subtaskID := e.CreateSubTask(ctx, sessionUserName, taskId)
		resp := response.RPC{}
		resp.StatusCode = 200
		resp.Body = subscriptionPost
		services.SaveEventSubscriptionID(ctx, taskId, subscriptionId)

		e.UpdateTask(ctx, fillTaskData(subtaskID, originResource, "", resp, common.Completed,
			common.OK, 100, http.MethodDelete))
		return isStatusAccepted, nil
	}
	target, _, err := e.getTargetDetails(originResource)
	if err != nil {
		return isStatusAccepted, err
	}

	plugin, errs := e.GetPluginData(target.PluginID)
	if errs != nil {
		return isStatusAccepted, errs
	}
	postBody, err := json.Marshal(subscriptionPost)
	if err != nil {
		return isStatusAccepted, fmt.Errorf("error while marshalling subscription details: %s", err)
	}
	target.PostBody = postBody
	err = e.GetSubscriptionLocation(ctx, target)
	if err != nil {
		return isStatusAccepted, err
	}

	// if deleteFlag is true then only one document is there
	// so don't re subscribe again
	if deleteFlag {
		return isStatusAccepted, nil
	}

	var contactRequest evcommon.PluginContactRequest
	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(ctx, plugin)
		if token == "" {
			return isStatusAccepted, fmt.Errorf("error: Unable to create session with plugin " + plugin.ID)
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredential = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	// Note: commenting delete subscription API call and calling patch API instead
	// a subtask is created for each patch call

	subtaskID := e.CreateSubTask(ctx, sessionUserName, taskId)
	contactRequest.URL = "/ODIM/v1/Subscriptions"
	contactRequest.HTTPMethodType = http.MethodPatch
	contactRequest.PostBody = target
	createResponse, loc, _, pluginIP, err := e.PluginCall(ctx, contactRequest)
	if err != nil {
		e.UpdateTask(ctx, fillTaskData(subtaskID, contactRequest.URL, "", createResponse, common.Exception,
			common.Critical, 100, http.MethodPatch))
		return isStatusAccepted, err
	}
	if createResponse.StatusCode == http.StatusAccepted {
		isStatusAccepted = true
		services.SavePluginTaskInfo(ctx, pluginIP, plugin.IP,
			subtaskID, loc)
		return isStatusAccepted, nil
	}

	e.UpdateTask(ctx, fillTaskData(subtaskID, contactRequest.URL, "", createResponse, common.Completed,
		common.OK, 100, http.MethodPatch))

	// Update Location to all destination of device if already subscribed to the device
	var resp response.RPC
	deviceIPAddress, err := common.GetIPFromHostName(target.ManagerAddress)
	if err != nil {
		msgArgs := []interface{}{"Host", target.ManagerAddress}
		evcommon.GenErrorResponse(err.Error(), response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		l.LogWithFields(ctx).Error(err.Error())
	}
	searchKey := evcommon.GetSearchKey(deviceIPAddress, evmodel.DeviceSubscriptionIndex)
	devSub, err := e.GetDeviceSubscriptions(searchKey)
	if err != nil {
		return isStatusAccepted, err
	}
	deviceSub := common.DeviceSubscription{
		EventHostIP:     devSub.EventHostIP,
		Location:        loc,
		OriginResources: devSub.OriginResources,
	}
	return isStatusAccepted, e.UpdateDeviceSubscriptionLocation(deviceSub)

}

// DeleteFabricsSubscription will delete fabric subscription
func (e *ExternalInterfaces) DeleteFabricsSubscription(ctx context.Context, originResource string, plugin *common.Plugin) (response.RPC, error) {
	var resp response.RPC
	addr, err := GetIPFromHostNameFunc(plugin.IP)
	if err != nil {
		var msgArgs = []interface{}{"ManagerAddress", plugin.IP}
		evcommon.GenErrorResponse(err.Error(), response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
		l.LogWithFields(ctx).Error(err.Error())
		return resp, fmt.Errorf(err.Error())
	}
	searchKey := evcommon.GetSearchKey(addr, evmodel.DeviceSubscriptionIndex)
	devSub, err := e.GetDeviceSubscriptions(searchKey)
	if err != nil {

		errorMessage := "Error while get device subscription details: " + err.Error()
		if strings.Contains(err.Error(), "No data found for the key") {
			// retry the GetDeviceSubscription with plugin IP
			devSub, err = e.GetDeviceSubscriptions(plugin.IP)
			if err != nil {

				var msgArgs = []interface{}{plugin.ID + " Plugin", addr}
				evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
				l.LogWithFields(ctx).Error(errorMessage)
				return resp, err
			}
		} else {
			evcommon.GenErrorResponse(errorMessage, response.InternalError, http.StatusInternalServerError,
				[]interface{}{}, &resp)
			l.LogWithFields(ctx).Error(errorMessage)
			return resp, err
		}
	}

	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := e.getPluginToken(ctx, plugin)
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
	resp, _, _, _, err = e.PluginCall(ctx, contactRequest)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		resp, _, _, _, err = e.retryEventOperation(ctx, contactRequest)
		if err != nil {
			return resp, err
		}
	}
	return resp, nil
}

// resubscribeFabricsSubscription updates subscription fabric subscription details  by forming the super set of MessageIDs,EventTypes and ResourceTypes
func (e *ExternalInterfaces) resubscribeFabricsSubscription(ctx context.Context, subscriptionPost model.EventDestination,
	sessionUserName string, taskId string, origin string, deleteFlag bool) error {
	originResources := e.getSubordinateResourcesFromCollection(origin)
	for _, origin := range originResources {
		originResource := origin
		fabricID := getFabricID(originResource)
		if fabricID == "" {
			return nil
		}
		// get Fabrics Details
		fabric, dberr := e.GetFabricData(fabricID)
		if dberr != nil {
			errorMessage := "error while getting fabric data: " + dberr.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
		plugin, errs := e.GetPluginData(fabric.PluginID)
		if errs != nil {
			errorMessage := "error while getting plugin data: " + errs.Error()
			l.LogWithFields(ctx).Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
		// Deleting the fabric subscription
		resp, err := e.DeleteFabricsSubscription(ctx, origin, plugin)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}
			return err
		}
		// if deleteFlag is true then only one document is there
		// so don't re subscribe again
		if deleteFlag {
			return nil
		}
		var contactRequest evcommon.PluginContactRequest

		contactRequest.Plugin = plugin
		if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			token := e.getPluginToken(ctx, plugin)
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
		subscriptionPost.OriginResources = []model.Link{
			{
				Oid: originResource,
			},
		}
		postBody, _ := json.Marshal(subscriptionPost)
		var reqData string
		//replacing the request url with south bound translation URL
		for key, value := range config.Data.URLTranslation.SouthBoundURL {
			reqData = strings.Replace(string(postBody), key, value, -1)
		}

		subtaskID := e.CreateSubTask(ctx, sessionUserName, taskId)
		// recreating the subscription
		contactRequest.URL = "/ODIM/v1/Subscriptions"
		contactRequest.HTTPMethodType = http.MethodPost
		err = json.Unmarshal([]byte(reqData), &contactRequest.PostBody)
		if err != nil {
			return err
		}
		l.LogWithFields(ctx).Info("Resubscribe request" + reqData)
		response, loc, _, pluginIP, err := e.PluginCall(ctx, contactRequest)
		if err != nil {
			e.UpdateTask(ctx, fillTaskData(subtaskID, contactRequest.URL, "", response, common.Exception,
				common.Critical, 100, http.MethodPost))
			return err
		}

		if response.StatusCode == http.StatusAccepted {
			services.SavePluginTaskInfo(ctx, pluginIP, plugin.IP,
				subtaskID, loc)
			return nil
		}

		if response.StatusCode == http.StatusUnauthorized && strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			_, _, _, _, err = e.retryEventOperation(ctx, contactRequest)
			if err != nil {
				e.UpdateTask(ctx, fillTaskData(subtaskID, contactRequest.URL, "", response, common.Exception,
					common.Critical, 100, http.MethodPost))
				return err
			}
		}

		e.UpdateTask(ctx, fillTaskData(subtaskID, contactRequest.URL, "", response, common.ComputerSystemReset,
			common.OK, 100, http.MethodPost))

		addr, err := common.GetIPFromHostName(plugin.IP)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		searchKey := evcommon.GetSearchKey(addr, evmodel.DeviceSubscriptionIndex)
		// Update Location to all destination of device if already subscribed to the device
		devSub, err := e.GetDeviceSubscriptions(searchKey)
		if err != nil {
			return err
		}
		deviceSub := common.DeviceSubscription{
			EventHostIP:     devSub.EventHostIP,
			Location:        loc,
			OriginResources: devSub.OriginResources,
		}
		err = e.UpdateDeviceSubscriptionLocation(deviceSub)
		if err != nil {
			return err
		}
	}

	return nil
}

// getSubordinateResourcesFromCollection method return sub uri
func (e *ExternalInterfaces) getSubordinateResourcesFromCollection(originResources string) []string {
	data, _, collectionPresentFlag, _, _, _ := e.checkCollection(originResources)
	if !collectionPresentFlag {
		return []string{originResources}
	}
	return data
}

// getAllSubscriptions return list of subscription fabric resource
func (e *ExternalInterfaces) getAllSubscriptions(origin string, subscriptionDetails []evmodel.SubscriptionResource) []evmodel.SubscriptionResource {
	if origin == "/redfish/v1/Fabrics" {
		return subscriptionDetails
	}

	searchKey := "/redfish/v1/Fabrics"
	subscriptions, err := e.GetEvtSubscriptions(searchKey)
	if err != nil {
		return subscriptionDetails
	}
	// Checking the collection based subscription
	var collectionSubscription = make([]evmodel.SubscriptionResource, 0)
	for _, evtSubscription := range subscriptions {
		for _, originResource := range evtSubscription.EventDestination.OriginResources {
			if strings.Contains(origin, "Fabrics") && originResource.Oid == "/redfish/v1/Fabrics" {
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
func removeDuplicatesFromSubscription(subscriptions []evmodel.SubscriptionResource) []evmodel.SubscriptionResource {
	uniqueElementsDs := make(map[string]bool)
	var uniqueElementsList []evmodel.SubscriptionResource
	for _, sub := range subscriptions {
		if exist := uniqueElementsDs[sub.SubscriptionID]; !exist {
			uniqueElementsList = append(uniqueElementsList, sub)
			uniqueElementsDs[sub.SubscriptionID] = true
		}
	}
	return uniqueElementsList
}

// DeleteAggregateSubscriptions it will add subscription for newly Added system in aggregate
func (e *ExternalInterfaces) DeleteAggregateSubscriptions(ctx context.Context, req *eventsproto.EventUpdateRequest, isRemove bool) error {
	var aggregateID = req.AggregateId
	searchKeyAgg := evcommon.GetSearchKey(aggregateID, evmodel.SubscriptionIndex)
	subscriptionList, err := e.GetEvtSubscriptions(searchKeyAgg)
	if err != nil {
		l.LogWithFields(ctx).Error("no aggregate subscription found ", err)
		return err
	}
	for _, evtSubscription := range subscriptionList {
		evtSubscription.Hosts = removeElement(evtSubscription.Hosts, aggregateID)
		evtSubscription.EventDestination.OriginResources = removeLinks(evtSubscription.EventDestination.OriginResources, "/redfish/v1/AggregationService/Aggregates/"+aggregateID)

		if len(evtSubscription.EventDestination.OriginResources) == 0 {
			err = e.DeleteEvtSubscription(evtSubscription.SubscriptionID)
			if err != nil {
				errorMessage := "error while delete event subscription : " + err.Error()
				l.LogWithFields(ctx).Error(errorMessage)
				return err
			}
		} else {
			err = e.UpdateEventSubscription(evtSubscription)
			if err != nil {
				errorMessage := "error while updating event subscription : " + err.Error()
				l.LogWithFields(ctx).Error(errorMessage)
				return err
			}
		}
	}
	return nil
}

// getAggregateSystemList return  list of system to corresponding aggregate
func getAggregateSystemList(origin string, sessionToken string) ([]model.Link, error) {
	conn, err := services.ODIMService.Client(services.Aggregator)
	if err != nil {
		return nil, err
	}
	aggregator := aggregatorproto.NewAggregatorClient(conn)
	var req aggregatorproto.AggregatorRequest

	req.URL = origin
	req.SessionToken = sessionToken
	resp, err := aggregator.GetAggregate(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %v", err)
	}
	var data evmodel.Aggregate
	err = json.Unmarshal(resp.Body, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid json: %v", err)
	}
	defer conn.Close()
	return data.Elements, nil
}

// resubscribeAggregateSubscription method subscribe event for
// aggregate system members
func (e *ExternalInterfaces) resubscribeAggregateSubscription(ctx context.Context, subscriptionPost model.EventDestination,
	origin string, deleteFlag bool, sessionToken string, sessionUserName string, taskId, subscriptionId string) error {
	originResource := origin
	systems, err := getAggregateSystemList(originResource, sessionToken)
	if err != nil {
		return nil
	}
	for _, system := range systems {
		_, err = e.subscribe(ctx, subscriptionPost, system.Oid, deleteFlag, sessionToken, sessionUserName, taskId, subscriptionId)
		if err != nil {
			return err
		}
	}
	return nil
}
