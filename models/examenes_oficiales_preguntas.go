package models

import "github.com/globalsign/mgo/bson"

type ExamenesOficialesPreguntas struct {
	Id              bson.ObjectId `bson:"_id" json:"id"`
	IdExamenOficial bson.ObjectId `bson:"id_examen_oficial" json:"id_examen_oficial"`
	IdPregunta      bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
	Orden           int           `bson:"orden" json:"orden"`
}

func (n *ExamenesOficialesPreguntas) InsertExamenOficialPregunta(examen_oficial_pregunta ExamenesOficialesPreguntas) error {
	return Insert(db, "examenes_oficiales_preguntas", examen_oficial_pregunta)
}

func (n *ExamenesOficialesPreguntas) FindAllExamenesOficialesPreguntas() ([]ExamenesOficialesPreguntas, error) {
	var result []ExamenesOficialesPreguntas
	err := FindAll(db, "examenes_oficiales_preguntas", nil, nil, &result)
	return result, err
}

func (n *ExamenesOficialesPreguntas) FindPreguntaByIdExamenOficial(ids_examenes_oficiales []bson.ObjectId) ([]ExamenesOficialesPreguntas, error) {
	var result []ExamenesOficialesPreguntas
	err := FindAll(db, "examenes_oficiales_preguntas", bson.M{"id_examen_oficial": bson.M{"$in": ids_examenes_oficiales}}, nil, &result)
	return result, err
}

func (n *ExamenesOficialesPreguntas) FindPreguntaByIdExamenOficialOrder(ids_examenes_oficiales []bson.ObjectId) ([]ExamenesOficialesPreguntas, error) {
	var result []ExamenesOficialesPreguntas
	err := FindAllOrder(db, "examenes_oficiales_preguntas", bson.M{"id_examen_oficial": bson.M{"$in": ids_examenes_oficiales}}, nil, &result, "orden")
	return result, err
}

func (n *ExamenesOficialesPreguntas) FindExamenOficialPreguntaById(id string) (ExamenesOficialesPreguntas, error) {
	var result ExamenesOficialesPreguntas
	err := FindOne(db, "examenes_oficiales_preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *ExamenesOficialesPreguntas) FindExamenOficialPreguntaByIdHistorial(id_historial string) ([]ExamenesOficialesPreguntas, error) {
	var result []ExamenesOficialesPreguntas
	err := FindAll(db, "examenes_oficiales_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial)}, nil, &result)
	return result, err
}

func (n *ExamenesOficialesPreguntas) UpdateExamenOficialPregunta(examen_oficial_pregunta ExamenesOficialesPreguntas) error {
	return Update(db, "examenes_oficiales_preguntas", bson.M{"_id": examen_oficial_pregunta.Id}, examen_oficial_pregunta)
}

func (n *ExamenesOficialesPreguntas) RemoveExamenOficialPregunta(id string) error {
	return Remove(db, "examenes_oficiales_preguntas", bson.M{"_id": bson.ObjectIdHex(id)})
}
