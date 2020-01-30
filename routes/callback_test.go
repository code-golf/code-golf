package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestCallbackNoCode(t *testing.T) {
	w := httptest.NewRecorder()

	Router.ServeHTTP(w, httptest.NewRequest("", "/callback", nil))

	assert.Equal(t, w.Code, http.StatusBadRequest)
}
