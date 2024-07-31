package controllers

import (
	"net/http"
	"test/helper"
	"test/models"
)

var (
	daoExamenesOficiales          = models.ExamenesOficiales{}
	daoExamenesOficialesPreguntas = models.ExamenesOficialesPreguntas{}
)

func GetOficiales(w http.ResponseWriter, r *http.Request) {

	var oficiales_simulacros models.OficialesYSimulacros

	examenes_oficiales, err := daoExamenesOficiales.FindAllExamenesOficiales()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}
	simulacros, err := daoSimulacros.FindAllSimulacros()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}
	niveles, err := daoNiveles.FindAllNiveles()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}

	oficiales_simulacros.Oficiales = examenes_oficiales
	oficiales_simulacros.Simulacros = simulacros
	oficiales_simulacros.Niveles = niveles

	helper.ResponseWithJson(w, http.StatusOK, oficiales_simulacros)
}
