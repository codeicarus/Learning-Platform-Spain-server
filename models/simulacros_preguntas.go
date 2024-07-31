package models

import "github.com/globalsign/mgo/bson"

type SimulacrosPreguntas struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	IdSimulacro bson.ObjectId `bson:"id_simulacro" json:"id_simulacro"`
	IdPregunta  bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
	Orden       int           `bson:"orden" json:"orden"`
}

func (n *SimulacrosPreguntas) InsertSimulacroPregunta(simulacro_pregunta SimulacrosPreguntas) error {
	return Insert(db, "simulacros_preguntas", simulacro_pregunta)
}

func (n *SimulacrosPreguntas) FindAllSimulacrosPreguntas() ([]SimulacrosPreguntas, error) {
	var result []SimulacrosPreguntas
	err := FindAll(db, "simulacros_preguntas", nil, nil, &result)
	return result, err
}

func (n *SimulacrosPreguntas) FindPreguntaByIdSimulacro(ids_simulacros []bson.ObjectId) ([]SimulacrosPreguntas, error) {
	var result []SimulacrosPreguntas
	err := FindAll(db, "simulacros_preguntas", bson.M{"id_simulacro": bson.M{"$in": ids_simulacros}}, nil, &result)
	return result, err
}

func (n *SimulacrosPreguntas) FindPreguntaByIdSimulacroOrder(ids_simulacros []bson.ObjectId) ([]SimulacrosPreguntas, error) {
	var result []SimulacrosPreguntas
	err := FindAllOrder(db, "simulacros_preguntas", bson.M{"id_simulacro": bson.M{"$in": ids_simulacros}}, nil, &result, "orden")
	return result, err
}

func (n *SimulacrosPreguntas) FindSimulacroPreguntaById(id string) (SimulacrosPreguntas, error) {
	var result SimulacrosPreguntas
	err := FindOne(db, "simulacros_preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *SimulacrosPreguntas) FindSimulacroPreguntaByIdHistorial(id_historial string) ([]SimulacrosPreguntas, error) {
	var result []SimulacrosPreguntas
	err := FindAll(db, "simulacros_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial)}, nil, &result)
	return result, err
}

func (n *SimulacrosPreguntas) UpdateSimulacroPregunta(simulacro_pregunta SimulacrosPreguntas) error {
	return Update(db, "simulacros_preguntas", bson.M{"_id": simulacro_pregunta.Id}, simulacro_pregunta)
}

func (n *SimulacrosPreguntas) RemoveSimulacroPregunta(id string) error {
	return Remove(db, "simulacros_preguntas", bson.M{"_id": bson.ObjectIdHex(id)})
}
