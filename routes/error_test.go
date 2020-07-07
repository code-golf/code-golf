package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestForbidden(t *testing.T) {
	w := httptest.NewRecorder()

	Forbidden(w, httptest.NewRequest("", "/", nil))

	assert.Equal(t, w.Code, http.StatusForbidden)
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	NotFound(w, httptest.NewRequest("", "/", nil))

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}
