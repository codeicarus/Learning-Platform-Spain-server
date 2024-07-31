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
	daoAreas = models.Areas{}
)

func AllAreas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var areas []models.Areas
	areas, err := daoAreas.FindAllAreas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, areas)

}

func AllAreasWLegislaciones(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var areasWLegislaciones []models.AreasWLegislaciones
	var areas []models.Areas
	areas, err := daoAreas.FindAllAreas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	for _, area := range areas {
		result, err := daoLegislaciones.FindLegislacionByIdArea(area.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var areaWLegislacion models.AreasWLegislaciones
		areaWLegislacion.Id = area.Id
		areaWLegislacion.Name = area.Name
		areaWLegislacion.Legislaciones = result

		areasWLegislaciones = append(areasWLegislaciones, areaWLegislacion)
	}
	helper.ResponseWithJson(w, http.StatusOK, areasWLegislaciones)

}
func AllAreasWBasicosConfundidos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var areasWBasicosConfundidos []models.AreasWBasicosConfundidos
	var areas []models.Areas
	areas, err := daoAreas.FindAllAreas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	for _, area := range areas {
		result, err := daoBasicosConfundidos.FindBasicoConfundidoByIdArea(area.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var areaWLegislacion models.AreasWBasicosConfundidos
		areaWLegislacion.Id = area.Id
		areaWLegislacion.Name = area.Name
		areaWLegislacion.BasicosConfundidos = result

		areasWBasicosConfundidos = append(areasWBasicosConfundidos, areaWLegislacion)
	}

	helper.ResponseWithJson(w, http.StatusOK, areasWBasicosConfundidos)

}
func AllAreasWTemas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var areasWTemas []models.AreasWTemas
	var areas []models.Areas
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
	helper.ResponseWithJson(w, http.StatusOK, areasWTemas)

}
func AllAreasWBloques(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var areasWBloques []models.AreasWBloques
	var areas []models.Areas
	areas, err := daoAreas.FindAllAreas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	for _, area := range areas {
		result, err := daoBloques.FindBloqueByIdArea(area.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var areaWBloque models.AreasWBloques
		areaWBloque.Id = area.Id
		areaWBloque.Name = area.Name
		areaWBloque.Bloques = result

		areasWBloques = append(areasWBloques, areaWBloque)
	}
	helper.ResponseWithJson(w, http.StatusOK, areasWBloques)

}

func FindArea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := daoAreas.FindAreaById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, result)
}

// func FindAreaName(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	name := vars["name"]
// 	result, err := daoAreas.FindAreaByName(name)
// 	if err != nil {
// 		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
// 		return
// 	}
// 	helper.ResponseWithJson(w, http.StatusOK, result)
// }

func CreateArea(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var area models.Areas
	if err := json.NewDecoder(r.Body).Decode(&area); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	area.Id = bson.NewObjectId()
	if err := daoAreas.InsertArea(area); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusCreated, area)
}

func UpdateArea(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Areas
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	if err := daoAreas.UpdateArea(params); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteArea(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := daoAreas.RemoveArea(id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
