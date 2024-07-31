package controllers

import "net/http"

func TestSuccess(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/TestSuccess.html")
}
