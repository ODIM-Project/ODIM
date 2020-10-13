package plugin

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Get(uri string) (Response, sresponse.Error) {
	args := c.Called(uri)

	var r Response
	if arg0 := args.Get(0); arg0 != nil {
		r = arg0.(Response)
	}

	var e sresponse.Error
	if arg1 := args.Get(1); arg1 != nil {
		e = arg1.(sresponse.Error)
	}
	return r, e
}

func (c *ClientMock) Post(uri string, body interface{}) (Response, sresponse.Error) {
	panic("implement me")
}

type ResponseMock struct {
	mock.Mock
}

func (r *ResponseMock) JSON(t interface{}) error {
	return r.Called(t).Error(0)
}

func (r *ResponseMock) AsRPCResponse() response.RPC {
	return r.Called().Get(0).(response.RPC)
}
