package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestRobots(t *testing.T) {
	w := httptest.NewRecorder()

	Robots(w, httptest.NewRequest("", "/robots.txt", nil))

	assert.Equal(t, w.Code, http.StatusNoContent)
}
