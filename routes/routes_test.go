package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestAbout(t *testing.T) {
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, httptest.NewRequest("", "/about", nil))

	assert.Equal(t, w.Code, http.StatusOK)
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}

func TestLogOut(t *testing.T) {
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, httptest.NewRequest("", "/log-out", nil))

	assert.Equal(t, w.Code, http.StatusFound)
	assert.Equal(t, w.Header().Get("Location"), "/")
	assert.Equal(t, w.Header().Get("Set-Cookie"), "__Host-user=;MaxAge=0;Path=/;Secure")
}

func TestNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, httptest.NewRequest("", "/foo", nil))

	assert.Equal(t, w.Code, http.StatusNotFound)
	assert.Equal(t, w.Header().Get("Content-Type"), "text/html; charset=utf-8")
}

func TestRandom(t *testing.T) {
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, httptest.NewRequest("", "/random", nil))

	assert.Equal(t, w.Code, http.StatusFound)
}

func TestRobots(t *testing.T) {
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, httptest.NewRequest("", "/robots.txt", nil))

	assert.Equal(t, w.Code, http.StatusNoContent)
}
