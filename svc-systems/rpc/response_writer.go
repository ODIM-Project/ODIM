package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

func jsonResponseWriter(res http.Response, bodyTransformers ...func(string) string) (statusCode int32, body []byte, headers map[string]string) {
	jsonBodyBytes := new(json.RawMessage)
	err := json.NewDecoder(res.Body).Decode(jsonBodyBytes)
	if err != nil {
		ge := common.GeneralError(
			http.StatusInternalServerError,
			response.GeneralError,
			fmt.Sprintf("Cannot read response: %v", err), nil, nil)

		return ge.StatusCode, generateResponse(ge.Body), ge.Header
	}

	bodyToBeTransformed := string(*jsonBodyBytes)
	for _, t := range bodyTransformers {
		bodyToBeTransformed = t(bodyToBeTransformed)
	}

	headersToBeReturned := map[string]string{}
	for k := range res.Header {
		if hasToBeSkipped(k) {
			continue
		}
		for _, t := range bodyTransformers {
			headersToBeReturned[k] = t(res.Header.Get(k))
		}
	}

	return int32(res.StatusCode), []byte(bodyToBeTransformed), headersToBeReturned
}

func hasToBeSkipped(header string) bool {
	return header == "Content-Length"
}
