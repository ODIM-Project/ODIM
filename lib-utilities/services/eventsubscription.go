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

package services

import (
	"context"
	"fmt"

	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
)

// SubscribeToEMB method will subscribe to respective  event queue of the plugin
func SubscribeToEMB(pluginID string, queueList []string) error {
	conn, errConn := ODIMService.Client(Events)
	if errConn != nil {
		return fmt.Errorf("Failed to create client connection: %s", errConn.Error())
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)
	_, err := events.SubsribeEMB(context.TODO(), &eventsproto.SubscribeEMBRequest{
		PluginID:     pluginID,
		EMBQueueName: queueList,
	})
	if err != nil {
		return fmt.Errorf("error subscribing to EMB  %s", err.Error())
	}
	return nil
}

// DeleteSubscription  calls the event service and delete all subscription realated to that server
func DeleteSubscription(uuid string) (*eventsproto.EventSubResponse, error) {
	var resp eventsproto.EventSubResponse
	req := eventsproto.EventRequest{
		UUID: uuid,
	}
	conn, errConn := ODIMService.Client(Events)
	if errConn != nil {
		return &resp, fmt.Errorf("Failed to create client connection: %v", errConn)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	return events.DeleteEventSubscription(context.TODO(), &req)
}
