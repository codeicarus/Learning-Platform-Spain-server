package controllers

import (
	"net/http"
	"test/helper"
	"test/models"
)

var (
	daoExamenesCasosPracticos  = models.ExamenesCasosPracticos{}
	daoCasosPracticos          = models.CasosPracticos{}
	daoCasosPracticosPreguntas = models.CasosPracticosPreguntas{}
)

func GetCasosPracticos(w http.ResponseWriter, r *http.Request) {

	casos_practicos, err := daoExamenesCasosPracticos.FindAllExamenesCasosPracticos()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, casos_practicos)
}
