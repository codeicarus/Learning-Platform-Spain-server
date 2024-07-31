package controllers

import (
	"net/http"
)

func SendEmail(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/verifyEmail.html")
}
