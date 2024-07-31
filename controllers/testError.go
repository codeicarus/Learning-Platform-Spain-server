package controllers

import "net/http"

func TestError(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/TestError.html")
}
