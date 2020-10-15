package chassis

import (
	"net/http"
	"testing"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewGetHandler(t *testing.T) {
	managedChassis := dmtf.Chassis{}

	sut := NewGetHandler(
		nil,
		func(table, key string, r interface{}) *errors.Error {
			r = managedChassis
			return nil
		},
	)

	getChassisRPCRequest := chassisproto.GetChassisRequest{}
	r := sut.Handle(&getChassisRPCRequest)
	require.EqualValues(t, http.StatusOK, r.StatusCode)
	require.Equal(t, managedChassis, r.Body)
}

func TestNewGetHandler_WhenManagedChassisFinderRespondsWithError(t *testing.T) {
	sut := NewGetHandler(
		nil,
		func(table, key string, r interface{}) *errors.Error {
			return errors.PackError(errors.JSONUnmarshalFailed, "error")
		},
	)

	r := sut.Handle(&chassisproto.GetChassisRequest{})
	require.EqualValues(t, http.StatusInternalServerError, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func TestNewGetHandler_WhenPluginClientFactoryReturnsNotFoundError(t *testing.T) {
	sut := NewGetHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.DBKeyNotFound, "urp plugin not found")
		},
		func(table, key string, r interface{}) *errors.Error {
			return errors.PackError(errors.DBKeyNotFound, "there is no managed chassis with specified key")
		})

	r := sut.Handle(&chassisproto.GetChassisRequest{})
	require.EqualValues(t, http.StatusNotFound, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
	require.EqualValues(t, response.ResourceNotFound, r.Body.(response.CommonError).Error.MessageExtendedInfo[0].MessageID)
}

func TestNewGetHandler_WhenPluginClientFactoryReturnsUnexpectedError(t *testing.T) {
	sut := NewGetHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.InvalidAuthToken, "urp plugin not found")
		},
		func(table, key string, r interface{}) *errors.Error {
			return errors.PackError(errors.DBKeyNotFound, "there is no managed chassis with specified key")
		})

	r := sut.Handle(&chassisproto.GetChassisRequest{})
	require.EqualValues(t, http.StatusInternalServerError, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
	require.EqualValues(t, response.InternalError, r.Body.(response.CommonError).Error.MessageExtendedInfo[0].MessageID)
}

func TestNewGetHandler_WhenPluginClientReturnsError(t *testing.T) {
	ppc := new(plugin.ClientMock)
	ppc.On("Get", mock.AnythingOfType("string")).
		Return(internalError)
	sut := NewGetHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return ppc, nil
		},
		func(table, key string, r interface{}) *errors.Error {
			return errors.PackError(errors.DBKeyNotFound, "there is no managed chassis with specified key")
		})

	r := sut.Handle(&chassisproto.GetChassisRequest{})
	require.EqualValues(t, http.StatusInternalServerError, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
	require.EqualValues(t, response.InternalError, r.Body.(response.CommonError).Error.MessageExtendedInfo[0].MessageID)
}

func TestNewGetHandler_WhenPluginClientReturnsNonErrorResponse(t *testing.T) {
	ppc := new(plugin.ClientMock)
	ppc.On("Get", mock.AnythingOfType("string")).
		Return(
			response.RPC{
				StatusCode: http.StatusOK,
				Body:       dmtf.Chassis{},
			},
		)

	sut := NewGetHandler(
		func(name string) (plugin.Client, *errors.Error) {
			return ppc, nil
		},
		func(table, key string, r interface{}) *errors.Error {
			return errors.PackError(errors.DBKeyNotFound, "there is no managed chassis with specified key")
		})

	r := sut.Handle(&chassisproto.GetChassisRequest{})
	require.EqualValues(t, http.StatusOK, r.StatusCode)
	require.IsType(t, dmtf.Chassis{}, r.Body)
}
