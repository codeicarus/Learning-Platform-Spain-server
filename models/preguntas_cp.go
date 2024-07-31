package models

import "github.com/globalsign/mgo/bson"

type PreguntasCP struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Pregunta    string        `bson:"pregunta" json:"pregunta"`
	Explicacion string        `bson:"explicacion" json:"explicacion"`
}
type SupuestosCPADevolver struct {
	Id          bson.ObjectId          `bson:"_id" json:"id"`
	Name        string                 `bson:"name" json:"name"`
	Oficial     bool                   `bson:"oficial" json:"oficial"`
	AnioOficial string                 `bson:"anio_oficial" json:"anio_oficial"`
	Texto       string                 `bson:"texto" json:"texto"`
	Preguntas   []PreguntasCPADevolver `bson:"preguntas" json:"preguntas"`
}
type PreguntasCPADevolver struct {
	Id          bson.ObjectId           `bson:"_id" json:"id"`
	Pregunta    string                  `bson:"pregunta" json:"pregunta"`
	Explicacion string                  `bson:"explicacion" json:"explicacion"`
	Respuestas  []RespuestasCPADevolver `bson:"respuestas" json:"respuestas"`
}
type PreguntasCPRespuestas struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Pregunta   string        `bson:"pregunta" json:"pregunta"`
	Respuestas []string      `bson:"respuestas" json:"respuestas"`
}
type SavePreguntaCP struct {
	Email    string                   `bson:"email" json:"email"`
	Pregunta PreguntaCPRecibidaEditar `bson:"pregunta" json:"pregunta"`
}
type PreguntaCPRecibidaEditar struct {
	Id          bson.ObjectId           `bson:"_id" json:"id"`
	Pregunta    string                  `bson:"pregunta" json:"pregunta"`
	Explicacion string                  `bson:"explicacion" json:"explicacion"`
	Respuestas  []RespuestasCPADevolver `bson:"respuestas" json:"respuestas"`
}
type CreateTestCP struct {
	CasosPracticos      []string `bson:"casos_practicos" json:"casos_practicos"`
	NumeroPreguntas     int      `bson:"numero_preguntas" json:"numero_preguntas"`
	RespuestaAutomatica bool     `bson:"respuesta_automatica" json:"respuesta_automatica"`
	Tipo                string   `bson:"tipo" json:"tipo"`
	Email               string   `bson:"email" json:"email"`
}
type CorregirTestCP struct {
	IdTest             string `bson:"id_test" json:"id_test"`
	RespuestasMarcadas string `bson:"respuestas_marcadas" json:"respuestas_marcadas"`
	Tiempo             int    `bson:"tiempo" json:"tiempo"`
}

func (n *PreguntasCP) InsertPregunta(pregunta PreguntasCP) error {
	return Insert(db, "preguntas_cp", pregunta)
}

func (n *PreguntasCP) FindAllPreguntas() ([]PreguntasCP, error) {
	var result []PreguntasCP
	err := FindAll(db, "preguntas_cp", nil, nil, &result)
	return result, err
}

func (n *PreguntasCP) FindPreguntaById(id string) (PreguntasCP, error) {
	var result PreguntasCP
	err := FindOne(db, "preguntas_cp", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *PreguntasCP) FindPreguntaByPregunta(pregunta string) (PreguntasCP, error) {
	var result PreguntasCP
	err := FindOne(db, "preguntas_cp", bson.M{"pregunta": pregunta}, nil, &result)
	return result, err
}

func (n *PreguntasCP) UpdatePregunta(pregunta PreguntasCP) error {
	return Update(db, "preguntas_cp", bson.M{"_id": pregunta.Id}, pregunta)
}

func (n *PreguntasCP) RemovePregunta(id string) error {
	return Remove(db, "preguntas_cp", bson.M{"_id": bson.ObjectIdHex(id)})
}
