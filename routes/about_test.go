package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestAbout(t *testing.T) {
	w := httptest.NewRecorder()

	about(w, httptest.NewRequest("", "/", nil), nil)

	assert.Equal(t, w.Code, http.StatusOK)
}
