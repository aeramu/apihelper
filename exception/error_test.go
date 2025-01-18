package exception_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aeramu/apihelper/exception"
	"github.com/aeramu/apihelper/httphelper"
	"github.com/stretchr/testify/assert"
)

func TestSoftError_HTTP(t *testing.T) {
	err := exception.SoftError("TEST_CODE", "message", nil)
	var httpErr httphelper.HTTPError

	errors.As(err, &httpErr)

	assert.Equal(t, "TEST_CODE", httpErr.Code())
	assert.Equal(t, "message", httpErr.Message())
	assert.Equal(t, http.StatusOK, httpErr.HTTPStatus())
}

func TestInvalidArgumentError(t *testing.T) {
	err := exception.InvalidRequest("TEST_CODE", "message", nil)
	var httpErr httphelper.HTTPError

	errors.As(err, &httpErr)

	assert.Equal(t, "TEST_CODE", httpErr.Code())
	assert.Equal(t, "message", httpErr.Message())
	assert.Equal(t, http.StatusBadRequest, httpErr.HTTPStatus())
}

func TestWrap(t *testing.T) {
	repo := func() error {
		return exception.ErrorNotFound
	}

	repo2 := func() error {
		return exception.Wrap(repo(), "message")
	}

	assert.True(t, errors.Is(repo(), exception.ErrorNotFound))
	assert.True(t, errors.Is(repo2(), exception.ErrorNotFound))
}
