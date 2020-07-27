package tmessagebus

import (
	"encoding/json"
	"log"

	dc "github.com/ODIM-Project/ODIM/lib-messagebus/datacommunicator"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	uuid "github.com/satori/go.uuid"
)

//Publish will takes the taskURI, messageID, Event type and publishes the data to message bus
func Publish(taskURI string, messageID string, eventType string) {

	k, err := dc.Communicator(dc.KAFKA, config.Data.MessageQueueConfigFilePath)
	if err != nil {
		log.Println("Unable to connect to kafka", err)
		return
	}

	defer k.Close()
	var event = common.Event{
		EventID:           uuid.NewV4().String(),
		MessageID:         messageID,
		EventType:         eventType,
		OriginOfCondition: taskURI,
	}
	var events = []common.Event{event}
	var messageData = common.MessageData{
		Name:      "Resource Event",
		Context:   "/redfish/v1/$metadata#Event.Event",
		OdataType: "#Event.v1_0_0.Event",
		Events:    events,
	}
	data, _ := json.Marshal(messageData)
	var mbevent = common.Events{
		IP:      "TasksCollection",
		Request: data,
	}

	if err := k.Distribute("REDFISH-EVENTS-TOPIC", mbevent); err != nil {
		log.Println("Unable Publish events to kafka", err)
		return
	}
	log.Printf("debug: %s Event Published", messageID)

}
