package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Parallel()
	request := httptest.NewRequest("GET", "/ping", nil)
	response := httptest.NewRecorder()

	Ping(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "application/json", response.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"message":"pong"}`, response.Body.String())
}

func TestHealthcheck(t *testing.T) {
	t.Parallel()
	request := httptest.NewRequest("GET", "/health", nil)
	response := httptest.NewRecorder()

	Healthcheck(response, request)

	assert.Equal(t, http.StatusNoContent, response.Code)
	assert.Empty(t, response.Body.String(), "Healthcheck should not return any body")
}
