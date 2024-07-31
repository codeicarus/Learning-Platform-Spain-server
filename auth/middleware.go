package auth

import (
	"fmt"
	"net/http"
	"time"

	"test/helper"
	"test/models"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(user *models.Usuarios) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Email":    user.Email,
		"Password": user.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte("secret"))
}

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("authorization")
		if tokenStr == "" {
			helper.ResponseWithJson(w, http.StatusUnauthorized,
				helper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
		} else {
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					helper.ResponseWithJson(w, http.StatusUnauthorized,
						helper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
					return nil, fmt.Errorf("not authorization")
				}
				return []byte("secret"), nil
			})
			if !token.Valid {
				helper.ResponseWithJson(w, http.StatusUnauthorized,
					helper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
			} else {
				next.ServeHTTP(w, r)
			}
		}
	})
}
