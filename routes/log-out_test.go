package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SeerUK/assert"
)

func TestLogOutGET(t *testing.T) {
	w := httptest.NewRecorder()

	logOutPost(w, httptest.NewRequest("POST", "/", nil))

	res := w.Result()
	res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusSeeOther)
	assert.Equal(t, res.Header.Get("Location"), "/")

	var cookies []string
	for _, cookie := range res.Cookies() {
		cookies = append(cookies, cookie.String())
	}

	assert.Equal(t, cookies, []string{
		"__Host-session=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Lax",
	})
}
