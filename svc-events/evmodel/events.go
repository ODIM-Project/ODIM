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

// Package evmodel have the struct models and DB functionalities
package evmodel

import (
	"encoding/json"
	"fmt"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

const (
	// EventFormatType is set to Event (MetricReport is not supporting now
	EventFormatType = "Event"

	// SubscriptionType is set to RedfishEvent (make it as array of SubscriptionType
	SubscriptionType = "RedfishEvent"

	// Context is set to default if its empty
	Context = "Default"

	// SubscriptionName is set to default name incase if its empty
	SubscriptionName = "Event Subscription"

	// SubscriptionIndex is a index name which required for indexing of event subscriptions
	SubscriptionIndex = common.SubscriptionIndex

	// DeviceSubscriptionIndex is a index name which required for indexing
	// subscription of device
	DeviceSubscriptionIndex = common.DeviceSubscriptionIndex

	// UndeliveredEvents holds table for UndeliveredEvent
	UndeliveredEvents = "UndeliveredEvents"

	// ReadInProgress holds table for ReadInProgress
	ReadInProgress = "ReadInProgress"
	// DeliveryRetryPolicy is set to default value incase if its empty
	DeliveryRetryPolicy = "RetryForever"

	// AggregateSubscriptionIndex is a index name which required for indexing
	// subscription of device
	AggregateSubscriptionIndex = common.AggregateSubscriptionIndex
)

var (
	//GetDbConnection alias for common.GetDBConnection
	GetDbConnection = common.GetDBConnection
)

// SubscriptionResource is a model to store the subscription details
type SubscriptionResource struct {
	EventDestination *dmtf.EventDestination `json:"EventDestination"`
	EventHostIP      string                 `json:"EventHostIP,omitempty"`
	Hosts            []string               `json:"Hosts"`
	SubscriptionID   string                 `json:"SubscriptionID"`
	UserName         string                 `json:"UserName"`
	Location         string                 `json:"location,omitempty"`
}

// Fabric is the model for fabrics information
type Fabric struct {
	FabricUUID string
	PluginID   string
}

// EventPost is the model for post data to client
type EventPost struct {
	Destination        string
	EventID            string
	UndeliveredEventID string
	Message            []byte
}

// Aggregate is the model for Aggregate information
type Aggregate struct {
	Elements []dmtf.Link `json:"Elements"`
}

// GetResource fetches a resource from database using table and key
func GetResource(Table, key string) (string, *errors.Error) {
	conn, err := GetDbConnection(common.InMemory)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), err)
	}
	resourceData, err := conn.Read(Table, key)
	if err != nil {
		return "", errors.PackError(err.ErrNo(), "error while trying to get resource details: ", err.Error())
	}
	var resource string
	if errs := json.Unmarshal([]byte(resourceData), &resource); errs != nil {
		return "", errors.PackError(errors.UndefinedErrorType, errs)
	}
	return resource, nil
}

// GetTarget fetches the System(Target Device Credentials) table details
func GetTarget(deviceUUID string) (*common.Target, error) {
	var target common.Target
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}

	data, err := conn.Read("System", deviceUUID)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get compute details: %v", err.Error())
	}
	if errs := json.Unmarshal([]byte(data), &target); errs != nil {
		return nil, errs
	}
	return &target, nil

}

// GetPluginData will fetch plugin details
func GetPluginData(pluginID string) (*common.Plugin, *errors.Error) {
	var plugin common.Plugin

	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}

	plugindata, err := conn.Read("Plugin", pluginID)
	if err != nil {
		return nil, errors.PackError(err.ErrNo(), "error while trying to fetch plugin data: ", err.Error())
	}

	if err := json.Unmarshal([]byte(plugindata), &plugin); err != nil {
		return nil, errors.PackError(errors.JSONUnmarshalFailed, err)
	}

	bytepw, errs := common.DecryptWithPrivateKey([]byte(plugin.Password))
	if errs != nil {
		return nil, errors.PackError(errors.DecryptionFailed, "error: "+pluginID+
			" plugin password decryption failed: "+errs.Error())
	}
	plugin.Password = bytepw

	return &plugin, nil
}

// GetAllPlugins gets all the Plugin from the db
func GetAllPlugins() ([]common.Plugin, *errors.Error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("Plugin")
	if err != nil {
		return nil, err
	}
	var plugins []common.Plugin
	for _, key := range keys {
		var plugin common.Plugin
		plugindata, err := conn.Read("Plugin", key)
		if err != nil {
			return nil, errors.PackError(err.ErrNo(), "error while trying to fetch plugin data: ", err.Error())
		}

		if err := json.Unmarshal([]byte(plugindata), &plugin); err != nil {
			return nil, errors.PackError(errors.JSONUnmarshalFailed, err)
		}

		bytepw, errs := common.DecryptWithPrivateKey([]byte(plugin.Password))
		if errs != nil {
			return nil, errors.PackError(errors.DecryptionFailed, "error: "+plugin.ID+
				" plugin password decryption failed: "+errs.Error())
		}
		plugin.Password = bytepw

		plugins = append(plugins, plugin)

	}
	return plugins, nil
}

// GetAllKeysFromTable return all matching data give table name
func GetAllKeysFromTable(table string) ([]string, error) {
	conn, err := GetDbConnection(common.InMemory)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails(table)
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table - %v: %v", table, err.Error())
	}
	return keysArray, nil
}

// GetAllSystems retrieves all the compute systems in odimra
func GetAllSystems() ([]string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails("System")
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table - System: %v", err)
	}
	return keysArray, nil
}

// GetSingleSystem retrieves specific compute system in odimra based on the ID
func GetSingleSystem(id string) (string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return "", errors.PackError(errors.UndefinedErrorType, err)
	}

	data, rerr := conn.Read("System", id)
	if rerr != nil {
		return "", errors.PackError(rerr.ErrNo(), "error while trying to get compute details: ", rerr.Error())
	}
	return data, nil
}

// GetFabricData  will fetch fabric details
func GetFabricData(fabricID string) (Fabric, error) {
	var fabric Fabric

	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return fabric, err
	}

	fabricdata, err := conn.Read("Fabric", fabricID)
	if err != nil {
		return fabric, fmt.Errorf("error while trying to get user: %v", err.Error())
	}

	if errs := json.Unmarshal([]byte(fabricdata), &fabric); errs != nil {
		return fabric, errs
	}

	return fabric, nil
}

// GetAggregateData  will fetch aggregate details
func GetAggregateData(aggregateKey string) (Aggregate, error) {
	var aggregate Aggregate
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return aggregate, err
	}
	aggregateData, err := conn.Read("Aggregate", aggregateKey)
	if err != nil {
		return aggregate, fmt.Errorf("error while trying to get user: %v", err.Error())
	}
	if errs := json.Unmarshal([]byte(aggregateData), &aggregate); errs != nil {
		return aggregate, errs
	}

	return aggregate, nil
}

// GetAllFabrics return all Fabrics
func GetAllFabrics() ([]string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails("Fabric")
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table -Fabric: %v", err.Error())
	}
	for i := 0; i < len(keysArray); i++ {
		keysArray[i] = "/redfish/v1/Fabrics/" + keysArray[i]
	}
	return keysArray, nil
}

// GetDeviceSubscriptions is to get subscription details of device
func GetDeviceSubscriptions(hostIP string) (*common.DeviceSubscription, error) {

	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	devSubscription, gerr := conn.GetDeviceSubscription(DeviceSubscriptionIndex, hostIP+"*")
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get subscription of device %v", gerr.Error())
	}
	devSub := strings.Split(devSubscription[0], "||")
	var deviceSubscription = &common.DeviceSubscription{
		EventHostIP:     devSub[0],
		Location:        devSub[1],
		OriginResources: getSliceFromString(devSub[2]),
	}

	return deviceSubscription, nil
}

// UpdateDeviceSubscriptionLocation is to update subscription details of device
func UpdateDeviceSubscriptionLocation(devSubscription common.DeviceSubscription) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	updateErr := conn.UpdateDeviceSubscription(DeviceSubscriptionIndex, devSubscription.EventHostIP,
		devSubscription.Location, devSubscription.OriginResources)
	if updateErr != nil {
		return fmt.Errorf("error while trying to update subscription of device %v", updateErr.Error())
	}
	return nil
}

// SaveDeviceSubscription is to save subscription details of device
func SaveDeviceSubscription(devSubscription common.DeviceSubscription) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	cerr := conn.CreateDeviceSubscriptionIndex(DeviceSubscriptionIndex, devSubscription.EventHostIP,
		devSubscription.Location, devSubscription.OriginResources)
	if cerr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cerr.Error())
	}
	return nil
}

// DeleteDeviceSubscription is to delete subscription details of device
func DeleteDeviceSubscription(hostIP string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	derr := conn.DeleteDeviceSubscription(DeviceSubscriptionIndex, hostIP)
	if derr != nil {
		return fmt.Errorf("error while trying to delete subscription of device %v", derr.Error())
	}
	return nil
}

// getSliceFromString is to convert the string to array
func getSliceFromString(sliceString string) []string {
	// EX : array stored in db in string("[alert statuschange]")
	// to convert into an array removing "[" ,"]" and splitting
	slice := strings.Replace(sliceString, "[", "", -1)
	slice = strings.Replace(slice, "]", "", -1)
	if len(slice) < 1 {
		return []string{}
	}
	return strings.Split(slice, " ")
}

// SaveEventSubscription is to save event subscription details in db
func SaveEventSubscription(evtSubscription SubscriptionResource) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	subscription, marshalErr := json.Marshal(evtSubscription)
	if marshalErr != nil {
		return fmt.Errorf("error while trying marshall event subscriptions %v", marshalErr.Error())
	}
	createErr := conn.CreateEvtSubscriptionIndex(SubscriptionIndex, string(subscription))
	if createErr != nil {
		return fmt.Errorf("error while trying to save event subscriptions %v", createErr.Error())
	}
	return nil
}

// GetEvtSubscriptions is to get event subscription details
func GetEvtSubscriptions(searchKey string) ([]SubscriptionResource, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	evtSub, gerr := conn.GetEvtSubscriptions(SubscriptionIndex, "*"+searchKey+"*")
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get subscription of device %v", gerr.Error())
	}
	var eventSubscriptions []SubscriptionResource
	for _, value := range evtSub {
		var eventSub SubscriptionResource
		if err := json.Unmarshal([]byte(value), &eventSub); err != nil {
			return nil, fmt.Errorf("error while unmarshalling event subscriptions: %v", err.Error())
		}
		eventSubscriptions = append(eventSubscriptions, eventSub)
	}

	return eventSubscriptions, nil
}

// DeleteEvtSubscription is to delete event subscription details
func DeleteEvtSubscription(key string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	derr := conn.DeleteEvtSubscriptions(SubscriptionIndex, "*"+key+"*")
	if derr != nil {
		return fmt.Errorf("error while trying to delete subscription of device %v", derr.Error())
	}
	return nil
}

// UpdateEventSubscription is to update event subscription details
func UpdateEventSubscription(evtSubscription SubscriptionResource) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	subscription, merr := json.Marshal(evtSubscription)
	if merr != nil {
		return fmt.Errorf("error while trying marshall event subscriptions %v", merr.Error())
	}
	uerr := conn.UpdateEvtSubscriptions(SubscriptionIndex, "*"+evtSubscription.SubscriptionID+"*", string(subscription))
	if uerr != nil {
		return fmt.Errorf("error while trying to update subscription of device %v", uerr.Error())
	}
	return nil
}

// GetAllMatchingDetails accepts the table name ,pattern and DB type and return all the keys which matches the pattern
func GetAllMatchingDetails(table, pattern string, dbtype common.DbType) ([]string, *errors.Error) {
	conn, err := GetDbConnection(dbtype)
	if err != nil {
		return []string{}, err
	}
	return conn.GetAllMatchingDetails(table, pattern)
}

// SaveUndeliveredEvents accepts the undelivered event and destination with unique eventId and saves it
func SaveUndeliveredEvents(key string, event []byte) error {
	connPool, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.AddResourceData(UndeliveredEvents, key, string(event)); err != nil {
		return fmt.Errorf("error while trying to add Undelivered Events to DB: %v", err.Error())
	}
	return nil
}

// GetUndeliveredEvents read the undelivered events for the destination
func GetUndeliveredEvents(destination string) (string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return "", fmt.Errorf("error: while trying to create connection with DB: %v", err.Error())
	}

	eventData, err := conn.GetKeyValue(destination)
	if err != nil {
		return "", fmt.Errorf("error: while trying to fetch details: %v", err.Error())
	}

	return eventData, nil
}

// DeleteUndeliveredEvents deletes the undelivered events for the destination
func DeleteUndeliveredEvents(destination string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error: while trying to create connection with DB: %v", err.Error())
	}
	if err := conn.DeleteKey(destination); err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}

// SetUndeliveredEventsFlag will set the flag to maintain one instance already picked up
// the undelivered events for the destination
func SetUndeliveredEventsFlag(destination string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error: while trying to create connection with DB: %v", err.Error())
	}
	if err = conn.AddResourceData(ReadInProgress, destination, "true"); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", ReadInProgress, err.Error())
	}
	_, err = conn.Read(ReadInProgress, destination)
	if err != nil {
		return err
	}
	return nil
}

// GetUndeliveredEventsFlag will get the flag to maintain one instance already picked up
// the undelivered events for the destination
func GetUndeliveredEventsFlag(destination string) (bool, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return false, fmt.Errorf("error: while trying to create connection with DB: %v", err.Error())
	}
	_, err = conn.Read(ReadInProgress, destination)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteUndeliveredEventsFlag deletes the PickUpUndeliveredEventsFlag key from the DB, return error if any
func DeleteUndeliveredEventsFlag(destination string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error: while trying to create connection with DB: %v", err.Error())
	}
	if err := conn.Delete(ReadInProgress, destination); err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}

// SaveAggregateSubscription is to save subscription details of device
func SaveAggregateSubscription(aggregateID string, hostIP []string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	cerr := conn.CreateAggregateHostIndex(AggregateSubscriptionIndex, aggregateID, hostIP)
	if cerr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cerr.Error())
	}
	return nil
}

// UpdateAggregateHosts is to update aggregate hosts details of device
func UpdateAggregateHosts(aggregateID string, hostIP []string) error {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return err
	}
	cerr := conn.UpdateAggregateHosts(AggregateSubscriptionIndex, aggregateID, hostIP)
	if cerr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cerr.Error())
	}
	return nil
}

// GetAggregateHosts is to get subscription details of device
func GetAggregateHosts(aggregateID string) ([]string, error) {

	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	aggregateList, gerr := conn.GetAggregateHosts(AggregateSubscriptionIndex, aggregateID+"[^0-9]*")
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get aggregate host of device %v", gerr.Error())
	}
	devSub := strings.Split(aggregateList[0], "||")
	hostsIP := getSliceFromString(devSub[1])
	return hostsIP, nil
}

// GetAggregateList  will fetch aggregate list
func GetAggregateList(hostIP string) ([]string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	aggregateList, gerr := conn.GetAggregateHosts(AggregateSubscriptionIndex, "*"+hostIP+"*")
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get aggregate host list of device %v", gerr.Error())
	}
	aggregates := []string{}
	for _, v := range aggregateList {
		devSub := strings.Split(v, "||")
		if devSub[0] == "0" {
			continue
		}
		aggregates = append(aggregates, devSub[0])
	}
	return aggregates, nil
}

// GetAggregate fetches the aggregate info for the given aggregateURI
func GetAggregate(aggregateURI string) (Aggregate, *errors.Error) {
	var aggregate Aggregate
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return aggregate, err
	}
	const table string = "Aggregate"
	data, err := conn.Read(table, aggregateURI)
	if err != nil {
		return aggregate, errors.PackError(err.ErrNo(), "error: while trying to fetch connection method data: ", err.Error())
	}
	if err := json.Unmarshal([]byte(data), &aggregate); err != nil {
		return aggregate, errors.PackError(errors.JSONUnmarshalFailed, err)
	}
	return aggregate, nil
}

// GetAllAggregates return all aggregate url added in DB
func GetAllAggregates() ([]string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keysArray, err := conn.GetAllDetails("Aggregate")
	if err != nil {
		return nil, fmt.Errorf("error while trying to get all keys from table - %v: %v", "Aggregate", err.Error())
	}
	return keysArray, nil
}

// GetAllDeviceSubscriptions is to get subscription details of device
func GetAllDeviceSubscriptions() ([]string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	devSubscription, gerr := conn.GetAllDataByIndex(DeviceSubscriptionIndex)
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get subscription of device %v", gerr.Error())
	}
	return devSubscription, nil
}

// GetSliceFromString is to convert the string to array
func GetSliceFromString(sliceString string) []string {
	// EX : array stored in db in string("[alert statusChange]")
	// to convert into an array removing "[" ,"]" and splitting
	r := strings.NewReplacer(
		"[", "",
		"]", "",
	)
	return strings.Split(r.Replace(sliceString), " ")
}

// GetAllEvtSubscriptions is to get all event subscription details
func GetAllEvtSubscriptions() ([]string, error) {
	conn, err := GetDbConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	evtSub, gerr := conn.GetAllDataByIndex(SubscriptionIndex)
	if gerr != nil {
		return nil, fmt.Errorf("error while trying to get subscription of device %v", gerr.Error())
	}
	return evtSub, nil
}

// GetUndeliveredEventsKeyList accepts the table name ,pattern ,cursor value
// and DB type and return all the keys which matches the pattern
func GetUndeliveredEventsKeyList(table, pattern string, dbType common.DbType, nextCursor int) ([]string, int, *errors.Error) {
	conn, err := GetDbConnection(dbType)
	if err != nil {
		return []string{}, 0, err
	}
	return conn.GetAllKeysFromDb(table, pattern, nextCursor)
}
