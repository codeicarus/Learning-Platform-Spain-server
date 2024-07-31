package models

import (
	"github.com/globalsign/mgo/bson"
)

type Temas struct {
	Id                 bson.ObjectId `bson:"_id" json:"id"`
	Name               string        `bson:"name" json:"name"`
	AbreviacionPublica string        `bson:"abreviacion_publica" json:"abreviacion_publica"`
	Abreviacion        string        `bson:"abreviacion" json:"abreviacion"`
	IdArea             bson.ObjectId `bson:"id_area" json:"id_area"`
}

func (n *Temas) InsertTema(tema Temas) error {
	return Insert(db, "temas", tema)
}

func (n *Temas) FindAllTemas() ([]Temas, error) {
	var result []Temas
	err := FindAll(db, "temas", nil, nil, &result)
	return result, err
}

func (n *Temas) FindTemaById(id string) (Temas, error) {
	var result Temas
	err := FindOne(db, "temas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Temas) FindTemaByIdArea(id_area string) ([]Temas, error) {
	var result []Temas
	err := FindAll(db, "temas", bson.M{"id_area": bson.ObjectIdHex(id_area)}, nil, &result)
	return result, err
}

func (n *Temas) FindBasicoConfundidoByAbreviacion(abreviacion string) (Temas, error) {
	var result Temas
	err := FindOne(db, "temas", bson.M{"abreviacion": abreviacion}, nil, &result)
	return result, err
}

func (n *Temas) UpdateTema(tema Temas) error {
	return Update(db, "temas", bson.M{"_id": tema.Id}, tema)
}

func (n *Temas) RemoveTema(id string) error {
	return Remove(db, "temas", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n *Temas) FindTemaByIdAreaAndAbreviacion(id_area string, abreviacion string) ([]Temas, error) {
	var result []Temas
	err := FindAll(db, "temas", bson.M{"id_area": bson.ObjectIdHex(id_area), "abreviacion": abreviacion}, nil, &result)
	return result, err
}
