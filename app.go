package main

import (
	"log"
	"net/http"
	"test/helper"
	"test/routes"

	"github.com/gorilla/handlers"
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Println("dispatchRequest :: panic :: ", err)
		}
	}()

	r := routes.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 50*1024*1024)
			next.ServeHTTP(w, r)
		})
	})

	origins := handlers.AllowedOrigins([]string{"https://penitenciarios.com", "https://www.penitenciarios.com", "https://dev.penitenciarios.com", "https://www.dev.penitenciarios.com", "http://localhost:3000", "http://127.0.0.1:3000"})

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})
	corsOptions := handlers.CORS(origins, allowedHeaders, allowedMethods)

	// updateData.UpdateData() // Funcion para actuaizar los datos de todos los usuarios

	config_vars := helper.GetConfigVars()
	http.ListenAndServe(":"+config_vars["LISTEN_PORT"], corsOptions(r))

}
