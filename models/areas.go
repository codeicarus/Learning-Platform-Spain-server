package models

import (
	"github.com/globalsign/mgo/bson"
)

type Areas struct {
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
}
type AreasWLegislaciones struct {
	Id            bson.ObjectId   `bson:"_id" json:"id"`
	Name          string          `bson:"name" json:"name"`
	Legislaciones []Legislaciones `bson:"legislaciones" json:"legislaciones"`
}
type AreasWBasicosConfundidos struct {
	Id                 bson.ObjectId        `bson:"_id" json:"id"`
	Name               string               `bson:"name" json:"name"`
	BasicosConfundidos []BasicosConfundidos `bson:"basicosconfundidos" json:"basicosconfundidos"`
}
type AreasWTemas struct {
	Id    bson.ObjectId `bson:"_id" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Temas []Temas       `bson:"temas" json:"temas"`
}
type AreasWBloques struct {
	Id      bson.ObjectId `bson:"_id" json:"id"`
	Name    string        `bson:"name" json:"name"`
	Bloques []Bloques     `bson:"bloques" json:"bloques"`
}

func (n *Areas) InsertArea(area Areas) error {
	return Insert(db, "areas", area)
}

func (n *Areas) FindAllAreas() ([]Areas, error) {
	var result []Areas
	err := FindAllOrder(db, "areas", nil, nil, &result, "ordenacion")
	return result, err
}

func (n *Areas) FindAreaById(id string) (Areas, error) {
	var result Areas
	err := FindOne(db, "areas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Areas) FindAreaByName(name string) (Areas, error) {
	var result Areas
	err := FindOne(db, "areas", bson.M{"name": name}, nil, &result)
	return result, err
}

func (n *Areas) UpdateArea(area Areas) error {
	return Update(db, "areas", bson.M{"_id": area.Id}, area)
}

func (n *Areas) RemoveArea(id string) error {
	return Remove(db, "areas", bson.M{"_id": bson.ObjectIdHex(id)})
}
