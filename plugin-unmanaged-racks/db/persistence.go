package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const (
	DB_ERR_ALREADY_EXISTS DBError = iota + 1
	DB_ERR_GENERAL
)

type DBError int

type Error struct {
	Code DBError
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

func NewConnectionManager(protocol, host, port string) *ConnectionManager {
	return &ConnectionManager{
		pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial(protocol, host+":"+port)
			},
		},
		domain: "URP",
	}
}

type ConnectionManager struct {
	pool   *redis.Pool
	domain string
}

func (c *ConnectionManager) Delete(schema, key string) (bool, error) {
	cs := c.pool.Get()
	defer cs.Close()

	numOfRemovedKeys, err := cs.Do("DEL", schema+":"+key)
	if err != nil {
		return false, err
	}
	if numOfRemovedKeys.(int64) == 0 {
		return false, nil
	}
	return true, nil
}

// Returns object related with requested key(schema + key) or nil in case when requested key does not exist
func (c *ConnectionManager) FindByKey(key Key) (interface{}, error) {
	cs := c.pool.Get()
	defer cs.Close()

	return cs.Do("GET", key)
}

func (c *ConnectionManager) GetConnection() redis.Conn {
	return c.pool.Get()
}

type transactional func(c redis.Conn)

func (c *ConnectionManager) DoInTransaction(syncKey string, callback transactional) error {
	conn := c.GetConnection()
	defer conn.Close()

	_, err := conn.Do("WATCH", syncKey)
	if err != nil {
		return err
	}
	err = conn.Send("MULTI")
	if err != nil {
		return err
	}
	callback(conn)
	r, err := conn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %v", err)
	}
	if r == redis.ErrNil {
		return fmt.Errorf("transaction aborted for unknown reason: %v", err)
	}
	return nil
}

func (c *ConnectionManager) Create(key Key, data []byte) *Error {
	cs := c.pool.Get()
	defer cs.Close()

	r, err := cs.Do("SETNX", key, data)
	if err != nil {
		return &Error{DB_ERR_GENERAL, err.Error()}
	}

	v, ok := r.(int64)
	if !ok {
		return &Error{DB_ERR_GENERAL, "unexpected response type"}
	}

	switch v {
	case 0:
		return &Error{DB_ERR_ALREADY_EXISTS, fmt.Sprintf("specified key(%s) already exists", key)}
	default:
		return nil
	}
}

func (c *ConnectionManager) Update(schema, key string, data []byte) error {
	cs := c.pool.Get()
	defer cs.Close()

	pk := c.CreateKey(c.domain, schema, key)
	s, err := redis.String(cs.Do("SET", pk, data))
	if err != nil {
		return err
	}

	switch s {
	case "-1":
		return fmt.Errorf("specified key(%s) does not exists", pk)
	default:
		return nil
	}
}

func (c *ConnectionManager) GetAllKeys(schema string) ([]string, error) {
	cs := c.pool.Get()
	defer cs.Close()

	patternKey := c.CreateKey(schema)
	keys, err := cs.Do("KEYS", patternKey+"*")
	if err != nil {
		return nil, err
	}

	var result []string
	for _, key := range keys.([]interface{}) {
		result = append(result, strings.TrimPrefix(string(key.([]uint8)), string(patternKey)))
	}
	return result, nil

}

type Key string

func (c *ConnectionManager) CreateChassisContainsKey(chassisOid string) Key {
	return c.CreateKey("CONTAINS", chassisOid)
}

func (c *ConnectionManager) CreateContainedInKey(oid string) Key {
	return c.CreateKey("CONTAINEDIN", oid)
}

func (c *ConnectionManager) CreateKey(keys ...string) Key {
	return Key(strings.Join(append([]string{c.domain}, keys...), ":"))
}

func NewConnectionCloser(conn *redis.Conn) func() {
	return func() {
		err := (*conn).Close()
		if err != nil {
			log.Print("Error: ", err)
		}
	}
}
