package models

import "github.com/globalsign/mgo/bson"

type FrasesMotivadoras struct {
	Id    bson.ObjectId `bson:"_id" json:"id"`
	Frase string        `bson:"frase" json:"frase"`
	Link  string        `bson:"link" json:"link"`
	Turno string        `bson:"turno" json:"turno"`
}

func (n *FrasesMotivadoras) InsertFraseMotivadora(frase_motivadora FrasesMotivadoras) error {
	return Insert(db, "frases_motivadoras", frase_motivadora)
}
func (n *FrasesMotivadoras) FindFraseMotivadoraByTurno(turnos []string) ([]FrasesMotivadoras, error) {
	var result []FrasesMotivadoras
	err := FindAll(db, "frases_motivadoras", bson.M{"turno": bson.M{"$in": turnos}}, nil, &result)
	return result, err
}

func (n *FrasesMotivadoras) FindAllFrasesMotivadoras() ([]FrasesMotivadoras, error) {
	var result []FrasesMotivadoras
	err := FindAll(db, "frases_motivadoras", nil, nil, &result)
	return result, err
}

func (n *FrasesMotivadoras) FindFraseMotivadoraById(id string) (FrasesMotivadoras, error) {
	var result FrasesMotivadoras
	err := FindOne(db, "frases_motivadoras", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *FrasesMotivadoras) UpdateFraseMotivadora(frase_motivadora FrasesMotivadoras) error {
	return Update(db, "frases_motivadoras", bson.M{"_id": frase_motivadora.Id}, frase_motivadora)
}

func (n *FrasesMotivadoras) RemoveFraseMotivadora(id string) error {
	return Remove(db, "frases_motivadoras", bson.M{"_id": bson.ObjectIdHex(id)})
}
