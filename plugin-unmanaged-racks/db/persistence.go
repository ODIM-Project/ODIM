package db

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
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
	return &ConnectionManager{&redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial(protocol, host+":"+port)
		},
	}}
}

type ConnectionManager struct {
	pool *redis.Pool
}

func (c *ConnectionManager) Delete(schema, key string) (bool, error) {
	cs := c.pool.Get()
	defer cs.Close()

	numOfRemovedKeys, err := cs.Do("DEL", schema+":"+key)
	if err != nil {
		return false, err
	}
	if numOfRemovedKeys.(int) == 0 {
		return false, fmt.Errorf("requested asset(%s) does not exist", schema+":"+key)
	}
	return true, nil
}

func (c *ConnectionManager) FindByKey(schema, key string) (interface{}, error) {
	cs := c.pool.Get()
	defer cs.Close()

	v, e := cs.Do("GET", schema+":"+key)
	if e != nil {
		return nil, e
	}
	if v == nil {
		return nil, fmt.Errorf("requested object(%s) does not exist", schema+":"+key)
	}
	return v, nil
}

func (c *ConnectionManager) Create(schema, key string, data []byte) *Error {
	cs := c.pool.Get()
	defer cs.Close()

	pk := strings.Title(schema) + ":" + key
	r, err := cs.Do("SETNX", pk, data)
	if err != nil {
		return &Error{DB_ERR_GENERAL, err.Error()}
	}

	v, ok := r.(int64)
	if !ok {
		return &Error{DB_ERR_GENERAL, "unexpected response type"}
	}

	switch v {
	case 0:
		return &Error{DB_ERR_ALREADY_EXISTS, fmt.Sprintf("specified key(%s) already exists", pk)}
	default:
		return nil
	}
}

func (c *ConnectionManager) Update(schema, key string, data []byte) error {
	cs := c.pool.Get()
	defer cs.Close()

	pk := strings.Title(schema) + ":" + key
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

	schema = formatSchemaSufix(schema)

	keys, err := cs.Do("KEYS", schema+"*")
	if err != nil {
		return nil, err
	}

	var result []string
	for _, key := range keys.([]interface{}) {
		result = append(result, strings.TrimPrefix(string(key.([]uint8)), schema))
	}
	return result, nil

}

func formatSchemaSufix(schema string) string {
	if strings.HasSuffix(schema, ":") {
		return schema
	} else {
		return schema + ":"
	}
}
