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
/*

Package datacommunicator facilitates communication between different Microservice
components. It allows user to select the underlying communication mediums. For now,
Kafka is the only communication middlewares supported by MessageBus in this release.

This component provides the basic building block for enabling communication by allowing
users / clients to send any user defined information over to the remote applications / systems

Primary users of this library are all the plugins and odimra - Event service component

There are two small sample applications implemented to behave as MessageBus Producer and
Consumer. Both are existing under "messagebus/mqproducer" and "messagebus/mqconsumer"
directories defined under "messagebus" folder.

Because of the mocking issues faced with these communication middleware libraries implementation
in Golang, developer should be able to use these two applications for defining all the
testing procedures for MessageBus component

Usage Example

KAFKA Producer Client :
	dc.Enable(Person{})
	K, _ := dc.Communicator(dc.KAFKA, nil) // This obj won't handle Sync call
	defer K.Close()
	K.Distribute("example.topic", P)  // P object of Person{}

KAFKA Consumer Client :
	dc.Enable(Person{})
	K, _ := dc.Communicator(dc.KAFKA, nil) // KAFKA doesn't support Sync calls
	N.Accept("example.topic", KAfkaHandler)
	// Incoming Message would be handled in "KafkaHandler" API in client side

*/

package datacommunicator
