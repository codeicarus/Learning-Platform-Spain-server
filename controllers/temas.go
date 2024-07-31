package controllers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"test/helper"
	"test/models"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/xuri/excelize/v2"
)

var (
	daoTemas = models.Temas{}
)

func AllTemas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var temas []models.Temas
	temas, err := daoTemas.FindAllTemas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, temas)

}

func FindTema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := daoTemas.FindTemaById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, result)
}

func CreateTema(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var tema models.Temas
	if err := json.NewDecoder(r.Body).Decode(&tema); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	tema.Id = bson.NewObjectId()
	if err := daoTemas.InsertTema(tema); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusCreated, tema)
}

func UpdateTema(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Temas
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	if err := daoTemas.UpdateTema(params); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteTema(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := daoTemas.RemoveTema(id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func FindDataTema(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	tema, err := daoTemas.FindTemaById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	var preguntasID []bson.ObjectId
	preguntasID = append(preguntasID, tema.Id)
	preguntas, err := daoPreguntas.FindPreguntaByIdsTemas(preguntasID)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"result": "success", "tema": tema, "preguntas": preguntas})
}

func UpdateNameTema (w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data models.UpdateName

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	tema, err := daoTemas.FindTemaById(data.Id.Hex())

	if  err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return 
	}

	tema.Name = data.NewName

	if err := daoTemas.UpdateTema(tema); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}


	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteAllPreguntasTema(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)

	id := vars["id"]

	preguntas, err := daoPreguntas.FindAllPreguntasByTema(id)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	for _, pregunta := range preguntas {
		if err := daoRespuestas.RemoveRespuestasByIdPregunta(pregunta.Id.Hex()); err != nil{
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		if err := daoPreguntas.RemovePregunta(pregunta.Id.Hex()); err != nil{
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}
	}


	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func UploadQuestionTema(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)

	id := vars["id"]

	thisTema, err := daoTemas.FindTemaById(id)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}


	var importadores models.Importadores

	if err := json.NewDecoder(r.Body).Decode(&importadores); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	dec, err := base64.StdEncoding.DecodeString(importadores.File)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	f, err := os.Create("excel.xlsx")
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	if err := f.Sync(); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	fileRead, err := excelize.OpenFile("excel.xlsx")
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	// Get value from cell by given worksheet name and axis.
	sheets := fileRead.GetSheetList()

	var cont = 0

	var errores_importacion = ""
	for _, sheet := range sheets {
		cont = 0
		rows, err := fileRead.GetRows(sheet)

		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		for _, row := range rows {

			if cont > 0 && len(row) >= 9 {
				if row[1] != "" && row[5] != "" {
					if row[0] != "" {
						_, err := daoNiveles.FindNivelByAbreviacion(row[0])
						if err != nil {
							errores_importacion = errores_importacion + "Nivel '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					// if row[2] != "" {
					// 	_, err := daoNiveles.FindNivelByAbreviacion(row[2])
					// 	if err != nil {
					// 		errores_importacion = errores_importacion + "Nivel '" + row[1] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					// 	}
					// }

					// if row[2] != "" && row[2] != "NO" {
					// 	_, err := daoTemas.FindBasicoConfundidoByAbreviacion(row[2])
					// 	if err != nil {
					// 		errores_importacion = errores_importacion + "Tema  '" + row[2] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					// 	}

					// }


					respuestas := strings.Split(row[3], "\n")

					// a := []rune(row[7])
					// letter := string(a[0:1])

					a := []rune(strings.ToUpper(row[4]))

					letter := string(a[0:1])

					var alguna_true = false

					if len(respuestas) == 4 {
						var contR = 0
						for _ = range respuestas {

							if contR == 0 && letter == "A" {
								alguna_true = true
							} else if contR == 1 && letter == "B" {
								alguna_true = true
							} else if contR == 2 && letter == "C" {
								alguna_true = true
							} else if contR == 3 && letter == "D" {
								alguna_true = true
							}
							contR = contR + 1
						}

					} else {
						errores_importacion = errores_importacion + "No hay 4 respuestas en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					}
					if !alguna_true {
						errores_importacion = errores_importacion + "No hay ninguna respuesta correcta en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					}
				}

			}
			cont = cont + 1
		}
	}

	if errores_importacion != "" {

		//#########################################

		// config_vars := helper.GetConfigVars()

		// res, err := http.Get(config_vars["URL"] + "email/testError")
		// if err != nil {
		// 	log.Println(err)
		// }
		// texto_email_bytes, err := ioutil.ReadAll(res.Body)
		// res.Body.Close()
		// texto_email := string(texto_email_bytes)
		// texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		// texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		// helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")
		log.Println("ERROR")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		// 	//#########################################

		return
	}

	cont = 0
	for _, sheet := range sheets {
		cont = 0
		rows, err := fileRead.GetRows(sheet)
		if err != nil {
			log.Println("ERROR AQUI EN LA LINEA 326 de Temas")
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		for _, row := range rows {
			if cont > 0 && len(row) >= 9 {
				if row[1] != "" && row[5] != "" {
					var preguntaAGuardar models.Preguntas
					preguntaAGuardar.Id = bson.NewObjectId()
					preguntaAGuardar.IdTema = thisTema.Id

					if row[0] != "" {
						result, err := daoNiveles.FindNivelByAbreviacion(row[0])
						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdNivel = result.Id
					}

					// if row[1] != "" {

					// 	result, err := daoNiveles.FindNivelByAbreviacion(row[1])
					// 	if err != nil {
					// 		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel '" + row[1] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
					// 		return
					// 	}
					// 	preguntaAGuardar.IdNivel = result.Id
					// }

					preguntaAGuardar.Pregunta = row[2]
					preguntaAGuardar.Explicacion = row[5]
					preguntaAGuardar.Oficial = false
					preguntaAGuardar.AnioOficial = ""

					respuestas := strings.Split(row[3], "\n")

					a := []rune(strings.ToUpper(row[4]))

					letter := string(a[0:1])

					if len(respuestas) == 4 {
						var contR = 0
						for _, respuesta := range respuestas {
							var respuestaAGuardar models.Respuestas
							respuestaAGuardar.Id = bson.NewObjectId()
							respuestaAGuardar.Respuesta = respuesta[3:]
							if contR == 0 && letter == "A" {
								respuestaAGuardar.Correcta = true
							} else if contR == 1 && letter == "B" {
								respuestaAGuardar.Correcta = true
							} else if contR == 2 && letter == "C" {
								respuestaAGuardar.Correcta = true
							} else if contR == 3 && letter == "D" {
								respuestaAGuardar.Correcta = true
							} else {
								respuestaAGuardar.Correcta = false
							}
							respuestaAGuardar.IdPregunta = preguntaAGuardar.Id
							if err := daoRespuestas.InsertRespuesta(respuestaAGuardar); err != nil {
								helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
								return
							}
							contR = contR + 1
						}

						if err := daoPreguntas.InsertPregunta(preguntaAGuardar); err != nil {
							log.Println("pregunta noguardada")
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
							return
						}
					} else {

						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "No hay 4 respuestas en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
						return
					}
				}
			}
			cont = cont + 1
		}
	}
	//#########################################

	// config_vars := helper.GetConfigVars()

	// result, err := http.Get(config_vars["URL"] + "email/testSuccess")
	// if err != nil {
	// 	log.Println(err)
	// }
	// texto_email_bytes, err := ioutil.ReadAll(result.Body)
	// result.Body.Close()
	// texto_email := string(texto_email_bytes)
	// texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)

	// helper.SendEmail("PENITENCIARIOS.COM :: Test Subido con exito", texto_email, importadores.Email, "inpenitenciariosfo@penitenciarios.com")

	//#########################################
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
