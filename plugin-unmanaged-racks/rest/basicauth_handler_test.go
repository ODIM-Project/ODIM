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

package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/stretchr/testify/require"
)

func Test_authorization_header_is_missing(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)

	rw := httptest.NewRecorder()
	sut := newBasicAuthHandler(testConfig.UserName, testConfig.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusUnauthorized)
}

func Test_provided_username_is_invalid(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	req.SetBasicAuth("invalid-user", "Od!m12$4")

	rw := httptest.NewRecorder()
	sut := newBasicAuthHandler(testConfig.UserName, testConfig.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusUnauthorized)
}

func Test_provided_password_is_invalid(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	req.SetBasicAuth("admin", "invalid-password")

	rw := httptest.NewRecorder()
	sut := newBasicAuthHandler(testConfig.UserName, testConfig.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusUnauthorized)
}

func Test_provided_credentials_are_valid(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	req.SetBasicAuth("admin", "Od!m12$4")

	rw := httptest.NewRecorder()
	sut := newBasicAuthHandler(testConfig.UserName, testConfig.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusOK)
}

func Do(w http.ResponseWriter, r *http.Request, handler iris.Handler) {
	app := iris.New()
	app.ContextPool = context.New(func() context.Context {
		return context.NewContext(app)
	})

	ctx := app.ContextPool.Acquire(w, r)
	handler(ctx)
	app.ContextPool.Release(ctx)
}
