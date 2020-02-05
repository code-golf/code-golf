package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestAbout(t *testing.T) {
	w := httptest.NewRecorder()

	About(w, httptest.NewRequest("", "/about", nil))

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}
