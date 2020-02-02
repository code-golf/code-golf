package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestRandom(t *testing.T) {
	w := httptest.NewRecorder()

	Random(w, httptest.NewRequest("", "/random", nil))

	assert.Equal(t, w.Code, http.StatusFound)
}
