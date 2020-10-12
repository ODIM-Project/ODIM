package rpc

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

func auth(sessionToken string, callback func() response.RPC) response.RPC {
	if sessionToken == "" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "X-Auth-Token header is missing", nil, nil)
	}

	status, msg := services.IsAuthorized(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if status == http.StatusOK {
		return callback()
	}
	return common.GeneralError(status, response.NoValidSession, msg, nil, nil)
}
