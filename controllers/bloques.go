package controllers

import (
	"encoding/json"
	"net/http"
	"test/helper"
	"test/models"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	daoBloques = models.Bloques{}
)

func AllBloques(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var bloques []models.Bloques
	bloques, err := daoBloques.FindAllBloques()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, bloques)

}

func FindBloque(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := daoBloques.FindBloqueById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, result)
}

func CreateBloque(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var bloque models.Bloques
	if err := json.NewDecoder(r.Body).Decode(&bloque); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	bloque.Id = bson.NewObjectId()
	if err := daoBloques.InsertBloque(bloque); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusCreated, bloque)
}

func UpdateBloque(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Bloques
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	if err := daoBloques.UpdateBloque(params); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteBloque(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := daoBloques.RemoveBloque(id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
