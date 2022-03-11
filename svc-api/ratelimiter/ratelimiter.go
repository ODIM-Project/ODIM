package ratelimiter

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// Limit is max number of requests can be processed by ODIM at parallelly
// TODO: Limit should be set proper value, after scale test
const Limit int = 100 

// RateLimiter is for limiting the requests
// here we will check if count which is added db against the session token
// if count is exceded the limit then return the response with too many requests already in the systems
// if its not exceded then incremen the counter
// and when the request completes the task then decrement the counter
func RateLimiter(sessionToken string) error {
	if sessionToken != "" {
		count, _ := IncrementCounter(sessionToken)
		if count > Limit {
			DecrementCounter(sessionToken)
			errorMessage := "too many requests, retry after some time"
			log.Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
	}
	return nil
}

// IncrementCounter will increment the count
func IncrementCounter(sessionToken string) (int, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return 0, err
	}
	return conn.Incr("SessionRateLimit", sessionToken)
}

// DecrementCounter will decrement the count
func DecrementCounter(sessionToken string) (int, *errors.Error) {
	conn, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return 0, err
	}
	return conn.Decr("SessionRateLimit", sessionToken)
}
