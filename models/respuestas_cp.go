package models

import "github.com/globalsign/mgo/bson"

type RespuestasCP struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	Respuesta  string        `bson:"respuesta" json:"respuesta"`
	Correcta   bool          `bson:"correcta" json:"correcta"`
	IdPregunta bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
}
type RespuestasCPADevolver struct {
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Respuesta string        `bson:"respuesta" json:"respuesta"`
	Correcta  bool          `bson:"correcta" json:"correcta"`
}

func (n *RespuestasCP) InsertRespuesta(respuesta RespuestasCP) error {
	return Insert(db, "respuestas_cp", respuesta)
}

func (n *RespuestasCP) FindAllRespuestasCP() ([]RespuestasCP, error) {
	var result []RespuestasCP
	err := FindAll(db, "respuestas_cp", nil, nil, &result)
	return result, err
}

func (n *RespuestasCP) FindRespuestaById(id string) (RespuestasCP, error) {
	var result RespuestasCP
	err := FindOne(db, "respuestas_cp", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *RespuestasCP) FindRespuestaCorrectaByIdPregunta(id_pregunta string) (RespuestasCP, error) {
	var result RespuestasCP
	err := FindOne(db, "respuestas_cp", bson.M{"id_pregunta": bson.ObjectIdHex(id_pregunta), "correcta": true}, nil, &result)
	return result, err
}

func (n *RespuestasCP) FindRespuestaByIdPregunta(id_pregunta string) ([]RespuestasCP, error) {
	var result []RespuestasCP
	err := FindAll(db, "respuestas_cp", bson.M{"id_pregunta": bson.ObjectIdHex(id_pregunta)}, nil, &result)
	return result, err
}

func (n *RespuestasCP) UpdateRespuesta(respuesta RespuestasCP) error {
	return Update(db, "respuestas_cp", bson.M{"_id": respuesta.Id}, respuesta)
}

func (n *RespuestasCP) RemoveRespuesta(id string) error {
	return Remove(db, "respuestas_cp", bson.M{"_id": bson.ObjectIdHex(id)})
}
