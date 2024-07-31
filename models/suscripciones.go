package models

import (
	"github.com/globalsign/mgo/bson"
)

type Suscripcion struct {
	Id          bson.ObjectId   `bson:"_id" json:"id"`
	Suscription string          `bson:"suscription" json:"suscription"`
	Email       string          `bson:"email" json:"email"`
	UserID      bson.ObjectId   `bson:"user_id" json:"user_id"`
	Type        string          `bson:"type" json:"type"`
	Concepto    string          `bson:"concepto" json:"concepto"`
	Status      string          `bson:"status" json:"status"`
	Mount 			string          `bson:"mount"  json:"mount"`
}


type User struct {
	Id                   bson.ObjectId `bson:"_id" json:"id"`
	Name                 string        `bson:"name" json:"name"`
	Email                string        `bson:"email" json:"email"`
	MaximoDiaSuscripcion int64         `bson:"maximo_dia_suscripcion" json:"maximo_dia_suscripcion"`
}

type SuscripcionList struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Suscription string          `bson:"suscription" json:"suscription"`
	Email       string          `bson:"email" json:"email"`
	UserID      bson.ObjectId   `bson:"user_id" json:"user_id"`
	Type        string          `bson:"type" json:"type"`
	Concepto    string          `bson:"concepto" json:"concepto"`
	Status      string          `bson:"status" json:"status"`
	Mount 			string          `bson:"mount"  json:"mount"`
	User        User            `bson:"user" json:"user"`	
}

type SuscripcionRechazar struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	UserID      bson.ObjectId   `bson:"user_id" json:"user_id"`
	Message    string   `bson:"message"   json:"message" `
}

func (n *Suscripcion) InsertSuscripcion(suscripcion Suscripcion) error {
	return Insert(db, "suscripcion", suscripcion)
}

func (n *Suscripcion) FindSuscripcionByID(id string) (Suscripcion, error) {
	var result Suscripcion
	err := FindOne(db, "suscripcion", bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}


func (n *Suscripcion) FindSuscripcionByIdUser(id string) (Suscripcion, error) {
		var result Suscripcion
		err := FindOne(db, "suscripcion", bson.M{"user_id": bson.ObjectIdHex(id), "status": bson.M{"$in": []string{"PENDIENTE", "RECHAZADO"}}}, nil, &result)
		return result, err
}

func (n *Suscripcion) FindAllSuscripcion() ([]Suscripcion, error) {
	var result []Suscripcion
	err := FindAll(db, "suscripcion", bson.M{"status": bson.M{"$in": []string{"PENDIENTE", "RECHAZADO"}}}, nil, &result)
	// err := FindAll(db, "suscripcion", bson.M{"status": "PENDIENTE"}, nil, &result)
	return result, err
}

func (n *Suscripcion) UpdateSuscripcion(suscripcion Suscripcion) error {
	return Update(db, "suscripcion", bson.M{"_id": suscripcion.Id}, suscripcion)
}

func (n *Suscripcion) DeleteSuscripcion(id string) error {
	return Remove(db, "suscripcion", bson.M{"_id": bson.ObjectIdHex(id)})
}