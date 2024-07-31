package models

import "github.com/globalsign/mgo/bson"

type HistorialesPreguntas struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	IdHistorial bson.ObjectId `bson:"id_historial" json:"id_historial"`
	IdPregunta  bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
	IdRespuesta bson.ObjectId `bson:"id_respuesta,omitempty" json:"id_respuesta,omitempty"`
	Correcta    bool          `bson:"correcta,omitempty" json:"correcta,omitempty"`
}

type HistorialesPreguntasPipe struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	IdHistorial bson.ObjectId `bson:"id_historial" json:"id_historial"`
	IdPregunta  bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
	IdRespuesta bson.ObjectId `bson:"id_respuesta,omitempty" json:"id_respuesta,omitempty"`
	Correcta    bool          `bson:"correcta,omitempty" json:"correcta,omitempty"`
	Preguntas   []Preguntas   `bson:"preguntas" json:"preguntas"`
}

func (n *HistorialesPreguntas) InsertHistorialPregunta(historial_pregunta HistorialesPreguntas) error {
	return Insert(db, "historiales_preguntas", historial_pregunta)
}

func (n *HistorialesPreguntas) FindAllHistorialesPreguntas() ([]HistorialesPreguntas, error) {
	var result []HistorialesPreguntas
	err := FindAll(db, "historiales_preguntas", nil, nil, &result)
	return result, err
}

func (n *HistorialesPreguntas) FindHistorialPreguntaById(id string) (HistorialesPreguntas, error) {
	var result HistorialesPreguntas
	err := FindOne(db, "historiales_preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}
func (n *HistorialesPreguntas) FindHistorialPreguntaByIdsHistoriales(ids_historiales []bson.ObjectId) ([]HistorialesPreguntas, error) {
	var result []HistorialesPreguntas
	err := FindAllOrder(db, "historiales_preguntas", bson.M{"id_historial": bson.M{"$in": ids_historiales}}, bson.M{"_id": 0}, &result, "-id_historial")
	return result, err
}

func (n *HistorialesPreguntas) FindHistorialPreguntaByIdHistorial(id_historial string) ([]HistorialesPreguntas, error) {
	var result []HistorialesPreguntas
	err := FindAllOrder(db, "historiales_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial)}, nil, &result, "_id")
	return result, err
}
func (n *HistorialesPreguntas) FindHistorialPreguntaByIdHistorialPipe(id_historial string) ([]HistorialesPreguntasPipe, error) {
	var result []HistorialesPreguntasPipe
	err := FindAllWPreguntas(db, "historiales_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial)}, nil, &result)
	return result, err
}

func (n *HistorialesPreguntas) FindHistorialPreguntaByIdHistorialIdPregunta(id_historial string, id_pregunta string) (HistorialesPreguntas, error) {
	var result HistorialesPreguntas
	err := FindOne(db, "historiales_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial), "id_pregunta": bson.ObjectIdHex(id_pregunta)}, nil, &result)
	return result, err
}
func (n *HistorialesPreguntas) UpdateHistorialPregunta(historial_pregunta HistorialesPreguntas) error {
	return Update(db, "historiales_preguntas", bson.M{"_id": historial_pregunta.Id}, historial_pregunta)
}

func (n *HistorialesPreguntas) RemoveHistorialPregunta(id string) error {
	return Remove(db, "historiales_preguntas", bson.M{"_id": bson.ObjectIdHex(id)})
}
