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
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"gopkg.in/go-playground/validator.v9"
)

// ExternalInterfaces struct to inject the pmb client function into the handlers
type ExternalInterfaces struct {
	External
	DB
}

var (
	// UpdateTaskService function  pointer for calling the files
	UpdateTaskService = services.UpdateTask
	// IOUtilReadAllFunc function  pointer for calling the files
	IOUtilReadAllFunc = ioutil.ReadAll
)

// External struct to inject the contact external function into the handlers
type External struct {
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth            func(string, []string, []string) (response.RPC, error)
	CreateTask      func(string) (string, error)
	UpdateTask      func(common.TaskData) error
	CreateChildTask func(string, string) (string, error)
}

// DB struct to inject the contact DB function into the handlers
type DB struct {
	GetSessionUserName               func(sessionToken string) (string, error)
	GetEvtSubscriptions              func(string) ([]evmodel.SubscriptionResource, error)
	SaveEventSubscription            func(evmodel.SubscriptionResource) error
	GetPluginData                    func(string) (*evmodel.Plugin, *errors.Error)
	GetDeviceSubscriptions           func(string) (*evmodel.DeviceSubscription, error)
	GetTarget                        func(string) (*evmodel.Target, error)
	GetAllKeysFromTable              func(string) ([]string, error)
	GetAllFabrics                    func() ([]string, error)
	GetAllMatchingDetails            func(string, string, common.DbType) ([]string, *errors.Error)
	UpdateDeviceSubscriptionLocation func(evmodel.DeviceSubscription) error
	GetFabricData                    func(string) (evmodel.Fabric, error)
	DeleteEvtSubscription            func(string) error
	DeleteDeviceSubscription         func(hostIP string) error
	UpdateEventSubscription          func(evmodel.SubscriptionResource) error
	SaveUndeliveredEvents            func(string, []byte) error
	SaveDeviceSubscription           func(evmodel.DeviceSubscription) error
	GetUndeliveredEvents             func(string) (string, error)
	DeleteUndeliveredEvents          func(string) error
	GetUndeliveredEventsFlag         func(string) (bool, error)
	SetUndeliveredEventsFlag         func(string) error
	DeleteUndeliveredEventsFlag      func(string) error
	GetAggregateData                 func(string) (evmodel.Aggregate, error)
	SaveAggregateSubscription        func(aggregateId string, hostIP []string) error
	GetAggregateHosts                func(aggregateIP string) ([]string, error)
	UpdateAggregateHosts             func(aggregateId string, hostIP []string) error
	GetAggregateList                 func(hostIP string) ([]string, error)
}

// fillTaskData is to fill task information in TaskData struct
func fillTaskData(taskID, targetURI, request string, resp errResponse.RPC, taskState string, taskStatus string, percentComplete int32, httpMethod string) common.TaskData {
	return common.TaskData{
		TaskID:          taskID,
		TargetURI:       targetURI,
		Response:        resp,
		TaskRequest:     request,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		PercentComplete: percentComplete,
		HTTPMethod:      httpMethod,
	}
}

// UpdateTaskData update the task with the given data
func UpdateTaskData(taskData common.TaskData) error {
	respBody, _ := json.Marshal(taskData.Response.Body)
	payLoad := &taskproto.Payload{
		HTTPHeaders:   taskData.Response.Header,
		HTTPOperation: taskData.HTTPMethod,
		JSONBody:      taskData.TaskRequest,
		StatusCode:    taskData.Response.StatusCode,
		TargetURI:     taskData.TargetURI,
		ResponseBody:  respBody,
	}

	err := UpdateTaskService(taskData.TaskID, taskData.TaskState, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		UpdateTaskService(taskData.TaskID, common.Cancelled, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
		if taskData.PercentComplete == 0 {
			return fmt.Errorf("error while starting the task: %v", err)
		}
		l.Log.Error("error: task update for " + taskData.TaskID + " failed with err: " + err.Error())
		runtime.Goexit()
	}
	return nil
}

// this function is for to create array of originofresources without odata id
func removeOdataIDfromOriginResources(originResources []evmodel.OdataIDLink) []string {
	var originRes []string
	for _, origin := range originResources {
		originRes = append(originRes, origin.OdataID)
	}
	return originRes
}

// remove duplicate elements in string slice.
// Takes string slice and length, and updates the same with new values
func removeDuplicatesFromSlice(slc *[]string, slcLen *int) {
	if *slcLen > 1 {
		uniqueElementsDs := make(map[string]bool)
		var uniqueElementsList []string
		for _, element := range *slc {
			if exist := uniqueElementsDs[element]; !exist {
				uniqueElementsList = append(uniqueElementsList, element)
				uniqueElementsDs[element] = true
			}
		}
		// length of uniqueElementsList will be less than passed string slice,
		// only if duplicates existed, so will assign slc with modified list and update length
		if len(uniqueElementsList) < *slcLen {
			*slc = uniqueElementsList
			*slcLen = len(*slc)
		}
	}
	return
}

// removeElement will remove the element from the slice return
// slice of remaining elements
func removeElement(slice []string, element string) []string {
	var elements []string
	for _, val := range slice {
		if val != element {
			elements = append(elements, val)
		}
	}
	return elements
}

// PluginCall method is to call to given url and method
// and validate the response and return
func (e *ExternalInterfaces) PluginCall(req evcommon.PluginContactRequest) (errResponse.RPC, string, string, error) {
	var resp errResponse.RPC
	response, err := e.callPlugin(req)
	if err != nil {
		if evcommon.GetPluginStatus(req.Plugin) {
			response, err = e.callPlugin(req)
		}
		if err != nil {
			errorMessage := "Error : " + err.Error()
			evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
				[]interface{}{}, &resp)
			l.Log.Error(errorMessage)
			return resp, "", "", err
		}
	}
	defer response.Body.Close()
	body, err := IOUtilReadAllFunc(response.Body)
	if err != nil {
		errorMessage := "error while trying to read response body: " + err.Error()
		evcommon.GenErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			[]interface{}{}, &resp)
		l.Log.Error(errorMessage)
		return resp, "", "", err
	}
	if !(response.StatusCode == http.StatusCreated || response.StatusCode == http.StatusOK) {
		resp.StatusCode = int32(response.StatusCode)
		resp.Body = string(body)
		return resp, "", "", err
	}
	var outBody interface{}
	json.Unmarshal(body, &outBody)
	resp.StatusCode = int32(response.StatusCode)
	resp.Body = outBody
	return resp, response.Header.Get("location"), response.Header.Get("X-Auth-Token"), nil
}

// validateFields is for validating subscription parameters
func validateFields(request *evmodel.RequestBody) (int32, string, []interface{}, error) {
	validEventFormatTypes := map[string]bool{"Event": true, "MetricReport": true}
	validEventTypes := map[string]bool{"Alert": true, "MetricReport": true, "ResourceAdded": true, "ResourceRemoved": true, "ResourceUpdated": true, "StatusChange": true, "Other": true}

	validate := validator.New()

	// if any of the mandatory fields missing in the struct, then it return an error
	err := validate.Struct(request)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{err.Field()}, fmt.Errorf(err.Field() + " field is missing")
		}
	}
	if request.EventFormatType == "" {
		request.EventFormatType = "Event"
	}

	if _, ok := validEventFormatTypes[request.EventFormatType]; !ok {
		return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{request.EventFormatType, "EventFormatType"}, fmt.Errorf("Invalid EventFormatType")
	}

	if len(request.EventTypes) == 0 && request.EventFormatType == "MetricReport" {
		request.EventTypes = []string{"MetricReport"}
	}

	for _, eventType := range request.EventTypes {
		if _, ok := validEventTypes[eventType]; !ok {
			return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{eventType, "EventTypes"}, fmt.Errorf("Invalid EventTypes")
		}
	}

	if request.EventFormatType == "MetricReport" {
		if len(request.EventTypes) > 1 {
			return http.StatusBadRequest, errResponse.PropertyValueFormatError, []interface{}{request.EventFormatType, "EventTypes"}, fmt.Errorf("Unsupported EventType")
		}
		if request.EventTypes[0] != "MetricReport" {
			return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{request.EventTypes[0], "EventType"}, fmt.Errorf("Unsupported EventType")
		}
	}

	// if request.SubscriptionType == "" {
	// 	request.SubscriptionType = evmodel.SubscriptionType
	// } else {
	// 	isValid, errorMessage := request.SubscriptionType.IsValidSubscriptionType()
	// 	if !isValid {
	// 		return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf(errorMessage)
	// 	}
	// }
	if !request.SubscriptionType.IsSubscriptionTypeSupported() {
		if request.SubscriptionType.IsValidSubscriptionType() {
			return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf("unsupported SubscriptionType")
		} else {
			return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf("invalid SubscriptionType")
		}

	}

	if request.Context == "" {
		request.Context = evmodel.Context
	}

	if request.DeliveryRetryPolicy == "" {
		request.DeliveryRetryPolicy = evmodel.DeliveryRetryPolicy
	} else {
		if !request.DeliveryRetryPolicy.IsDeliveryRetryPolicyTypeSupported() {
			if request.DeliveryRetryPolicy.IsValidDeliveryRetryPolicyType() {
				return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf("unsupported DeliveryRetryPolicy")

			} else {
				return http.StatusBadRequest, errResponse.PropertyMissing, []interface{}{"SubscriptionType"}, fmt.Errorf("invalid DeliveryRetryPolicy")
			}
		}
	}
	availableProtocols := []string{"Redfish"}
	var validProtocol bool
	validProtocol = false
	for _, protocol := range availableProtocols {
		if request.Protocol == protocol {
			validProtocol = true
		}
	}
	if !validProtocol {
		return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{request.Protocol, "Protocol"}, fmt.Errorf("Protocol %v is invalid", request.Protocol)
	}

	// check the All ResourceTypes are supported
	for _, resourceType := range request.ResourceTypes {
		if _, ok := common.ResourceTypes[resourceType]; !ok {
			return http.StatusBadRequest, errResponse.PropertyValueNotInList, []interface{}{resourceType, "ResourceType"}, fmt.Errorf("Unsupported ResourceType")
		}
	}

	return http.StatusOK, common.OK, []interface{}{}, nil
}

// GetUUID fetches the UUID from the Origin Resource
func getUUID(origin string) (string, error) {
	var uuid string
	requestData := strings.SplitN(origin, ".", 2)
	if len(requestData) <= 1 {
		return "", fmt.Errorf("error: SystemUUID not found")
	}
	resource := requestData[0]
	uuid = resource[strings.LastIndexByte(resource, '/')+1:]
	return uuid, nil
}

func createEventSubscriptionResponse() interface{} {
	return errors.ErrorClass{
		MessageExtendedInfo: []errors.MsgExtendedInfo{
			errors.MsgExtendedInfo{
				MessageID: response.Created,
			},
		},
		Code:    errResponse.Created,
		Message: "See @Message.ExtendedInfo for more information.",
	}
}

// getPluginToken will verify the if any token present to the plugin else it will create token for the new plugin
func (e *ExternalInterfaces) getPluginToken(plugin *evmodel.Plugin) string {
	authToken := evcommon.Token.GetToken(plugin.ID)
	if authToken == "" {
		return e.createToken(plugin)
	}
	return authToken
}

func (e *ExternalInterfaces) createToken(plugin *evmodel.Plugin) string {
	var contactRequest evcommon.PluginContactRequest

	contactRequest.Plugin = plugin
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = map[string]interface{}{
		"Username": plugin.Username,
		"Password": string(plugin.Password),
	}
	contactRequest.URL = "/ODIM/v1/Sessions"
	_, _, token, err := e.PluginCall(contactRequest)
	if err != nil {
		l.Log.Error(err.Error())
	}
	pluginToken := evcommon.PluginToken{
		Tokens: make(map[string]string),
	}
	if token != "" {
		pluginToken.StoreToken(plugin.ID, token)
	}
	return token
}

func (e *ExternalInterfaces) retryEventOperation(req evcommon.PluginContactRequest) (errResponse.RPC, string, string, error) {
	var resp errResponse.RPC
	var token = e.createToken(req.Plugin)
	if token == "" {
		evcommon.GenErrorResponse("error: Unable to create session with plugin "+req.Plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
			[]interface{}{}, &resp)
		return resp, "", "", fmt.Errorf("error: Unable to create session with plugin")
	}
	req.Token = token
	return e.PluginCall(req)
}

func (e *ExternalInterfaces) retryEventSubscriptionOperation(req evcommon.PluginContactRequest) (*http.Response, evresponse.EventResponse, error) {
	var resp evresponse.EventResponse
	var token = e.createToken(req.Plugin)
	if token == "" {
		evcommon.GenEventErrorResponse("error: Unable to create session with plugin "+req.Plugin.ID, errResponse.NoValidSession, http.StatusUnauthorized,
			&resp, []interface{}{})
		return nil, resp, fmt.Errorf("error: Unable to create session with plugin")
	}
	req.Token = token

	response, err := e.callPlugin(req)
	if err != nil {
		errorMessage := "error while unmarshaling the body : " + err.Error()
		evcommon.GenEventErrorResponse(errorMessage, errResponse.InternalError, http.StatusInternalServerError,
			&resp, []interface{}{})
		l.Log.Error(errorMessage)
		return nil, resp, err
	}
	return response, resp, err
}

// isHostPresent will check if hostip present in the hosts slice
func isHostPresent(hosts []string, hostip string) bool {

	if len(hosts) < 1 {
		return false
	}

	front := 0
	rear := len(hosts) - 1
	for front <= rear {
		if hosts[front] == hostip || hosts[rear] == hostip {
			return true
		}
		front++
		rear--
	}
	return false
}

func getFabricID(origin string) string {
	data := strings.Split(origin, "/redfish/v1/Fabrics/")
	if len(data) > 1 {
		fabricData := strings.Split(data[1], "/")
		return fabricData[0]
	}
	return ""
}
func getAggregateID(origin string) string {
	data := strings.Split(origin, "/redfish/v1/AggregationService/Aggregates/")
	if len(data) > 1 {
		fabricData := strings.Split(data[1], "/")
		return fabricData[0]
	}
	return ""
}

// callPlugin check the given request url and PreferAuth type plugin
func (e *ExternalInterfaces) callPlugin(req evcommon.PluginContactRequest) (*http.Response, error) {
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.URL
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		return e.ContactClient(reqURL, req.HTTPMethodType, "", "", req.PostBody, req.LoginCredential)
	}
	return e.ContactClient(reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, nil)
}

// checkCollection verifies if the given origin is collection and extracts all the suboridinate resources
func (e *ExternalInterfaces) checkCollection(origin string) ([]string, string, bool, string, bool, error) {
	switch origin {
	case "/redfish/v1/Systems":
		collection, err := e.GetAllKeysFromTable("ComputerSystem")
		return collection, "SystemsCollection", true, "", false, err
	case "/redfish/v1/Chassis":
		return []string{}, "ChassisCollection", true, "", false, nil
	case "/redfish/v1/Managers":
		//TODO:After Managers implemention need to get all Managers data
		return []string{}, "ManagerCollection", true, "", false, nil
	case "/redfish/v1/Fabrics":
		collection, err := e.GetAllFabrics()
		return collection, "FabricsCollection", true, "", false, err
	case "/redfish/v1/TaskService/Tasks":
		return []string{}, "TasksCollection", true, "", false, nil
	}
	if strings.Contains(origin, "/AggregationService/Aggregates/") {
		aggregateCollection, err := e.GetAggregateData(origin)
		if err != nil {
			return []string{}, "AggregateCollections", true, "", false, err
		}
		var collection []string = []string{}
		for _, system := range aggregateCollection.Elements {
			var systemID string = system.OdataID
			collection = append(collection, systemID)
		}
		return collection, "AggregateCollections", true, origin, true, err
	}

	return []string{}, "", false, "", false, nil
}

// isHostPresentInEventForward will check if hostip present in the hosts slice
func isHostPresentInEventForward(hosts []string, hostip string) bool {

	if len(hosts) == 0 {
		return true
	}

	front := 0
	rear := len(hosts) - 1
	for front <= rear {
		if hosts[front] == hostip || hosts[rear] == hostip || strings.Contains(hosts[rear], "Collection") || strings.Contains(hosts[front], "Collection") {
			return true
		}
		front++
		rear--
	}
	return false
}

// updateOriginResourceswithOdataID is for to create array of odata id
func updateOriginResourcesWithOdataID(originResources []string) []evresponse.ListMember {
	var originRes []evresponse.ListMember
	for _, origin := range originResources {
		originRes = append(originRes, evresponse.ListMember{OdataID: origin})
	}
	return originRes
}
