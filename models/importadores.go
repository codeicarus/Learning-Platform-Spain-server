package models

type Importadores struct {
	Email string `bson:"email" json:"email"`
	Name  string `bson:"string" json:"name"`
	File  string `bson:"file" json:"file"`
}
