package models

import (
	"github.com/globalsign/mgo/bson"
)

type BasicosConfundidos struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Abreviacion string        `bson:"abreviacion" json:"abreviacion"`
	IdArea      bson.ObjectId `bson:"id_area" json:"id_area"`
}

func (n *BasicosConfundidos) InsertBasicoConfundido(basico_confundido BasicosConfundidos) error {
	return Insert(db, "basicos_confundidos", basico_confundido)
}

func (n *BasicosConfundidos) FindAllBasicosConfundidos() ([]BasicosConfundidos, error) {
	var result []BasicosConfundidos
	err := FindAll(db, "basicos_confundidos", nil, nil, &result)
	return result, err
}

func (n *BasicosConfundidos) FindBasicoConfundidoById(id string) (BasicosConfundidos, error) {
	var result BasicosConfundidos
	err := FindOne(db, "basicos_confundidos", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *BasicosConfundidos) FindBasicoConfundidoByIdArea(id_area string) ([]BasicosConfundidos, error) {
	var result []BasicosConfundidos
	err := FindAll(db, "basicos_confundidos", bson.M{"id_area": bson.ObjectIdHex(id_area)}, nil, &result)
	return result, err
}

func (n *BasicosConfundidos) FindBasicoConfundidoByAbreviacion(abreviacion string) (BasicosConfundidos, error) {
	var result BasicosConfundidos
	err := FindOne(db, "basicos_confundidos", bson.M{"abreviacion": abreviacion}, nil, &result)
	return result, err
}

func (n *BasicosConfundidos) FindBasicoConfundidoByAbreviacionAndArea(abreviacion string, id_area string) (BasicosConfundidos, error) {
	var result BasicosConfundidos
	err := FindOne(db, "basicos_confundidos", bson.M{"abreviacion": abreviacion, "id_area": id_area}, nil, &result)
	return result, err
}

func (n *BasicosConfundidos) UpdateBasicoConfundido(basico_confundido BasicosConfundidos) error {
	return Update(db, "basicos_confundidos", bson.M{"_id": basico_confundido.Id}, basico_confundido)
}

func (n *BasicosConfundidos) RemoveBasicoConfundido(id string) error {
	return Remove(db, "basicos_confundidos", bson.M{"_id": bson.ObjectIdHex(id)})
}
