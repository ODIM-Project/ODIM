package chassis

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

var internalError = common.GeneralError(http.StatusInternalServerError, response.InternalError, "error", nil, nil)
