package chassis

import (
	"fmt"
	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func Test_GetCollectionHandler_WhenMultipleSourcesAreAvailable(t *testing.T) {
	source1 := new(sourceMock)
	source1.On("read").Return([]dmtfmodel.Link{{"1"}, {"3"}}, nil)
	source2 := new(sourceMock)
	source2.On("read").Return([]dmtfmodel.Link{{"2"}, {"4"}}, nil)

	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{source1, source2}, nil)
	sut := GetCollectionHandler{cspMock}

	r := sut.Handle()
	require.EqualValues(t, http.StatusOK, r.StatusCode)
	require.IsType(t, sresponse.NewChassisCollection(), r.Body)
	require.Equal(t, []dmtfmodel.Link{{"1"}, {"3"}, {"2"}, {"4"}}, r.Body.(sresponse.Collection).Members)
	require.Equal(t, map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}, r.Header)
}

func Test_GetCollectionHandler_WhenCollectionSourcesCannotBeDetermined(t *testing.T) {
	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{}, &sresponse.UnknownErrorWrapper{fmt.Errorf("error"), 500})
	sut := GetCollectionHandler{cspMock}

	r := sut.Handle()
	require.NotEqual(t, http.StatusOK, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func Test_GetCollectionHandler_WhenFirstSourceReturnsError(t *testing.T) {
	source1 := new(sourceMock)
	source1.On("read").Return([]dmtfmodel.Link{}, &sresponse.UnknownErrorWrapper{Error: fmt.Errorf("error"), StatusCode: 500})
	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{source1}, nil)
	sut := GetCollectionHandler{cspMock}

	r := sut.Handle()
	require.NotEqual(t, http.StatusOK, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func Test_GetCollectionHandler_WhenNonFirstSourceReturnsError(t *testing.T) {
	source1 := new(sourceMock)
	source1.On("read").Return([]dmtfmodel.Link{{"1"}}, nil)

	source2 := new(sourceMock)
	source2.On("read").Return([]dmtfmodel.Link{}, &sresponse.UnknownErrorWrapper{Error: fmt.Errorf("error"), StatusCode: 500})
	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{source1, source2}, nil)
	sut := GetCollectionHandler{cspMock}

	r := sut.Handle()
	require.NotEqual(t, http.StatusOK, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func Test_collectionSourceProvider_whenURPIsNotRegistered(t *testing.T) {
	sut := sourceProviderImpl{
		getPluginConfig: func(pluginID string) (smodel.Plugin, *errors.Error) {
			return smodel.Plugin{}, errors.PackError(errors.DBKeyNotFound, "plugin not found")
		},
	}

	r, e := sut.findSources()
	require.Nil(t, e)
	require.Len(t, r, 1)
	require.IsType(t, &managedChassisProvider{}, r[0])
}

func Test_collectionSourceProvider_whenURPIsRegistered(t *testing.T) {
	sut := sourceProviderImpl{
		getPluginConfig: func(pluginID string) (smodel.Plugin, *errors.Error) {
			return smodel.Plugin{}, nil
		},
	}

	r, e := sut.findSources()
	require.Nil(t, e)
	require.Len(t, r, 2)
	require.IsType(t, &managedChassisProvider{}, r[0])
	require.IsType(t, &unmanagedChassisProvider{}, r[1])
}

func Test_collectionSourceProvider_whenURPIsRegisteredAndUnderlyingDBReturnsError(t *testing.T) {
	sut := sourceProviderImpl{
		getPluginConfig: func(pluginID string) (smodel.Plugin, *errors.Error) {
			return smodel.Plugin{}, errors.PackError(errors.UndefinedErrorType, "unexpected error")
		},
	}

	_, e := sut.findSources()
	require.NotNil(t, e)
}

func Test_managedChassisProvider_WhenUnderlyingDBReturnsError(t *testing.T) {
	sut := managedChassisProvider{
		func(table string) ([]string, error) {
			return nil, fmt.Errorf("error")
		},
	}

	_, e := sut.read()
	require.NotNil(t, e)
}

func Test_managedChassisProvider_WhenUnderlyingDBReturnsNoKeys(t *testing.T) {
	sut := managedChassisProvider{
		func(table string) ([]string, error) {
			return []string{}, nil
		},
	}

	r, e := sut.read()
	require.Nil(t, e)
	require.Len(t, r, 0)
}

func Test_managedChassisProvider_WhenUnderlyingDBReturnsSomeKeys(t *testing.T) {
	sut := managedChassisProvider{
		func(table string) ([]string, error) {
			return []string{
				"first", "second", "third",
			}, nil
		},
	}

	r, e := sut.read()
	require.Nil(t, e)
	require.Len(t, r, 3)
	require.Equal(t, []dmtfmodel.Link{
		{Oid: "first"}, {"second"}, {"third"},
	}, r)
}

type collectionSourceProviderMock struct {
	mock.Mock
}

func (c *collectionSourceProviderMock) findSources() ([]source, sresponse.Error) {
	args := c.Mock.Called()
	return args.Get(0).([]source), getErrorOrNil(args.Get(1))
}

func getErrorOrNil(a interface{}) sresponse.Error {
	if a == nil {
		return nil
	}
	return a.(sresponse.Error)
}

type sourceMock struct {
	mock.Mock
}

func (s *sourceMock) read() ([]dmtfmodel.Link, sresponse.Error) {
	args := s.Mock.Called()
	return args.Get(0).([]dmtfmodel.Link), getErrorOrNil(args.Get(1))
}
