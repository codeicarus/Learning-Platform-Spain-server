package controllers

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"test/helper"
// 	"test/models"
// )

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"test/helper"
	"test/models"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	daoSuscripciones = models.Suscripcion{}
)

func EmailPayStarted(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/NewSuscripcion.html")
}

func EmailAccepPayStarted(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/AcceptSuscripcion.html")
}

func EmailRechazarPayStarted(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./email/RechazarSuscripcion.html")
}

func StartPayment(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close();

	var data models.Suscripcion

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	user,err := daoUsuarios.FindUsuarioByEmail(data.Email);

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al buscar el usuario"})
		return
	}

	sub, err := daoSuscripciones.FindSuscripcionByIdUser(user.Id.Hex())

	if err == nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Ya tiene una solicitud realizada de "+sub.Mount+"€"})
		return
	}

	data.Status = "PENDIENTE";
	data.UserID = user.Id
	data.Id = bson.NewObjectId()

	if err := daoSuscripciones.InsertSuscripcion(data); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al insertar la suscripcion"})
		return
	}



	config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/new_suscripcion")
		if err != nil {
			log.Println(err)
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "ocurrio un al obtener plantilla de email"})
			return
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "###CONCEPTO###", data.Concepto)
		texto_email = strings.ReplaceAll(texto_email, "###METHOD###", data.Type)
		texto_email = strings.ReplaceAll(texto_email, "###SUB###", data.Mount)

		helper.SendEmail("PENITENCIARIOS.COM :: Pago iniciado", texto_email, "penitenciarios@penitenciarios.com", "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusOK, map[string]string{"message": "Solicitud enviada"})

		return

}


func ViewAllPayments(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	sus, err := daoSuscripciones.FindAllSuscripcion()


	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al buscar las suscripciones"})
		return 
	}

	var result []models.SuscripcionList

	for _, sub := range sus {
		user, err := daoUsuarios.FindUsuarioById(sub.UserID.Hex())

		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al buscar el usuario"})
			return
		}
		var usuario models.User
		usuario.Email = user.Email
		usuario.Name = user.Name
		usuario.Id = user.Id
		usuario.MaximoDiaSuscripcion = user.MaximoDiaSuscripcion


		var data models.SuscripcionList
		data.Id = bson.ObjectId(sub.Id)
		data.Concepto = sub.Concepto
		data.Suscription = sub.Suscription
  	data.Email = sub.Email
  	data.UserID = sub.UserID 
  	data.Type = sub.Type
  	data.Concepto = sub.Concepto   
  	data.Status = sub.Status
  	data.Mount = sub.Mount
  	data.User =  	usuario




		result = append(result, data)
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"message":"success","data": result})
}


func AcceptPayment(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	var data models.Suscripcion

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	var suscriptionTime string

	if data.Suscription == "0990" {
		suscriptionTime = "por 1 Mes."
	}else if data.Suscription == "2790" {
		suscriptionTime = "por 3 Mes."
	}else if data.Suscription == "salida" {
		suscriptionTime = "hasta el 15 de Septiembre."
	}

	user, err := daoUsuarios.FindUsuarioById(data.UserID.Hex())

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al buscar el usuario"})
		return
	}



	var dia_maximo int64 = 0
		dia_maximo_actual := time.Now().Unix()

		if dia_maximo_actual < user.MaximoDiaSuscripcion {
			dia_maximo_actual = user.MaximoDiaSuscripcion
		}

		dia_maximo_actual_date := time.Unix(dia_maximo_actual, 0)

		if data.Suscription == "0990" {
			after := dia_maximo_actual_date.AddDate(0, 1, 0)
			dia_maximo = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()
		} else if data.Suscription == "2790" {
			after := dia_maximo_actual_date.AddDate(0, 3, 0)
			dia_maximo = time.Date(after.Year(), after.Month(), after.Day(), 23, 59, 59, 0, time.UTC).Unix()
		} else if data.Suscription == "salida" {
			dia_maximo = 1694833199
		}

		user.MaximoDiaSuscripcion = dia_maximo

		if err := daoUsuarios.UpdateUsuario(user); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		data.Status = "APROBADO"

		if err := daoSuscripciones.UpdateSuscripcion(data); err != nil {
			log.Println(err.Error())
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al actualizar la sub"})
			return
		}

	config_vars := helper.GetConfigVars()

	res, err := http.Get(config_vars["URL"] + "email/accept_suscripcion")
	if err != nil {
		log.Println(err)
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "ocurrio un al obtener plantilla de email"})
		return
	}
	texto_email_bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	texto_email := string(texto_email_bytes)
	texto_email = strings.ReplaceAll(texto_email, "###NAME###", user.Name)
	texto_email = strings.ReplaceAll(texto_email, "###TIME###", suscriptionTime)

	helper.SendEmail("PENITENCIARIOS.COM :: PAGO ACEPTADO", texto_email, user.Email, "penitenciarios@penitenciarios.com")

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"message": "Pago aceptado con exito"})
	return
}

func RechazarPago(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	var data models.SuscripcionRechazar

	if err:= json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}


	user, err := daoUsuarios.FindUsuarioById(data.UserID.Hex())

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al buscar el usuario"})
		return
	}

	sus, err := daoSuscripciones.FindSuscripcionByID(data.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	sus.Status = "RECHAZADO"

	if err := daoSuscripciones.UpdateSuscripcion(sus); err != nil {
		log.Println(err.Error())
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al actualizar la sub"})
		return
	}

	config_vars := helper.GetConfigVars()

	res, err := http.Get(config_vars["URL"] + "email/rechazar_suscripcion")
	if err != nil {
		log.Println(err)
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "ocurrio un al obtener plantilla de email"})
		return
	}
	texto_email_bytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	texto_email := string(texto_email_bytes)
	texto_email = strings.ReplaceAll(texto_email, "###NAME###", user.Name)
	texto_email = strings.ReplaceAll(texto_email, "###MESSAGE###", data.Message)

	helper.SendEmail("PENITENCIARIOS.COM :: PAGO ACEPTADO", texto_email, user.Email, "penitenciarios@penitenciarios.com")

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"message": "Pago Rechazado"})
	return
}


func DeleteSuscripcion(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	subID := mux.Vars(r)["id"]

	
	if err := daoSuscripciones.DeleteSuscripcion(subID); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"message": "Pago Eliminado"})
}
