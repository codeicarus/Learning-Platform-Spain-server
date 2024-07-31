package models

import (
	"github.com/globalsign/mgo/bson"
)

type HistorialesCP struct {
	Id                  bson.ObjectId `bson:"_id" json:"id"`
	Tipo                string        `bson:"tipo" json:"tipo"`
	Fecha               int64         `bson:"fecha" json:"fecha"`
	CasosPracticos      []string      `bson:"casos_practicos" json:"casos_practicos"`
	IdUsuario           bson.ObjectId `bson:"id_usuario" json:"id_usuario"`
	IdCasoPractico      bson.ObjectId `bson:"id_caso_practico" json:"id_caso_practico"`
	NumeroPreguntas     int           `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool          `bson:"respuesta_automatica" json:"respuesta_automatica"`
	TiempoTranscurrido  int           `bson:"tiempo_transcurrido" json:"tiempo_transcurrido"`
	Terminado           bool          `bson:"terminado" json:"terminado"`
	PreguntasTotales    int           `bson:"preguntas_totales" json:"preguntas_totales"`
	PreguntasAcertadas  int           `bson:"preguntas_acertadas" json:"preguntas_acertadas"`
	PreguntasFalladas   int           `bson:"preguntas_falladas" json:"preguntas_falladas"`
	PreguntasBlancas    int           `bson:"preguntas_blancas" json:"preguntas_blancas"`
	Puntuacion          float32       `bson:"puntuacion" json:"puntuacion"`
}
type HistorialesCPADevolver struct {
	Historial          HistorialesCP          `bson:"historial" json:"historial"`
	CasosPracticos     []SupuestosCPADevolver `bson:"casos_practicos" json:"casos_practicos"`
	Preguntas          []PreguntasCPADevolver `bson:"preguntas" json:"preguntas"`
	RespuestasMarcadas map[string]string      `bson:"respuestas_marcadas" json:"respuestas_marcadas"`
}

func (n *HistorialesCP) InsertHistorial(historial HistorialesCP) error {
	return Insert(db, "historiales_cp", historial)
}

func (n *HistorialesCP) FindAllHistorialesCP() ([]HistorialesCP, error) {
	var result []HistorialesCP
	err := FindAll(db, "historiales_cp", nil, nil, &result)
	return result, err
}

func (n *HistorialesCP) FindHistorialById(id string) (HistorialesCP, error) {
	var result HistorialesCP
	err := FindOne(db, "historiales_cp", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *HistorialesCP) FindHistorialByIdUsuario(id_usuario string, selector interface{}) ([]HistorialesCP, error) {
	var result []HistorialesCP
	err := FindAll(db, "historiales_cp", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario)}, selector, &result)
	return result, err
}

func (n *HistorialesCP) UpdateHistorial(historial HistorialesCP) error {
	return Update(db, "historiales_cp", bson.M{"_id": historial.Id}, historial)
}

func (n *HistorialesCP) RemoveHistorial(id string) error {
	return Remove(db, "historiales_cp", bson.M{"_id": bson.ObjectIdHex(id)})
}


func (n *HistorialesCP) RemovePreguntaHistorialCP(id string) error {
	return RemoveAll(db, "historiales_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id)})
}