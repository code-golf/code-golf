package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomGET(t *testing.T) {
	w := httptest.NewRecorder()

	randomGET(w, httptest.NewRequest("", "/random", nil))

	assert.Equal(t, w.Code, http.StatusFound)
}
