package models

import "github.com/globalsign/mgo/bson"

type HistorialesCPPreguntas struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	IdHistorial bson.ObjectId `bson:"id_historial" json:"id_historial"`
	IdPregunta  bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
	IdRespuesta bson.ObjectId `bson:"id_respuesta,omitempty" json:"id_respuesta,omitempty"`
	Correcta    bool          `bson:"correcta,omitempty" json:"correcta,omitempty"`
}

func (n *HistorialesCPPreguntas) InsertHistorialPregunta(historial_pregunta HistorialesCPPreguntas) error {
	return Insert(db, "historiales_cp_preguntas", historial_pregunta)
}

func (n *HistorialesCPPreguntas) FindAllHistorialesCPPreguntas() ([]HistorialesCPPreguntas, error) {
	var result []HistorialesCPPreguntas
	err := FindAll(db, "historiales_cp_preguntas", nil, nil, &result)
	return result, err
}

func (n *HistorialesCPPreguntas) FindHistorialPreguntaById(id string) (HistorialesCPPreguntas, error) {
	var result HistorialesCPPreguntas
	err := FindOne(db, "historiales_cp_preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *HistorialesCPPreguntas) FindHistorialPreguntaByIdsHistorialesCP(ids_historiales []bson.ObjectId) ([]HistorialesCPPreguntas, error) {
	var result []HistorialesCPPreguntas
	err := FindAllOrder(db, "historiales_cp_preguntas", bson.M{"id_historial": bson.M{"$in": ids_historiales}}, nil, &result, "-id_historial")
	return result, err
}

func (n *HistorialesCPPreguntas) FindHistorialPreguntaByIdHistorial(id_historial string) ([]HistorialesCPPreguntas, error) {
	var result []HistorialesCPPreguntas
	err := FindAll(db, "historiales_cp_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial)}, nil, &result)
	return result, err
}

func (n *HistorialesCPPreguntas) FindHistorialPreguntaByIdHistorialIdPregunta(id_historial string, id_pregunta string) (HistorialesCPPreguntas, error) {
	var result HistorialesCPPreguntas
	err := FindOne(db, "historiales_cp_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial), "id_pregunta": bson.ObjectIdHex(id_pregunta)}, nil, &result)
	return result, err
}
func (n *HistorialesCPPreguntas) UpdateHistorialPregunta(historial_pregunta HistorialesCPPreguntas) error {
	return Update(db, "historiales_cp_preguntas", bson.M{"_id": historial_pregunta.Id}, historial_pregunta)
}

func (n *HistorialesCPPreguntas) RemoveHistorialPregunta(id string) error {
	return Remove(db, "historiales_cp_preguntas", bson.M{"_id": bson.ObjectIdHex(id)})
}
