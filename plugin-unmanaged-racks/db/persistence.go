package db

import (
	"fmt"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var ErrAlreadyExists = redis.Error("already exists")

func NewConnectionManager(protocol, host, port string) *ConnectionManager {
	return &ConnectionManager{
		pool: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial(protocol, host+":"+port)
			},
		},
	}
}

type ConnectionManager struct {
	pool *redis.Pool
}

func (c *ConnectionManager) GetConnection() redis.Conn {
	return c.pool.Get()
}

func (c *ConnectionManager) DoInTransaction(callback func(c redis.Conn) error, syncKeys ...string) error {
	conn := c.GetConnection()
	defer NewConnectionCloser(&conn)

	for _, sk := range syncKeys {
		_, err := conn.Do("WATCH", sk)
		if err != nil {
			return err
		}
	}
	err := conn.Send("MULTI")
	if err != nil {
		return err
	}
	err = callback(conn)
	if err != nil {
		return err
	}
	r, err := conn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("cannot commit transaction: %v", err)
	}
	if r == redis.ErrNil {
		return fmt.Errorf("transaction aborted for unknown reason: %v", err)
	}
	return nil
}

type Key string

func (k Key) String() string {
	return string(k)
}

func (k Key) Prefix() string {
	return string(k) + ":"
}

func (k Key) TrimWildcard() Key {
	return Key(strings.TrimSuffix(string(k), "*"))
}

func (k Key) WithWildcard() Key {
	return k + "*"
}

func CreateContainsKey(tokens ...string) Key {
	return CreateKey(append([]string{"CONTAINS"}, tokens...)...)
}

func CreateContainedInKey(tokens ...string) Key {
	return CreateKey(append([]string{"CONTAINEDIN"}, tokens...)...)
}

func CreateKey(keys ...string) Key {
	return Key(strings.Join(keys, ":"))
}

func NewConnectionCloser(conn *redis.Conn) {
	func() {
		err := (*conn).Close()
		if err != nil {
			logging.Errorf("Cannot close DB connection: %s", err)
		}
	}()
}
