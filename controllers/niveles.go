package controllers

import (
	"encoding/json"
	"net/http"
	"test/helper"
	"test/models"

	// "github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	daoNiveles = models.Niveles{}
)

func AllNiveles(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var niveles []models.Niveles
	niveles, err := daoNiveles.FindAllNiveles()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, niveles)

}

func FindNivel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := daoNiveles.FindNivelById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, result)
}

func CreateNivel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var nivel models.Niveles
	if err := json.NewDecoder(r.Body).Decode(&nivel); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	// nivel.Id = bson.NewObjectId()
	if err := daoNiveles.InsertNivel(nivel); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusCreated, nivel)
}

func UpdateNivel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Niveles
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	if err := daoNiveles.UpdateNivel(params); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteNivel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := daoNiveles.RemoveNivel(id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
