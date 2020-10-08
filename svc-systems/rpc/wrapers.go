package rpc

import (
	"context"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"

	"github.com/micro/go-micro/metadata"
)

func auth(ctx context.Context, callback func() response.RPC) response.RPC {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "X-Auth-Token header is missing", nil, nil)
	}
	sessionToken, ok := md["X-Auth-Token"]
	if !ok {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "X-Auth-Token header is missing", nil, nil)
	}

	status, msg := services.IsAuthorized(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if status == http.StatusOK {
		return callback()
	}
	return common.GeneralError(status, response.NoValidSession, msg, nil, nil)
}
