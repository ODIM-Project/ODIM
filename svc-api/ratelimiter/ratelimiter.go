package ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	iris "github.com/kataras/iris/v12"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

const (

	// SessionRateLimit is table name to limit the number of requests per session
	SessionRateLimit = "SessionRateLimit"

	// UserRateLimit is table name to limit the number of sessions per user
	UserRateLimit = "UserRateLimit"

	// ResourceRateLimit is table name to limit the resource on time bound
	ResourceRateLimit = "ResourceRateLimit"
)

// RequestRateLimiter is for limiting number of requests per session
// here we will check if count which is added db against the session token
// if count is exceded the limit then return the response with too many requests, retry after some time
// if its not exceded then increment the counter
// and when the request completes the task then decrement the counter
func RequestRateLimiter(ctx context.Context, sessionToken string) error {
	if sessionToken != "" {
		count, _ := IncrementCounter(sessionToken, SessionRateLimit)
		if count > config.Data.RequestLimitCountPerSession {
			DecrementCounter(sessionToken, SessionRateLimit)
			errorMessage := "too many requests, retry after some time"
			l.LogWithFields(ctx).Error(errorMessage)
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
func SessionRateLimiter(ctx context.Context, userid string) error {
	if userid != "" {
		count, _ := IncrementCounter(userid, UserRateLimit)
		if count > config.Data.SessionLimitCountPerUser {
			DecrementCounter(userid, UserRateLimit)
			errorMessage := "too many requests, retry after some time"
			l.LogWithFields(ctx).Error(errorMessage)
			return fmt.Errorf(errorMessage)
		}
		//IncrementCounter(userid, UserRateLimit)
	}
	return nil
}

// ResourceRateLimiter will Limit the get on resource untill previous get completed the task
func ResourceRateLimiter(ctx iris.Context) {
	ctxt := ctx.Request().Context()
	uri := ctx.Request().RequestURI
	for _, val := range config.Data.ResourceRateLimit {
		resourceLimit := strings.Split(val, ":")
		if len(resourceLimit) > 1 && resourceLimit[1] != "" {
			rLimit, _ := strconv.Atoi(resourceLimit[1])
			resource := strings.Replace(resourceLimit[0], "{id}", "[a-zA-Z0-9._-]+", -1)
			regex := regexp.MustCompile(resource)
			if regex.MatchString(uri) {
				conn, err := common.GetDBConnection(common.InMemory)
				if err != nil {
					l.LogWithFields(ctxt).Error(err.Error())
					response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
					common.SetResponseHeader(ctx, response.Header)
					ctx.StatusCode(http.StatusInternalServerError)
					ctx.JSON(&response.Body)
					return
				}
				// convert millisecond to second
				expiretime := rLimit / 1000
				if err = conn.SetExpire("ResourceRateLimit", uri, "", expiretime); err != nil {
					errorMessage := "too many requests, retry after some time"
					l.LogWithFields(ctxt).Error(errorMessage)
					response := common.GeneralError(http.StatusServiceUnavailable, response.RateLimitExceeded, errorMessage, nil, nil)
					remainTime, _ := conn.TTL(ResourceRateLimit, uri)
					if remainTime > 0 {
						ctx.ResponseWriter().Header().Set("Retry-After", strconv.Itoa(remainTime))
					}
					common.SetResponseHeader(ctx, response.Header)
					ctx.StatusCode(http.StatusServiceUnavailable)
					ctx.JSON(&response.Body)
					return
				}
			}
		}
	}
	ctx.Next()
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
