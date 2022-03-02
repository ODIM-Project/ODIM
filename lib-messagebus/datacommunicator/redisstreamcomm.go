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
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// RedisStreamsPacket defines the RedisStreamsPacket Message Packet Object. Apart from Base Packet, it
// will contain Redis Connection Object
type RedisStreamsPacket struct {
	Packet
	pipe string
}

func getDBConnection() *redis.Client {
	var dbConn *redis.Client

	if len(MQ.RedisStreams.SentinalAddress) > 0 {
		dbConn = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    MQ.RedisStreams.SentinalAddress,
			SentinelAddrs: []string{fmt.Sprintf("%s:%s", MQ.RedisStreams.RedisServerAddress, MQ.RedisStreams.RedisServerPort)},
			MaxRetries:    -1,
		})
	} else {
		dbConn = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", MQ.RedisStreams.RedisServerAddress, MQ.RedisStreams.RedisServerPort),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
	return dbConn
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
		log.Error(e.Error())
		return e
	}
	redisClient := getDBConnection()
	defer redisClient.Close()
	_, rerr := redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: rp.pipe,
		Values: map[string]interface{}{"data": b},
	}).Result()

	if rerr != nil {
		log.Error("Unable to publish event to redis, got: " + rerr.Error())
		if rerr != nil {
			if strings.Contains(rerr.Error(), " connection timed out") {
				redisClient = getDBConnection()
				redisClient.XAdd(ctx, &redis.XAddArgs{
					Stream: rp.pipe,
					Values: map[string]interface{}{"data": b},
				}).Result()
			}
		}
		return rerr
	}

	return nil
}

// Accept implmentation need to be added
func (rp *RedisStreamsPacket) Accept(fn MsgProcess) error {
	redisClient := getDBConnection()
	var id = uuid.NewV4().String()
	rerr := redisClient.XGroupCreateMkStream(context.Background(),
		rp.pipe, EVENTREADERGROUPNAME, "$").Err()
	if rerr != nil {
		log.Error("Unable to create the group ", rerr)
		if strings.Contains(rerr.Error(), " connection timed out") {
			redisClient = getDBConnection()
		}

	}

	// create a unique consumer id for the  instance
	go rp.checkUnacknowledgedEvents(fn, id)
	go func() {
		for {
			events, err := redisClient.XReadGroup(context.Background(),
				&redis.XReadGroupArgs{
					Group:    EVENTREADERGROUPNAME,
					Consumer: id,
					Count:    1,
					Streams:  []string{rp.pipe, ">"},
				}).Result()
			if err != nil {
				log.Error("Unable to get data from the group ", err)
				if strings.Contains(err.Error(), " connection timed out") {
					redisClient = getDBConnection()
				}
			} else {

				if len(events) > 0 && len(events[0].Messages) > 0 {
					messageID := events[0].Messages[0].ID
					evtStr := events[0].Messages[0].Values["data"].(string)
					var evt interface{}
					Decode([]byte(evtStr), &evt)
					fn(evt)
					redisClient.XAck(context.Background(), rp.pipe, EVENTREADERGROUPNAME, messageID)
				}
			}
		}
	}()
	return nil
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
	return nil
}
func (rp *RedisStreamsPacket) checkUnacknowledgedEvents(fn MsgProcess, id string) {
	redisClient := getDBConnection()
	for {
		events, _, err := redisClient.XAutoClaim(context.Background(), &redis.XAutoClaimArgs{
			Stream:   rp.pipe,
			Group:    EVENTREADERGROUPNAME,
			Consumer: id,
			MinIdle:  time.Minute * 10,
			Count:    100,
			Start:    "0-0",
		}).Result()
		if err != nil {
			if strings.Contains(err.Error(), " connection timed out") {
				redisClient = getDBConnection()
			}
		}
		for _, event := range events {
			messageID := event.ID
			evtStr := event.Values["data"].(string)
			var evt interface{}
			Decode([]byte(evtStr), &evt)
			fn(evt)
			redisClient.XAck(context.Background(), rp.pipe, EVENTREADERGROUPNAME, messageID)
		}
		time.Sleep(time.Minute * 10)
	}
}
