package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	daoPreguntasWRespuestas = models.PreguntasWRespuestas{}
	daoPreguntas            = models.Preguntas{}
	daoHistoriales          = models.Historiales{}
	daoHistorialesPreguntas = models.HistorialesPreguntas{}
	daoRespuestas           = models.Respuestas{}
)

func SavePregunta(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.SavePregunta
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
			pregunta, err := daoPreguntas.FindPreguntaById(params.Pregunta.Id.Hex())
			if err != nil {
			} else {
				pregunta.Pregunta = params.Pregunta.Pregunta
				pregunta.Explicacion = params.Pregunta.Explicacion
				pregunta.IdNivel = bson.ObjectIdHex(params.Pregunta.Nivel)

				if err := daoPreguntas.UpdatePregunta(pregunta); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}

				for _, respuesta_recibida := range params.Pregunta.Respuestas {
					respuesta, err := daoRespuestas.FindRespuestaById(respuesta_recibida.Id.Hex())
					if err != nil {
					} else {
						respuesta.Respuesta = respuesta_recibida.Respuesta
						respuesta.Correcta = respuesta_recibida.Correcta

						if err := daoRespuestas.UpdateRespuesta(respuesta); err != nil {
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
func CreateTestFalladasBlancasGuardadas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_legislaciones []string
		var ids_oficiales []string
		var ids_simulacros []string
		var ids_basicosconfundidos []string
		var ids_temas []string
		var historial models.Historiales

		respuestaAutomatica := false

		if params["respuesta_automatica"] == "1" {
			respuestaAutomatica = true
		}

		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
		historial.IdUsuario = usuario.Id
		historial.NumeroPreguntas = 0
		historial.RespuestaAutomatica = respuestaAutomatica
		historial.Tipo = params["tipo"]
		historial.TiempoTranscurrido = 0
		historial.Terminado = false
		historial.PreguntasTotales = 0
		historial.PreguntasAcertadas = 0
		historial.PreguntasFalladas = 0
		historial.PreguntasBlancas = 0
		historial.Puntuacion = 0.0

		numero_preguntas, _ := strconv.Atoi(params["numero_preguntas"])

		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		var ids_historiales []bson.ObjectId
		historiales_usuario, err := daoHistoriales.FindHistorialByIdUsuario(usuario.Id.Hex(), bson.M{"_id": 1, "terminado": 1})
		if err != nil {
		} else {
			for _, historial_usuario := range historiales_usuario {
				if historial_usuario.Terminado {
					ids_historiales = append(ids_historiales, historial_usuario.Id)
				}
			}
		}

		var ids_preguntas_historial_final []bson.ObjectId

		var ids_preguntas_historiales_usuario_recorridos []bson.ObjectId

		var ids_preguntas_guardadas_historiales_usuario []bson.ObjectId
		preguntas_favoritas_usuario, err := daoPreguntasFavoritas.FindAllPreguntasFavoritasByIdUsuario(usuario.Id.Hex())
		if err != nil {
		} else {
			for _, pregunta_favorita_usuario := range preguntas_favoritas_usuario {

				exists, _ := helper.InArray(pregunta_favorita_usuario.IdPregunta, ids_preguntas_historiales_usuario_recorridos)
				if exists {
				} else {
					ids_preguntas_guardadas_historiales_usuario = append(ids_preguntas_guardadas_historiales_usuario, pregunta_favorita_usuario.IdPregunta)
					ids_preguntas_historiales_usuario_recorridos = append(ids_preguntas_historiales_usuario_recorridos, pregunta_favorita_usuario.IdPregunta)
				}
			}
		}

		rand.Seed(time.Now().UnixNano())
		for i := len(ids_preguntas_guardadas_historiales_usuario) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			ids_preguntas_guardadas_historiales_usuario[i], ids_preguntas_guardadas_historiales_usuario[j] = ids_preguntas_guardadas_historiales_usuario[j], ids_preguntas_guardadas_historiales_usuario[i]
		}

		var ids_preguntas_blancas_historiales_usuario []bson.ObjectId
		var ids_preguntas_falladas_historiales_usuario []bson.ObjectId

		preguntas_historiales_usuario, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdsHistoriales(ids_historiales)
		if err != nil {
		} else {
			for _, pregunta_historial_usuario := range preguntas_historiales_usuario {

				exists, _ := helper.InArray(pregunta_historial_usuario.IdPregunta, ids_preguntas_historiales_usuario_recorridos)
				if exists {
				} else {
					if pregunta_historial_usuario.IdRespuesta == "" { //BLANCA
						ids_preguntas_blancas_historiales_usuario = append(ids_preguntas_blancas_historiales_usuario, pregunta_historial_usuario.IdPregunta)
					} else if pregunta_historial_usuario.Correcta != true { //FALLADA
						ids_preguntas_falladas_historiales_usuario = append(ids_preguntas_falladas_historiales_usuario, pregunta_historial_usuario.IdPregunta)
					}
					ids_preguntas_historiales_usuario_recorridos = append(ids_preguntas_historiales_usuario_recorridos, pregunta_historial_usuario.IdPregunta)
				}
			}
		}

		rand.Seed(time.Now().UnixNano())
		for i := len(ids_preguntas_blancas_historiales_usuario) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			ids_preguntas_blancas_historiales_usuario[i], ids_preguntas_blancas_historiales_usuario[j] = ids_preguntas_blancas_historiales_usuario[j], ids_preguntas_blancas_historiales_usuario[i]
		}
		rand.Seed(time.Now().UnixNano())
		for i := len(ids_preguntas_falladas_historiales_usuario) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			ids_preguntas_falladas_historiales_usuario[i], ids_preguntas_falladas_historiales_usuario[j] = ids_preguntas_falladas_historiales_usuario[j], ids_preguntas_falladas_historiales_usuario[i]
		}

		porcentaje_extra := 0

		numero_slice_blancas := 0
		porcentaje_blancas := 0
		quedan_blancas := false

		if params["blancas"] == "1" && len(ids_preguntas_blancas_historiales_usuario) > 0 {
			porcentaje_blancas_que_tenemos := 100 / numero_preguntas * len(ids_preguntas_blancas_historiales_usuario)

			if params["falladas"] == "1" && params["guardadas"] == "1" {
				porcentaje_blancas = 34
			} else if params["falladas"] == "0" && params["guardadas"] == "0" {
				porcentaje_blancas = 101
			} else {
				porcentaje_blancas = 51
			}
			if porcentaje_blancas_que_tenemos < porcentaje_blancas {
				porcentaje_blancas = porcentaje_blancas_que_tenemos
				porcentaje_extra = porcentaje_blancas - porcentaje_blancas_que_tenemos
				for _, id_pregunta := range ids_preguntas_blancas_historiales_usuario {
					ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
				}
			} else {
				quedan_blancas = true
				numero_slice_blancas = int(math.RoundToEven(float64(numero_preguntas) / 100.0 * float64(porcentaje_blancas)))
				for _, id_pregunta := range ids_preguntas_blancas_historiales_usuario[0:numero_slice_blancas] {
					ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
				}
			}

		}
		numero_slice_falladas := 0
		porcentaje_falladas := 0
		quedan_falladas := false
		if params["falladas"] == "1" && len(ids_preguntas_falladas_historiales_usuario) > 0 {
			porcentaje_falladas_que_tenemos := 100 / numero_preguntas * len(ids_preguntas_falladas_historiales_usuario)

			if params["blancas"] == "1" && params["guardadas"] == "1" {
				porcentaje_falladas = 34 + (porcentaje_extra / 2)
			} else if params["blancas"] == "0" && params["guardadas"] == "0" {
				porcentaje_falladas = 101
			} else {
				porcentaje_falladas = 51 + porcentaje_extra
			}
			if porcentaje_falladas_que_tenemos < porcentaje_falladas {
				porcentaje_falladas = porcentaje_falladas_que_tenemos
				porcentaje_extra = porcentaje_falladas - porcentaje_falladas_que_tenemos
				for _, id_pregunta := range ids_preguntas_falladas_historiales_usuario {
					ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
				}
			} else {
				quedan_falladas = true
				numero_slice_falladas = int(math.RoundToEven(float64(numero_preguntas) / 100.0 * float64(porcentaje_falladas)))
				for _, id_pregunta := range ids_preguntas_falladas_historiales_usuario[0:numero_slice_falladas] {
					ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
				}
			}

		}
		numero_slice_guardadas := 0
		porcentaje_guardadas := 0
		quedan_guardadas := false
		if params["guardadas"] == "1" && len(ids_preguntas_guardadas_historiales_usuario) > 0 {
			porcentaje_guardadas_que_tenemos := 100 / numero_preguntas * len(ids_preguntas_guardadas_historiales_usuario)

			if params["blancas"] == "1" && params["guardadas"] == "1" {
				porcentaje_guardadas = 34 + (porcentaje_extra / 2)
			} else if params["blancas"] == "0" && params["guardadas"] == "0" {
				porcentaje_guardadas = 101
			} else {
				porcentaje_guardadas = 51 + porcentaje_extra
			}
			if porcentaje_guardadas_que_tenemos < porcentaje_guardadas {
				porcentaje_guardadas = porcentaje_guardadas_que_tenemos
				porcentaje_extra = porcentaje_guardadas - porcentaje_guardadas_que_tenemos
				for _, id_pregunta := range ids_preguntas_guardadas_historiales_usuario {
					ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
				}
			} else {
				quedan_guardadas = true
				numero_slice_guardadas = int(math.RoundToEven(float64(numero_preguntas) / 100.0 * float64(porcentaje_guardadas)))
				for _, id_pregunta := range ids_preguntas_guardadas_historiales_usuario[0:numero_slice_guardadas] {
					ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
				}
			}

		}

		if len(ids_preguntas_historial_final) < numero_preguntas { /* REAJUSTE */
			numero_preguntas_faltan := numero_preguntas - len(ids_preguntas_historial_final)
			if params["blancas"] == "1" && quedan_blancas && numero_preguntas_faltan > 0 {
				if len(ids_preguntas_blancas_historiales_usuario)-numero_slice_blancas >= numero_preguntas_faltan {
					for _, id_pregunta := range ids_preguntas_blancas_historiales_usuario[numero_slice_blancas:(numero_slice_blancas + numero_preguntas_faltan)] {
						ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
						numero_preguntas_faltan--
					}
				} else {
					for _, id_pregunta := range ids_preguntas_blancas_historiales_usuario[numero_slice_blancas:(len(ids_preguntas_blancas_historiales_usuario) - 1)] {
						ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
						numero_preguntas_faltan--
					}
				}
			}
			if params["falladas"] == "1" && quedan_falladas && numero_preguntas_faltan > 0 {
				if len(ids_preguntas_falladas_historiales_usuario)-numero_slice_falladas >= numero_preguntas_faltan {
					for _, id_pregunta := range ids_preguntas_falladas_historiales_usuario[numero_slice_falladas:(numero_slice_falladas + numero_preguntas_faltan)] {
						ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
						numero_preguntas_faltan--
					}
				} else {
					for _, id_pregunta := range ids_preguntas_falladas_historiales_usuario[numero_slice_falladas:(len(ids_preguntas_falladas_historiales_usuario) - 1)] {
						ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
						numero_preguntas_faltan--
					}
				}
			}
			if params["guardadas"] == "1" && quedan_guardadas && numero_preguntas_faltan > 0 {
				if len(ids_preguntas_guardadas_historiales_usuario)-numero_slice_guardadas >= numero_preguntas_faltan {
					for _, id_pregunta := range ids_preguntas_guardadas_historiales_usuario[numero_slice_guardadas:(numero_slice_guardadas + numero_preguntas_faltan)] {
						ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
						numero_preguntas_faltan--
					}
				} else {
					for _, id_pregunta := range ids_preguntas_guardadas_historiales_usuario[numero_slice_guardadas:(len(ids_preguntas_guardadas_historiales_usuario) - 1)] {
						ids_preguntas_historial_final = append(ids_preguntas_historial_final, id_pregunta)
						numero_preguntas_faltan--
					}
				}
			}
		}

		if len(ids_preguntas_historial_final) > 0 {

			rand.Seed(time.Now().UnixNano())
			for i := len(ids_preguntas_historial_final) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				ids_preguntas_historial_final[i], ids_preguntas_historial_final[j] = ids_preguntas_historial_final[j], ids_preguntas_historial_final[i]
			}

			var cont_preguntas = 0
			for _, pregunta := range ids_preguntas_historial_final {
				if cont_preguntas <= numero_preguntas {
					var historial_pregunta models.HistorialesPreguntas
					historial_pregunta.Id = bson.NewObjectId()
					historial_pregunta.IdHistorial = historial.Id
					historial_pregunta.IdPregunta = pregunta

					if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
						return
					}
					cont_preguntas = cont_preguntas + 1
				}
			}

			historial.NumeroPreguntas = cont_preguntas
			historial.PreguntasTotales = cont_preguntas
			historial.PreguntasBlancas = cont_preguntas

			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}
func GetDataFalladasBlancasGuardadas(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	{
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_historiales []bson.ObjectId
		historiales_usuario, err := daoHistoriales.FindHistorialByIdUsuario(usuario.Id.Hex(), bson.M{"_id": 1, "terminado": 1})
		if err != nil {
		} else {
			for _, historial_usuario := range historiales_usuario {
				if historial_usuario.Terminado {
					ids_historiales = append(ids_historiales, historial_usuario.Id)
				}
			}
		}

		var ids_preguntas_historiales_usuario_recorridos []bson.ObjectId

		var ids_preguntas_guardadas_historiales_usuario []bson.ObjectId
		preguntas_favoritas_usuario, err := daoPreguntasFavoritas.FindAllPreguntasFavoritasByIdUsuario(usuario.Id.Hex())
		if err != nil {
		} else {
			for _, pregunta_favorita_usuario := range preguntas_favoritas_usuario {

				exists, _ := helper.InArray(pregunta_favorita_usuario.IdPregunta, ids_preguntas_historiales_usuario_recorridos)
				if exists {
				} else {
					ids_preguntas_guardadas_historiales_usuario = append(ids_preguntas_guardadas_historiales_usuario, pregunta_favorita_usuario.IdPregunta)
					ids_preguntas_historiales_usuario_recorridos = append(ids_preguntas_historiales_usuario_recorridos, pregunta_favorita_usuario.IdPregunta)
				}
			}
		}

		rand.Seed(time.Now().UnixNano())
		for i := len(ids_preguntas_guardadas_historiales_usuario) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			ids_preguntas_guardadas_historiales_usuario[i], ids_preguntas_guardadas_historiales_usuario[j] = ids_preguntas_guardadas_historiales_usuario[j], ids_preguntas_guardadas_historiales_usuario[i]
		}

		var ids_preguntas_blancas_historiales_usuario []bson.ObjectId
		var ids_preguntas_falladas_historiales_usuario []bson.ObjectId

		preguntas_historiales_usuario, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdsHistoriales(ids_historiales)
		if err != nil {
		} else {
			for _, pregunta_historial_usuario := range preguntas_historiales_usuario {

				exists, _ := helper.InArray(pregunta_historial_usuario.IdPregunta, ids_preguntas_historiales_usuario_recorridos)
				if exists {
				} else {
					if pregunta_historial_usuario.IdRespuesta == "" { //BLANCA
						ids_preguntas_blancas_historiales_usuario = append(ids_preguntas_blancas_historiales_usuario, pregunta_historial_usuario.IdPregunta)
					} else if pregunta_historial_usuario.Correcta != true { //FALLADA
						ids_preguntas_falladas_historiales_usuario = append(ids_preguntas_falladas_historiales_usuario, pregunta_historial_usuario.IdPregunta)
					}
					ids_preguntas_historiales_usuario_recorridos = append(ids_preguntas_historiales_usuario_recorridos, pregunta_historial_usuario.IdPregunta)
				}
			}
		}

		helper.ResponseWithJson(w, http.StatusOK, map[string]string{
			"result":    "success",
			"falladas":  strconv.Itoa(len(ids_preguntas_falladas_historiales_usuario)),
			"blancas":   strconv.Itoa(len(ids_preguntas_blancas_historiales_usuario)),
			"guardadas": strconv.Itoa(len(ids_preguntas_guardadas_historiales_usuario)),
		})
	}
}
func CreateTestNivel(w http.ResponseWriter, r *http.Request) {
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
		historial_a_copiar, err := daoHistoriales.FindHistorialById("5f23fafa1e881b248f718097")
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial_a_copiar"})
			return
		}

		var ids_temas []string
		var ids_legislaciones []string
		var ids_oficiales []string
		var ids_simulacros []string
		var ids_basicosconfundidos []string
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
		historial.IdUsuario = usuario.Id
		historial.NumeroPreguntas = historial_a_copiar.NumeroPreguntas
		historial.RespuestaAutomatica = historial_a_copiar.RespuestaAutomatica
		historial.Tipo = historial_a_copiar.Tipo
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

		historiales_preguntas_a_copiar, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial_a_copiar.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas_a_copiar"})
			return
		}

		var preguntas_mezcladas []models.Preguntas

		for _, historial_pregunta_a_copiar := range historiales_preguntas_a_copiar {
			pregunta_a_copiar, err := daoPreguntas.FindPreguntaById(historial_pregunta_a_copiar.IdPregunta.Hex())
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
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "UpdateHistorial"})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				var historial_pregunta models.HistorialesPreguntas
				historial_pregunta.Id = bson.NewObjectId()
				historial_pregunta.IdHistorial = historial.Id
				historial_pregunta.IdPregunta = pregunta.Id

				if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
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
func CreateTest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.CreateTest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {
		now := time.Now()
		year, month, day := now.Date()
		hoy_unix := time.Date(year, month, day, 0, 0, 0, 0, now.Location()).Unix()
		if usuario.MaximoDiaSuscripcion <= hoy_unix {

			historiales_hoy, err := daoHistoriales.FindHistorialByIdUsuarioAndDate(usuario.Id.Hex(), hoy_unix, nil)
			if err != nil {
			} else {
				if len(historiales_hoy) > 0 {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Ya has realizado un test hoy. Suscríbete para poder tener el acceso completo a la plataforma"})
					return
				}
			}
		}

		var ids_legislaciones []string
		var ids_oficiales []string
		var ids_simulacros []string
		var ids_basicosconfundidos []string
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = params.Temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
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
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		// var ids_temas []bson.ObjectId
		var preguntas_mezcladas []models.Preguntas
		// var ids_preguntas_seleccionadas []bson.ObjectId

		if len(params.Temas) == 1 {
			preguntas_bbdd, err := daoPreguntas.FindRandomPreguntasByTemas(params.NumeroPreguntas, params.Temas[0])
			if err != nil {
			} else {

				for _, pregunta := range preguntas_bbdd {
					preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
				}
			}
		} else {

			num_preguntas := params.NumeroPreguntas/len(params.Temas) + 5

			for _, tema := range params.Temas {
				preguntas_bbdd, err := daoPreguntas.FindRandomPreguntasByTemas(num_preguntas, tema)
				// log.Println(len(preguntas_bbdd))
				if err != nil {
				} else {
					for _, pregunta := range preguntas_bbdd {
						preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
					}
				}
			}
		}

		//#####################
		//#####################
		//#####################
		//#####################

		// for _, tema := range params.Temas {
		// 	ids_temas = append(ids_temas, bson.ObjectIdHex(tema))

		// 	var ids_temas_a_buscar []bson.ObjectId
		// 	ids_temas_a_buscar = append(ids_temas_a_buscar, bson.ObjectIdHex(tema))

		// if usuario.IdNivel.Hex() == "5ea6ad3dbb5c000045007637" { //BASICO

		// var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 80.0 / float64(len(params.Temas))))
		// preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)

		// 	if err != nil {
		// 		log.Println("ERROR")
		// 	} else {
		// 		log.Println(preguntas_basicas)
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad4abb5c000045007638" { //Intermedio
		// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 65.0 / float64(len(params.Temas))))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 			ids_preguntas_seleccionadas = append(ids_preguntas_seleccionadas, pregunta.Id)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 25.0 / float64(len(params.Temas))))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 			ids_preguntas_seleccionadas = append(ids_preguntas_seleccionadas, pregunta.Id)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 10.0 / float64(len(params.Temas))))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 			ids_preguntas_seleccionadas = append(ids_preguntas_seleccionadas, pregunta.Id)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad51bb5c000045007639" { //Avanzado
		// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0 / float64(len(params.Temas))))

		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 			ids_preguntas_seleccionadas = append(ids_preguntas_seleccionadas, pregunta.Id)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 30.0 / float64(len(params.Temas))))

		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 			ids_preguntas_seleccionadas = append(ids_preguntas_seleccionadas, pregunta.Id)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 13.0 / float64(len(params.Temas))))

		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimit(ids_temas_a_buscar, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 			ids_preguntas_seleccionadas = append(ids_preguntas_seleccionadas, pregunta.Id)
		// }
		// }
		// }

		//#####################
		//#####################
		//#####################
		//#####################

		// if params.NumeroPreguntas > len(ids_preguntas_seleccionadas) {

		// 	var numero_preguntas_reajuste = params.NumeroPreguntas - len(ids_preguntas_seleccionadas)
		// if usuario.IdNivel.Hex() == "5ea6ad3dbb5c000045007637" { //BASICO
		// 	var n_preguntas_basicas = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 80.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 15.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 5.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad4abb5c000045007638" { //Intermedio
		// 	var n_preguntas_basicas = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 65.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 25.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 10.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad51bb5c000045007639" { //Avanzado
		// 	var n_preguntas_basicas = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 57.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 30.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(numero_preguntas_reajuste) / 100.0 * 13.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsTemasIdNivelLimitNotInPreguntas(ids_temas, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas, ids_preguntas_seleccionadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// }

		//#####################
		//#####################
		//#####################

		rand.Seed(time.Now().UnixNano())
		for i := len(preguntas_mezcladas) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			preguntas_mezcladas[i], preguntas_mezcladas[j] = preguntas_mezcladas[j], preguntas_mezcladas[i]
		}

		if len(preguntas_mezcladas) > 0 {
			if len(preguntas_mezcladas) > params.NumeroPreguntas {
				historial.PreguntasTotales = params.NumeroPreguntas
				historial.PreguntasBlancas = params.NumeroPreguntas
			} else {
				historial.PreguntasTotales = len(preguntas_mezcladas)
				historial.PreguntasBlancas = len(preguntas_mezcladas)
			}
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				if cont_preguntas >= params.NumeroPreguntas {
				} else {
					var historial_pregunta models.HistorialesPreguntas
					historial_pregunta.Id = bson.NewObjectId()
					historial_pregunta.IdHistorial = historial.Id
					historial_pregunta.IdPregunta = pregunta.Id

					if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
						return
					}
					cont_preguntas = cont_preguntas + 1
				}
			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}

func CreateTestLegislaciones(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.CreateTestLegislacion
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_temas []string
		var ids_oficiales []string
		var ids_simulacros []string
		var ids_basicosconfundidos []string
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = params.Legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
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
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		var ids_legislaciones []bson.ObjectId

		for _, legislacion := range params.Legislaciones {
			ids_legislaciones = append(ids_legislaciones, bson.ObjectIdHex(legislacion))
		}

		var preguntas_mezcladas []models.Preguntas

		// if usuario.IdNivel.Hex() == "5ea6ad3dbb5c000045007637" { //BASICO
		// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 80.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 15.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 5.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad4abb5c000045007638" { //Intermedio
		// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 65.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 25.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 10.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad51bb5c000045007639" { //Avanzado

		if len(params.Legislaciones) == 1 {
			preguntas_bbdd, err := daoPreguntas.FindRandomPreguntasByLegislacion(params.NumeroPreguntas, params.Legislaciones[0])
			if err != nil {
			} else {

				for _, pregunta := range preguntas_bbdd {
					preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
				}
			}
		} else {

			num_preguntas := params.NumeroPreguntas/len(params.Legislaciones) + 5

			for _, legislacion := range params.Legislaciones {
				preguntas_bbdd, err := daoPreguntas.FindRandomPreguntasByLegislacion(num_preguntas, legislacion)
				// log.Println(len(preguntas_bbdd))
				if err != nil {
				} else {
					for _, pregunta := range preguntas_bbdd {
						preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
					}
				}
			}

			// var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
			// if err != nil {
			// } else {
			// 	for _, pregunta := range preguntas_basicas {
			// 		preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 	}
			// }

			// // var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 30.0))
			// var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
			// if err != nil {
			// } else {
			// 	for _, pregunta := range preguntas_intermedias {
			// 		preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 	}
			// }

			// var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// // var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 13.0))
			// preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
			// if err != nil {
			// } else {
			// 	for _, pregunta := range preguntas_avanzadas {
			// 		preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 	}
			// }

			// if len(preguntas_mezcladas) < params.NumeroPreguntas {
			// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
			// 	if err != nil {
			// 	} else {
			// 		for _, pregunta := range preguntas_basicas {
			// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 		}
			// 	}

			// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// 	// var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 30.0))
			// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
			// 	if err != nil {
			// 	} else {
			// 		for _, pregunta := range preguntas_intermedias {
			// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 		}
			// 	}

			// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// 	// var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 13.0))
			// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
			// 	if err != nil {
			// 	} else {
			// 		for _, pregunta := range preguntas_avanzadas {
			// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 		}
			// 	}
			// }

			// if len(preguntas_mezcladas) < params.NumeroPreguntas {
			// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
			// 	if err != nil {
			// 	} else {
			// 		for _, pregunta := range preguntas_basicas {
			// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 		}
			// 	}

			// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// 	// var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 30.0))
			// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
			// 	if err != nil {
			// 	} else {
			// 		for _, pregunta := range preguntas_intermedias {
			// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 		}
			// 	}

			// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
			// 	// var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 13.0))
			// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsLegislacionesIdNivelLimit(ids_legislaciones, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
			// 	if err != nil {
			// 	} else {
			// 		for _, pregunta := range preguntas_avanzadas {
			// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			// 		}
			// 	}
			// }
		}

		rand.Seed(time.Now().UnixNano())
		for i := len(preguntas_mezcladas) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			preguntas_mezcladas[i], preguntas_mezcladas[j] = preguntas_mezcladas[j], preguntas_mezcladas[i]
		}

		if len(preguntas_mezcladas) > 0 {

			if len(preguntas_mezcladas) > params.NumeroPreguntas {
				historial.PreguntasTotales = params.NumeroPreguntas
				historial.PreguntasBlancas = params.NumeroPreguntas
			} else {
				historial.PreguntasTotales = len(preguntas_mezcladas)
				historial.PreguntasBlancas = len(preguntas_mezcladas)
			}
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				if cont_preguntas >= params.NumeroPreguntas {
				} else {
					var historial_pregunta models.HistorialesPreguntas
					historial_pregunta.Id = bson.NewObjectId()
					historial_pregunta.IdHistorial = historial.Id
					historial_pregunta.IdPregunta = pregunta.Id

					if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
						return
					}
					cont_preguntas = cont_preguntas + 1
				}
			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}
func CreateTestBasicosConfundidos(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.CreateTestBasicoConfundido
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_temas []string
		var ids_legislaciones []string
		var ids_oficiales []string
		var ids_simulacros []string
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = params.BasicosConfundidos
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
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		var ids_basicosconfundidos []bson.ObjectId

		for _, basicoconfundido := range params.BasicosConfundidos {
			ids_basicosconfundidos = append(ids_basicosconfundidos, bson.ObjectIdHex(basicoconfundido))
		}

		var preguntas_mezcladas []models.Preguntas

		// if usuario.IdNivel.Hex() == "5ea6ad3dbb5c000045007637" { //BASICO
		// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 80.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 15.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 5.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad4abb5c000045007638" { //Intermedio
		// 	var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 65.0))
		// 	preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_basicas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 25.0))
		// 	preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_intermedias {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// 	var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 10.0))
		// 	preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		// 	if err != nil {
		// 	} else {
		// 		for _, pregunta := range preguntas_avanzadas {
		// 			preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
		// 		}
		// 	}

		// } else if usuario.IdNivel.Hex() == "5ea6ad51bb5c000045007639" { //Avanzado
		var n_preguntas_basicas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 57.0))
		preguntas_basicas, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad3dbb5c000045007637"), n_preguntas_basicas)
		if err != nil {
		} else {
			for _, pregunta := range preguntas_basicas {
				preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			}
		}

		var n_preguntas_intermedias = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 30.0))
		preguntas_intermedias, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad4abb5c000045007638"), n_preguntas_intermedias)
		if err != nil {
		} else {
			for _, pregunta := range preguntas_intermedias {
				preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			}
		}

		var n_preguntas_avanzadas = int(math.Ceil(float64(params.NumeroPreguntas) / 100.0 * 13.0))
		preguntas_avanzadas, err := daoPreguntas.FindPreguntaByIdsBasicosConfundidosIdNivelLimit(ids_basicosconfundidos, bson.ObjectIdHex("5ea6ad51bb5c000045007639"), n_preguntas_avanzadas)
		if err != nil {
		} else {
			for _, pregunta := range preguntas_avanzadas {
				preguntas_mezcladas = append(preguntas_mezcladas, pregunta)
			}
			// }

		}

		rand.Seed(time.Now().UnixNano())
		for i := len(preguntas_mezcladas) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			preguntas_mezcladas[i], preguntas_mezcladas[j] = preguntas_mezcladas[j], preguntas_mezcladas[i]
		}

		if len(preguntas_mezcladas) > 0 {

			if len(preguntas_mezcladas) > params.NumeroPreguntas {
				historial.PreguntasTotales = params.NumeroPreguntas
				historial.PreguntasBlancas = params.NumeroPreguntas
			} else {
				historial.PreguntasTotales = len(preguntas_mezcladas)
				historial.PreguntasBlancas = len(preguntas_mezcladas)
			}
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error_code": "2"})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				if cont_preguntas >= params.NumeroPreguntas {
				} else {
					var historial_pregunta models.HistorialesPreguntas
					historial_pregunta.Id = bson.NewObjectId()
					historial_pregunta.IdHistorial = historial.Id
					historial_pregunta.IdPregunta = pregunta.Id

					if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error_code": "1"})
						return
					}
					cont_preguntas = cont_preguntas + 1
				}
			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}
func CreateTestOficiales(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.CreateTestOficiales
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_temas []string
		var ids_basicosconfundidos []string
		var ids_simulacros []string
		var ids_legislaciones []string
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = params.Oficiales
		historial.Simulacros = ids_simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
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
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		var ids_oficiales []bson.ObjectId

		for _, simulacrooficial := range params.Oficiales {
			ids_oficiales = append(ids_oficiales, bson.ObjectIdHex(simulacrooficial))
		}

		var preguntas_mezcladas []models.Preguntas

		preguntas_oficiales, err := daoExamenesOficialesPreguntas.FindPreguntaByIdExamenOficialOrder(ids_oficiales)
		if err != nil {
		} else {
			for _, pregunta := range preguntas_oficiales {
				pregunta_real, err := daoPreguntas.FindPreguntaById(pregunta.IdPregunta.Hex())
				if err != nil {
				} else {
					preguntas_mezcladas = append(preguntas_mezcladas, pregunta_real)
				}
			}
		}

		if len(preguntas_mezcladas) > 0 {

			if len(preguntas_mezcladas) > params.NumeroPreguntas {
				historial.PreguntasTotales = params.NumeroPreguntas
				historial.PreguntasBlancas = params.NumeroPreguntas
			} else {
				historial.PreguntasTotales = len(preguntas_mezcladas)
				historial.PreguntasBlancas = len(preguntas_mezcladas)
			}
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				if cont_preguntas >= params.NumeroPreguntas {
				} else {
					var historial_pregunta models.HistorialesPreguntas
					historial_pregunta.Id = bson.NewObjectId()
					historial_pregunta.IdHistorial = historial.Id
					historial_pregunta.IdPregunta = pregunta.Id

					if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
						return
					}
					cont_preguntas = cont_preguntas + 1
				}
			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}
func CreateTestSimulacros(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.CreateTestSimulacros
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}
	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params.Email)))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		var ids_temas []string
		var ids_basicosconfundidos []string
		var ids_oficiales []string
		var ids_legislaciones []string
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = ids_temas
		historial.Legislaciones = ids_legislaciones
		historial.Oficiales = ids_oficiales
		historial.Simulacros = params.Simulacros
		historial.BasicosConfundidos = ids_basicosconfundidos
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
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		var ids_simulacros []bson.ObjectId

		for _, simulacrooficial := range params.Simulacros {
			ids_simulacros = append(ids_simulacros, bson.ObjectIdHex(simulacrooficial))
		}

		var preguntas_mezcladas []models.Preguntas

		preguntas_simulacros, err := daoSimulacrosPreguntas.FindPreguntaByIdSimulacroOrder(ids_simulacros)
		if err != nil {
		} else {
			for _, pregunta := range preguntas_simulacros {
				pregunta_real, err := daoPreguntas.FindPreguntaById(pregunta.IdPregunta.Hex())
				if err != nil {
				} else {
					preguntas_mezcladas = append(preguntas_mezcladas, pregunta_real)
				}
			}
		}

		if len(preguntas_mezcladas) > 0 {

			if len(preguntas_mezcladas) > params.NumeroPreguntas {
				historial.PreguntasTotales = params.NumeroPreguntas
				historial.PreguntasBlancas = params.NumeroPreguntas
			} else {
				historial.PreguntasTotales = len(preguntas_mezcladas)
				historial.PreguntasBlancas = len(preguntas_mezcladas)
			}
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				if cont_preguntas >= params.NumeroPreguntas {
				} else {
					var historial_pregunta models.HistorialesPreguntas
					historial_pregunta.Id = bson.NewObjectId()
					historial_pregunta.IdHistorial = historial.Id
					historial_pregunta.IdPregunta = pregunta.Id

					if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
						helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
						return
					}
					cont_preguntas = cont_preguntas + 1
				}
			}
		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Preguntas no encontradas"})
			return
		}
		helper.ResponseWithJson(w, http.StatusOK, historial)
	}
}
func GetTest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var historialADevolver models.HistorialesADevolver
	var preguntasADevolver []models.PreguntasADevolver
	var respuestas_marcadas = make(map[string]string)
	var ids_preguntas []bson.ObjectId

	historial, err := daoHistoriales.FindHistorialById(id)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}

	historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}

	for _, historial_pregunta := range historiales_preguntas {
		if historial_pregunta.IdRespuesta.Hex() != "" {
			respuestas_marcadas[historial_pregunta.IdPregunta.Hex()] = historial_pregunta.IdRespuesta.Hex()
		}
		ids_preguntas = append(ids_preguntas, historial_pregunta.IdPregunta)
	}
	preguntas_del_historial, err := daoPreguntas.FindPreguntasByIds(ids_preguntas, nil)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}
	respuestas_del_historial, err := daoRespuestas.FindRespuestasByIdsPregunta(ids_preguntas)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}

	for _, id_pregunta := range ids_preguntas {
		var la_pregunta models.Preguntas
		for _, pregunta_del_historial := range preguntas_del_historial {
			if pregunta_del_historial.Id == id_pregunta {
				la_pregunta = pregunta_del_historial
				if len(pregunta_del_historial.SchemaId) > 0 {
					files, err := daoPreguntas.FindFileByPreguntaId(id_pregunta.Hex())
					if err != nil {
						log.Println(err.Error())
					}
					la_pregunta.Schema = files
					// log.Println(la_pregunta.Schema)
					}
			}
		}


		var respuestasADevolver []models.RespuestasADevolver

		for _, respuesta := range respuestas_del_historial {
			if respuesta.IdPregunta == id_pregunta {
				var respuestaADevolver models.RespuestasADevolver
				respuestaADevolver.Id = respuesta.Id
				respuestaADevolver.Respuesta = respuesta.Respuesta
				if helper.CheckUser(params["email"]) {
					respuestaADevolver.Correcta = respuesta.Correcta
				} else {
					respuestaADevolver.Correcta = false
				}
				respuestasADevolver = append(respuestasADevolver, respuestaADevolver)
			}
		}

		
		rand.Seed(time.Now().UnixNano())
		for i := len(respuestasADevolver) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			respuestasADevolver[i], respuestasADevolver[j] = respuestasADevolver[j], respuestasADevolver[i]
		}
		
		var preguntaADevolver models.PreguntasADevolver
		preguntaADevolver.Id = la_pregunta.Id
		preguntaADevolver.Schema = la_pregunta.Schema
		preguntaADevolver.Pregunta = la_pregunta.Pregunta

		// if len(preguntaADevolver.Schema) > 0 {
		// 	log.Println(preguntaADevolver.Id)
		// }

		if helper.CheckUser(params["email"]) {
			preguntaADevolver.Explicacion = la_pregunta.Explicacion
		} else {
			preguntaADevolver.Explicacion = ""
		}
		preguntaADevolver.Respuestas = respuestasADevolver
		preguntaADevolver.IdTema = la_pregunta.IdTema
		preguntaADevolver.Oficial = la_pregunta.Oficial
		preguntaADevolver.AnioOficial = la_pregunta.AnioOficial
		preguntasADevolver = append(preguntasADevolver, preguntaADevolver)

	}

	historialADevolver.Historial = historial
	historialADevolver.Preguntas = preguntasADevolver
	historialADevolver.RespuestasMarcadas = respuestas_marcadas

	helper.ResponseWithJson(w, http.StatusOK, historialADevolver)
}
func CambiarTiempoTranscurrido(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	historial, err := daoHistoriales.FindHistorialById(pregunta_a_corregir["id_test"])
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

	if err := daoHistoriales.UpdateHistorial(historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"status": "true"})
}
func CorregirPregunta(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	historial, err := daoHistoriales.FindHistorialById(pregunta_a_corregir["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
		return
	}

	pregunta, err := daoPreguntas.FindPreguntaById(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	respuesta, err := daoRespuestas.FindRespuestaById(pregunta_a_corregir["id_respuesta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	respuesta_correcta, err := daoRespuestas.FindRespuestaCorrectaByIdPregunta(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	historial_pregunta, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorialIdPregunta(historial.Id.Hex(), pregunta_a_corregir["id_pregunta"])
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

	if err := daoHistoriales.UpdateHistorial(historial); err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	historial_pregunta.IdRespuesta = respuesta.Id

	if respuesta.Id == respuesta_correcta.Id {
		historial_pregunta.Correcta = true

		historial.PreguntasAcertadas = historial.PreguntasAcertadas + 1
		historial.PreguntasBlancas = historial.PreguntasBlancas - 1
		historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)
		if err := daoHistoriales.UpdateHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

	} else {
		historial_pregunta.Correcta = false

		historial.PreguntasFalladas = historial.PreguntasFalladas + 1
		historial.PreguntasBlancas = historial.PreguntasBlancas - 1
		historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)
		if err := daoHistoriales.UpdateHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}
	}

	if err := daoHistorialesPreguntas.UpdateHistorialPregunta(historial_pregunta); err != nil {
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
func CorregirPreguntaYaContestada(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	historial, err := daoHistoriales.FindHistorialById(pregunta_a_corregir["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
		return
	}

	pregunta, err := daoPreguntas.FindPreguntaById(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	respuesta, err := daoRespuestas.FindRespuestaById(pregunta_a_corregir["id_respuesta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	respuesta_correcta, err := daoRespuestas.FindRespuestaCorrectaByIdPregunta(pregunta_a_corregir["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
		return
	}

	historial_pregunta, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorialIdPregunta(historial.Id.Hex(), pregunta_a_corregir["id_pregunta"])
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
func CorregirTest(w http.ResponseWriter, r *http.Request) {
	var respuestas_marcadas map[string]string
	var respuestasCorrectas = make(map[string]models.Respuestas)
	var preguntas = make(map[string]models.Preguntas)
	var correlacion_temas_areas = make(map[string]string)
	var falladas_por_area = make(map[string]int)
	var blancas_por_area = make(map[string]int)
	var temas_a_enviar = make(map[string]string)
	temas_a_enviar["5ea6ae74bb5c00004500763a"] = "Conducta Humana"
	temas_a_enviar["5ea6ae86bb5c00004500763b"] = "Derecho Administrativo"
	temas_a_enviar["5ea6ae8dbb5c00004500763c"] = "Derecho Penal"
	temas_a_enviar["5ea6ae97bb5c00004500763d"] = "Derecho Penitenciario"
	falladas_por_area["5ea6ae74bb5c00004500763a"] = 0
	falladas_por_area["5ea6ae86bb5c00004500763b"] = 0
	falladas_por_area["5ea6ae8dbb5c00004500763c"] = 0
	falladas_por_area["5ea6ae97bb5c00004500763d"] = 0
	blancas_por_area["5ea6ae74bb5c00004500763a"] = 0
	blancas_por_area["5ea6ae86bb5c00004500763b"] = 0
	blancas_por_area["5ea6ae8dbb5c00004500763c"] = 0
	blancas_por_area["5ea6ae97bb5c00004500763d"] = 0
	temas, err := daoTemas.FindAllTemas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}
	for _, tema := range temas {
		correlacion_temas_areas[tema.Id.Hex()] = tema.IdArea.Hex()
	}

	var params models.CorregirTest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	if err := json.Unmarshal([]byte(params.RespuestasMarcadas), &respuestas_marcadas); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
	}

	historial, err := daoHistoriales.FindHistorialById(params.IdTest)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioById(historial.IdUsuario.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "usuario"})
		return
	}

	historial.PreguntasAcertadas = 0
	historial.PreguntasFalladas = 0
	historial.PreguntasBlancas = historial.PreguntasTotales

	historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}

	for _, historial_pregunta := range historiales_preguntas {

		pregunta, err := daoPreguntas.FindPreguntaById(historial_pregunta.IdPregunta.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta"})
			return
		}
		preguntas[pregunta.Id.Hex()] = pregunta

		respuesta_correcta, err := daoRespuestas.FindRespuestaCorrectaByIdPregunta(pregunta.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Respuestas no encontradas"})
			return
		}
		respuestasCorrectas[pregunta.Id.Hex()] = respuesta_correcta

		if respuesta_a_la_pregunta, ok := respuestas_marcadas[pregunta.Id.Hex()]; ok {

			if respuesta_a_la_pregunta != "" {

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
					falladas_por_area[correlacion_temas_areas[pregunta.IdTema.Hex()]]++
				}

				if err := daoHistorialesPreguntas.UpdateHistorialPregunta(historial_pregunta); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
					return
				}
			} else {
				blancas_por_area[correlacion_temas_areas[pregunta.IdTema.Hex()]]++
			}
		} else {
			blancas_por_area[correlacion_temas_areas[pregunta.IdTema.Hex()]]++
		}

	}

	puntaje := float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) / 3)

	if puntaje < 0 {
		historial.Puntuacion = 0
	} else {
		historial.Puntuacion = puntaje
	}
	// historial.Puntuacion = float32(historial.PreguntasAcertadas) - (float32(historial.PreguntasFalladas) * 0.33)
	historial.Terminado = true
	historial.TiempoTranscurrido = params.Tiempo

	if err := daoHistoriales.UpdateHistorial(historial); err != nil {
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

	var respuesta_corregir = map[string]string{"respuestas_correctas": string(json_respuestasCorrectas), "preguntas": string(json_preguntas)}

	json_falladas_por_area, err := json.Marshal(falladas_por_area)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	json_blancas_por_area, err := json.Marshal(blancas_por_area)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	json_temas_a_enviar, err := json.Marshal(temas_a_enviar)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	respuesta_corregir["temas_a_enviar"] = string(json_temas_a_enviar)
	respuesta_corregir["falladas_por_area"] = string(json_falladas_por_area)
	respuesta_corregir["blancas_por_area"] = string(json_blancas_por_area)
	respuesta_corregir["nivel"] = ""
	if strings.Compare(historial.Tipo, "Test de nivel") == 0 {

		if historial.PreguntasAcertadas >= 40 {
			usuario.IdNivel = bson.ObjectIdHex("5ea6ad51bb5c000045007639")
			respuesta_corregir["nivel"] = "Avanzado"
		} else if historial.PreguntasAcertadas >= 28 {
			usuario.IdNivel = bson.ObjectIdHex("5ea6ad4abb5c000045007638")
			respuesta_corregir["nivel"] = "Intermedio"
		} else {
			usuario.IdNivel = bson.ObjectIdHex("5ea6ad3dbb5c000045007637")
			respuesta_corregir["nivel"] = "Básico"
		}
		usuario.FirstLogin = true
		if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
			return
		}

	}

	helper.ResponseWithJson(w, http.StatusOK, respuesta_corregir)
}
func CorregirTestYaContestado(w http.ResponseWriter, r *http.Request) {

	var respuestas_marcadas map[string]string
	var respuestasCorrectas = make(map[string]models.Respuestas)
	var preguntas = make(map[string]models.Preguntas)
	var correlacion_temas_areas = make(map[string]string)
	var falladas_por_area = make(map[string]int)
	var blancas_por_area = make(map[string]int)
	var temas_a_enviar = make(map[string]string)
	temas_a_enviar["5ea6ae74bb5c00004500763a"] = "Conducta Humana"
	temas_a_enviar["5ea6ae86bb5c00004500763b"] = "Derecho Administrativo"
	temas_a_enviar["5ea6ae8dbb5c00004500763c"] = "Derecho Penal"
	temas_a_enviar["5ea6ae97bb5c00004500763d"] = "Derecho Penitenciario"
	falladas_por_area["5ea6ae74bb5c00004500763a"] = 0
	falladas_por_area["5ea6ae86bb5c00004500763b"] = 0
	falladas_por_area["5ea6ae8dbb5c00004500763c"] = 0
	falladas_por_area["5ea6ae97bb5c00004500763d"] = 0
	blancas_por_area["5ea6ae74bb5c00004500763a"] = 0
	blancas_por_area["5ea6ae86bb5c00004500763b"] = 0
	blancas_por_area["5ea6ae8dbb5c00004500763c"] = 0
	blancas_por_area["5ea6ae97bb5c00004500763d"] = 0
	temas, err := daoTemas.FindAllTemas()
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}
	for _, tema := range temas {
		correlacion_temas_areas[tema.Id.Hex()] = tema.IdArea.Hex()
	}

	var params models.CorregirTest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	if err := json.Unmarshal([]byte(params.RespuestasMarcadas), &respuestas_marcadas); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
	}

	historial, err := daoHistoriales.FindHistorialById(params.IdTest)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
		return
	}

	historial.PreguntasAcertadas = 0
	historial.PreguntasFalladas = 0
	historial.PreguntasBlancas = historial.PreguntasTotales

	historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
		return
	}
	for _, historial_pregunta := range historiales_preguntas {

		pregunta, err := daoPreguntas.FindPreguntaById(historial_pregunta.IdPregunta.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "pregunta"})
			return
		}
		preguntas[pregunta.Id.Hex()] = pregunta

		respuesta_correcta, err := daoRespuestas.FindRespuestaCorrectaByIdPregunta(pregunta.Id.Hex())
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
				falladas_por_area[correlacion_temas_areas[pregunta.IdTema.Hex()]]++
			}

		} else {
			blancas_por_area[correlacion_temas_areas[pregunta.IdTema.Hex()]]++
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

	json_falladas_por_area, err := json.Marshal(falladas_por_area)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	json_blancas_por_area, err := json.Marshal(blancas_por_area)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}
	json_temas_a_enviar, err := json.Marshal(temas_a_enviar)
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{
		"respuestas_correctas": string(json_respuestasCorrectas),
		"preguntas":            string(json_preguntas),
		"temas_a_enviar":       string(json_temas_a_enviar),
		"falladas_por_area":    string(json_falladas_por_area),
		"blancas_por_area":     string(json_blancas_por_area),
	})
}
func HayError(w http.ResponseWriter, r *http.Request) {

	var datos_error = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&datos_error); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	pregunta, err := daoPreguntas.FindPreguntaById(datos_error["id_pregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(datos_error["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario no encontrado"})
		return
	} else {

		historialRepot, err := daoHistoriales.FindHistorialById(datos_error["id_test"])

		if err != nil{
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Historial no encontrado"})
			return	
		}

		tema, err := daoTemas.FindTemaById(pregunta.IdTema.Hex())

		if err != nil{
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al buscar el tema de la pregunta"})
			return	
		}

		if !helper.SendEmail("PENITENCIARIOS.COM :: Hay error", "El usuario '"+usuario.Name+"' <br /> Email: '"+usuario.Email+"' <br /> ID '"+usuario.Id.Hex()+"' <br /> " +"Id Test: "+ datos_error["id_test"] +" <br />"+ "Tipo: "+ historialRepot.Tipo +"<br /> "+ "Tema: "+ tema.Name + "<br /><br /> Ha encontrado un error en la pregunta '"+pregunta.Id.Hex()+"'<br /><br /><b>DETALLES:</b>"+datos_error["detalles"]+" ", "penitenciarios@penitenciarios.com", usuario.Email) {
			// if !helper.SendEmail("PENITENCIARIOS.COM :: Hay error", "El usuario '"+usuario.Name+"' <br /> Email: '"+usuario.Email+"' <br /> ID '"+usuario.Id.Hex()+"' <br /> " +"Id Test: "+ datos_error["id_test"] +" <br />"+ "Tipo: "+ historialRepot.Tipo +"<br /> "+ "Tema: "+ tema.Name + "<br /><br /> Ha encontrado un error en la pregunta '"+pregunta.Id.Hex()+"'<br /><br /><b>DETALLES:</b>"+datos_error["detalles"]+" ", "toledof764@gmail.com", usuario.Email) {
			// if !helper.SendEmail("PENITENCIARIOS.COM :: Hay error", "El usuario '"+usuario.Name+"' <br /> Email: '"+usuario.Email+"' <br /> ID '"+usuario.Id.Hex()+"' <br /> " +"Id Test: "+ datos_error["id_test"] +" <br />"+ "Tipo: "+ historialRepot.Tipo +"<br /><br />" +" Ha encontrado un error en la pregunta '"+pregunta.Id.Hex()+"'<br /><br /><b>DETALLES:</b>"+datos_error["detalles"]+" ", "penitenciarios@penitenciarios.com", usuario.Email) {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Se ha producido un error al enviar el email"})
			return
		}
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "true"})
}
func TieneDuda(w http.ResponseWriter, r *http.Request) {

	var datos_error = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&datos_error); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	pregunta, err := daoPreguntas.FindPreguntaById(datos_error["id_pregunta"])
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
func GuardarComoSimulacro(w http.ResponseWriter, r *http.Request) {

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

		if helper.CheckUser(usuario.Email) {
			/* CREAMOS EL OBJETO SIMULACRO */
			var simulacroAGuardar models.Simulacros
			simulacroAGuardar.Id = bson.NewObjectId()
			simulacroAGuardar.Name = params["nombre_simulacro"]
			simulacroAGuardar.IdNivel = params["id_nivel"]
			// simulacroAGuardar.IdNivel = usuario.IdNivel
			if err := daoSimulacros.InsertSimulacro(simulacroAGuardar); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'Simulacros'"})
				return
			}

			/* BUSCAMOS EL HISTORIAL */
			historial, err := daoHistoriales.FindHistorialById(params["id_test"])
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial"})
				return
			}
			/* BUSCAMOS LAS PREGUNTAS DEL HISTORIAL */
			historiales_preguntas, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial.Id.Hex())
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas"})
				return
			}

			var cont = 0

			for _, historial_pregunta := range historiales_preguntas {

				//historial_pregunta.IdPregunta

				var simulacroPreguntaAGuardar models.SimulacrosPreguntas
				simulacroPreguntaAGuardar.Id = bson.NewObjectId()
				simulacroPreguntaAGuardar.IdSimulacro = simulacroAGuardar.Id
				simulacroPreguntaAGuardar.IdPregunta = historial_pregunta.IdPregunta
				simulacroPreguntaAGuardar.Orden = cont
				if err := daoSimulacrosPreguntas.InsertSimulacroPregunta(simulacroPreguntaAGuardar); err != nil {
					helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al guardar el registro 'SimulacrosPreguntas'"})
					return
				}
				cont = cont + 1

			}

			helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})

		} else {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "No tienes permisos"})
			return
		}
	}

}
func BuscarPreguntaParaEditar(w http.ResponseWriter, r *http.Request) {

	var pregunta_a_corregir = map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&pregunta_a_corregir); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	pregunta, err := daoPreguntasWRespuestas.FindPreguntaByIdWRespuestas(pregunta_a_corregir["formIdPregunta"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Pregunta no encontrada"})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, pregunta)

}
func RepetirTest(w http.ResponseWriter, r *http.Request) {
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
		historial_a_copiar, err := daoHistoriales.FindHistorialById(params["id"])
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historial_a_copiar"})
			return
		}
		if historial_a_copiar.IdUsuario != usuario.Id {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error historial no válido", "error2": "historial_a_copiar"})
			return
		}
		var historial models.Historiales
		historial.Id = bson.NewObjectId()
		historial.Tipo = historial_a_copiar.Tipo
		historial.Fecha = helper.MakeTimestamp()
		historial.Temas = historial_a_copiar.Temas
		historial.Legislaciones = historial_a_copiar.Legislaciones
		historial.BasicosConfundidos = historial_a_copiar.BasicosConfundidos
		historial.Oficiales = historial_a_copiar.Oficiales
		historial.Simulacros = historial_a_copiar.Simulacros
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
		if err := daoHistoriales.InsertHistorial(historial); err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Error al insertar el historial"})
			return
		}

		historiales_preguntas_a_copiar, err := daoHistorialesPreguntas.FindHistorialPreguntaByIdHistorial(historial_a_copiar.Id.Hex())
		if err != nil {
			helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "historiales_preguntas_a_copiar"})
			return
		}

		var preguntas_mezcladas []models.Preguntas

		for _, historial_pregunta_a_copiar := range historiales_preguntas_a_copiar {
			pregunta_a_copiar, err := daoPreguntas.FindPreguntaById(historial_pregunta_a_copiar.IdPregunta.Hex())
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
			if err := daoHistoriales.UpdateHistorial(historial); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error(), "error2": "UpdateHistorial"})
				return
			}

			var cont_preguntas = 0
			for _, pregunta := range preguntas_mezcladas {
				var historial_pregunta models.HistorialesPreguntas
				historial_pregunta.Id = bson.NewObjectId()
				historial_pregunta.IdHistorial = historial.Id
				historial_pregunta.IdPregunta = pregunta.Id

				if err := daoHistorialesPreguntas.InsertHistorialPregunta(historial_pregunta); err != nil {
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

func UpdatePuntajetest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
	}

	hisorialResut, err := daoHistoriales.FindHistorialById(params["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurrio un error al buscar el historial"})
		return
	}

	numero, err := strconv.ParseFloat(params["newPuntaje"], 32)
	if err != nil {
		fmt.Println("Error al convertir el número:", err)
		return
	}
	hisorialResut.Puntuacion = float32(numero)

	if err := daoHistoriales.UpdateHistorial(hisorialResut); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurio un error al actualizar el puntaje"})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "ok"})
}

func UpdatePuntajetestCP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
	}

	hisorialResut, err := daoHistorialesCP.FindHistorialById(params["id_test"])
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurrio un error al buscar el historial"})
		return
	}

	numero, err := strconv.ParseFloat(params["newPuntaje"], 32)
	if err != nil {
		fmt.Println("Error al convertir el número:", err)
		return
	}
	hisorialResut.Puntuacion = float32(numero)

	if err := daoHistorialesCP.UpdateHistorial(hisorialResut); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Ocurio un error al actualizar el puntaje"})
		return
	}
	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "ok"})
}

func GetAllPreguntas (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	// var preguntas []models.Preguntas

	preguntas, err := daoPreguntas.FindAllPreguntas()
	
	if err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]interface{}{"result": preguntas})
}

func DeletePreguntaAndRespuesta (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	vars := mux.Vars(r)

	pregunta_id := vars["id"]

	if err := daoPreguntas.RemovePregunta(pregunta_id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	if err := daoRespuestas.RemoveRespuestasByIdPregunta(pregunta_id); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": err.Error()})
		return
	}

	helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func DeleteAllPreguntasLegislacion(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	vars := mux.Vars(r)

	leg_id := vars["id"]

	preguntas, err := daoPreguntas.FindAllPreguntasByLegislacion(leg_id)

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