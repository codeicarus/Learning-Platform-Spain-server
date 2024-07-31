package models

import (
	"github.com/globalsign/mgo/bson"
)

type Simulacros struct {
	Id      bson.ObjectId `bson:"_id" json:"id"`
	Name    string        `bson:"name" json:"name"`
	IdNivel string        `bson:"id_nivel" json:"id_nivel"`
}

func (n *Simulacros) InsertSimulacro(simulacro Simulacros) error {
	return Insert(db, "simulacros", simulacro)
}

func (n *Simulacros) FindAllSimulacros() ([]Simulacros, error) {
	var result []Simulacros
	err := FindAll(db, "simulacros", nil, nil, &result)
	return result, err
}

func (n *Simulacros) FindSimulacroById(id string) (Simulacros, error) {
	var result Simulacros
	err := FindOne(db, "simulacros", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Simulacros) UpdateSimulacro(simulacro Simulacros) error {
	return Update(db, "simulacros", bson.M{"_id": simulacro.Id}, simulacro)
}

func (n *Simulacros) RemoveSimulacro(id string) error {
	return Remove(db, "simulacros", bson.M{"_id": bson.ObjectIdHex(id)})
}
