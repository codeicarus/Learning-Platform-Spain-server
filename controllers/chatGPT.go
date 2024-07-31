package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"test/helper"
	"test/models"
)

var (
	daoQuestion = models.QuetionGPT{}
)

func GetInfoGPT(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var question models.QuetionGPT

	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	pregunta, err := daoPreguntas.FindPreguntaById(question.PreguntaID)

	if err != nil {
		log.Println(err)
		helper.ResponseWithJson(w, http.StatusOK, "Ocurrio un error al obtener la pregunta")
		return
	}

	opciones, err := daoRespuestas.FindRespuestaByIdPregunta(question.PreguntaID)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusOK, "Ocurrio un error al obtener las opciones de la pregunta")
		return
	}

	correcto, err := daoRespuestas.FindRespuestaCorrectaByIdPregunta(question.PreguntaID)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusOK, "Ocurrio un error al obtener la respuesta correcta de la pregunta")
		return
	}

	mirespuesta, err := daoRespuestas.FindRespuestaById(question.RespuestaID)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusOK, "Ocurrio un error al obtener respuesta correcta de la pregunta")
		return
	}

	opt := "Estas son las opciones " + opciones[0].Respuesta + "\n" + opciones[1].Respuesta + "\n" + opciones[2].Respuesta + "\n" + opciones[3].Respuesta

	var response string = "Segun el BOE responde si esta bien o no, si esta mal justifica, de lo contrario dame una explicacion simple pero clara, quiero que las respuestas sean cortas pero entendibles. \n La pregunta: " + pregunta.Pregunta + "\n" + opt + "\n" + "Mi respuesta es: " + mirespuesta.Respuesta + "\n" + "la respuesta correcta es : " + correcto.Respuesta

	question.Contexto = append(question.Contexto, response)

	formData := url.Values{}
	formData.Set("prompt", response)

	req, err := http.NewRequest("POST", "http://localhost:5000/initial_chat", strings.NewReader(formData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusOK, "Ocurrio un error al obtener respuesta de GPT")
		return
	}

	defer res.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println("Error al leer la respuesta:", err)
		return
	}

	decodedBody, err := unquoteUnicode(string(body))

	if err != nil {
		log.Println("Error al decodificar la respuesta:", err)
		return
	}

	re := regexp.MustCompile(`\\"([^\\"]+)\\"`)
	decodedBody = re.ReplaceAllString(decodedBody, `<b>$1</b>`)

	responseMap := make(map[string]interface{})
	responseMap["result"] = "success"
	responseMap["response"] = decodedBody
	responseMap["Contexto"] = question.Contexto

	helper.ResponseWithJson(w, http.StatusOK, responseMap)

}

func unquoteUnicode(s string) (string, error) {
	result := ""
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == 'u' {
			code, err := strconv.ParseInt(s[i+2:i+6], 16, 32)
			if err != nil {
				return "", err
			}
			result += string(code)
			i += 5
		} else {
			result += string(s[i])
		}
	}
	return result, nil
}

func ContinueChatGPT(w http.ResponseWriter, r *http.Request) {
	var continueContext models.ContinueChat

	if err := json.NewDecoder(r.Body).Decode(&continueContext); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	// log.Println(continueContext.Context)

	ultimoElemento := continueContext.Context[len(continueContext.Context)-1]

	// Seleccionar todos los elementos excepto el último
	elementosAnteriores := continueContext.Context[:len(continueContext.Context)-1]

	data := struct {
		Responde_A      string   `json:"responde_a"`
		Ten_de_contexto []string `json:"Toma_contexto"`
	}{
		Responde_A:      ultimoElemento,
		Ten_de_contexto: elementosAnteriores,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	formData := url.Values{}
	formData.Set("prompt", string(jsonData))

	// log.Println(formData)

	req, err := http.NewRequest("POST", "http://localhost:5000/chat", strings.NewReader(formData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		helper.ResponseWithJson(w, http.StatusOK, "Ocurrio un error al obtener respuesta de GPT")
		return
	}

	defer res.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Println("Error al leer la respuesta:", err)
		return
	}

	decodedBody, err := unquoteUnicode(string(body))

	if err != nil {
		log.Println("Error al decodificar la respuesta:", err)
		return
	}

	re := regexp.MustCompile(`\\"([^\\"]+)\\"`)
	decodedBody = re.ReplaceAllString(decodedBody, `<b>$1</b>`)

	helper.ResponseWithJson(w, http.StatusOK, decodedBody)
	return

}

// helper.ResponseWithJson(w, http.StatusOK, "HELLO XDXDXDDDD")
// responseMap := make(map[string]interface{})
// responseMap["result"] = "success"
// responseMap["response"] = decodedBody
// responseMap["Contexto"] = question.Contexto
