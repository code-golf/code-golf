package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestLogOut(t *testing.T) {
	w := httptest.NewRecorder()

	LogOut(w, httptest.NewRequest("", "/log-out", nil))

	assert.Equal(t, w.Code, http.StatusFound)
	assert.Equal(t, w.Header().Get("Location"), "/")
	assert.Equal(t, w.Header().Get("Set-Cookie"), "__Host-user=;MaxAge=0;Path=/;Secure")
}
