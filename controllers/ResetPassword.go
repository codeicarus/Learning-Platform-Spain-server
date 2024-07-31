package controllers

import (
	"net/http"
)

func Resetpassword(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/Resetpassword.html")
}
