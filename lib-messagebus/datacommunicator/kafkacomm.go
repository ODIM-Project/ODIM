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

// Package datacommunicator ...
package datacommunicator

// -----------------------------------------------------------------------------
// IMPORT Section
// -----------------------------------------------------------------------------
import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaPacket defines the KAFKA Message Object. This one conains all the required
// KAFKA-GO related identifiers to maintain connection with KAFKA servers. For
// Publishing and Consuming two different Connection used towards Kafka as we are
// using Reader and Writer IO Stream Integration with RPC call. Because of the way
// Kafka communication works, we are storing these IO objects as Value for a map &
// mapped to Channel name for which these objects are created. Apart of Reader or
// Writer maps, It also maintains the Dialer Object for initial Kafka connection.
// Current Active Server name too maintained as part of KafkaPacket Object.
type KafkaPacket struct {

	// All common base function objects are defined in this object. This
	// object will support only Publishing and Subscriptions based on KAFKA
	// support. We use KAFKA 2.2.0 with Scala 2.12.
	Packet

	// Following are the map definition of both KAFKA reader and writers with Topic name.
	// Instead of using low level Conn Object from KAFKA-GO, we are using this high level
	// handle to make sure it does provide and help us with additional features like (Retry
	// or Reconnect in case of errors, Configurable distribution of messages based on
	// partitions, Sync and Async messaging, Flushing of messages in case of App shutdown.)
	// Some of the features are for Future Expansion.

	// Readers would maintain a mapping between the Kafka Reader pointer
	// and the Topic which is handled in that reader.
	Readers map[string]*kafka.Reader

	// Writers defines the mapping between KAFKA Writer pointer reference
	// and the Topic which is handled in that Writer
	Writers map[string]*kafka.Writer

	// DialerConn defines the member which can be used for single connection
	// towards KAFKA
	DialerConn *kafka.Dialer

	// ServerInfo  defines list of the KAFKA server with port
	ServersInfo []string
}

// TLS creates the TLS Configuration object to used by any Broker for Auth and
// Encryption. The Certficate and Key files are created from Java Keytool
// generated JKS format files. Please look into README for more information
// In case of Kafka, we generate the Server certificate files in JKS format.
// We do the same for Clients as well. Then we convert those files into PEM
// format.
func TLS(cCert, cKey, caCert string) (*tls.Config, error) {

	tlsConfig := tls.Config{}

	// Load client cert
	cert, e1 := tls.LoadX509KeyPair(cCert, cKey)
	if e1 != nil {
		return &tlsConfig, e1
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCertR, e2 := ioutil.ReadFile(caCert)
	if e2 != nil {
		return &tlsConfig, e2
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertR)
	tlsConfig.RootCAs = caCertPool

	tlsConfig.BuildNameToCertificate()
	return &tlsConfig, e2
}

// KafkaConnect defines the connection procedure for KAFKA Server. For now, we are
// taking only one server as input. TLS for client send would be formed as TLS
// object and same would be passed to the Server for connnection request. Common
// Dialer object will be used for both Reader and Writer objects. These objects
// would be updated if there is a request coming for specific Pipe, that specific
// Pipe name and Connection object would be stored as part of this map pair.
func KafkaConnect(kp *KafkaPacket, messageQueueConfigPath string) error {

	// Using MQF details, connecting to the KAFKA Server.
	kp.ServersInfo = mq.KServersInfo

	// Creation of TLS Config and Dialer
	tls, e := TLS(mq.KAFKACertFile, mq.KAFKAKeyFile, mq.KAFKACAFile)
	if e != nil {
		log.Error(e.Error())
		return e
	}
	kp.DialerConn = &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       tls,
	}

	// Initialize the connection map for both Reader and Writer for KAFKA
	if kp.Readers == nil {
		kp.Readers = make(map[string]*kafka.Reader)
	}
	if kp.Writers == nil {
		kp.Writers = make(map[string]*kafka.Writer)
	}

	return nil
}

// Distribute defines the Producer / Publisher role and functionality. Writer
// would be created for each Pipe comes-in for communication. If Writer already
// exists, that connection would be used for this call. Before publishing the
// message in the specified Pipe, it will be converted into Byte stream using
// "Encode" API. Encryption is enabled for the message via TLS.
func (kp *KafkaPacket) Distribute(pipe string, d interface{}) error {

	// Check for existing Writers. If not existing for this specific Pipe,
	// then we would create this Writer object for sending the message.
	if _, a := kp.Writers[pipe]; a == false {

		kp.Writers[pipe] = kafka.NewWriter(kafka.WriterConfig{
			Brokers:       kp.ServersInfo,
			Topic:         pipe,
			Balancer:      &kafka.RoundRobin{},
			BatchSize:     1,
			QueueCapacity: 1,
			Async:         true,
			Dialer:        kp.DialerConn,
		})
	}

	// Encode the message before appending into KAFKA Message struct
	b, e := Encode(d)
	if e != nil {
		log.Error(e.Error())
		return e
	}

	// Place the byte stream into Kafka.Message
	km := kafka.Message{
		Key:   []byte(pipe),
		Value: b,
	}

	// Write the messgae in the specified Pipe.
	if e = kp.Writers[pipe].WriteMessages(context.Background(), km); e != nil {
		log.Error(e.Error())
		return e
	}

	return nil
}

// Accept function defines the Consumer or Subscriber functionality for KAFKA.
// If Reader object for the specified Pipe is not available, New Reader Object
// would be created. From this function Goroutine "Read" will be invoked to
// handle the incoming messages.
func (kp *KafkaPacket) Accept(pipe string, fn MsgProcess) error {

	// If for the Reader Object for pipe and create one if required.
	if _, a := kp.Readers[pipe]; a == false {

		kp.Readers[pipe] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        kp.ServersInfo,
			GroupID:        pipe,
			Topic:          pipe,
			MinBytes:       10e1,
			MaxBytes:       10e6,
			CommitInterval: 1 * time.Second,
			Dialer:         kp.DialerConn,
		})
	}

	kp.Read(pipe, fn)
	return nil
}

// Read would access the KAFKA messages in a infinite loop. Callback method
// access is existing only in "goka" library.  Not available in "kafka-go".
func (kp *KafkaPacket) Read(p string, fn MsgProcess) error {

	// This interface should be defined outside the inner level to make sure
	// we are making the ToData API to work. Otherwise we would get exception
	// of having local scope interface pointer into passing to remote one
	var d interface{}
	c := context.Background()

	// Infinite loop to make sure we are constantly reading the messages
	// from KAFKA.
	for {
		// ReadMessages is also possible.  Here in this case, we are
		// explicitly committing the messages
		m, e := kp.Readers[p].ReadMessage(c)
		if e != nil {
			log.Error(e.Error())
			return e
		}

		// Decode the message before passing it to Callback
		if e = Decode(m.Value, &d); e != nil {
			log.Error(e.Error())
			return e
		}
		// Callback Function call.
		fn(d)
	}
}

// Get - Not supported for now in Kafka from Message Bus side due to limitations
// on the quality of the go library implementation. Will be taken-up in future.
func (kp *KafkaPacket) Get(pipe string, d interface{}) interface{} {

	return nil
}

// Remove will just remove the existing subscription. This API would check just
// the Reader map as to Distribute / Publish messages, we don't need subscription
func (kp *KafkaPacket) Remove(pipe string) error {

	es, ok := kp.Readers[pipe]
	if ok == false {
		e := fmt.Errorf("specified pipe is not subscribed yet. please check the pipe name passed")
		return e
	}
	es.Close()
	delete(kp.Readers, pipe)

	return nil
}

// Close will disconnect KAFKA Connection. This API should be called when client
// is completely closing Kafka connection, both Reader and Writer objects. We don't
// close just one channel subscription using this API. For that we would be have
// different APIs defined, called "Remove".
func (kp *KafkaPacket) Close() {

	// Closing all opened Readers Connections
	for rp, rc := range kp.Readers {
		rc.Close()
		delete(kp.Readers, rp)
	}

	// Closing all opened Writers Connections
	for wp, wc := range kp.Writers {
		wc.Close()
		delete(kp.Writers, wp)
	}
}
