package ratelimiter

import (
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

const (

	// SessionRateLimit is table name to limit the number of requests per session
	SessionRateLimit = "SessionRateLimit"

	// UserRateLimit is table name to limit the number of sessions per user
	UserRateLimit = "UserRateLimit"
)

// RequestRateLimiter is for limiting number of requests per session
// here we will check if count which is added db against the session token
// if count is exceded the limit then return the response with too many requests, retry after some time
// if its not exceded then increment the counter
// and when the request completes the task then decrement the counter
func RequestRateLimiter(sessionToken string) error {
	if sessionToken != "" {
		count, _ := IncrementCounter(sessionToken, SessionRateLimit)
		if count > config.Data.RequestLimitCountPerSession {
			DecrementCounter(sessionToken, SessionRateLimit)
			errorMessage := "too many requests, retry after some time"
			log.Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
	}
	return nil
}

// SessionRateLimiter is for limiting number of sessions per user
// here we will check if count which is added db against the user id
// if count is exceded the limit then return the response with too many requests, retry after some time
// if its not exceded then incremen the counter
// and when the request completes the task then decrement the counter
func SessionRateLimiter(userid string) error {
	if userid != "" {
		count, _ := IncrementCounter(userid, UserRateLimit)
		if count > config.Data.SessionLimitCountPerUser {
			DecrementCounter(userid, UserRateLimit)
			errorMessage := "too many requests, retry after some time"
			log.Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
		//IncrementCounter(userid, UserRateLimit)
	}
	return nil
}

// IncrementCounter will increment the count
func IncrementCounter(key, table string) (int, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return 0, err
	}
	return conn.Incr(table, key)
}

// DecrementCounter will decrement the count
func DecrementCounter(key, table string) (int, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return 0, err
	}
	return conn.Decr(table, key)
}
