package models

import (
	"github.com/globalsign/mgo/bson"
)

type Historiales struct {
	Id                  bson.ObjectId `bson:"_id" json:"id"`
	Tipo                string        `bson:"tipo" json:"tipo"`
	Fecha               int64         `bson:"fecha" json:"fecha"`
	Temas               []string      `bson:"temas" json:"temas"`
	Legislaciones       []string      `bson:"legislaciones" json:"legislaciones"`
	BasicosConfundidos  []string      `bson:"basicosconfundidos" json:"basicosconfundidos"`
	Oficiales           []string      `bson:"oficiales" json:"oficiales"`
	Simulacros          []string      `bson:"simulacros" json:"simulacros"`
	IdUsuario           bson.ObjectId `bson:"id_usuario" json:"id_usuario"`
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
type HistorialesPipe struct {
	Id                  bson.ObjectId          `bson:"_id" json:"id"`
	Tipo                string                 `bson:"tipo" json:"tipo"`
	Fecha               int64                  `bson:"fecha" json:"fecha"`
	Temas               []string               `bson:"temas" json:"temas"`
	Legislaciones       []string               `bson:"legislaciones" json:"legislaciones"`
	BasicosConfundidos  []string               `bson:"basicosconfundidos" json:"basicosconfundidos"`
	Oficiales           []string               `bson:"oficiales" json:"oficiales"`
	Simulacros          []string               `bson:"simulacros" json:"simulacros"`
	IdUsuario           bson.ObjectId          `bson:"id_usuario" json:"id_usuario"`
	NumeroPreguntas     int                    `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool                   `bson:"respuesta_automatica" json:"respuesta_automatica"`
	TiempoTranscurrido  int                    `bson:"tiempo_transcurrido" json:"tiempo_transcurrido"`
	Terminado           bool                   `bson:"terminado" json:"terminado"`
	PreguntasTotales    int                    `bson:"preguntas_totales" json:"preguntas_totales"`
	PreguntasAcertadas  int                    `bson:"preguntas_acertadas" json:"preguntas_acertadas"`
	PreguntasFalladas   int                    `bson:"preguntas_falladas" json:"preguntas_falladas"`
	PreguntasBlancas    int                    `bson:"preguntas_blancas" json:"preguntas_blancas"`
	Puntuacion          float32                `bson:"puntuacion" json:"puntuacion"`
	Constestaciones     []HistorialesPreguntas `bson:"contestaciones" json:"contestaciones"`
}
type HistorialesADevolver struct {
	Historial          Historiales          `bson:"historial" json:"historial"`
	Preguntas          []PreguntasADevolver `bson:"preguntas" json:"preguntas"`
	RespuestasMarcadas map[string]string    `bson:"respuestas_marcadas" json:"respuestas_marcadas"`
}

func (n *Historiales) InsertHistorial(historial Historiales) error {
	return Insert(db, "historiales", historial)
}

func (n *Historiales) FindAllHistoriales() ([]Historiales, error) {
	var result []Historiales
	err := FindAll(db, "historiales", nil, nil, &result)
	return result, err
}

func (n *Historiales) FindHistorialById(id string) (Historiales, error) {
	var result Historiales
	err := FindOne(db, "historiales", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *Historiales) FindHistorialByIdUsuario(id_usuario string, selector interface{}) ([]Historiales, error) {
	var result []Historiales
	err := FindAll(db, "historiales", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario)}, selector, &result)
	return result, err
}
func (n *Historiales) FindHistorialByIdUsuarioAndDate(id_usuario string, date int64, selector interface{}) ([]Historiales, error) {
	var result []Historiales
	err := FindAll(db, "historiales", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario), "fecha": bson.M{"$gte": date}}, selector, &result)
	return result, err
}

func (n *Historiales) UpdateHistorial(historial Historiales) error {
	return Update(db, "historiales", bson.M{"_id": historial.Id}, historial)
}

func (n *Historiales) RemoveHistorial(id string) error {
	return Remove(db, "historiales", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n *Historiales) RemovePreguntaHistorial(id string) error {
	return RemoveAll(db, "historiales_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id)})
}
