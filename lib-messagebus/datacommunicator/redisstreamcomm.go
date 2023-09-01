//(C) Copyright [2021] Hewlett Packard Enterprise Development LP
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

package datacommunicator

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
)

var dbConn *redis.Client

const (
	// DefaultTLSMinVersion is default minimum version for tls
	DefaultTLSMinVersion = tls.VersionTLS12
	// TimeoutErrMsg is the connection time out error message
	TimeoutErrMsg string = " connection timed out"
)

// RedisStreamsPacket defines the RedisStreamsPacket Message Packet Object. Apart from Base Packet, it
// will contain Redis Connection Object
type RedisStreamsPacket struct {
	client *redis.Client
	Packet
	pipe string
}

func (rp *RedisStreamsPacket) getDBConnection() error {

	// Assigning the existing the db connection to RedisStreamsPacket.
	// This connection will be validated. If it is a bad connection,
	// a new connection will be created and assigned to RedisStreamsPacket
	rp.client = dbConn
	if rp.Ping() {
		return nil
	}

	// closing the existing connection as it is corrupted or nil.
	// A new connection will be created and assigned to RedisStreamsPacket
	rp.Close()

	tlsConfig, e := TLS(MQ.RedisStreams.RedisCertFile, MQ.RedisStreams.RedisKeyFile, MQ.RedisStreams.RedisCAFile)
	if e != nil {
		return fmt.Errorf("error while trying to get DB connection: %s", e.Error())
	}

	tlsConfig.MinVersion = DefaultTLSMinVersion

	if len(MQ.RedisStreams.SentinalAddress) > 0 {
		dbConn = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       MQ.RedisStreams.SentinalAddress,
			SentinelAddrs:    []string{fmt.Sprintf("%s:%s", MQ.RedisStreams.RedisServerAddress, MQ.RedisStreams.RedisServerPort)},
			MaxRetries:       -1,
			TLSConfig:        tlsConfig,
			SentinelPassword: string(MQ.RedisStreams.RedisInMemoryPassword),
			Password:         string(MQ.RedisStreams.RedisInMemoryPassword),
		})
	}

	dbConn = redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%s:%s", MQ.RedisStreams.RedisServerAddress, MQ.RedisStreams.RedisServerPort),
		TLSConfig: tlsConfig,
		Password:  string(MQ.RedisStreams.RedisInMemoryPassword),
		DB:        0, // use default DB
	})

	rp.client = dbConn
	return nil
}

// Ping function is used to test the db connection with ping command
func (rp *RedisStreamsPacket) Ping() bool {
	if rp.client != nil {
		_, err := rp.client.Ping(rp.client.Context()).Result()
		return err == nil
	}
	return false
}

// Distribute defines the Producer / Publisher role and functionality. Writer
// would be created for each Pipe comes-in for communication. If Writer already
// exists, that connection would be used for this call. Before publishing the
// message in the specified Pipe, it will be converted into Byte stream using
func (rp *RedisStreamsPacket) Distribute(data interface{}) error {
	ctx := context.Background()
	// Encode the message before appending into Redis Message struct
	b, e := Encode(data)
	if e != nil {
		return fmt.Errorf("while trying to encode message: %s", e.Error())
	}

	_, rerr := rp.client.XAdd(ctx, &redis.XAddArgs{
		Stream: rp.pipe,
		Values: map[string]interface{}{"data": b},
	}).Result()

	if rerr != nil {
		if strings.Contains(rerr.Error(), TimeoutErrMsg) {
			err := rp.getDBConnection()
			if err != nil {
				return err
			}
			rp.client.XAdd(ctx, &redis.XAddArgs{
				Stream: rp.pipe,
				Values: map[string]interface{}{"data": b},
			}).Result()
		}
		return fmt.Errorf("unable to publish event to redis, got: %s", rerr.Error())
	}

	return nil
}

// Accept implmentation need to be added
func (rp *RedisStreamsPacket) Accept(fn MsgProcess) error {

	// create a unique consumer id for the  instance
	var err error
	var id = uuid.NewV4().String()
	rerr := rp.client.XGroupCreateMkStream(context.Background(),
		rp.pipe, EVENTREADERGROUPNAME, "$").Err()
	if rerr != nil && strings.Contains(rerr.Error(), TimeoutErrMsg) {
		if err := rp.getDBConnection(); err != nil {
			return err
		}
	}
	// errChan to hold the errors faced in the  below go-rotines
	errChan := make(chan error)
	go rp.checkUnacknowledgedEvents(fn, id, errChan)
	if err = <-errChan; err != nil {
		return err
	}

	go func() {
		for {
			events, err := rp.client.XReadGroup(context.Background(),
				&redis.XReadGroupArgs{
					Group:    EVENTREADERGROUPNAME,
					Consumer: id,
					Count:    1,
					Streams:  []string{rp.pipe, ">"},
				}).Result()
			if err != nil {
				errChan <- fmt.Errorf("unable to get data from the group %s", err.Error())
				if strings.Contains(err.Error(), TimeoutErrMsg) {
					err := rp.getDBConnection()
					if err != nil {
						errChan <- err
						return
					}
				}
			} else {
				processEvent(rp, events, errChan, fn)
			}
		}
	}()

	// channel to handle the errors occured during go routines
	if err = <-errChan; err != nil {
		return err
	}
	return nil
}

// processEvent processes the redis stream events and decode the data
func processEvent(rp *RedisStreamsPacket, events []redis.XStream, errChan chan<- error, fn MsgProcess) {
	if len(events) > 0 && len(events[0].Messages) > 0 {
		messageID := events[0].Messages[0].ID
		evtStr := events[0].Messages[0].Values["data"].(string)
		var evt interface{}
		err := Decode([]byte(evtStr), &evt)
		if err != nil {
			errChan <- err
			return
		}
		fn(evt)
		rp.client.XAck(context.Background(), rp.pipe, EVENTREADERGROUPNAME, messageID)
	}
}

// Read implmentation need to be added
func (rp *RedisStreamsPacket) Read(fn MsgProcess) error {
	return nil
}

// Get - Not supported for now in RedisStreams from Message Bus side due to limitations
func (rp *RedisStreamsPacket) Get(pipe string, d interface{}) interface{} {

	return nil
}

// Remove implmentation need to be added
func (rp *RedisStreamsPacket) Remove() error {
	return nil
}

// Close implmentation need to be added
func (rp *RedisStreamsPacket) Close() error {
	if rp.client != nil {
		return rp.client.Close()
	}
	return nil
}

func (rp *RedisStreamsPacket) checkUnacknowledgedEvents(fn MsgProcess, id string, errChan chan<- error) {

	for {
		events, _, err := rp.client.XAutoClaim(context.Background(), &redis.XAutoClaimArgs{
			Stream:   rp.pipe,
			Group:    EVENTREADERGROUPNAME,
			Consumer: id,
			MinIdle:  time.Minute * 10,
			Count:    100,
			Start:    "0-0",
		}).Result()
		if err != nil && strings.Contains(err.Error(), TimeoutErrMsg) {
			if err = rp.getDBConnection(); err != nil {
				errChan <- err
				return
			}
		}
		for _, event := range events {
			messageID := event.ID
			evtStr := event.Values["data"].(string)
			var evt interface{}
			Decode([]byte(evtStr), &evt)
			if err != nil {
				errChan <- err
				return
			}
			fn(evt)
			rp.client.XAck(context.Background(), rp.pipe, EVENTREADERGROUPNAME, messageID)
		}
		// Pass the nil to errChan when no error encountered
		errChan <- nil
		time.Sleep(time.Minute * 10)
	}
}
