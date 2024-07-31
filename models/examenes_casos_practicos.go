package models

import (
	"github.com/globalsign/mgo/bson"
)

type ExamenesCasosPracticos struct {
	Id             bson.ObjectId   `bson:"_id" json:"id"`
	Name           string          `bson:"name" json:"name"`
	Oficial        bool            `bson:"oficial" json:"oficial"`
	Texto          string          `bson:"texto" json:"texto"`
	CasosPracticos []bson.ObjectId `bson:"casos_practicos" json:"casos_practicos"`
}

func (n *ExamenesCasosPracticos) InsertExamenCasoPractico(caso_practico ExamenesCasosPracticos) error {
	return Insert(db, "examenes_casos_practicos", caso_practico)
}

func (n *ExamenesCasosPracticos) FindAllExamenesCasosPracticos() ([]ExamenesCasosPracticos, error) {
	var result []ExamenesCasosPracticos
	err := FindAll(db, "examenes_casos_practicos", nil, nil, &result)
	return result, err
}

func (n *ExamenesCasosPracticos) FindExamenCasoPracticoById(id string) (ExamenesCasosPracticos, error) {
	var result ExamenesCasosPracticos
	err := FindOne(db, "examenes_casos_practicos", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (n *ExamenesCasosPracticos) FindExamenCasoPracticoByName(name string) (ExamenesCasosPracticos, error) {
	var result ExamenesCasosPracticos
	err := FindOne(db, "examenes_casos_practicos", bson.M{"name": name}, nil, &result)
	return result, err
}

func (n *ExamenesCasosPracticos) UpdateExamenCasoPractico(caso_practico ExamenesCasosPracticos) error {
	return Update(db, "examenes_casos_practicos", bson.M{"_id": caso_practico.Id}, caso_practico)
}

func (n *ExamenesCasosPracticos) RemoveExamenCasoPractico(id string) error {
	return Remove(db, "examenes_casos_practicos", bson.M{"_id": bson.ObjectIdHex(id)})
}
