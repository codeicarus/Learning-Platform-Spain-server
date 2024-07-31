package models

import "github.com/globalsign/mgo/bson"

type Bloques struct {
	Id     bson.ObjectId   `bson:"_id" json:"id"`
	Name   string          `bson:"name" json:"name"`
	IdArea bson.ObjectId   `bson:"id_area" json:"id_area"`
	Temas  []bson.ObjectId `bson:"temas" json:"temas"`
}

func (n *Bloques) InsertBloque(bloque Bloques) error {
	return Insert(db, "bloques", bloque)
}

func (n *Bloques) FindAllBloques() ([]Bloques, error) {
	var result []Bloques
	err := FindAll(db, "bloques", nil, nil, &result)
	return result, err
}

func (n *Bloques) FindBloqueById(id string) (Bloques, error) {
	var result Bloques
	err := FindOne(db, "bloques", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Bloques) FindBloqueByIdArea(id_area string) ([]Bloques, error) {
	var result []Bloques
	err := FindAll(db, "bloques", bson.M{"id_area": bson.ObjectIdHex(id_area)}, nil, &result)
	return result, err
}

func (n *Bloques) FindBasicoConfundidoByAbreviacion(abreviacion string) (Bloques, error) {
	var result Bloques
	err := FindOne(db, "bloques", bson.M{"abreviacion": abreviacion}, nil, &result)
	return result, err
}

func (n *Bloques) UpdateBloque(bloque Bloques) error {
	return Update(db, "bloques", bson.M{"_id": bloque.Id}, bloque)
}

func (n *Bloques) RemoveBloque(id string) error {
	return Remove(db, "bloques", bson.M{"_id": bson.ObjectIdHex(id)})
}
