package models

import "github.com/globalsign/mgo/bson"

type PreguntasCPFavoritas struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	IdUsuario  bson.ObjectId `bson:"id_usuario" json:"id_usuario"`
	IdPregunta bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
}

func (n *PreguntasCPFavoritas) InsertPreguntaCPFavorita(pregunta PreguntasCPFavoritas) error {
	return Insert(db, "preguntas_cp_favoritas", pregunta)
}

func (n *PreguntasCPFavoritas) FindAllPreguntasCPFavoritas() ([]PreguntasCPFavoritas, error) {
	var result []PreguntasCPFavoritas
	err := FindAll(db, "preguntas_cp_favoritas", nil, nil, &result)
	return result, err
}

func (n *PreguntasCPFavoritas) FindPreguntaCPFavoritaById(id string) (PreguntasCPFavoritas, error) {
	var result PreguntasCPFavoritas
	err := FindOne(db, "preguntas_cp_favoritas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *PreguntasCPFavoritas) FindAllPreguntasCPFavoritasByIdUsuario(id_usuario string) ([]PreguntasCPFavoritas, error) {
	var result []PreguntasCPFavoritas
	err := FindAll(db, "preguntas_cp_favoritas", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario)}, nil, &result)
	return result, err
}

func (n *PreguntasCPFavoritas) UpdatePreguntaCPFavorita(pregunta PreguntasCPFavoritas) error {
	return Update(db, "preguntas_cp_favoritas", bson.M{"_id": pregunta.Id}, pregunta)
}

func (n *PreguntasCPFavoritas) RemovePreguntaCPFavorita(id string) error {
	return Remove(db, "preguntas_cp_favoritas", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n *PreguntasCPFavoritas) RemovePreguntaCPFavoritaByIdUsuarioIdPregunta(id_usuario string, id_pregunta string) error {
	return Remove(db, "preguntas_cp_favoritas", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario), "id_pregunta": bson.ObjectIdHex(id_pregunta)})
}
