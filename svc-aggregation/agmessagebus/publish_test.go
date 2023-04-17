package agmessagebus

import (
	"context"
	"fmt"
	"testing"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	common "github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/stretchr/testify/assert"
)

// Define an interface for MQBusCommunicator
func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}

// Define a mock MQBus object that implements the MQBus interface
var (
	communicateerror bool
	distributeerror  bool
)

type MockMQBus struct {
	Name string
}

func (m *MockMQBus) Distribute(event interface{}) error {
	// Do nothing for now
	if distributeerror {
		return fmt.Errorf("Error while posting the message to message bus")

	}
	return nil
}

func (m *MockMQBus) Accept(process dc.MsgProcess) error {
	// Do nothing for now)
	return nil
}
func (m *MockMQBus) Close() error {
	// Do nothing for now
	return nil
}
func (m *MockMQBus) Get(key string, value interface{}) interface{} {
	// Do nothing for now
	return value
}
func (m *MockMQBus) Remove() error {
	// Do nothing for now
	return nil
}

// Define a mock MQBusCommunicator that returns a mock MQBus object
type MockMQBusCommunicator struct{}

func MockCommunicator(bt string, messageQueueConfigPath string, pipe string) (dc.MQBus, error) {
	if communicateerror {
		return nil, fmt.Errorf("Error while connecting to messagebus")

	}
	mp := new(MockMQBus)
	mp.Name = "mock"
	return mp, nil

}

func TestPublish(t *testing.T) {
	// Set up test data
	systemID := "systemID"
	eventType := "ResourceAdded"
	collectionType := "collectionType"
	config.SetUpMockConfig(t)
	tests := []struct {
		name            string
		wantErr         bool
		communicatorErr bool
		distributedErr  bool
	}{
		{
			name:            "Positive Case",
			wantErr:         false,
			communicatorErr: false,
			distributedErr:  false,
		},
		{
			name:            "Kafka COnnection Failure",
			wantErr:         true,
			communicatorErr: true,
			distributedErr:  false,
		},
		{
			name:            "Failure While passing message to message bus",
			wantErr:         true,
			communicatorErr: false,
			distributedErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MQBusCommunicatorMock := MQBusCommunicator{
				Communicator: MockCommunicator,
			}
			if tt.distributedErr {
				distributeerror = true

			}
			if tt.communicatorErr {
				communicateerror = true
			}
			ctx := mockContext()
			err := Publish(ctx, systemID, eventType, collectionType, MQBusCommunicatorMock)
			if tt.wantErr {
				assert.NotEqual(t, nil, err, "Error should not be Nil for this scenario")

			} else {
				assert.Equal(t, nil, err, "Error should be Nil for this scenario")

			}

		})
	}

}

// ////////////////////////////////
func TestPublishCtrlMsg(t *testing.T) {
	communicateerror = false
	distributeerror = false
	msgType := 0
	msg := "test message"
	config.SetUpMockConfig(t)
	test := []struct {
		name1           string
		wantErr1        bool
		communicatorErr bool
		distributionErr bool
	}{
		{
			name1:           "Positive Case",
			wantErr1:        false,
			communicatorErr: false,
			distributionErr: false,
		},
		{
			name1:           "Kafka COnnection Failure",
			wantErr1:        true,
			communicatorErr: true,
			distributionErr: false,
		},
		{
			name1:           "Failure While passing message to message bus",
			wantErr1:        true,
			communicatorErr: false,
			distributionErr: true,
		},
	}
	for _, tt := range test {
		t.Run(tt.name1, func(t *testing.T) {
			MQBusCommunicatorMock := MQBusCommunicator{
				Communicator: MockCommunicator,
			}
			if tt.distributionErr {
				distributeerror = true

			}
			if tt.communicatorErr {
				communicateerror = true
			}
			err := PublishCtrlMsg(common.ControlMessage(msgType), msg, MQBusCommunicatorMock)
			if tt.wantErr1 {
				assert.NotEqual(t, nil, err, "Error should not be Nil for this scenario")

			} else {
				assert.Equal(t, nil, err, "Error should be Nil for this scenario")

			}

		})
	}

}
