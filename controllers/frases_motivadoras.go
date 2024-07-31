package controllers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"test/helper"
	"test/models"
	"time"
)

var (
	daoFrasesMotivadoras = models.FrasesMotivadoras{}
)

func GetFraseMotivadora(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params map[string]string
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		helper.ResponseWithJson(w, http.StatusBadRequest, map[string]string{"result": "error", "error": "Invalid request payload"})
		return
	}

	usuario, err := daoUsuarios.FindUsuarioByEmail(strings.ToLower(strings.TrimSpace(params["email"])))
	if err != nil {
		helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": "Usuario / Password incorrectos"})
		return
	} else {

		currentTime := time.Now().Unix()

		if currentTime > usuario.LastFrase+172800 {
			//if currentTime > usuario.LastFrase+1 {
			usuario.LastFrase = currentTime
			if err := daoUsuarios.UpdateUsuario(usuario); err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

			var frases_motivadoras []models.FrasesMotivadoras

			var turnos_a_buscar []string

			t := time.Now()
			h := t.Hour()

			if h >= 7 && h <= 11 {
				turnos_a_buscar = []string{"O", "L", "M"}
			} else if h >= 15 && h <= 20 {
				turnos_a_buscar = []string{"O", "L", "T"}
			} else {
				turnos_a_buscar = []string{"O", "L"}
			}

			frases_motivadoras, err := daoFrasesMotivadoras.FindFraseMotivadoraByTurno(turnos_a_buscar)
			if err != nil {
				helper.ResponseWithJson(w, http.StatusInternalServerError, map[string]string{"result": "error", "error": err.Error()})
				return
			}

			rand.Seed(time.Now().UnixNano())
			for i := len(frases_motivadoras) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				frases_motivadoras[i], frases_motivadoras[j] = frases_motivadoras[j], frases_motivadoras[i]
			}
			helper.ResponseWithJson(w, http.StatusOK, frases_motivadoras[0])
		} else {
			helper.ResponseWithJson(w, http.StatusOK, map[string]string{"result": "no_need"})
		}
	}

}
