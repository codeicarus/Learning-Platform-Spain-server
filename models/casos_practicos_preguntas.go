package models

import (
	"github.com/globalsign/mgo/bson"
)

type CasosPracticosPreguntas struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`
	IdCasoPractico bson.ObjectId `bson:"id_caso_practico" json:"id_caso_practico"`
	IdPregunta     bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
	Orden          int           `bson:"orden" json:"orden"`
}

func (n *CasosPracticosPreguntas) InsertCasoPracticoPregunta(caso_practico_pregunta CasosPracticosPreguntas) error {
	// log.Println(caso_practico_pregunta)
	return Insert(db, "casos_practicos_preguntas", caso_practico_pregunta)
}

func (n *CasosPracticosPreguntas) FindAllCasosPracticosPreguntas() ([]CasosPracticosPreguntas, error) {
	var result []CasosPracticosPreguntas
	err := FindAll(db, "casos_practicos_preguntas", nil, nil, &result)
	return result, err
}

func (n *CasosPracticosPreguntas) FindPreguntaByIdCasoPractico(ids_casos_practicos []bson.ObjectId) ([]CasosPracticosPreguntas, error) {
	var result []CasosPracticosPreguntas
	err := FindAll(db, "casos_practicos_preguntas", bson.M{"id_caso_practico": bson.M{"$in": ids_casos_practicos}}, nil, &result)
	return result, err
}

func (n *CasosPracticosPreguntas) FindPreguntaByIdCasoPracticoOrder(ids_casos_practicos []bson.ObjectId) ([]CasosPracticosPreguntas, error) {
	var result []CasosPracticosPreguntas
	err := FindAllOrder(db, "casos_practicos_preguntas", bson.M{"id_caso_practico": bson.M{"$in": ids_casos_practicos}}, nil, &result, "orden")
	return result, err
}

func (n *CasosPracticosPreguntas) FindCasoPracticoPreguntaById(id string) (CasosPracticosPreguntas, error) {
	var result CasosPracticosPreguntas
	err := FindOne(db, "casos_practicos_preguntas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *CasosPracticosPreguntas) FindCasoPracticoPreguntaByIdHistorial(id_historial string) ([]CasosPracticosPreguntas, error) {
	var result []CasosPracticosPreguntas
	err := FindAll(db, "casos_practicos_preguntas", bson.M{"id_historial": bson.ObjectIdHex(id_historial)}, nil, &result)
	return result, err
}

func (n *CasosPracticosPreguntas) UpdateCasoPracticoPregunta(caso_practico_pregunta CasosPracticosPreguntas) error {
	return Update(db, "casos_practicos_preguntas", bson.M{"_id": caso_practico_pregunta.Id}, caso_practico_pregunta)
}

func (n *CasosPracticosPreguntas) RemoveCasoPracticoPregunta(id string) error {
	return Remove(db, "casos_practicos_preguntas", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n* CasosPracticosPreguntas) FindCasoPracticoByCasoPracticoPregunta(id_caso_practico string, id_pregunta string) (error) {
	var result CasosPracticosPreguntas
	err := FindOne(db, "casos_practicos_preguntas", bson.M{"id_caso_practico": bson.ObjectIdHex(id_caso_practico), "id_pregunta": bson.ObjectIdHex(id_pregunta) }, nil, &result)
	return err
}