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

// Package persistencemgr provides an  interfaces for database communication
package persistencemgr

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"sync/atomic"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	redis "github.com/go-redis/redis"
)

var inMemDBConnPool *ConnPool
var onDiskDBConnPool *ConnPool

const (
	errorCollectingData   string = "error while trying to collect data: "
	count                 int    = 1000
	writeToDBJSONErrMsg   string = "Write to DB in json form failed: "
	writeToDBErrMsg       string = "Write to DB failed : "
	redigoNil             string = "redigo: nil returned"
	noDataErrMsg          string = "no data with the with key "
	writePoolDeleteErrMsg string = "error while trying to delete data: WritePool is nil "
	writePoolNil          string = "WritePool is nil"
	errMsg                string = "error: data with key "
	stringConvErrMsg      string = "error while trying to convert the data into string: "
	deleteErrMsg          string = "error while trying to delete data: "
	dataRetrivalErrMsg    string = "error while trying to get data: "
	writeDataErrMsg       string = "error while trying to Write Transaction data: WritePool is nil"
	notFoundErrMsg        string = "no data found for the key: %v"
	intConvErrMsg         string = "error while trying to convert the data into int: "
	foundStr              string = " found"
)

// DbType is a alias name for int32
type DbType int32

const (
	//InMemory - To select in-memory db connection pool
	InMemory DbType = iota
	// OnDisk - To select in-disk db connection pool
	OnDisk
)

// Conn contains the write connection instance retrieved from the connection pool
type Conn struct {
	WriteConn *redis.Client
}

// RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	newSentinelClient(opt *redis.Options) *redis.SentinelClient
	getMasterAddrByName(mset string, snlClient *redis.SentinelClient) []string
}

type redisExtCallsImp struct{}

func (r redisExtCallsImp) newSentinelClient(opt *redis.Options) *redis.SentinelClient {
	return redis.NewSentinelClient(opt)
}

func (r redisExtCallsImp) getMasterAddrByName(masterSet string, snlClient *redis.SentinelClient) []string {
	return snlClient.GetMasterAddrByName(masterSet).Val()
}

// NewRedisExternalCalls is Constructor for RedisExternalCalls
func NewRedisExternalCalls() RedisExternalCalls {
	return &redisExtCallsImp{}
}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
}

func sentinelNewClient(dbConfig *Config) (*redis.SentinelClient, error) {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("error while trying to get tls configuration : %s", err.Error())
	}
	rdb := redisExtCalls.newSentinelClient(&redis.Options{
		Addr:      dbConfig.SentinelHost + ":" + dbConfig.SentinelPort,
		DB:        0, // use default DB
		TLSConfig: tlsConfig,
		Password:  dbConfig.Password,
	})
	return rdb, nil
}

var (
	goroutineCreated uint32
)

// GetCurrentMasterHostPort is to get the current Redis Master IP and Port from Sentinel.
func GetCurrentMasterHostPort(dbConfig *Config) (string, string, error) {
	sentinelClient, err := sentinelNewClient(dbConfig)
	if err != nil {
		return "", "", err
	}
	stringSlice := redisExtCalls.getMasterAddrByName(dbConfig.MasterSet, sentinelClient)
	var masterIP string
	var masterPort string
	if len(stringSlice) == 2 {
		masterIP = stringSlice[0]
		masterPort = stringSlice[1]
	}
	if atomic.CompareAndSwapUint32(&goroutineCreated, 0, 1) {
		go monitorFailureOver(sentinelClient)
	}

	return masterIP, masterPort, nil
}

func monitorFailureOver(sentinelClient *redis.SentinelClient) {
	pub := sentinelClient.Subscribe("+switch-master")
	var err *errors.Error

	for {
		pub.Receive()
		time.Sleep(10 * time.Second)
		if inMemDBConnPool != nil && onDiskDBConnPool != nil {
			config := getInMemoryDBConfig()
			inMemDBConnPool, err = config.Connection()
			if err != nil {
				continue
			}

			config1 := getOnDiskDBConfig()
			onDiskDBConnPool, err = config1.Connection()
			if err != nil {
				continue
			}
			resetDBWriteConnection(InMemory)
			resetDBWriteConnection(OnDisk)
		}
	}

}

// resetDBWriteConnection is used to reset the WriteConnection Pool (inmemory / OnDisk).
func resetDBWriteConnection(dbFlag DbType) error {
	switch dbFlag {
	case InMemory:
		config := getInMemoryDBConfig()
		inMemDBConnPool.Mux.Lock()
		defer inMemDBConnPool.Mux.Unlock()
		if inMemDBConnPool.WritePool != nil {
			return nil
		}
		err := inMemDBConnPool.setWritePool(config)
		if err != nil {
			return fmt.Errorf("reset of inMemory write pool failed: %s", err.Error())
		}
		return nil
	case OnDisk:
		config := getOnDiskDBConfig()
		onDiskDBConnPool.Mux.Lock()
		defer onDiskDBConnPool.Mux.Unlock()
		if onDiskDBConnPool.WritePool != nil {
			return nil
		}
		err := onDiskDBConnPool.setWritePool(config)
		if err != nil {
			return fmt.Errorf("reset of onDisk write pool failed: %s", err.Error())
		}
		return nil
	default:
		return nil
	}
}

func (p *ConnPool) setWritePool(c *Config) error {
	currentMasterIP := c.Host
	currentMasterPort := c.Port
	if config.Data.DBConf.RedisHAEnabled {
		currentMasterIP, currentMasterPort = retryForMasterIP(p, c)
	}
	if currentMasterIP == "" {
		return fmt.Errorf("unable to retrieve master ip from sentinel master election")
	}

	writePool, _ := goRedisNewClient(c, currentMasterIP, currentMasterPort)
	if writePool == nil {
		return fmt.Errorf("write pool creation failed")
	}
	p.WritePool = writePool
	p.MasterIP = currentMasterIP
	return nil
}

func retryForMasterIP(pool *ConnPool, config *Config) (currentMasterIP, currentMasterPort string) {
	for i := 0; i < 120; i++ {
		currentMasterIP, currentMasterPort, _ = GetCurrentMasterHostPort(config)
		if currentMasterIP != "" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return
}

func getInMemoryDBConfig() *Config {
	return &Config{
		Port:         config.Data.DBConf.InMemoryPort,
		Protocol:     config.Data.DBConf.Protocol,
		Host:         config.Data.DBConf.InMemoryHost,
		SentinelHost: config.Data.DBConf.InMemorySentinelHost,
		SentinelPort: config.Data.DBConf.InMemorySentinelPort,
		MasterSet:    config.Data.DBConf.InMemoryPrimarySet,
		Password:     string(config.Data.DBConf.RedisInMemoryPassword),
	}
}

func getOnDiskDBConfig() *Config {
	return &Config{
		Port:         config.Data.DBConf.OnDiskPort,
		Protocol:     config.Data.DBConf.Protocol,
		Host:         config.Data.DBConf.OnDiskHost,
		SentinelHost: config.Data.DBConf.OnDiskSentinelHost,
		SentinelPort: config.Data.DBConf.OnDiskSentinelPort,
		MasterSet:    config.Data.DBConf.OnDiskPrimarySet,
		Password:     string(config.Data.DBConf.RedisOnDiskPassword),
	}
}

// GetDBConnection is used to get the new Connection Pool for Inmemory/OnDisk DB
func GetDBConnection(dbFlag DbType) (*ConnPool, *errors.Error) {
	var err *errors.Error
	switch dbFlag {
	case InMemory:
		inMemDBConnPool, err := getInmemoryDBConnection()
		if err != nil {
			return nil, err
		}
		return inMemDBConnPool, err

	case OnDisk:
		onDiskDBConnPool, err = getOnDiskDBConnection()
		if err != nil {
			return nil, err
		}
		return onDiskDBConnPool, err

	default:
		return nil, errors.PackError(errors.UndefinedErrorType, "error invalid db type selection")
	}
}

func goRedisNewClient(dbConfig *Config, host, port string) (*redis.Client, error) {
	tlsConfig, err := getTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("error while trying to get tls configuration : %s", err.Error())
	}
	client := redis.NewClient(&redis.Options{
		Addr:            host + ":" + port,
		Password:        dbConfig.Password,
		DB:              0,
		TLSConfig:       tlsConfig,
		MaxRetries:      3,
		MaxRetryBackoff: time.Millisecond * 500,
		IdleTimeout:     time.Duration(config.Data.DBConf.MaxIdleConns) * time.Second,
		PoolSize:        config.Data.DBConf.MaxActiveConns,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
// getInmemoryDBConnection return an in-memory db connection pool
func getInmemoryDBConnection() (*ConnPool, *errors.Error) {
	var err *errors.Error
	if inMemDBConnPool == nil || inMemDBConnPool.ReadPool == nil {
		config := getInMemoryDBConfig()
		inMemDBConnPool, err = config.Connection()
		if err != nil {
			return nil, errors.PackError(err.ErrNo(), err.Error())
		}
	}
	if inMemDBConnPool.WritePool == nil {
		resetDBWriteConnection(InMemory)
	}

	return inMemDBConnPool, err

}

// getOnDiskDBConnection returns an On-Disk db connection pool
func getOnDiskDBConnection() (*ConnPool, *errors.Error) {
	var err *errors.Error
	if onDiskDBConnPool == nil || onDiskDBConnPool.ReadPool == nil {
		config := getOnDiskDBConfig()
		onDiskDBConnPool, err = config.Connection()
		if err != nil {
			return nil, errors.PackError(err.ErrNo(), err.Error())
		}
	}
	if onDiskDBConnPool.WritePool == nil {
		resetDBWriteConnection(OnDisk)
	}
	return onDiskDBConnPool, err
}


// GetWriteConnection retrieve a write connection from the connection pool
func (p *ConnPool) GetWriteConnection() (*Conn, *errors.Error) {
	if p.WritePool != nil {
		return &Conn{
			WriteConn: p.WritePool,
		}, nil
	}
	return nil, errors.PackError(errors.DBConnFailed)
}

func getTLSConfig() (*tls.Config, error) {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(config.Data.KeyCertConf.RootCACertificate)
	cert, err := tls.X509KeyPair(config.Data.KeyCertConf.RPCCertificate, config.Data.KeyCertConf.RPCPrivateKey)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		RootCAs:      pool,
		MinVersion:   config.DefaultTLSMinVersion,
		Certificates: []tls.Certificate{cert},
	}
	return cfg, nil
}

// Connection returns connection pool
// Connection does not take any input and returns a connection object used to interact with the DB
func (c *Config) Connection() (*ConnPool, *errors.Error) {
	var err error
	var masterIP, masterPort string
	connPools := &ConnPool{}
	masterIP = c.Host
	masterPort = c.Port
	if config.Data.DBConf.RedisHAEnabled {
		masterIP, masterPort, err = GetCurrentMasterHostPort(c)
		if err != nil {
			return nil, errors.PackError(errors.UndefinedErrorType, err.Error())
		}
	}
	connPools.ReadPool, err = goRedisNewClient(c, c.Host, c.Port)
	if err != nil {
		return nil, errors.PackError(errors.DBConnFailed, err.Error())
	}
	connPools.WritePool, err = goRedisNewClient(c, masterIP, masterPort)
	if err != nil {
		return nil, errors.PackError(errors.DBConnFailed, err.Error())
	}
	connPools.MasterIP = masterIP

	return connPools, nil
}

// Create will make an entry into the database with the given values
/* Create takes the following keys as input:
1."table" is a string which is used identify what kind of data we are storing.
2."data" is of type interface and is the userdata sent to be stored in DB.
3."key" is a string which acts as a unique ID to the data entry.
*/
func (p *ConnPool) Create(table, key string, data interface{}) *errors.Error {
	saveID := table + ":" + key
	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, writeToDBJSONErrMsg+err.Error())
	}
	value, createErr := p.WritePool.SetNX(saveID, jsondata, 0).Result()
	if createErr != nil {
		return errors.PackError(errors.UndefinedErrorType, writeToDBErrMsg+createErr.Error())
	}

	if !value {
		return errors.PackError(errors.DBKeyAlreadyExist, errMsg, key, " already exists")
	}

	return nil
}

//Update data
/* Update take the following keys as input:
1."uid" is a string which acts as a unique ID to fetch the data from the DB
2."data" is userdata which is of type interface sent by the user to update/patch the already existing data
*/
func (p *ConnPool) Update(table, key string, data interface{}) (string, *errors.Error) {

	if _, readErr := p.Read(table, key); readErr != nil {
		if errors.DBKeyNotFound == readErr.ErrNo() {
			return "", errors.PackError(readErr.ErrNo(), errMsg, key, " does not exist")
		}
		return "", readErr
	}
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return "", errors.PackError(errors.UndefinedErrorType, err.Error())
	}

	createErr := p.WritePool.Set(saveID, jsondata, 0).Err()
	if createErr != nil {
		return "", errors.PackError(errors.UndefinedErrorType, writeToDBErrMsg+createErr.Error())
	}

	return saveID, nil
}

// Upsert will insert new data in DB if the key does not exist.
// else it will update the existing key
/* Upsert take the following leys as input:
1."uid" is a string which acts as a unique ID to fetch the data from the DB
2."data" is user data which is of type interface sent by the user to update/patch the already existing data
*/
func (p *ConnPool) Upsert(table, key string, data interface{}) *errors.Error {
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, err.Error())
	}
	createErr := p.WritePool.Set(saveID, jsondata, 0).Err()
	if createErr != nil {
		return errors.PackError(errors.UndefinedErrorType, writeToDBErrMsg+createErr.Error())
	}

	return nil
}


// Read is for getting singular data
// Read takes "key" sting as input which acts as a unique ID to fetch specific data from DB
func (p *ConnPool) Read(table, key string) (string, *errors.Error) {
	value, err := p.ReadPool.Get(table + ":" + key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return "", errors.PackError(errors.DBKeyNotFound, noDataErrMsg, key, foundStr)
		}
		if errs, aye := isDbConnectError(err); aye {
			return "", errs
		}
		return "", errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	if value == "" {
		return "", errors.PackError(errors.DBKeyNotFound, noDataErrMsg, key, foundStr)
	}

	return value, nil
}

// ReadMultipleKeys function is used to read data for multiple keys from DB
func (p *ConnPool) ReadMultipleKeys(key []string) ([]string, *errors.Error) {
	value, err := p.ReadPool.MGet(key...).Result()
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return nil, errs
		}
		return nil, errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	if len(value) < 1 {
		return nil, errors.PackError(errors.DBKeyNotFound, noDataErrMsg, key, foundStr)
	}
		if value == nil {
		return nil, errors.PackError(errors.DBKeyNotFound,noDataErrMsg, key, foundStr)
	}
	strArr := make([]string, len(value))
	for i, v := range value {
		if s, ok := v.(string); ok {
			strArr[i] = s
		}
	}

	return strArr, nil
}

// FindOrNull is a wrapper for Read function. If requested asset doesn't exist errors.DBKeyNotFound error returned by Read is converted to nil
func (p *ConnPool) FindOrNull(table, key string) (string, error) {
	r, e := p.Read(table, key)
	if e != nil {
		switch e.ErrNo() {
		case errors.DBKeyNotFound:
			return "", nil
		default:
			return "", e
		}
	}
	return r, nil
}

var cursor uint64

// GetAllDetails will fetch all the keys present in the database
func (p *ConnPool) GetAllDetails(table string) ([]string, *errors.Error) {
	var IDs []string
	iter := p.ReadPool.Scan(cursor, table+":*", int64(count)).Iterator()
	if err := iter.Err(); err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}

	for iter.Next() {
		ID := strings.TrimPrefix(iter.Val(), table+":")
		IDs = append(IDs, ID)
	}

	return IDs, nil
}

// Delete data entry
// Read takes "key" sting as input which acts as a unique ID to delete specific data from DB
func (p *ConnPool) Delete(table, key string) *errors.Error {
	_, readErr := p.Read(table, key)
	if readErr != nil {
		return errors.PackError(errors.DBKeyNotFound, readErr.Error())
	}

	doErr := p.WritePool.Del(table + ":" + key).Err()
	if doErr != nil {
		if errs, aye := isDbConnectError(doErr); aye {
			return errors.PackError(errors.DBKeyNotFound, errs.Error())
		}
		return errors.PackError(errors.UndefinedErrorType, deleteErrMsg, doErr)
	}

	return nil
}

// DeleteMultipleKeys data entry takes "keys" array of sting as input to delete data from DB at once
func (p *ConnPool) DeleteMultipleKeys(keys []string) *errors.Error {
	tx := p.WritePool.TxPipeline()
	tx.Del(keys...)
	_, err := tx.Exec()
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data", err.Error())
	}

	return nil
}

// CleanUpDB will delete all database entries
// The flush command will be executed without warnings please be cautious in using this
func (p *ConnPool) CleanUpDB() *errors.Error {
	err := p.WritePool.FlushDB().Err()
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	return nil
}

// DeleteServer data entry without table
// Read takes "key" sting as input which acts as a unique ID to delete specific data from DB
func (p *ConnPool) DeleteServer(key string) *errors.Error {
	keys, err := p.ReadPool.Keys(key).Result()
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return errors.PackError(errors.DBKeyNotFound, errs.Error())
		}
		return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	tx := p.WritePool.TxPipeline()
	tx.Del(keys...)
	_, err = tx.Exec()
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data", err.Error())
	}

	return nil
}

// GetAllMatchingDetails will fetch all the keys which matches pattern present in the database
func (p *ConnPool) GetAllMatchingDetails(table, pattern string) ([]string, *errors.Error) {
	var IDs []string
	iter := p.ReadPool.Scan(cursor, table+":*"+pattern+"*", int64(count)).Iterator()
	if err := iter.Err(); err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}

	for iter.Next() {
		ID := strings.TrimPrefix(iter.Val(), table+":")
		IDs = append(IDs, ID)
	}
	return IDs, nil
}

// Transaction is to do a atomic operation using optimistic lock
func (p *ConnPool) Transaction(ctx context.Context, key string, cb func(context.Context, string) error) *errors.Error {
	err := p.WritePool.Watch(func(tx *redis.Tx) error {
		_, err := tx.Get(key).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		_, err = tx.TxPipelined(func(pipe redis.Pipeliner) error {
			if err := cb(ctx, key); err != nil {
				return errors.PackError(errors.UndefinedErrorType, err)
			}
			return nil
		})
		return err
	}, key)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}

	return nil
}

// isDbConnectError is for checking if error is dial connection error
func isDbConnectError(err error) (*errors.Error, bool) {
	if strings.HasSuffix(err.Error(), "connect: connection refused") || err.Error() == "EOF" {
		return errors.PackError(errors.DBConnFailed, err), true
	}
	return nil, false
}

// SaveBMCInventory function save all bmc inventory data togeter using the transaction model
func (p *ConnPool) SaveBMCInventory(data map[string]interface{}) *errors.Error {
	tx := p.WritePool.TxPipeline()
	for key, val := range data {
		jsondata, err := json.Marshal(val)
		if err != nil {
			return errors.PackError(errors.UndefinedErrorType, "Write to DB in json form failed: "+err.Error())
		}
		createErr := tx.Set(key, jsondata, 0).Err()
		if createErr != nil {
			return errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
		}
	}

	_, err := tx.Exec()
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	return nil

}

// Close closes the write connection retrieved from the connection pool
func (c *Conn) Close() {
	if c.WriteConn != nil {
		c.WriteConn.Close()
	}
}

// IsBadConn checks if the connection to DB is active or not
func (c *Conn) IsBadConn() bool {
	if err := c.WriteConn.Ping(); err != nil {
		return false
	}
	return true
}

func getSortedMapKeys(m interface{}) []string {
	var keys []string
	switch m := m.(type) {
	case map[string]interface{}:
		keys = make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
	case map[string]int64:
		keys = make([]string, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
	default:
		return nil
	}
	sort.Strings(keys)
	return keys
}

// UpdateTransaction will update the database using pipelined transaction
/* UpdateTransaction takes the following keys as input:
1."data" is of type map[string]interface{} and is the user data sent to be updated in DB.
key of map should be the key in database.
*/
func (c *Conn) UpdateTransaction(data map[string]interface{}) *errors.Error {
	var partialFailure bool = false
	if c.WriteConn == nil {
		return errors.PackError(errors.DBConnFailed)
	}
	tx := c.WriteConn.TxPipeline()
	keys := getSortedMapKeys(data)
	for _, key := range keys {
		jsondata, err := json.Marshal(data[key])
		if err != nil {
			delete(data, key)
			continue
		}
		updateErr := tx.Set(key, jsondata, 0).Err()
		if updateErr != nil {
			if isTimeOutError(updateErr) {
				return errors.PackError(errors.TimeoutError, updateErr.Error())
			}
			return errors.PackError(errors.DBUpdateFailed, updateErr.Error())
		}
	}

	cmd, err := tx.Exec()
	if err != nil {
		if isTimeOutError(err) {
			return errors.PackError(errors.TimeoutError, err.Error())
		}
		return errors.PackError(errors.DBUpdateFailed, err.Error())
	}

	for i, key := range keys {
		if cmd[i].Err() != nil {
			partialFailure = true
		} else {
			delete(data, key)
		}
	}

	if partialFailure {
		return errors.PackError(errors.TransactionPartiallyFailed, "TransactionPartiallyFailed : All keys in transaction are not updated in DB")
	}

	return nil
}


// SetExpiryTimeForKeys will create the expiry time using pipelined transaction
/* SetExpiryTimeForKeys takes the taskID  as input:
 */
func (c *Conn) SetExpiryTimeForKeys(taskKeys map[string]int64, keyExpiryInterval int) *errors.Error {
	tx := c.WriteConn.TxPipeline()
	var partialFailure bool = false
	members := getSortedMapKeys(taskKeys)
	duration := time.Duration(keyExpiryInterval) * time.Second
	for _, taskkey := range members {
		createErr := tx.Expire(taskkey, duration).Err() //check if it should be seconds or milliseconds
		if createErr != nil {
			if isTimeOutError(createErr) {
				return errors.PackError(errors.TimeoutError, createErr.Error())
			}
			return errors.PackError(errors.DBUpdateFailed, createErr.Error())
		}
	}

	cmd, err := tx.Exec()
	if err != nil {
		if isTimeOutError(err) {
			return errors.PackError(errors.TimeoutError, err.Error())
		}
		return errors.PackError(errors.DBUpdateFailed, err.Error())
	}

	for i, key := range members {
		if cmd[i].Err() != nil {
			partialFailure = true
		} else {
			delete(taskKeys, key)
		}
	}

	if partialFailure {
		return errors.PackError(errors.TransactionPartiallyFailed, "TransactionPartiallyFailed : All indices for the key are not created in DB")
	}
	return nil
}

// IsRetriable checks fi the redis db operation can be retried or not by validating the error returned by redis
func IsRetriable(err error) bool {
	if err == nil {
		return false
	}

	e := err.Error()
	redisErrorPrefixes := []string{
		"LOADING ",
		"READONLY ",
		"MOVED ",
		"TRYAGAIN ",
	}

	switch e {
	case io.EOF.Error(), io.ErrUnexpectedEOF.Error():
		return true
	case "ERR max number of clients reached":
		return true
	}

	for _, prefix := range redisErrorPrefixes {
		if strings.Contains(e, prefix) {
			return true
		}
	}

	// if instance of Error struct in errors package of lib-utilities is passed as the error,
	// conversion to timeout error would not be possible
	// So actual error should be passed to check if it is timeout
	return isTimeOutError(err)
}

func isTimeOutError(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}
	return false
}

// GetResourceDetails will fetch the key and also fetch the data
func (p *ConnPool) GetResourceDetails(key string) (string, *errors.Error) {
	var ID []string

	iter := p.ReadPool.Scan(cursor, "*"+key, int64(count)).Iterator()
	if err := iter.Err(); err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return "", errs
		}
		return "", errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}

	for iter.Next() {
		ID = strings.SplitN(iter.Val(), ":", 2)
	}
	if len(ID) < 1 {
		return "", errors.PackError(errors.DBKeyNotFound, noDataErrMsg, key)
	}
	return p.Read(ID[0], ID[1])
}

// AddResourceData will make an entry into the database with the given values
/* AddResourceData takes the following keys as input:
1."table" is a string which is used identify what kind of data we are storing.
2."data" is of type interface and is the userdata sent to be stored in DB.
3."key" is a string which acts as a unique ID to the data entry.
*/
func (p *ConnPool) AddResourceData(table, key string, data interface{}) *errors.Error {
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, writeToDBJSONErrMsg+err.Error())
	}
	createErr := p.WritePool.Set(saveID, jsondata, 0).Err()
	if createErr != nil {
		return errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
	}

	return nil
}

// SaveUndeliveredEvents method store undelivered event data in db
// takes table name ,key data and connection pool as input
func (p *ConnPool) SaveUndeliveredEvents(table, key string, data []byte) *errors.Error {
	saveID := table + ":" + key
	createErr := p.WritePool.Set(saveID, data, 0).Err()
	if createErr != nil {
return errors.PackError(errors.UndefinedErrorType, writeToDBErrMsg+createErr.Error())
		
	}

	return nil
}

// Ping will check the DB connection health
func (p *ConnPool) Ping() error {
	if err := p.WritePool.Ping().Err(); err != nil {
		return fmt.Errorf("error while pinging DB with read connection")
	}

	return nil
}

func convertToFloat(val interface{}) float64 {
	var f float64
	switch v := val.(type) {
	case int:
		f = float64(v)
	case int8:
		f = float64(v)
	case int16:
		f = float64(v)
	case int32:
		f = float64(v)
	case int64:
		f = float64(v)
	case uint:
		f = float64(v)
	case uint8:
		f = float64(v)
	case uint16:
		f = float64(v)
	case uint32:
		f = float64(v)
	case uint64:
		f = float64(v)
	case float32:
		f = float64(v)
	case float64:
		f = v
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			f = parsed
		}
	default:
		fmt.Printf("Cannot convert type %T to float64\n", val)
	}
	return f
}

// CreateIndex is used to create and save secondary index
/* CreateIndex take the following keys are input:
1. form is a map of the index to be created and the data along with it
2. uuid is the resource id with witch the value is stored
*/
func (p *ConnPool) CreateIndex(form map[string]interface{}, uuid string) error {
	for index, value := range form {
		var key string
		var val interface{}
		switch v := value.(type) {
		case int:
			key = strconv.Itoa(value.(int)) + "::" + uuid
			val = value
		case float64:
			key = strconv.FormatFloat(value.(float64), 'f', -1, 64) + "::" + uuid
			val = value
		case string:
			val = 0
			value = strings.ToLower(value.(string))
			key = value.(string) + "::" + uuid
		case []string:
			val = 0
			sliceString := strings.Join(value.([]string), " ")
			sliceString = "[" + sliceString + "]"
			sliceString = strings.ToLower(sliceString)
			key = sliceString + "::" + uuid
		case []float64:
			val = 0
			var floatString []string
			for _, v := range value.([]float64) {
				vs := strconv.FormatFloat(v, 'f', -1, 64)
				floatString = append(floatString, vs)
			}
			sliceString := strings.Join(floatString, " ")
			sliceString = "[" + sliceString + "]"
			key = sliceString + "::" + uuid
		default:
			return fmt.Errorf("error while saving index, unsupported value type %v", v)
		}
		createErr := p.WritePool.ZAdd(index, redis.Z{Score: convertToFloat(val), Member: key}).Err()
		if createErr != nil {
			return createErr
		}
	}
	return nil
}

//CreateTaskIndex is used to create secondary indexing for task service
/*Following are the input parameters for creating task index:
1. index name
2. value takes the Endtime for sorting with range
3. key if of the format `UserName::Endtime::TaskID`
*/
func (p *ConnPool) CreateTaskIndex(index string, value int64, key string) error {
	createErr := p.WritePool.ZAdd(index, redis.Z{Score: convertToFloat(value), Member: key}).Err()
	if createErr != nil {
		return createErr
	}
	return nil
}

/*
AddMemberToSet add a member to the redis set
Following are the input parameters for adding member to redis set:
1. key - redis set name
2. member - member id that to be added to the redis set
*/
func (p *ConnPool) AddMemberToSet(key string, member string) *errors.Error {
	createErr := p.WritePool.SAdd(key, member).Err()
	if createErr != nil {
		return errors.PackError(errors.DBUpdateFailed, createErr.Error())
	}
	return nil
}

/*
GetAllMembersInSet get all members in a redis set
Following are the input parameters to get embers from redis set:
1. key - redis set name
*/
func (p *ConnPool) GetAllMembersInSet(key string) ([]string, *errors.Error) {
	members, err := p.ReadPool.SMembers(key).Result()
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return members, errs
		}
		return members, errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}
	return members, nil
}

/*
RemoveMemberFromSet removes a member from the redis set
Following are the input parameters for removing member from redis set:
1. key - redis set name
2. member - member id that to be added to the redis set
*/
func (p *ConnPool) RemoveMemberFromSet(key string, member string) *errors.Error {
	deleteErr := p.WritePool.SRem(key, member).Err()
	if deleteErr != nil {
		return errors.PackError(errors.DBUpdateFailed, deleteErr.Error())
	}
	return nil
}

// GetString is used to retrive index values of type string
/* Inputs:
1. index is the index name to search with
2. cursor is the redis db cursor value
3. match is the value to match with
*/
func (p *ConnPool) GetString(index string, cursor float64, match string, regexFlag bool) ([]string, error) {
	var getList []string
	currentCursor := cursor
	match = strings.ToLower(match)
	for {
		d, cursor, getErr := p.ReadPool.ZScan(index, uint64(currentCursor), match,  int64(count)).Result()
		if getErr != nil {
			return []string{}, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
		}
		if len(d) < 1 {
			return []string{}, fmt.Errorf("no data found for the key: %v", match)

		}
		if len(d) > 1 {
			for i := 0; i < len(d); i++ {
				if d[i] != "0" {
					if regexFlag {
						getList = append(getList, d[i])
					} else {
						getList = append(getList, strings.Split(d[i], "::")[1])
					}
				}
			}
		}
		// stop when the cursor is 0
		if cursor == 0 {
			break
		}

	}
	return getList, nil
}

// // getListItems forms a list of index values from the retrieved data from database
// func getListItems(getList []string, getErr error, regexFlag bool, d interface{}) ([]string, error) {
// 	if len(d.([]interface{})) > 1 {
// 		data, err := redis.Strings(d.([]interface{})[1], getErr)
// 		if err != nil {
// 			return []string{}, fmt.Errorf(dataRetrivalErrMsg + err.Error())
// 		}
// 		for i := 0; i < len(data); i++ {
// 			if data[i] != "0" {
// 				if regexFlag {
// 					getList = append(getList, data[i])
// 				} else {
// 					getList = append(getList, strings.Split(data[i], "::")[1])
// 				}
// 			}
// 		}
// 	}
// 	return getList, nil
// }

func getUniqueSlice(inputSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, val := range inputSlice {
		if _, value := keys[val]; !value {
			keys[val] = true
			list = append(list, val)
		}
	}
	return list
}

// GetStorageList is used to storage list of capacity
/*
1.index name to search with
2. cursor is the redis cursor value
3. match is the search for list float type
4. condition is the value for condition operation
*/
func (p *ConnPool) GetStorageList(index string, cursor, match float64, condition string, regexFlag bool) ([]string, error) {
	var getList, storeList []string
	currentCursor := cursor
	for {
		data, cursor, getErr := p.ReadPool.ZScan(index, uint64(currentCursor), "*", int64(count)).Result()
		if getErr != nil {
			return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
		}
		if len(data) < 1 {
			return []string{}, fmt.Errorf("no data found for the key: %v", match)

		}
		if len(data) > 1 {
			for _, j := range data {
				if j != "0" {
					getList = append(getList, j)
				}
			}
		}
		if cursor == 0 {
			break
		}
	}
	if regexFlag {
		return getList, nil
	}
	var err error
	if storeList, err = getStoreListItemsFromList(getList, match, condition); err != nil {
		return nil, err
	}
	storeList = getUniqueSlice(storeList)
	return storeList, nil
}

// getListForStorage creates a list from the data retrieved from database
// func getListForStorage(getList []string, getErr error, d interface{}) ([]string, error) {
// 	if len(d.([]interface{})) > 1 {
// 		data, err := redis.Strings(d.([]interface{})[1], getErr)
// 		if err != nil {
// 			return nil, fmt.Errorf(dataRetrivalErrMsg + err.Error())
// 		}
// 		for _, j := range data {
// 			if j != "0" {
// 				getList = append(getList, j)
// 			}
// 		}
// 	}
// 	return getList, nil
// }

// getStoreListItemsFromList returns a storage list of capacity based on the condition satisfied
func getStoreListItemsFromList(list []string, match float64, condition string) ([]string, error) {
	var storeList []string
	for _, k := range list {
		values := strings.Split(k, "::")[0]
		id := strings.Split(k, "::")[1]
		values = strings.Replace(values, "]", "", -1)
		values = strings.Replace(values, "[", "", -1)
		valuesList := strings.Split(values, " ")
		for _, value := range valuesList {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			switch condition {
			case "eq":
				if result := big.NewFloat(match).Cmp(big.NewFloat(v)); result == 0 {
					storeList = append(storeList, id)
				}
			case "gt":
				if v > match {
					storeList = append(storeList, id)
				}
			case "ge":
				if v >= match {
					storeList = append(storeList, id)
				}
			case "lt":
				if v < match {
					storeList = append(storeList, id)
				}
			case "le":
				if v <= match {
					storeList = append(storeList, id)
				}
			}
		}
	}
	return storeList, nil
}

// GetRange is used to range over float type values
/*
1. index is the name of the index to search under
2. min is the minimum value for the search
3. max is the maximum value for the search
*/
func (p *ConnPool) GetRange(index string, min, max int, regexFlag bool) ([]string, error) {
	data, getErr := p.ReadPool.ZRangeByScore(index, redis.ZRangeBy{Min: strconv.Itoa(min), Max: strconv.Itoa(max)}).Result()
	if getErr != nil {
		return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
	}
	var getList = []string{}
	if regexFlag {
		return data, nil
	}
	for i := 0; i < len(data); i++ {

		getList = append(getList, strings.Split(data[i], "::")[1])
	}
	return getList, nil
}

// GetTaskList is used to range over float type values
/*
1. index is the name of the index to search under
2. min is the minimum value for the search
3. max is the maximum value for the search
*/
func (p *ConnPool) GetTaskList(index string, min, max int) ([]string, error) {
	data, getErr := p.ReadPool.ZRange(index, int64(min), int64(max)).Result()
	if getErr != nil {
		return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
	}
	return data, nil
}

// Del is used to delete the index key
/*
1. index is the name of the index under which the key needs to be deleted
2. key is the id of the resource to be deleted under an index
*/
func (p *ConnPool) Del(index string, k string) error {
	var currentCursor uint64
	key := "*" + k
	for {
		data, cursor, getErr := p.ReadPool.ZScan(index, currentCursor, key, int64(count)).Result()
		if getErr != nil {
			return fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
		}
		if len(data) < 1 {
			return fmt.Errorf("no data with ID found")
		}

		for _, resource := range data {
			if resource != "0" {
				if delErr := p.WritePool.ZRem(index, resource).Err(); delErr != nil {
					if errs, aye := isDbConnectError(delErr); aye {
						return errs
					}
					return fmt.Errorf("error while trying to delete data: " + delErr.Error())
				}
			}
		}

		// stop when the cursor is 0
		if cursor == 0 {
			break
		}
	}

	return nil
}

// // deleteData deletes the resource under an index from the database
// func deleteData(data []string, writeConn redis.Conn, index string, p *ConnPool) error {
// 	for _, resource := range data {
// 		if resource != "0" {
// 			_, delErr := writeConn.Do("ZREM", index, resource)
// 			if delErr != nil {
// 				if errs, aye := isDbConnectError(delErr); aye {
// 					atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
// 					return errs
// 				}
// 				return fmt.Errorf(deleteErrMsg + delErr.Error())
// 			}
// 		}
// 	}
// 	return nil
// }

// CreateEvtSubscriptionIndex is used to create and save secondary index
/* CreateSubscriptionIndex take the following keys are input:
1. index is the name of the index to be created
2. key and value are the key value pair for the index
*/
func (p *ConnPool) CreateEvtSubscriptionIndex(index string, key interface{}) error {
	const value = 0

	matchKey := strings.Replace(key.(string), "[", "\\[", -1)
	matchKey = strings.Replace(matchKey, "]", "\\]", -1)
	val, _ := p.GetEvtSubscriptions(index, matchKey)
	if len(val) > 0 {
		return fmt.Errorf("data Already Exist for the index: %v", index)
	}
	createErr := p.WritePool.ZAdd(index, redis.Z{Score: convertToFloat(value), Member: key}).Err()
	if createErr != nil {
		return createErr
	}
	return nil
}

// GetEvtSubscriptions is for to get subscription details
// 1. index is the name of the index to be created
// 2. searchKey is for search
// TODO: Add support for cursors and multiple data
func (p *ConnPool) GetEvtSubscriptions(index, searchKey string) ([]string, error) {
	var getList []string
	const cursor float64 = 0
	currentCursor := cursor
	countData, getErr := p.ReadPool.ZCount(index, "0", "0").Result()
	if getErr != nil {
		return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
	}
	d, _, getErr := p.ReadPool.ZScan(index, uint64(currentCursor), searchKey, countData).Result()
	if getErr != nil {
		return []string{}, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
	}
	if len(d) > 1 {
		for i := 0; i < len(d); i++ {
			if d[i] != "0" {
				getList = append(getList, d[i])
			}
		}
	}

	return getList, nil
}

// DeleteEvtSubscriptions is for to Delete subscription details
// 1. index is the name of the index to be created
// 2. removeKey is string parameter for remove
func (p *ConnPool) DeleteEvtSubscriptions(index, removeKey string) error {
	matchKey := strings.Replace(removeKey, "[", "\\[", -1)
	matchKey = strings.Replace(matchKey, "]", "\\]", -1)
	value, err := p.GetEvtSubscriptions(index, matchKey)
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf("no data found for the key: %v", matchKey)
	}
	for _, data := range value {
		delErr := p.WritePool.ZRem(index, data).Err()
		if delErr != nil {
			if errs, aye := isDbConnectError(delErr); aye {
				return errs
			}
		}
	}
	return nil
}

// UpdateEvtSubscriptions is for to Update subscription details
// 1. index is the name of the index to be created
// 2. key and value are the key value pair for the index
func (p *ConnPool) UpdateEvtSubscriptions(index, subscritionID string, key interface{}) error {
	err := p.DeleteEvtSubscriptions(index, subscritionID)
	if err != nil {
		return err
	}
	err = p.CreateEvtSubscriptionIndex(index, key)
	if err != nil {
		return fmt.Errorf("error while updating subscriptions")
	}
	return nil
}

// CreateDeviceSubscriptionIndex is used to create and save secondary index
/* CreateDeviceSubscriptionIndex take the following keys are input:
1. index is the name of the index to be created
2. key is for the index
*/
func (p *ConnPool) CreateDeviceSubscriptionIndex(index, hostIP, location string, originResources []string) error {
	const value = 0
	originResourceStr := "[" + strings.Join(originResources, " ") + "]"
	key := hostIP + "||" + location + "||" + originResourceStr
	// escape the square brackets before scanning
	searchKey := strings.Replace(key, "[", "\\[", -1)
	searchKey = strings.Replace(searchKey, "]", "\\]", -1)
	val, _ := p.GetDeviceSubscription(index, searchKey)
	if len(val) > 0 {
		return fmt.Errorf("data Already Exist for the index: %v", index)
	}
	createErr := p.WritePool.ZAdd(index, redis.Z{Score: convertToFloat(value), Member: key}).Err()
	if createErr != nil {
		return createErr
	}
	return nil
}

// GetDeviceSubscription is used to retrive index values of type string
/* Inputs:
1. index is the index name to search with
2. match is the value to match with
*/
// TODO : Handle cursor
func (p *ConnPool) GetDeviceSubscription(index string, match string) ([]string, error) {
	const cursor uint64 = 0
	currentCursor := cursor
	countData, getErr := p.ReadPool.ZCount(index, "0", "0").Result()
	if getErr != nil {
		return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
	}
	d, _, getErr := p.ReadPool.ZScan(index, currentCursor, match, countData).Result()
	if getErr != nil {
		return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
	}
	if len(d) < 1 {
		return []string{}, fmt.Errorf("no data found for the key: %v", match)
	}
	return d, nil
}

// DeleteDeviceSubscription is for to Delete subscription details of Device
// 1. index is the name of the index to be created
// 2. removeKey is string parameter for remove
func (p *ConnPool) DeleteDeviceSubscription(index, hostIP string) error {
	value, err := p.GetDeviceSubscription(index, hostIP+"*")
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf("no data found for the key: %v", hostIP)
	}
	for _, data := range value {
		delErr := p.WritePool.ZRem(index, data).Err()
		if delErr != nil {
			if errs, aye := isDbConnectError(delErr); aye {
				return errs
			}
		}
	}
	return nil
}

// UpdateDeviceSubscription is for to Update subscription details
// 1. index is the name of the index to be created
// 2. key and value are the key value pair for the index
func (p *ConnPool) UpdateDeviceSubscription(index, hostIP, location string, originResources []string) error {
	_, err := p.GetDeviceSubscription(index, hostIP+"[^0-9]*")
	if err != nil {
		return err
	}
	// host ip will be unique on each index in subscription of device
	// so there will be only one data
	err = p.DeleteDeviceSubscription(index, hostIP+"[^0-9]")
	if err != nil {
		return err
	}
	err = p.CreateDeviceSubscriptionIndex(index, hostIP, location, originResources)
	if err != nil {
		return fmt.Errorf("error while updating subscriptions")
	}
	return nil
}

// UpdateResourceIndex is used to update the resource inforamtion which is indexed
// form contains index name and value:key for the index
func (p *ConnPool) UpdateResourceIndex(form map[string]interface{}, uuid string) error {
	for index := range form {
		err := p.Del(index, uuid)
		if (err != nil) && (err.Error() != "no data with ID found") {
			return fmt.Errorf("error while updating index: %v", err)
		}
	}
	err := p.CreateIndex(form, uuid)
	if err != nil {
		return fmt.Errorf("error while updating index: %v", err)
	}
	return nil
}

// Incr is for incrementing the count
// Incr takes "key" string as input which acts as a unique ID to increment the count and return same
func (p *ConnPool) Incr(table, key string) (int, *errors.Error) {
	var count int
	value, err := p.WritePool.Incr(table + ":" + key).Result()
	if err != nil {

		if err.Error() == "redis: nil" {
			return count, errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
		}
		if errs, aye := isDbConnectError(err); aye {
			return count, errs
		}
		return count, errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	if reflect.TypeOf(value) == nil {
		return count, errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
	}

	return int(value), nil
}

// Decr is for decrementing the count
// Decr takes "key" string as input which acts as a unique ID to decrement the count and return same
func (p *ConnPool) Decr(table, key string) (int, *errors.Error) {
	var count int
	value, err := p.WritePool.Decr(table + ":" + key).Result()
	if err != nil {

		if err.Error() == "redis: nil" {
			return count, errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
		}
		if errs, aye := isDbConnectError(err); aye {
			return count, errs
		}
		return count, errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	if reflect.TypeOf(value) == nil {
		return count, errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
	}

	return int(value), nil
}

// SetExpire key to hold the string value and set key to timeout after a given number of seconds
/* SetExpire takes the following keys as input:
1."table" is a string which is used identify what kind of data we are storing.
2."data" is of type interface and is the userdata sent to be stored in DB.
3."key" is a string which acts as a unique ID to the data entry.
4. "expiretime" is of type int, which acts as expiry time for the key
*/
func (p *ConnPool) SetExpire(table, key string, data interface{}, expiretime int) *errors.Error {
	value, readErr := p.Read(table, key)
	if readErr != nil && readErr.ErrNo() == errors.DBConnFailed {
		return errors.PackError(readErr.ErrNo(), "error: db connection failed")
	}
	if value != "" {
		return errors.PackError(errors.DBKeyAlreadyExist, errMsg, key, " already exists")
	}
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, writeToDBJSONErrMsg+err.Error())
	}
	createErr := p.WritePool.Set(saveID, jsondata, time.Duration(expiretime)*time.Second).Err()
	if createErr != nil {
		return errors.PackError(errors.UndefinedErrorType, writeToDBErrMsg+createErr.Error())
	}

	return nil
}

// TTL is for getting singular data
// TTL takes "key" string as input which acts as a unique ID to fetch time left
func (p *ConnPool) TTL(table, key string) (int, *errors.Error) {
	value, err := p.ReadPool.TTL(table + ":" + key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, errors.PackError(errors.DBKeyNotFound, noDataErrMsg, key, foundStr)
		}
		if errs, aye := isDbConnectError(err); aye {
			return 0, errs
		}
		return 0, errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	return int(value), nil
}

// CreateAggregateHostIndex is used to create and save secondary index
/* CreateAggregateHostIndex take the following keys are input:
1. index is the name of the index to be created
2. key is for the index
*/
func (p *ConnPool) CreateAggregateHostIndex(index, aggregateID string, hostIP []string) error {
	const value = 0
	originResourceStr := "[" + strings.Join(hostIP, " ") + "]"
	key := aggregateID + "||" + originResourceStr
	createErr := p.WritePool.ZAdd(index, redis.Z{Score: convertToFloat(value), Member: key}).Err()
	if createErr != nil {
		return createErr
	}
	return nil
}

// GetAggregateHosts is used to retrive index values of type string
/* Inputs:
1. index is the index name to search with
2. match is the value to match with
*/
// TODO : Handle cursor
func (p *ConnPool) GetAggregateHosts(index string, match string) ([]string, error) {
	var data []string
	var (
		cursor uint64
	)

	for {
		data, cursorVal, getErr := p.ReadPool.ZScan(index, uint64(cursor), match, int64(count)).Result()
		if getErr != nil {
			return nil, fmt.Errorf(dataRetrivalErrMsg + getErr.Error())
		}
		if len(data) < 1 {
			return []string{}, fmt.Errorf("no data found for the key: %v", match)

		}
		if len(data) > 1 {
			return data, nil
		}

		// stop when the cursor is 0
		if cursorVal == 0 {
			break
		}
		cursor = cursorVal
	}
	return data, nil
}

// UpdateAggregateHosts is for to Update subscription details
// 1. index is the name of the index to be created
// 2. key and value are the key value pair for the index
func (p *ConnPool) UpdateAggregateHosts(index, aggregateID string, hostIP []string) error {
	err := p.DeleteAggregateHosts(index, aggregateID+"[^0-9]")
	if err != nil {
		return err
	}
	err = p.CreateAggregateHostIndex(index, aggregateID, hostIP)
	if err != nil {
		return fmt.Errorf("error while updating aggregate host ")
	}
	return nil
}

// DeleteAggregateHosts is for to Delete subscription details of aggregate
// 1. index is the name of the index to be created
// 2. removeKey is string parameter for remove
func (p *ConnPool) DeleteAggregateHosts(index, aggregateID string) error {
	value, err := p.GetAggregateHosts(index, aggregateID+"[^0-9]*")
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf(notFoundErrMsg, aggregateID)
	}
	for _, data := range value {
		delErr := p.WritePool.ZRem(index, data).Err()
		if delErr != nil {
			if errs, aye := isDbConnectError(delErr); aye {
				return errs
			}
		}
	}
	return nil
}

// GetAllDataByIndex - This function retrieves all data for a given index from sorted sets
// This maybe used to get all event/device subscriptions and aggregate hosts
func (p *ConnPool) GetAllDataByIndex(index string) ([]string, error) {
	dList, cursor, size, ferror := p.getAllDataFromSortedList(index)
	if ferror != nil {
		return []string{}, ferror
	}
	EvtSubscriptions, _, extracterr := getDataAsStringList(dList, cursor, size)
	if extracterr != nil {
		return []string{}, ferror
	}

	return EvtSubscriptions, nil
}

// getAllDataFromSortedList function read all member from index
func (p *ConnPool) getAllDataFromSortedList(index string) (data interface{}, cursr uint64, size int, err error) {
	var cursor uint64
	countData, getErr := p.ReadPool.ZCount(index, "0", "0").Result()
	if getErr != nil {
		return nil, 0, 0, fmt.Errorf("Unable to fetch count of data for index : " +
			index + " : " +
			getErr.Error())
	}
	data, cursor, getErr = p.ReadPool.ZScan(index, uint64(cursor), "*", countData).Result()
	if getErr != nil {
		return []string{}, 0, 0,
			fmt.Errorf("Error while fetching data for " + index + " : " + getErr.Error())
	}
	return data, cursor, int(countData), nil
}

// getDataAsStringList function convert list of interface into string
func getDataAsStringList(d interface{}, nextCursor uint64, size int) ([]string, int, error) {
	dataList := make([]string, 0, size)
	if len(d.([]string)) > 1 {
		for i := 0; i < len(d.([]string)); i++ {
			if d.([]string)[i] != "0" {
				dataList = append(dataList, d.([]string)[i])
			}
		}
	}
	return dataList, int(nextCursor), nil
}

// GetAllKeysFromDb will fetch all the keys which matches pattern present
// in the database using scan command, return list of key and nextCursor
func (p *ConnPool) GetAllKeysFromDb(table, pattern string, nextCursor int) ([]string, int, *errors.Error) {
	count := 500
	data, cursor, getErr := p.ReadPool.Scan(uint64(nextCursor), table+":*"+pattern+"*", int64(count)).Result()
	if getErr != nil {
		return []string{}, 0, errors.PackError(errors.UndefinedErrorType, "Error while fetching data for", getErr.Error())
	}
	keys, nextCursor, err := getDataAsStringList(data, cursor, count)
	if err != nil {
		return []string{}, 0, errors.PackError(errors.JSONUnmarshalFailed, err.Error())
	}
	return keys, int(nextCursor), nil
}

// GetKeyValue takes "key" sting as input which acts as a unique ID to fetch specific data from DB
func (p *ConnPool) GetKeyValue(key string) (string, *errors.Error) {
	value, err := p.ReadPool.Get(key).Result()
	if err != nil {

		if err.Error() == "redis: nil" {
			return "", errors.PackError(errors.DBKeyNotFound, noDataErrMsg, key, foundStr)
		}
		if errs, aye := isDbConnectError(err); aye {
			return "", errs
		}
		return "", errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	if value == "" {
		return "", errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
	}

	return value, nil
}


// DeleteKey takes "key" sting as input which acts as a unique ID
// to delete specific data from DB
func (p *ConnPool) DeleteKey(key string) *errors.Error {
	doErr := p.WritePool.Del(key).Err()
	if doErr != nil {
		if errs, aye := isDbConnectError(doErr); aye {
			return errors.PackError(errors.DBKeyNotFound, errs.Error())
		}
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data: ", doErr)
	}
	return nil
}

// EnableKeySpaceNotifier enable keyspace event notifications
// takes notifierType ad filterType as input to set value to filter redis event
func (p *ConnPool) EnableKeySpaceNotifier(notifierType, filterType string) *errors.Error {
	doErr := p.WritePool.ConfigSet(notifierType, filterType).Err()
	if doErr != nil {
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data: ", doErr)
	}
	return nil
}