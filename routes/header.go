package routes

import (
	"net/http"
	"net/url"

	"github.com/JRaspass/code-golf/cookie"
)

func printHeader(w http.ResponseWriter, r *http.Request, status int) int {
	header := w.Header()

	header["Content-Language"] = []string{"en"}
	header["Content-Type"] = []string{"text/html;charset=utf8"}
	header["Referrer-Policy"] = []string{"no-referrer"}
	header["X-Content-Type-Options"] = []string{"nosniff"}
	header["X-Frame-Options"] = []string{"DENY"}
	header["Content-Security-Policy"] = []string{
		"base-uri 'none';" +
			"connect-src 'self';" +
			"default-src 'none';" +
			"form-action 'none';" +
			"frame-ancestors 'none';" +
			"img-src 'self' data: avatars.githubusercontent.com;" +
			"script-src 'self';" +
			"style-src 'self'",
	}

	w.WriteHeader(status)

	var logInOrOut string

	userID, login := cookie.Read(r)

	if userID == 0 {
		logInOrOut = `"//github.com/login/oauth/authorize?client_id=7f6709819023e9215205&scope=user:email&redirect_uri=https://code-golf.io/callback?redirect_uri%3D` + url.QueryEscape(url.QueryEscape(r.RequestURI)) + `" id=login>Login with GitHub`
	} else {
		logInOrOut = `/logout id=logout title=Logout></a><a href="/users/` + login + `" id=me><img src="//avatars.githubusercontent.com/` + login + `?s=30">` + login
	}

	w.Write([]byte(
		"<!doctype html>" +
			"<link rel=stylesheet href=" + commonCssPath + ">" +
			`<meta name=description content="` +
			"Code Golf is a game designed to let you show off your code-fu " +
			`by solving problems in the least number of characters.">` +
			"<meta name=theme-color content=#222>" +
			`<meta name=viewport content="maximum-scale=1,user-scalable=0,width=device-width">` +
			"<title>Code-Golf</title>" +
			"<header><nav>" +
			"<a href=/>Home</a>" +
			"<a href=/about>About</a>" +
			"<a href=/scores>Scores</a>" +
			"<a href=/stats>Stats</a>" +
			"<a href=" + logInOrOut + "</a>" +
			"</nav></header>",
	))

	return userID
}
