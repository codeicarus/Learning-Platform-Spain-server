package controllers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"test/helper"
	"test/models"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

var (
	daoPreguntasCP            = models.PreguntasCP{}
	daoHistorialesCP          = models.HistorialesCP{}
	daoHistorialesCPPreguntas = models.HistorialesCPPreguntas{}
	daoRespuestasCP           = models.RespuestasCP{}
)

func SavePreguntaCP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.SavePreguntaCP
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		if helper.CheckUser(usuario.Email) {

			pregunta, err := daoPreguntasCP.FindPreguntaById(params.Pregunta.Id.Hex())
			if err != nil {
			} else {
				pregunta.Pregunta = params.Pregunta.Pregunta
				pregunta.Explicacion = params.Pregunta.Explicacion

				if err := daoPreguntasCP.UpdatePregunta(pregunta); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}

				for _, respuesta_recibida := range params.Pregunta.Respuestas {
					respuesta, err := daoRespuestasCP.FindRespuestaById(respuesta_recibida.Id.Hex())
					if err != nil {
					} else {
						respuesta.Respuesta = respuesta_recibida.Respuesta
						respuesta.Correcta = respuesta_recibida.Correcta

						if err := daoRespuestasCP.UpdateRespuesta(respuesta); err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
							return
						}
					}
				}
			}

			helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})

		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "No tienes permisos"})
			return
		}
	}
}

func CreateTestCasosPracticos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.CreateTestCP
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_casos_practicos_string []string
		var id_caso_practico bson.ObjectId

		for _, examencasopractico := range params.CasosPracticos {
			examen_caso_practico, err := daoExamenesCasosPracticos.FindExamenCasoPracticoById(examencasopractico)
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
				return
			}

			for _, casopractico := range examen_caso_practico.CasosPracticos {
				ids_casos_practicos_string = append(ids_casos_practicos_string, casopractico.Hex())
			}
			id_caso_practico = examen_caso_practico.Id
		}

		var historial models.HistorialesCP
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.CasosPracticos = ids_casos_practicos_string
		historial.IdCasoPractico = id_caso_practico
		historial.IdUsuario = usuario.Id
		historial.NumeroPreguntas = params.NumeroPreguntas
		historial.RespuestaAutomatica = params.RespuestaAutomatica
		historial.Tipo = params.Tipo
		historial.TiempoTranscurrido = 0
		historial.Terminado = false
		historial.PreguntasTotales = 0
		historial.PreguntasAcertadas = 0
		historial.PreguntasFalladas = 0
		historial.PreguntasBlancas = 0
		historial.Puntuacion = 0.0
		if err := daoHistorialesCP.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		var total_total_preguntas = 0

		for _, casopractico := range ids_casos_practicos_string {
			var ids_casos_practicos []bson.ObjectId
			ids_casos_practicos = append(ids_casos_practicos, bson.ObjectIdHex(casopractico))

			var preguntas_mezcladas []models.PreguntasCP

			preguntas_oficiales, err := daoCasosPracticosPreguntas.FindPreguntaByIdCasoPracticoOrder(ids_casos_practicos)
			if err != nil {
			} else {
				for _, pregunta := range preguntas_oficiales {
					pregunta_real, err := daoPreguntasCP.FindPreguntaById(pregunta.IdPregunta.Hex())
					if err != nil {
					} else {
						preguntas_mezcladas = append(preguntas_mezcladas, pregunta_real)
					}
				}
			}

			log.Println(len(preguntas_oficiales))

			if len(preguntas_mezcladas) > 0 {

				var cont_preguntas = 0
				for _, pregunta := range preguntas_mezcladas {
					if cont_preguntas >= params.NumeroPreguntas {
					} else {
						var historial_pregunta models.HistorialesCPPreguntas
						historial_pregunta.Id = bson.NewObjectId()
						historial_pregunta.IdHistorial = historial.Id
						historial_pregunta.IdPregunta = pregunta.Id

						if err := daoHistorialesCPPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
							helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
							return
						}
						cont_preguntas = cont_preguntas + 1
						total_total_preguntas = total_total_preguntas + 1
					}
				}
			} else {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
				return
			}
		}

		if total_total_preguntas > params.NumeroPreguntas {
			historial.PreguntasTotales = params.NumeroPreguntas
			historial.PreguntasBlancas = params.NumeroPreguntas
		} else {
			historial.PreguntasTotales = total_total_preguntas
			historial.PreguntasBlancas = total_total_preguntas
		}
		if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}

func CambiarTiempoTranscurridoCP(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	historial, err := daoHistorialesCP.FindHistorialById(pregunta_a_corregir["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
		return
	}

	timepo_transcurrido, err := strconv.Atoi(pregunta_a_corregir["tiempo"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	historial.TiempoTranscurrido = timepo_transcurrido

	if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"status": "true"})
}

func GetTestCP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var historialADevolver models.HistorialesCPADevolver
	var supuestosADevolver []models.SupuestosCPADevolver
	var preguntasADevolverFinal []models.PreguntasCPADevolver
	var respuestas_marcadas = make(map[string]string)

	historial, err := daoHistorialesCP.FindHistorialById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return

	}

	historiales_preguntas, err := daoHistorialesCPPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}

	for _, historial_pregunta := range historiales_preguntas {
		pregunta, err := daoPreguntasCP.FindPreguntaById(historial_pregunta.IdPregunta.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta"})
			return
		}

		if historial_pregunta.IdRespuesta.Hex() != "" {
			respuestas_marcadas[pregunta.Id.Hex()] = historial_pregunta.IdRespuesta.Hex()
		}
	}

	for _, id_caso_practico := range historial.CasosPracticos {
		var preguntasADevolver []models.PreguntasCPADevolver
		var unSupuesto models.SupuestosCPADevolver
		caso_practico, err := daoCasosPracticos.FindCasoPracticoById(id_caso_practico)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
			return
		}

		var for_earch []bson.ObjectId
		for_earch = append(for_earch, bson.ObjectIdHex(id_caso_practico))

		historiales_preguntas, err := daoCasosPracticosPreguntas.FindPreguntaByIdCasoPracticoOrder(for_earch)
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
			return
		}

		for _, historial_pregunta := range historiales_preguntas {
			pregunta, err := daoPreguntasCP.FindPreguntaById(historial_pregunta.IdPregunta.Hex())
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta"})
				return
			}

			var respuestasADevolver []models.RespuestasCPADevolver

			respuestas, err := daoRespuestasCP.FindRespuestaByIdPregunta(pregunta.Id.Hex())
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "respuestas"})
				return
			}

			for _, respuesta := range respuestas {

				var respuestaADevolver models.RespuestasCPADevolver
				respuestaADevolver.Id = respuesta.Id
				respuestaADevolver.Respuesta = respuesta.Respuesta
				if helper.CheckUser(params["email"]) {
					respuestaADevolver.Correcta = respuesta.Correcta
				} else {
					respuestaADevolver.Correcta = false
				}
				respuestasADevolver = append(respuestasADevolver, respuestaADevolver)
			}

			rand.Seed(time.Now().UnixNano())
			for i := len(respuestasADevolver) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				respuestasADevolver[i], respuestasADevolver[j] = respuestasADevolver[j], respuestasADevolver[i]
			}

			var preguntaADevolver models.PreguntasCPADevolver
			preguntaADevolver.Id = pregunta.Id
			preguntaADevolver.Pregunta = pregunta.Pregunta
			if helper.CheckUser(params["email"]) {
				preguntaADevolver.Explicacion = pregunta.Explicacion
			} else {
				preguntaADevolver.Explicacion = ""
			}
			preguntaADevolver.Respuestas = respuestasADevolver
			preguntasADevolver = append(preguntasADevolver, preguntaADevolver)
			preguntasADevolverFinal = append(preguntasADevolverFinal, preguntaADevolver)
		}

		unSupuesto.Id = caso_practico.Id
		unSupuesto.Name = caso_practico.Name
		unSupuesto.Oficial = caso_practico.Oficial
		unSupuesto.AnioOficial = caso_practico.AnioOficial
		unSupuesto.Texto = caso_practico.Texto
		unSupuesto.Preguntas = preguntasADevolver

		supuestosADevolver = append(supuestosADevolver, unSupuesto)

	}

	historialADevolver.Historial = historial
	historialADevolver.CasosPracticos = supuestosADevolver
	historialADevolver.RespuestasMarcadas = respuestas_marcadas
	historialADevolver.Preguntas = preguntasADevolverFinal

	helper.ResponseWithJson(w, http.StatusOK, historialADevolver)
}

func CorregirPreguntaCP(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	historial, err := daoHistorialesCP.FindHistorialById(pregunta_a_corregir["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
		return
	}

	pregunta, err := daoPreguntasCP.FindPreguntaById(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	respuesta, err := daoRespuestasCP.FindRespuestaById(pregunta_a_corregir["id_respuesta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	respuesta_correcta, err := daoRespuestasCP.FindRespuestaCorrectaByIdPregunta(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	historial_pregunta, err := daoHistorialesCPPreguntas.FindHistorialPreguntaByIdHistorialIdPregunta(historial.Id.Hex(), pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Esa pregunta no se encuentra en ese historial"})
		return
	}

	timepo_transcurrido, err := strconv.Atoi(pregunta_a_corregir["tiempo"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	historial.TiempoTranscurrido = timepo_transcurrido

	if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	historial_pregunta.IdRespuesta = respuesta.Id

	if respuesta.Id == respuesta_correcta.Id {
		historial_pregunta.Correcta = true

		historial.PreguntasAcertadas = historial.PreguntasAcertadas + 1
		historial.PreguntasBlancas = historial.PreguntasBlancas - 1
		historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)
		if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

	} else {
		historial_pregunta.Correcta = false

		historial.PreguntasFalladas = historial.PreguntasFalladas + 1
		historial.PreguntasBlancas = historial.PreguntasBlancas - 1
		historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)
		if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}
	}

	if err := daoHistorialesCPPreguntas.UpdateHistorialPregunta(historial_pregunta); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	json_respuesta_correcta, err := json.Marshal(respuesta_correcta)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	json_pregunta, err := json.Marshal(pregunta)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"respuesta_correcta": string(json_respuesta_correcta), "pregunta": string(json_pregunta)})
}

func CorregirPreguntaCPYaContestada(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	historial, err := daoHistorialesCP.FindHistorialById(pregunta_a_corregir["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
		return
	}

	pregunta, err := daoPreguntasCP.FindPreguntaById(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	respuesta, err := daoRespuestasCP.FindRespuestaById(pregunta_a_corregir["id_respuesta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	respuesta_correcta, err := daoRespuestasCP.FindRespuestaCorrectaByIdPregunta(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	historial_pregunta, err := daoHistorialesCPPreguntas.FindHistorialPreguntaByIdHistorialIdPregunta(historial.Id.Hex(), pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Esa pregunta no se encuentra en ese historial"})
		return
	}

	historial_pregunta.IdRespuesta = respuesta.Id

	if respuesta.Id == respuesta_correcta.Id {
		historial_pregunta.Correcta = true
	} else {
		historial_pregunta.Correcta = false
	}

	json_respuesta_correcta, err := json.Marshal(respuesta_correcta)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	json_pregunta, err := json.Marshal(pregunta)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"respuesta_correcta": string(json_respuesta_correcta), "pregunta": string(json_pregunta)})
}

func CorregirTestCP(w http.ResponseWriter, r *http.Request) {

	var respuestas_marcadas map[string]string
	var respuestasCorrectas = make(map[string]models.RespuestasCP)
	var preguntas = make(map[string]models.PreguntasCP)

	var params models.CorregirTestCP
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	if err := json.Unmarshal([]byte(params.RespuestasMarcadas), &respuestas_marcadas); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
	}

	historial, err := daoHistorialesCP.FindHistorialById(params.IdTest)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}

	historial.PreguntasAcertadas = 0
	historial.PreguntasFalladas = 0
	historial.PreguntasBlancas = historial.PreguntasTotales

	historiales_preguntas, err := daoHistorialesCPPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}
	for _, historial_pregunta := range historiales_preguntas {

		pregunta, err := daoPreguntasCP.FindPreguntaById(historial_pregunta.IdPregunta.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta"})
			return
		}
		preguntas[pregunta.Id.Hex()] = pregunta

		respuesta_correcta, err := daoRespuestasCP.FindRespuestaCorrectaByIdPregunta(pregunta.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
			return
		}
		respuestasCorrectas[pregunta.Id.Hex()] = respuesta_correcta

		if respuesta_a_la_pregunta, ok := respuestas_marcadas[pregunta.Id.Hex()]; ok {

			//CONTESTADA
			historial_pregunta.IdRespuesta = bson.ObjectIdHex(respuesta_a_la_pregunta)

			if bson.ObjectIdHex(respuesta_a_la_pregunta) == respuesta_correcta.Id {
				historial_pregunta.Correcta = true

				historial.PreguntasAcertadas = historial.PreguntasAcertadas + 1
				historial.PreguntasBlancas = historial.PreguntasBlancas - 1

			} else {
				historial_pregunta.Correcta = false

				historial.PreguntasFalladas = historial.PreguntasFalladas + 1
				historial.PreguntasBlancas = historial.PreguntasBlancas - 1
			}

			if err := daoHistorialesCPPreguntas.UpdateHistorialPregunta(historial_pregunta); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

		} else {

			//NO CONTESTADA
		}

	}

	// historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)

	puntaje := float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) / 3)

	if puntaje < 0 {
		historial.Puntuacion = 0
	} else {
		historial.Puntuacion = puntaje
	}


	historial.Terminado = true
	historial.TiempoTranscurrido = params.Tiempo

	if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	json_preguntas, err := json.Marshal(preguntas)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	json_respuestasCorrectas, err := json.Marshal(respuestasCorrectas)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"respuestas_correctas": string(json_respuestasCorrectas), "preguntas": string(json_preguntas)})
}

func CorregirTestCPYaContestado(w http.ResponseWriter, r *http.Request) {

	var respuestas_marcadas map[string]string
	var respuestasCorrectas = make(map[string]models.RespuestasCP)
	var preguntas = make(map[string]models.PreguntasCP)

	var params models.CorregirTestCP
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	if err := json.Unmarshal([]byte(params.RespuestasMarcadas), &respuestas_marcadas); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
	}

	historial, err := daoHistorialesCP.FindHistorialById(params.IdTest)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}

	historial.PreguntasAcertadas = 0
	historial.PreguntasFalladas = 0
	historial.PreguntasBlancas = historial.PreguntasTotales

	historiales_preguntas, err := daoHistorialesCPPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}
	for _, historial_pregunta := range historiales_preguntas {

		pregunta, err := daoPreguntasCP.FindPreguntaById(historial_pregunta.IdPregunta.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta"})
			return
		}
		preguntas[pregunta.Id.Hex()] = pregunta

		respuesta_correcta, err := daoRespuestasCP.FindRespuestaCorrectaByIdPregunta(pregunta.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
			return
		}
		respuestasCorrectas[pregunta.Id.Hex()] = respuesta_correcta

		if respuesta_a_la_pregunta, ok := respuestas_marcadas[pregunta.Id.Hex()]; ok {

			//CONTESTADA
			historial_pregunta.IdRespuesta = bson.ObjectIdHex(respuesta_a_la_pregunta)

			if bson.ObjectIdHex(respuesta_a_la_pregunta) == respuesta_correcta.Id {
				historial_pregunta.Correcta = true

				historial.PreguntasAcertadas = historial.PreguntasAcertadas + 1
				historial.PreguntasBlancas = historial.PreguntasBlancas - 1

			} else {
				historial_pregunta.Correcta = false

				historial.PreguntasFalladas = historial.PreguntasFalladas + 1
				historial.PreguntasBlancas = historial.PreguntasBlancas - 1
			}

		} else {

			//NO CONTESTADA
		}

	}

	historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)
	historial.Terminado = true

	json_preguntas, err := json.Marshal(preguntas)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	json_respuestasCorrectas, err := json.Marshal(respuestasCorrectas)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"respuestas_correctas": string(json_respuestasCorrectas), "preguntas": string(json_preguntas)})
}

func HayErrorCP(w http.ResponseWriter, r *http.Request) {

	var datos_error = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&datos_error); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	pregunta, err := daoPreguntasCP.FindPreguntaById(datos_error["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(datos_error["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {
		historialRepot, err := daoHistorialesCP.FindHistorialById(datos_error["id_test"])

		if err != nil{
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
			return	
		}
		if !helper.SendEmail("PENITENCIARIOS.COM :: Hay error", "El usuario '"+usuario.Name+"' <br /> Email: '"+usuario.Email+"' <br /> ID '"+usuario.Id.Hex()+"' <br /> " +"Id Test: "+ datos_error["id_test"] +" <br />"+ "Tipo: "+ historialRepot.Tipo +"<br /><br />" +" Ha encontrado un error en la pregunta '"+pregunta.Id.Hex()+"'<br /><br /><b>DETALLES:</b>"+datos_error["detalles"]+" ", "penitenciarios@penitenciarios.com", usuario.Email) {

			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
			return
		}
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "true"})
}
func TieneDudaCP(w http.ResponseWriter, r *http.Request) {

	var datos_error = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&datos_error); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	pregunta, err := daoPreguntasCP.FindPreguntaById(datos_error["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(datos_error["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {
		if !helper.SendEmail("PENITENCIARIOS.COM :: Tiene duda", "El usuario '"+usuario.Id.Hex()+"' tiene una duda en la pregunta '"+pregunta.Id.Hex()+"'<br><br><b>DETALLES:</b>"+datos_error["detalles"]+" ", "penitenciarios@penitenciarios.com", usuario.Email) {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
			return
		}
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "true"})
}

func RepetirTestCP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {
		historial_a_copiar, err := daoHistorialesCP.FindHistorialById(params["id"])
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial_a_copiar"})
			return
		}
		if historial_a_copiar.IdUsuario != usuario.Id {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error historial no válido", "error2": "historial_a_copiar"})
			return
		}
		var historial models.HistorialesCP
		historial.Id = bson.NewObjectId()
		historial.Tipo = historial_a_copiar.Tipo
		historial.Fecha = helper.MakeTimestamp()
		historial.CasosPracticos = historial_a_copiar.CasosPracticos
		historial.IdUsuario = usuario.Id
		historial.NumeroPreguntas = historial_a_copiar.NumeroPreguntas
		historial.RespuestaAutomatica = historial_a_copiar.RespuestaAutomatica
		historial.TiempoTranscurrido = 0
		historial.Terminado = false
		historial.PreguntasTotales = 0
		historial.PreguntasAcertadas = 0
		historial.PreguntasFalladas = 0
		historial.PreguntasBlancas = 0
		historial.Puntuacion = 0.0
		historial.IdCasoPractico = historial_a_copiar.IdCasoPractico
		if err := daoHistorialesCP.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		historiales_preguntas_a_copiar, err := daoHistorialesCPPreguntas.FindHistorialPreguntaByIdHistorial(historial_a_copiar.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas_a_copiar"})
			return
		}

		var preguntas_mezcladas []models.PreguntasCP

		for _, historial_pregunta_a_copiar := range historiales_preguntas_a_copiar {
			pregunta_a_copiar, err := daoPreguntasCP.FindPreguntaById(historial_pregunta_a_copiar.IdPregunta.Hex())
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta_a_copiar"})
				return
			} else {
				preguntas_mezcladas = append(preguntas_mezcladas, pregunta_a_copiar)
			}
		}

		rand.Seed(time.Now().UnixNano())
		for i := len(preguntas_mezcladas) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			preguntas_mezcladas[i], preguntas_mezcladas[j] = preguntas_mezcladas[j], preguntas_mezcladas[i]
		}
		if len(preguntas_mezcladas) > 0 {

			historial.PreguntasTotales = len(preguntas_mezcladas)
			historial.PreguntasBlancas = len(preguntas_mezcladas)
			if err := daoHistorialesCP.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "UpdateHistorial"})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				var historial_pregunta models.HistorialesCPPreguntas
				historial_pregunta.Id = bson.NewObjectId()
				historial_pregunta.IdHistorial = historial.Id
				historial_pregunta.IdPregunta = pregunta.Id

				if err := daoHistorialesCPPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "InsertHistorialPregunta"})
					return
				}
				cont_preguntas = cont_preguntas + 1

			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}
