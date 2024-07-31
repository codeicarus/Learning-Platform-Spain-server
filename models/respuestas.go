package models

import "github.com/globalsign/mgo/bson"

type Respuestas struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Respuesta  string        `bson:"respuesta" json:"respuesta"`
	Correcta   bool          `bson:"correcta" json:"correcta"`
	IdPregunta bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
}
type RespuestasADevolver struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Respuesta string        `bson:"respuesta" json:"respuesta"`
	Correcta  bool          `bson:"correcta" json:"correcta"`
}

func (n *Respuestas) InsertRespuesta(respuesta Respuestas) error {
	return Insert(db, "respuestas", respuesta)
}

func (n *Respuestas) FindAllRespuestas() ([]Respuestas, error) {
	var result []Respuestas
	err := FindAll(db, "respuestas", nil, nil, &result)
	return result, err
}

func (n *Respuestas) FindRespuestaById(id string) (Respuestas, error) {
	var result Respuestas
	err := FindOne(db, "respuestas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Respuestas) FindRespuestaCorrectaByIdPregunta(id_pregunta string) (Respuestas, error) {
	var result Respuestas
	err := FindOne(db, "respuestas", bson.M{"id_pregunta": bson.ObjectIdHex(id_pregunta), "correcta": true}, nil, &result)
	return result, err
}

func (n *Respuestas) FindRespuestasByIdsPregunta(ids []bson.ObjectId) ([]Respuestas, error) {
	var result []Respuestas
	err := FindAll(db, "respuestas", bson.M{"id_pregunta": bson.M{"$in": ids}}, nil, &result)
	return result, err
}
func (n *Respuestas) FindRespuestaByIdPregunta(id_pregunta string) ([]Respuestas, error) {
	var result []Respuestas
	err := FindAll(db, "respuestas", bson.M{"id_pregunta": bson.ObjectIdHex(id_pregunta)}, nil, &result)
	return result, err
}

func (n *Respuestas) UpdateRespuesta(respuesta Respuestas) error {
	return Update(db, "respuestas", bson.M{"_id": respuesta.Id}, respuesta)
}

func (n *Respuestas) RemoveRespuesta(id string) error {
	return Remove(db, "respuestas", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n *Respuestas) RemoveRespuestasByIdPregunta(id_pregunta string) error {
	return RemoveAll(db, "respuestas", bson.M{"id_pregunta": bson.ObjectIdHex(id_pregunta)})
}