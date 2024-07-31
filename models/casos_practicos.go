package models

import (
	"github.com/globalsign/mgo/bson"
)

type CasosPracticos struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Oficial     bool          `bson:"oficial" json:"oficial"`
	AnioOficial string        `bson:"anio_oficial" json:"anio_oficial"`
	Texto       string        `bson:"texto" json:"texto"`
}

func (n *CasosPracticos) InsertCasoPractico(caso_practico CasosPracticos) error {
	return Insert(db, "casos_practicos", caso_practico)
}

func (n *CasosPracticos) FindAllCasosPracticos() ([]CasosPracticos, error) {
	var result []CasosPracticos
	err := FindAll(db, "casos_practicos", nil, nil, &result)
	return result, err
}

func (n *CasosPracticos) FindCasoPracticoById(id string) (CasosPracticos, error) {
	var result CasosPracticos
	err := FindOne(db, "casos_practicos", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *CasosPracticos) FindCasoPracticoByName(name string) (CasosPracticos, error) {
	var result CasosPracticos
	err := FindOne(db, "casos_practicos", bson.M{"name": name}, nil, &result)
	return result, err
}

func (n *CasosPracticos) UpdateCasoPractico(caso_practico CasosPracticos) error {
	return Update(db, "casos_practicos", bson.M{"_id": caso_practico.Id}, caso_practico)
}

func (n *CasosPracticos) RemoveCasoPractico(id string) error {
	return Remove(db, "casos_practicos", bson.M{"_id": bson.ObjectIdHex(id)})
}
