package controllers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"

	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"test/helper"
	"test/models"

	// "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/globalsign/mgo/bson"
	"github.com/xuri/excelize/v2"
)

var (
	daoImportadores = models.Importadores{}

	daoBasicosConfundidos = models.BasicosConfundidos{}
)

func extractPattern(str string) string {
	re := regexp.MustCompile(`^([A-Za-z])(\d+)`)
	match := re.FindStringSubmatch(str)
	if len(match) == 3 {
		return match[1] + match[2]
	}
	return ""
}

func DoImport(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
				if row[5] != "" && row[6] != "" && row[8] != "" {
					if row[0] != "" {
						_, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(row[0])
						if err != nil {
							errores_importacion = errores_importacion + "Basico confundido '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					if row[1] != "" {
						_, err := daoLegislaciones.FindLegislacionByAbreviacion(row[1])
						if err != nil {
							errores_importacion = errores_importacion + "Legislacion '" + row[1] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					if row[2] != "" && row[2] != "NO" {
						_, err := daoTemas.FindBasicoConfundidoByAbreviacion(row[2])
						if err != nil {
							errores_importacion = errores_importacion + "Tema  '" + row[2] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}

					}

					if row[3] != "" {
						_, err := daoNiveles.FindNivelByAbreviacion(row[3])
						if err != nil {
							errores_importacion = errores_importacion + "Nivel '" + row[3] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					} else {
						errores_importacion = errores_importacion + "Nivel no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					}

					respuestas := strings.Split(row[6], "\n")

					// a := []rune(row[7])
					// letter := string(a[0:1])

					a := []rune(strings.ToUpper(row[7]))

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

		config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/testError")
		if err != nil {
			log.Println(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		// 	//#########################################

		return
	}

	cont = 0
	for _, sheet := range sheets {
		cont = 0
		rows, err := fileRead.GetRows(sheet)
		if err != nil {
			log.Println("ERROR AQUI EN LA LINEA 196 de importadores")
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		for _, row := range rows {
			if cont > 0 && len(row) >= 9 {
				if row[5] != "" && row[6] != "" && row[8] != "" {
					var preguntaAGuardar models.Preguntas
					preguntaAGuardar.Id = bson.NewObjectId()

					if row[0] != "" {
						result, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(row[0])
						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Basico confundido '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdBasicosConfundidos = result.Id
					}

					if row[1] != "" {

						result, err := daoLegislaciones.FindLegislacionByAbreviacion(row[1])
						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Legislacion '" + row[1] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdLegislacion = result.Id
					}

					if row[2] != "" && row[2] != "NO" {
						result, err := daoTemas.FindBasicoConfundidoByAbreviacion(row[2])

						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Tema  2131232'" + row[2] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdTema = result.Id
					}

					if row[3] != "" {
						result, err := daoNiveles.FindNivelByAbreviacion(row[3])
						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel '" + row[3] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdNivel = result.Id
					} else {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
						return
					}

					preguntaAGuardar.Pregunta = row[5]
					preguntaAGuardar.Explicacion = row[8]
					preguntaAGuardar.Oficial = false
					preguntaAGuardar.AnioOficial = ""

					respuestas := strings.Split(row[6], "\n")

					// a := []rune(row[7])
					// letter := string(a[0:1])
					a := []rune(strings.ToUpper(row[7]))

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

	config_vars := helper.GetConfigVars()

	result, err := http.Get(config_vars["URL"] + "email/testSuccess")
	if err != nil {
		log.Println(err)
	}
	texto_email_bytes, err := ioutil.ReadAll(result.Body)
	result.Body.Close()
	texto_email := string(texto_email_bytes)
	texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)

	helper.SendEmail("PENITENCIARIOS.COM :: Test Subido con exito", texto_email, importadores.Email, "inpenitenciariosfo@penitenciarios.com")

	//#########################################
}

func DoImportOficiales(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
				if row[5] != "" && row[6] != "" && row[8] != "" {
					if row[0] != "" {
						_, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(row[0])
						if err != nil {
							errores_importacion = errores_importacion + "Basico confundido '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					if row[1] != "" {
						_, err := daoLegislaciones.FindLegislacionByAbreviacion(row[1])
						if err != nil {
							errores_importacion = errores_importacion + "Legislacion '" + row[1] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					if row[2] != "" && row[2] != "NO" {
						_, err := daoTemas.FindBasicoConfundidoByAbreviacion(row[2])
						if err != nil {
							errores_importacion = errores_importacion + "Tema '" + row[2] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					if row[3] != "" {
						_, err := daoNiveles.FindNivelByAbreviacion(row[3])
						if err != nil {
							errores_importacion = errores_importacion + "Nivel '" + row[3] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					} else {
						errores_importacion = errores_importacion + "Nivel no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					}

					respuestas := strings.Split(row[6], "\n")

					var alguna_true = false

					// a := []rune(row[7])
					// letter := string(a[0:1])
					a := []rune(strings.ToUpper(row[7]))

					letter := string(a[0:1])

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

		config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/testError")
		if err != nil {
			log.Println(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		//#########################################
		return
	}

	cont = 0
	for _, sheet := range sheets {
		var examenOficialAGuardar models.ExamenesOficiales
		examenOficialAGuardar.Id = bson.NewObjectId()
		examenOficialAGuardar.Name = sheet
		if err := daoExamenesOficiales.InsertExamenesOficial(examenOficialAGuardar); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'ExamenesOficiales'"})
			return
		}

		cont = 0
		rows, err := fileRead.GetRows(sheet)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		for _, row := range rows {
			if cont > 0 && len(row) >= 9 {
				if row[5] != "" && row[6] != "" && row[8] != "" {
					var preguntaAGuardar models.Preguntas
					preguntaAGuardar.Id = bson.NewObjectId()

					if row[0] != "" {
						result, _ := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(row[0])
						preguntaAGuardar.IdBasicosConfundidos = result.Id
					}

					if row[1] != "" {
						result, _ := daoLegislaciones.FindLegislacionByAbreviacion(row[1])
						preguntaAGuardar.IdLegislacion = result.Id
					}

					if row[2] != "" && row[2] != "NO" {
						result, _ := daoTemas.FindBasicoConfundidoByAbreviacion(row[2])
						preguntaAGuardar.IdTema = result.Id
					}

					if row[3] != "" {
						result, _ := daoNiveles.FindNivelByAbreviacion(row[3])
						preguntaAGuardar.IdNivel = result.Id
					}

					preguntaAGuardar.Pregunta = row[5]
					preguntaAGuardar.Explicacion = row[8]
					preguntaAGuardar.Oficial = true
					preguntaAGuardar.AnioOficial = sheet

					respuestas := strings.Split(row[6], "\n")

					// a := []rune(row[7])
					// letter := string(a[0:1])
					a := []rune(strings.ToUpper(row[7]))

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
								helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'Respuestas'"})
								return
							}
							contR = contR + 1
						}

						if err := daoPreguntas.InsertPregunta(preguntaAGuardar); err != nil {
							log.Println(err)
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'Preguntas'"})
							return
						}

						var examenOficialPreguntaAGuardar models.ExamenesOficialesPreguntas
						examenOficialPreguntaAGuardar.Id = bson.NewObjectId()
						examenOficialPreguntaAGuardar.IdExamenOficial = examenOficialAGuardar.Id
						examenOficialPreguntaAGuardar.IdPregunta = preguntaAGuardar.Id
						examenOficialPreguntaAGuardar.Orden = cont
						if err := daoExamenesOficialesPreguntas.InsertExamenOficialPregunta(examenOficialPreguntaAGuardar); err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'ExamenesOficialesPreguntas'"})
							return
						}

					}
				}
			}
			cont = cont + 1
		}
	}
}

func DoImportCP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
			if cont > 1 && len(row) >= 7 {
				if row[3] != "" && row[4] != "" && row[5] != "" && row[6] != "" {

					respuestas := strings.Split(row[4], "\n")

					// a := []rune(row[5])
					// letter := string(a[0:1])
					a := []rune(strings.ToUpper(row[5]))

					letter := string(a[0:1])

					var alguna_true = false

					if len(respuestas) >= 4 {
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

		config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/testError")
		if err != nil {
			log.Println(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		//#########################################
		return
	}

	cont = 0
	for _, sheet := range sheets {


		var examenCasoPracticoAGuardar models.ExamenesCasosPracticos
		var casos_practicos_del_examen = []bson.ObjectId{}
		
		ex, err := daoExamenesCasosPracticos.FindExamenCasoPracticoByName(sheet)

		if err != nil {
			if err.Error() == "not found" {
				log.Println("Creando nuevo examen")
				examenCasoPracticoAGuardar.Id = bson.NewObjectId()
				examenCasoPracticoAGuardar.Name = sheet
				examenCasoPracticoAGuardar.Oficial = false
				examenCasoPracticoAGuardar.CasosPracticos = casos_practicos_del_examen

				if err := daoExamenesCasosPracticos.InsertExamenCasoPractico(examenCasoPracticoAGuardar); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'ExamenesCasosPracticos'"})
					return
				}
			}
		}

		

		cont = 0
		rows, err := fileRead.GetRows(sheet)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var lastId = bson.NewObjectId()

		for _, row := range rows {

			if cont >= 1 && len(row) >= 7 {
				
				if row[3] != "" && row[4] != "" && row[5] != "" && row[6] != "" {
					

					if row[0] != "" {
						_, err := daoCasosPracticos.FindCasoPracticoByName(row[0])
						
						if err != nil {
							if err.Error() == "not found" {
								var casoPracticoAGuardar models.CasosPracticos
								casoPracticoAGuardar.Id = bson.NewObjectId()
								casoPracticoAGuardar.Name = row[0]
								casoPracticoAGuardar.Oficial = false
								casoPracticoAGuardar.AnioOficial = ""
								casoPracticoAGuardar.Texto = row[1]
								if err := daoCasosPracticos.InsertCasoPractico(casoPracticoAGuardar); err != nil {
									helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'CasosPracticos'"})
									return
								}
	
								casos_practicos_del_examen = append(casos_practicos_del_examen, casoPracticoAGuardar.Id)
	
								lastId = casoPracticoAGuardar.Id
	
							}
						}
					}

					preguntaCP, err := daoPreguntasCP.FindPreguntaByPregunta(row[3])

					if err != nil {
						if err.Error() == "not found" {
							
							var preguntaAGuardar models.PreguntasCP
						preguntaAGuardar.Id = bson.NewObjectId()

						preguntaAGuardar.Pregunta = row[3]
						preguntaAGuardar.Explicacion = row[6]

						respuestas := strings.Split(strings.TrimSpace(row[4]), "\n")


						a := []rune(strings.ToUpper(row[5]))

						letter := string(a[0:1])

						// log.Println(respuestas)

						if len(respuestas) == 4 {
							// log.Println("##############################")
							var contR = 0
							for _, respuesta := range respuestas {
								var respuestaAGuardar models.RespuestasCP
								respuestaAGuardar.Id = bson.NewObjectId()
									respuestaAGuardar.Respuesta = respuesta[3:]
									// log.Println(respuesta)
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
								// log.Println(respuestaAGuardar)
								if err := daoRespuestasCP.InsertRespuesta(respuestaAGuardar); err != nil {
									helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'RespuestasCP'"})
									return
								}
								contR = contR + 1
							}

							if err := daoPreguntasCP.InsertPregunta(preguntaAGuardar); err != nil {
								helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'PreguntasCP'"})
								return
							}
							

							var examenOficialPreguntaAGuardar models.CasosPracticosPreguntas
							examenOficialPreguntaAGuardar.Id = bson.NewObjectId()
							examenOficialPreguntaAGuardar.IdCasoPractico = lastId
							examenOficialPreguntaAGuardar.IdPregunta = preguntaAGuardar.Id
							examenOficialPreguntaAGuardar.Orden = cont


							if err := daoCasosPracticosPreguntas.InsertCasoPracticoPregunta(examenOficialPreguntaAGuardar); err != nil {
								helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'CasosPracticosPreguntas'"})
								return
							}

						}


						}
					}else {
						cp, err := daoCasosPracticos.FindCasoPracticoByName(sheet)

						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'CasosPracticosPreguntas'"})
							return
						}

						if err := daoCasosPracticosPreguntas.FindCasoPracticoByCasoPracticoPregunta(cp.Id.Hex(), preguntaCP.Id.Hex()); err != nil {
							if err.Error() == "not found" {

								var examenOficialPreguntaAGuardar models.CasosPracticosPreguntas
								examenOficialPreguntaAGuardar.Id = bson.NewObjectId()
								examenOficialPreguntaAGuardar.IdCasoPractico = cp.Id
							examenOficialPreguntaAGuardar.IdPregunta = preguntaCP.Id
							examenOficialPreguntaAGuardar.Orden = cont
							
							
							if err := daoCasosPracticosPreguntas.InsertCasoPracticoPregunta(examenOficialPreguntaAGuardar); err != nil {
								log.Println(err.Error())
								helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'CasosPracticosPreguntas'"})
								return
							}
							log.Println("Se creo una nueva relacion de caso practico pregunta")
						}
						}

							
							
					}
				

					
				}

			}
			cont = cont + 1
		}

		examenCasoPracticoAGuardar.CasosPracticos = casos_practicos_del_examen
		

		ex2, err := daoExamenesCasosPracticos.FindExamenCasoPracticoByName(ex.Name)

		if err != nil {
			if err.Error() == "not found" {
				log.Println("Creando nuevo examen")
				if err := daoExamenesCasosPracticos.UpdateExamenCasoPractico(examenCasoPracticoAGuardar); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'ExamenesCasosPracticos'"})
					return
				}
			}
		}else {	
			log.Println("El examen " + ex2.Name + "ya existe")
		}
	}
}

func DoImportCPOficiales(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

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
			if cont > 1 && len(row) >= 7 {
				if row[3] != "" && row[4] != "" && row[5] != "" && row[6] != "" {

					respuestas := strings.Split(row[4], "\n")
					// respuestas := strings.Split(row[6], "\n")

					// a := []rune(row[5])
					// letter := string(a[0:1])
					a := []rune(strings.ToUpper(row[5]))

					letter := string(a[0:1])

					var alguna_true = false

					if len(respuestas) >= 4 {
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

		config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/testError")
		if err != nil {
			log.Println(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		//#########################################
		return
	}

	cont = 0
	for _, sheet := range sheets {

		var casos_practicos_del_examen = []bson.ObjectId{}

		var examenCasoPracticoAGuardar models.ExamenesCasosPracticos
		examenCasoPracticoAGuardar.Id = bson.NewObjectId()
		examenCasoPracticoAGuardar.Name = sheet
		examenCasoPracticoAGuardar.Oficial = true
		examenCasoPracticoAGuardar.CasosPracticos = casos_practicos_del_examen
		if err := daoExamenesCasosPracticos.InsertExamenCasoPractico(examenCasoPracticoAGuardar); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'ExamenesCasosPracticos'"})
			return
		}

		cont = 0
		rows, err := fileRead.GetRows(sheet)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var lastId = bson.NewObjectId()

		for _, row := range rows {

			if cont >= 1 && len(row) >= 7 {

				if row[3] != "" && row[4] != "" && row[5] != "" && row[6] != "" {

					if row[0] != "" {
						var casoPracticoAGuardar models.CasosPracticos
						casoPracticoAGuardar.Id = bson.NewObjectId()
						casoPracticoAGuardar.Name = row[0]
						casoPracticoAGuardar.Oficial = true
						casoPracticoAGuardar.AnioOficial = sheet
						casoPracticoAGuardar.Texto = row[1]

						if err := daoCasosPracticos.InsertCasoPractico(casoPracticoAGuardar); err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'CasosPracticos'"})
							return
						}

						casos_practicos_del_examen = append(casos_practicos_del_examen, casoPracticoAGuardar.Id)

						lastId = casoPracticoAGuardar.Id

					}

					var preguntaAGuardar models.PreguntasCP
					preguntaAGuardar.Id = bson.NewObjectId()

					preguntaAGuardar.Pregunta = row[3]
					preguntaAGuardar.Explicacion = row[6]

					respuestas := strings.Split(row[4], "\n")

					// a := []rune(row[5])
					// letter := string(a[0:1])

					a := []rune(strings.ToUpper(row[5]))

					letter := string(a[0:1])

					if len(respuestas) == 4 {
						var contR = 0
						for _, respuesta := range respuestas {
							var respuestaAGuardar models.RespuestasCP
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
							if err := daoRespuestasCP.InsertRespuesta(respuestaAGuardar); err != nil {
								helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'RespuestasCP'"})
								return
							}
							contR = contR + 1
						}

						if err := daoPreguntasCP.InsertPregunta(preguntaAGuardar); err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'PreguntasCP'"})
							return
						}

						var examenOficialPreguntaAGuardar models.CasosPracticosPreguntas
						examenOficialPreguntaAGuardar.Id = bson.NewObjectId()
						examenOficialPreguntaAGuardar.IdCasoPractico = lastId
						examenOficialPreguntaAGuardar.IdPregunta = preguntaAGuardar.Id
						examenOficialPreguntaAGuardar.Orden = cont
						if err := daoCasosPracticosPreguntas.InsertCasoPracticoPregunta(examenOficialPreguntaAGuardar); err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'CasosPracticosPreguntas'"})
							return
						}

					}
				}

			}
			cont = cont + 1
		}

		examenCasoPracticoAGuardar.CasosPracticos = casos_practicos_del_examen
		if err := daoExamenesCasosPracticos.UpdateExamenCasoPractico(examenCasoPracticoAGuardar); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'ExamenesCasosPracticos'"})
			return
		}
	}

}

func DoImportTestNivel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
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
	var sheet = sheets[0]
	cont = 0
	rows, err := fileRead.GetRows(sheet)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	for _, row := range rows {
		if cont > 0 && len(row) >= 6 {
			if row[0] != "" && row[2] != "" && row[3] != "" && row[4] != "" && row[5] != "" {

				if row[0] != "" {
					_, err := daoNiveles.FindNivelByAbreviacion(row[0])
					if err != nil {
						errores_importacion = errores_importacion + "Nivel '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					}
				} else {
					errores_importacion = errores_importacion + "Nivel no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
				}

				respuestas := strings.Split(row[3], "\n")

				var alguna_true = false

				// a := []rune(row[4])
				// letter := string(a[0:1])

				a := []rune(strings.ToUpper(row[4]))

				letter := string(a[0:1])

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

	if errores_importacion != "" {
		//#########################################

		config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/testError")
		if err != nil {
			log.Println(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		//#########################################
		return
	}

	var preguntas_mezcladas []models.Preguntas

	cont = 0
	for _, row := range rows {
		if cont > 0 && len(row) >= 6 {
			if row[0] != "" && row[2] != "" && row[3] != "" && row[4] != "" && row[5] != "" {
				var preguntaAGuardar models.Preguntas
				preguntaAGuardar.Id = bson.NewObjectId()
				result, err := daoTemas.FindBasicoConfundidoByAbreviacion("NO")
				if err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Tema 'NO' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
					return
				}
				preguntaAGuardar.IdTema = result.Id
				if row[0] != "" {
					result, err := daoNiveles.FindNivelByAbreviacion(row[0])
					if err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
						return
					}
					preguntaAGuardar.IdNivel = result.Id
				} else {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
					return
				}

				preguntaAGuardar.Pregunta = row[2]
				preguntaAGuardar.Explicacion = row[5]
				preguntaAGuardar.Oficial = false
				preguntaAGuardar.AnioOficial = ""

				respuestas := strings.Split(row[3], "\n")

				// a := []rune(row[4])
				// letter := string(a[0:1])
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
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
						return
					}

					preguntas_mezcladas = append(preguntas_mezcladas, preguntaAGuardar)
				} else {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "No hay 4 respuestas en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
					return
				}
			}
		}
		cont = cont + 1
	}

	var historial_test models.Historiales

	historial, err := daoHistoriales.FindHistorialById("5f23fafa1e881b248f718097")
	if err != nil {
		var ids_temas []string
		var ids_basicosconfundidos []string
		var ids_oficiales []string
		var ids_legislaciones []string
		var ids_simulacros []string
		var historial models.Historiales
		historial.Id = bson.ObjectIdHex("5f23fafa1e881b248f718097")
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
		historial.IdUsuario = bson.ObjectIdHex("5ebb8cf5b8710000f6006502")
		historial.NumeroPreguntas = 50
		historial.RespuestaAutomatica = false
		historial.Tipo = "Test de nivel"
		historial.TiempoTranscurrido = 0
		historial.Terminado = false
		historial.PreguntasTotales = 0
		historial.PreguntasAcertadas = 0
		historial.PreguntasFalladas = 0
		historial.PreguntasBlancas = 0
		historial.Puntuacion = 0.0
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}
		historial_test = historial
	} else {
		historial_test = historial
	}
	historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial_test.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}
	for _, historial_pregunta := range historiales_preguntas {
		daoHistorialesPreguntas.RemoveHistorialPregunta(historial_pregunta.Id.Hex())
	}

	if len(preguntas_mezcladas) > 0 {

		historial_test.PreguntasTotales = len(preguntas_mezcladas)
		historial_test.PreguntasBlancas = len(preguntas_mezcladas)
		if err := daoHistoriales.UpdateHistorial(historial_test); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

		var cont_preguntas = 0
		for _, pregunta := range preguntas_mezcladas {
			var historial_pregunta models.HistorialesPreguntas
			historial_pregunta.Id = bson.NewObjectId()
			historial_pregunta.IdHistorial = historial_test.Id
			historial_pregunta.IdPregunta = pregunta.Id

			if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
			cont_preguntas = cont_preguntas + 1
		}
	} else {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
		return
	}

}

//#########################################################################
//#########################################################################
//#########################################################################

func DoImportAreas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

		count := 0
		countP := 0

		var bloque models.Bloques

		var Area models.Areas

		areaName := strings.ToLower(sheet)

		if sheet != "niveles" && sheet != "BASICOS CONFUNDIDOS" && sheet != "legislaciones" {
			areaRes, areaErr := daoAreas.FindAreaByName(areaName)

			if areaErr != nil {
				AreaID := bson.NewObjectId()
				Area.Id = AreaID
				if areaName == "administracion del estado" {
					Area.Name = "derecho administrativo y función pública"
				} else {
					Area.Name = areaName
				}

				err := daoAreas.InsertArea(Area)
				if err != nil {
					errores_importacion = errores_importacion + "Ocurrio un error al subir el area '" + sheet + " de la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
				}
			} else {
				Area = areaRes
			}
		}
		
		for _, row := range rows {
			if cont >= 0 && len(row) >= 2 {
				if row[0] != "" && row[1] != "" {
					if row[0] == "" || row[1] == "" {
						errores_importacion = errores_importacion + "No se encontro tema o nombre de tema en la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
					} else {
						if sheet == "niveles" {
							_, err := daoNiveles.FindNivelByAbreviacion(row[0])
							if err != nil {
								var nivel models.Niveles
								if row[0] == "B" {
									nivel.Id = bson.ObjectIdHex("5ea6ad3dbb5c000045007637")
								}
								if row[0] == "I" {
									nivel.Id = bson.ObjectIdHex("5ea6ad4abb5c000045007638")
								}
								if row[0] == "A" {
									nivel.Id = bson.ObjectIdHex("5ea6ad51bb5c000045007639")
								}
								nivel.Abreviacion = row[0]
								nivel.Name = row[1]

								nivelErr := daoNiveles.InsertNivel(nivel)
								if nivelErr != nil {
									errores_importacion = errores_importacion + "Ocurrio un error al subir el Nivel '" + row[1] + " de la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
								}
							}
						} else {
							if sheet == "BASICOS CONFUNDIDOS" {
								_, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(row[1])
								if err != nil {
									var basico models.BasicosConfundidos

									basico.Abreviacion = row[1]
									basico.Id = bson.NewObjectId()
									basico.Name = row[0]
									basico.IdArea = Area.Id

									basicoErrr := daoBasicosConfundidos.InsertBasicoConfundido(basico)
									if basicoErrr != nil {
										errores_importacion = errores_importacion + "Ocurrio un error al subir Basico Confundido '" + row[0] + "Con codigo " + row[1] + " de la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
									}
								}
							} else if sheet == "legislaciones" {

								_, err := daoLegislaciones.FindLegislacionByAbreviacion(row[0])
								if err != nil {
									var legislacion models.Legislaciones

									legislacion.Abreviacion = row[0]
									legislacion.Id = bson.NewObjectId()
									legislacion.Name = row[1]
									legislacion.IdArea = Area.Id

									legislacionErrr := daoLegislaciones.InsertLegislacion(legislacion)
									if legislacionErrr != nil {
										errores_importacion = errores_importacion + "Ocurrio un error al subir Basico Confundido '" + row[0] + "Con codigo " + row[1] + " de la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
									}
								}
							} else {

								if len(row) == 3 {

									_, err := daoLegislaciones.FindLegislacionByAbreviacion(row[1])
									if err != nil {

										var leg models.Legislaciones
										leg.Id = bson.NewObjectId()
										leg.Abreviacion = row[1]
										leg.Name = row[2]
										leg.IdArea = Area.Id

										legislacionErr := daoLegislaciones.InsertLegislacion(leg)
										if legislacionErr != nil {
											errores_importacion = errores_importacion + "Ocurrio un error al subir legislacion '" + row[1] + " de la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
										}
									}

								} else {

									res, _ := daoTemas.FindTemaByIdAreaAndAbreviacion(Area.Id.Hex(), row[0])

									if len(res) == 0 {
										abrevviacionP := extractPattern(row[0])

										var tema models.Temas
										tema.Id = bson.NewObjectId()
										tema.Abreviacion = row[0]
										tema.AbreviacionPublica = abrevviacionP
										tema.IdArea = Area.Id
										tema.Name = row[1]

										temaErr := daoTemas.InsertTema(tema)
										if temaErr != nil {
											errores_importacion = errores_importacion + "Ocurrio un error al subir el Tema '" + row[1] + " de la fila '" + strconv.Itoa(cont+1) + "de la hoja '" + sheet + "' <br>"
										}

										countP = countP + 1
										abr := extractPattern(row[0])

										if count < 2 {
											if count == 0 {
												bloque.Id = bson.NewObjectId()
												bloque.Name = "Temas " + abr
											}
											bloque.Temas = append(bloque.Temas, tema.Id)
											bloque.IdArea = Area.Id
											count = count + 1
										} else if count == 2 {
											bloque.Temas = append(bloque.Temas, tema.Id)
											bloque.IdArea = Area.Id
											count = 0
										}

										if countP == 3 || !(cont+1 < len(rows)){
											if !strings.Contains(bloque.Name, abr) {
													bloque.Name = bloque.Name + " - " + abr
											}
											err := daoBloques.InsertBloque(bloque)
											countP = 0
											bloque.Temas = []bson.ObjectId{}

											if err != nil {
												log.Println("OCURRIO UN ERROR AL INSERTAR UN BLOQUE linea 1513")
												log.Println(err)
												errores_importacion = errores_importacion + "Ocurrio un error al subir un bloque <br>"

											}
										}
										// if !(cont+1 < len(rows)){
										// 	bloque.Name = 
										// 	err := daoBloques.InsertBloque(bloque)
										// 	countP = 0
										// 	bloque.Temas = []bson.ObjectId{}

										// 	if err != nil {
										// 		log.Println("OCURRIO UN ERROR AL INSERTAR UN BLOQUE linea 1513")
										// 		log.Println(err)
										// 		errores_importacion = errores_importacion + "Ocurrio un error al subir un bloque <br>"

										// 	}
										// }

									}
								}
							}
						}
					}
				}
			} else if len(row) == 1 {
				areaN := strings.ToLower(row[0])
				if areaN == "administracion del estado" {
					areaN = "derecho administrativo y función pública"
				}
				areaRes, areaErr := daoAreas.FindAreaByName(areaN)
				if areaErr != nil {
					errores_importacion = errores_importacion + "Ocurrio un error al actualizar area para subir basico confundido " + areaErr.Error() + "'<br>"
					log.Println(areaErr)
				} else {
					Area = areaRes
				}
			}

			cont = cont + 1
		}
	}

	if errores_importacion != "" {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})
		return
	}
}

func DoImportBC(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
				if row[5] != "" && row[6] != "" && row[8] != "" {
					if row[0] != "" {
						_, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(sheet)
						// _, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(row[0])
						if err != nil {
							errores_importacion = errores_importacion + "Basico confundido '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					}

					if row[3] != "" {
						_, err := daoNiveles.FindNivelByAbreviacion(row[3])
						if err != nil {
							errores_importacion = errores_importacion + "Nivel'" + row[3] + "' no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
						}
					} else {
						errores_importacion = errores_importacion + "Nivel no encontrado en la fila '" + strconv.Itoa(cont+1) + "' de la hoja '" + sheet + "' <br>"
					}

					respuestas := strings.Split(row[6], "\n")

					// a := []rune(row[7])
					a := []rune(strings.ToUpper(row[7]))
					letter := string(a[0:1])

					var alguna_true = false

					if len(respuestas) >= 4 {
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

							// if contR == 0 && letter == "A" {
							// 	alguna_true = true
							// } else if contR == 1 && letter == "B" {
							// 	alguna_true = true
							// } else if contR == 2 && letter == "C" {
							// 	alguna_true = true
							// } else if contR == 3 && letter == "D" {
							// 	alguna_true = true
							// }

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

		config_vars := helper.GetConfigVars()

		res, err := http.Get(config_vars["URL"] + "email/testError")
		if err != nil {
			log.Println(err)
		}
		texto_email_bytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		texto_email := string(texto_email_bytes)
		texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)
		texto_email = strings.ReplaceAll(texto_email, "##PROBLEMAS##", errores_importacion)

		helper.SendEmail("PENITENCIARIOS.COM :: Error al subir el archivo", texto_email, importadores.Email, "penitenciarios@penitenciarios.com")

		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": errores_importacion})

		// 	//#########################################

		return
	}

	cont = 0
	for _, sheet := range sheets {
		cont = 0
		rows, err := fileRead.GetRows(sheet)
		if err != nil {
			log.Println("ERROR AQUI EN LA LINEA 1703 de importadores")
			helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		for _, row := range rows {
			if cont > 0 && len(row) >= 9 {
				if row[5] != "" && row[6] != "" && row[8] != "" {
					var preguntaAGuardar models.Preguntas
					preguntaAGuardar.Id = bson.NewObjectId()

					if row[0] != "" {
						result, err := daoBasicosConfundidos.FindBasicoConfundidoByAbreviacion(sheet)
						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Basico confundido asdasd '" + row[0] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdBasicosConfundidos = result.Id
					}

					if row[3] != "" {
						result, err := daoNiveles.FindNivelByAbreviacion(row[3])
						if err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel '" + row[3] + "' no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
							return
						}
						preguntaAGuardar.IdNivel = result.Id
					} else {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Nivel no encontrado en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
						return
					}

					preguntaAGuardar.Pregunta = row[5]
					preguntaAGuardar.Explicacion = row[8]
					preguntaAGuardar.Oficial = false
					preguntaAGuardar.AnioOficial = ""

					respuestas := strings.Split(row[6], "\n")

					// a := []rune(row[7])
					a := []rune(strings.ToUpper(row[7]))

					letter := string(a[0:1])
					// log.Println(respuestas)

					if len(respuestas) >= 4 {
						var contR = 0
						for _, respuesta := range respuestas[:4] {

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
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
							return
						}
						// helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "awdoianwodiwd'" + strconv.Itoa(cont) + "' awdawd '" + sheet + "'"})
						// return
					} else {

						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "No hay as 4 respuestas en la fila '" + strconv.Itoa(cont) + "' de la hoja '" + sheet + "'"})
						return
					}
				}
			}
			cont = cont + 1
		}
	}
	//#########################################

	config_vars := helper.GetConfigVars()

	result, err := http.Get(config_vars["URL"] + "email/testSuccess")
	if err != nil {
		log.Println(err)
	}
	texto_email_bytes, err := ioutil.ReadAll(result.Body)
	result.Body.Close()
	texto_email := string(texto_email_bytes)
	texto_email = strings.ReplaceAll(texto_email, "##TEST##", importadores.Name)

	helper.SendEmail("PENITENCIARIOS.COM :: Test Subido con exito", texto_email, importadores.Email, "inpenitenciariosfo@penitenciarios.com")

	//#########################################
}
