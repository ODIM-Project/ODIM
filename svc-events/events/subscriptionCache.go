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
// and corresponding unit test cases
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"

	redis "github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

var (
	subscriptionsCache                    map[string]dmtf.EventDestination
	systemToSubscriptionsMap              map[string]map[string]bool
	aggregateIDToSubscriptionsMap         map[string]map[string]bool
	collectionToSubscriptionsMap          map[string]map[string]bool
	emptyOriginResourceToSubscriptionsMap map[string]bool
	systemIDToAggregateIdsMap             map[string]map[string]bool
	eventSourceToManagerIDMap             map[string]string
	managerIDToSystemIDsMap               map[string][]string
	managerIDToChassisIDsMap              map[string][]string

	eventConsumerActionID   = "218"
	eventConsumerActionName = "EventConsumer"
	logging                 *logrus.Entry
)

// temporary variable hold data till reading process happening
var (
	systemToSubscriptionsMapTemp      map[string]map[string]bool
	aggregateIDToSubscriptionsMapTemp map[string]map[string]bool
	collectionToSubscriptionsMapTemp  map[string]map[string]bool
	reAttemptInQueue                  = make(map[string]int)
	reattemptLock                     sync.Mutex
)

// eventForwardingChanel channel is used for communicate
// between event forwarding go routine pools
var eventForwardingChanel = make(chan evmodel.EventPost)

// saveEventChanel channel is used for communicate
// between event forwarding go routine pools
var saveEventChanel = make(chan evmodel.EventPost)

// LoadSubscriptionData method calls whenever service is started
// Here we load Subscription, DeviceSubscription, AggregateToHost
// table data into cache memory
func (e *ExternalInterfaces) LoadSubscriptionData(ctx context.Context) error {
	ctx = context.WithValue(ctx, common.ActionName, eventConsumerActionName)
	ctx = context.WithValue(ctx, common.ActionID, eventConsumerActionID)
	transactionID := uuid.NewV4().String()
	ctx = context.WithValue(ctx, common.TransactionID, transactionID)
	logging = l.LogWithFields(ctx)

	logging.Debug("Event cache is initialized")
	err := getAllSubscriptions(ctx)
	if err != nil {
		return err
	}
	err = getAllAggregates(ctx)
	if err != nil {
		return err
	}
	err = getAllDeviceSubscriptions(ctx)
	if err != nil {
		return err
	}
	threadID := 1
	ctx = context.WithValue(ctx, common.ThreadID, strconv.Itoa(threadID))
	go initializeDbObserver(ctx)
	go e.forwardUndeliveredEventToClient(ctx)
	// create event forwarding worker pool
	for i := 0; i < config.Data.EventForwardingWorkerPoolCount; i++ {
		go e.runEventForwardingWorkers()
	}
	// init event save undelivered worker pool
	for i := 0; i < config.Data.EventSaveWorkerPoolCount; i++ {
		go e.saveEventWorkers()
	}
	return nil
}

// getAllSubscriptions this method read data from Subscription table and
// load in corresponding cache
func getAllSubscriptions(ctx context.Context) error {
	subscriptions, err := evmodel.GetAllEvtSubscriptions()
	if err != nil {
		logging.Error("Error while reading all subscription data ", err)
		return err
	}
	systemToSubscriptionsMapTemp = make(map[string]map[string]bool)
	aggregateIDToSubscriptionsMapTemp = make(map[string]map[string]bool)
	collectionToSubscriptionsMapTemp = make(map[string]map[string]bool)
	emptyOriginResourceToSubscriptionsMapTemp := make(map[string]bool)
	subscriptionsCache = make(map[string]dmtf.EventDestination, len(subscriptions))
	for _, subscription := range subscriptions {
		var sub evmodel.SubscriptionResource
		err = json.Unmarshal([]byte(subscription), &sub)
		if err != nil {
			continue
		}
		subCache := sub.EventDestination
		subCache.ID = sub.SubscriptionID
		subscriptionsCache[subCache.ID] = *subCache
		if len(sub.EventDestination.OriginResources) == 0 && sub.SubscriptionID != evcommon.DefaultSubscriptionID {
			emptyOriginResourceToSubscriptionsMapTemp[sub.SubscriptionID] = true
		} else {
			loadSubscriptionCacheData(sub.SubscriptionID, sub.Hosts)
		}
	}
	emptyOriginResourceToSubscriptionsMap = emptyOriginResourceToSubscriptionsMapTemp
	systemToSubscriptionsMap = systemToSubscriptionsMapTemp
	aggregateIDToSubscriptionsMap = aggregateIDToSubscriptionsMapTemp
	collectionToSubscriptionsMap = collectionToSubscriptionsMapTemp

	logging.Debug("Subscriptions cache updated ")
	return nil
}

// getAllDeviceSubscriptions method fetch data from DeviceSubscription table
func getAllDeviceSubscriptions(ctx context.Context) error {
	deviceSubscriptionList, err := evmodel.GetAllDeviceSubscriptions()
	if err != nil {
		l.LogWithFields(ctx).Error("error while reading all aggregate data ", err)
		return err
	}
	eventSourceToManagerIDMapTemp := make(map[string]string, len(deviceSubscriptionList))
	for _, device := range deviceSubscriptionList {
		devSub := strings.Split(device, "||")
		updateCatchDeviceSubscriptionData(devSub[0], evmodel.GetSliceFromString(devSub[2]), eventSourceToManagerIDMapTemp)
	}
	eventSourceToManagerIDMap = eventSourceToManagerIDMapTemp
	l.LogWithFields(ctx).Debug("DeviceSubscription cache updated ")
	return nil
}

// updateCatchDeviceSubscriptionData update eventSourceToManagerMap for each key with their system IDs
func updateCatchDeviceSubscriptionData(key string, originResources []string, cacheMap map[string]string) {
	systemID := originResources[0][strings.LastIndexByte(originResources[0], '/')+1:]
	cacheMap[key] = systemID
}

// loadSubscriptionCacheData update collectionToSubscriptionsMap,
// systemToSubscriptionsMap , aggregateToSystemIdsMap again subscription details
func loadSubscriptionCacheData(id string, hosts []string) {
	for _, host := range hosts {
		addSubscriptionCache(host, id)
	}
}

// addSubscriptionCache add subscription in corresponding cache based on key type
// collectionToSubscriptionsMap, aggregateIdToSubscriptionsMap, systemToSubscriptionsMap
func addSubscriptionCache(key string, subscriptionID string) {
	if strings.Contains(key, "Collection") {
		updateCacheMaps(key, subscriptionID, collectionToSubscriptionsMapTemp)
		return
	}
	_, err := uuid.FromString(key)
	if err == nil {
		updateCacheMaps(key, subscriptionID, aggregateIDToSubscriptionsMapTemp)
		return
	}
	updateCacheMaps(key, subscriptionID, systemToSubscriptionsMapTemp)

}

// getAllAggregates method will read all aggregate from db and
// update systemIdToAggregateIdsMap to corresponding member in aggregate
func getAllAggregates(ctx context.Context) error {
	systemIDToAggregateIdsMapTemp := make(map[string]map[string]bool)
	aggregateUrls, err := evmodel.GetAllAggregates()
	if err != nil {
		logging.Debug("error occurred while getting aggregate list ", err)
		return err
	}
	for _, aggregateURL := range aggregateUrls {
		aggregate, err := evmodel.GetAggregate(aggregateURL)
		if err != nil {
			continue
		}
		aggregateID := aggregateURL[strings.LastIndexByte(aggregateURL, '/')+1:]
		addSystemIDToAggregateCache(aggregateID, aggregate, systemIDToAggregateIdsMapTemp)
	}
	systemIDToAggregateIdsMap = systemIDToAggregateIdsMapTemp
	logging.Debug("AggregateToHost cache updated ")
	return nil
}

// addSystemIDToAggregateCache update cache for each aggregate member
func addSystemIDToAggregateCache(aggregateID string, aggregate evmodel.Aggregate, cacheMap map[string]map[string]bool) {
	for _, ids := range aggregate.Elements {
		ids.Oid = ids.Oid[strings.LastIndexByte(strings.TrimSuffix(ids.Oid, "/"), '/')+1:]
		updateCacheMaps(ids.Oid, aggregateID, cacheMap)
	}
}

// getSourceID function return system id corresponding host, if not found then return host
func getSourceID(host string) (string, error) {
	data, isExists := eventSourceToManagerIDMap[host]
	if !isExists {
		if strings.Contains(host, "Collection") {
			return host, nil
		}
		return "", fmt.Errorf("invalid source")
	}
	return data, nil
}

// updateCacheMaps update map value corresponding key
func updateCacheMaps(key, value string, cacheData map[string]map[string]bool) {
	elements, isExists := cacheData[key]
	if isExists {
		elements[value] = true
		cacheData[key] = elements
	} else {
		cacheData[key] = map[string]bool{value: true}
	}
}

// getSubscriptions return list of subscription from cache corresponding to originOfCondition
func getSubscriptions(originOfCondition, systemID, hostIP string) (subs []dmtf.EventDestination) {
	subs = append(subs, getSystemSubscriptionList(hostIP)...)
	subs = append(subs, getAggregateSubscriptionList(systemID)...)
	subs = append(subs, getEmptyOriginResourceSubscriptionList()...)
	subs = append(subs, getCollectionSubscriptionList(originOfCondition, hostIP)...)
	return
}

// getSystemSubscriptionList return list of subscription corresponding to host
func getSystemSubscriptionList(hostIP string) (subs []dmtf.EventDestination) {
	systemSubscription, isExists := systemToSubscriptionsMap[hostIP]
	if isExists {
		for subID := range systemSubscription {
			sub, isValidSubID := getSubscriptionDetails(subID)
			if isValidSubID {
				subs = append(subs, sub)
			}

		}
	}
	return
}

// getAggregateSubscriptionList return list of subscription corresponding to system
// is members of different aggregate
func getAggregateSubscriptionList(systemID string) (subs []dmtf.EventDestination) {
	aggregateList, isExists := systemIDToAggregateIdsMap[systemID]
	if isExists {
		for aggregateID := range aggregateList {
			subscriptions, isValidAggregateID := aggregateIDToSubscriptionsMap[aggregateID]
			if isValidAggregateID {
				for subID := range subscriptions {
					sub, isValidSubID := getSubscriptionDetails(subID)
					sub.OriginResources = append(sub.OriginResources, model.Link{Oid: "/redfish/v1/Systems/" + systemID})
					if isValidSubID {
						subs = append(subs, sub)
					}
				}
			}
		}
	}
	return
}

// getCollectionSubscriptionList return list of subscription against
// originOfCondition type
func getCollectionSubscriptionList(originOfCondition, hostIP string) (subs []dmtf.EventDestination) {
	collectionsKey := getCollectionKey(originOfCondition, hostIP)
	collectionSubscription, isExists := collectionToSubscriptionsMap[collectionsKey]
	if isExists {
		for subID := range collectionSubscription {
			sub, isValidSubID := getSubscriptionDetails(subID)
			if isValidSubID {
				subs = append(subs, sub)
			}
		}
	}
	return
}

// getEmptyOriginResourceSubscriptionList return list of subscription
// whose originResources is empty
func getEmptyOriginResourceSubscriptionList() (subs []dmtf.EventDestination) {
	for subID := range emptyOriginResourceToSubscriptionsMap {
		sub, isValidSubID := getSubscriptionDetails(subID)
		if isValidSubID {
			subs = append(subs, sub)
		}
	}
	return
}

// getSubscriptionDetails this method return subscription details corresponding subscription Id
func getSubscriptionDetails(subscriptionID string) (sub dmtf.EventDestination, status bool) {
	if sub, isExists := subscriptionsCache[subscriptionID]; isExists {
		return sub, true
	}
	return dmtf.EventDestination{}, false
}

// getCollectionKey return collection key corresponding originOfCondition uri
func getCollectionKey(oid, host string) (key string) {
	if strings.Contains(oid, "Systems") && host != "SystemsCollection" {
		key = "SystemsCollection"
	} else if strings.Contains(oid, "Chassis") && host != "ChassisCollection" {
		key = "ChassisCollection"
	} else if strings.Contains(oid, "Managers") && host != "ManagerCollection" {
		key = "ManagerCollection"
	} else if strings.Contains(oid, "Fabrics") && host != "FabricsCollection" {
		key = "FabricsCollection"
	}
	return
}

// initializeDbObserver function subscribe redis keySpace notifier
// function notify by channel if any update happened subscribed key
func initializeDbObserver(ctx context.Context) {
START:
	logging.Info("Initializing observer ")
	conn, errDbConn := common.GetDBConnection(common.OnDisk)
	if errDbConn != nil {
		l.Log.Error("error while getDbConnection  ", errDbConn)
		goto START
	}
	err := conn.EnableKeySpaceNotifier(evcommon.RedisNotifierType, evcommon.RedisNotifierFilterKey)
	if err != nil {
		l.LogWithFields(ctx).Error("error occurred configuring key event ", err)
		time.Sleep(time.Second * 1)
		goto START
	}
	psc := redis.PubSub(*conn.RedisClient.Subscribe())
	psc.Subscribe(evcommon.AggregateToHostChannelKey, evcommon.DeviceSubscriptionChannelKey,
		evcommon.SubscriptionChannelKey)

	for {
		data, _ := psc.Receive()
		switch v := data.(type) {

		case redis.Message:
			switch string(v.Pattern) {
			case evcommon.DeviceSubscriptionChannelKey:
				err := getAllDeviceSubscriptions(ctx)
				if err != nil {
					l.LogWithFields(ctx).Error(err)
				}
			case evcommon.SubscriptionChannelKey:
				err := getAllSubscriptions(ctx)
				if err != nil {
					l.LogWithFields(ctx).Error(err)
				}
			case evcommon.AggregateToHostChannelKey:
				err := getAllAggregates(ctx)
				if err != nil {
					l.LogWithFields(ctx).Error(err)
				}
			}
		case error:
			logging.Error("Error occurred in redis keySpace notifier publisher ", v)
			goto START
		}
	}
}

// forwardUndeliveredEventToClient function read the undelivered event from db
// and forward event to destination after specified interval
func (e *ExternalInterfaces) forwardUndeliveredEventToClient(ctx context.Context) {
	for {
		for _, sub := range subscriptionsCache {
			if sub.Destination != "" {
				go e.checkUndeliveredEvents(sub.Destination)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

// RunEventForwardingWorkers will create a worker pool for forwarding
// event to destination
func (e *ExternalInterfaces) runEventForwardingWorkers() {
	for job := range eventForwardingChanel {
		e.postEvent(job)
	}
}

// RunEventForwardingWorkers will create a worker pool for forwarding
// event to destination
func (e *ExternalInterfaces) saveEventWorkers() {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		logging.Error("error occurred saveEventWorker ", err.Error())
		e.saveEventWorkers()
		return
	}
	if err != nil {
		logging.Error("error occurred get write pool in saveEventWorker ", err.Error())
		e.saveEventWorkers()
		return
	}
	for job := range saveEventChanel {

		err := conn.SaveUndeliveredEvents(evmodel.UndeliveredEvents, job.UndeliveredEventID, job.Message)
		if err != nil {
			logging.Error("error while save undelivered event ", err)
			time.Sleep(time.Second)
			if err != nil {
				continue
			}
			err = conn.SaveUndeliveredEvents(evmodel.UndeliveredEvents, job.UndeliveredEventID, job.Message)
			if err != nil {
				logging.Error("error occurred while save saveEventWorker ", err.Error())
			}
		}
	}
}
