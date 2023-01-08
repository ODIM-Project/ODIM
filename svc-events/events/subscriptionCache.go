package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	uuid "github.com/satori/go.uuid"
)

var (
	subscriptionsCache                    = make(map[string]evmodel.Subscription)
	systemToSubscriptionsMap              = make(map[string]map[string]bool)
	aggregateIdToSubscriptionsMap         = make(map[string]map[string]bool)
	collectionToSubscriptionsMap          = make(map[string]map[string]bool)
	emptyOriginResourceToSubscriptionsMap = make(map[string]map[string]bool)
	systemIdToAggregateIdsMap             = make(map[string]map[string]bool)
	eventSourceToManagerIDMap             = make(map[string]string)
	managerIDToSystemIDsMap               = make(map[string][]string)
	managerIDToChassisIDsMap              = make(map[string][]string)
)

func LoadSubscriptionData() {
	l.Log.Debug("Event data load initialized ")
	t := time.Now()
	defer l.Log.Debug("Time take to read Complete LoadSubscriptionData ", time.Since(t))
	getAllSubscriptions()
	getAllAggregates()
	getAllDeviceSubscriptions()
}
func getAllSubscriptions() {
	subscriptionsCache = make(map[string]evmodel.Subscription)
	systemToSubscriptionsMap = make(map[string]map[string]bool)
	aggregateIdToSubscriptionsMap = make(map[string]map[string]bool)
	collectionToSubscriptionsMap = make(map[string]map[string]bool)
	emptyOriginResourceToSubscriptionsMap = make(map[string]map[string]bool)
	subscriptions, err := evmodel.GetAllEvtSubscriptions()
	if err != nil {
		l.Log.Error("Error while reading all subscription data ", err)
		return
	}
	for _, subscription := range subscriptions {
		var sub evmodel.SubscriptionResource
		err = json.Unmarshal([]byte(subscription), &sub)
		if err != nil {
			continue
		}
		loadSubscriptionCacheData(sub)
	}

}

func getAllDeviceSubscriptions() {
	eventSourceToManagerIDMap = make(map[string]string)
	t := time.Now()
	defer l.Log.Debug("Time take to read complete aggregateToHost ", time.Since(t))
	deviceSubscriptionList, err := evmodel.GetAllDeviceSubscriptions()
	if err != nil {
		l.Log.Error("Error while reading all aggregate data ", err)
		return
	}
	for _, device := range deviceSubscriptionList {
		devSub := strings.Split(device, "||")
		if strings.Contains(devSub[0], "Collection") {
			continue
		} else {
			updateCatchDeviceSubscriptionData(devSub[0], evmodel.GetSliceFromString(devSub[2]))
		}
	}
}
func updateCatchDeviceSubscriptionData(key string, originResources []string) {
	systemId := originResources[0][strings.LastIndexByte(originResources[0], '/')+1:]
	eventSourceToManagerIDMap[key] = systemId
}

func loadSubscriptionCacheData(sub evmodel.SubscriptionResource) {
	if len(sub.EventDestination.OriginResources) == 0 && sub.SubscriptionID != "0" {
		subCache := evmodel.Subscription{
			Id:                   sub.SubscriptionID,
			Destination:          sub.EventDestination.Destination,
			EventTypes:           sub.EventDestination.EventTypes,
			MessageIds:           sub.EventDestination.MessageIds,
			SubordinateResources: sub.EventDestination.SubordinateResources,
			ResourceTypes:        sub.EventDestination.ResourceTypes,
			SubscriptionType:     sub.EventDestination.SubscriptionType,
			OriginResources:      sub.EventDestination.OriginResources,
			DeliveryRetryPolicy:  sub.EventDestination.DeliveryRetryPolicy,
		}
		addEmptyOriginSubscriptionCache(subCache.Id)
		subscriptionsCache[subCache.Id] = subCache
	} else {
		for _, host := range sub.Hosts {
			subCache := evmodel.Subscription{
				Id:                   sub.SubscriptionID,
				Destination:          sub.EventDestination.Destination,
				EventTypes:           sub.EventDestination.EventTypes,
				MessageIds:           sub.EventDestination.MessageIds,
				SubordinateResources: sub.EventDestination.SubordinateResources,
				ResourceTypes:        sub.EventDestination.ResourceTypes,
				SubscriptionType:     sub.EventDestination.SubscriptionType,
				OriginResources:      sub.EventDestination.OriginResources,
				DeliveryRetryPolicy:  sub.EventDestination.DeliveryRetryPolicy,
			}
			addSubscriptionCache(host, subCache.Id)
			subscriptionsCache[subCache.Id] = subCache
		}
	}
}

//addSubscriptionCache add subscription in corresponding cache based on key type
func addSubscriptionCache(key string, subId string) {
	if strings.Contains(key, "Collection") {
		data, isExists := collectionToSubscriptionsMap[key]
		if isExists {
			data[subId] = true
			collectionToSubscriptionsMap[key] = data
		} else {
			data := make(map[string]bool)
			data[subId] = true
			collectionToSubscriptionsMap[key] = data
		}
		return
	} else {
		_, err := uuid.FromString(key)
		if err == nil {
			addAggregateSubscriptionCache(key, subId)
			return
		} else {
			data, isExists := systemToSubscriptionsMap[key]
			if isExists {
				data[subId] = true
				systemToSubscriptionsMap[key] = data
			} else {
				data := make(map[string]bool)
				data[subId] = true
				systemToSubscriptionsMap[key] = data
			}
			return
		}
	}
}

func addEmptyOriginSubscriptionCache(subscriptionId string) {
	data, isExists := emptyOriginResourceToSubscriptionsMap["0"]
	if isExists {
		data[subscriptionId] = true
		emptyOriginResourceToSubscriptionsMap["0"] = data
	} else {
		emptyOriginResourceToSubscriptionsMap["0"] = map[string]bool{subscriptionId: true}
	}
}
func addAggregateSubscriptionCache(key, subId string) {
	data, isExists := aggregateIdToSubscriptionsMap[key]
	if isExists {
		data[subId] = true
		aggregateIdToSubscriptionsMap[key] = data
	} else {
		aggregateIdToSubscriptionsMap[key] = map[string]bool{subId: true}
	}
}

func getAllAggregates() {
	systemIdToAggregateIdsMap = make(map[string]map[string]bool)

	aggregateUrls, err := evmodel.GetAllAggregates()
	if err != nil {
		l.Log.Debug("Exception getting aggregates url list ", err)
		return
	}
	if len(aggregateUrls) == 0 {
		l.Log.Debug("No Aggregates found ", aggregateUrls)
		return
	}
	for _, aggregateUrl := range aggregateUrls {
		aggregate, err := evmodel.GetAggregate(aggregateUrl)
		if err != nil {
			return
		}
		aggregateId := aggregateUrl[strings.LastIndexByte(aggregateUrl, '/')+1:]
		addSystemIdToAggregateCache(aggregateId, aggregate)
	}
}
func addSystemIdToAggregateCache(aggregateUrl string, aggregate evmodel.Aggregate) {
	for _, ids := range aggregate.Elements {
		ids.OdataID = ids.OdataID[strings.LastIndexByte(strings.TrimSuffix(ids.OdataID, "/"), '/')+1:]
		aggregateIds, isExists := systemIdToAggregateIdsMap[ids.OdataID]
		if isExists {
			aggregateIds[aggregateUrl] = true
			systemIdToAggregateIdsMap[ids.OdataID] = aggregateIds
		} else {
			systemIdToAggregateIdsMap[ids.OdataID] = map[string]bool{aggregateUrl: true}
		}
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
