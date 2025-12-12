package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallbackGETNoCode(t *testing.T) {
	w := httptest.NewRecorder()

	callbackGET(w, httptest.NewRequest("", "/callback", nil))

	assert.Equal(t, w.Code, http.StatusBadRequest)
}
