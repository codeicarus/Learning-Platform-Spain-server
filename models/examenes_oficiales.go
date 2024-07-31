package models

import (
	"github.com/globalsign/mgo/bson"
)

type ExamenesOficiales struct {
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}
type OficialesYSimulacros struct {
	Oficiales  []ExamenesOficiales `bson:"oficiales" json:"oficiales"`
	Simulacros []Simulacros        `bson:"simulacros" json:"simulacros"`
	Niveles    []Niveles           `bson:"niveles" json:"niveles"`
}

func (n *ExamenesOficiales) InsertExamenesOficial(examen_oficial ExamenesOficiales) error {
	return Insert(db, "examenes_oficiales", examen_oficial)
}

func (n *ExamenesOficiales) FindAllExamenesOficiales() ([]ExamenesOficiales, error) {
	var result []ExamenesOficiales
	err := FindAll(db, "examenes_oficiales", nil, nil, &result)
	return result, err
}

func (n *ExamenesOficiales) FindExamenesOficialById(id string) (ExamenesOficiales, error) {
	var result ExamenesOficiales
	err := FindOne(db, "examenes_oficiales", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *ExamenesOficiales) UpdateExamenesOficial(examen_oficial ExamenesOficiales) error {
	return Update(db, "examenes_oficiales", bson.M{"_id": examen_oficial.Id}, examen_oficial)
}

func (n *ExamenesOficiales) RemoveExamenesOficial(id string) error {
	return Remove(db, "examenes_oficiales", bson.M{"_id": bson.ObjectIdHex(id)})
}
