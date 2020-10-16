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
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/gomodule/redigo/redis"
)

const (
	errorCollectingData string = "error while trying to collect data: "
	count               int    = 1000
)

// Connection returns connection pool
// Connection does not take any input and returns a connection object used to interact with the DB
func (c *Config) Connection() (*ConnPool, *errors.Error) {
	var err error
	p := &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: config.Data.DBConf.MaxIdleConns,
		// max number of connections
		MaxActive: config.Data.DBConf.MaxActiveConns,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(c.Protocol, c.Host+":"+c.Port)
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
	//Check if any connection error occured
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return nil, errs
		}
		return nil, errors.PackError(errors.UndefinedErrorType, err)
	}

	return &ConnPool{pool: p}, nil
}

// Create will make an entry into the database with the given values
/* Create takes the following keys as input:
1."table" is a string which is used identify what kind of data we are storing.
2."data" is of type interface and is the userdata sent to be stored in DB.
3."key" is a string which acts as a unique ID to the data entry.
*/
func (p *ConnPool) Create(table, key string, data interface{}) *errors.Error {
	conn := p.pool.Get()
	defer conn.Close()

	value, readErr := p.Read(table, key)
	if readErr.ErrNo() == errors.DBConnFailed  {
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
	_, createErr := conn.Do("SET", saveID, jsondata)
	if createErr != nil {
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
	conn := p.pool.Get()
	defer conn.Close()

	if _, readErr := p.Read(table, key); readErr != nil {
		if errors.DBKeyNotFound == readErr.ErrNo() {
			return "", errors.PackError(readErr.ErrNo(), "error: data with key ", key, " does not exist")
		}
		return "", readErr
	}
	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return "", errors.PackError(errors.UndefinedErrorType, "Write to DB in json form failed: "+err.Error())
	}
	_, createErr := conn.Do("SET", saveID, jsondata)
	if createErr != nil {
		return "", errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
	}

	return saveID, nil
}

//Read is for getting singular data
// Read takes "key" sting as input which acts as a unique ID to fetch specific data from DB
func (p *ConnPool) Read(table, key string) (string, *errors.Error) {
	c := p.pool.Get()
	defer c.Close()
	var (
		value interface{}
		err   error
	)

	value, err = c.Do("Get", table+":"+key)

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

//GetAllDetails will fetch all the keys present in the database
func (p *ConnPool) GetAllDetails(table string) ([]string, *errors.Error) {
	c := p.pool.Get()
	defer c.Close()
	keys, err := c.Do("KEYS", table+":*")
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
	c := p.pool.Get()
	defer c.Close()
	_, readErr := p.Read(table, key)
	if readErr != nil {
		return readErr
	}

	_, doErr := c.Do("DEL", table+":"+key)
	if doErr != nil {
		if errs, aye := isDbConnectError(doErr); aye {
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, "error while trying to delete data: ", doErr)
	}

	return nil
}

//CleanUpDB will delete all database entries
//The flush command will be executed without warnings please be cautious in using this
func (p *ConnPool) CleanUpDB() *errors.Error {
	c := p.pool.Get()
	defer c.Close()
	_, err := c.Do("FLUSHALL")
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
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
	c := p.pool.Get()
	defer c.Close()
	keys, err := c.Do("KEYS", key)
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
	}
	for _, data := range keys.([]interface{}) {
		delkey := string(data.([]uint8))
		_, err := c.Do("DEL", delkey)
		if err != nil {
			if errs, aye := isDbConnectError(err); aye {
				return errs
			}
			return errors.PackError(errors.UndefinedErrorType, errorCollectingData, err)
		}
	}
	return nil
}

//GetAllMatchingDetails will fetch all the keys which matches pattern present in the database
func (p *ConnPool) GetAllMatchingDetails(table, pattern string) ([]string, *errors.Error) {
	c := p.pool.Get()
	defer c.Close()
	keys, err := c.Do("KEYS", table+":*"+pattern+"*")
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
	c := p.pool.Get()
	defer c.Close()
	if _, err := c.Do("WATCH", key); err != nil {
		if errs, aye := isDbConnectError(err); aye {
			return errs
		}
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	c.Send("MULTI")
	if err := cb(key); err != nil {
		return errors.PackError(errors.UndefinedErrorType, err)
	}
	_, err := c.Do("EXEC")
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
	c := p.pool.Get()
	defer c.Close()
	keys, err := c.Do("KEYS", "*"+key)
	if err != nil {
		if errs, aye := isDbConnectError(err); aye {
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
	conn := p.pool.Get()
	defer conn.Close()

	saveID := table + ":" + key

	jsondata, err := json.Marshal(data)
	if err != nil {
		return errors.PackError(errors.UndefinedErrorType, "Write to DB in json form failed: "+err.Error())
	}
	_, createErr := conn.Do("SET", saveID, jsondata)
	if createErr != nil {
		return errors.PackError(errors.UndefinedErrorType, "Write to DB failed : "+createErr.Error())
	}

	return nil
}

// Ping will check the DB connection health
func (p *ConnPool) Ping() error {
	conn := p.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		return fmt.Errorf("error while pinging DB")
	}
	return nil
}

// CreateIndex is used to create and save secondary index
/* CreateIndex take the following keys are input:
1. form is a map of the index to be created and the data along with it
2. uuid is the resource id with witch the value is stored
*/
func (p *ConnPool) CreateIndex(form map[string]interface{}, uuid string) error {
	c := p.pool.Get()
	defer c.Close()
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
		createErr := c.Send("ZADD", index, val, key)
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
	c := p.pool.Get()
	defer c.Close()
	createErr := c.Send("ZADD", index, value, key)
	if createErr != nil {
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
	c := p.pool.Get()
	defer c.Close()
	currentCursor := cursor
	match = strings.ToLower(match)
	for {
		d, getErr := c.Do("ZSCAN", index, currentCursor, "MATCH", match, "COUNT", count)
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
	c := p.pool.Get()
	defer c.Close()
	currentCursor := cursor
	for {
		d, getErr := c.Do("ZSCAN", index, currentCursor, "MATCH", "*", "COUNT", count)
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
	c := p.pool.Get()
	defer c.Close()
	data, getErr := redis.Strings(c.Do("ZRANGEBYSCORE", index, min, max))
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
	c := p.pool.Get()
	defer c.Close()
	data, getErr := redis.Strings(c.Do("ZRANGE", index, min, max))
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
func (p *ConnPool) Del(index, key string) error {
	c := p.pool.Get()
	defer c.Close()
	k := "*" + key
	d, e := c.Do("ZSCAN", index, 0, "MATCH", k)
	if e != nil {
		return fmt.Errorf("error while trying to get data: " + e.Error())
	}
	if len(d.([]interface{})) > 1 {
		data, err := redis.Strings(d.([]interface{})[1], e)
		if err != nil {
			return fmt.Errorf("error while trying to get value of ID: " + err.Error())
		}
		if len(data) < 1 {
			return fmt.Errorf("no data with ID found")
		}
		for _, resource := range data {
			if resource != "0" {
				_, delErr := c.Do("ZREM", index, resource)
				if delErr != nil {
					return fmt.Errorf("error while trying to delete data: " + delErr.Error())
				}
			}
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
	c := p.pool.Get()
	defer c.Close()
	const value = 0
	val, _ := p.GetEvtSubscriptions(index, key.(string))
	if len(val) > 0 {
		return fmt.Errorf("Data Already Exist for the index: %v", index)
	}
	createErr := c.Send("ZADD", index, value, key)
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
	c := p.pool.Get()
	defer c.Close()
	const cursor float64 = 0
	currentCursor := cursor

	matchKey := strings.Replace(searchKey, "[", "\\[", -1)
	matchKey = strings.Replace(matchKey, "]", "\\]", -1)

	for {
		d, getErr := c.Do("ZSCAN", index, currentCursor, "MATCH", matchKey, "COUNT", count)
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
	c := p.pool.Get()
	defer c.Close()

	value, err := p.GetEvtSubscriptions(index, removeKey)
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf("No data found for the key: %v", removeKey)
	}
	for _, data := range value {
		c.Send("ZREM", index, data)
	}
	return nil
}

// UpdateEvtSubscriptions is for to Update subscription details
// 1. index is the name of the index to be created
// 2. key and value are the key value pair for the index
func (p *ConnPool) UpdateEvtSubscriptions(index, subscritionID string, key interface{}) error {
	c := p.pool.Get()
	defer c.Close()

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
	c := p.pool.Get()
	defer c.Close()
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
	createErr := c.Send("ZADD", index, value, key)
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
	c := p.pool.Get()
	defer c.Close()
	const cursor float64 = 0
	currentCursor := cursor
	for {
		d, getErr := c.Do("ZSCAN", index, currentCursor, "MATCH", match, "COUNT", count)
		if getErr != nil {
			return nil, fmt.Errorf("error while trying to get data: " + getErr.Error())
		}
		if len(d.([]interface{})) > 1 {
			var err error
			data, err = redis.Strings(d.([]interface{})[1], getErr)
			if err != nil {
				return []string{}, err
			}
			log.Println("No of data records for get device subscription query : ", len(data))
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
	c := p.pool.Get()
	defer c.Close()
	value, err := p.GetDeviceSubscription(index, hostIP+"*")
	if err != nil {
		return err
	}
	if len(value) < 1 {
		return fmt.Errorf("No data found for the key: %v", hostIP)
	}
	for _, data := range value {
		c.Send("ZREM", index, data)
	}
	return nil
}

// UpdateDeviceSubscription is for to Update subscription details
// 1. index is the name of the index to be created
// 2. key and value are the key value pair for the index
func (p *ConnPool) UpdateDeviceSubscription(index, hostIP, location string, originResources []string) error {
	c := p.pool.Get()
	defer c.Close()
	_, err := p.GetDeviceSubscription(index, hostIP+"*")
	if err != nil {
		return err
	}
	// host ip will be unique on each index in subscription of device
	// so there will be only one data
	err = p.DeleteDeviceSubscription(index, hostIP)
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
