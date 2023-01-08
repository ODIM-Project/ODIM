package events

import (
	"fmt"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

func Test_getAllSubscriptions(t *testing.T) {
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	config.SetUpMockConfig(t)

	mockData()
	tests := []struct {
		name string
	}{
		{
			name: "Positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllSubscriptions()
		})
	}
}

func Test_getAllAggregates(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name string
	}{
		{
			name: "Positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllAggregates()
		})
	}
}

func Test_getAllDeviceSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)
	mockData()
	tests := []struct {
		name string
	}{
		{
			name: "Positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAllDeviceSubscriptions()
		})
	}
}

func TestLoadSubscriptionData(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name string
	}{
		{
			name: "Positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadSubscriptionData()
		})
	}
}

func mockData() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	cErr := connPool.CreateDeviceSubscriptionIndex(evmodel.DeviceSubscriptionIndex, "10.10.10.10", "https://odim.2.com/EventService/Subscriptions/1", []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"})
	if cErr != nil {
		return fmt.Errorf("error while trying to save device subscription of device %v", cErr.Error())
	}
	cErr = connPool.CreateEvtSubscriptionIndex(evmodel.SubscriptionIndex, "{\"EventDestination\":{\"Context\":\"ABCDEFGHJLKJ\",\"EventTypes\":[\"Alert\"],\"EventFormatType\":\"Event\",\"DeliveryRetryPolicy\":\"RetryForever\",\"Destination\":\"https://node.odim.com:8080/Destination\",\"MessageIds\":[],\"Name\":\"Bruce\",\"OriginResources\":[\"/redfish/v1/Systems/e2616735-aa1f-49d9-9e03-bb1823b3100e.1\"],\"Protocol\":\"Redfish\",\"ResourceTypes\":[],\"SubscriptionType\":\"RedfishEvent\",\"SubordinateResources\":true},\"Hosts\":[\"10.24.0.5\"],\"SubscriptionID\":\"3ce177bf-42a4-4335-b1c1-41540a4b65d7\",\"UserName\":\"admin\"}")
	if cErr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cErr.Error())
	}
	cErr = connPool.CreateEvtSubscriptionIndex(evmodel.SubscriptionIndex, "{\"EventDestination\":{\"Context\":\"ABCDEFGHJLKJ\",\"EventTypes\":[\"Alert\"],\"EventFormatType\":\"Event\",\"DeliveryRetryPolicy\":\"RetryForever\",\"Destination\":\"https://node.odim.com:8081/Destination\",\"MessageIds\":[],\"Name\":\"Bruce\",\"OriginResources\":[\"/redfish/v1/Systems\"],\"Protocol\":\"Redfish\",\"ResourceTypes\":[],\"SubscriptionType\":\"RedfishEvent\",\"SubordinateResources\":true},\"Hosts\":[\"SystemsCollection\"],\"SubscriptionID\":\"fb496acb-7948-463b-a3d2-2206cd1f0b85\",\"UserName\":\"admin\"}")

	if cErr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cErr.Error())
	}

	cErr = connPool.CreateEvtSubscriptionIndex(evmodel.SubscriptionIndex, "{\"EventDestination\":{\"Context\":\"ABCDEFGHJLKJ\",\"EventTypes\":[\"Alert\"],\"EventFormatType\":\"Event\",\"DeliveryRetryPolicy\":\"RetryForever\",\"Destination\":\"https://node.odim.com:8082/Destination\",\"MessageIds\":[],\"Name\":\"Bruce\",\"OriginResources\":[],\"Protocol\":\"Redfish\",\"ResourceTypes\":[],\"SubscriptionType\":\"RedfishEvent\",\"SubordinateResources\":true},\"Hosts\":[],\"SubscriptionID\":\"df3f3450-bda7-4e3d-bde1-4c338be59cc7\",\"UserName\":\"admin\"}")
	if cErr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cErr.Error())
	}

	cErr = connPool.CreateEvtSubscriptionIndex(evmodel.SubscriptionIndex, "{\"EventDestination\":{\"Context\":\"ABCDEFGHJLKJ\",\"EventTypes\":[\"Alert\"],\"EventFormatType\":\"Event\",\"DeliveryRetryPolicy\":\"RetryForever\",\"Destination\":\"https://node.odim.com:8084/Destination\",\"MessageIds\":[],\"Name\":\"Bruce\",\"OriginResources\":[\"/redfish/v1/AggregationService/Aggregates/b98ab95b-9187-442a-817f-b9ec60046575\"],\"Protocol\":\"Redfish\",\"ResourceTypes\":[],\"SubscriptionType\":\"RedfishEvent\",\"SubordinateResources\":true},\"Hosts\":[\"b98ab95b-9187-442a-817f-b9ec60046575\"],\"SubscriptionID\":\"f2916a4d-f142-4179-a16c-8efd15ee6d7f\",\"UserName\":\"admin\"}")
	if cErr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cErr.Error())
	}
	cErr = connPool.CreateEvtSubscriptionIndex(evmodel.SubscriptionIndex, "{\"UserName\":\"\",\"SubscriptionID\":\"0\",\"Hosts\":[],\"EventDestination\":{\"DeliveryRetryPolicy\":\"RetryForever\",\"Destination\":\"\",\"Name\":\"default\",\"Context\":\"\",\"EventTypes\":[\"Alert\"],\"MessageIds\":[],\"Protocol\":\"Redfish\",\"SubscriptionType\":\"RedfishEvent\",\"EventFormatType\":\"\",\"SubordinateResources\":true,\"ResourceTypes\":[],\"OriginResources\":[]}}")
	if cErr != nil {
		return fmt.Errorf("error while trying to save subscription of device %v", cErr.Error())
	}
	aggregateData := evmodel.Aggregate{
		Elements: []evmodel.OdataIDLink{
			{
				OdataID: "/redfish/v1/Systems/e2616735-aa1f-49d9-9e03-bb1823b3100e.1",
			},
		},
	}
	err = connPool.Create("Aggregate", "/redfish/v1/AggregationService/Aggregates/b98ab95b-9187-442a-817f-b9ec60046575", aggregateData)
	if err != nil {
		return fmt.Errorf("error while trying to save Aggregate %v", cErr.Error())
	}
	return nil
}

func Test_getSourceId(t *testing.T) {
	type args struct {
		host string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Valid host",
			args: args{
				host: "10.10.10.10",
			},
			want: "6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
		},
		{
			name: "Positive Test SystemCollection",
			args: args{
				host: "SystemCollection",
			},
			want: "SystemCollection",
		},
		{
			name: "Invalid Host Name",
			args: args{
				host: "test.com",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := getSourceId(tt.args.host)
			if got != tt.want {
				t.Errorf("getSourceId() = %v, want %v", got, tt.want)
			}
		})
	}
}
