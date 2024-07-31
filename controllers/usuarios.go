package controllers

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"test/helper"
	"test/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	daoUsuarios      = models.Usuarios{}
	daoAcademias     = models.Academias{}
	daoTransacciones = models.Transacciones{}
)

func FindUsuario(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var params models.FindUsuario
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	result, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {
		helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"result": "success", "usuario": result})
	}
}
func LoginUsuario(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var params models.Login
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	// result, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))

	// if err != nil {
	// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario o Contraseña incorrectos"})
	// 	return
	// } else {

	// 	if result.Estado == "VERIFICADO" {

	// 		if helper.GetMD5Hash(params.Password) == result.Password {

	// 			if result.Connected == "online" || result.Connected == "pending" {
	// 				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario conectado en otro dispositivo"})
	// 				return
	// 			}

	// 			// if result.LastIp == params.Ip {
	// 			// 	result.LastHeartBeat = helper.MakeTimestamp()
	// 			// 	result.LastIp = params.Ip
	// 			// 	if err := daoUsuarios.UpdateUsuario(result); err != nil {
	// 			// 		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
	// 			// 		return
	// 			// 	}
	// 			// } else {
	// 			// 	if result.LastHeartBeat+60 < helper.MakeTimestamp() {
	// 			// 		result.LastHeartBeat = helper.MakeTimestamp()
	// 			// 		result.LastIp = params.Ip
	// 			// 		if err := daoUsuarios.UpdateUsuario(result); err != nil {
	// 			// 			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
	// 			// 			return
	// 			// 		}
	// 			// 	} else {

	// 			// texto_email := "El usuario ##NOMBRE## (##EMAIL##) ha producido una conexion doble "
	// 			// texto_email = strings.ReplaceAll(texto_email, "##NOMBRE##", result.Name)
	// 			// texto_email = strings.ReplaceAll(texto_email, "##EMAIL##", result.Email)

	// 			// if !helper.SendEmail("PENITENCIARIOS.COM :: Error conexión doble", texto_email, "penitenciarios@penitenciarios.com", "penitenciarios@penitenciarios.com") {
	// 			// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
	// 			// 	return
	// 			// }

	// 			// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Parece que estas conectado a este usuario desdo otro equipo. Si no es asi, espere 2 minutos y vuelva a intentarlo"})
	// 			// 	return
	// 			// }
	// 			// }

	// 			token, _ := auth.GenerateToken(&result)
	// 			helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"result": "success", "token": token, "first_login": result.FirstLogin, "email": result.Email, "maximo_dia_suscripcion": result.MaximoDiaSuscripcion})
	// 		} else {
	// 			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
	// 			return
	// 		}
	// 	} else {
	// 		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no verificado"})
	// 		return
	// 	}
	// }
}

func RegisterUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var datos_usuario = map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&datos_usuario); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	var usuario models.Usuarios

	usuario.Name = datos_usuario["nombre"]
	usuario.Email = datos_usuario["email"]
	usuario.Password = datos_usuario["password"]
	usuario.Estado = "PENDIENTE"
	usuario.FirstLogin = false
	usuario.IdNivel = bson.ObjectIdHex("5ea6ad3dbb5c000045007637")
	usuario.LastFrase = 0
	usuario.CodigoPromocional = datos_usuario["codigo_promocional"]

	if datos_usuario["n_letter"] == "1" {
		usuario.Newsletter = true
	} else {
		usuario.Newsletter = false
	}

	dia_maximo_actual := time.Now().Unix()
	dia_maximo_actual_date := time.Unix(dia_maximo_actual, 0)
	after := dia_maximo_actual_date.AddDate(0, 0, 3)
	usuario.MaximoDiaSuscripcion = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()

	// usuario.MaximoDiaSuscripcion = 1689476399 // 15 de julio de 2023

	if usuario.CodigoPromocional != "" {
		// _, err := daoAcademias.FindAcademiaByCodigo(usuario.CodigoPromocional)
		// if err != nil {
		// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Código promocional no válido"})
		// 	return
		// }
	}
	// _, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(usuario.Email)))

	// if err != nil {
	// 	usuario.Id = bson.NewObjectId()
	// 	if err := daoUsuarios.InsertUsuario(usuario); err != nil {
	// 		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
	// 		return
	// 	}
	// 	config_vars := helper.GetConfigVars()

	// 	res, err := http.Get(config_vars["URL"] + "email/validate_email")
	// 	// res, err := http.Get(config_vars["URL"] + "EMAILS/validate_email/index.html")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	texto_email_bytes, err := ioutil.ReadAll(res.Body)
	// 	res.Body.Close()
	// 	texto_email := string(texto_email_bytes)
	// 	texto_email = strings.ReplaceAll(texto_email, "##NOMBRE##", usuario.Name)

	// 	texto_email = strings.ReplaceAll(texto_email, "##LINK##", config_vars["URLF"]+"verify/"+helper.GetMD5Hash("VERIFICAREMAIL"+strings.ToLower(strings.TrimSpace(usuario.Email))))

	// 	// if !helper.SendEmail("PENITENCIARIOS.COM :: Verificación de email", texto_email, usuario.Email, "penitenciarios@penitenciarios.com") {
	// 	// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
	// 	// 	return
	// 	// }
	// 	// helper.ResponseWithJson(w, http.StatusCreated, usuario)
	// } else {
	// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Ya existe un usuario con ese email en nuestra BBDD"})
	// 	return
	// }

}

func CheckMDS(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))

	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {
		helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"result": "success", "penitenciarios_mds": usuario.MaximoDiaSuscripcion})
	}
}

func VerifyUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Verify
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByVerificador(params.Verificador)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Url de verificación erronea"})
		return
	} else {
		usuario.Estado = "VERIFICADO"
		if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	}
}

func VerifyUsuarioAcount(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.VerifyAcountUser
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioById(params.Id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "usuario no encontrado"})
		return
	}
	usuario.Estado = "VERIFICADO"
	if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteUserAcount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := daoUsuarios.RemoveUsuario(id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func SetNivelUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.SetNivelUsuario
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(params.Email)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {
		usuario.IdNivel = params.IdNivel
		usuario.FirstLogin = true
		usuario.Detalle = params.Preguntas

		if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	}
}
func UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Usuarios
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	if err := daoUsuarios.UpdateUsuario(params); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func CheckUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	tokenStr := r.Header.Get("authorization")

	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			helper.ResponseWithJson(w, http.StatusUnauthorized,
				helper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
			return nil, fmt.Errorf("not authorization")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Acceder al campo "Email" del token
		email := claims["Email"].(string)

		// Mostrar el email en la consola

		helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success", "email": email})
	} else {
		helper.ResponseWithJson(w, http.StatusUnauthorized,
			helper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
	}
}

func CheckIpUsuario(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {

		if usuario.LastIp == params["ip"] {

			if params["token"] == usuario.Token {
				usuario.LastHeartBeat = helper.MakeTimestamp()

				if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
				return
			} else {
				// if usuario.LastHeartBeat+40 > helper.MakeTimestamp() {
				// 	log.Println("Conexion del usuario es mas reciente")
				// 	usuario.LastHeartBeat = helper.MakeTimestamp()
				// 	usuario.LastIp = params["ip"]
				// 	usuario.Token = params["token"]

				// 	if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
				// 		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				// 		return
				// 	}
				// 	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
				// 	return
				// } else {

				// log.Println("El usuario tiene token distinto")
				// log.Println(usuario.LastHeartBeat+30 > helper.MakeTimestamp())
				// log.Println("Ultima conexion +30: ", usuario.LastHeartBeat+30)
				// log.Println("tiempo actual: ", helper.MakeTimestamp())
				// helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Te hemos deslogueado porque parece que estas conectado a este usuario desdo otro equipo. Si no es asi, espere 2 minutos y vuelva a intentarlo"})
				// return

				// }
			}

			usuario.LastHeartBeat = helper.MakeTimestamp()
			usuario.LastIp = params["ip"]
			usuario.Token = params["token"]

			if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
			helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
			return
		} else {
			log.Println("No son iguales las IP")
			log.Println("el tiempo actual: ", helper.MakeTimestamp())
			log.Println("La antigua conexion +30: ", usuario.LastHeartBeat+30)
			log.Println("La ultima conexion es mayor a la nueva: ", usuario.LastHeartBeat+30 < helper.MakeTimestamp())

			if usuario.LastHeartBeat+40 < helper.MakeTimestamp() {
				usuario.LastHeartBeat = helper.MakeTimestamp()
				usuario.LastIp = params["ip"]
				if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
				return
			} else {

				texto_email := "El usuario ##NOMBRE## (##EMAIL##) ha producido una conexion doble "
				texto_email = strings.ReplaceAll(texto_email, "##NOMBRE##", usuario.Name)
				texto_email = strings.ReplaceAll(texto_email, "##EMAIL##", usuario.Email)

				// if !helper.SendEmail("PENITENCIARIOS.COM :: Error conexión doble", texto_email, "penitenciarios@penitenciarios.com", "penitenciarios@penitenciarios.com") {
				// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
				// 	return
				// }
				// helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Te hemos deslogueado porque parece que estas conectado a este usuario desdo otro equipo. Si no es asi, espere 2 minutos y vuelva a intentarlo"})
				// return
			}
		}
	}
}
func CheckSuperUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	if helper.CheckUser(params["email"]) {
		helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
	} else {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {

		const charset = "abcdefghijklmnopqrstuvwxyz" +
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

		var seededRand *rand.Rand = rand.New(
			rand.NewSource(time.Now().UnixNano()))

		b := make([]byte, 10)
		for i := range b {
			b[i] = charset[seededRand.Intn(len(charset))]
		}

		var new_pass = string(b)

		usuario.Password = helper.GetMD5Hash(new_pass)
		if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		config_vars := helper.GetConfigVars()
		// res, err := http.Get(config_vars["URL"] + "EMAILS/reset_password/index.html")
		res, err := http.Get(config_vars["URL"] + "email/reset_password")
		if err != nil {
			log.Fatal(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##NOMBRE##", usuario.Name)
		texto_email = strings.ReplaceAll(texto_email, "##CONTRASENA##", new_pass)

		texto_email = strings.ReplaceAll(texto_email, "##LINK##", config_vars["URL"]+"login")

		// if !helper.SendEmail("PENITENCIARIOS.COM :: Resetear password", texto_email, usuario.Email, "penitenciarios@penitenciarios.com") {
		// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
		// 	return
		// }
		// helper.ResponseWithJson(w, http.StatusCreated, usuario)

	}
}

func SetNewPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
	}

	usuario, err := daoUsuarios.FindUsuarioById(params["id"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	}

	usuario.Password = params["new_password"]

	if err := daoUsuarios.UpdateUsuarioPassword(usuario); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problemas al actualizar la contraseña"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func Contact(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	// texto_email := "Nombre: " + params["nombre"] + "<br>Email: " + params["email"] + "<br>Mensaje: " + params["mensaje"]

	// if !helper.SendEmail("PENITENCIARIOS.COM :: Formulario de contacto", texto_email, "penitenciarios@penitenciarios.com", params["email"]) {
	// 	helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
	// 	return
	// }
	// helper.ResponseWithJson(w, http.StatusCreated, "")
}
func ChangeMisDatos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email_user_logued"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {

		if strings.ToLower(strings.TrimSpace(params["email_user_logued"])) == strings.ToLower(strings.TrimSpace(params["email"])) {

		} else {
			_, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
			if err != nil {
			} else {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Ya existe un usuario con ese email"})
				return
			}
		}

		usuario.Name = strings.TrimSpace(params["nombre"])
		usuario.Email = strings.ToLower(strings.TrimSpace(params["email"]))
		if params["password"] != "" {
			usuario.Password = params["password"]
		}
		usuario.IdNivel = bson.ObjectIdHex(params["id_nivel"])

		if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		helper.ResponseWithJson(w, http.StatusCreated, "")
	}
}
func CheckTransacction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	config_vars := helper.GetConfigVars()

	var PASSWORD = config_vars["PAYCOMET_PASSWORD"]
	var NotificationHash_Calculed = sha512.Sum512([]byte(r.Form.Get("AccountCode") + r.Form.Get("TpvID") + r.Form.Get("TransactionType") + r.Form.Get("Order") + r.Form.Get("Amount") + r.Form.Get("Currency") + helper.GetMD5Hash(PASSWORD) + r.Form.Get("BankDateTime") + r.Form.Get("Response")))
	var NotificationHash_Calculed_string = fmt.Sprintf("%x", NotificationHash_Calculed)

	var id_tran, err = strconv.ParseInt(r.Form.Get("Order"), 10, 32)

	transaccion, err := daoTransacciones.FindTransaccionByIdTransaccion(id_tran)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	}
	if transaccion.Estado != 0 {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Transaccion ya validada"})
		return
	}
	if r.Form.Get("NotificationHash") == NotificationHash_Calculed_string {

		if r.Form.Get("Response") == "KO" {
			transaccion.Estado = -1
		} else {
			transaccion.Estado = 1
			transaccion.FechaCobro = helper.MakeTimestamp()
			transaccion.Token = r.Form.Get("Signature")
		}
	} else {
		transaccion.Estado = -2
	}

	if err := daoTransacciones.UpdateTransaccion(transaccion); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	if transaccion.Estado == 1 {

		usuario, err := daoUsuarios.FindUsuarioById(transaccion.IdUsuario.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
			return
		}

		var dia_maximo int64 = 0
		dia_maximo_actual := time.Now().Unix()

		if dia_maximo_actual < usuario.MaximoDiaSuscripcion {
			dia_maximo_actual = usuario.MaximoDiaSuscripcion
		}

		dia_maximo_actual_date := time.Unix(dia_maximo_actual, 0)

		if transaccion.Producto == "promocion_salida" {
			// dia_maximo = 2500000000
			dia_maximo = 1688180403 // 1 de julio de 2023
		} else if transaccion.Producto == "suscripcion_0990" {
			after := dia_maximo_actual_date.AddDate(0, 1, 0)
			dia_maximo = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()
		} else if transaccion.Producto == "suscripcion_2790" {
			after := dia_maximo_actual_date.AddDate(0, 3, 0)
			dia_maximo = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()
		} else if transaccion.Producto == "suscripcion_5590" {
			after := dia_maximo_actual_date.AddDate(0, 6, 0)
			dia_maximo = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()
		} else if transaccion.Producto == "suscripcion_9990" {
			after := dia_maximo_actual_date.AddDate(1, 0, 0)
			dia_maximo = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()
		}

		usuario.MaximoDiaSuscripcion = dia_maximo

		if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}
	}

	helper.ResponseWithJson(w, http.StatusCreated, "")

}
func CreateTransacction(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	config_vars := helper.GetConfigVars()

	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {

		var lastIdTransaccion = int64(0)

		lastTransaccion, err := daoTransacciones.FindTransaccionLastTransaction()
		if err != nil {
		} else {
			lastIdTransaccion = lastTransaccion.IdTransaccion + 1
		}

		var precio = 9999.99
		var descripcion_producto = ""
		var descriptor_producto = ""
		if params["producto"] == "promocion_salida" {
			// precio = 39.90
			precio = 0
			descripcion_producto = "Promoción exclusiva de salida penitenciarios.com"
			descriptor_producto = "PENITENCIARIOS.COM - PROMOCION SALIDA"
		} else if params["producto"] == "suscripcion_0990" {
			precio = 9.90
			descripcion_producto = "Suscripción 1 mes penitenciarios.com"
			descriptor_producto = "PENITENCIARIOS.COM - 1 MES"
		} else if params["producto"] == "suscripcion_2790" {
			precio = 27.90
			descripcion_producto = "Suscripción 3 mes penitenciarios.com"
			descriptor_producto = "PENITENCIARIOS.COM - 3 MES"
		} else if params["producto"] == "suscripcion_5590" {
			precio = 55.90
			descripcion_producto = "Suscripción 6 mes penitenciarios.com"
			descriptor_producto = "PENITENCIARIOS.COM - 6 MES"
		} else if params["producto"] == "suscripcion_9990" {
			precio = 99.90
			descripcion_producto = "Suscripción 12 mes penitenciarios.com"
			descriptor_producto = "PENITENCIARIOS.COM - 12 MES"
		}

		if helper.CheckUser(params["email"]) {
			precio = 0.01
		}

		var transaccion models.Transacciones
		transaccion.Id = bson.NewObjectId()
		transaccion.IdTransaccion = lastIdTransaccion
		transaccion.Fecha = helper.MakeTimestamp()
		transaccion.Estado = 0
		transaccion.Producto = params["producto"]
		transaccion.Precio = precio
		transaccion.FechaCobro = 0
		transaccion.Token = ""
		transaccion.IdUsuario = usuario.Id

		if err := daoTransacciones.InsertTransaccion(transaccion); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var MERCHANT_MERCHANTCODE = config_vars["PAYCOMET_MERCHANTCODE"]
		var MERCHANT_TERMINAL = config_vars["PAYCOMET_TERMINAL"]
		var OPERATION = config_vars["PAYCOMET_OPERATION"]
		var LANGUAGE = config_vars["PAYCOMET_LANGUAGE"]
		var MERCHANT_ORDER = fmt.Sprintf("%012d", transaccion.IdTransaccion)
		var MERCHANT_AMOUNT = strconv.Itoa(int(precio * 100))
		var MERCHANT_CURRENCY = config_vars["PAYCOMET_CURRENCY"]
		var PASSWORD = config_vars["PAYCOMET_PASSWORD"]
		var MERCHANT_MERCHANTSIGNATURE = sha512.Sum512([]byte(MERCHANT_MERCHANTCODE + MERCHANT_TERMINAL + OPERATION + MERCHANT_ORDER + MERCHANT_AMOUNT + MERCHANT_CURRENCY + helper.GetMD5Hash(PASSWORD)))

		var url = config_vars["PAYCOMET_URL"]
		url += "?MERCHANT_MERCHANTCODE=" + MERCHANT_MERCHANTCODE
		url += "&MERCHANT_TERMINAL=" + MERCHANT_TERMINAL
		url += "&OPERATION=" + OPERATION
		url += "&LANGUAGE=" + LANGUAGE
		url += "&MERCHANT_MERCHANTSIGNATURE=" + fmt.Sprintf("%x", MERCHANT_MERCHANTSIGNATURE)
		url += "&MERCHANT_ORDER=" + MERCHANT_ORDER
		url += "&MERCHANT_AMOUNT=" + MERCHANT_AMOUNT
		url += "&MERCHANT_CURRENCY=" + MERCHANT_CURRENCY
		url += "&MERCHANT_PRODUCTDESCRIPTION=" + descripcion_producto
		url += "&MERCHANT_DESCRIPTOR=" + descriptor_producto

		helper.ResponseWithJson(w, http.StatusCreated, map[string]string{"url": url})
	}
}

func StatsUsuario(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {
		historiales, err := daoHistoriales.FindHistorialByIdUsuario(usuario.Id.Hex(), nil)
		historiales_cp, err := daoHistorialesCP.FindHistorialByIdUsuario(usuario.Id.Hex(), nil)

		var historiales_preguntas_array []models.HistorialesPipe
		var ids_preguntas []bson.ObjectId
		var ids_historiales []bson.ObjectId
		var preguntas_devolver = make(map[string]models.Preguntas)

		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
			return
		} else {
			for _, historial := range historiales {
				ids_historiales = append(ids_historiales, historial.Id)
			}

			historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdsHistoriales(ids_historiales)
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
				return
			}
			for _, historial_pregunta := range historiales_preguntas {
				ids_preguntas = append(ids_preguntas, historial_pregunta.IdPregunta)
			}

			for _, historial := range historiales {

				var historiales_preguntas_de_este_examen []models.HistorialesPreguntas

				for _, historial_pregunta := range historiales_preguntas {
					if historial.Id == historial_pregunta.IdHistorial {
						historiales_preguntas_de_este_examen = append(historiales_preguntas_de_este_examen, historial_pregunta)
					}
				}

				var params models.HistorialesPipe
				params.Id = historial.Id
				params.Tipo = historial.Tipo
				params.Fecha = historial.Fecha
				params.Temas = historial.Temas
				params.Legislaciones = historial.Legislaciones
				params.BasicosConfundidos = historial.BasicosConfundidos
				params.Oficiales = historial.Oficiales
				params.Simulacros = historial.Simulacros
				params.IdUsuario = historial.IdUsuario
				params.NumeroPreguntas = historial.NumeroPreguntas
				params.RespuestaAutomatica = historial.RespuestaAutomatica
				params.TiempoTranscurrido = historial.TiempoTranscurrido
				params.Terminado = historial.Terminado
				params.PreguntasTotales = historial.PreguntasTotales
				params.PreguntasAcertadas = historial.PreguntasAcertadas
				params.PreguntasFalladas = historial.PreguntasFalladas
				params.PreguntasBlancas = historial.PreguntasBlancas
				params.Puntuacion = historial.Puntuacion
				params.Constestaciones = historiales_preguntas_de_este_examen

				historiales_preguntas_array = append(historiales_preguntas_array, params)
			}

			preguntas_w_respuestas, err := daoPreguntas.FindPreguntasByIds(ids_preguntas, bson.M{"pregunta": 0, "explicacion": 0})

			for _, pregunta_w_respuestas := range preguntas_w_respuestas {
				preguntas_devolver[pregunta_w_respuestas.Id.Hex()] = pregunta_w_respuestas
			}

			var areasWTemas []models.AreasWTemas
			areas, err := daoAreas.FindAllAreas()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
			for _, area := range areas {
				result, err := daoTemas.FindTemaByIdArea(area.Id.Hex())
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}

				var areaWTema models.AreasWTemas
				areaWTema.Id = area.Id
				areaWTema.Name = area.Name
				areaWTema.Temas = result

				areasWTemas = append(areasWTemas, areaWTema)
			}

			examenes_oficiales, err := daoExamenesOficiales.FindAllExamenesOficiales()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			casos_practicos, err := daoExamenesCasosPracticos.FindAllExamenesCasosPracticos()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			simulacros, err := daoSimulacros.FindAllSimulacros()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}
			basicosconfundidos, err := daoBasicosConfundidos.FindAllBasicosConfundidos()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			legislaciones, err := daoLegislaciones.FindAllLegislaciones()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}
			niveles, err := daoNiveles.FindAllNiveles()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			var json_preguntas_devolver = ""
			if len(preguntas_devolver) == 0 {
				json_preguntas_devolver = "[]"
			} else {
				prejson_preguntas_devolver, err := json.Marshal(preguntas_devolver)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_preguntas_devolver = string(prejson_preguntas_devolver)
			}

			var json_historiales_preguntas_array = ""
			if len(historiales_preguntas_array) == 0 {
				json_historiales_preguntas_array = "[]"
			} else {
				prejson_historiales_preguntas_array, err := json.Marshal(historiales_preguntas_array)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_historiales_preguntas_array = string(prejson_historiales_preguntas_array)
			}

			var json_areasWTemas = ""
			if len(areasWTemas) == 0 {
				json_areasWTemas = "[]"
			} else {
				prejson_areasWTemas, err := json.Marshal(areasWTemas)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_areasWTemas = string(prejson_areasWTemas)
			}

			var json_examenes_oficiales = ""
			if len(examenes_oficiales) == 0 {
				json_examenes_oficiales = "[]"
			} else {
				prejson_examenes_oficiales, err := json.Marshal(examenes_oficiales)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_examenes_oficiales = string(prejson_examenes_oficiales)
			}

			var json_simulacros = ""
			if len(simulacros) == 0 {
				json_simulacros = "[]"
			} else {
				prejson_simulacros, err := json.Marshal(simulacros)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_simulacros = string(prejson_simulacros)
			}

			var json_niveles = ""
			if len(niveles) == 0 {
				json_niveles = "[]"
			} else {
				prejson_niveles, err := json.Marshal(niveles)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_niveles = string(prejson_niveles)
			}

			var json_usuario = ""
			prejson_usuario, err := json.Marshal(usuario)
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
			json_usuario = string(prejson_usuario)

			var json_historiales_cp = ""
			if len(historiales_cp) == 0 {
				json_historiales_cp = "[]"
			} else {
				prejson_historiales_cp, err := json.Marshal(historiales_cp)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_historiales_cp = string(prejson_historiales_cp)
			}

			var json_casos_practicos = ""
			if len(casos_practicos) == 0 {
				json_casos_practicos = "[]"
			} else {
				prejson_casos_practicos, err := json.Marshal(casos_practicos)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_casos_practicos = string(prejson_casos_practicos)
			}
			var json_basicosconfundidos = ""
			if len(basicosconfundidos) == 0 {
				json_basicosconfundidos = "[]"
			} else {
				prejson_basicosconfundidos, err := json.Marshal(basicosconfundidos)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_basicosconfundidos = string(prejson_basicosconfundidos)
			}
			var json_legislaciones = ""
			if len(legislaciones) == 0 {
				json_legislaciones = "[]"
			} else {
				prejson_legislaciones, err := json.Marshal(legislaciones)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_legislaciones = string(prejson_legislaciones)
			}
			helper.ResponseWithJson(w, http.StatusCreated, map[string]string{
				"preguntas":          json_preguntas_devolver,
				"historiales":        json_historiales_preguntas_array,
				"areasWTemas":        json_areasWTemas,
				"examenes_oficiales": json_examenes_oficiales,
				"simulacros":         json_simulacros,
				"niveles":            json_niveles,
				"usuario":            json_usuario,
				"historiales_cp":     json_historiales_cp,
				"casos_practicos":    json_casos_practicos,
				"basicosconfundidos": json_basicosconfundidos,
				"legislaciones":      json_legislaciones,
			})
		}

	}
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	userCheck := helper.CheckUser(strings.TrimSpace(params["email"]))
	if !userCheck {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no autorizado"})
		return
	}
	usuarios, err := daoUsuarios.FindAllUsuarios()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al obtener los usuarios"})
		return
	} else {

		helper.ResponseWithJson(w, http.StatusCreated, map[string]interface{}{"result": "success", "data": usuarios})
	}
}

func SearchUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	userCheck := helper.CheckUser(strings.TrimSpace(params["email"]))
	if !userCheck {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no autorizado"})
		return
	}
	typeSearch := strings.TrimSpace(params["type"])
	dataSearch := strings.TrimSpace(params["data"])

	if typeSearch == "email" {
		result, err := daoUsuarios.FindUsuariosByEmail(dataSearch)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurrio un error al buscar usuarios por email"})
			return
		}
		helper.ResponseWithJson(w, http.StatusCreated, map[string]interface{}{"result": "success", "data": result})
		return
	} else {
		result, err := daoUsuarios.FindUsuarioByName(dataSearch)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurrio un error al buscar usuarios por Nombre"})
			return
		}
		helper.ResponseWithJson(w, http.StatusCreated, map[string]interface{}{"result": "success", "data": result})
		return
	}
}

func StatsUsuarioInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	userCheck := helper.CheckUser(strings.TrimSpace(params["email"]))
	if !userCheck {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no autorizado"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioById(strings.ToLower(strings.TrimSpace(params["userid"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {
		historiales, err := daoHistoriales.FindHistorialByIdUsuario(usuario.Id.Hex(), nil)
		historiales_cp, err := daoHistorialesCP.FindHistorialByIdUsuario(usuario.Id.Hex(), nil)

		var historiales_preguntas_array []models.HistorialesPipe
		var ids_preguntas []bson.ObjectId
		var ids_historiales []bson.ObjectId
		var preguntas_devolver = make(map[string]models.Preguntas)

		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
			return
		} else {
			for _, historial := range historiales {
				ids_historiales = append(ids_historiales, historial.Id)
			}

			historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdsHistoriales(ids_historiales)
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
				return
			}
			for _, historial_pregunta := range historiales_preguntas {
				ids_preguntas = append(ids_preguntas, historial_pregunta.IdPregunta)
			}

			for _, historial := range historiales {

				var historiales_preguntas_de_este_examen []models.HistorialesPreguntas

				for _, historial_pregunta := range historiales_preguntas {
					if historial.Id == historial_pregunta.IdHistorial {
						historiales_preguntas_de_este_examen = append(historiales_preguntas_de_este_examen, historial_pregunta)
					}
				}

				var params models.HistorialesPipe
				params.Id = historial.Id
				params.Tipo = historial.Tipo
				params.Fecha = historial.Fecha
				params.Temas = historial.Temas
				params.Legislaciones = historial.Legislaciones
				params.BasicosConfundidos = historial.BasicosConfundidos
				params.Oficiales = historial.Oficiales
				params.Simulacros = historial.Simulacros
				params.IdUsuario = historial.IdUsuario
				params.NumeroPreguntas = historial.NumeroPreguntas
				params.RespuestaAutomatica = historial.RespuestaAutomatica
				params.TiempoTranscurrido = historial.TiempoTranscurrido
				params.Terminado = historial.Terminado
				params.PreguntasTotales = historial.PreguntasTotales
				params.PreguntasAcertadas = historial.PreguntasAcertadas
				params.PreguntasFalladas = historial.PreguntasFalladas
				params.PreguntasBlancas = historial.PreguntasBlancas
				params.Puntuacion = historial.Puntuacion
				params.Constestaciones = historiales_preguntas_de_este_examen

				historiales_preguntas_array = append(historiales_preguntas_array, params)
			}

			preguntas_w_respuestas, err := daoPreguntas.FindPreguntasByIds(ids_preguntas, bson.M{"pregunta": 0, "explicacion": 0})

			for _, pregunta_w_respuestas := range preguntas_w_respuestas {
				preguntas_devolver[pregunta_w_respuestas.Id.Hex()] = pregunta_w_respuestas
			}

			var areasWTemas []models.AreasWTemas
			areas, err := daoAreas.FindAllAreas()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
			for _, area := range areas {
				result, err := daoTemas.FindTemaByIdArea(area.Id.Hex())
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}

				var areaWTema models.AreasWTemas
				areaWTema.Id = area.Id
				areaWTema.Name = area.Name
				areaWTema.Temas = result

				areasWTemas = append(areasWTemas, areaWTema)
			}

			examenes_oficiales, err := daoExamenesOficiales.FindAllExamenesOficiales()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			casos_practicos, err := daoExamenesCasosPracticos.FindAllExamenesCasosPracticos()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			simulacros, err := daoSimulacros.FindAllSimulacros()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}
			basicosconfundidos, err := daoBasicosConfundidos.FindAllBasicosConfundidos()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			legislaciones, err := daoLegislaciones.FindAllLegislaciones()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}
			niveles, err := daoNiveles.FindAllNiveles()
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}

			var json_preguntas_devolver = ""
			if len(preguntas_devolver) == 0 {
				json_preguntas_devolver = "[]"
			} else {
				prejson_preguntas_devolver, err := json.Marshal(preguntas_devolver)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_preguntas_devolver = string(prejson_preguntas_devolver)
			}

			var json_historiales_preguntas_array = ""
			if len(historiales_preguntas_array) == 0 {
				json_historiales_preguntas_array = "[]"
			} else {
				prejson_historiales_preguntas_array, err := json.Marshal(historiales_preguntas_array)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_historiales_preguntas_array = string(prejson_historiales_preguntas_array)
			}

			var json_areasWTemas = ""
			if len(areasWTemas) == 0 {
				json_areasWTemas = "[]"
			} else {
				prejson_areasWTemas, err := json.Marshal(areasWTemas)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_areasWTemas = string(prejson_areasWTemas)
			}

			var json_examenes_oficiales = ""
			if len(examenes_oficiales) == 0 {
				json_examenes_oficiales = "[]"
			} else {
				prejson_examenes_oficiales, err := json.Marshal(examenes_oficiales)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_examenes_oficiales = string(prejson_examenes_oficiales)
			}

			var json_simulacros = ""
			if len(simulacros) == 0 {
				json_simulacros = "[]"
			} else {
				prejson_simulacros, err := json.Marshal(simulacros)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_simulacros = string(prejson_simulacros)
			}

			var json_niveles = ""
			if len(niveles) == 0 {
				json_niveles = "[]"
			} else {
				prejson_niveles, err := json.Marshal(niveles)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_niveles = string(prejson_niveles)
			}

			var json_usuario = ""
			prejson_usuario, err := json.Marshal(usuario)
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
			json_usuario = string(prejson_usuario)

			var json_historiales_cp = ""
			if len(historiales_cp) == 0 {
				json_historiales_cp = "[]"
			} else {
				prejson_historiales_cp, err := json.Marshal(historiales_cp)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_historiales_cp = string(prejson_historiales_cp)
			}

			var json_casos_practicos = ""
			if len(casos_practicos) == 0 {
				json_casos_practicos = "[]"
			} else {
				prejson_casos_practicos, err := json.Marshal(casos_practicos)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_casos_practicos = string(prejson_casos_practicos)
			}
			var json_basicosconfundidos = ""
			if len(basicosconfundidos) == 0 {
				json_basicosconfundidos = "[]"
			} else {
				prejson_basicosconfundidos, err := json.Marshal(basicosconfundidos)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_basicosconfundidos = string(prejson_basicosconfundidos)
			}
			var json_legislaciones = ""
			if len(legislaciones) == 0 {
				json_legislaciones = "[]"
			} else {
				prejson_legislaciones, err := json.Marshal(legislaciones)
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
				json_legislaciones = string(prejson_legislaciones)
			}
			helper.ResponseWithJson(w, http.StatusCreated, map[string]string{
				"preguntas":          json_preguntas_devolver,
				"historiales":        json_historiales_preguntas_array,
				"areasWTemas":        json_areasWTemas,
				"examenes_oficiales": json_examenes_oficiales,
				"simulacros":         json_simulacros,
				"niveles":            json_niveles,
				"usuario":            json_usuario,
				"historiales_cp":     json_historiales_cp,
				"casos_practicos":    json_casos_practicos,
				"basicosconfundidos": json_basicosconfundidos,
				"legislaciones":      json_legislaciones,
			})
		}
	}
}

func ChangeDaySuscripcion(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.NewDayUserUpdate

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Invalid Payload"})
	}

	usuario, err := daoUsuarios.FindUsuarioById(data.Id)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problema al buscar el usuario"})
		return
	}

	usuario.MaximoDiaSuscripcion = data.Maximo_dia_suscripcion

	if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problema al actualizar el usuario"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteHistorial(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var historial models.DeleteHistorial

	if err := json.NewDecoder(r.Body).Decode(&historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Invalid Payload"})
		return
	}

	if err := daoHistoriales.RemoveHistorial(historial.Id); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problema al eliminar el historial"})
		return
	}

	if err := daoHistoriales.RemovePreguntaHistorial(historial.Id); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problema al eliminar las preguntas del historial"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "Historial eliminado con exito"})
}

func DeleteHistorialCP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var historial models.DeleteHistorial

	if err := json.NewDecoder(r.Body).Decode(&historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Invalid Payload"})
		return
	}

	if err := daoHistorialesCP.RemoveHistorial(historial.Id); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problema al eliminar el historial"})
		return
	}

	if err := daoHistorialesCP.RemovePreguntaHistorialCP(historial.Id); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Problema al eliminar las preguntas del historial"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "Historial eliminado con exito"})
}
