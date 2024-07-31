package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"test/helper"
	"test/models"

	"github.com/globalsign/mgo/bson"
)

var (
	daoPreguntasCPFavoritas = models.PreguntasCPFavoritas{}
)

func GetPreguntasCPFavoritas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var preguntas_favoritas_a_devolver = map[string]string{}
	var params = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha porducido un error"})
		return
	} else {
		preguntas_favoritas, _ := daoPreguntasCPFavoritas.FindAllPreguntasCPFavoritasByIdUsuario(usuario.Id.Hex())

		if len(preguntas_favoritas) > 0 {
			for _, pregunta_favorita := range preguntas_favoritas {
				preguntas_favoritas_a_devolver[pregunta_favorita.IdPregunta.Hex()] = pregunta_favorita.IdPregunta.Hex()
			}
		}

		helper.ResponseWithJson(w, http.StatusOK, preguntas_favoritas_a_devolver)
	}
}
func SetPreguntasCPFavoritas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var preguntas_favoritas_a_devolver = map[string]string{}
	var params = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha porducido un error"})
		return
	} else {

		var pregunta_favorita models.PreguntasCPFavoritas
		pregunta_favorita.Id = bson.NewObjectId()
		pregunta_favorita.IdUsuario = usuario.Id
		pregunta_favorita.IdPregunta = bson.ObjectIdHex(params["id_pregunta"])
		if err := daoPreguntasCPFavoritas.InsertPreguntaCPFavorita(pregunta_favorita); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha porducido un error"})
			return
		}

		preguntas_favoritas, _ := daoPreguntasCPFavoritas.FindAllPreguntasCPFavoritasByIdUsuario(usuario.Id.Hex())

		if len(preguntas_favoritas) > 0 {
			for _, pregunta_favorita := range preguntas_favoritas {
				preguntas_favoritas_a_devolver[pregunta_favorita.IdPregunta.Hex()] = pregunta_favorita.IdPregunta.Hex()
			}
		}

		helper.ResponseWithJson(w, http.StatusOK, preguntas_favoritas_a_devolver)
	}
}
func DeletePreguntasCPFavoritas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var preguntas_favoritas_a_devolver = map[string]string{}
	var params = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha porducido un error"})
		return
	} else {
		if err := daoPreguntasCPFavoritas.RemovePreguntaCPFavoritaByIdUsuarioIdPregunta(usuario.Id.Hex(), params["id_pregunta"]); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha porducido un error"})
			return
		}

		preguntas_favoritas, _ := daoPreguntasCPFavoritas.FindAllPreguntasCPFavoritasByIdUsuario(usuario.Id.Hex())

		if len(preguntas_favoritas) > 0 {
			for _, pregunta_favorita := range preguntas_favoritas {
				preguntas_favoritas_a_devolver[pregunta_favorita.IdPregunta.Hex()] = pregunta_favorita.IdPregunta.Hex()
			}
		}

		helper.ResponseWithJson(w, http.StatusOK, preguntas_favoritas_a_devolver)
	}
}
