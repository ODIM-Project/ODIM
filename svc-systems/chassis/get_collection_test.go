/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chassis

import (
	"fmt"
	"net/http"
	"testing"

	dmtfmodel "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_GetCollectionHandler_WhenMultipleSourcesAreAvailable(t *testing.T) {
	source1 := new(sourceMock)
	source1.On("read").Return([]dmtfmodel.Link{dmtfmodel.Link{Oid: "1"}, dmtfmodel.Link{Oid: "3"}}, nil)
	source2 := new(sourceMock)
	source2.On("read").Return([]dmtfmodel.Link{{Oid: "2"}, {Oid: "4"}}, nil)

	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{source1, source2}, nil)
	sut := GetCollection{cspMock}

	r := sut.Handle()
	require.EqualValues(t, http.StatusOK, r.StatusCode)
	require.IsType(t, sresponse.NewChassisCollection(), r.Body)
	require.Equal(t, []dmtfmodel.Link{{Oid: "1"}, {Oid: "3"}, {Oid: "2"}, {Oid: "4"}, {Oid: "5"}}, r.Body.(sresponse.Collection).Members)
}

func Test_GetCollectionHandler_WhenCollectionSourcesCannotBeDetermined(t *testing.T) {
	cspMock := new(collectionSourceProviderMock)

	cspMock.On("findSources").Return([]source{}, &internalError)
	sut := GetCollection{cspMock}

	r := sut.Handle()
	require.NotEqual(t, http.StatusOK, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func Test_GetCollectionHandler_WhenFirstSourceReturnsError(t *testing.T) {
	source1 := new(sourceMock)
	source1.On("read").Return([]dmtfmodel.Link{}, &internalError)
	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{source1}, nil)
	sut := GetCollection{cspMock}

	r := sut.Handle()
	require.NotEqual(t, http.StatusOK, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func Test_GetCollectionHandler_WhenNonFirstSourceReturnsError(t *testing.T) {
	source1 := new(sourceMock)
	source1.On("read").Return([]dmtfmodel.Link{{Oid: "1"}}, nil)

	source2 := new(sourceMock)
	source2.On("read").Return([]dmtfmodel.Link{}, &internalError)
	cspMock := new(collectionSourceProviderMock)
	cspMock.On("findSources").Return([]source{source1, source2}, nil)
	sut := GetCollection{cspMock}

	r := sut.Handle()
	require.NotEqual(t, http.StatusOK, r.StatusCode)
	require.IsType(t, response.CommonError{}, r.Body)
}

func Test_collectionSourceProvider_whenURPIsNotRegistered(t *testing.T) {
	sut := sourceProviderImpl{
		pluginClientFactory: func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.DBKeyNotFound, "plugin not found")
		},
	}

	r, e := sut.findSources()
	require.Nil(t, e)
	require.Len(t, r, 1)
	require.IsType(t, &managedChassisProvider{}, r[0])
}

func Test_collectionSourceProvider_whenURPIsRegistered(t *testing.T) {
	cm := new(plugin.ClientMock)
	sut := sourceProviderImpl{
		pluginClientFactory: func(name string) (plugin.Client, *errors.Error) {
			return cm, nil
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
		pluginClientFactory: func(name string) (plugin.Client, *errors.Error) {
			return nil, errors.PackError(errors.UndefinedErrorType, "unexpected error")
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
		{Oid: "first"}, {Oid: "second"}, {Oid: "third"},
	}, r)
}

type collectionSourceProviderMock struct {
	mock.Mock
}

func (c *collectionSourceProviderMock) findSources() ([]source, *response.RPC) {
	args := c.Mock.Called()
	return args.Get(0).([]source), getErrorOrNil(args.Get(1))
}

func (c *collectionSourceProviderMock) findFabricChassis(col *sresponse.Collection) {
	link := dmtfmodel.Link{
		Oid: "5",
	}
	col.AddMember(link)
}

func getErrorOrNil(a interface{}) *response.RPC {
	if a == nil {
		return nil
	}
	return a.(*response.RPC)
}

type sourceMock struct {
	mock.Mock
}

func (s *sourceMock) read() ([]dmtfmodel.Link, *response.RPC) {
	args := s.Mock.Called()
	return args.Get(0).([]dmtfmodel.Link), getErrorOrNil(args.Get(1))
}

func TestNewGetCollectionHandler(t *testing.T) {
	config := config.URLTranslation{NorthBoundURL: map[string]string{
		"ODIM": "redfish",
	},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		}}
	fun := func(table string) ([]string, error) {
		return []string{}, nil
	}
	NewGetCollectionHandler(plugin.NewClientFactory(&config), fun)
}
