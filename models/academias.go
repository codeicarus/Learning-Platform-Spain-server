package models

import "github.com/globalsign/mgo/bson"

type Academias struct {
	Id                   bson.ObjectId `bson:"_id" json:"id"`
	Nombre               string        `bson:"nombre" json:"nombre"`
	Codigo               string        `bson:"codigo" json:"codigo"`
	PorcentajeEstudiante float32       `bson:"porcentaje_estudiante" json:"porcentaje_estudiante"`
	PorcentajeAcademia   float32       `bson:"porcentaje_academia" json:"porcentaje_academia"`
	DuracionMeses        float32       `bson:"duracion_meses" json:"duracion_meses"`
	InicioContrato       int64         `bson:"inicio_contrato" json:"inicio_contrato"`
}

func (n *Academias) InsertAcademia(academia Academias) error {
	return Insert(db, "academias", academia)
}

func (n *Academias) FindAllAcademias() ([]Academias, error) {
	var result []Academias
	err := FindAll(db, "academias", nil, nil, &result)
	return result, err
}

func (n *Academias) FindAcademiaByCodigo(codigo string) (Academias, error) {
	var result Academias
	err := FindOne(db, "academias", bson.M{"codigo": codigo}, nil, &result)
	return result, err
}
func (n *Academias) FindAcademiaById(id string) (Academias, error) {
	var result Academias
	err := FindOne(db, "academias", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *Academias) UpdateAcademia(academia Academias) error {
	return Update(db, "academias", bson.M{"_id": academia.Id}, academia)
}

func (n *Academias) RemoveAcademia(id string) error {
	return Remove(db, "academias", bson.M{"_id": bson.ObjectIdHex(id)})
}
