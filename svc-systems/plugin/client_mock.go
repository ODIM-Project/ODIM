package plugin

import (
	"encoding/json"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Get(uri string) response.RPC {
	args := c.Called(uri)

	var r response.RPC
	if arg0 := args.Get(0); arg0 != nil {
		r = arg0.(response.RPC)
	}
	return r
}

func (c *ClientMock) Post(uri string, body *json.RawMessage) response.RPC {
	panic("implement me")
}

func (c *ClientMock) Delete(uri string) response.RPC {
	panic("implement me")
}

func (c *ClientMock) Patch(uri string, body *json.RawMessage) response.RPC {
	panic("implement me")
}
