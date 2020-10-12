package rpc

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

type authenticator func(sessionToken string, privileges, oemPrivileges []string) (int32, string)

func auth(authenticate authenticator, sessionToken string, callback func() response.RPC) response.RPC {
	if sessionToken == "" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "X-Auth-Token header is missing", nil, nil)
	}

	status, msg := authenticate(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if status == http.StatusOK {
		return callback()
	}
	return common.GeneralError(status, response.NoValidSession, msg, nil, nil)
}
