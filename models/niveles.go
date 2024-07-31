package models

import (
	"github.com/globalsign/mgo/bson"
)

type Niveles struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Abreviacion string        `bson:"abreviacion" json:"abreviacion"`
}

func (n *Niveles) InsertNivel(nivel Niveles) error {
	return Insert(db, "niveles", nivel)
}

func (n *Niveles) FindAllNiveles() ([]Niveles, error) {
	var result []Niveles
	err := FindAll(db, "niveles", nil, nil, &result)
	return result, err
}

func (n *Niveles) FindNivelById(id string) (Niveles, error) {
	var result Niveles
	err := FindOne(db, "niveles", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Niveles) FindNivelByAbreviacion(abreviacion string) (Niveles, error) {
	var result Niveles
	err := FindOne(db, "niveles", bson.M{"abreviacion": abreviacion}, nil, &result)
	return result, err
}

func (n *Niveles) UpdateNivel(nivel Niveles) error {
	return Update(db, "niveles", bson.M{"_id": nivel.Id}, nivel)
}

func (n *Niveles) RemoveNivel(id string) error {
	return Remove(db, "niveles", bson.M{"_id": bson.ObjectIdHex(id)})
}
