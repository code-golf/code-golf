package routes

import (
	"net/http"

	"github.com/jraspass/code-golf/cookie"
)

func printHeader(w http.ResponseWriter, r *http.Request, status int) int {
	w.Header()["Content-Type"] = []string{"text/html;charset=utf8"}
	w.Header()["Content-Security-Policy"] = []string{
		"connect-src 'self';" +
			"default-src 'none';" +
			"img-src 'self' data: https://avatars.githubusercontent.com;" +
			"script-src 'self';" +
			"style-src 'self'",
	}

	w.WriteHeader(status)

	var logInOrOut string

	userID, login := cookie.Read(r)

	if userID == 0 {
		logInOrOut = `<a href="//github.com/login/oauth/authorize?client_id=7f6709819023e9215205&scope=user:email" id=login>Login with GitHub</a>`
	} else {
		logInOrOut = `<a href=/logout id=logout title=Logout></a><a href="/users/` + login + `" id=me><img src="//avatars.githubusercontent.com/` + login + `?s=30">` + login + "</a>"
	}

	w.Write([]byte(
		"<!doctype html>" +
			"<link rel=stylesheet href=" + commonCssPath + ">" +
			`<meta name=description content="` +
			"Code Golf is a game designed to let you show off your code-fu " +
			`by solving problems in the least number of characters.">` +
			"<meta name=theme-color content=#222>" +
			`<meta name=viewport content="width=device-width">` +
			"<title>Code-Golf</title>" +
			"<header><nav>" +
			"<a href=/>Home</a>" +
			"<a href=/about>About</a>" +
			"<a href=/scores/all/all>Scores</a>" +
			logInOrOut +
			"</nav></header>",
	))

	return userID
}
