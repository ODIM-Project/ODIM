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

package system

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agcommon"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmessagebus"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

//Device struct to define the response from plugin for UUID
type Device struct {
	ServerIP   string `json:"ServerIP"`
	Username   string `json:"Username"`
	DeviceUUID string `json:"device_UUID"`
}

// ExternalInterface struct holds the function pointers all outboud services
type ExternalInterface struct {
	ContactClient            func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	Auth                     func(string, []string, []string) response.RPC
	GetSessionUserName       func(string) (string, error)
	CreateChildTask          func(string, string) (string, error)
	CreateTask               func(string) (string, error)
	UpdateTask               func(common.TaskData) error
	CreateSubcription        func([]string)
	PublishEvent             func([]string, string)
	PublishEventMB           func(string, string, string)
	GetPluginStatus          func(agmodel.Plugin) bool
	SubscribeToEMB           func(string, []string)
	EncryptPassword          func([]byte) ([]byte, error)
	DecryptPassword          func([]byte) ([]byte, error)
	DeleteComputeSystem      func(int, string) *errors.Error
	DeleteSystem             func(string) *errors.Error
	DeleteEventSubscription  func(string) (*eventsproto.EventSubResponse, error)
	EventNotification        func(string, string, string)
	GetAllKeysFromTable      func(string) ([]string, error)
	GetConnectionMethod      func(string) (agmodel.ConnectionMethod, *errors.Error)
	UpdateConnectionMethod   func(agmodel.ConnectionMethod, string) *errors.Error
	GetPluginMgrAddr         func(string) (agmodel.Plugin, *errors.Error)
	GetAggregationSourceInfo func(string) (agmodel.AggregationSource, *errors.Error)
	GenericSave              func([]byte, string, string) error
	CheckActiveRequest       func(string) (bool, *errors.Error)
	DeleteActiveRequest      func(string) *errors.Error
}

type responseStatus struct {
	StatusCode    int32
	StatusMessage string
	MsgArgs       []interface{}
}

type getResourceRequest struct {
	Data              []byte
	Username          string
	Password          string
	SystemID          string
	DeviceUUID        string
	DeviceInfo        interface{}
	LoginCredentials  map[string]string
	ParentOID         string
	OID               string
	ContactClient     func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	OemFlag           bool
	Plugin            agmodel.Plugin
	TaskRequest       string
	HTTPMethodType    string
	Token             string
	StatusPoll        bool
	CreateSubcription func([]string)
	PublishEvent      func([]string, string)
	GetPluginStatus   func(agmodel.Plugin) bool
	UpdateFlag        bool
	TargetURI         string
	UpdateTask        func(common.TaskData) error
}

type respHolder struct {
	ErrorMessage   string
	StatusCode     int32
	StatusMessage  string
	MsgArgs        []interface{}
	lock           sync.Mutex
	SystemURL      []string
	PluginResponse string
	TraversedLinks map[string]bool
}

//AddResourceRequest is payload of adding a  resource
type AddResourceRequest struct {
	ManagerAddress   string            `json:"ManagerAddress"`
	UserName         string            `json:"UserName"`
	Password         string            `json:"Password"`
	ConnectionMethod *ConnectionMethod `json:"ConnectionMethod"`
}

//ConnectionMethod struct definition for @odata.id
type ConnectionMethod struct {
	OdataID string `json:"@odata.id"`
}

// TaskData holds the data of the Task
type TaskData struct {
	TaskID          string
	TargetURI       string
	Resp            response.RPC
	TaskState       string
	TaskStatus      string
	PercentComplete int32
	HTTPMethod      string
}

// ActiveRequestsSet holds details of ongoing requests
type ActiveRequestsSet struct {
	// ReqRecord holds data of ongoing requests
	ReqRecord map[string]interface{}
	// UpdateMu is the mutex for protecting OngoingReqs
	UpdateMu sync.Mutex
}

var southBoundURL = "southboundurl"
var northBoundURL = "northboundurl"

// AggregationSource  payload of adding a  AggregationSource
type AggregationSource struct {
	HostName string `json:"HostName"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	Links    *Links `json:"Links,omitempty"`
}

// Links holds information of Oem
type Links struct {
	ConnectionMethod *ConnectionMethod `json:"ConnectionMethod,omitempty"`
}

type connectionMethodVariants struct {
	PluginType        string
	PreferredAuthType string
	PluginID          string
	FirmwareVersion   string
}

func getIPAndPortFromAddress(address string) (string, string) {
	ip, port, err := net.SplitHostPort(address)
	if err != nil {
		ip = address
	}
	return ip, port
}

func getKeyFromManagerAddress(managerAddress string) string {
	ipAddr, host, port, err := agcommon.LookupHost(managerAddress)
	if err != nil {
		ipAddr = host
	}
	if port != "" {
		return net.JoinHostPort(host, port)
	}
	return ipAddr
}

func fillTaskData(taskID, targetURI, request string, resp response.RPC, taskState string, taskStatus string, percentComplete int32, httpMethod string) common.TaskData {
	return common.TaskData{
		TaskID:          taskID,
		TargetURI:       targetURI,
		TaskRequest:     request,
		Response:        resp,
		TaskState:       taskState,
		TaskStatus:      taskStatus,
		PercentComplete: percentComplete,
		HTTPMethod:      httpMethod,
	}
}

//genError generates error response so as to reduce boiler plate code
func genError(errorMessage string, respPtr *response.RPC, httpStatusCode int32, StatusMessage string, header map[string]string) {
	respPtr.StatusCode = httpStatusCode
	respPtr.StatusMessage = StatusMessage
	respPtr.Body = errors.CreateErrorResponse(respPtr.StatusMessage, errorMessage)
	respPtr.Header = header
	log.Error(errorMessage)
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

	err := services.UpdateTask(taskData.TaskID, taskData.TaskState, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
	if err != nil && (err.Error() == common.Cancelling) {
		// We cant do anything here as the task has done it work completely, we cant reverse it.
		//Unless if we can do opposite/reverse action for delete server which is add server.
		services.UpdateTask(taskData.TaskID, common.Cancelled, taskData.TaskStatus, taskData.PercentComplete, payLoad, time.Now())
		if taskData.PercentComplete == 0 {
			return fmt.Errorf("error while starting the task: %v", err)
		}
		runtime.Goexit()
		//		return fmt.Errorf(common.Cancelled)
	}
	return nil
}

func contactPlugin(req getResourceRequest, errorMessage string) ([]byte, string, responseStatus, error) {
	var resp responseStatus
	pluginResp, err := callPlugin(req)
	if err != nil {
		if req.StatusPoll {
			if req.GetPluginStatus(req.Plugin) {
				pluginResp, err = callPlugin(req)
			}
		}
		if err != nil {
			errorMessage = errorMessage + err.Error()
			resp.StatusCode = http.StatusServiceUnavailable
			resp.StatusMessage = response.CouldNotEstablishConnection
			resp.MsgArgs = []interface{}{"https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.OID}
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
	}

	defer pluginResp.Body.Close()
	body, err := ioutil.ReadAll(pluginResp.Body)
	if err != nil {
		errorMessage := "error while trying to read plugin response body: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		return nil, "", resp, fmt.Errorf(errorMessage)
	}

	if pluginResp.StatusCode != http.StatusCreated && pluginResp.StatusCode != http.StatusOK {
		if pluginResp.StatusCode == http.StatusUnauthorized {
			errorMessage += "error: invalid resource username/password"
			resp.StatusCode = int32(pluginResp.StatusCode)
			resp.StatusMessage = response.ResourceAtURIUnauthorized
			resp.MsgArgs = []interface{}{"https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.OID}
			log.Error(errorMessage)
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
		errorMessage += string(body)
		resp.StatusCode = int32(pluginResp.StatusCode)
		resp.StatusMessage = response.InternalError
		log.Error(errorMessage)
		return body, "", resp, fmt.Errorf(errorMessage)
	}

	data := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range getTranslationURL(northBoundURL) {
		data = strings.Replace(data, key, value, -1)
	}
	resp.StatusCode = int32(pluginResp.StatusCode)
	return []byte(data), pluginResp.Header.Get("X-Auth-Token"), resp, nil
}

// keyFormation is to form the key to insert in DB
func keyFormation(oid, systemID, DeviceUUID string) string {
	if oid[len(oid)-1:] == "/" {
		oid = oid[:len(oid)-1]
	}
	str := strings.Split(oid, "/")
	var key []string
	for i, id := range str {
		if id == systemID && (strings.EqualFold(str[i-1], "Systems") || strings.EqualFold(str[i-1], "Chassis") || strings.EqualFold(str[i-1], "Managers") || strings.EqualFold(str[i-1], "FirmwareInventory") || strings.EqualFold(str[i-1], "SoftwareInventory")) {
			key = append(key, DeviceUUID+":"+id)
			continue
		}
		key = append(key, id)
	}
	return strings.Join(key, "/")
}

func (h *respHolder) getAllSystemInfo(taskID string, progress int32, alottedWork int32, req getResourceRequest) (string, string, int32, error) {
	var computeSystemID, resourceURI string
	body, _, getResponse, err := contactPlugin(req, "error while trying to get system collection details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()

		h.StatusMessage = getResponse.StatusMessage
		h.StatusCode = getResponse.StatusCode
		h.MsgArgs = getResponse.MsgArgs
		h.lock.Unlock()
		log.Error(err)
		return computeSystemID, resourceURI, progress, err
	}
	h.SystemURL = make([]string, 0)
	h.PluginResponse = string(body)
	systemsMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &systemsMap)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal systems collection: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		log.Error("error while trying unmarshal systems collection: " + err.Error())
		return computeSystemID, resourceURI, progress, err
	}
	systemMembers := systemsMap["Members"]
	// Loop through System collection members and discover all of them
	errorMessage := "error : get system collection members failed for ["
	foundErr := false
	for _, object := range systemMembers.([]interface{}) {
		estimatedWork := alottedWork / int32(len(systemMembers.([]interface{})))
		oDataID := object.(map[string]interface{})["@odata.id"].(string)
		req.OID = oDataID
		if computeSystemID, resourceURI, progress, err = h.getSystemInfo(taskID, progress, estimatedWork, req); err != nil {
			errorMessage += oDataID + ":err-" + err.Error() + "; "
			foundErr = true
		}
	}
	if foundErr {
		return computeSystemID, resourceURI, progress, fmt.Errorf("%s]", errorMessage)
	}
	return computeSystemID, resourceURI, progress, nil
}

//Registries Discovery function
func (h *respHolder) getAllRegistries(taskID string, progress int32, alottedWork int32, req getResourceRequest) int32 {

	// Get all available file names in the registry store directory in a list
	registryStore := config.Data.RegistryStorePath
	regFiles, err := ioutil.ReadDir(registryStore)
	if err != nil {
		log.Error("error while reading the files from directory " + registryStore + ": " + err.Error())
		log.Fatal(err)
	}
	//Construct the list of file names available
	var standardFiles []string
	for _, regFile := range regFiles {
		standardFiles = append(standardFiles, regFile.Name())
	}

	body, _, getResponse, err := contactPlugin(req, "error while trying to get the Registries collection  details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		h.StatusMessage = getResponse.StatusMessage
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		log.Error(err)
		return progress
	}
	registriesMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &registriesMap)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal Registries collection: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		log.Error("error while trying to unmarshal Registries collection: " + err.Error())
		return progress

	}
	registriesMembers := registriesMap["Members"]
	// Loop through all the registry members collection and discover all of them
	for _, object := range registriesMembers.([]interface{}) {
		estimatedWork := alottedWork / int32(len(registriesMembers.([]interface{})))
		if object == nil {
			progress = progress + estimatedWork
			continue
		}
		oDataIDInterface := object.(map[string]interface{})["@odata.id"]
		if oDataIDInterface == nil {
			progress = progress + estimatedWork
			continue
		}
		oDataID := oDataIDInterface.(string)
		req.OID = oDataID
		progress = h.getRegistriesInfo(taskID, progress, estimatedWork, standardFiles, req)
	}
	return progress
}

func (h *respHolder) getRegistriesInfo(taskID string, progress int32, allotedWork int32, standardFiles []string, req getResourceRequest) int32 {
	body, _, getResponse, err := contactPlugin(req, "error while trying to get Registry fileinfo details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		h.StatusMessage = getResponse.StatusMessage
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		return progress
	}
	var registryFileInfo map[string]interface{}
	err = json.Unmarshal(body, &registryFileInfo)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal response body: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return progress
	}
	uri := ""
	/* '#' charactor in the begining of the registryfile name is giving some issue
	* during api routing. So getting Id instead of Registry name if it has '#' char as a
	* prefix.
	 */
	registryNameInterface := registryFileInfo["Registry"]
	// If Registry field is not present, then nothing to discover.So return progress.
	if registryNameInterface == nil {
		return progress + allotedWork
	}
	registryName := registryNameInterface.(string)
	if strings.HasPrefix(registryName, "#") {
		registryName = registryFileInfo["Id"].(string)
	}
	// Check if file not exist go get ut and store in DB
	if isFileExist(standardFiles, registryName+".json") == true {
		return progress + allotedWork
	}
	locations := registryFileInfo["Location"]
	for _, location := range locations.([]interface{}) {
		if location == nil {
			continue
		}
		languageInterface := location.(map[string]interface{})["Language"]
		if languageInterface == nil {
			continue
		}
		language := languageInterface.(string)
		if language == "en" {
			uriInterface := location.(map[string]interface{})["Uri"]
			//if  Uri object type is map then we skip, as we dont know how to proceed
			// with processing the document.
			if reflect.ValueOf(uriInterface).Kind() == reflect.Map {
				continue
			}
			if uriInterface != nil {
				uri = uriInterface.(string)
			}
			break
		}
	}
	if uri == "" {
		/*
			h.lock.Lock()
			h.ErrorMessage = "error while Registry file Uri is empty"
			h.StatusMessage = response.InternalError
			h.StatusCode = http.StatusInternalServerError
			h.lock.Unlock()
		*/
		return progress + allotedWork
	}
	req.OID = uri
	h.getRegistryFile(registryName, req)
	// File already exist retrun progress here
	return progress + allotedWork

}

func (h *respHolder) getRegistryFile(registryName string, req getResourceRequest) {
	body, _, getResponse, err := contactPlugin(req, "error while trying to get Registry file: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		h.StatusMessage = getResponse.StatusMessage
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		return
	}

	// Store the file in DB
	err = agmodel.SaveRegistryFile(body, "Registries", registryName+".json")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying to save data: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return
	}
}

func isFileExist(existingFiles []string, substr string) bool {
	fileExist := false

	for _, existingFile := range existingFiles {
		index := strings.Index(existingFile, substr)
		if index != -1 {
			return true
		}
	}
	// Check if the file is present in DB
	_, err := agmodel.GetRegistryFile("Registries", substr)
	if err == nil {
		fileExist = true
	}
	return fileExist
}

func (h *respHolder) getAllRootInfo(taskID string, progress int32, alottedWork int32, req getResourceRequest) int32 {
	resourceName := req.OID
	body, _, getResponse, err := contactPlugin(req, "error while trying to get the"+resourceName+"collection details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		h.StatusMessage = getResponse.StatusMessage
		h.StatusCode = getResponse.StatusCode
		h.MsgArgs = getResponse.MsgArgs
		h.lock.Unlock()
		log.Error(err)
		return progress
	}

	resourceMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(body), &resourceMap)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal " + resourceName + " " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		log.Error("error while trying to unmarshal " + resourceName + ": " + err.Error())
		return progress

	}

	resourceMembers := resourceMap["Members"]
	// Loop through all the resource members collection and discover all of them
	for _, object := range resourceMembers.([]interface{}) {
		estimatedWork := alottedWork / int32(len(resourceMembers.([]interface{})))
		oDataID := object.(map[string]interface{})["@odata.id"].(string)
		req.OID = oDataID
		progress = h.getIndivdualInfo(taskID, progress, estimatedWork, req)
	}
	return progress
}

func (h *respHolder) getSystemInfo(taskID string, progress int32, alottedWork int32, req getResourceRequest) (string, string, int32, error) {
	var computeSystemID, oidKey string
	body, _, getResponse, err := contactPlugin(req, "error while trying to get system collection details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		if strings.Contains(h.ErrorMessage, errors.SystemNotSupportedErrString) {
			h.StatusMessage = response.ActionNotSupported
		} else {
			h.StatusMessage = getResponse.StatusMessage
		}
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		return computeSystemID, oidKey, progress, err
	}

	var computeSystem map[string]interface{}
	err = json.Unmarshal(body, &computeSystem)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal response body: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return computeSystemID, oidKey, progress, err
	}

	oid := computeSystem["@odata.id"].(string)
	computeSystemID = computeSystem["Id"].(string)
	computeSystemUUID := computeSystem["UUID"].(string)
	oidKey = keyFormation(oid, computeSystemID, req.DeviceUUID)
	if !req.UpdateFlag {
		indexList, err := agmodel.GetString("UUID", computeSystemUUID)
		if err != nil {
			log.Error(err.Error())
			h.lock.Lock()
			h.StatusCode = http.StatusInternalServerError
			h.StatusMessage = response.InternalError
			h.lock.Unlock()
			return computeSystemID, oidKey, progress, err
		}
		if len(indexList) > 0 {
			h.lock.Lock()
			h.StatusCode = http.StatusConflict
			h.StatusMessage = response.ResourceAlreadyExists
			h.ErrorMessage = "Resource already exists"
			h.MsgArgs = []interface{}{"ComputerSystem", "ComputerSystem", "ComputerSystem"}
			h.lock.Unlock()
			return computeSystemID, oidKey, progress, fmt.Errorf(h.ErrorMessage)
		}

	}
	updatedResourceData := updateResourceDataWithUUID(string(body), req.DeviceUUID)
	// persist the response with table ComputerSystem  and key as system UUID + Oid Needs relook TODO
	err = agmodel.GenericSave([]byte(updatedResourceData), "ComputerSystem", oidKey)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying to save data: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return computeSystemID, oidKey, progress, err
	}
	h.TraversedLinks[req.OID] = true
	h.SystemURL = append(h.SystemURL, oidKey)
	var retrievalLinks = make(map[string]bool)
	getLinks(computeSystem, retrievalLinks, false)

	removeRetrievalLinks(retrievalLinks, oid, config.Data.AddComputeSkipResources.SystemCollection, h.TraversedLinks)

	req.SystemID = computeSystemID
	req.ParentOID = oid
	for resourceOID, oemFlag := range retrievalLinks {
		estimatedWork := alottedWork / int32(len(retrievalLinks))
		req.OID = resourceOID
		req.OemFlag = oemFlag
		progress = h.getResourceDetails(taskID, progress, estimatedWork, req)
	}
	json.Unmarshal([]byte(updatedResourceData), &computeSystem)
	searchForm := createServerSearchIndex(computeSystem, oidKey, req.DeviceUUID)
	//save the final search form here
	if req.UpdateFlag {
		err = agmodel.UpdateIndex(searchForm, oidKey, computeSystemUUID)
	} else {
		err = agmodel.SaveIndex(searchForm, oidKey, computeSystemUUID)
	}
	if err != nil {
		h.ErrorMessage = "error while trying save index values: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		return computeSystemID, oidKey, progress, err
	}
	return computeSystemID, oidKey, progress, nil
}

// getStorageInfo is used to rediscover storage data from a system
func (h *respHolder) getStorageInfo(progress int32, alottedWork int32, req getResourceRequest) (string, int32, error) {
	body, _, getResponse, err := contactPlugin(req, "error while trying to get system storage collection details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		if strings.Contains(h.ErrorMessage, errors.SystemNotSupportedErrString) {
			h.StatusMessage = response.ActionNotSupported
		} else {
			h.StatusMessage = getResponse.StatusMessage
		}
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		return "", progress, err
	}

	var computeSystem map[string]interface{}
	err = json.Unmarshal(body, &computeSystem)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal response body of system storage: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return "", progress, err
	}

	// Read system data from DB
	systemURI := strings.Replace(req.OID, "/Storage", "", -1)
	systemURI = strings.Replace(systemURI, "/Systems/", "/Systems/"+req.DeviceUUID+":", -1)
	data, dbErr := agmodel.GetResource("ComputerSystem", systemURI)
	if dbErr != nil {
		log.Error("error while getting the systems data" + dbErr.Error())
		return "", progress, err
	}
	// unmarshall the systems data
	var systemData map[string]interface{}
	err = json.Unmarshal([]byte(data), &systemData)
	if err != nil {
		log.Error("Error while unmarshaling system's data" + err.Error())
		return "", progress, err
	}

	oid := computeSystem["@odata.id"].(string)
	computeSystemID := systemData["Id"].(string)
	computeSystemUUID := systemData["UUID"].(string)
	oidKey := keyFormation(oid, computeSystemID, req.DeviceUUID)

	updatedResourceData := updateResourceDataWithUUID(string(body), req.DeviceUUID)
	// persist the response with table Storage
	resourceName := getResourceName(req.OID, true)
	err = agmodel.GenericSave([]byte(updatedResourceData), resourceName, oidKey)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying to save data: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return oidKey, progress, err
	}
	h.TraversedLinks[req.OID] = true
	h.SystemURL = append(h.SystemURL, oidKey)
	var retrievalLinks = make(map[string]bool)
	getLinks(computeSystem, retrievalLinks, false)

	removeRetrievalLinks(retrievalLinks, oid, config.Data.AddComputeSkipResources.SystemCollection, h.TraversedLinks)

	req.SystemID = computeSystemID
	req.ParentOID = oid
	for resourceOID, oemFlag := range retrievalLinks {
		estimatedWork := alottedWork / int32(len(retrievalLinks))
		req.OID = resourceOID
		req.OemFlag = oemFlag
		// Passing taskid as empty string
		progress = h.getResourceDetails("", progress, estimatedWork, req)
	}
	json.Unmarshal([]byte(updatedResourceData), &computeSystem)
	searchForm := createServerSearchIndex(computeSystem, systemURI, req.DeviceUUID)
	//save the final search form here
	if req.UpdateFlag {
		err = agmodel.SaveIndex(searchForm, systemURI, computeSystemUUID)
	}
	if err != nil {
		h.ErrorMessage = "error while trying save index values: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		return oidKey, progress, err
	}
	return oidKey, progress, nil
}

func createServerSearchIndex(computeSystem map[string]interface{}, oidKey, deviceUUID string) map[string]interface{} {
	var searchForm = make(map[string]interface{})

	if val, ok := computeSystem["MemorySummary"]; ok {
		memSum := val.(map[string]interface{})
		searchForm["MemorySummary/TotalSystemMemoryGiB"] = memSum["TotalSystemMemoryGiB"].(float64)
		if _, ok := memSum["TotalSystemPersistentMemoryGiB"]; ok {
			searchForm["MemorySummary/TotalSystemPersistentMemoryGiB"] = memSum["TotalSystemPersistentMemoryGiB"].(float64)
		}
	}
	if _, ok := computeSystem["SystemType"]; ok {
		searchForm["SystemType"] = computeSystem["SystemType"].(string)
	}
	if val, ok := computeSystem["ProcessorSummary"]; ok {
		procSum := val.(map[string]interface{})
		searchForm["ProcessorSummary/Count"] = procSum["Count"].(float64)
		searchForm["ProcessorSummary/sockets"] = procSum["Count"].(float64)
		searchForm["ProcessorSummary/Model"] = procSum["Model"].(string)
	}
	if _, ok := computeSystem["PowerState"]; ok {
		searchForm["PowerState"] = computeSystem["PowerState"].(string)
	}

	// saving the firmware version
	if !strings.Contains(oidKey, "/Storage") {
		if firmwareVersion := getFirmwareVersion(oidKey); firmwareVersion != "" {
			searchForm["FirmwareVersion"] = firmwareVersion
		}
	}

	// saving storage drive quantity/capacity/type
	if val, ok := computeSystem["Storage"]; ok || strings.Contains(oidKey, "/Storage") {
		var storageCollectionOdataID string
		if strings.Contains(oidKey, "/Storage") {
			storageCollectionOdataID = oidKey
		} else {
			storage := val.(map[string]interface{})
			storageCollectionOdataID = storage["@odata.id"].(string)
		}
		storageCollection := agcommon.GetStorageResources(strings.TrimSuffix(storageCollectionOdataID, "/"))
		storageMembers := storageCollection["Members"]
		if storageMembers != nil {
			// Loop through all the storage members collection and discover all of them
			for _, object := range storageMembers.([]interface{}) {
				storageODataID := object.(map[string]interface{})["@odata.id"].(string)
				storageRes := agcommon.GetStorageResources(strings.TrimSuffix(storageODataID, "/"))
				drives := storageRes["Drives"]
				if drives != nil {
					quantity := len(drives.([]interface{}))
					var capacity []float64
					var types []string
					for _, drive := range drives.([]interface{}) {
						driveODataID := drive.(map[string]interface{})["@odata.id"].(string)
						driveRes := agcommon.GetStorageResources(strings.TrimSuffix(driveODataID, "/"))
						capInBytes := driveRes["CapacityBytes"]
						// convert bytes to gb in decimal format
						if capInBytes != nil {
							capInGbs := capInBytes.(float64) / 1000000000
							capacity = append(capacity, capInGbs)
						}
						mediaType := driveRes["MediaType"]
						if mediaType != nil {
							types = append(types, mediaType.(string))
						}
					}

					searchForm["Storage/Drives/Quantity"] = quantity
					searchForm["Storage/Drives/Capacity"] = capacity
					searchForm["Storage/Drives/Type"] = types
				}
			}
		}
	}
	return searchForm
}
func (h *respHolder) getIndivdualInfo(taskID string, progress int32, alottedWork int32, req getResourceRequest) int32 {
	resourceName := getResourceName(req.OID, false)
	body, _, getResponse, err := contactPlugin(req, "error while trying to get "+resourceName+" details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		h.StatusMessage = getResponse.StatusMessage
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		return progress
	}
	var resource map[string]interface{}
	err = json.Unmarshal(body, &resource)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal response body: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return progress
	}
	oid := resource["@odata.id"].(string)
	resourceID := resource["Id"].(string)

	oidKey := keyFormation(oid, resourceID, req.DeviceUUID)

	//replacing the uuid while saving the data
	updatedResourceData := updateResourceDataWithUUID(string(body), req.DeviceUUID)
	// persist the response with table resource and key as system UUID + Oid Needs relook TODO
	err = agmodel.GenericSave([]byte(updatedResourceData), resourceName, oidKey)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying to save data: " + err.Error()
		h.StatusMessage = response.InternalError
		h.StatusCode = http.StatusInternalServerError
		h.lock.Unlock()
		return progress
	}
	h.TraversedLinks[req.OID] = true
	var retrievalLinks = make(map[string]bool)
	getLinks(resource, retrievalLinks, false)

	removeRetrievalLinks(retrievalLinks, oid, config.Data.AddComputeSkipResources.ChassisCollection, h.TraversedLinks)

	req.SystemID = resourceID
	req.ParentOID = oid
	for resourceOID, oemFlag := range retrievalLinks {
		estimatedWork := alottedWork / int32(len(retrievalLinks))
		req.OID = resourceOID
		req.OemFlag = oemFlag
		progress = h.getResourceDetails(taskID, progress, estimatedWork, req)
	}
	return progress
}

func (h *respHolder) getResourceDetails(taskID string, progress int32, alottedWork int32, req getResourceRequest) int32 {
	h.TraversedLinks[req.OID] = true
	body, _, getResponse, err := contactPlugin(req, "error while trying to get the "+req.OID+" details: ")
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = err.Error()
		h.StatusMessage = getResponse.StatusMessage
		h.MsgArgs = getResponse.MsgArgs
		h.StatusCode = getResponse.StatusCode
		h.lock.Unlock()
		return progress
	}
	var resourceData map[string]interface{}
	err = json.Unmarshal(body, &resourceData)
	if err != nil {
		h.lock.Lock()
		h.ErrorMessage = "error while trying unmarshal : " + err.Error()
		h.StatusCode = http.StatusInternalServerError
		h.StatusMessage = response.InternalError
		log.Error(h.ErrorMessage)
		h.lock.Unlock()
		return progress
	}
	oidKey := keyFormation(req.OID, req.SystemID, req.DeviceUUID)
	var memberFlag bool
	if _, ok := resourceData["Members"]; ok {
		memberFlag = true
	}
	resourceName := getResourceName(req.OID, memberFlag)

	//replacing the uuid while saving the data
	updatedResourceData := updateResourceDataWithUUID(string(body), req.DeviceUUID)
	// persist the response with table resourceName and key as system UUID + Oid Needs relook TODO
	err = agmodel.GenericSave([]byte(updatedResourceData), resourceName, oidKey)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return progress
		}
		h.lock.Lock()
		h.ErrorMessage = "error while trying to save data: " + err.Error()
		h.StatusCode = http.StatusInternalServerError
		h.StatusMessage = response.InternalError
		log.Error(h.ErrorMessage)
		h.lock.Unlock()
		return progress
	}
	var retrievalLinks = make(map[string]bool)
	getLinks(resourceData, retrievalLinks, req.OemFlag)
	/* Loop through  Collection members and discover all of them*/
	for oid, oemFlag := range retrievalLinks {
		// skipping the Retrieval if oid mathches the parent oid
		if checkRetrieval(oid, req.OID, h.TraversedLinks) {
			estimatedWork := alottedWork / int32(len(retrievalLinks))
			childReq := req
			childReq.OID = oid
			childReq.ParentOID = req.OID
			childReq.OemFlag = oemFlag
			progress = h.getResourceDetails(taskID, progress, estimatedWork, childReq)
		}
	}

	progress = progress + alottedWork
	var task = fillTaskData(taskID, req.TargetURI, req.TaskRequest, response.RPC{}, common.Running, common.OK, progress, http.MethodPost)
	err = req.UpdateTask(task)

	if err != nil && (err.Error() == common.Cancelling) {
		var task = fillTaskData(taskID, req.TargetURI, req.TaskRequest, response.RPC{}, common.Cancelled, common.OK, progress, http.MethodPost)
		req.UpdateTask(task)

	}
	return progress
}

func getResourceName(oDataID string, memberFlag bool) string {
	str := strings.Split(oDataID, "/")
	if memberFlag {
		return str[len(str)-1] + "Collection"
	}
	if _, err := strconv.Atoi(str[len(str)-2]); err == nil {
		return str[len(str)-1]
	}
	return str[len(str)-2]
}

// getLinks recursively finds and stores all the  @odata.id whcih is present in the request
func getLinks(data map[string]interface{}, retrievalLinks map[string]bool, oemFlag bool) {
	for key, value := range data {
		switch value.(type) {
		// condition to validate the map data
		case map[string]interface{}:
			if strings.EqualFold(key, "Oem") {
				oemFlag = true
			}
			getLinks(value.(map[string]interface{}), retrievalLinks, oemFlag)
		// condition to validate the array data
		case []interface{}:
			memberData := value.([]interface{})
			for _, v := range memberData {
				switch v.(type) {
				case map[string]interface{}:
					if strings.EqualFold(key, "Oem") {
						oemFlag = true
					}
					getLinks(v.(map[string]interface{}), retrievalLinks, oemFlag)
				}
			}
		default:
			// stores value of @odata.id
			if key == "@odata.id" {
				link := strings.TrimSuffix(value.(string), "/")
				retrievalLinks[link] = oemFlag
			}
		}

	}
}

func checkRetrieval(oid, parentoid string, traversedLinks map[string]bool) bool {
	if _, ok := traversedLinks[oid]; ok {
		return false
	}
	//skiping the Retrieval if oid mathches the parent oid
	if strings.EqualFold(parentoid, oid) || strings.EqualFold(parentoid+"/", oid) {
		return false
	}
	//skiping the Retrieval if parent oid contains links in other resource of config
	// TODO : beyond second level Retrieval need to be taken from config it will be implemented in RUCE-1239
	for _, resourceName := range config.Data.AddComputeSkipResources.OtherCollection {
		if strings.Contains(parentoid, resourceName) {
			return false
		}
	}
	return true
}

func removeRetrievalLinks(retrievalLinks map[string]bool, parentoid string, resourceList []string, traversedLinks map[string]bool) {
	for resoureOID := range retrievalLinks {
		// check if oid is already traversed
		if _, ok := traversedLinks[resoureOID]; ok {
			delete(retrievalLinks, resoureOID)
			continue
		}
		// removing the oid if matches parent oid
		if strings.EqualFold(parentoid, resoureOID) || strings.EqualFold(parentoid+"/", resoureOID) {
			delete(retrievalLinks, resoureOID)
			continue
		}
		for i := 0; i < len(resourceList); i++ {
			// removing the oid if it is present list which contains all resoure name  which need to be ignored
			if strings.Contains(resoureOID, resourceList[i]) {
				delete(retrievalLinks, resoureOID)
				continue
			}
		}
	}
	return
}

func callPlugin(req getResourceRequest) (*http.Response, error) {
	var oid string
	for key, value := range getTranslationURL(southBoundURL) {
		oid = strings.Replace(req.OID, key, value, -1)
	}
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + oid
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		return req.ContactClient(reqURL, req.HTTPMethodType, "", oid, req.DeviceInfo, req.LoginCredentials)
	}
	return req.ContactClient(reqURL, req.HTTPMethodType, req.Token, oid, req.DeviceInfo, nil)
}

func updateManagerName(data []byte, pluginID string) []byte {
	var managersMap map[string]interface{}
	json.Unmarshal(data, &managersMap)
	managersMap["Name"] = pluginID
	data, _ = json.Marshal(managersMap)
	return data
}

func getFirmwareVersion(oid string) string {
	// replace the system with the manager
	managerID := strings.Replace(oid, "Systems", "Managers", -1)
	data, dbErr := agmodel.GetResource("Managers", managerID)
	if dbErr != nil {
		log.Error("error while getting the managers data" + dbErr.Error())
		return ""
	}
	// unmarshall the managers data
	var managersData map[string]interface{}
	err := json.Unmarshal([]byte(data), &managersData)
	if err != nil {
		log.Error("Error while unmarshaling  the data" + err.Error())
		return ""
	}
	var firmwareVersion string
	var ok bool
	if firmwareVersion, ok = managersData["FirmwareVersion"].(string); !ok {
		return ""
	}
	return firmwareVersion
}

// CreateDefaultEventSubscription will create default events subscriptions
func CreateDefaultEventSubscription(systemID []string) {
	log.Error("Creation of default subscriptions for " + strings.Join(systemID, ", ") + " are initiated.")
	events := eventsproto.NewEventsService(services.Events, services.Service.Client())
	_, err := events.CreateDefaultEventSubscription(context.TODO(), &eventsproto.DefaultEventSubRequest{
		SystemID:      systemID,
		EventTypes:    []string{"Alert"},
		MessageIDs:    []string{},
		ResourceTypes: []string{},
		Protocol:      "Redfish",
	})
	if err != nil {
		log.Error("error while creating default events: " + err.Error())
		return
	}
}

// PublishEvent will publish default events
func PublishEvent(systemIDs []string, collectionName string) {
	for i := 0; i < len(systemIDs); i++ {
		agmessagebus.Publish(systemIDs[i], "ResourceAdded", collectionName)
	}
}

// PublishPluginStatusOKEvent is for notifying active status of a plugin
// and indicating to resubscribe the EMB of the plugin
func PublishPluginStatusOKEvent(plugin string, msgQueues []string) {
	data := common.SubscribeEMBData{
		PluginID:  plugin,
		EMBQueues: msgQueues,
	}
	if err := agmessagebus.PublishCtrlMsg(common.SubscribeEMB, data); err != nil {
		log.Error("failed to publish resubscribe to " + plugin + " EMB event: " + err.Error())
		return
	}
	log.Info("Published event to resubscribe to " + plugin + " EMB")
}

func getIDsFromURI(uri string) (string, string, error) {
	lastChar := uri[len(uri)-1:]
	if lastChar == "/" {
		uri = uri[:len(uri)-1]
	}
	uriParts := strings.Split(uri, "/")
	ids := strings.Split(uriParts[len(uriParts)-1], ":")
	if len(ids) != 2 {
		return "", "", fmt.Errorf("error: no system id is found in %v", uri)
	}
	return ids[0], ids[1], nil
}

// rollbackInMemory will delete all InMemory data with the resourceURI
// passed. This function is used for rollback the InMemoryDB data
// if any error happens while adding a server
func (e *ExternalInterface) rollbackInMemory(resourceURI string) {
	if resourceURI != "" {
		index := strings.LastIndexAny(resourceURI, "/")
		e.DeleteComputeSystem(index, resourceURI)
	}
}

func updateResourceDataWithUUID(resourceData, uuid string) string {
	//replacing the uuid while saving the data
	//to replace the id of system
	var updatedResourceData = strings.Replace(resourceData, "/redfish/v1/Systems/", "/redfish/v1/Systems/"+uuid+":", -1)
	updatedResourceData = strings.Replace(updatedResourceData, "/redfish/v1/systems/", "/redfish/v1/systems/"+uuid+":", -1)
	// to replace the id in managers
	updatedResourceData = strings.Replace(updatedResourceData, "/redfish/v1/Managers/", "/redfish/v1/Managers/"+uuid+":", -1)
	// to replace id in chassis
	return strings.Replace(updatedResourceData, "/redfish/v1/Chassis/", "/redfish/v1/Chassis/"+uuid+":", -1)

}

// check plugin type is supported
func isPluginTypeSupported(pluginType string) bool {
	for _, pType := range config.Data.SupportedPluginTypes {
		if pType == pluginType {
			return true
		}
	}
	return false
}

func getTranslationURL(translationURL string) map[string]string {
	common.MuxLock.Lock()
	defer common.MuxLock.Unlock()
	if translationURL == southBoundURL {
		return config.Data.URLTranslation.SouthBoundURL
	}
	return config.Data.URLTranslation.NorthBoundURL
}

func checkStatus(pluginContactRequest getResourceRequest, req AddResourceRequest, cmVariants connectionMethodVariants, taskInfo *common.TaskUpdateInfo) (response.RPC, int32, []string) {

	var queueList = make([]string, 0)
	ipData := strings.Split(req.ManagerAddress, ":")
	var ip, port string
	ip = ipData[0]
	if len(ipData) > 1 {
		port = ipData[1]
	}
	var plugin = agmodel.Plugin{
		IP:                ip,
		Port:              port,
		Username:          req.UserName,
		Password:          []byte(req.Password),
		ID:                cmVariants.PluginID,
		PluginType:        cmVariants.PluginType,
		PreferredAuthType: cmVariants.PreferredAuthType,
	}
	pluginContactRequest.Plugin = plugin
	pluginContactRequest.StatusPoll = true
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		pluginContactRequest.HTTPMethodType = http.MethodPost
		pluginContactRequest.DeviceInfo = map[string]interface{}{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}
		pluginContactRequest.OID = "/ODIM/v1/Sessions"
		_, token, getResponse, err := contactPlugin(pluginContactRequest, "error while creating the session: ")
		if err != nil {
			errMsg := err.Error()
			log.Error(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), getResponse.StatusCode, queueList
		}
		pluginContactRequest.Token = token
	} else {
		pluginContactRequest.LoginCredentials = map[string]string{
			"UserName": plugin.Username,
			"Password": string(plugin.Password),
		}
	}

	// Verfiying the plugin Status
	pluginContactRequest.HTTPMethodType = http.MethodGet
	pluginContactRequest.OID = "/ODIM/v1/Status"
	body, _, getResponse, err := contactPlugin(pluginContactRequest, "error while getting the details "+pluginContactRequest.OID+": ")
	if err != nil {
		errMsg := err.Error()
		log.Error(errMsg)
		if getResponse.StatusCode == http.StatusNotFound {
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil), getResponse.StatusCode, queueList
		}
		return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, taskInfo), getResponse.StatusCode, queueList
	}
	// extracting the EMB Type and EMB Queue name
	var statusResponse common.StatusResponse
	err = json.Unmarshal(body, &statusResponse)
	if err != nil {
		errMsg := err.Error()
		log.Error(errMsg)
		getResponse.StatusCode = http.StatusInternalServerError
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, taskInfo), getResponse.StatusCode, queueList
	}

	// check the firmware version of plugin is matched with connection method variant version
	if statusResponse.Version != cmVariants.FirmwareVersion {
		errMsg := fmt.Sprintf("Provided firmware version %s does not match supported firmware version %s of the plugin %s", cmVariants.FirmwareVersion, statusResponse.Version, cmVariants.PluginID)
		log.Error(errMsg)
		getResponse.StatusCode = http.StatusBadRequest
		return common.GeneralError(http.StatusBadRequest, response.PropertyValueNotInList, errMsg, []interface{}{"FirmwareVersion", statusResponse.Version}, taskInfo), getResponse.StatusCode, queueList
	}
	if statusResponse.EventMessageBus != nil {
		for i := 0; i < len(statusResponse.EventMessageBus.EmbQueue); i++ {
			queueList = append(queueList, statusResponse.EventMessageBus.EmbQueue[i].QueueName)
		}
	}
	return response.RPC{}, getResponse.StatusCode, queueList
}

func getConnectionMethodVariants(connectionMethodVariant string) connectionMethodVariants {
	// Split the connectionmethodvariant and get the PluginType, PreferredAuthType, PluginID and FirmwareVersion.
	// Example: Compute:BasicAuth:GRF_v1.0.0
	cm := strings.Split(connectionMethodVariant, ":")
	firmwareVersion := strings.Split(cm[2], "_")
	return connectionMethodVariants{
		PluginType:        cm[0],
		PreferredAuthType: cm[1],
		PluginID:          cm[2],
		FirmwareVersion:   firmwareVersion[1],
	}
}
