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

//Package persistencemgr provides an  interfaces for database communication
package persistencemgr

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	redisSentinel "github.com/go-redis/redis"
	"github.com/gomodule/redigo/redis"
)

var inMemDBConnPool *ConnPool
var onDiskDBConnPool *ConnPool

const (
	errorCollectingData string = "error while trying to collect data: "
	count               int    = 1000
)

// DbType is a alias name for int32
type DbType int32

const (
	//InMemory - To select in-memory db connection pool
	InMemory DbType = iota
	// OnDisk - To select in-disk db connection pool
	OnDisk
)

//RedisExternalCalls containes the methods to make calls to external client libraries of Redis DB
type RedisExternalCalls interface {
	newSentinelClient(opt *redisSentinel.Options) *redisSentinel.SentinelClient
	getMasterAddrByName(mset string, snlClient *redisSentinel.SentinelClient) []string
}

type redisExtCallsImp struct{}

func (r redisExtCallsImp) newSentinelClient(opt *redisSentinel.Options) *redisSentinel.SentinelClient {
	return redisSentinel.NewSentinelClient(opt)
}

func (r redisExtCallsImp) getMasterAddrByName(masterSet string, snlClient *redisSentinel.SentinelClient) []string {
	return snlClient.GetMasterAddrByName(masterSet).Val()
}

//NewRedisExternalCalls is Constructor for RedisExternalCalls
func NewRedisExternalCalls() RedisExternalCalls {
	return &redisExtCallsImp{}
}

var redisExtCalls RedisExternalCalls

func init() {
	redisExtCalls = redisExtCallsImp{}
}

func sentinelNewClient(dbConfig *Config) *redisSentinel.SentinelClient {
	rdb := redisExtCalls.newSentinelClient(&redisSentinel.Options{
		Addr:     dbConfig.Host + ":" + dbConfig.SentinelPort,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}

//GetCurrentMasterHostPort is to get the current Redis Master IP and Port from Sentinel.
func GetCurrentMasterHostPort(dbConfig *Config) (string, string) {
	sentinelClient := sentinelNewClient(dbConfig)
	stringSlice := redisExtCalls.getMasterAddrByName(dbConfig.MasterSet, sentinelClient)
	var masterIP string
	var masterPort string
	if len(stringSlice) == 2 {
		masterIP = stringSlice[0]
		masterPort = stringSlice[1]
	}
	log.Info("GetCurrentMasterHostPort masterIP : "+masterIP, ", masterPort : "+masterPort)
	return masterIP, masterPort
}

//resetDBWriteConection is used to reset the WriteConnection Pool (inmemory / OnDisk).
func resetDBWriteConection(dbFlag DbType) {
	switch dbFlag {
	case InMemory:
		if config.Data.DBConf.RedisHAEnabled {
			config := getInMemoryDBConfig()
			inMemDBConnPool.Mux.Lock()
			defer inMemDBConnPool.Mux.Unlock()
			if inMemDBConnPool.WritePool != nil {
				return
			}
			err := inMemDBConnPool.setWritePool(config)
			if err != nil {
				log.Error("Reset of inMemory write pool failed: " + err.Error())
				return
			}
			log.Info("New inMemory connection pool created")
		}
		return
	case OnDisk:
		if config.Data.DBConf.RedisHAEnabled {
			config := getOnDiskDBConfig()
			onDiskDBConnPool.Mux.Lock()
			defer onDiskDBConnPool.Mux.Unlock()
			if onDiskDBConnPool.WritePool != nil {
				return
			}
			err := onDiskDBConnPool.setWritePool(config)
			if err != nil {
				log.Error("Reset of onDisk write pool failed: " + err.Error())
				return
			}
			log.Info("New onDisk connection pool created")
		}
		return
	default:
		return
	}
}

func (p *ConnPool) setWritePool(config *Config) error {
	currentMasterIP, currentMasterPort := retryForMasterIP(p, config)
	if currentMasterIP == "" {
		return fmt.Errorf("unable to retrieve master ip from sentinel master election")
	}
	log.Info("new write pool master IP found: " + currentMasterIP)
	writePool, _ := getPool(currentMasterIP, currentMasterPort)
	if writePool == nil {
		return fmt.Errorf("write pool creation failed")
	}

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), unsafe.Pointer(writePool))
	p.MasterIP = currentMasterIP
	p.PoolUpdatedTime = time.Now()
	return nil
}

func retryForMasterIP(pool *ConnPool, config *Config) (currentMasterIP, currentMasterPort string) {
	for i := 0; i < 120; i++ {
		currentMasterIP, currentMasterPort = GetCurrentMasterHostPort(config)
		if currentMasterIP != "" && pool.MasterIP != currentMasterIP {
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
		SentinelPort: config.Data.DBConf.InMemorySentinelPort,
		MasterSet:    config.Data.DBConf.InMemoryMasterSet,
	}
}

func getOnDiskDBConfig() *Config {
	return &Config{
		Port:         config.Data.DBConf.OnDiskPort,
		Protocol:     config.Data.DBConf.Protocol,
		Host:         config.Data.DBConf.OnDiskHost,
		SentinelPort: config.Data.DBConf.OnDiskSentinelPort,
		MasterSet:    config.Data.DBConf.OnDiskMasterSet,
	}
}

//GetDBConnection is used to get the new Connection Pool for Inmemory/OnDisk DB
func GetDBConnection(dbFlag DbType) (*ConnPool, *errors.Error) {
	var err *errors.Error
	switch dbFlag {
	case InMemory:
		// In this case this function return in-memory db connection pool
		if inMemDBConnPool == nil || inMemDBConnPool.ReadPool == nil {
			log.Info("GetDBConnection : inMemDBConnPool OR inMemDBConnPool.ReadPool is nil")
			config := getInMemoryDBConfig()
			inMemDBConnPool, err = config.Connection()
			if err != nil {
				log.Error("error while trying to get Inmemory Readpool connection : " + err.Error())
			}
			inMemDBConnPool.PoolUpdatedTime = time.Now()
		}
		if inMemDBConnPool.WritePool == nil {
			log.Info("GetDBConnection : inMemDBConnPool.WritePool is nil, invoking resetDBWriteConection ")
			resetDBWriteConection(InMemory)
		}

		return inMemDBConnPool, err

	case OnDisk:
		// In this case this function returns On-Disk db connection pool
		if onDiskDBConnPool == nil || onDiskDBConnPool.ReadPool == nil {
			config := getOnDiskDBConfig()
			onDiskDBConnPool, err = config.Connection()
			if err != nil {
				log.Error("error while trying to get Ondisk Readpool connection : " + err.Error())
			}
			onDiskDBConnPool.PoolUpdatedTime = time.Now()
		}
		if onDiskDBConnPool.WritePool == nil {
			log.Info("GetDBConnection : onDiskDBConnPool.WritePool is nil, invoking resetDBWriteConection ")
			resetDBWriteConection(OnDisk)
		}
		return onDiskDBConnPool, err
	default:
		return nil, errors.PackError(errors.UndefinedErrorType, "error invalid db type selection")
	}
}

//getPool is used is utility function to get the Connection Pool from DB.
func getPool(host, port string) (*redis.Pool, error) {
	protocol := config.Data.DBConf.Protocol
	p := &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: config.Data.DBConf.MaxIdleConns,
		// max number of connections
		MaxActive: config.Data.DBConf.MaxActiveConns,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(protocol, host+":"+port)
			return c, err
		},
		/*TestOnBorrow is an optional application supplied function to
		  check the health of an idle connection before the connection is
		  used again by the application. Argument t is the time that the
		  connection was returned to the pool.This function PINGs
		  connections that have been idle more than a minute.
		  If the function returns an error, then the connection is closed.
		*/
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return p, nil
}

// Connection returns connection pool
// Connection does not take any input and returns a connection object used to interact with the DB
func (c *Config) Connection() (*ConnPool, *errors.Error) {
	var err error
	var masterIP string
	var masterPort string
	connPools := &ConnPool{}
	masterIP = c.Host
	masterPort = c.Port
	if config.Data.DBConf.RedisHAEnabled {
		masterIP, masterPort = GetCurrentMasterHostPort(c)
	}

	connPools.ReadPool, err = getPool(c.Host, c.Port)
	//Check if any connection error occured
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			log.Error("error while trying to get Readpool connection : " + errs.Error())
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, err)
	}
	connPools.WritePool, err = getPool(masterIP, masterPort)
	//Check if any connection error occured
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			log.Error("error while trying to get Writepool connection : " + errs.Error())
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, err)
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
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		log.Info("Create : WritePool nil")
		return errors.PackError(errors.UndefinedErrorType, "Create : WritePool is nil ")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()

	value, readErr := p.Read(table, key)
	if readErr != nil && readErr.ErrNo() == errors.DBConnFailed {
		return errors.PackError(readErr.ErrNo(), "error: db connection failed")
	}
	if value != "" {
		return errors.PackError(errors.DBKeyAlreadyExist, "error: data with key ", key, " already exists")
	}
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, "Write to DB in json form failed: "+err.Error())
	}
	_, createErr := writeConn.Do("SET", saveID, jsondata)
	if createErr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
		return errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
	}

	return nil
}

//Update data
/* Update take the following leys as input:
1."uid" is a string which acts as a unique ID to fetch the data from the DB
2."data" is userdata which is of type interface sent by the user to update/patch the already existing data
*/
func (p *ConnPool) Update(table, key string, data interface{}) (string, *errors.Error) {

	if _, readErr := p.Read(table, key); readErr != nil {
		if errors.DBKeyNotFound == readErr.ErrNo() {
			return "", errors.PackError(readErr.ErrNo(), "error: data with key ", key, " does not exist")
		}
		return "", readErr
	}
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		log.Error("Update : error in masrshalling json data", err.Error())
		return "", errors.PackError(errors.UndefinedErrorType, "Write to DB in json form failed: "+err.Error())
	}

	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		log.Info("Update : WritePool nil")
		return "", errors.PackError(errors.UndefinedErrorType, "Update : Writepool is nil ")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	_, createErr := writeConn.Do("SET", saveID, jsondata)
	if createErr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
		return "", errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
	}

	return saveID, nil
}

//Read is for getting singular data
// Read takes "key" sting as input which acts as a unique ID to fetch specific data from DB
func (p *ConnPool) Read(table, key string) (string, *errors.Error) {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	var (
		value interface{}
		err   error
	)

	value, err = readConn.Do("Get", table+":"+key)

	if err != nil {

		if err.Error() == "redigo: nil returned" {
			return "", errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
		}
		if errs, aye := isDbConnectError(err); aye {
			return "", errs
		}
		return "", errors.PackError(errors.DBKeyFetchFailed, errorCollectingData, err)
	}

	if value == nil {
		return "", errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
	}
	data, err := redis.String(value, err)
	if err != nil {
		return "", errors.PackError(errors.UndefinedErrorType, "error while trying to convert the data into string: ", err)
	}
	return string(data), nil
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

//GetAllDetails will fetch all the keys present in the database
func (p *ConnPool) GetAllDetails(table string) ([]string, *errors.Error) {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	keys, err := readConn.Do("KEYS", table+":*")
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	var IDs []string
	for _, data := range keys.([]interface{}) {
		key := string(data.([]uint8))
		ID := strings.TrimPrefix(key, table+":")
		IDs = append(IDs, ID)
	}
	return IDs, nil
}

//Delete data entry
// Read takes "key" sting as input which acts as a unique ID to delete specific data from DB
func (p *ConnPool) Delete(table, key string) *errors.Error {

	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		log.Info("Delete : WritePool nil")
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data: WritePool is nil ")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	_, readErr := p.Read(table, key)
	if readErr != nil {
		log.Error("Error while deleting data : " + readErr.Error())
		return readErr
	}

	_, doErr := writeConn.Do("DEL", table+":"+key)
	if doErr != nil {
		if errs, aye := isDbConnectError(doErr); aye {
			log.Error("Error while getting connection in Delete : " + errs.Error())
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data: ", doErr)
	}

	return nil
}

//CleanUpDB will delete all database entries
//The flush command will be executed without warnings please be cautious in using this
func (p *ConnPool) CleanUpDB() *errors.Error {
	writeConn := p.WritePool.Get()
	defer writeConn.Close()
	_, err := writeConn.Do("FLUSHALL")
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	return nil
}

/*
//FilterSearch to search resource with given filter
func (p *ConnPool) FilterSearch(table, key, path string) (interface{}, *errors.Error) {
    c := p.pool.Get()
    defer c.Close()
    rh := rejson.NewReJSONHandler()
    rh.SetRedigoClient(c)
    value, err := redis.Bytes(rh.JSONGet(table+":"+key, path))

    if err != nil {
        if errs, aye := isDbConnectError(err); aye {
            return "", errs
        }
        return "", errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
    }
    return value, nil
}

*/

//DeleteServer data entry without table
// Read takes "key" sting as input which acts as a unique ID to delete specific data from DB
func (p *ConnPool) DeleteServer(key string) *errors.Error {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	keys, err := readConn.Do("KEYS", key)
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			log.Error("Error while getting connection in Delete Server : " + errs.Error())
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	log.Info("Delete Server : WritePool value : ", writePool)
	if writePool == nil {
		log.Info("Delete Server: WritePool nil")
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data: WritePool is nil ")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	for _, data := range keys.([]interface{}) {
		delkey := string(data.([]uint8))
		_, err := writeConn.Do("DEL", delkey)
		if err != nil {
			if errs, aye := isDbConnectError(err); aye {
				log.Error("Error while Deleting Server : " + errs.Error())
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
				return errs
			}
			//          return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
		}
	}
	return nil
}

//GetAllMatchingDetails will fetch all the keys which matches pattern present in the database
func (p *ConnPool) GetAllMatchingDetails(table, pattern string) ([]string, *errors.Error) {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	keys, err := readConn.Do("KEYS", table+":*"+pattern+"*")
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	var IDs []string
	for _, data := range keys.([]interface{}) {
		key := string(data.([]uint8))
		ID := strings.TrimPrefix(key, table+":")
		IDs = append(IDs, ID)
	}
	return IDs, nil
}

//Transaction is to do a atomic operation using optimistic lock
func (p *ConnPool) Transaction(key string, cb func(string) error) *errors.Error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		return errors.PackError(errors.UndefinedErrorType, "error while trying to Write Transaction data: WritePool is nil")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	if _, err := writeConn.Do("WATCH", key); err != nil {
		if errs, aye := isDbConnectError(err); aye {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	writeConn.Send("MULTI")
	if err := cb(key); err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	_, err := writeConn.Do("EXEC")
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	/*
	   if queued != nil {
	       result = members[0]
	   }
	*/
	return nil
}

// isDbConnectError is for checking if error is dial connection error
func isDbConnectError(err error) (*errors.Error, bool) {
	if strings.HasSuffix(err.Error(), "connect: connection refused") || err.Error() == "EOF" {
		return errors.PackError(errors.DBConnFailed, err), true
	}
	return nil, false
}

//GetResourceDetails will fetch the key and also fetch the data
func (p *ConnPool) GetResourceDetails(key string) (string, *errors.Error) {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	keys, err := readConn.Do("KEYS", "*"+key)
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			log.Error("Error while getting resource datails  : " + errs.Error())
			return "", errs
		}
		return "", errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	var dkey string
	// keys array always 1
	for _, data := range keys.([]interface{}) {
		dkey = string(data.([]uint8))
	}
	if dkey == "" {
		return "", errors.PackError(errors.DBKeyNotFound, "no data with the with key ", key, " found")
	}

	params := strings.SplitN(dkey, ":", 2)
	return p.Read(params[0], params[1])
}

// AddResourceData will make an entry into the database with the given values
/* AddResourceData takes the following keys as input:
1."table" is a string which is used identify what kind of data we are storing.
2."data" is of type interface and is the userdata sent to be stored in DB.
3."key" is a string which acts as a unique ID to the data entry.
*/
func (p *ConnPool) AddResourceData(table, key string, data interface{}) *errors.Error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		return errors.PackError(errors.UndefinedErrorType, "WritePool is nil")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()

	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, "Write to DB in json form failed: "+err.Error())
	}
	_, createErr := writeConn.Do("SET", saveID, jsondata)
	if createErr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
		return errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
	}

	return nil
}

// Ping will check the DB connection health
func (p *ConnPool) Ping() error {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	if _, err := readConn.Do("PING"); err != nil {
		return fmt.Errorf("error while pinging DB with read connection")
	}
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	log.Info("Ping : WritePool value : ", writePool)
	if writePool == nil {
		log.Info("Ping : WritePool nil")
		return errors.PackError(errors.UndefinedErrorType, "error while pinging DB with read connection: WritePool is nil ")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	if _, err := writeConn.Do("PING"); err != nil {
		log.Error("error while pinging DB with write connection : " + err.Error())
		return errors.PackError(errors.UndefinedErrorType, "error while pinging DB with read connection: WritePool is nil "+err.Error())

	}
	return nil
}

// CreateIndex is used to create and save secondary index
/* CreateIndex take the following keys are input:
1. form is a map of the index to be created and the data along with it
2. uuid is the resource id with witch the value is stored
*/
func (p *ConnPool) CreateIndex(form map[string]interface{}, uuid string) error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	log.Info("CreateIndex : WritePool value : ", writePool)
	if writePool == nil {
		log.Info("CreateIndex : WritePool nil")
		return errors.PackError(errors.UndefinedErrorType, "error while Creating index: WritePool is nil ")

	}
	writeConn := writePool.Get()
	defer writeConn.Close()
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
		createErr := writeConn.Send("ZADD", index, val, key)
		if createErr != nil {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
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
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))

	if writePool == nil {
		log.Info("CreateTaskIndex : WritePool nil")
		return errors.PackError(errors.UndefinedErrorType, "error while creating task index: WritePool is nil ")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	createErr := writeConn.Send("ZADD", index, value, key)
	if createErr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
		return createErr
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
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	currentCursor := cursor
	match = strings.ToLower(match)
	for {
		d, getErr := readConn.Do("ZSCAN", index, currentCursor, "MATCH", match, "COUNT", count)
		if getErr != nil {
			return []string{}, fmt.Errorf("error while trying to get data: " + getErr.Error())
		}
		if len(d.([]interface{})) > 1 {
			data, err := redis.Strings(d.([]interface{})[1], getErr)
			if err != nil {
				return []string{}, fmt.Errorf("error while trying to get data: " + err.Error())
			}
			for i := 0; i < len(data); i++ {
				if data[i] != "0" {
					if regexFlag {
						getList = append(getList, data[i])
					} else {
						getList = append(getList, strings.Split(data[i], "::")[1])
					}
				}
			}
		}
		stringCursor := string(d.([]interface{})[0].([]uint8))
		if stringCursor == "0" {
			break
		}
		currentCursor, getErr = strconv.ParseFloat(stringCursor, 64)
		if getErr != nil {
			return []string{}, getErr
		}
	}
	return getList, nil
}

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
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	currentCursor := cursor
	for {
		d, getErr := readConn.Do("ZSCAN", index, currentCursor, "MATCH", "*", "COUNT", count)
		if getErr != nil {
			return nil, fmt.Errorf("error while trying to get data: " + getErr.Error())
		}
		if len(d.([]interface{})) > 1 {
			data, err := redis.Strings(d.([]interface{})[1], getErr)
			if err != nil {
				return nil, fmt.Errorf("error while trying to get data: " + err.Error())
			}
			for _, j := range data {
				if j != "0" {
					getList = append(getList, j)
				}
			}
		}
		stringCursor := string(d.([]interface{})[0].([]uint8))
		if stringCursor == "0" {
			break
		}
		currentCursor, getErr = strconv.ParseFloat(stringCursor, 64)
		if getErr != nil {
			return []string{}, getErr
		}
	}

	if regexFlag {
		return getList, nil
	}

	for _, k := range getList {
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
	storeList = getUniqueSlice(storeList)
	return storeList, nil
}

// GetRange is used to range over float type values
/*
1. index is the name of the index to search under
2. min is the minimum value for the search
3. max is the maximum value for the search
*/
func (p *ConnPool) GetRange(index string, min, max int, regexFlag bool) ([]string, error) {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	data, getErr := redis.Strings(readConn.Do("ZRANGEBYSCORE", index, min, max))
	if getErr != nil {
		return nil, fmt.Errorf("error while trying to get data: " + getErr.Error())
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
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	data, getErr := redis.Strings(readConn.Do("ZRANGE", index, min, max))
	if getErr != nil {
		return nil, fmt.Errorf("error while trying to get data: " + getErr.Error())
	}
	return data, nil
}

// Del is used to delete the index key
/*
1. index is the name of the index under which the key needs to be deleted
2. key is the id of the resource to be deleted under an index
*/
func (p *ConnPool) Del(index string, k string) error {
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	currentCursor := 0
	key := "*" + k
	for {
		d, getErr := readConn.Do("ZSCAN", index, currentCursor, "MATCH", key, "COUNT", count)
		if getErr != nil {
			return fmt.Errorf("error while trying to get data: " + getErr.Error())
		}
		if len(d.([]interface{})) > 1 {
			data, err := redis.Strings(d.([]interface{})[1], getErr)

			if err != nil {
				return fmt.Errorf("error while trying to get data: " + err.Error())
			}
			if len(data) < 1 {
				return fmt.Errorf("no data with ID found")
			}
			writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
			if writePool == nil {
				return fmt.Errorf("WritePool is nil")
			}
			writeConn := writePool.Get()
			defer writeConn.Close()
			for _, resource := range data {
				if resource != "0" {
					_, delErr := writeConn.Do("ZREM", index, resource)
					if delErr != nil {
						if errs, aye := isDbConnectError(delErr); aye {
							atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
							return errs
						}
						return fmt.Errorf("error while trying to delete data: " + delErr.Error())
					}
				}
			}
		}
		stringCursor := string(d.([]interface{})[0].([]uint8))
		if stringCursor == "0" {
			break
		}
		currentCursor, getErr = strconv.Atoi(stringCursor)
		if getErr != nil {
			return getErr
		}
	}

	return nil
}

// CreateEvtSubscriptionIndex is used to create and save secondary index
/* CreateSubscriptionIndex take the following keys are input:
1. index is the name of the index to be created
2. key and value are the key value pair for the index
*/
func (p *ConnPool) CreateEvtSubscriptionIndex(index string, key interface{}) error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		return fmt.Errorf("WritePool is nil")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	const value = 0

	matchKey := strings.Replace(key.(string), "[", "\\[", -1)
	matchKey = strings.Replace(matchKey, "]", "\\]", -1)
	val, _ := p.GetEvtSubscriptions(index, matchKey)
	if len(val) > 0 {
		return fmt.Errorf("Data Already Exist for the index: %v", index)
	}
	createErr := writeConn.Send("ZADD", index, value, key)
	if createErr != nil {
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
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
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	const cursor float64 = 0
	currentCursor := cursor

	for {
		d, getErr := readConn.Do("ZSCAN", index, currentCursor, "MATCH", searchKey, "COUNT", count)
		if getErr != nil {
			return []string{}, fmt.Errorf("error while trying to get data: " + getErr.Error())
		}
		if len(d.([]interface{})) > 1 {
			data, err := redis.Strings(d.([]interface{})[1], getErr)
			if err != nil {
				return []string{}, fmt.Errorf("error while trying to get data: " + err.Error())
			}
			for i := 0; i < len(data); i++ {
				if data[i] != "0" {
					getList = append(getList, data[i])
				}
			}
		}
		stringCursor := string(d.([]interface{})[0].([]uint8))
		if stringCursor == "0" {
			break
		}
		currentCursor, getErr = strconv.ParseFloat(stringCursor, 64)
		if getErr != nil {
			return []string{}, getErr
		}
	}
	return getList, nil
}

// DeleteEvtSubscriptions is for to Delete subscription details
// 1. index is the name of the index to be created
// 2. removeKey is string parameter for remove
func (p *ConnPool) DeleteEvtSubscriptions(index, removeKey string) error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		return fmt.Errorf("WritePool is nil")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()

	value, err := p.GetEvtSubscriptions(index, removeKey)
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf("No data found for the key: %v", removeKey)
	}
	for _, data := range value {
		delErr := writeConn.Send("ZREM", index, data)
		if delErr != nil {
			if errs, aye := isDbConnectError(delErr); aye {
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
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
		return fmt.Errorf("Error while updating subscriptions")
	}
	return nil
}

// CreateDeviceSubscriptionIndex is used to create and save secondary index
/* CreateDeviceSubscriptionIndex take the following keys are input:
1. index is the name of the index to be created
2. key is for the index
*/
func (p *ConnPool) CreateDeviceSubscriptionIndex(index, hostIP, location string, originResources []string) error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		return fmt.Errorf("WritePool is nil")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	const value = 0
	originResourceStr := "[" + strings.Join(originResources, " ") + "]"
	key := hostIP + "::" + location + "::" + originResourceStr
	// escape the square brackets before scanning
	searchKey := strings.Replace(key, "[", "\\[", -1)
	searchKey = strings.Replace(searchKey, "]", "\\]", -1)
	val, _ := p.GetDeviceSubscription(index, searchKey)
	if len(val) > 0 {
		return fmt.Errorf("Data Already Exist for the index: %v", index)
	}
	createErr := writeConn.Send("ZADD", index, value, key)
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
	var data []string
	readConn := p.ReadPool.Get()
	defer readConn.Close()
	const cursor float64 = 0
	currentCursor := cursor
	for {
		d, getErr := readConn.Do("ZSCAN", index, currentCursor, "MATCH", match, "COUNT", count)
		if getErr != nil {
			return nil, fmt.Errorf("error while trying to get data: " + getErr.Error())
		}
		if len(d.([]interface{})) > 1 {
			var err error
			data, err = redis.Strings(d.([]interface{})[1], getErr)
			if err != nil {
				return []string{}, err
			}
			log.Info("No of data records for get device subscription query : " + strconv.Itoa(len(data)))
			if len(data) < 1 {
				return []string{}, fmt.Errorf("No data found for the key: %v", match)
			}
			return data, nil
		}
		stringCursor := string(d.([]interface{})[0].([]uint8))
		if stringCursor == "0" {
			break
		}
		currentCursor, getErr = strconv.ParseFloat(stringCursor, 64)
		if getErr != nil {
			return []string{}, getErr
		}
	}
	return data, nil
}

// DeleteDeviceSubscription is for to Delete subscription details of Device
// 1. index is the name of the index to be created
// 2. removeKey is string parameter for remove
func (p *ConnPool) DeleteDeviceSubscription(index, hostIP string) error {
	writePool := (*redis.Pool)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool))))
	if writePool == nil {
		return fmt.Errorf("WritePool is nil")
	}
	writeConn := writePool.Get()
	defer writeConn.Close()
	value, err := p.GetDeviceSubscription(index, hostIP+"*")
	if err != nil {
		return err
	}

	if len(value) < 1 {
		return fmt.Errorf("No data found for the key: %v", hostIP)
	}
	for _, data := range value {
		delErr := writeConn.Send("ZREM", index, data)
		if delErr != nil {
			if errs, aye := isDbConnectError(delErr); aye {
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&p.WritePool)), nil)
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
		return fmt.Errorf("Error while updating subscriptions")
	}
	return nil
}

//UpdateResourceIndex is used to update the resource inforamtion which is indexed
// form contains index name and value:key for the index
func (p *ConnPool) UpdateResourceIndex(form map[string]interface{}, uuid string) error {
	for index := range form {
		err := p.Del(index, uuid)
		if (err != nil) && (err.Error() != "no data with ID found") {
			return fmt.Errorf("Error while updating index: %v", err)
		}
	}
	err := p.CreateIndex(form, uuid)
	if err != nil {
		return fmt.Errorf("Error while updating index: %v", err)
	}
	return nil
}
