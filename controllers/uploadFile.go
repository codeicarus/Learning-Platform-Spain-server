package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"test/helper"
	"test/models"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	daoFiles = models.Files{}
)

func UploadFile2(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	
	file, header, err := r.FormFile("file")

	vars := mux.Vars(r)
	
	id := vars["id"]
	
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener el archivo"})
		return
	}

	fileExiste, err := daoFiles.FindFileByName(header.Filename)

	if err != nil && err.Error() != "not found" {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "El archivo " + fileExiste.FileName + " ya existe"})
			return
	} else if fileExiste != (models.Files{}) {
	  	helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "El archivo " + fileExiste.FileName + " ya existe"})
			return
	}

	fileData, err := ioutil.ReadAll(file)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al leer el archivo"})
		return
	}

	dir, err := os.Getwd()
	config_vars := helper.GetConfigVars()

	urlBack := config_vars["URL"]

	filePath := filepath.Join(dir, "files", header.Filename)

	filePathURL := "files/" + url.PathEscape(header.Filename)
	fileURL := urlBack + filePathURL

	err = ioutil.WriteFile(filePath, fileData, 0644)

	if err != nil {
		log.Println(err.Error())
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al guardar el archivo"})
		return
	}

	var fileDB models.Files
	pageLink := r.FormValue("PageLink")

	fileDB.Id = bson.NewObjectId()
	fileDB.FileName = header.Filename
	fileDB.FileSize = header.Size
	fileDB.Url = fileURL
	fileDB.IdPregunta = bson.ObjectIdHex(id)
	fileDB.PageLink = pageLink
	
	if err := daoFiles.InsertFile(fileDB); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al guardar el archivo"})
		return
	}

	pregunta, err := daoPreguntas.FindPreguntaById(id)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al buscar la pregunta"})
		return
	}

	pregunta.SchemaId = append(pregunta.SchemaId, fileDB.Id)

	if err := daoPreguntas.UpdatePregunta(pregunta); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al actualizar la pregunta"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"succes": fileDB})
}


func ViewFile(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	clientReferer := []string{
		"https://penitenciarios.com/", 
		"https://www.penitenciarios.com/", 
		"https://dev.penitenciarios.com/",
		"https://www.dev.penitenciarios.com/",
		"http://localhost:3000/", 
		"http://127.0.0.1:3000/"}


		referer := r.Referer()

	validReferer := false
	for _, allowedReferer := range clientReferer {
		if referer == allowedReferer {
			validReferer = true
			break
		}
	}

	if !validReferer {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "No autorizado"})
		return
	}

	archivoID := mux.Vars(r)["id"]

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filePath := filepath.Join(dir, "files", archivoID)

	// Verificar si el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Archivo no encontrado", http.StatusNotFound)
		return
	}

	// Servir el archivo
	http.ServeFile(w, r, filePath)
}

func GetViewAllFiles(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	files, err := daoFiles.FindAllFiles()

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener los archivos"})
		return
	}

	var filesView []models.FileData

	for _, file := range files {
		pregunta, err := daoPreguntas.FindPreguntaById(file.IdPregunta.Hex())
		if err != nil {
			log.Println(err.Error())
		}
		var fileData models.FileData
		fileData.FileData = file
		fileData.Pregunta = pregunta

		filesView = append(filesView, fileData)

	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"files": filesView})
}

func DeleteFile(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]

	file, err := daoFiles.FindFileByID(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener el archivo de la bbdd"})
		return
	}

	if err := daoFiles.DeleteFileByID(id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al eliminar el archivo"})
		return
	}


	filePath := filepath.Join("files", file.FileName)

	if err := os.Remove(filePath); err != nil {
		log.Println(err.Error())
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"error": "Error al eliminar el archivo"})
		return
}

	files, err := daoFiles.FindAllFiles()

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener los archivos"})
		return
	}

	var filesView []models.FileData

	for _, file := range files {
		pregunta, err := daoPreguntas.FindPreguntaById(file.IdPregunta.Hex())
		if err != nil {
			log.Println(err.Error())
		}
		var fileData models.FileData
		fileData.FileData = file
		fileData.Pregunta = pregunta

		filesView = append(filesView, fileData)

	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"files": filesView})
}

func DeleteQuestionAndSchema (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	vars := mux.Vars(r)

	// idSchema := vars["idSchema"]
	idQuestion := vars["idQuestion"]

	// file, err := daoFiles.FindFileByID(idSchema)
	// if err != nil {
	// 	helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener el archivo de la bbdd"})
	// 	return
	// }

	filesByPreguntaId, err := daoFiles.FindFilesByIdPregunta(idQuestion)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener los archivos"})
		return
	}

	for _, fileInfo := range filesByPreguntaId {
		dir, err := os.Getwd()		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		filePath := filepath.Join(dir,"files", fileInfo.FileName)
		
		if err := os.Remove(filePath); err != nil {
			log.Println(err.Error())
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"error": "Error al eliminar el archivo"})
			return
		}
	}

if err := daoFiles.DeleteFilesByIdPregunta(idQuestion); err != nil {
	helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al eliminar los archivos"})
	return
}

	if err := daoPreguntas.RemovePregunta(idQuestion); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al eliminar la pregunta"})
		return
	}

	if err := daoRespuestas.RemoveRespuestasByIdPregunta(idQuestion); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al eliminar las respuestas"})
		return
	}


	files, err := daoFiles.FindAllFiles()

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"error": "Error al obtener los archivos"})
		return
	}

	var filesView []models.FileData

	for _, file := range files {
		pregunta, err := daoPreguntas.FindPreguntaById(file.IdPregunta.Hex())
		if err != nil {
			log.Println(err.Error())
		}
		var fileData models.FileData
		fileData.FileData = file
		fileData.Pregunta = pregunta

		filesView = append(filesView, fileData)

	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"files": filesView})
}