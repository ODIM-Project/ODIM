//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package sresponse ...
package sresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

//CommonError struct definition
type CommonError struct {
	Error ErrorClass `json:"error"`
}

//ErrorClass struct definition
type ErrorClass struct {
	Code                string            `json:"code"`
	Message             string            `json:"message"`
	MessageExtendedInfo []MsgExtendedInfo `json:"@Message.ExtendedInfo"`
}

//MsgExtendedInfo struct definition
type MsgExtendedInfo struct {
	OdataType  string `json:"@odata.type"`
	MessageID  string `json:"MessageId"`
	Message    string `json:"Message"`
	Severity   string `json:"Severity"`
	Resolution string `json:"Resolution"`
}

type Error interface {
	AsRPCResponse() response.RPC
}

type RPCErrorWrapper struct {
	response.RPC
}

func (r *RPCErrorWrapper) AsRPCResponse() response.RPC {
	return r.RPC
}

type UnknownErrorWrapper struct {
	Error      error
	StatusCode int
}

func (e *UnknownErrorWrapper) AsRPCResponse() response.RPC {
	return common.GeneralError(int32(e.StatusCode), response.InternalError, e.Error.Error(), nil, nil)
}
