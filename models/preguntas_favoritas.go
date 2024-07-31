package models

import "github.com/globalsign/mgo/bson"

type PreguntasFavoritas struct {
	Id         bson.ObjectId `bson:"_id" json:"id"`
	IdUsuario  bson.ObjectId `bson:"id_usuario" json:"id_usuario"`
	IdPregunta bson.ObjectId `bson:"id_pregunta" json:"id_pregunta"`
}

func (n *PreguntasFavoritas) InsertPreguntaFavorita(pregunta PreguntasFavoritas) error {
	return Insert(db, "preguntas_favoritas", pregunta)
}

func (n *PreguntasFavoritas) FindAllPreguntasFavoritas() ([]PreguntasFavoritas, error) {
	var result []PreguntasFavoritas
	err := FindAll(db, "preguntas_favoritas", nil, nil, &result)
	return result, err
}

func (n *PreguntasFavoritas) FindPreguntaFavoritaById(id string) (PreguntasFavoritas, error) {
	var result PreguntasFavoritas
	err := FindOne(db, "preguntas_favoritas", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *PreguntasFavoritas) FindAllPreguntasFavoritasByIdUsuario(id_usuario string) ([]PreguntasFavoritas, error) {
	var result []PreguntasFavoritas
	err := FindAll(db, "preguntas_favoritas", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario)}, nil, &result)
	return result, err
}

func (n *PreguntasFavoritas) UpdatePreguntaFavorita(pregunta PreguntasFavoritas) error {
	return Update(db, "preguntas_favoritas", bson.M{"_id": pregunta.Id}, pregunta)
}

func (n *PreguntasFavoritas) RemovePreguntaFavorita(id string) error {
	return Remove(db, "preguntas_favoritas", bson.M{"_id": bson.ObjectIdHex(id)})
}

func (n *PreguntasFavoritas) RemovePreguntaFavoritaByIdUsuarioIdPregunta(id_usuario string, id_pregunta string) error {
	return Remove(db, "preguntas_favoritas", bson.M{"id_usuario": bson.ObjectIdHex(id_usuario), "id_pregunta": bson.ObjectIdHex(id_pregunta)})
}
