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
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

var (
	subscriptionsCache                    map[string]dmtf.EventDestination
	systemToSubscriptionsMap              map[string]map[string]bool
	aggregateIdToSubscriptionsMap         map[string]map[string]bool
	collectionToSubscriptionsMap          map[string]map[string]bool
	emptyOriginResourceToSubscriptionsMap map[string]bool
	systemIdToAggregateIdsMap             map[string]map[string]bool
	eventSourceToManagerIDMap             map[string]string
	managerIDToSystemIDsMap               map[string][]string
	managerIDToChassisIDsMap              map[string][]string
	subscribeCacheLock                    sync.Mutex
)

// LoadSubscriptionData method calls whenever service is started
// Here we load Subscription, DeviceSubscription, AggregateToHost
// table data into cache memory
func LoadSubscriptionData(ctx context.Context) {
	l.LogWithFields(ctx).Debug("Event cache is initialized")
	getAllSubscriptions(ctx)
	getAllAggregates(ctx)
	getAllDeviceSubscriptions(ctx)
	threadID := 1
	ctx = context.WithValue(ctx, common.ThreadID, strconv.Itoa(threadID))
	go initializeDbObserver(ctx)
}

// getAllSubscriptions this method read data from Subscription table and
// load in corresponding cache
func getAllSubscriptions(ctx context.Context) {
	subscribeCacheLock.Lock()
	defer subscribeCacheLock.Unlock()
	systemToSubscriptionsMap = make(map[string]map[string]bool)
	aggregateIdToSubscriptionsMap = make(map[string]map[string]bool)
	collectionToSubscriptionsMap = make(map[string]map[string]bool)
	emptyOriginResourceToSubscriptionsMap = make(map[string]bool)
	subscriptions, err := evmodel.GetAllEvtSubscriptions()
	if err != nil {
		l.LogWithFields(ctx).Error("Error while reading all subscription data ", err)
		return
	}

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
			emptyOriginResourceToSubscriptionsMap[sub.SubscriptionID] = true
		} else {
			loadSubscriptionCacheData(sub.SubscriptionID, sub.Hosts)
		}
	}
	l.LogWithFields(ctx).Debug("Subscriptions cache updated ")
}

// getAllDeviceSubscriptions method fetch data from DeviceSubscription table
func getAllDeviceSubscriptions(ctx context.Context) {
	subscribeCacheLock.Lock()
	defer subscribeCacheLock.Unlock()
	deviceSubscriptionList, err := evmodel.GetAllDeviceSubscriptions()
	if err != nil {
		l.LogWithFields(ctx).Error("Error while reading all aggregate data ", err)
		return
	}
	eventSourceToManagerIDMap = make(map[string]string, len(deviceSubscriptionList))
	for _, device := range deviceSubscriptionList {
		devSub := strings.Split(device, "||")
		updateCatchDeviceSubscriptionData(devSub[0], evmodel.GetSliceFromString(devSub[2]))
	}
	l.LogWithFields(ctx).Debug("DeviceSubscription cache updated ")
}

// updateCatchDeviceSubscriptionData update eventSourceToManagerMap for each key with their system IDs
func updateCatchDeviceSubscriptionData(key string, originResources []string) {
	systemId := originResources[0][strings.LastIndexByte(originResources[0], '/')+1:]
	eventSourceToManagerIDMap[key] = systemId
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
func addSubscriptionCache(key string, subscriptionId string) {
	if strings.Contains(key, "Collection") {
		updateCacheMaps(key, subscriptionId, collectionToSubscriptionsMap)
		return
	} else {
		_, err := uuid.FromString(key)
		if err == nil {
			updateCacheMaps(key, subscriptionId, aggregateIdToSubscriptionsMap)
			return
		} else {
			updateCacheMaps(key, subscriptionId, systemToSubscriptionsMap)
			return
		}
	}
}

// getAllAggregates method will read all aggregate from db and
// update systemIdToAggregateIdsMap to corresponding member in aggregate
func getAllAggregates(ctx context.Context) {
	subscribeCacheLock.Lock()
	defer subscribeCacheLock.Unlock()

	systemIdToAggregateIdsMap = make(map[string]map[string]bool)
	aggregateUrls, err := evmodel.GetAllAggregates()
	if err != nil {
		l.LogWithFields(ctx).Debug("error occurred while getting aggregate list ", err)
		return
	}
	for _, aggregateUrl := range aggregateUrls {
		aggregate, err := evmodel.GetAggregate(aggregateUrl)
		if err != nil {
			continue
		}
		aggregateId := aggregateUrl[strings.LastIndexByte(aggregateUrl, '/')+1:]
		addSystemIdToAggregateCache(aggregateId, aggregate)
	}
	l.LogWithFields(ctx).Debug("AggregateToHost cache updated ")
}

// addSystemIdToAggregateCache update cache for each aggregate member
func addSystemIdToAggregateCache(aggregateId string, aggregate evmodel.Aggregate) {
	for _, ids := range aggregate.Elements {
		ids.Oid = ids.Oid[strings.LastIndexByte(strings.TrimSuffix(ids.Oid, "/"), '/')+1:]
		updateCacheMaps(ids.Oid, aggregateId, systemIdToAggregateIdsMap)
	}
}

//getSourceId function return system id corresponding host, if not found then return host
func getSourceId(host string) (string, error) {
	data, isExists := eventSourceToManagerIDMap[host]
	if !isExists {
		if strings.Contains(host, "Collection") {
			return host, nil
		} else {
			return "", fmt.Errorf("invalid source")
		}
	}
	return data, nil
}

//updateCacheMaps update map value corresponding key
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
func getSubscriptions(originOfCondition, systemId, hostIp string) (subs []dmtf.EventDestination) {
	subs = append(subs, getSystemSubscriptionList(hostIp)...)
	subs = append(subs, getAggregateSubscriptionList(systemId)...)
	subs = append(subs, getEmptyOriginResourceSubscriptionList()...)
	subs = append(subs, getCollectionSubscriptionList(originOfCondition, hostIp)...)
	return
}

//getSystemSubscriptionList return list of subscription corresponding to host
func getSystemSubscriptionList(hostIp string) (subs []dmtf.EventDestination) {
	systemSubscription, isExists := systemToSubscriptionsMap[hostIp]
	if isExists {
		for subId, _ := range systemSubscription {
			sub, isValidSubId := getSubscriptionDetails(subId)
			if isValidSubId {
				subs = append(subs, sub)
			}

		}
	}
	return
}

// getAggregateSubscriptionList return list of subscription corresponding to system
// is members of different aggregate
func getAggregateSubscriptionList(systemId string) (subs []dmtf.EventDestination) {
	aggregateList, isExists := systemIdToAggregateIdsMap[systemId]
	if isExists {
		for aggregateID := range aggregateList {
			subscriptions, isValidAggregateId := aggregateIdToSubscriptionsMap[aggregateID]
			if isValidAggregateId {
				for subId := range subscriptions {
					sub, isValidSubId := getSubscriptionDetails(subId)
					sub.OriginResources = append(sub.OriginResources, model.Link{Oid: "/redfish/v1/Systems/" + systemId})
					if isValidSubId {
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
func getCollectionSubscriptionList(originOfCondition, hostIp string) (subs []dmtf.EventDestination) {
	collectionsKey := getCollectionKey(originOfCondition, hostIp)
	collectionSubscription, isExists := collectionToSubscriptionsMap[collectionsKey]
	if isExists {
		for subId := range collectionSubscription {
			sub, isValidSubId := getSubscriptionDetails(subId)
			if isValidSubId {
				subs = append(subs, sub)
			}
		}
	}
	return
}

// getEmptyOriginResourceSubscriptionList return list of subscription
// whose originResources is empty
func getEmptyOriginResourceSubscriptionList() (subs []dmtf.EventDestination) {
	for subId := range emptyOriginResourceToSubscriptionsMap {
		sub, isValidSubId := getSubscriptionDetails(subId)
		if isValidSubId {
			subs = append(subs, sub)
		}
	}
	return
}

//getSubscriptionDetails this method return subscription details corresponding subscription Id
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

// initializeDbObserver function subscribe redis keyspace notifier
func initializeDbObserver(ctx context.Context) {
	l.LogWithFields(ctx).Debug("Initializing observer ")
START:
	conn, _ := common.GetDBConnection(common.OnDisk)
	writeConn := conn.WritePool.Get()
	defer writeConn.Close()
	_, err := writeConn.Do("CONFIG", "SET", evcommon.RedisNotifierType, evcommon.RedisNotifierFilterKey)
	if err != nil {
		l.LogWithFields(ctx).Error("error occurred configuring key event ", err)
		time.Sleep(time.Second * 1)
		goto START
	}
	psc := redis.PubSubConn{Conn: writeConn}

	psc.PSubscribe(evcommon.AggregateToHostChannelKey, evcommon.DeviceSubscriptionChannelKey,
		evcommon.SubscriptionChannelKey)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			switch string(v.Pattern) {
			case evcommon.DeviceSubscriptionChannelKey:
				getAllDeviceSubscriptions(ctx)
			case evcommon.SubscriptionChannelKey:
				getAllSubscriptions(ctx)
			case evcommon.AggregateToHostChannelKey:
				getAllAggregates(ctx)
			}
		case error:
			l.LogWithFields(ctx).Error("Error occurred in redis keyspace notifier publisher ", v)
			goto START
		}
	}
}
