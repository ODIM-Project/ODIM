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
	sut := NewBasicAuthHandler(TEST_CONFIG.UserName, TEST_CONFIG.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusUnauthorized)
}

func Test_provided_username_is_invalid(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	req.SetBasicAuth("invalid-user", "Od!m12$4")

	rw := httptest.NewRecorder()
	sut := NewBasicAuthHandler(TEST_CONFIG.UserName, TEST_CONFIG.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusUnauthorized)
}

func Test_provided_password_is_invalid(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	req.SetBasicAuth("admin", "invalid-password")

	rw := httptest.NewRecorder()
	sut := NewBasicAuthHandler(TEST_CONFIG.UserName, TEST_CONFIG.Password)

	Do(rw, req, sut)
	require.Equal(t, rw.Result().StatusCode, http.StatusUnauthorized)
}

func Test_provided_credentials_are_valid(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	req.SetBasicAuth("admin", "Od!m12$4")

	rw := httptest.NewRecorder()
	sut := NewBasicAuthHandler(TEST_CONFIG.UserName, TEST_CONFIG.Password)

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
