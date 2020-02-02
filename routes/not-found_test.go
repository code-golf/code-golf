package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	NotFound(w, httptest.NewRequest("", "/foo", nil))

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}
