package models

import (
	"github.com/globalsign/mgo/bson"
)

type Legislaciones struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Abreviacion string        `bson:"abreviacion" json:"abreviacion"`
	IdArea      bson.ObjectId `bson:"id_area" json:"id_area"`
}

type UpdateName struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	NewName     string        `bson:"newname" json:"newname"`
}

func (n *Legislaciones) InsertLegislacion(legislacion Legislaciones) error {
	return Insert(db, "legislaciones", legislacion)
}

func (n *Legislaciones) FindAllLegislaciones() ([]Legislaciones, error) {
	var result []Legislaciones
	err := FindAll(db, "legislaciones", nil, nil, &result)
	return result, err
}

func (n *Legislaciones) FindLegislacionById(id string) (Legislaciones, error) {
	var result Legislaciones
	err := FindOne(db, "legislaciones", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Legislaciones) FindLegislacionByIdArea(id_area string) ([]Legislaciones, error) {
	var result []Legislaciones
	err := FindAll(db, "legislaciones", bson.M{"id_area": bson.ObjectIdHex(id_area)}, nil, &result)
	return result, err
}

func (n *Legislaciones) FindLegislacionByAbreviacion(abreviacion string) (Legislaciones, error) {
	var result Legislaciones
	err := FindOne(db, "legislaciones", bson.M{"abreviacion": abreviacion}, nil, &result)
	return result, err
}

func (n *Legislaciones) FindLegislacionByAbreviacionAndAreaId(abreviacion string, id_area string) (Legislaciones, error) {
	var result Legislaciones
	err := FindOne(db, "legislaciones", bson.M{"abreviacion": abreviacion, "id_area": bson.ObjectIdHex(id_area)}, nil, &result)
	return result, err
}

func (n *Legislaciones) UpdateLegislacion(legislacion Legislaciones) error {
	return Update(db, "legislaciones", bson.M{"_id": legislacion.Id}, legislacion)
}

func (n *Legislaciones) RemoveLegislacion(id string) error {
	return Remove(db, "legislaciones", bson.M{"_id": bson.ObjectIdHex(id)})
}
