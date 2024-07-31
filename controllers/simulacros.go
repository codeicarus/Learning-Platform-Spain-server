package controllers

import (
	"encoding/json"
	"net/http"
	"test/helper"
	"test/models"
)

var (
	daoSimulacros          = models.Simulacros{}
	daoSimulacrosPreguntas = models.SimulacrosPreguntas{}
)

func DeleteSimulacro(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var params models.Simulacros
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})	
	}
	err := daoSimulacros.RemoveSimulacro(params.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurrio un error al eliminar el simulacro"})	
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "Simulacro eliminado con exito"})
}

func UpdateNameSimulacro(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
	}

	simulacro, err := daoSimulacros.FindSimulacroById(params["id"])

	if err !=nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "no se encontro el simulacro"})
	}

	simulacro.Name = params["newName"]

	if err := daoSimulacros.UpdateSimulacro(simulacro); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Error al actualizar el simulacro"})
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "Simulacro actualizado con exito"})


}